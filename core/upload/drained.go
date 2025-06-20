package upload

import (
	"analabit/core"
	"analabit/core/drainer"
	"analabit/core/ent"
	"analabit/core/utils"
	"context"
	"fmt"
)

const (
	drainedResultsLockID = 3
)

func DrainedResults(ctx context.Context, client *ent.Client, origin *core.VarsityCalculator, results []drainer.DrainedResult) error {
	h := &helper{
		origin: origin,
		client: client,
	}

	return h.doUploadDrained(ctx, results)
}

func (u *helper) doUploadDrained(ctx context.Context, results []drainer.DrainedResult) (err error) {
	return utils.WithTx(ctx, u.client, func(tx *ent.Tx) error {
		if err := lock(ctx, tx, drainedResultsLockID); err != nil {
			return err
		}
		defer unlock(ctx, tx, drainedResultsLockID)

		txu := &helper{
			origin: u.origin,
			client: tx.Client(),
		}

		return txu.uploadDrained(ctx, results)
	})
}

func (u *helper) uploadDrained(ctx context.Context, results []drainer.DrainedResult) error {
	var v []struct {
		Max int `json:"max"`
	}
	if err := u.client.DrainedResult.Query().Aggregate(ent.Max("iteration")).Scan(ctx, &v); err != nil {
		return fmt.Errorf("failed to get max drained result iteration: %w", err)
	}
	nextIteration := v[0].Max + 1

	for _, result := range results {
		h, err := u.heading(ctx, result.Heading)

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
