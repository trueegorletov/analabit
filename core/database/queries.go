package database

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
)

// ApplicationFlags represents the structure of application flags data
type ApplicationFlags struct {
	ApplicationID         int  `db:"application_id"`
	RunID                 int  `db:"run_id"`
	StudentID             string `db:"student_id"`
	Priority              int  `db:"priority"`
	OriginalSubmitted     bool `db:"original_submitted"`
	HeadingID             int  `db:"heading_id"`
	PassingToMorePriority bool `db:"passing_to_more_priority"`
	PassingNow            bool `db:"passing_now"`
	OriginalQuit          bool `db:"original_quit"`
	AnotherVarsitiesCount int  `db:"another_varsities_count"`
}

// GetApplicationFlags retrieves application flags for given application IDs
func (c *Client) GetApplicationFlags(ctx context.Context, applicationIDs []int) (map[int]ApplicationFlags, error) {
	if len(applicationIDs) == 0 {
		return make(map[int]ApplicationFlags), nil
	}

	// Convert int slice to interface slice for Squirrel
	ids := make([]interface{}, len(applicationIDs))
	for i, id := range applicationIDs {
		ids[i] = id
	}

	query := c.builder.Select(
		"application_id",
		"run_id",
		"student_id",
		"priority",
		"original_submitted",
		"heading_id",
		"passing_to_more_priority",
		"passing_now",
		"original_quit",
		"another_varsities_count",
	).From("application_flags").Where(squirrel.Eq{"application_id": ids})

	rows, err := c.QueryRows(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query application flags: %w", err)
	}
	defer rows.Close()

	result := make(map[int]ApplicationFlags)
	for rows.Next() {
		var flags ApplicationFlags
		err := rows.Scan(
			&flags.ApplicationID,
			&flags.RunID,
			&flags.StudentID,
			&flags.Priority,
			&flags.OriginalSubmitted,
			&flags.HeadingID,
			&flags.PassingToMorePriority,
			&flags.PassingNow,
			&flags.OriginalQuit,
			&flags.AnotherVarsitiesCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan application flags: %w", err)
		}
		result[flags.ApplicationID] = flags
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating application flags: %w", err)
	}

	return result, nil
}

// GetApplicationFlagsByStudentID retrieves application flags for a specific student
func (c *Client) GetApplicationFlagsByStudentID(ctx context.Context, studentID string, runID int) ([]ApplicationFlags, error) {
	query := c.builder.Select(
		"application_id",
		"run_id",
		"student_id",
		"priority",
		"original_submitted",
		"heading_id",
		"passing_to_more_priority",
		"passing_now",
		"original_quit",
		"another_varsities_count",
	).From("application_flags").Where(squirrel.Eq{
		"student_id": studentID,
		"run_id":     runID,
	})

	rows, err := c.QueryRows(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query application flags by student: %w", err)
	}
	defer rows.Close()

	var result []ApplicationFlags
	for rows.Next() {
		var flags ApplicationFlags
		err := rows.Scan(
			&flags.ApplicationID,
			&flags.RunID,
			&flags.StudentID,
			&flags.Priority,
			&flags.OriginalSubmitted,
			&flags.HeadingID,
			&flags.PassingToMorePriority,
			&flags.PassingNow,
			&flags.OriginalQuit,
			&flags.AnotherVarsitiesCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan application flags: %w", err)
		}
		result = append(result, flags)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating application flags: %w", err)
	}

	return result, nil
}