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

func Do(ctx context.Context, client *ent.Client, origin *core.VarsityCalculator, results []core.CalculationResult) error {
	tx := &transaction{
		origin:       origin,
		calculations: results,
		client:       client,
	}

	if err := tx.initMetadata(ctx); err != nil {
		return fmt.Errorf("failed to initialize metadata: %w", err)
	}

	if err := tx.initTx(ctx); err != nil {
		return fmt.Errorf("failed to initialize transaction: %w", err)
	}

	if err := tx.do(ctx); err != nil {
		return err
	}

	return nil
}

type transaction struct {
	origin       *core.VarsityCalculator
	calculations []core.CalculationResult
	client       *ent.Client

	metadata *ent.Metadata
	tx       *ent.Tx
}

func (u *transaction) initMetadata(ctx context.Context) error {
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
		return fmt.Errorf("uploading lock is active, cannot proceed with transaction")
	}

	u.metadata = metadata
	return nil
}

func createMetadata(ctx context.Context, client *ent.Client) (*ent.Metadata, error) {
	metadata, err := client.Metadata.Create().
		SetLastApplicationsIteration(0).
		SetLastCalculationsIteration(0).
		SetUploadingLock(false).
		Save(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to create metadata: %w", err)
	}
	slog.Info("metadata created successfully")
	return metadata, nil
}

func (u *transaction) initTx(ctx context.Context) error {
	if u.tx != nil {
		return fmt.Errorf("transaction is already initialized")
	}

	tx, err := u.client.Tx(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	u.tx = tx
	slog.Info("transaction initialized successfully")
	return nil
}

func (u *transaction) do(ctx context.Context) (err error) {
	if err := u.initMetadata(ctx); err != nil {
		return err
	}

	if err := utils.WithTx(ctx, u.client, func(tx *ent.Tx) error {
		return u.upload(ctx, tx.Client())
	}); err != nil {
		return err
	}

	slog.Info("upload transaction completed successfully")
	return nil
}

func (u *transaction) lockMetadata(ctx context.Context) (err error) {
	u.metadata, err = u.metadata.Update().SetUploadingLock(true).Save(ctx)
	return
}

func (u *transaction) unlockMetadata(ctx context.Context) (err error) {
	u.metadata, err = u.metadata.Update().SetUploadingLock(false).Save(ctx)
	return
}

func (u *transaction) upload(ctx context.Context, client *ent.Client) error {
	if err := u.lockMetadata(ctx); err != nil {
		return fmt.Errorf("failed to lock metadata: %w", err)
	}

	if err := u.uploadApplications(ctx, client); err != nil {
		return fmt.Errorf("failed to upload applications: %w", err)
	}

	if err := u.uploadCalculations(ctx, client); err != nil {
		return fmt.Errorf("failed to upload calculations: %w", err)
	}

	if err := u.unlockMetadata(ctx); err != nil {
		return fmt.Errorf("failed to unlock metadata: %w", err)
	}

	return nil
}

func (u *transaction) uploadApplications(ctx context.Context, client *ent.Client) error {
	nextIteration := u.metadata.LastApplicationsIteration + 1

	for _, student := range u.origin.GetStudents() {
		applications := student.Applications()

		for i := range applications {
			a := &applications[i]

			err := client.Application.Create().
				SetStudentID(a.StudentID()).
				SetPriority(a.Priority()).
				SetCompetitionType(a.CompetitionType()).
				SetRatingPlace(a.RatingPlace()).
				SetScore(a.Score()).
				SetIteration(nextIteration).Exec(ctx)

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

func (u *transaction) uploadCalculations(ctx context.Context, client *ent.Client) error {
	nextIteration := u.metadata.LastCalculationsIteration + 1

	for i := range u.calculations {
		result := &u.calculations[i]

		h, err := client.Heading.Query().Where(
			heading.Code(result.Heading.FullCode())).First(ctx)

		if err != nil && !ent.IsNotFound(err) {
			return fmt.Errorf("failed to query heading %s: %w", result.Heading.FullCode(), err)
		}

		if h == nil {
			h, err = u.createHeading(ctx, client, result.Heading)

			if err != nil {
				return err
			}
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

func (u *transaction) createHeading(ctx context.Context, client *ent.Client, h *core.Heading) (*ent.Heading, error) {
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

func (u *transaction) createVarsity(ctx context.Context, client *ent.Client, code, prettyName string) (*ent.Varsity, error) {
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
