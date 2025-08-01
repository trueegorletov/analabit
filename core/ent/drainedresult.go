// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/trueegorletov/analabit/core/ent/drainedresult"
	"github.com/trueegorletov/analabit/core/ent/heading"
	"github.com/trueegorletov/analabit/core/ent/run"
)

// DrainedResult is the model entity for the DrainedResult schema.
type DrainedResult struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// DrainedPercent holds the value of the "drained_percent" field.
	DrainedPercent int `json:"drained_percent,omitempty"`
	// AvgPassingScore holds the value of the "avg_passing_score" field.
	AvgPassingScore int `json:"avg_passing_score,omitempty"`
	// MinPassingScore holds the value of the "min_passing_score" field.
	MinPassingScore int `json:"min_passing_score,omitempty"`
	// MaxPassingScore holds the value of the "max_passing_score" field.
	MaxPassingScore int `json:"max_passing_score,omitempty"`
	// MedPassingScore holds the value of the "med_passing_score" field.
	MedPassingScore int `json:"med_passing_score,omitempty"`
	// AvgLastAdmittedRatingPlace holds the value of the "avg_last_admitted_rating_place" field.
	AvgLastAdmittedRatingPlace int `json:"avg_last_admitted_rating_place,omitempty"`
	// MinLastAdmittedRatingPlace holds the value of the "min_last_admitted_rating_place" field.
	MinLastAdmittedRatingPlace int `json:"min_last_admitted_rating_place,omitempty"`
	// MaxLastAdmittedRatingPlace holds the value of the "max_last_admitted_rating_place" field.
	MaxLastAdmittedRatingPlace int `json:"max_last_admitted_rating_place,omitempty"`
	// MedLastAdmittedRatingPlace holds the value of the "med_last_admitted_rating_place" field.
	MedLastAdmittedRatingPlace int `json:"med_last_admitted_rating_place,omitempty"`
	// RunID holds the value of the "run_id" field.
	RunID int `json:"run_id,omitempty"`
	// RegularsAdmitted holds the value of the "regulars_admitted" field.
	RegularsAdmitted bool `json:"regulars_admitted,omitempty"`
	// IsVirtual holds the value of the "is_virtual" field.
	IsVirtual bool `json:"is_virtual,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the DrainedResultQuery when eager-loading is set.
	Edges                   DrainedResultEdges `json:"edges"`
	heading_drained_results *int
	selectValues            sql.SelectValues
}

// DrainedResultEdges holds the relations/edges for other nodes in the graph.
type DrainedResultEdges struct {
	// Heading holds the value of the heading edge.
	Heading *Heading `json:"heading,omitempty"`
	// Run holds the value of the run edge.
	Run *Run `json:"run,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
}

// HeadingOrErr returns the Heading value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e DrainedResultEdges) HeadingOrErr() (*Heading, error) {
	if e.Heading != nil {
		return e.Heading, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: heading.Label}
	}
	return nil, &NotLoadedError{edge: "heading"}
}

// RunOrErr returns the Run value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e DrainedResultEdges) RunOrErr() (*Run, error) {
	if e.Run != nil {
		return e.Run, nil
	} else if e.loadedTypes[1] {
		return nil, &NotFoundError{label: run.Label}
	}
	return nil, &NotLoadedError{edge: "run"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*DrainedResult) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case drainedresult.FieldRegularsAdmitted, drainedresult.FieldIsVirtual:
			values[i] = new(sql.NullBool)
		case drainedresult.FieldID, drainedresult.FieldDrainedPercent, drainedresult.FieldAvgPassingScore, drainedresult.FieldMinPassingScore, drainedresult.FieldMaxPassingScore, drainedresult.FieldMedPassingScore, drainedresult.FieldAvgLastAdmittedRatingPlace, drainedresult.FieldMinLastAdmittedRatingPlace, drainedresult.FieldMaxLastAdmittedRatingPlace, drainedresult.FieldMedLastAdmittedRatingPlace, drainedresult.FieldRunID:
			values[i] = new(sql.NullInt64)
		case drainedresult.ForeignKeys[0]: // heading_drained_results
			values[i] = new(sql.NullInt64)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the DrainedResult fields.
