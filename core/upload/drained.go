package upload

import (
	"analabit/core"
	"analabit/core/drainer"
	"analabit/core/ent"
	"analabit/core/utils"
	"context"
	"fmt"
)

func DrainedResults(ctx context.Context, client *ent.Client, origin *core.VarsityCalculator, results []drainer.DrainedResult) error {
	tx := &helper{
		origin: origin,
		client: client,
	}

	if err := tx.initMetadata(ctx); err != nil {
		return fmt.Errorf("failed to initialize metadata: %w", err)
	}

	if err := tx.initTx(ctx); err != nil {
		return fmt.Errorf("failed to initialize helper: %w", err)
	}

	if err := tx.doUploadDrained(ctx, results); err != nil {
		return err
	}

	return nil
}

func (u *helper) doUploadDrained(ctx context.Context, results []drainer.DrainedResult) (err error) {
	if err := u.initMetadata(ctx); err != nil {
		return err
	}

	if err := utils.WithTx(ctx, u.client, func(tx *ent.Tx) error {
		return u.uploadDrained(ctx, tx.Client(), results)
	}); err != nil {
		return err
	}

	return nil
}

func (u *helper) uploadDrained(ctx context.Context, client *ent.Client, results []drainer.DrainedResult) error {
	if err := u.lockMetadata(ctx); err != nil {
		return fmt.Errorf("failed to lock metadata: %w", err)
	}

	nextIteration := u.metadata.LastDrainedResultsIteration + 1

	for _, result := range results {
		h, err := u.heading(result.Heading)

		if err != nil {
			return err
		}

		passingScore := result.AvgPassingScore
		lastAdmittedRatingPlace := result.AvgLastAdmittedRatingPlace
		drainedPercent := result.DrainedPercent

		err = client.DrainedResult.Create().
			SetDrainedPercent(drainedPercent).
			SetPassingScore(passingScore).
			SetLastAdmittedRatingPlace(lastAdmittedRatingPlace).
			SetIteration(nextIteration).
			SetHeading(h).
			Exec(ctx)

		if err != nil {
			return fmt.Errorf("failed to create drained result: %w", err)
		}
	}

	var err error
	u.metadata, err = u.metadata.Update().SetLastDrainedResultsIteration(nextIteration).Save(ctx)

	if err != nil {
		return fmt.Errorf("failed to update metadata: %w", err)
	}

	if err := u.unlockMetadata(ctx); err != nil {
		return fmt.Errorf("failed to unlock metadata: %w", err)
	}

	return nil
}
