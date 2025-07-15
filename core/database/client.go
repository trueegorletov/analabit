package database

import (
	"compress/gzip"
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/trueegorletov/analabit/core/ent"
	"github.com/trueegorletov/analabit/core/ent/run"
	"github.com/trueegorletov/analabit/core/metrics"
)

// BackupError represents an error that occurred during backup but should not stop cleanup
type BackupError struct {
	Message string
	Cause   error
}

func (e *BackupError) Error() string {
	return fmt.Sprintf("backup failed: %s", e.Message)
}

func (e *BackupError) Unwrap() error {
	return e.Cause
}

// IsBackupError checks if an error is a backup error
func IsBackupError(err error) bool {
	_, ok := err.(*BackupError)
	return ok
}

// Client wraps the Ent client with additional database functionality
type Client struct {
	*ent.Client
	builder squirrel.StatementBuilderType
}

// NewClient creates a new database client
func NewClient(entClient *ent.Client) (*Client, error) {
	return &Client{
		Client:  entClient,
		builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}, nil
}

// ExecContext executes a query without returning any rows
func (c *Client) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return c.Client.ExecContext(ctx, query, args...)
}

// QueryRowContext executes a query that is expected to return at most one row
func (c *Client) QueryRowContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return c.Client.QueryContext(ctx, query, args...)
}

// Builder returns the Squirrel query builder
func (c *Client) Builder() squirrel.StatementBuilderType {
	return c.builder
}

// ExecQuery executes a query built with Squirrel with metrics
func (c *Client) ExecQuery(ctx context.Context, query squirrel.Sqlizer) (sql.Result, error) {
	start := time.Now()
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	result, err := c.Client.ExecContext(ctx, sql, args...)
	if metrics.DBMetrics != nil {
		metrics.DBMetrics.RecordQuery("exec", "unknown", time.Since(start), err)
	}
	return result, err
}

// QueryRows executes a query and returns rows with metrics
func (c *Client) QueryRows(ctx context.Context, query squirrel.Sqlizer) (*sql.Rows, error) {
	start := time.Now()
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := c.Client.QueryContext(ctx, sql, args...)
	if metrics.DBMetrics != nil {
		metrics.DBMetrics.RecordQuery("select", "unknown", time.Since(start), err)
	}
	return rows, err
}

// errorRow implements squirrel.RowScanner for error cases
type errorRow struct {
	err error
}

func (e *errorRow) Scan(dest ...interface{}) error {
	return e.err
}

// BackupDataToBeDeleted creates a backup of data that will be deleted during cleanup
func (c *Client) BackupDataToBeDeleted(ctx context.Context, backupDir string, thresholdRunID int) error {
	// Ensure backup directory exists
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return &BackupError{
			Message: fmt.Sprintf("failed to create backup directory %s", backupDir),
			Cause:   err,
		}
	}

	// Tables that contain run_id and will have data deleted
	tables := []string{"applications", "calculations", "drained_results"}
	backupFile := fmt.Sprintf("%s/cleanup_backup_%d.csv.gz", backupDir, time.Now().Unix())

	file, err := os.Create(backupFile)
	if err != nil {
		return &BackupError{
			Message: fmt.Sprintf("failed to create backup file %s", backupFile),
			Cause:   err,
		}
	}
	defer file.Close()

	gzipWriter := gzip.NewWriter(file)
	defer gzipWriter.Close()

	csvWriter := csv.NewWriter(gzipWriter)
	defer csvWriter.Flush()

	// Write header indicating table structure
	header := []string{"table_name", "data"}
	if err := csvWriter.Write(header); err != nil {
		return &BackupError{
			Message: "failed to write CSV header",
			Cause:   err,
		}
	}

	for _, tableName := range tables {
		// Query only data that will be deleted (run_id < thresholdRunID)
		query := fmt.Sprintf("SELECT * FROM %s WHERE run_id < %d", tableName, thresholdRunID)
		dataRows, err := c.Client.QueryContext(ctx, query)
		if err != nil {
			return &BackupError{
				Message: fmt.Sprintf("failed to query table %s for backup", tableName),
				Cause:   err,
			}
		}

		columns, err := dataRows.Columns()
		if err != nil {
			dataRows.Close()
			return &BackupError{
				Message: fmt.Sprintf("failed to get columns for table %s", tableName),
				Cause:   err,
			}
		}

		// Write table header row
		tableHeader := append([]string{tableName + "_header"}, columns...)
		if err := csvWriter.Write(tableHeader); err != nil {
			dataRows.Close()
			return &BackupError{
				Message: fmt.Sprintf("failed to write table header for %s", tableName),
				Cause:   err,
			}
		}

		// Write data rows
		for dataRows.Next() {
			values := make([]interface{}, len(columns))
			valuePtrs := make([]interface{}, len(columns))
			for i := range columns {
				valuePtrs[i] = &values[i]
			}

			if err := dataRows.Scan(valuePtrs...); err != nil {
				dataRows.Close()
				return &BackupError{
					Message: fmt.Sprintf("failed to scan row for table %s", tableName),
					Cause:   err,
				}
			}

			// Convert values to strings for CSV
			row := make([]string, len(columns)+1)
			row[0] = tableName
			for i, val := range values {
				if val == nil {
					row[i+1] = ""
				} else if b, ok := val.([]byte); ok {
					row[i+1] = string(b)
				} else {
					row[i+1] = fmt.Sprintf("%v", val)
				}
			}

			if err := csvWriter.Write(row); err != nil {
				dataRows.Close()
				return &BackupError{
					Message: fmt.Sprintf("failed to write data row for table %s", tableName),
					Cause:   err,
				}
			}
		}
		dataRows.Close()
	}

	return nil
}

// PerformBackupAndCleanup performs database backup and cleans up old runs
// Returns a BackupError if only backup fails, allowing cleanup to continue
func (c *Client) PerformBackupAndCleanup(ctx context.Context, retention int, backupDir string) error {
	// Find the most recent finished run
	latestRun, err := c.Client.Run.Query().
		Where(run.Finished(true)).
		Order(ent.Desc(run.FieldID)).
		First(ctx)
	if err != nil {
		return fmt.Errorf("failed to find latest finished run: %w", err)
	}

	thresholdRunID := latestRun.ID - retention
	if thresholdRunID < 1 {
		thresholdRunID = 1
	}

	// Attempt backup - if it fails, continue with cleanup but return backup error
	var backupErr error
	if err := c.BackupDataToBeDeleted(ctx, backupDir, thresholdRunID); err != nil {
		backupErr = err
	}

	// Cleanup old data regardless of backup success
	tables := []string{"applications", "calculations", "drained_results"}
	for _, table := range tables {
		query := fmt.Sprintf("DELETE FROM %s WHERE run_id < %d", table, thresholdRunID)
		_, err := c.Client.ExecContext(ctx, query)
		if err != nil {
			return fmt.Errorf("failed to cleanup table %s: %w", table, err)
		}
	}

	// Return backup error if it occurred (cleanup was successful)
	return backupErr
}