func (dr *DrainedResult) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case drainedresult.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			dr.ID = int(value.Int64)
		case drainedresult.FieldDrainedPercent:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field drained_percent", values[i])
			} else if value.Valid {
				dr.DrainedPercent = int(value.Int64)
			}
		case drainedresult.FieldAvgPassingScore:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field avg_passing_score", values[i])
			} else if value.Valid {
				dr.AvgPassingScore = int(value.Int64)
			}
		case drainedresult.FieldMinPassingScore:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field min_passing_score", values[i])
			} else if value.Valid {
				dr.MinPassingScore = int(value.Int64)
			}
		case drainedresult.FieldMaxPassingScore:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field max_passing_score", values[i])
			} else if value.Valid {
				dr.MaxPassingScore = int(value.Int64)
			}
		case drainedresult.FieldMedPassingScore:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field med_passing_score", values[i])
			} else if value.Valid {
				dr.MedPassingScore = int(value.Int64)
			}
		case drainedresult.FieldAvgLastAdmittedRatingPlace:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field avg_last_admitted_rating_place", values[i])
			} else if value.Valid {
				dr.AvgLastAdmittedRatingPlace = int(value.Int64)
			}
		case drainedresult.FieldMinLastAdmittedRatingPlace:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field min_last_admitted_rating_place", values[i])
			} else if value.Valid {
				dr.MinLastAdmittedRatingPlace = int(value.Int64)
			}
		case drainedresult.FieldMaxLastAdmittedRatingPlace:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field max_last_admitted_rating_place", values[i])
			} else if value.Valid {
				dr.MaxLastAdmittedRatingPlace = int(value.Int64)
			}
		case drainedresult.FieldMedLastAdmittedRatingPlace:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field med_last_admitted_rating_place", values[i])
			} else if value.Valid {
				dr.MedLastAdmittedRatingPlace = int(value.Int64)
			}
		case drainedresult.FieldRunID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field run_id", values[i])
			} else if value.Valid {
				dr.RunID = int(value.Int64)
			}
		case drainedresult.FieldRegularsAdmitted:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field regulars_admitted", values[i])
			} else if value.Valid {
				dr.RegularsAdmitted = value.Bool
			}
		case drainedresult.FieldIsVirtual:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field is_virtual", values[i])
			} else if value.Valid {
				dr.IsVirtual = value.Bool
			}
		case drainedresult.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for edge-field heading_drained_results", value)
			} else if value.Valid {
				dr.heading_drained_results = new(int)
				*dr.heading_drained_results = int(value.Int64)
			}
		default:
			dr.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the DrainedResult.
// This includes values selected through modifiers, order, etc.
func (dr *DrainedResult) Value(name string) (ent.Value, error) {
	return dr.selectValues.Get(name)
}

// QueryHeading queries the "heading" edge of the DrainedResult entity.
func (dr *DrainedResult) QueryHeading() *HeadingQuery {
	return NewDrainedResultClient(dr.config).QueryHeading(dr)
}

// QueryRun queries the "run" edge of the DrainedResult entity.
func (dr *DrainedResult) QueryRun() *RunQuery {
	return NewDrainedResultClient(dr.config).QueryRun(dr)
}

// Update returns a builder for updating this DrainedResult.
// Note that you need to call DrainedResult.Unwrap() before calling this method if this DrainedResult
// was returned from a transaction, and the transaction was committed or rolled back.
func (dr *DrainedResult) Update() *DrainedResultUpdateOne {
	return NewDrainedResultClient(dr.config).UpdateOne(dr)
}

// Unwrap unwraps the DrainedResult entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (dr *DrainedResult) Unwrap() *DrainedResult {
	_tx, ok := dr.config.driver.(*txDriver)
	if !ok {
		panic("ent: DrainedResult is not a transactional entity")
	}
	dr.config.driver = _tx.drv
	return dr
}

// String implements the fmt.Stringer.
func (dr *DrainedResult) String() string {
	var builder strings.Builder
	builder.WriteString("DrainedResult(")
	builder.WriteString(fmt.Sprintf("id=%v, ", dr.ID))
	builder.WriteString("drained_percent=")
	builder.WriteString(fmt.Sprintf("%v", dr.DrainedPercent))
	builder.WriteString(", ")
	builder.WriteString("avg_passing_score=")
	builder.WriteString(fmt.Sprintf("%v", dr.AvgPassingScore))
	builder.WriteString(", ")
	builder.WriteString("min_passing_score=")
	builder.WriteString(fmt.Sprintf("%v", dr.MinPassingScore))
	builder.WriteString(", ")
	builder.WriteString("max_passing_score=")
	builder.WriteString(fmt.Sprintf("%v", dr.MaxPassingScore))
	builder.WriteString(", ")
	builder.WriteString("med_passing_score=")
	builder.WriteString(fmt.Sprintf("%v", dr.MedPassingScore))
	builder.WriteString(", ")
	builder.WriteString("avg_last_admitted_rating_place=")
	builder.WriteString(fmt.Sprintf("%v", dr.AvgLastAdmittedRatingPlace))
	builder.WriteString(", ")
	builder.WriteString("min_last_admitted_rating_place=")
	builder.WriteString(fmt.Sprintf("%v", dr.MinLastAdmittedRatingPlace))
	builder.WriteString(", ")
	builder.WriteString("max_last_admitted_rating_place=")
	builder.WriteString(fmt.Sprintf("%v", dr.MaxLastAdmittedRatingPlace))
	builder.WriteString(", ")
	builder.WriteString("med_last_admitted_rating_place=")
	builder.WriteString(fmt.Sprintf("%v", dr.MedLastAdmittedRatingPlace))
	builder.WriteString(", ")
	builder.WriteString("run_id=")
	builder.WriteString(fmt.Sprintf("%v", dr.RunID))
	builder.WriteString(", ")
	builder.WriteString("regulars_admitted=")
	builder.WriteString(fmt.Sprintf("%v", dr.RegularsAdmitted))
	builder.WriteString(", ")
	builder.WriteString("is_virtual=")
	builder.WriteString(fmt.Sprintf("%v", dr.IsVirtual))
	builder.WriteByte(')')
	return builder.String()
}

// DrainedResults is a parsable slice of DrainedResult.
type DrainedResults []*DrainedResult
