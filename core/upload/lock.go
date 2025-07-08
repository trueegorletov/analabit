package upload

import (
	"github.com/trueegorletov/analabit/core/ent"
	"context"
	"fmt"
	"log/slog"
	"time"
)

const lockTimeout = 5 * time.Minute

func lock(ctx context.Context, tx *ent.Tx, lockID int64) error {
	slog.Info("waiting to acquire advisory lock", "lockID", lockID, "timeout", lockTimeout)

	timeoutCtx, cancel := context.WithTimeout(ctx, lockTimeout)
	defer cancel()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeoutCtx.Done():
			return fmt.Errorf("timed out waiting for advisory lock %d: %w", lockID, timeoutCtx.Err())
		case <-ticker.C:
			var locked bool
			rows, err := tx.QueryContext(timeoutCtx, "SELECT pg_try_advisory_lock($1)", lockID)
			if err != nil {
				return fmt.Errorf("failed to try to acquire advisory lock %d: %w", lockID, err)
			}

			if !rows.Next() {
				rows.Close()
				return fmt.Errorf("no rows returned from pg_try_advisory_lock for lock %d", lockID)
			}

			if err := rows.Scan(&locked); err != nil {
				rows.Close()
				return fmt.Errorf("failed to scan lock status for advisory lock %d: %w", lockID, err)
			}
			rows.Close()

			if locked {
				slog.Info("acquired advisory lock", "lockID", lockID)
				return nil
			}
			slog.Info("advisory lock is busy, retrying...", "lockID", lockID)
		}
	}
}

func unlock(ctx context.Context, tx *ent.Tx, lockID int64) {
	if _, err := tx.ExecContext(ctx, "SELECT pg_advisory_unlock($1)", lockID); err != nil {
		slog.Error("failed to release advisory lock", "lockID", lockID, "err", err)
	} else {
		slog.Info("released advisory lock", "lockID", lockID)
	}
}
