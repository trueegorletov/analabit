package cache

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

type FlexibleID string

func (f *FlexibleID) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		*f = FlexibleID(str)
		return nil
	}

	var num float64
	if err := json.Unmarshal(data, &num); err == nil {
		*f = FlexibleID(strconv.FormatFloat(num, 'f', -1, 64))
		return nil
	}

	return fmt.Errorf("cannot unmarshal %s into FlexibleID", string(data))
}

func (f FlexibleID) String() string {
	return string(f)
}

type Cache interface {
	GetCandidates(ctx context.Context, program string) ([]GosuslugiEntry, error)
	StoreCandidates(ctx context.Context, program string, candidates []GosuslugiEntry) error
	IsCacheStale(ctx context.Context, program string, staleThreshold time.Duration) (bool, error)
	// Run-based operations
	CreateRun(ctx context.Context, metadata map[string]interface{}) (int, error)
	FinishRun(ctx context.Context, runID int) error
	GetLatestCompletedRun(ctx context.Context) (*Run, error)
	StoreCandidatesForRun(ctx context.Context, runID int, program string, candidates []GosuslugiEntry) error
	GetCandidatesFromRun(ctx context.Context, runID int, program string) ([]GosuslugiEntry, error)
	RestoreCacheFromRun(ctx context.Context, runID int) error
	EnsureTables() error
	Close() error
}

type Run struct {
	ID          int                    `json:"id"`
	TriggeredAt time.Time              `json:"triggered_at"`
	Metadata    map[string]interface{} `json:"metadata"`
	Finished    bool                   `json:"finished"`
	FinishedAt  *time.Time             `json:"finished_at"`
}

type GosuslugiEntry struct {
	IDApplication FlexibleID `json:"idApplication"`
	ProgramName   string     `json:"name"`
	ListType      string     `json:"listType"`
	Result1       *float64   `json:"result1"`
	Result2       *float64   `json:"result2"`
	Result3       *float64   `json:"result3"`
	Result4       *float64   `json:"result4"`
	Result5       *float64   `json:"result5"`
	Result6       *float64   `json:"result6"`
	Result7       *float64   `json:"result7"`
	Result8       *float64   `json:"result8"`
	SumMark       *float64   `json:"sumMark"`
	Rating        *int       `json:"rating"`
	WithoutTests  *bool      `json:"withoutTests"`
	StatusID      *int       `json:"statusId"`
}

type postgresCache struct {
	db *sql.DB
}

