// Code generated by ent, DO NOT EDIT.

package heading

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/trueegorletov/analabit/core/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.Heading {
	return predicate.Heading(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.Heading {
	return predicate.Heading(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.Heading {
	return predicate.Heading(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.Heading {
	return predicate.Heading(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.Heading {
	return predicate.Heading(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.Heading {
	return predicate.Heading(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.Heading {
	return predicate.Heading(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.Heading {
	return predicate.Heading(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.Heading {
	return predicate.Heading(sql.FieldLTE(FieldID, id))
}

// RegularCapacity applies equality check predicate on the "regular_capacity" field. It's identical to RegularCapacityEQ.
func RegularCapacity(v int) predicate.Heading {
	return predicate.Heading(sql.FieldEQ(FieldRegularCapacity, v))
}

// TargetQuotaCapacity applies equality check predicate on the "target_quota_capacity" field. It's identical to TargetQuotaCapacityEQ.
func TargetQuotaCapacity(v int) predicate.Heading {
	return predicate.Heading(sql.FieldEQ(FieldTargetQuotaCapacity, v))
}

// DedicatedQuotaCapacity applies equality check predicate on the "dedicated_quota_capacity" field. It's identical to DedicatedQuotaCapacityEQ.
func DedicatedQuotaCapacity(v int) predicate.Heading {
	return predicate.Heading(sql.FieldEQ(FieldDedicatedQuotaCapacity, v))
}

// SpecialQuotaCapacity applies equality check predicate on the "special_quota_capacity" field. It's identical to SpecialQuotaCapacityEQ.
func SpecialQuotaCapacity(v int) predicate.Heading {
	return predicate.Heading(sql.FieldEQ(FieldSpecialQuotaCapacity, v))
}

// Code applies equality check predicate on the "code" field. It's identical to CodeEQ.
func Code(v string) predicate.Heading {
	return predicate.Heading(sql.FieldEQ(FieldCode, v))
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.Heading {
	return predicate.Heading(sql.FieldEQ(FieldName, v))
}

// RegularCapacityEQ applies the EQ predicate on the "regular_capacity" field.
func RegularCapacityEQ(v int) predicate.Heading {
	return predicate.Heading(sql.FieldEQ(FieldRegularCapacity, v))
}

// RegularCapacityNEQ applies the NEQ predicate on the "regular_capacity" field.
func RegularCapacityNEQ(v int) predicate.Heading {
	return predicate.Heading(sql.FieldNEQ(FieldRegularCapacity, v))
}

// RegularCapacityIn applies the In predicate on the "regular_capacity" field.
func RegularCapacityIn(vs ...int) predicate.Heading {
	return predicate.Heading(sql.FieldIn(FieldRegularCapacity, vs...))
}

// RegularCapacityNotIn applies the NotIn predicate on the "regular_capacity" field.
func RegularCapacityNotIn(vs ...int) predicate.Heading {
	return predicate.Heading(sql.FieldNotIn(FieldRegularCapacity, vs...))
}

// RegularCapacityGT applies the GT predicate on the "regular_capacity" field.
func RegularCapacityGT(v int) predicate.Heading {
	return predicate.Heading(sql.FieldGT(FieldRegularCapacity, v))
}

// RegularCapacityGTE applies the GTE predicate on the "regular_capacity" field.
func RegularCapacityGTE(v int) predicate.Heading {
	return predicate.Heading(sql.FieldGTE(FieldRegularCapacity, v))
}

// RegularCapacityLT applies the LT predicate on the "regular_capacity" field.
func RegularCapacityLT(v int) predicate.Heading {
	return predicate.Heading(sql.FieldLT(FieldRegularCapacity, v))
}

// RegularCapacityLTE applies the LTE predicate on the "regular_capacity" field.
func RegularCapacityLTE(v int) predicate.Heading {
	return predicate.Heading(sql.FieldLTE(FieldRegularCapacity, v))
}

// TargetQuotaCapacityEQ applies the EQ predicate on the "target_quota_capacity" field.
func TargetQuotaCapacityEQ(v int) predicate.Heading {
	return predicate.Heading(sql.FieldEQ(FieldTargetQuotaCapacity, v))
}

// TargetQuotaCapacityNEQ applies the NEQ predicate on the "target_quota_capacity" field.
func TargetQuotaCapacityNEQ(v int) predicate.Heading {
	return predicate.Heading(sql.FieldNEQ(FieldTargetQuotaCapacity, v))
}

// TargetQuotaCapacityIn applies the In predicate on the "target_quota_capacity" field.
func TargetQuotaCapacityIn(vs ...int) predicate.Heading {
	return predicate.Heading(sql.FieldIn(FieldTargetQuotaCapacity, vs...))
}

// TargetQuotaCapacityNotIn applies the NotIn predicate on the "target_quota_capacity" field.
func TargetQuotaCapacityNotIn(vs ...int) predicate.Heading {
	return predicate.Heading(sql.FieldNotIn(FieldTargetQuotaCapacity, vs...))
}

// TargetQuotaCapacityGT applies the GT predicate on the "target_quota_capacity" field.
func TargetQuotaCapacityGT(v int) predicate.Heading {
	return predicate.Heading(sql.FieldGT(FieldTargetQuotaCapacity, v))
}

// TargetQuotaCapacityGTE applies the GTE predicate on the "target_quota_capacity" field.
func TargetQuotaCapacityGTE(v int) predicate.Heading {
	return predicate.Heading(sql.FieldGTE(FieldTargetQuotaCapacity, v))
}

// TargetQuotaCapacityLT applies the LT predicate on the "target_quota_capacity" field.
func TargetQuotaCapacityLT(v int) predicate.Heading {
	return predicate.Heading(sql.FieldLT(FieldTargetQuotaCapacity, v))
}

// TargetQuotaCapacityLTE applies the LTE predicate on the "target_quota_capacity" field.
func TargetQuotaCapacityLTE(v int) predicate.Heading {
	return predicate.Heading(sql.FieldLTE(FieldTargetQuotaCapacity, v))
}

// DedicatedQuotaCapacityEQ applies the EQ predicate on the "dedicated_quota_capacity" field.
func DedicatedQuotaCapacityEQ(v int) predicate.Heading {
	return predicate.Heading(sql.FieldEQ(FieldDedicatedQuotaCapacity, v))
}

// DedicatedQuotaCapacityNEQ applies the NEQ predicate on the "dedicated_quota_capacity" field.
func DedicatedQuotaCapacityNEQ(v int) predicate.Heading {
	return predicate.Heading(sql.FieldNEQ(FieldDedicatedQuotaCapacity, v))
}

// DedicatedQuotaCapacityIn applies the In predicate on the "dedicated_quota_capacity" field.
func DedicatedQuotaCapacityIn(vs ...int) predicate.Heading {
	return predicate.Heading(sql.FieldIn(FieldDedicatedQuotaCapacity, vs...))
}

// DedicatedQuotaCapacityNotIn applies the NotIn predicate on the "dedicated_quota_capacity" field.
func DedicatedQuotaCapacityNotIn(vs ...int) predicate.Heading {
	return predicate.Heading(sql.FieldNotIn(FieldDedicatedQuotaCapacity, vs...))
}

// DedicatedQuotaCapacityGT applies the GT predicate on the "dedicated_quota_capacity" field.
func DedicatedQuotaCapacityGT(v int) predicate.Heading {
	return predicate.Heading(sql.FieldGT(FieldDedicatedQuotaCapacity, v))
}

// DedicatedQuotaCapacityGTE applies the GTE predicate on the "dedicated_quota_capacity" field.
func DedicatedQuotaCapacityGTE(v int) predicate.Heading {
	return predicate.Heading(sql.FieldGTE(FieldDedicatedQuotaCapacity, v))
}

// DedicatedQuotaCapacityLT applies the LT predicate on the "dedicated_quota_capacity" field.
func DedicatedQuotaCapacityLT(v int) predicate.Heading {
	return predicate.Heading(sql.FieldLT(FieldDedicatedQuotaCapacity, v))
}

// DedicatedQuotaCapacityLTE applies the LTE predicate on the "dedicated_quota_capacity" field.
func DedicatedQuotaCapacityLTE(v int) predicate.Heading {
	return predicate.Heading(sql.FieldLTE(FieldDedicatedQuotaCapacity, v))
}

// SpecialQuotaCapacityEQ applies the EQ predicate on the "special_quota_capacity" field.
func SpecialQuotaCapacityEQ(v int) predicate.Heading {
	return predicate.Heading(sql.FieldEQ(FieldSpecialQuotaCapacity, v))
}

// SpecialQuotaCapacityNEQ applies the NEQ predicate on the "special_quota_capacity" field.
func SpecialQuotaCapacityNEQ(v int) predicate.Heading {
	return predicate.Heading(sql.FieldNEQ(FieldSpecialQuotaCapacity, v))
}

// SpecialQuotaCapacityIn applies the In predicate on the "special_quota_capacity" field.
func SpecialQuotaCapacityIn(vs ...int) predicate.Heading {
	return predicate.Heading(sql.FieldIn(FieldSpecialQuotaCapacity, vs...))
}

// SpecialQuotaCapacityNotIn applies the NotIn predicate on the "special_quota_capacity" field.
func SpecialQuotaCapacityNotIn(vs ...int) predicate.Heading {
	return predicate.Heading(sql.FieldNotIn(FieldSpecialQuotaCapacity, vs...))
}

// SpecialQuotaCapacityGT applies the GT predicate on the "special_quota_capacity" field.
func SpecialQuotaCapacityGT(v int) predicate.Heading {
	return predicate.Heading(sql.FieldGT(FieldSpecialQuotaCapacity, v))
}

// SpecialQuotaCapacityGTE applies the GTE predicate on the "special_quota_capacity" field.
func SpecialQuotaCapacityGTE(v int) predicate.Heading {
	return predicate.Heading(sql.FieldGTE(FieldSpecialQuotaCapacity, v))
}

// SpecialQuotaCapacityLT applies the LT predicate on the "special_quota_capacity" field.
func SpecialQuotaCapacityLT(v int) predicate.Heading {
	return predicate.Heading(sql.FieldLT(FieldSpecialQuotaCapacity, v))
}

// SpecialQuotaCapacityLTE applies the LTE predicate on the "special_quota_capacity" field.
func SpecialQuotaCapacityLTE(v int) predicate.Heading {
	return predicate.Heading(sql.FieldLTE(FieldSpecialQuotaCapacity, v))
}

// CodeEQ applies the EQ predicate on the "code" field.
func CodeEQ(v string) predicate.Heading {
	return predicate.Heading(sql.FieldEQ(FieldCode, v))
}

// CodeNEQ applies the NEQ predicate on the "code" field.
func CodeNEQ(v string) predicate.Heading {
	return predicate.Heading(sql.FieldNEQ(FieldCode, v))
}

// CodeIn applies the In predicate on the "code" field.
func CodeIn(vs ...string) predicate.Heading {
	return predicate.Heading(sql.FieldIn(FieldCode, vs...))
}

// CodeNotIn applies the NotIn predicate on the "code" field.
func CodeNotIn(vs ...string) predicate.Heading {
	return predicate.Heading(sql.FieldNotIn(FieldCode, vs...))
}

// CodeGT applies the GT predicate on the "code" field.
func CodeGT(v string) predicate.Heading {
	return predicate.Heading(sql.FieldGT(FieldCode, v))
}

// CodeGTE applies the GTE predicate on the "code" field.
func CodeGTE(v string) predicate.Heading {
	return predicate.Heading(sql.FieldGTE(FieldCode, v))
}

// CodeLT applies the LT predicate on the "code" field.
func CodeLT(v string) predicate.Heading {
	return predicate.Heading(sql.FieldLT(FieldCode, v))
}

// CodeLTE applies the LTE predicate on the "code" field.
func CodeLTE(v string) predicate.Heading {
	return predicate.Heading(sql.FieldLTE(FieldCode, v))
}

// CodeContains applies the Contains predicate on the "code" field.
func CodeContains(v string) predicate.Heading {
	return predicate.Heading(sql.FieldContains(FieldCode, v))
}

// CodeHasPrefix applies the HasPrefix predicate on the "code" field.
func CodeHasPrefix(v string) predicate.Heading {
	return predicate.Heading(sql.FieldHasPrefix(FieldCode, v))
}

// CodeHasSuffix applies the HasSuffix predicate on the "code" field.
func CodeHasSuffix(v string) predicate.Heading {
	return predicate.Heading(sql.FieldHasSuffix(FieldCode, v))
}

// CodeEqualFold applies the EqualFold predicate on the "code" field.
func CodeEqualFold(v string) predicate.Heading {
	return predicate.Heading(sql.FieldEqualFold(FieldCode, v))
}

// CodeContainsFold applies the ContainsFold predicate on the "code" field.
func CodeContainsFold(v string) predicate.Heading {
	return predicate.Heading(sql.FieldContainsFold(FieldCode, v))
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.Heading {
	return predicate.Heading(sql.FieldEQ(FieldName, v))
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.Heading {
	return predicate.Heading(sql.FieldNEQ(FieldName, v))
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.Heading {
	return predicate.Heading(sql.FieldIn(FieldName, vs...))
}

// NameNotIn applies the NotIn predicate on the "name" field.
func NameNotIn(vs ...string) predicate.Heading {
	return predicate.Heading(sql.FieldNotIn(FieldName, vs...))
}

// NameGT applies the GT predicate on the "name" field.
func NameGT(v string) predicate.Heading {
	return predicate.Heading(sql.FieldGT(FieldName, v))
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.Heading {
	return predicate.Heading(sql.FieldGTE(FieldName, v))
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.Heading {
	return predicate.Heading(sql.FieldLT(FieldName, v))
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.Heading {
	return predicate.Heading(sql.FieldLTE(FieldName, v))
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.Heading {
	return predicate.Heading(sql.FieldContains(FieldName, v))
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.Heading {
	return predicate.Heading(sql.FieldHasPrefix(FieldName, v))
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.Heading {
	return predicate.Heading(sql.FieldHasSuffix(FieldName, v))
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.Heading {
	return predicate.Heading(sql.FieldEqualFold(FieldName, v))
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.Heading {
	return predicate.Heading(sql.FieldContainsFold(FieldName, v))
}

// HasVarsity applies the HasEdge predicate on the "varsity" edge.
func HasVarsity() predicate.Heading {
	return predicate.Heading(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, VarsityTable, VarsityColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasVarsityWith applies the HasEdge predicate on the "varsity" edge with a given conditions (other predicates).
func HasVarsityWith(preds ...predicate.Varsity) predicate.Heading {
	return predicate.Heading(func(s *sql.Selector) {
		step := newVarsityStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasApplications applies the HasEdge predicate on the "applications" edge.
func HasApplications() predicate.Heading {
	return predicate.Heading(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, ApplicationsTable, ApplicationsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasApplicationsWith applies the HasEdge predicate on the "applications" edge with a given conditions (other predicates).
func HasApplicationsWith(preds ...predicate.Application) predicate.Heading {
	return predicate.Heading(func(s *sql.Selector) {
		step := newApplicationsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasCalculations applies the HasEdge predicate on the "calculations" edge.
func HasCalculations() predicate.Heading {
	return predicate.Heading(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, CalculationsTable, CalculationsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasCalculationsWith applies the HasEdge predicate on the "calculations" edge with a given conditions (other predicates).
func HasCalculationsWith(preds ...predicate.Calculation) predicate.Heading {
	return predicate.Heading(func(s *sql.Selector) {
		step := newCalculationsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasDrainedResults applies the HasEdge predicate on the "drained_results" edge.
func HasDrainedResults() predicate.Heading {
	return predicate.Heading(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, DrainedResultsTable, DrainedResultsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasDrainedResultsWith applies the HasEdge predicate on the "drained_results" edge with a given conditions (other predicates).
func HasDrainedResultsWith(preds ...predicate.DrainedResult) predicate.Heading {
	return predicate.Heading(func(s *sql.Selector) {
		step := newDrainedResultsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Heading) predicate.Heading {
	return predicate.Heading(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Heading) predicate.Heading {
	return predicate.Heading(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Heading) predicate.Heading {
	return predicate.Heading(sql.NotPredicates(p))
}
