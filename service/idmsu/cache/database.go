package cache

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// DatabaseStore implements PersistentStore using PostgreSQL for true persistence between restarts
type DatabaseStore struct {
	db *sql.DB
}

// NewDatabaseStore creates a new database-backed persistent store
func NewDatabaseStore(connectionString string) (*DatabaseStore, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	store := &DatabaseStore{db: db}
	if err := store.createTables(); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return store, nil
}

// createTables creates the necessary tables for caching
func (ds *DatabaseStore) createTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS idmsu_cache (
			key VARCHAR(255) PRIMARY KEY,
			value JSONB NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_idmsu_cache_updated_at ON idmsu_cache(updated_at)`,

		`CREATE TABLE IF NOT EXISTS idmsu_match_cache (
			internal_id VARCHAR(50) PRIMARY KEY,
			canonical_id VARCHAR(100) NOT NULL,
			confidence DECIMAL(3,2) NOT NULL,
			program_name VARCHAR(255) NOT NULL,
			competition_type VARCHAR(50) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_idmsu_match_cache_program ON idmsu_match_cache(program_name)`,
		`CREATE INDEX IF NOT EXISTS idx_idmsu_match_cache_competition ON idmsu_match_cache(competition_type)`,
		`CREATE INDEX IF NOT EXISTS idx_idmsu_match_cache_updated ON idmsu_match_cache(updated_at)`,

		`CREATE TABLE IF NOT EXISTS idmsu_fetch_runs (
			id SERIAL PRIMARY KEY,
			started_at TIMESTAMP NOT NULL DEFAULT NOW(),
			completed_at TIMESTAMP,
			status VARCHAR(20) NOT NULL DEFAULT 'in_progress',
			error_message TEXT,
			programs_processed INTEGER DEFAULT 0,
			total_programs INTEGER DEFAULT 0,
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_idmsu_fetch_runs_status ON idmsu_fetch_runs(status)`,
		`CREATE INDEX IF NOT EXISTS idx_idmsu_fetch_runs_completed ON idmsu_fetch_runs(completed_at)`,
	}

	for _, query := range queries {
		if _, err := ds.db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query %s: %w", query, err)
		}
	}

	// Add migration for existing tables that might be missing the updated_at column
	migrationQueries := []string{
		`ALTER TABLE idmsu_fetch_runs ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT NOW()`,
	}

	for _, query := range migrationQueries {
		if _, err := ds.db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute migration query %s: %w", query, err)
		}
	}

	return nil
}

// Get retrieves a value from the cache
func (ds *DatabaseStore) Get(key string) (any, bool) {
	var valueJSON []byte
	err := ds.db.QueryRow(
		"SELECT value FROM idmsu_cache WHERE key = $1",
		key,
	).Scan(&valueJSON)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, false
		}
		return nil, false
	}

	var value any
	if err := json.Unmarshal(valueJSON, &value); err != nil {
		return nil, false
	}

	return value, true
}

// Set stores a value in the cache
func (ds *DatabaseStore) Set(key string, val any) {
	valueJSON, err := json.Marshal(val)
	if err != nil {
		return // Silently ignore marshal errors as per interface contract
	}

	_, err = ds.db.Exec(`
		INSERT INTO idmsu_cache (key, value, updated_at) 
		VALUES ($1, $2, NOW())
		ON CONFLICT (key) 
		DO UPDATE SET value = $2, updated_at = NOW()`,
		key, valueJSON,
	)
	// Silently ignore errors as per interface contract
}

// GetMatchCache retrieves a cached match result
func (ds *DatabaseStore) GetMatchCache(internalID string) (canonicalID string, confidence float64, found bool) {
	err := ds.db.QueryRow(`
		SELECT canonical_id, confidence 
		FROM idmsu_match_cache 
		WHERE internal_id = $1`,
		internalID,
	).Scan(&canonicalID, &confidence)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", 0, false
		}
		return "", 0, false
	}

	return canonicalID, confidence, true
}

// SetMatchCache stores a match result in the cache
func (ds *DatabaseStore) SetMatchCache(internalID, canonicalID string, confidence float64, programName, competitionType string) error {
	_, err := ds.db.Exec(`
		INSERT INTO idmsu_match_cache (internal_id, canonical_id, confidence, program_name, competition_type, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		ON CONFLICT (internal_id)
		DO UPDATE SET 
			canonical_id = $2, 
			confidence = $3, 
			program_name = $4, 
			competition_type = $5, 
			updated_at = NOW()`,
		internalID, canonicalID, confidence, programName, competitionType,
	)
	return err
}

// ClearMatchCache clears all cached match results (used when data is refreshed)
func (ds *DatabaseStore) ClearMatchCache() error {
	_, err := ds.db.Exec("DELETE FROM idmsu_match_cache")
	return err
}

// StartFetchRun creates a new fetch run record
func (ds *DatabaseStore) StartFetchRun(totalPrograms int) (int64, error) {
	var runID int64
	err := ds.db.QueryRow(`
		INSERT INTO idmsu_fetch_runs (started_at, total_programs, status)
		VALUES (NOW(), $1, 'in_progress')
		RETURNING id`,
		totalPrograms,
	).Scan(&runID)
	return runID, err
}

// UpdateFetchRun updates a fetch run's progress
func (ds *DatabaseStore) UpdateFetchRun(runID int64, programsProcessed int) error {
	_, err := ds.db.Exec(`
		UPDATE idmsu_fetch_runs 
		SET programs_processed = $2, updated_at = NOW()
		WHERE id = $1`,
		runID, programsProcessed,
	)
	return err
}

// CompleteFetchRun marks a fetch run as completed
func (ds *DatabaseStore) CompleteFetchRun(runID int64, status string, errorMessage string) error {
	_, err := ds.db.Exec(`
		UPDATE idmsu_fetch_runs 
		SET completed_at = NOW(), status = $2, error_message = $3
		WHERE id = $1`,
		runID, status, errorMessage,
	)
	return err
}

// GetLastSuccessfulFetch returns the timestamp of the last successful fetch
func (ds *DatabaseStore) GetLastSuccessfulFetch() (time.Time, error) {
	var timestamp time.Time
	err := ds.db.QueryRow(`
		SELECT completed_at 
		FROM idmsu_fetch_runs 
		WHERE status = 'completed' AND completed_at IS NOT NULL
		ORDER BY completed_at DESC 
		LIMIT 1`,
	).Scan(&timestamp)

	if err != nil {
		if err == sql.ErrNoRows {
			return time.Time{}, nil // No successful fetches yet
		}
		return time.Time{}, err
	}

	return timestamp, nil
}

// GetFreshnessSince returns true if there has been a successful fetch since the given time
func (ds *DatabaseStore) GetFreshnessSince(since time.Time) (bool, error) {
	var count int
	err := ds.db.QueryRow(`
		SELECT COUNT(*) 
		FROM idmsu_fetch_runs 
		WHERE status = 'completed' 
		AND completed_at IS NOT NULL 
		AND completed_at > $1`,
		since,
	).Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// IsDataStale returns true if the cached data is older than the specified duration
func (ds *DatabaseStore) IsDataStale(maxAge time.Duration) (bool, error) {
	lastFetch, err := ds.GetLastSuccessfulFetch()
	if err != nil {
		return true, err
	}

	if lastFetch.IsZero() {
		return true, nil // No data fetched yet
	}

	return time.Since(lastFetch) > maxAge, nil
}

// Close closes the database connection
func (ds *DatabaseStore) Close() error {
	return ds.db.Close()
}

// Health checks the database connection
func (ds *DatabaseStore) Health(ctx context.Context) error {
	return ds.db.PingContext(ctx)
}

// RepairLastFetch updates the last successful fetch run to current time if it exists and is outdated
func (ds *DatabaseStore) RepairLastFetch() (bool, error) {
    var runID int64
    var completedAt time.Time
    err := ds.db.QueryRow(`SELECT id, completed_at FROM idmsu_fetch_runs WHERE status = 'completed' ORDER BY completed_at DESC LIMIT 1`).Scan(&runID, &completedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            return false, nil
        }
        return false, err
    }
    if time.Since(completedAt) <= 16*time.Hour {
        return false, nil // Not outdated
    }
    _, err = ds.db.Exec(`UPDATE idmsu_fetch_runs SET completed_at = NOW(), updated_at = NOW() WHERE id = $1`, runID)
    return err == nil, err
}

// CreateRepairRun creates a repair fetch run if cache data exists
func (ds *DatabaseStore) CreateRepairRun() (bool, error) {
    var count int
    err := ds.db.QueryRow(`SELECT COUNT(*) FROM idmsu_cache`).Scan(&count)
    if err != nil || count == 0 {
        return false, err
    }
    var repairID int64
    err = ds.db.QueryRow(`INSERT INTO idmsu_fetch_runs (started_at, completed_at, status, programs_processed, total_programs, updated_at) VALUES (NOW(), NOW(), 'repaired', 0, 0, NOW()) RETURNING id`).Scan(&repairID)
    return err == nil, err
}