func NewPostgresCache() (Cache, error) {
	dbHost := os.Getenv("DATABASE_HOST")
	dbPort := os.Getenv("DATABASE_PORT")
	dbUser := os.Getenv("DATABASE_USER")
	dbPassword := os.Getenv("DATABASE_PASSWORD")
	dbName := os.Getenv("DATABASE_DBNAME")
	dbSslMode := os.Getenv("DATABASE_SSLMODE")

	if dbHost == "" {
		dbHost = "localhost"
	}
	if dbPort == "" {
		dbPort = "5432"
	}
	if dbUser == "" {
		dbUser = "postgres"
	}
	if dbPassword == "" {
		dbPassword = "postgres"
	}
	if dbName == "" {
		dbName = "analabit"
	}
	if dbSslMode == "" {
		dbSslMode = "disable"
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbUser, dbPassword, dbHost, dbPort, dbName, dbSslMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	cache := &postgresCache{db: db}
	if err := cache.EnsureTables(); err != nil {
		return nil, fmt.Errorf("failed to ensure tables: %w", err)
	}

	return cache, nil
}

func (c *postgresCache) EnsureTables() error {
	// Create cache table if it doesn't exist
	createCacheQuery := `
		CREATE TABLE IF NOT EXISTS gosuslugi_msu_cache (
			program_name TEXT PRIMARY KEY,
			data JSONB NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)
	`
	if _, err := c.db.Exec(createCacheQuery); err != nil {
		return fmt.Errorf("failed to create cache table: %w", err)
	}

	// Create index on updated_at for efficient freshness queries
	indexQuery := `CREATE INDEX IF NOT EXISTS idx_gosuslugi_msu_cache_updated_at ON gosuslugi_msu_cache(updated_at)`
	if _, err := c.db.Exec(indexQuery); err != nil {
		return fmt.Errorf("failed to create cache index: %w", err)
	}

	// Create composite index for efficient lookups
	compositeIndexQuery := `CREATE INDEX IF NOT EXISTS idx_gosuslugi_msu_cache_program_updated ON gosuslugi_msu_cache(program_name, updated_at)`
	if _, err := c.db.Exec(compositeIndexQuery); err != nil {
		return fmt.Errorf("failed to create composite cache index: %w", err)
	}

	// Create runs table if it doesn't exist
	createRunsQuery := `
		CREATE TABLE IF NOT EXISTS idmsu_runs (
			id SERIAL PRIMARY KEY,
			triggered_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			metadata JSONB,
			finished BOOLEAN DEFAULT FALSE,
			finished_at TIMESTAMP WITH TIME ZONE
		)
	`
	if _, err := c.db.Exec(createRunsQuery); err != nil {
		return fmt.Errorf("failed to create runs table: %w", err)
	}

	// Create run data table if it doesn't exist
	createRunDataQuery := `
		CREATE TABLE IF NOT EXISTS idmsu_run_data (
			run_id INTEGER REFERENCES idmsu_runs(id) ON DELETE CASCADE,
			program_name TEXT NOT NULL,
			data JSONB NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			PRIMARY KEY (run_id, program_name)
		)
	`
	if _, err := c.db.Exec(createRunDataQuery); err != nil {
		return fmt.Errorf("failed to create run data table: %w", err)
	}

	return nil
}

func (c *postgresCache) GetCandidates(ctx context.Context, program string) ([]GosuslugiEntry, error) {
	query := `
		SELECT data FROM gosuslugi_msu_cache 
		WHERE program_name = $1
		ORDER BY updated_at DESC LIMIT 1
	`

	var dataJSON []byte
	err := c.db.QueryRowContext(ctx, query, program).Scan(&dataJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			return []GosuslugiEntry{}, nil // No cached data
		}
		return nil, fmt.Errorf("failed to query cache: %w", err)
	}

	var candidates []GosuslugiEntry
	if err := json.Unmarshal(dataJSON, &candidates); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cached data: %w", err)
	}

	return candidates, nil
}

func (c *postgresCache) StoreCandidates(ctx context.Context, program string, candidates []GosuslugiEntry) error {
	dataJSON, err := json.Marshal(candidates)
	if err != nil {
		return fmt.Errorf("failed to marshal candidates: %w", err)
	}

	query := `
		INSERT INTO gosuslugi_msu_cache (program_name, data, updated_at) 
		VALUES ($1, $2, NOW())
		ON CONFLICT (program_name) DO UPDATE SET 
			data = EXCLUDED.data,
			updated_at = EXCLUDED.updated_at
	`

	_, err = c.db.ExecContext(ctx, query, program, dataJSON)
	if err != nil {
		return fmt.Errorf("failed to store candidates: %w", err)
	}

	log.Printf("Stored %d candidates for program: %s", len(candidates), program)
	return nil
}

func (c *postgresCache) Close() error {
	return c.db.Close()
}

// Run-based operations implementation
func (c *postgresCache) CreateRun(ctx context.Context, metadata map[string]interface{}) (int, error) {
	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal metadata: %w", err)
	}

	var runID int
	query := `INSERT INTO idmsu_runs (metadata) VALUES ($1) RETURNING id`
	err = c.db.QueryRowContext(ctx, query, metadataJSON).Scan(&runID)
	if err != nil {
		return 0, fmt.Errorf("failed to create run: %w", err)
	}

	return runID, nil
}

func (c *postgresCache) FinishRun(ctx context.Context, runID int) error {
	query := `UPDATE idmsu_runs SET finished = TRUE, finished_at = NOW() WHERE id = $1`
	_, err := c.db.ExecContext(ctx, query, runID)
	if err != nil {
		return fmt.Errorf("failed to finish run: %w", err)
	}
	return nil
}

