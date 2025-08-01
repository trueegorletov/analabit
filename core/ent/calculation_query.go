// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"math"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/trueegorletov/analabit/core/ent/calculation"
	"github.com/trueegorletov/analabit/core/ent/heading"
	"github.com/trueegorletov/analabit/core/ent/predicate"
	"github.com/trueegorletov/analabit/core/ent/run"
)

// CalculationQuery is the builder for querying Calculation entities.
type CalculationQuery struct {
	config
	ctx         *QueryContext
	order       []calculation.OrderOption
	inters      []Interceptor
	predicates  []predicate.Calculation
	withHeading *HeadingQuery
	withRun     *RunQuery
	withFKs     bool
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the CalculationQuery builder.
func (cq *CalculationQuery) Where(ps ...predicate.Calculation) *CalculationQuery {
	cq.predicates = append(cq.predicates, ps...)
	return cq
}

// Limit the number of records to be returned by this query.
func (cq *CalculationQuery) Limit(limit int) *CalculationQuery {
	cq.ctx.Limit = &limit
	return cq
}

// Offset to start from.
func (cq *CalculationQuery) Offset(offset int) *CalculationQuery {
	cq.ctx.Offset = &offset
	return cq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (cq *CalculationQuery) Unique(unique bool) *CalculationQuery {
	cq.ctx.Unique = &unique
	return cq
}

// Order specifies how the records should be ordered.
func (cq *CalculationQuery) Order(o ...calculation.OrderOption) *CalculationQuery {
	cq.order = append(cq.order, o...)
	return cq
}

// QueryHeading chains the current query on the "heading" edge.
func (cq *CalculationQuery) QueryHeading() *HeadingQuery {
	query := (&HeadingClient{config: cq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := cq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := cq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(calculation.Table, calculation.FieldID, selector),
			sqlgraph.To(heading.Table, heading.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, calculation.HeadingTable, calculation.HeadingColumn),
		)
		fromU = sqlgraph.SetNeighbors(cq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryRun chains the current query on the "run" edge.
func (cq *CalculationQuery) QueryRun() *RunQuery {
	query := (&RunClient{config: cq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := cq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := cq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(calculation.Table, calculation.FieldID, selector),
			sqlgraph.To(run.Table, run.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, calculation.RunTable, calculation.RunColumn),
		)
		fromU = sqlgraph.SetNeighbors(cq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first Calculation entity from the query.
// Returns a *NotFoundError when no Calculation was found.
func (cq *CalculationQuery) First(ctx context.Context) (*Calculation, error) {
	nodes, err := cq.Limit(1).All(setContextOp(ctx, cq.ctx, ent.OpQueryFirst))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{calculation.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (cq *CalculationQuery) FirstX(ctx context.Context) *Calculation {
	node, err := cq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first Calculation ID from the query.
// Returns a *NotFoundError when no Calculation ID was found.
func (cq *CalculationQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = cq.Limit(1).IDs(setContextOp(ctx, cq.ctx, ent.OpQueryFirstID)); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{calculation.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (cq *CalculationQuery) FirstIDX(ctx context.Context) int {
	id, err := cq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single Calculation entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one Calculation entity is found.
// Returns a *NotFoundError when no Calculation entities are found.
func (cq *CalculationQuery) Only(ctx context.Context) (*Calculation, error) {
	nodes, err := cq.Limit(2).All(setContextOp(ctx, cq.ctx, ent.OpQueryOnly))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{calculation.Label}
	default:
		return nil, &NotSingularError{calculation.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (cq *CalculationQuery) OnlyX(ctx context.Context) *Calculation {
	node, err := cq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only Calculation ID in the query.
// Returns a *NotSingularError when more than one Calculation ID is found.
// Returns a *NotFoundError when no entities are found.
func (cq *CalculationQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = cq.Limit(2).IDs(setContextOp(ctx, cq.ctx, ent.OpQueryOnlyID)); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{calculation.Label}
	default:
		err = &NotSingularError{calculation.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (cq *CalculationQuery) OnlyIDX(ctx context.Context) int {
	id, err := cq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of Calculations.
func (cq *CalculationQuery) All(ctx context.Context) ([]*Calculation, error) {
	ctx = setContextOp(ctx, cq.ctx, ent.OpQueryAll)
	if err := cq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*Calculation, *CalculationQuery]()
	return withInterceptors[[]*Calculation](ctx, cq, qr, cq.inters)
}

// AllX is like All, but panics if an error occurs.
func (cq *CalculationQuery) AllX(ctx context.Context) []*Calculation {
	nodes, err := cq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of Calculation IDs.
func (cq *CalculationQuery) IDs(ctx context.Context) (ids []int, err error) {
	if cq.ctx.Unique == nil && cq.path != nil {
		cq.Unique(true)
	}
	ctx = setContextOp(ctx, cq.ctx, ent.OpQueryIDs)
	if err = cq.Select(calculation.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (cq *CalculationQuery) IDsX(ctx context.Context) []int {
	ids, err := cq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (cq *CalculationQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, cq.ctx, ent.OpQueryCount)
	if err := cq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, cq, querierCount[*CalculationQuery](), cq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (cq *CalculationQuery) CountX(ctx context.Context) int {
	count, err := cq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (cq *CalculationQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, cq.ctx, ent.OpQueryExist)
	switch _, err := cq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (cq *CalculationQuery) ExistX(ctx context.Context) bool {
	exist, err := cq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the CalculationQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (cq *CalculationQuery) Clone() *CalculationQuery {
	if cq == nil {
		return nil
	}
	return &CalculationQuery{
		config:      cq.config,
		ctx:         cq.ctx.Clone(),
		order:       append([]calculation.OrderOption{}, cq.order...),
		inters:      append([]Interceptor{}, cq.inters...),
		predicates:  append([]predicate.Calculation{}, cq.predicates...),
		withHeading: cq.withHeading.Clone(),
		withRun:     cq.withRun.Clone(),
		// clone intermediate query.
		sql:  cq.sql.Clone(),
		path: cq.path,
	}
}

// WithHeading tells the query-builder to eager-load the nodes that are connected to
// the "heading" edge. The optional arguments are used to configure the query builder of the edge.
func (cq *CalculationQuery) WithHeading(opts ...func(*HeadingQuery)) *CalculationQuery {
	query := (&HeadingClient{config: cq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	cq.withHeading = query
	return cq
}

// WithRun tells the query-builder to eager-load the nodes that are connected to
// the "run" edge. The optional arguments are used to configure the query builder of the edge.
func (cq *CalculationQuery) WithRun(opts ...func(*RunQuery)) *CalculationQuery {
	query := (&RunClient{config: cq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	cq.withRun = query
	return cq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		StudentID string `json:"student_id,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.Calculation.Query().
//		GroupBy(calculation.FieldStudentID).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (cq *CalculationQuery) GroupBy(field string, fields ...string) *CalculationGroupBy {
	cq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &CalculationGroupBy{build: cq}
	grbuild.flds = &cq.ctx.Fields
	grbuild.label = calculation.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		StudentID string `json:"student_id,omitempty"`
//	}
//
//	client.Calculation.Query().
//		Select(calculation.FieldStudentID).
//		Scan(ctx, &v)
func (cq *CalculationQuery) Select(fields ...string) *CalculationSelect {
	cq.ctx.Fields = append(cq.ctx.Fields, fields...)
	sbuild := &CalculationSelect{CalculationQuery: cq}
	sbuild.label = calculation.Label
	sbuild.flds, sbuild.scan = &cq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a CalculationSelect configured with the given aggregations.
func (cq *CalculationQuery) Aggregate(fns ...AggregateFunc) *CalculationSelect {
	return cq.Select().Aggregate(fns...)
}

func (cq *CalculationQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range cq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, cq); err != nil {
				return err
			}
		}
	}
	for _, f := range cq.ctx.Fields {
		if !calculation.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if cq.path != nil {
		prev, err := cq.path(ctx)
		if err != nil {
			return err
		}
		cq.sql = prev
	}
	return nil
}

func (cq *CalculationQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*Calculation, error) {
	var (
		nodes       = []*Calculation{}
		withFKs     = cq.withFKs
		_spec       = cq.querySpec()
		loadedTypes = [2]bool{
			cq.withHeading != nil,
			cq.withRun != nil,
		}
	)
	if cq.withHeading != nil {
		withFKs = true
	}
	if withFKs {
		_spec.Node.Columns = append(_spec.Node.Columns, calculation.ForeignKeys...)
	}
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*Calculation).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &Calculation{config: cq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, cq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := cq.withHeading; query != nil {
		if err := cq.loadHeading(ctx, query, nodes, nil,
			func(n *Calculation, e *Heading) { n.Edges.Heading = e }); err != nil {
			return nil, err
		}
	}
	if query := cq.withRun; query != nil {
		if err := cq.loadRun(ctx, query, nodes, nil,
			func(n *Calculation, e *Run) { n.Edges.Run = e }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (cq *CalculationQuery) loadHeading(ctx context.Context, query *HeadingQuery, nodes []*Calculation, init func(*Calculation), assign func(*Calculation, *Heading)) error {
	ids := make([]int, 0, len(nodes))
	nodeids := make(map[int][]*Calculation)
	for i := range nodes {
		if nodes[i].heading_calculations == nil {
			continue
		}
		fk := *nodes[i].heading_calculations
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(heading.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "heading_calculations" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}
func (cq *CalculationQuery) loadRun(ctx context.Context, query *RunQuery, nodes []*Calculation, init func(*Calculation), assign func(*Calculation, *Run)) error {
	ids := make([]int, 0, len(nodes))
	nodeids := make(map[int][]*Calculation)
	for i := range nodes {
		fk := nodes[i].RunID
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(run.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "run_id" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}

func (cq *CalculationQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := cq.querySpec()
	_spec.Node.Columns = cq.ctx.Fields
	if len(cq.ctx.Fields) > 0 {
		_spec.Unique = cq.ctx.Unique != nil && *cq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, cq.driver, _spec)
}

func (cq *CalculationQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(calculation.Table, calculation.Columns, sqlgraph.NewFieldSpec(calculation.FieldID, field.TypeInt))
	_spec.From = cq.sql
	if unique := cq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if cq.path != nil {
		_spec.Unique = true
	}
	if fields := cq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, calculation.FieldID)
		for i := range fields {
			if fields[i] != calculation.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
		if cq.withRun != nil {
			_spec.Node.AddColumnOnce(calculation.FieldRunID)
		}
	}
	if ps := cq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := cq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := cq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := cq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (cq *CalculationQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(cq.driver.Dialect())
	t1 := builder.Table(calculation.Table)
	columns := cq.ctx.Fields
	if len(columns) == 0 {
		columns = calculation.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if cq.sql != nil {
		selector = cq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if cq.ctx.Unique != nil && *cq.ctx.Unique {
		selector.Distinct()
	}
	for _, p := range cq.predicates {
		p(selector)
	}
	for _, p := range cq.order {
		p(selector)
	}
	if offset := cq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := cq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// CalculationGroupBy is the group-by builder for Calculation entities.
type CalculationGroupBy struct {
	selector
	build *CalculationQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (cgb *CalculationGroupBy) Aggregate(fns ...AggregateFunc) *CalculationGroupBy {
	cgb.fns = append(cgb.fns, fns...)
	return cgb
}

// Scan applies the selector query and scans the result into the given value.
func (cgb *CalculationGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, cgb.build.ctx, ent.OpQueryGroupBy)
	if err := cgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*CalculationQuery, *CalculationGroupBy](ctx, cgb.build, cgb, cgb.build.inters, v)
}

func (cgb *CalculationGroupBy) sqlScan(ctx context.Context, root *CalculationQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(cgb.fns))
	for _, fn := range cgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*cgb.flds)+len(cgb.fns))
		for _, f := range *cgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*cgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := cgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// CalculationSelect is the builder for selecting fields of Calculation entities.
type CalculationSelect struct {
	*CalculationQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (cs *CalculationSelect) Aggregate(fns ...AggregateFunc) *CalculationSelect {
	cs.fns = append(cs.fns, fns...)
	return cs
}

// Scan applies the selector query and scans the result into the given value.
func (cs *CalculationSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, cs.ctx, ent.OpQuerySelect)
	if err := cs.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*CalculationQuery, *CalculationSelect](ctx, cs.CalculationQuery, cs, cs.inters, v)
}

func (cs *CalculationSelect) sqlScan(ctx context.Context, root *CalculationQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(cs.fns))
	for _, fn := range cs.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*cs.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := cs.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
