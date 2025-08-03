package core

// UploadPayload is the contract between producer → aggregator.
type UploadPayload struct {
	VarsityCode  string                     `json:"varsity_code"`
	VarsityName  string                     `json:"varsity_name"`
	Headings     []HeadingDTO               `json:"headings"`
	Students     []StudentDTO               `json:"students"`
	Applications []ApplicationDTO           `json:"applications"`
	Calculations []CalculationResultDTO     `json:"calculations"`
	Drained      map[int][]DrainedResultDTO `json:"drained"` // key = drainedPercent
}

// StudentDTO contains only essential data for an uploader.
type StudentDTO struct {
	ID                string `json:"id"`
	OriginalSubmitted bool   `json:"original_submitted"`
}

// ApplicationDTO is a lean version of core.Application.
type ApplicationDTO struct {
	StudentID       string      `json:"student_id"`
	HeadingCode     string      `json:"heading_code"`
	Priority        int         `json:"priority"`
	CompetitionType Competition `json:"competition_type"`
	RatingPlace     int         `json:"rating_place"`
	Score           int         `json:"score"`
	MSUInternalID   *string     `json:"msu_internal_id,omitempty"`
}

// HeadingDTO carries all essential information about a heading.
type HeadingDTO struct {
	Code                   string `json:"code"`
	Name                   string `json:"name"`
	RegularCapacity        int    `json:"regular_capacity"`
	TargetQuotaCapacity    int    `json:"target_quota_capacity"`
	DedicatedQuotaCapacity int    `json:"dedicated_quota_capacity"`
	SpecialQuotaCapacity   int    `json:"special_quota_capacity"`
}

// CalculationResultDTO is a lean version of core.CalculationResult.
type CalculationResultDTO struct {
	HeadingCode             string       `json:"heading_code"`
	Admitted                []StudentDTO `json:"admitted"`
	RegularsAdmitted        bool         `json:"regulars_admitted"`
	PassingScore            int          `json:"passing_score"`
	LastAdmittedRatingPlace int          `json:"last_admitted_rating_place"`
}

// DrainedResultDTO is a lean version of drainer.DrainedResult.
// It's a flattened version, as the original's fields are already primitive.
type DrainedResultDTO struct {
	HeadingCode                string `json:"heading_code"`
	DrainedPercent             int    `json:"drained_percent"`
	AvgPassingScore            int    `json:"avg_passing_score"`
	MinPassingScore            int    `json:"min_passing_score"`
	MaxPassingScore            int    `json:"max_passing_score"`
	MedPassingScore            int    `json:"med_passing_score"`
	AvgLastAdmittedRatingPlace int    `json:"avg_last_admitted_rating_place"`
	MinLastAdmittedRatingPlace int    `json:"min_last_admitted_rating_place"`
	MaxLastAdmittedRatingPlace int    `json:"max_last_admitted_rating_place"`
	MedLastAdmittedRatingPlace int    `json:"med_last_admitted_rating_place"`
	RegularsAdmitted           bool   `json:"regulars_admitted"`
	IsVirtual                  bool   `json:"is_virtual"`
}

// NewUploadPayloadFromCalculator creates an UploadPayload from a VarsityCalculator and its results
// Note: drainedResults parameter should be map[int][]drainer.DrainedResult but we avoid the import cycle
// msuInternalIDs parameter is optional map of studentID -> msuInternalID for MSU-specific data
func NewUploadPayloadFromCalculator(vc *VarsityCalculator, results []CalculationResult, drainedDTOs map[int][]DrainedResultDTO, msuInternalIDs map[string]string) *UploadPayload {

	payload := &UploadPayload{
		VarsityCode:  vc.code,
		VarsityName:  vc.prettyName,
		Headings:     make([]HeadingDTO, 0, len(vc.Headings())),
		Students:     make([]StudentDTO, 0),
		Applications: make([]ApplicationDTO, 0),
		Calculations: make([]CalculationResultDTO, 0, len(results)),
		Drained:      drainedDTOs,
	}

	// Create a map to track students we've already added
	studentMap := make(map[string]bool)

	// Convert headings
	for _, h := range vc.Headings() {
		payload.Headings = append(payload.Headings, HeadingDTO{
			Code:                   h.FullCode(),
			Name:                   h.PrettyName(),
			RegularCapacity:        h.Capacities().Regular,
			TargetQuotaCapacity:    h.Capacities().TargetQuota,
			DedicatedQuotaCapacity: h.Capacities().DedicatedQuota,
			SpecialQuotaCapacity:   h.Capacities().SpecialQuota,
		})
	}

	// Convert students and applications
	for _, student := range vc.Students() {
		if !studentMap[student.ID()] {
			payload.Students = append(payload.Students, StudentDTO{
				ID:                student.ID(),
				OriginalSubmitted: student.OriginalSubmitted(),
			})
			studentMap[student.ID()] = true
		}

		// Convert applications for this student
		for _, app := range student.Applications() {
			// Get MSU internal ID if available
			var msuInternalID *string
			if internalID, exists := msuInternalIDs[app.StudentID()]; exists {
				msuInternalID = &internalID
			}

			payload.Applications = append(payload.Applications, ApplicationDTO{
				StudentID:       app.StudentID(),
				HeadingCode:     app.Heading().FullCode(),
				Priority:        app.Priority(),
				CompetitionType: app.CompetitionType(),
				RatingPlace:     app.RatingPlace(),
				Score:           app.Score(),
				MSUInternalID:   msuInternalID,
			})
		}
	}

	// Convert calculation results
	for _, result := range results {
		admittedStudents := make([]StudentDTO, 0, len(result.Admitted))
		for _, student := range result.Admitted {
			admittedStudents = append(admittedStudents, StudentDTO{
				ID:                student.ID(),
				OriginalSubmitted: student.OriginalSubmitted(),
			})
		}

		passingScore, err := result.PassingScore()

		if err != nil {
			passingScore = 999
		}

		larp, err := result.LastAdmittedRatingPlace()

		if err != nil {
			larp = 0
		}

		payload.Calculations = append(payload.Calculations, CalculationResultDTO{
			HeadingCode:             result.Heading.FullCode(),
			Admitted:                admittedStudents,
			RegularsAdmitted:        result.CheckRegularsAdmitted(),
			PassingScore:            passingScore,
			LastAdmittedRatingPlace: larp,
		})
	}

	return payload
}
