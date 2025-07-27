package upload

import (
	"context"
	"fmt"

	"github.com/trueegorletov/analabit/core"
	"github.com/trueegorletov/analabit/core/ent"
	"github.com/trueegorletov/analabit/core/ent/heading"
)

const (
	drainedResultsLockID = 3
)

func DrainedResults(ctx context.Context, client *ent.Client, runID int, results []core.DrainedResultDTO) error {
	h := &helper{
		client: client,
		runID:  runID,
	}

	return h.doUploadDrained(ctx, results)
}

func (u *helper) doUploadDrained(ctx context.Context, results []core.DrainedResultDTO) (err error) {
	return WithTx(ctx, u.client, func(tx *ent.Tx) error {
		if err := lock(ctx, tx, drainedResultsLockID); err != nil {
			return err
		}
		defer unlock(ctx, tx, drainedResultsLockID)

		txu := &helper{
			client: tx.Client(),
			runID:  u.runID,
		}

		return txu.uploadDrained(ctx, results)
	})
}

func (u *helper) uploadDrained(ctx context.Context, results []core.DrainedResultDTO) error {
	for _, result := range results {
		h, err := u.headingByCodeSimple(ctx, result.HeadingCode)

		if err != nil {
			return err
		}

		err = u.client.DrainedResult.Create().
			SetDrainedPercent(result.DrainedPercent).
			SetAvgPassingScore(result.AvgPassingScore).
			SetMinPassingScore(result.MinPassingScore).
			SetMaxPassingScore(result.MaxPassingScore).
			SetMedPassingScore(result.MedPassingScore).
			SetAvgLastAdmittedRatingPlace(result.AvgLastAdmittedRatingPlace).
			SetMinLastAdmittedRatingPlace(result.MinLastAdmittedRatingPlace).
			SetMaxLastAdmittedRatingPlace(result.MaxLastAdmittedRatingPlace).
			SetMedLastAdmittedRatingPlace(result.MedLastAdmittedRatingPlace).
			SetIsVirtual(result.IsVirtual).
			SetRegularsAdmitted(result.RegularsAdmitted).
			SetRunID(u.runID).
			SetHeading(h).
			Exec(ctx)

		if err != nil {
			return fmt.Errorf("failed to create drained result: %w", err)
		}
	}

	return nil
}

// headingByCodeSimple finds a heading by its code, but doesn't create it if missing
// This is used for drained results where headings should already exist from primary upload
func (u *helper) headingByCodeSimple(ctx context.Context, headingCode string) (*ent.Heading, error) {
	if headingCode == "" {
		return nil, fmt.Errorf("heading code is empty")
	}

	// Try to find the heading by its code
	existingHeading, err := u.client.Heading.Query().Where(heading.Code(headingCode)).First(ctx)

	if err != nil && !ent.IsNotFound(err) {
		return nil, fmt.Errorf("failed to query heading %s: %w", headingCode, err)
	}

	if existingHeading != nil {
		return existingHeading, nil
	}

	// If not found, return error - headings should exist from primary upload
	return nil, fmt.Errorf("heading %s not found - should have been created during primary upload", headingCode)
}
