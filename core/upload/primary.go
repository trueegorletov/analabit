package upload

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/trueegorletov/analabit/core"
	"github.com/trueegorletov/analabit/core/ent"
	"github.com/trueegorletov/analabit/core/ent/heading"
	"github.com/trueegorletov/analabit/core/ent/varsity"
)

const (
	applicationsLockID = 1
	calculationsLockID = 2
)

func Primary(ctx context.Context, client *ent.Client, runID int, payload *core.UploadPayload) error {
	// Create a map for quick heading DTO lookups
	headingsMap := make(map[string]core.HeadingDTO, len(payload.Headings))
	for _, h := range payload.Headings {
		headingsMap[h.Code] = h
	}

	h := &helper{
		client:      client,
		payload:     payload,
		headingsMap: headingsMap,
		runID:       runID,
	}

	return h.doUploadPrimary(ctx)
}

type helper struct {
	client      *ent.Client
	payload     *core.UploadPayload
	headingsMap map[string]core.HeadingDTO // Map for efficient heading DTO lookup
	runID       int
}

func (u *helper) doUploadPrimary(ctx context.Context) (err error) {
	return WithTx(ctx, u.client, func(tx *ent.Tx) error {
		if err := lock(ctx, tx, applicationsLockID); err != nil {
			return err
		}
		defer unlock(ctx, tx, applicationsLockID)

		if err := lock(ctx, tx, calculationsLockID); err != nil {
			return err
		}
		defer unlock(ctx, tx, calculationsLockID)

		txu := &helper{
			client:      tx.Client(),
			payload:     u.payload,
			headingsMap: u.headingsMap, // Pass map to transaction helper
			runID:       u.runID,
		}

		if len(u.payload.Applications) > 0 {
			if err := txu.uploadApplications(ctx, u.payload.Applications, u.payload.Students); err != nil {
				return fmt.Errorf("failed to uploadPrimary applications: %w", err)
			}
		}

		if len(u.payload.Calculations) > 0 {
			if err := txu.uploadCalculations(ctx, u.payload.Calculations); err != nil {
				return fmt.Errorf("failed to uploadPrimary calculations: %w", err)
			}
		}

		return nil
	})
}

func (u *helper) uploadApplications(ctx context.Context, applications []core.ApplicationDTO, students []core.StudentDTO) error {
	// Create a map of student ID to OriginalSubmitted for fast lookup
	studentOriginalMap := make(map[string]bool)
	for _, student := range students {
		studentOriginalMap[student.ID] = student.OriginalSubmitted
	}

	for _, app := range applications {
		h, err := u.headingByCode(ctx, app.HeadingCode)

		if err != nil {
			return err
		}

		originalSubmitted := studentOriginalMap[app.StudentID]

		err = u.client.Application.Create().
			SetStudentID(app.StudentID).
			SetPriority(app.Priority).
			SetCompetitionType(app.CompetitionType).
			SetRatingPlace(app.RatingPlace).
			SetScore(app.Score).
			SetOriginalSubmitted(originalSubmitted).
			SetRunID(u.runID).
			SetHeading(h).
			Exec(ctx)

		if err != nil {
			return fmt.Errorf("failed to create application for student %s: %w", app.StudentID, err)
		}
	}

	return nil
}

func (u *helper) uploadCalculations(ctx context.Context, calculations []core.CalculationResultDTO) error {
	for _, result := range calculations {
		h, err := u.headingByCode(ctx, result.HeadingCode)

		if err != nil {
			return err
		}

		for j, student := range result.Admitted {
			admittedPlace := j + 1 // Places are 1-based

			err = u.client.Calculation.Create().
				SetStudentID(student.ID).
				SetAdmittedPlace(admittedPlace).
				SetRunID(u.runID).
				SetHeading(h).
				Exec(ctx)

			if err != nil {
				return fmt.Errorf("failed to create calculation for student %s in heading %s: %v",
					student.ID, h.Code, err)
			}
		}

	}

	return nil
}

