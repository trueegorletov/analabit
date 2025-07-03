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

const (
	applicationsLockID = 1
	calculationsLockID = 2
)

func Primary(ctx context.Context, client *ent.Client, origin *core.VarsityCalculator, results []core.CalculationResult) error {
	h := &helper{
		origin: origin,
		client: client,
	}

	return h.doUploadPrimary(ctx, results)
}

type helper struct {
	origin *core.VarsityCalculator
	client *ent.Client
}

func (u *helper) doUploadPrimary(ctx context.Context, calculations []core.CalculationResult) (err error) {
	return utils.WithTx(ctx, u.client, func(tx *ent.Tx) error {
		if err := lock(ctx, tx, applicationsLockID); err != nil {
			return err
		}
		defer unlock(ctx, tx, applicationsLockID)

		if err := lock(ctx, tx, calculationsLockID); err != nil {
			return err
		}
		defer unlock(ctx, tx, calculationsLockID)

		txu := &helper{
			origin: u.origin,
			client: tx.Client(),
		}

		if err := txu.uploadApplications(ctx); err != nil {
			return fmt.Errorf("failed to uploadPrimary applications: %w", err)
		}

		if err := txu.uploadCalculations(ctx, calculations); err != nil {
			return fmt.Errorf("failed to uploadPrimary calculations: %w", err)
		}

		return nil
	})
}

func (u *helper) uploadApplications(ctx context.Context) error {
	var v []struct {
		Max int `json:"max"`
	}
	if err := u.client.Application.Query().Aggregate(ent.Max("iteration")).Scan(ctx, &v); err != nil {
		return fmt.Errorf("failed to get max application iteration: %w", err)
	}
	nextIteration := v[0].Max + 1

	for _, student := range u.origin.Students() {
		applications := student.Applications()

		for i := range applications {
			a := &applications[i]

			h, err := u.heading(ctx, a.Heading())

			if err != nil {
				return err
			}

			err = u.client.Application.Create().
				SetStudentID(a.StudentID()).
				SetPriority(a.Priority()).
				SetCompetitionType(a.CompetitionType()).
				SetRatingPlace(a.RatingPlace()).
				SetScore(a.Score()).
				SetOriginalSubmitted(func() bool {
					st := u.origin.GetStudent(a.StudentID())
					if st != nil {
						return st.OriginalSubmitted()
					}
					return false
				}()).
				SetIteration(nextIteration).
				SetHeading(h).
				Exec(ctx)

			if err != nil {
				return fmt.Errorf("failed to create application for student %s: %w", a.StudentID(), err)
			}
		}
	}

	return nil
}

func (u *helper) uploadCalculations(ctx context.Context, calculations []core.CalculationResult) error {
	var v []struct {
		Max int `json:"max"`
	}
	if err := u.client.Calculation.Query().Aggregate(ent.Max("iteration")).Scan(ctx, &v); err != nil {
		return fmt.Errorf("failed to get max calculation iteration: %w", err)
	}
	nextIteration := v[0].Max + 1

	for i := range calculations {
		result := &calculations[i]

		h, err := u.heading(ctx, result.Heading)

		if err != nil {
			return err
		}

		for j, student := range result.Admitted {
			admittedPlace := j + 1 // Places are 1-based

			err = u.client.Calculation.Create().
				SetStudentID(student.ID()).
				SetAdmittedPlace(admittedPlace).
				SetIteration(nextIteration).
				SetHeading(h).
				Exec(ctx)

			if err != nil {
				return fmt.Errorf("failed to create calculation for student %s in heading %s: %v",
					student.ID(), h.Code, err)
			}
		}

	}

	return nil
}

func (u *helper) heading(ctx context.Context, h *core.Heading) (*ent.Heading, error) {
	if h == nil {
		return nil, fmt.Errorf("heading is nil")
	}

	// Try to find the heading by its full code
	existingHeading, err := u.client.Heading.Query().Where(heading.Code(h.FullCode())).First(ctx)

	if err != nil && !ent.IsNotFound(err) {
		return nil, fmt.Errorf("failed to query heading %s: %w", h.FullCode(), err)
	}

	if existingHeading != nil {
		return existingHeading, nil
	}

	// If not found, create a new heading
	return u.createHeading(ctx, h)
}

func (u *helper) varsity(ctx context.Context, code string) (*ent.Varsity, error) {
	if code == "" {
		return nil, fmt.Errorf("varsity code is empty")
	}

	// Try to find the varsity by its code
	v, err := u.client.Varsity.Query().Where(varsity.Code(code)).First(ctx)

	if err != nil && !ent.IsNotFound(err) {
		return nil, fmt.Errorf("failed to query varsity %s: %w", code, err)
	}

	if v != nil {
		return v, nil
	}

	// If not found, create a new varsity
	return u.createVarsity(ctx, code, code)
}

func (u *helper) createHeading(ctx context.Context, h *core.Heading) (*ent.Heading, error) {
	v, err := u.client.Varsity.Query().Where(varsity.Code(h.VarsityCode())).First(ctx)

	if err != nil && !ent.IsNotFound(err) {
		return nil, fmt.Errorf("failed to query varsity %s: %w", h.VarsityCode(), err)
	}

	if v == nil {
		v, err = u.createVarsity(ctx, h.VarsityCode(), h.VarsityPrettyName())

		if err != nil {
			return nil, err
		}
	}

	c := h.Capacities()

	err = u.client.Heading.Create().
		SetCode(h.FullCode()).
		SetName(h.PrettyName()).
		SetRegularCapacity(c.Regular).
		SetTargetQuotaCapacity(c.TargetQuota).
		SetDedicatedQuotaCapacity(c.DedicatedQuota).
		SetSpecialQuotaCapacity(c.SpecialQuota).
		SetVarsity(v).
		Exec(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to create heading %s: %w", h.FullCode(), err)
	}

	save, err := u.client.Heading.Query().Where(heading.Code(h.FullCode())).Only(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve created heading %s: %w", h.FullCode(), err)
	}

	slog.Info("created new heading", "code", save.Code, "name", save.Name)
	return save, nil
}

func (u *helper) createVarsity(ctx context.Context, code, prettyName string) (*ent.Varsity, error) {
	err := u.client.Varsity.Create().
		SetCode(code).
		SetName(prettyName).
		Exec(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to create varsity %s: %w", code, err)
	}

	save, err := u.client.Varsity.Query().Where(varsity.Code(code)).Only(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve created varsity %s: %w", code, err)
	}

	return save, nil
}
