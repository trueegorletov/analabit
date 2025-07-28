package upload

import (
	"context"
	"fmt"
	"log"

	"github.com/trueegorletov/analabit/core/database"
)

// RefreshMaterializedViews refreshes all materialized views
// Gracefully handles cases where materialized views don't exist yet
func RefreshMaterializedViews(ctx context.Context, dbClient *database.Client) error {
	// Check if application_flags materialized view exists
	checkQuery := `SELECT EXISTS (
		SELECT 1 FROM pg_matviews 
		WHERE matviewname = 'application_flags'
	)`
	
	rows, err := dbClient.QueryRowContext(ctx, checkQuery)
	if err != nil {
		return fmt.Errorf("failed to check if application_flags materialized view exists: %w", err)
	}
	defer rows.Close()
	
	var exists bool
	if rows.Next() {
		if err := rows.Scan(&exists); err != nil {
			return fmt.Errorf("failed to scan materialized view existence check: %w", err)
		}
	}
	
	if !exists {
		log.Println("application_flags materialized view does not exist yet, skipping refresh")
		return nil
	}
	
	// Refresh application_flags materialized view concurrently to avoid blocking reads
	query := "REFRESH MATERIALIZED VIEW CONCURRENTLY application_flags"
	_, err = dbClient.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to refresh application_flags materialized view: %w", err)
	}

	log.Println("Successfully refreshed materialized views")
	return nil
}
