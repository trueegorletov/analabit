package migrations

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/trueegorletov/analabit/core/database"
)

// Migration represents a single database migration
type Migration struct {
	Version     int
	Description string
	Up          string
	Down        string
}

// MigrationRunner handles database migrations
type MigrationRunner struct {
	client     *database.Client
	migrations []Migration
}

// NewMigrationRunner creates a new migration runner
func NewMigrationRunner(client *database.Client) *MigrationRunner {
	return &MigrationRunner{
		client:     client,
		migrations: getAllMigrations(),
	}
}

// Run executes all pending migrations
func (mr *MigrationRunner) Run(ctx context.Context) error {
	// Create migrations table if it doesn't exist
	if err := mr.createMigrationsTable(ctx); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get current migration version
	currentVersion, err := mr.getCurrentVersion(ctx)
	if err != nil {
		return fmt.Errorf("failed to get current migration version: %w", err)
	}

	// Sort migrations by version
	sort.Slice(mr.migrations, func(i, j int) bool {
		return mr.migrations[i].Version < mr.migrations[j].Version
	})

	// Execute pending migrations
	for _, migration := range mr.migrations {
		if migration.Version <= currentVersion {
			continue
		}

		log.Printf("Running migration %d: %s", migration.Version, migration.Description)
		if err := mr.executeMigration(ctx, migration); err != nil {
			return fmt.Errorf("failed to execute migration %d: %w", migration.Version, err)
		}

		if err := mr.recordMigration(ctx, migration); err != nil {
			return fmt.Errorf("failed to record migration %d: %w", migration.Version, err)
		}
	}

	return nil
}

// createMigrationsTable creates the schema_migrations table if it doesn't exist
func (mr *MigrationRunner) createMigrationsTable(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			description TEXT NOT NULL,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err := mr.client.ExecContext(ctx, query)
	return err
}

// getCurrentVersion gets the current migration version
func (mr *MigrationRunner) getCurrentVersion(ctx context.Context) (int, error) {
	query := "SELECT COALESCE(MAX(version), 0) FROM schema_migrations"
	rows, err := mr.client.QueryRowContext(ctx, query)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	
	var version int
	if rows.Next() {
		err = rows.Scan(&version)
		if err != nil {
			return 0, err
		}
	}
	return version, nil
}

// executeMigration executes a single migration
func (mr *MigrationRunner) executeMigration(ctx context.Context, migration Migration) error {
	// Split migration into individual statements
	statements := strings.Split(migration.Up, ";")
	
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		
		_, err := mr.client.ExecContext(ctx, stmt)
		if err != nil {
			return fmt.Errorf("failed to execute statement '%s': %w", stmt, err)
		}
	}
	
	return nil
}

// recordMigration records a migration as applied
func (mr *MigrationRunner) recordMigration(ctx context.Context, migration Migration) error {
	query := "INSERT INTO schema_migrations (version, description) VALUES ($1, $2)"
	_, err := mr.client.ExecContext(ctx, query, migration.Version, migration.Description)
	return err
}

