// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/trueegorletov/analabit/core"
	"github.com/trueegorletov/analabit/core/ent/application"
	"github.com/trueegorletov/analabit/core/ent/heading"
	"github.com/trueegorletov/analabit/core/ent/run"
)

// ApplicationCreate is the builder for creating a Application entity.
type ApplicationCreate struct {
	config
	mutation *ApplicationMutation
	hooks    []Hook
}

// SetStudentID sets the "student_id" field.
func (ac *ApplicationCreate) SetStudentID(s string) *ApplicationCreate {
	ac.mutation.SetStudentID(s)
	return ac
}

// SetPriority sets the "priority" field.
func (ac *ApplicationCreate) SetPriority(i int) *ApplicationCreate {
	ac.mutation.SetPriority(i)
	return ac
}

// SetCompetitionType sets the "competition_type" field.
func (ac *ApplicationCreate) SetCompetitionType(c core.Competition) *ApplicationCreate {
	ac.mutation.SetCompetitionType(c)
	return ac
}

// SetRatingPlace sets the "rating_place" field.
func (ac *ApplicationCreate) SetRatingPlace(i int) *ApplicationCreate {
	ac.mutation.SetRatingPlace(i)
	return ac
}

// SetScore sets the "score" field.
func (ac *ApplicationCreate) SetScore(i int) *ApplicationCreate {
	ac.mutation.SetScore(i)
	return ac
}

// SetRunID sets the "run_id" field.
func (ac *ApplicationCreate) SetRunID(i int) *ApplicationCreate {
	ac.mutation.SetRunID(i)
	return ac
}

// SetOriginalSubmitted sets the "original_submitted" field.
func (ac *ApplicationCreate) SetOriginalSubmitted(b bool) *ApplicationCreate {
	ac.mutation.SetOriginalSubmitted(b)
	return ac
}

// SetNillableOriginalSubmitted sets the "original_submitted" field if the given value is not nil.
func (ac *ApplicationCreate) SetNillableOriginalSubmitted(b *bool) *ApplicationCreate {
	if b != nil {
		ac.SetOriginalSubmitted(*b)
	}
	return ac
}

