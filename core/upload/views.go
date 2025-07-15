package upload

import (
	"context"
	"fmt"
	"log"

	"github.com/trueegorletov/analabit/core/database"
)

// RefreshMaterializedViews refreshes all materialized views
func RefreshMaterializedViews(ctx context.Context, dbClient *database.Client) error {
	// Refresh application_flags materialized view concurrently to avoid blocking reads
	query := "REFRESH MATERIALIZED VIEW CONCURRENTLY application_flags"
	_, err := dbClient.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to refresh application_flags materialized view: %w", err)
	}

	log.Println("Successfully refreshed materialized views")
	return nil
}