// getAllMigrations returns all available migrations
func getAllMigrations() []Migration {
	return []Migration{
		{
			Version:     4,
			Description: "Fix passing_to_more_priority to only consider applications within same varsity",
			Up: `
DROP MATERIALIZED VIEW IF EXISTS application_flags;

CREATE MATERIALIZED VIEW application_flags AS
SELECT
  a.id AS application_id,
  a.run_id,
  a.student_id,
  a.priority,
  a.original_submitted,
  a.heading_applications AS heading_id,
  EXISTS (SELECT 1 FROM applications a2
   JOIN headings h_a2 ON a2.heading_applications = h_a2.id
   JOIN headings h_current ON a.heading_applications = h_current.id
   WHERE a2.student_id = a.student_id
     AND a2.run_id = a.run_id
     AND a2.priority < a.priority
     AND a2.heading_applications != a.heading_applications
     AND h_a2.varsity_headings = h_current.varsity_headings
     AND EXISTS (SELECT 1 FROM calculations c
                 WHERE c.student_id = a2.student_id
                   AND c.heading_calculations = a2.heading_applications
                   AND c.run_id = a2.run_id)) AS passing_to_more_priority,
  EXISTS (SELECT 1 FROM calculations c
          WHERE c.student_id = a.student_id
            AND c.heading_calculations = a.heading_applications
            AND c.run_id = a.run_id) AS passing_now,
  EXISTS (SELECT 1 FROM applications a2
   JOIN headings h2 ON a2.heading_applications = h2.id
   WHERE a2.student_id = a.student_id
     AND a2.run_id = a.run_id
     AND a2.original_submitted = true
     AND h2.varsity_headings != (SELECT varsity_headings FROM headings h3 WHERE h3.id = a.heading_applications)) AS original_quit,
  (SELECT COUNT(DISTINCT h2.varsity_headings) FROM applications a2
   JOIN headings h2 ON a2.heading_applications = h2.id
   WHERE a2.student_id = a.student_id
     AND a2.run_id = a.run_id
     AND h2.varsity_headings != (SELECT varsity_headings FROM headings h3 WHERE h3.id = a.heading_applications))::int AS another_varsities_count
FROM applications a;

CREATE UNIQUE INDEX IF NOT EXISTS application_flags_pkey ON application_flags (application_id);
`,
			Down: "DROP MATERIALIZED VIEW IF EXISTS application_flags;",
		},
		{
			Version:     3,
			Description: "Ensure column order in application_flags materialized view",
			Up: `
DROP MATERIALIZED VIEW IF EXISTS application_flags;

CREATE MATERIALIZED VIEW application_flags AS
SELECT
  a.id AS application_id,
  a.run_id,
  a.student_id,
  a.priority,
  a.original_submitted,
  a.heading_applications AS heading_id,
  EXISTS (SELECT 1 FROM applications a2
   WHERE a2.student_id = a.student_id
     AND a2.run_id = a.run_id
     AND a2.priority < a.priority
     AND a2.heading_applications != a.heading_applications
     AND EXISTS (SELECT 1 FROM calculations c
                 WHERE c.student_id = a2.student_id
                   AND c.heading_calculations = a2.heading_applications
                   AND c.run_id = a2.run_id)) AS passing_to_more_priority,
  EXISTS (SELECT 1 FROM calculations c
          WHERE c.student_id = a.student_id
            AND c.heading_calculations = a.heading_applications
            AND c.run_id = a.run_id) AS passing_now,
  EXISTS (SELECT 1 FROM applications a2
   JOIN headings h2 ON a2.heading_applications = h2.id
   WHERE a2.student_id = a.student_id
     AND a2.run_id = a.run_id
     AND a2.original_submitted = true
     AND h2.varsity_headings != (SELECT varsity_headings FROM headings h3 WHERE h3.id = a.heading_applications)) AS original_quit,
  (SELECT COUNT(DISTINCT h2.varsity_headings) FROM applications a2
   JOIN headings h2 ON a2.heading_applications = h2.id
   WHERE a2.student_id = a.student_id
     AND a2.run_id = a.run_id
     AND h2.varsity_headings != (SELECT varsity_headings FROM headings h3 WHERE h3.id = a.heading_applications))::int AS another_varsities_count
FROM applications a;

CREATE UNIQUE INDEX IF NOT EXISTS application_flags_pkey ON application_flags (application_id);
`,
			Down: "DROP MATERIALIZED VIEW IF EXISTS application_flags;",
		},
		{
			Version:     1,
			Description: "Create application_flags materialized view",
			Up: `
DROP MATERIALIZED VIEW IF EXISTS application_flags;

CREATE MATERIALIZED VIEW application_flags AS
SELECT
  a.id AS application_id,
  a.run_id,
  a.student_id,
  a.priority,
  a.original_submitted,
  a.heading_applications AS heading_id,
  EXISTS (SELECT 1 FROM applications a2
   WHERE a2.student_id = a.student_id
     AND a2.run_id = a.run_id
     AND a2.priority < a.priority
     AND a2.heading_applications != a.heading_applications
     AND EXISTS (SELECT 1 FROM calculations c
                 WHERE c.student_id = a2.student_id
                   AND c.heading_calculations = a2.heading_applications
                   AND c.run_id = a2.run_id)) AS passing_to_more_priority,
  EXISTS (SELECT 1 FROM calculations c
          WHERE c.student_id = a.student_id
            AND c.heading_calculations = a.heading_applications
            AND c.run_id = a.run_id) AS passing_now,
  (SELECT COUNT(*) FROM applications a2
   JOIN headings h2 ON a2.heading_applications = h2.id
   WHERE a2.student_id = a.student_id
     AND a2.run_id = a.run_id
     AND h2.varsity_headings != (SELECT varsity_headings FROM headings h3 WHERE h3.id = a.heading_applications))::int AS another_varsities_count
FROM applications a;

CREATE UNIQUE INDEX IF NOT EXISTS application_flags_pkey ON application_flags (application_id);
`,
			Down: "DROP MATERIALIZED VIEW IF EXISTS application_flags;",
		},
		{
			Version:     2,
			Description: "Fix application_flags materialized view - add original_quit and fix another_varsities_count",
			Up: `
DROP MATERIALIZED VIEW IF EXISTS application_flags;

CREATE MATERIALIZED VIEW application_flags AS
SELECT
  a.id AS application_id,
  a.run_id,
  a.student_id,
  a.priority,
  a.original_submitted,
  a.heading_applications AS heading_id,
  EXISTS (SELECT 1 FROM applications a2
   WHERE a2.student_id = a.student_id
     AND a2.run_id = a.run_id
     AND a2.priority < a.priority
     AND a2.heading_applications != a.heading_applications
     AND EXISTS (SELECT 1 FROM calculations c
                 WHERE c.student_id = a2.student_id
                   AND c.heading_calculations = a2.heading_applications
                   AND c.run_id = a2.run_id)) AS passing_to_more_priority,
  EXISTS (SELECT 1 FROM calculations c
          WHERE c.student_id = a.student_id
            AND c.heading_calculations = a.heading_applications
            AND c.run_id = a.run_id) AS passing_now,
  EXISTS (SELECT 1 FROM applications a2
   JOIN headings h2 ON a2.heading_applications = h2.id
   JOIN headings h_current ON a.heading_applications = h_current.id
   WHERE a2.student_id = a.student_id
     AND a2.run_id = a.run_id
     AND a2.original_submitted = true
     AND h2.varsity_headings != h_current.varsity_headings) AS original_quit,
  (SELECT COUNT(DISTINCT h2.varsity_headings) FROM applications a2
   JOIN headings h2 ON a2.heading_applications = h2.id
   WHERE a2.student_id = a.student_id
     AND a2.run_id = a.run_id
     AND h2.varsity_headings != (SELECT varsity_headings FROM headings h3 WHERE h3.id = a.heading_applications))::int AS another_varsities_count
FROM applications a;

CREATE UNIQUE INDEX IF NOT EXISTS application_flags_pkey ON application_flags (application_id);
`,
			Down: "DROP MATERIALIZED VIEW IF EXISTS application_flags;",
		},
	}
}