// headingByCode finds or creates a heading by its code (which should be in FullCode format)
func (u *helper) headingByCode(ctx context.Context, headingCode string) (*ent.Heading, error) {
	if headingCode == "" {
		return nil, fmt.Errorf("heading code is empty")
	}

	// Try to find the heading by its code
	existingHeading, err := u.client.Heading.Query().Where(heading.Code(headingCode)).First(ctx)

	if err != nil && !ent.IsNotFound(err) {
		return nil, fmt.Errorf("failed to query heading %s: %w", headingCode, err)
	}

	if existingHeading != nil {
		// Check if we have DTO for this heading and update capacities if needed
		headingDTO, ok := u.headingsMap[headingCode]
		if !ok {
			// If no DTO found, return existing heading as-is
			return existingHeading, nil
		}

		// Update capacities if they differ from DTO values
		return u.updateHeadingCapacitiesIfNeeded(ctx, existingHeading, headingDTO)
	}

	// If not found, try to create it from the payload's DTOs
	headingDTO, ok := u.headingsMap[headingCode]
	if !ok {
		// This case should ideally not happen if the payload is well-formed.
		// It means an application or calculation refers to a heading not described in the payload.
		return nil, fmt.Errorf("heading %s not found in payload DTOs", headingCode)
	}

	// Extract varsity code from heading code
	var varsityCode string
	if colonIndex := strings.LastIndex(headingCode, ":"); colonIndex > 0 {
		varsityCode = headingCode[:colonIndex]
	} else {
		varsityCode = u.payload.VarsityCode // Fallback to the main varsity code
	}

	return u.createHeadingFromDTO(ctx, headingDTO, varsityCode)
}

// createHeadingFromDTO creates a new heading in the database from a DTO.
func (u *helper) createHeadingFromDTO(ctx context.Context, dto core.HeadingDTO, varsityCode string) (*ent.Heading, error) {
	// Ensure varsity exists
	v, err := u.varsity(ctx, varsityCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get or create varsity %s for heading %s: %w", varsityCode, dto.Code, err)
	}

	// Create heading with full information from the DTO
	err = u.client.Heading.Create().
		SetCode(dto.Code).
		SetName(dto.Name).
		SetRegularCapacity(dto.RegularCapacity).
		SetTargetQuotaCapacity(dto.TargetQuotaCapacity).
		SetDedicatedQuotaCapacity(dto.DedicatedQuotaCapacity).
		SetSpecialQuotaCapacity(dto.SpecialQuotaCapacity).
		SetVarsity(v).
		Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create heading %s from DTO: %w", dto.Code, err)
	}

	created, err := u.client.Heading.Query().Where(heading.Code(dto.Code)).Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve created heading %s: %w", dto.Code, err)
	}

	slog.Info("created heading from DTO", "code", created.Code, "name", created.Name)

	// Check and update capacities if needed
	_, err = u.updateHeadingCapacitiesIfNeeded(ctx, created, dto)
	if err != nil {
		return nil, err
	}

	return created, nil
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

	// If not found, create a new varsity using payload information
	prettyName := u.payload.VarsityName
	if prettyName == "" {
		prettyName = code // fallback to code if no pretty name
	}
	return u.createVarsity(ctx, code, prettyName)
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

// updateHeadingCapacitiesIfNeeded compares the existing heading capacities with DTO values
// and updates the heading if any capacities differ
func (u *helper) updateHeadingCapacitiesIfNeeded(ctx context.Context, existingHeading *ent.Heading, dto core.HeadingDTO) (*ent.Heading, error) {
	// Check if any capacity values differ
	needsUpdate := existingHeading.RegularCapacity != dto.RegularCapacity ||
		existingHeading.TargetQuotaCapacity != dto.TargetQuotaCapacity ||
		existingHeading.DedicatedQuotaCapacity != dto.DedicatedQuotaCapacity ||
		existingHeading.SpecialQuotaCapacity != dto.SpecialQuotaCapacity

	if !needsUpdate {
		return existingHeading, nil
	}

	// Log the capacity update
	slog.Info("updating heading capacities",
		"code", existingHeading.Code,
		"name", existingHeading.Name,
		"old_regular", existingHeading.RegularCapacity,
		"new_regular", dto.RegularCapacity,
		"old_target_quota", existingHeading.TargetQuotaCapacity,
		"new_target_quota", dto.TargetQuotaCapacity,
		"old_dedicated_quota", existingHeading.DedicatedQuotaCapacity,
		"new_dedicated_quota", dto.DedicatedQuotaCapacity,
		"old_special_quota", existingHeading.SpecialQuotaCapacity,
		"new_special_quota", dto.SpecialQuotaCapacity)

	// Update the heading capacities
	err := u.client.Heading.UpdateOneID(existingHeading.ID).
		SetRegularCapacity(dto.RegularCapacity).
		SetTargetQuotaCapacity(dto.TargetQuotaCapacity).
		SetDedicatedQuotaCapacity(dto.DedicatedQuotaCapacity).
		SetSpecialQuotaCapacity(dto.SpecialQuotaCapacity).
		Exec(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to update heading capacities for %s: %w", existingHeading.Code, err)
	}

	// Fetch and return the updated heading
	updatedHeading, err := u.client.Heading.Get(ctx, existingHeading.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated heading %s: %w", existingHeading.Code, err)
	}

	return updatedHeading, nil
}
