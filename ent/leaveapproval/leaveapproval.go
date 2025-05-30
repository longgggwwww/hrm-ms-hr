// Code generated by ent, DO NOT EDIT.

package leaveapproval

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the leaveapproval type in the database.
	Label = "leave_approval"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldComment holds the string denoting the comment field in the database.
	FieldComment = "comment"
	// FieldLeaveRequestID holds the string denoting the leave_request_id field in the database.
	FieldLeaveRequestID = "leave_request_id"
	// FieldReviewerID holds the string denoting the reviewer_id field in the database.
	FieldReviewerID = "reviewer_id"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// EdgeLeaveRequest holds the string denoting the leave_request edge name in mutations.
	EdgeLeaveRequest = "leave_request"
	// EdgeReviewer holds the string denoting the reviewer edge name in mutations.
	EdgeReviewer = "reviewer"
	// Table holds the table name of the leaveapproval in the database.
	Table = "leave_approvals"
	// LeaveRequestTable is the table that holds the leave_request relation/edge.
	LeaveRequestTable = "leave_approvals"
	// LeaveRequestInverseTable is the table name for the LeaveRequest entity.
	// It exists in this package in order to avoid circular dependency with the "leaverequest" package.
	LeaveRequestInverseTable = "leave_requests"
	// LeaveRequestColumn is the table column denoting the leave_request relation/edge.
	LeaveRequestColumn = "leave_request_id"
	// ReviewerTable is the table that holds the reviewer relation/edge.
	ReviewerTable = "leave_approvals"
	// ReviewerInverseTable is the table name for the Employee entity.
	// It exists in this package in order to avoid circular dependency with the "employee" package.
	ReviewerInverseTable = "employees"
	// ReviewerColumn is the table column denoting the reviewer relation/edge.
	ReviewerColumn = "reviewer_id"
)

// Columns holds all SQL columns for leaveapproval fields.
var Columns = []string{
	FieldID,
	FieldComment,
	FieldLeaveRequestID,
	FieldReviewerID,
	FieldCreatedAt,
	FieldUpdatedAt,
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

var (
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt func() time.Time
	// DefaultUpdatedAt holds the default value on creation for the "updated_at" field.
	DefaultUpdatedAt func() time.Time
	// UpdateDefaultUpdatedAt holds the default value on update for the "updated_at" field.
	UpdateDefaultUpdatedAt func() time.Time
)

// OrderOption defines the ordering options for the LeaveApproval queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByComment orders the results by the comment field.
func ByComment(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldComment, opts...).ToFunc()
}

// ByLeaveRequestID orders the results by the leave_request_id field.
func ByLeaveRequestID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldLeaveRequestID, opts...).ToFunc()
}

// ByReviewerID orders the results by the reviewer_id field.
func ByReviewerID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldReviewerID, opts...).ToFunc()
}

// ByCreatedAt orders the results by the created_at field.
func ByCreatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreatedAt, opts...).ToFunc()
}

// ByUpdatedAt orders the results by the updated_at field.
func ByUpdatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUpdatedAt, opts...).ToFunc()
}

// ByLeaveRequestField orders the results by leave_request field.
func ByLeaveRequestField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newLeaveRequestStep(), sql.OrderByField(field, opts...))
	}
}

// ByReviewerField orders the results by reviewer field.
func ByReviewerField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newReviewerStep(), sql.OrderByField(field, opts...))
	}
}
func newLeaveRequestStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(LeaveRequestInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, LeaveRequestTable, LeaveRequestColumn),
	)
}
func newReviewerStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(ReviewerInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, ReviewerTable, ReviewerColumn),
	)
}
