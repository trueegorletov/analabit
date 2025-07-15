package database

import (
	"context"
	"fmt"
	"log"
)

// ViewManager handles materialized view operations
type ViewManager struct {
	client *Client
}

// NewViewManager creates a new view manager
func NewViewManager(client *Client) *ViewManager {
	return &ViewManager{client: client}
}



// CreateApplicationFlagsView creates the application_flags materialized view
func (vm *ViewManager) CreateApplicationFlagsView(ctx context.Context) error {
	query := `
CREATE MATERIALIZED VIEW IF NOT EXISTS application_flags AS
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
`
	_, err := vm.client.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to create application_flags view: %w", err)
	}
	log.Println("Successfully created application_flags materialized view")
	return nil
}

// DropApplicationFlagsView drops the application_flags materialized view
func (vm *ViewManager) DropApplicationFlagsView(ctx context.Context) error {
	query := "DROP MATERIALIZED VIEW IF EXISTS application_flags"
	_, err := vm.client.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to drop application_flags view: %w", err)
	}
	log.Println("Successfully dropped application_flags materialized view")
	return nil
}

// ViewExists checks if a materialized view exists
func (vm *ViewManager) ViewExists(ctx context.Context, viewName string) (bool, error) {
	query := `
	SELECT EXISTS (
		SELECT 1 FROM pg_matviews 
		WHERE matviewname = $1
	)
	`
	
	rows, err := vm.client.QueryRowContext(ctx, query, viewName)
	if err != nil {
		return false, fmt.Errorf("failed to check view existence: %w", err)
	}
	defer rows.Close()
	
	var exists bool
	if rows.Next() {
		err = rows.Scan(&exists)
		if err != nil {
			return false, fmt.Errorf("failed to scan view existence: %w", err)
		}
	}
	return exists, nil
}

// GetViewSize returns the size of a materialized view in bytes
func (vm *ViewManager) GetViewSize(ctx context.Context, viewName string) (int64, error) {
	query := `
	SELECT pg_total_relation_size(c.oid)
	FROM pg_class c
	JOIN pg_namespace n ON n.oid = c.relnamespace
	WHERE c.relname = $1 AND n.nspname = 'public'
	`
	
	rows, err := vm.client.QueryRowContext(ctx, query, viewName)
	if err != nil {
		return 0, fmt.Errorf("failed to get view size: %w", err)
	}
	defer rows.Close()
	
	var size int64
	if rows.Next() {
		err = rows.Scan(&size)
		if err != nil {
			return 0, fmt.Errorf("failed to scan view size: %w", err)
		}
	}
	return size, nil
}

// GetViewRowCount returns the number of rows in a materialized view
func (vm *ViewManager) GetViewRowCount(ctx context.Context, viewName string) (int64, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", viewName)
	
	rows, err := vm.client.QueryRowContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("failed to get view row count: %w", err)
	}
	defer rows.Close()
	
	var count int64
	if rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return 0, fmt.Errorf("failed to scan view row count: %w", err)
		}
	}
	return count, nil
}