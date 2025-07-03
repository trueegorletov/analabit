package upload

import (
	"analabit/core"
	"analabit/core/ent"
	"analabit/core/ent/heading"
	"analabit/core/utils"
	"context"
	"fmt"
)

const (
	drainedResultsLockID = 3
)

func DrainedResults(ctx context.Context, client *ent.Client, results []core.DrainedResultDTO) error {
	h := &helper{
		client: client,
	}

	return h.doUploadDrained(ctx, results)
}

func (u *helper) doUploadDrained(ctx context.Context, results []core.DrainedResultDTO) (err error) {
	return utils.WithTx(ctx, u.client, func(tx *ent.Tx) error {
		if err := lock(ctx, tx, drainedResultsLockID); err != nil {
			return err
		}
		defer unlock(ctx, tx, drainedResultsLockID)

		txu := &helper{
			client: tx.Client(),
		}

		return txu.uploadDrained(ctx, results)
	})
}

func (u *helper) uploadDrained(ctx context.Context, results []core.DrainedResultDTO) error {
	var v []struct {
		Max int `json:"max"`
	}
	if err := u.client.DrainedResult.Query().Aggregate(ent.Max("iteration")).Scan(ctx, &v); err != nil {
		return fmt.Errorf("failed to get max drained result iteration: %w", err)
	}
	nextIteration := v[0].Max + 1

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
			SetIteration(nextIteration).
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