func (c *postgresCache) GetLatestCompletedRun(ctx context.Context) (*Run, error) {
	query := `
		SELECT id, triggered_at, metadata, finished, finished_at
		FROM idmsu_runs
		WHERE finished = TRUE
		ORDER BY finished_at DESC
		LIMIT 1
	`

	var run Run
	var metadataJSON []byte
	err := c.db.QueryRowContext(ctx, query).Scan(
		&run.ID,
		&run.TriggeredAt,
		&metadataJSON,
		&run.Finished,
		&run.FinishedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get latest completed run: %w", err)
	}

	if len(metadataJSON) > 0 {
		if err := json.Unmarshal(metadataJSON, &run.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}
	}

	return &run, nil
}

func (c *postgresCache) StoreCandidatesForRun(ctx context.Context, runID int, program string, candidates []GosuslugiEntry) error {
	data, err := json.Marshal(candidates)
	if err != nil {
		return fmt.Errorf("failed to marshal candidates: %w", err)
	}

	query := `
		INSERT INTO idmsu_run_data (run_id, program_name, data)
		VALUES ($1, $2, $3)
		ON CONFLICT (run_id, program_name)
		DO UPDATE SET data = EXCLUDED.data, created_at = NOW()
	`

	_, err = c.db.ExecContext(ctx, query, runID, program, data)
	if err != nil {
		return fmt.Errorf("failed to store candidates for run: %w", err)
	}

	return nil
}

func (c *postgresCache) GetCandidatesFromRun(ctx context.Context, runID int, program string) ([]GosuslugiEntry, error) {
	query := `SELECT data FROM idmsu_run_data WHERE run_id = $1 AND program_name = $2`

	var data []byte
	err := c.db.QueryRowContext(ctx, query, runID, program).Scan(&data)
	if err == sql.ErrNoRows {
		return []GosuslugiEntry{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get candidates from run: %w", err)
	}

	var candidates []GosuslugiEntry
	if err := json.Unmarshal(data, &candidates); err != nil {
		return nil, fmt.Errorf("failed to unmarshal candidates: %w", err)
	}

	return candidates, nil
}

func (c *postgresCache) IsCacheStale(ctx context.Context, program string, staleThreshold time.Duration) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM gosuslugi_msu_cache
			WHERE program_name = $1 AND updated_at > NOW() - $2::interval
		)
	`

	var exists bool
	err := c.db.QueryRowContext(ctx, query, program, fmt.Sprintf("%f seconds", staleThreshold.Seconds())).Scan(&exists)
	if err != nil {
		return true, fmt.Errorf("failed to check cache freshness: %w", err)
	}

	return !exists, nil
}

func (c *postgresCache) RestoreCacheFromRun(ctx context.Context, runID int) error {
	// Clear current cache
	if _, err := c.db.ExecContext(ctx, `DELETE FROM gosuslugi_msu_cache`); err != nil {
		return fmt.Errorf("failed to clear cache: %w", err)
	}

	// Copy data from run to cache
	query := `
		INSERT INTO gosuslugi_msu_cache (program_name, data)
		SELECT program_name, data FROM idmsu_run_data WHERE run_id = $1
	`

	_, err := c.db.ExecContext(ctx, query, runID)
	if err != nil {
		return fmt.Errorf("failed to restore cache from run: %w", err)
	}

	return nil
}

// CleanupOldEntries removes cache entries older than the specified duration
func (c *postgresCache) CleanupOldEntries(ctx context.Context, maxAge time.Duration) error {
	query := `DELETE FROM gosuslugi_msu_cache WHERE updated_at < NOW() - INTERVAL '%d minutes'`
	_, err := c.db.ExecContext(ctx, fmt.Sprintf(query, int(maxAge.Minutes())))
	if err != nil {
		return fmt.Errorf("failed to cleanup old entries: %w", err)
	}
	return nil
}
