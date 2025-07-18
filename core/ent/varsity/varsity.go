// Code generated by ent, DO NOT EDIT.

package varsity

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the varsity type in the database.
	Label = "varsity"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCode holds the string denoting the code field in the database.
	FieldCode = "code"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// EdgeHeadings holds the string denoting the headings edge name in mutations.
	EdgeHeadings = "headings"
	// Table holds the table name of the varsity in the database.
	Table = "varsities"
	// HeadingsTable is the table that holds the headings relation/edge.
	HeadingsTable = "headings"
	// HeadingsInverseTable is the table name for the Heading entity.
	// It exists in this package in order to avoid circular dependency with the "heading" package.
	HeadingsInverseTable = "headings"
	// HeadingsColumn is the table column denoting the headings relation/edge.
	HeadingsColumn = "varsity_headings"
)

// Columns holds all SQL columns for varsity fields.
var Columns = []string{
	FieldID,
	FieldCode,
	FieldName,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

// OrderOption defines the ordering options for the Varsity queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByCode orders the results by the code field.
func ByCode(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCode, opts...).ToFunc()
}

// ByName orders the results by the name field.
func ByName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldName, opts...).ToFunc()
}

// ByHeadingsCount orders the results by headings count.
func ByHeadingsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newHeadingsStep(), opts...)
	}
}

// ByHeadings orders the results by headings terms.
func ByHeadings(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newHeadingsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newHeadingsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(HeadingsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, HeadingsTable, HeadingsColumn),
	)
}
