package upload

import (
	"analabit/core"
	"analabit/core/ent"
	"analabit/core/ent/heading"
	"analabit/core/ent/varsity"
	"analabit/core/utils"
	"context"
	"fmt"
	"log/slog"
)

func Primary(ctx context.Context, client *ent.Client, origin *core.VarsityCalculator, results []core.CalculationResult) error {
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

	if err := tx.doUploadPrimary(ctx, results); err != nil {
		return err
	}

	return nil
}

type helper struct {
	origin *core.VarsityCalculator
	client *ent.Client

	metadata *ent.Metadata
	tx       *ent.Tx
}

func (u *helper) initMetadata(ctx context.Context) error {
	if u.metadata != nil {
		return fmt.Errorf("metadata is already initialized")
	}

	metadata, err := u.client.Metadata.Query().First(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			slog.Warn("metadata not found, creating initial one")
			if metadata, err = createMetadata(ctx, u.client); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("failed to query metadata: %w", err)
		}
	}

	if metadata.UploadingLock {
		return fmt.Errorf("uploading lock is active, cannot proceed with helper")
	}

	u.metadata = metadata
	return nil
}

func createMetadata(ctx context.Context, client *ent.Client) (*ent.Metadata, error) {
	metadata, err := client.Metadata.Create().
		SetLastApplicationsIteration(0).
		SetLastCalculationsIteration(0).
		SetLastDrainedResultsIteration(0).
		SetUploadingLock(false).
		Save(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to create metadata: %w", err)
	}
	slog.Info("metadata created successfully")
	return metadata, nil
}

func (u *helper) initTx(ctx context.Context) error {
	if u.tx != nil {
		return fmt.Errorf("helper is already initialized")
	}

	tx, err := u.client.Tx(ctx)
	if err != nil {
		return fmt.Errorf("failed to start helper: %w", err)
	}
	u.tx = tx
	slog.Info("helper initialized successfully")
	return nil
}

func (u *helper) doUploadPrimary(ctx context.Context, calculations []core.CalculationResult) (err error) {
	if err := u.initMetadata(ctx); err != nil {
		return err
	}

	if err := utils.WithTx(ctx, u.client, func(tx *ent.Tx) error {
		return u.uploadPrimary(ctx, tx.Client(), calculations)
	}); err != nil {
		return err
	}

	return nil
}

func (u *helper) lockMetadata(ctx context.Context) (err error) {
	u.metadata, err = u.metadata.Update().SetUploadingLock(true).Save(ctx)
	return
}

func (u *helper) unlockMetadata(ctx context.Context) (err error) {
	u.metadata, err = u.metadata.Update().SetUploadingLock(false).Save(ctx)
	return
}

func (u *helper) uploadPrimary(ctx context.Context, client *ent.Client, calculations []core.CalculationResult) error {
	if err := u.lockMetadata(ctx); err != nil {
		return fmt.Errorf("failed to lock metadata: %w", err)
	}

	if err := u.uploadApplications(ctx, client); err != nil {
		return fmt.Errorf("failed to uploadPrimary applications: %w", err)
	}

	if err := u.uploadCalculations(ctx, client, calculations); err != nil {
		return fmt.Errorf("failed to uploadPrimary calculations: %w", err)
	}

	if err := u.unlockMetadata(ctx); err != nil {
		return fmt.Errorf("failed to unlock metadata: %w", err)
	}

	return nil
}

func (u *helper) uploadApplications(ctx context.Context, client *ent.Client) error {
	nextIteration := u.metadata.LastApplicationsIteration + 1

	for _, student := range u.origin.Students() {
		applications := student.Applications()

		for i := range applications {
			a := &applications[i]

			h, err := u.heading(a.Heading())

			if err != nil {
				return err
			}

			err = client.Application.Create().
				SetStudentID(a.StudentID()).
				SetPriority(a.Priority()).
				SetCompetitionType(a.CompetitionType()).
				SetRatingPlace(a.RatingPlace()).
				SetScore(a.Score()).
				SetIteration(nextIteration).
				SetHeading(h).
				Exec(ctx)

			if err != nil {
				return fmt.Errorf("failed to create application for student %s: %w", a.StudentID(), err)
			}
		}
	}

	// Update metadata after uploading applications
	var err error
	u.metadata, err = u.metadata.Update().SetLastApplicationsIteration(nextIteration).Save(ctx)

	if err != nil {
		return fmt.Errorf("failed to update metadata after uploading applications: %w", err)
	}

	return nil
}

func (u *helper) uploadCalculations(ctx context.Context, client *ent.Client, calculations []core.CalculationResult) error {
	nextIteration := u.metadata.LastCalculationsIteration + 1

	for i := range calculations {
		result := &calculations[i]

		h, err := u.heading(result.Heading)

		if err != nil {
			return err
		}

		for j, student := range result.Admitted {
			admittedPlace := j + 1 // Places are 1-based

			err := client.Calculation.Create().
				SetStudentID(student.ID()).
				SetAdmittedPlace(admittedPlace).
				SetIteration(nextIteration).
				SetHeading(h)

			if err != nil {
				return fmt.Errorf("failed to create calculation for student %s in heading %s: %w",
					student.ID(), h.Code, err)
			}
		}

	}

	// Update metadata after uploading calculations
	var err error
	u.metadata, err = u.metadata.Update().SetLastCalculationsIteration(nextIteration).Save(ctx)

	if err != nil {
		return fmt.Errorf("failed to update metadata after uploading calculations: %w", err)
	}

	return nil
}

func (u *helper) heading(h *core.Heading) (*ent.Heading, error) {
	if h == nil {
		return nil, fmt.Errorf("heading is nil")
	}

	ctx := context.Background()
	client := u.client

	// Try to find the heading by its full code
	existingHeading, err := client.Heading.Query().Where(heading.Code(h.FullCode())).First(ctx)

	if err != nil && !ent.IsNotFound(err) {
		return nil, fmt.Errorf("failed to query heading %s: %w", h.FullCode(), err)
	}

	if existingHeading != nil {
		return existingHeading, nil
	}

	// If not found, create a new heading
	return u.createHeading(ctx, client, h)
}

func (u *helper) createHeading(ctx context.Context, client *ent.Client, h *core.Heading) (*ent.Heading, error) {
	v, err := client.Varsity.Query().Where(varsity.Code(h.VarsityCode())).First(ctx)

	if err != nil && !ent.IsNotFound(err) {
		return nil, fmt.Errorf("failed to query varsity %s: %w", h.VarsityCode(), err)
	}

	if v == nil {
		v, err = u.createVarsity(ctx, client, h.VarsityCode(), h.VarsityPrettyName())

		if err != nil {
			return nil, err
		}
	}

	c := h.Capacities()

	save, err := client.Heading.Create().
		SetCode(h.FullCode()).
		SetName(h.PrettyName()).
		SetRegularCapacity(c.Regular).
		SetTargetQuotaCapacity(c.TargetQuota).
		SetDedicatedQuotaCapacity(c.DedicatedQuota).
		SetSpecialQuotaCapacity(c.SpecialQuota).
		Save(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to create heading %s: %w", h.FullCode(), err)
	}

	slog.Info("created new heading", "code", save.Code, "name", save.Name)
	return save, nil
}

func (u *helper) createVarsity(ctx context.Context, client *ent.Client, code, prettyName string) (*ent.Varsity, error) {
	save, err := client.Varsity.Create().
		SetCode(code).
		SetName(prettyName).
		Save(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to create varsity %s: %w", code, err)
	}

	slog.Info("created new varsity", "code", save.Code, "name", save.Name)
	return save, nil
}