// SetUpdatedAt sets the "updated_at" field.
func (ac *ApplicationCreate) SetUpdatedAt(t time.Time) *ApplicationCreate {
	ac.mutation.SetUpdatedAt(t)
	return ac
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (ac *ApplicationCreate) SetNillableUpdatedAt(t *time.Time) *ApplicationCreate {
	if t != nil {
		ac.SetUpdatedAt(*t)
	}
	return ac
}

// SetHeadingID sets the "heading" edge to the Heading entity by ID.
func (ac *ApplicationCreate) SetHeadingID(id int) *ApplicationCreate {
	ac.mutation.SetHeadingID(id)
	return ac
}

// SetHeading sets the "heading" edge to the Heading entity.
func (ac *ApplicationCreate) SetHeading(h *Heading) *ApplicationCreate {
	return ac.SetHeadingID(h.ID)
}

// SetRun sets the "run" edge to the Run entity.
func (ac *ApplicationCreate) SetRun(r *Run) *ApplicationCreate {
	return ac.SetRunID(r.ID)
}

// Mutation returns the ApplicationMutation object of the builder.
func (ac *ApplicationCreate) Mutation() *ApplicationMutation {
	return ac.mutation
}

// Save creates the Application in the database.
func (ac *ApplicationCreate) Save(ctx context.Context) (*Application, error) {
	ac.defaults()
	return withHooks(ctx, ac.sqlSave, ac.mutation, ac.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (ac *ApplicationCreate) SaveX(ctx context.Context) *Application {
	v, err := ac.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ac *ApplicationCreate) Exec(ctx context.Context) error {
	_, err := ac.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ac *ApplicationCreate) ExecX(ctx context.Context) {
	if err := ac.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (ac *ApplicationCreate) defaults() {
	if _, ok := ac.mutation.OriginalSubmitted(); !ok {
		v := application.DefaultOriginalSubmitted
		ac.mutation.SetOriginalSubmitted(v)
	}
	if _, ok := ac.mutation.UpdatedAt(); !ok {
		v := application.DefaultUpdatedAt()
		ac.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (ac *ApplicationCreate) check() error {
	if _, ok := ac.mutation.StudentID(); !ok {
		return &ValidationError{Name: "student_id", err: errors.New(`ent: missing required field "Application.student_id"`)}
	}
	if _, ok := ac.mutation.Priority(); !ok {
		return &ValidationError{Name: "priority", err: errors.New(`ent: missing required field "Application.priority"`)}
	}
	if _, ok := ac.mutation.CompetitionType(); !ok {
		return &ValidationError{Name: "competition_type", err: errors.New(`ent: missing required field "Application.competition_type"`)}
	}
	if _, ok := ac.mutation.RatingPlace(); !ok {
		return &ValidationError{Name: "rating_place", err: errors.New(`ent: missing required field "Application.rating_place"`)}
	}
	if _, ok := ac.mutation.Score(); !ok {
		return &ValidationError{Name: "score", err: errors.New(`ent: missing required field "Application.score"`)}
	}
	if _, ok := ac.mutation.RunID(); !ok {
		return &ValidationError{Name: "run_id", err: errors.New(`ent: missing required field "Application.run_id"`)}
	}
	if _, ok := ac.mutation.OriginalSubmitted(); !ok {
		return &ValidationError{Name: "original_submitted", err: errors.New(`ent: missing required field "Application.original_submitted"`)}
	}
	if _, ok := ac.mutation.UpdatedAt(); !ok {
		return &ValidationError{Name: "updated_at", err: errors.New(`ent: missing required field "Application.updated_at"`)}
	}
	if len(ac.mutation.HeadingIDs()) == 0 {
		return &ValidationError{Name: "heading", err: errors.New(`ent: missing required edge "Application.heading"`)}
	}
	if len(ac.mutation.RunIDs()) == 0 {
		return &ValidationError{Name: "run", err: errors.New(`ent: missing required edge "Application.run"`)}
	}
	return nil
}

func (ac *ApplicationCreate) sqlSave(ctx context.Context) (*Application, error) {
	if err := ac.check(); err != nil {
		return nil, err
	}
	_node, _spec := ac.createSpec()
	if err := sqlgraph.CreateNode(ctx, ac.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	ac.mutation.id = &_node.ID
	ac.mutation.done = true
	return _node, nil
}

func (ac *ApplicationCreate) createSpec() (*Application, *sqlgraph.CreateSpec) {
	var (
		_node = &Application{config: ac.config}
		_spec = sqlgraph.NewCreateSpec(application.Table, sqlgraph.NewFieldSpec(application.FieldID, field.TypeInt))
	)
	if value, ok := ac.mutation.StudentID(); ok {
		_spec.SetField(application.FieldStudentID, field.TypeString, value)
		_node.StudentID = value
	}
	if value, ok := ac.mutation.Priority(); ok {
		_spec.SetField(application.FieldPriority, field.TypeInt, value)
		_node.Priority = value
	}
	if value, ok := ac.mutation.CompetitionType(); ok {
		_spec.SetField(application.FieldCompetitionType, field.TypeInt, value)
		_node.CompetitionType = value
	}
	if value, ok := ac.mutation.RatingPlace(); ok {
		_spec.SetField(application.FieldRatingPlace, field.TypeInt, value)
		_node.RatingPlace = value
	}
	if value, ok := ac.mutation.Score(); ok {
		_spec.SetField(application.FieldScore, field.TypeInt, value)
		_node.Score = value
	}
	if value, ok := ac.mutation.OriginalSubmitted(); ok {
		_spec.SetField(application.FieldOriginalSubmitted, field.TypeBool, value)
		_node.OriginalSubmitted = value
	}
	if value, ok := ac.mutation.UpdatedAt(); ok {
		_spec.SetField(application.FieldUpdatedAt, field.TypeTime, value)
		_node.UpdatedAt = value
	}
	if nodes := ac.mutation.HeadingIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   application.HeadingTable,
			Columns: []string{application.HeadingColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(heading.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.heading_applications = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := ac.mutation.RunIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   application.RunTable,
			Columns: []string{application.RunColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(run.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.RunID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// ApplicationCreateBulk is the builder for creating many Application entities in bulk.
type ApplicationCreateBulk struct {
	config
	err      error
	builders []*ApplicationCreate
}

// Save creates the Application entities in the database.
func (acb *ApplicationCreateBulk) Save(ctx context.Context) ([]*Application, error) {
	if acb.err != nil {
		return nil, acb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(acb.builders))
	nodes := make([]*Application, len(acb.builders))
	mutators := make([]Mutator, len(acb.builders))
	for i := range acb.builders {
		func(i int, root context.Context) {
			builder := acb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*ApplicationMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, acb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, acb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, acb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (acb *ApplicationCreateBulk) SaveX(ctx context.Context) []*Application {
	v, err := acb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (acb *ApplicationCreateBulk) Exec(ctx context.Context) error {
	_, err := acb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (acb *ApplicationCreateBulk) ExecX(ctx context.Context) {
	if err := acb.Exec(ctx); err != nil {
		panic(err)
	}
}
