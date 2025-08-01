// Code generated by ent, DO NOT EDIT.

package calculation

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/trueegorletov/analabit/core/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.Calculation {
	return predicate.Calculation(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.Calculation {
	return predicate.Calculation(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.Calculation {
	return predicate.Calculation(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.Calculation {
	return predicate.Calculation(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.Calculation {
	return predicate.Calculation(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.Calculation {
	return predicate.Calculation(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.Calculation {
	return predicate.Calculation(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.Calculation {
	return predicate.Calculation(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.Calculation {
	return predicate.Calculation(sql.FieldLTE(FieldID, id))
}

// StudentID applies equality check predicate on the "student_id" field. It's identical to StudentIDEQ.
func StudentID(v string) predicate.Calculation {
	return predicate.Calculation(sql.FieldEQ(FieldStudentID, v))
}

// AdmittedPlace applies equality check predicate on the "admitted_place" field. It's identical to AdmittedPlaceEQ.
func AdmittedPlace(v int) predicate.Calculation {
	return predicate.Calculation(sql.FieldEQ(FieldAdmittedPlace, v))
}

// RunID applies equality check predicate on the "run_id" field. It's identical to RunIDEQ.
func RunID(v int) predicate.Calculation {
	return predicate.Calculation(sql.FieldEQ(FieldRunID, v))
}

// UpdatedAt applies equality check predicate on the "updated_at" field. It's identical to UpdatedAtEQ.
func UpdatedAt(v time.Time) predicate.Calculation {
	return predicate.Calculation(sql.FieldEQ(FieldUpdatedAt, v))
}

// StudentIDEQ applies the EQ predicate on the "student_id" field.
func StudentIDEQ(v string) predicate.Calculation {
	return predicate.Calculation(sql.FieldEQ(FieldStudentID, v))
}

// StudentIDNEQ applies the NEQ predicate on the "student_id" field.
func StudentIDNEQ(v string) predicate.Calculation {
	return predicate.Calculation(sql.FieldNEQ(FieldStudentID, v))
}

// StudentIDIn applies the In predicate on the "student_id" field.
func StudentIDIn(vs ...string) predicate.Calculation {
	return predicate.Calculation(sql.FieldIn(FieldStudentID, vs...))
}

// StudentIDNotIn applies the NotIn predicate on the "student_id" field.
func StudentIDNotIn(vs ...string) predicate.Calculation {
	return predicate.Calculation(sql.FieldNotIn(FieldStudentID, vs...))
}

// StudentIDGT applies the GT predicate on the "student_id" field.
func StudentIDGT(v string) predicate.Calculation {
	return predicate.Calculation(sql.FieldGT(FieldStudentID, v))
}

// StudentIDGTE applies the GTE predicate on the "student_id" field.
func StudentIDGTE(v string) predicate.Calculation {
	return predicate.Calculation(sql.FieldGTE(FieldStudentID, v))
}

// StudentIDLT applies the LT predicate on the "student_id" field.
func StudentIDLT(v string) predicate.Calculation {
	return predicate.Calculation(sql.FieldLT(FieldStudentID, v))
}

// StudentIDLTE applies the LTE predicate on the "student_id" field.
func StudentIDLTE(v string) predicate.Calculation {
	return predicate.Calculation(sql.FieldLTE(FieldStudentID, v))
}

// StudentIDContains applies the Contains predicate on the "student_id" field.
func StudentIDContains(v string) predicate.Calculation {
	return predicate.Calculation(sql.FieldContains(FieldStudentID, v))
}

// StudentIDHasPrefix applies the HasPrefix predicate on the "student_id" field.
func StudentIDHasPrefix(v string) predicate.Calculation {
	return predicate.Calculation(sql.FieldHasPrefix(FieldStudentID, v))
}

// StudentIDHasSuffix applies the HasSuffix predicate on the "student_id" field.
func StudentIDHasSuffix(v string) predicate.Calculation {
	return predicate.Calculation(sql.FieldHasSuffix(FieldStudentID, v))
}

// StudentIDEqualFold applies the EqualFold predicate on the "student_id" field.
func StudentIDEqualFold(v string) predicate.Calculation {
	return predicate.Calculation(sql.FieldEqualFold(FieldStudentID, v))
}

// StudentIDContainsFold applies the ContainsFold predicate on the "student_id" field.
func StudentIDContainsFold(v string) predicate.Calculation {
	return predicate.Calculation(sql.FieldContainsFold(FieldStudentID, v))
}

// AdmittedPlaceEQ applies the EQ predicate on the "admitted_place" field.
func AdmittedPlaceEQ(v int) predicate.Calculation {
	return predicate.Calculation(sql.FieldEQ(FieldAdmittedPlace, v))
}

// AdmittedPlaceNEQ applies the NEQ predicate on the "admitted_place" field.
func AdmittedPlaceNEQ(v int) predicate.Calculation {
	return predicate.Calculation(sql.FieldNEQ(FieldAdmittedPlace, v))
}

// AdmittedPlaceIn applies the In predicate on the "admitted_place" field.
func AdmittedPlaceIn(vs ...int) predicate.Calculation {
	return predicate.Calculation(sql.FieldIn(FieldAdmittedPlace, vs...))
}

// AdmittedPlaceNotIn applies the NotIn predicate on the "admitted_place" field.
func AdmittedPlaceNotIn(vs ...int) predicate.Calculation {
	return predicate.Calculation(sql.FieldNotIn(FieldAdmittedPlace, vs...))
}

// AdmittedPlaceGT applies the GT predicate on the "admitted_place" field.
func AdmittedPlaceGT(v int) predicate.Calculation {
	return predicate.Calculation(sql.FieldGT(FieldAdmittedPlace, v))
}

// AdmittedPlaceGTE applies the GTE predicate on the "admitted_place" field.
func AdmittedPlaceGTE(v int) predicate.Calculation {
	return predicate.Calculation(sql.FieldGTE(FieldAdmittedPlace, v))
}

// AdmittedPlaceLT applies the LT predicate on the "admitted_place" field.
func AdmittedPlaceLT(v int) predicate.Calculation {
	return predicate.Calculation(sql.FieldLT(FieldAdmittedPlace, v))
}

// AdmittedPlaceLTE applies the LTE predicate on the "admitted_place" field.
func AdmittedPlaceLTE(v int) predicate.Calculation {
	return predicate.Calculation(sql.FieldLTE(FieldAdmittedPlace, v))
}

// RunIDEQ applies the EQ predicate on the "run_id" field.
func RunIDEQ(v int) predicate.Calculation {
	return predicate.Calculation(sql.FieldEQ(FieldRunID, v))
}

// RunIDNEQ applies the NEQ predicate on the "run_id" field.
func RunIDNEQ(v int) predicate.Calculation {
	return predicate.Calculation(sql.FieldNEQ(FieldRunID, v))
}

// RunIDIn applies the In predicate on the "run_id" field.
func RunIDIn(vs ...int) predicate.Calculation {
	return predicate.Calculation(sql.FieldIn(FieldRunID, vs...))
}

// RunIDNotIn applies the NotIn predicate on the "run_id" field.
func RunIDNotIn(vs ...int) predicate.Calculation {
	return predicate.Calculation(sql.FieldNotIn(FieldRunID, vs...))
}

// UpdatedAtEQ applies the EQ predicate on the "updated_at" field.
func UpdatedAtEQ(v time.Time) predicate.Calculation {
	return predicate.Calculation(sql.FieldEQ(FieldUpdatedAt, v))
}

// UpdatedAtNEQ applies the NEQ predicate on the "updated_at" field.
func UpdatedAtNEQ(v time.Time) predicate.Calculation {
	return predicate.Calculation(sql.FieldNEQ(FieldUpdatedAt, v))
}

// UpdatedAtIn applies the In predicate on the "updated_at" field.
func UpdatedAtIn(vs ...time.Time) predicate.Calculation {
	return predicate.Calculation(sql.FieldIn(FieldUpdatedAt, vs...))
}

// UpdatedAtNotIn applies the NotIn predicate on the "updated_at" field.
func UpdatedAtNotIn(vs ...time.Time) predicate.Calculation {
	return predicate.Calculation(sql.FieldNotIn(FieldUpdatedAt, vs...))
}

// UpdatedAtGT applies the GT predicate on the "updated_at" field.
func UpdatedAtGT(v time.Time) predicate.Calculation {
	return predicate.Calculation(sql.FieldGT(FieldUpdatedAt, v))
}

// UpdatedAtGTE applies the GTE predicate on the "updated_at" field.
func UpdatedAtGTE(v time.Time) predicate.Calculation {
	return predicate.Calculation(sql.FieldGTE(FieldUpdatedAt, v))
}

// UpdatedAtLT applies the LT predicate on the "updated_at" field.
func UpdatedAtLT(v time.Time) predicate.Calculation {
	return predicate.Calculation(sql.FieldLT(FieldUpdatedAt, v))
}

// UpdatedAtLTE applies the LTE predicate on the "updated_at" field.
func UpdatedAtLTE(v time.Time) predicate.Calculation {
	return predicate.Calculation(sql.FieldLTE(FieldUpdatedAt, v))
}

// HasHeading applies the HasEdge predicate on the "heading" edge.
func HasHeading() predicate.Calculation {
	return predicate.Calculation(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, HeadingTable, HeadingColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasHeadingWith applies the HasEdge predicate on the "heading" edge with a given conditions (other predicates).
func HasHeadingWith(preds ...predicate.Heading) predicate.Calculation {
	return predicate.Calculation(func(s *sql.Selector) {
		step := newHeadingStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasRun applies the HasEdge predicate on the "run" edge.
func HasRun() predicate.Calculation {
	return predicate.Calculation(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, RunTable, RunColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasRunWith applies the HasEdge predicate on the "run" edge with a given conditions (other predicates).
func HasRunWith(preds ...predicate.Run) predicate.Calculation {
	return predicate.Calculation(func(s *sql.Selector) {
		step := newRunStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Calculation) predicate.Calculation {
	return predicate.Calculation(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Calculation) predicate.Calculation {
	return predicate.Calculation(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Calculation) predicate.Calculation {
	return predicate.Calculation(sql.NotPredicates(p))
}
