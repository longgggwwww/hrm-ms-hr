// Code generated by ent, DO NOT EDIT.

package leaverequest

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/longgggwwww/hrm-ms-hr/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldLTE(FieldID, id))
}

// TotalDays applies equality check predicate on the "total_days" field. It's identical to TotalDaysEQ.
func TotalDays(v float64) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldEQ(FieldTotalDays, v))
}

// StartAt applies equality check predicate on the "start_at" field. It's identical to StartAtEQ.
func StartAt(v time.Time) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldEQ(FieldStartAt, v))
}

// EndAt applies equality check predicate on the "end_at" field. It's identical to EndAtEQ.
func EndAt(v time.Time) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldEQ(FieldEndAt, v))
}

// Reason applies equality check predicate on the "reason" field. It's identical to ReasonEQ.
func Reason(v string) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldEQ(FieldReason, v))
}

// OrgID applies equality check predicate on the "org_id" field. It's identical to OrgIDEQ.
func OrgID(v int) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldEQ(FieldOrgID, v))
}

// EmployeeID applies equality check predicate on the "employee_id" field. It's identical to EmployeeIDEQ.
func EmployeeID(v int) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldEQ(FieldEmployeeID, v))
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldEQ(FieldCreatedAt, v))
}

// UpdatedAt applies equality check predicate on the "updated_at" field. It's identical to UpdatedAtEQ.
func UpdatedAt(v time.Time) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldEQ(FieldUpdatedAt, v))
}

// TotalDaysEQ applies the EQ predicate on the "total_days" field.
func TotalDaysEQ(v float64) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldEQ(FieldTotalDays, v))
}

// TotalDaysNEQ applies the NEQ predicate on the "total_days" field.
func TotalDaysNEQ(v float64) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldNEQ(FieldTotalDays, v))
}

// TotalDaysIn applies the In predicate on the "total_days" field.
func TotalDaysIn(vs ...float64) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldIn(FieldTotalDays, vs...))
}

// TotalDaysNotIn applies the NotIn predicate on the "total_days" field.
func TotalDaysNotIn(vs ...float64) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldNotIn(FieldTotalDays, vs...))
}

// TotalDaysGT applies the GT predicate on the "total_days" field.
func TotalDaysGT(v float64) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldGT(FieldTotalDays, v))
}

// TotalDaysGTE applies the GTE predicate on the "total_days" field.
func TotalDaysGTE(v float64) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldGTE(FieldTotalDays, v))
}

// TotalDaysLT applies the LT predicate on the "total_days" field.
func TotalDaysLT(v float64) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldLT(FieldTotalDays, v))
}

// TotalDaysLTE applies the LTE predicate on the "total_days" field.
func TotalDaysLTE(v float64) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldLTE(FieldTotalDays, v))
}

// StartAtEQ applies the EQ predicate on the "start_at" field.
func StartAtEQ(v time.Time) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldEQ(FieldStartAt, v))
}

// StartAtNEQ applies the NEQ predicate on the "start_at" field.
func StartAtNEQ(v time.Time) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldNEQ(FieldStartAt, v))
}

// StartAtIn applies the In predicate on the "start_at" field.
func StartAtIn(vs ...time.Time) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldIn(FieldStartAt, vs...))
}

// StartAtNotIn applies the NotIn predicate on the "start_at" field.
func StartAtNotIn(vs ...time.Time) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldNotIn(FieldStartAt, vs...))
}

// StartAtGT applies the GT predicate on the "start_at" field.
func StartAtGT(v time.Time) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldGT(FieldStartAt, v))
}

// StartAtGTE applies the GTE predicate on the "start_at" field.
func StartAtGTE(v time.Time) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldGTE(FieldStartAt, v))
}

// StartAtLT applies the LT predicate on the "start_at" field.
func StartAtLT(v time.Time) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldLT(FieldStartAt, v))
}

// StartAtLTE applies the LTE predicate on the "start_at" field.
func StartAtLTE(v time.Time) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldLTE(FieldStartAt, v))
}

// EndAtEQ applies the EQ predicate on the "end_at" field.
func EndAtEQ(v time.Time) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldEQ(FieldEndAt, v))
}

// EndAtNEQ applies the NEQ predicate on the "end_at" field.
func EndAtNEQ(v time.Time) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldNEQ(FieldEndAt, v))
}

// EndAtIn applies the In predicate on the "end_at" field.
func EndAtIn(vs ...time.Time) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldIn(FieldEndAt, vs...))
}

// EndAtNotIn applies the NotIn predicate on the "end_at" field.
func EndAtNotIn(vs ...time.Time) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldNotIn(FieldEndAt, vs...))
}

// EndAtGT applies the GT predicate on the "end_at" field.
func EndAtGT(v time.Time) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldGT(FieldEndAt, v))
}

// EndAtGTE applies the GTE predicate on the "end_at" field.
func EndAtGTE(v time.Time) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldGTE(FieldEndAt, v))
}

// EndAtLT applies the LT predicate on the "end_at" field.
func EndAtLT(v time.Time) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldLT(FieldEndAt, v))
}

// EndAtLTE applies the LTE predicate on the "end_at" field.
func EndAtLTE(v time.Time) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldLTE(FieldEndAt, v))
}

// ReasonEQ applies the EQ predicate on the "reason" field.
func ReasonEQ(v string) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldEQ(FieldReason, v))
}

// ReasonNEQ applies the NEQ predicate on the "reason" field.
func ReasonNEQ(v string) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldNEQ(FieldReason, v))
}

// ReasonIn applies the In predicate on the "reason" field.
func ReasonIn(vs ...string) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldIn(FieldReason, vs...))
}

// ReasonNotIn applies the NotIn predicate on the "reason" field.
func ReasonNotIn(vs ...string) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldNotIn(FieldReason, vs...))
}

// ReasonGT applies the GT predicate on the "reason" field.
func ReasonGT(v string) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldGT(FieldReason, v))
}

// ReasonGTE applies the GTE predicate on the "reason" field.
func ReasonGTE(v string) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldGTE(FieldReason, v))
}

// ReasonLT applies the LT predicate on the "reason" field.
func ReasonLT(v string) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldLT(FieldReason, v))
}

// ReasonLTE applies the LTE predicate on the "reason" field.
func ReasonLTE(v string) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldLTE(FieldReason, v))
}

// ReasonContains applies the Contains predicate on the "reason" field.
func ReasonContains(v string) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldContains(FieldReason, v))
}

// ReasonHasPrefix applies the HasPrefix predicate on the "reason" field.
func ReasonHasPrefix(v string) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldHasPrefix(FieldReason, v))
}

// ReasonHasSuffix applies the HasSuffix predicate on the "reason" field.
func ReasonHasSuffix(v string) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldHasSuffix(FieldReason, v))
}

// ReasonIsNil applies the IsNil predicate on the "reason" field.
func ReasonIsNil() predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldIsNull(FieldReason))
}

// ReasonNotNil applies the NotNil predicate on the "reason" field.
func ReasonNotNil() predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldNotNull(FieldReason))
}

// ReasonEqualFold applies the EqualFold predicate on the "reason" field.
func ReasonEqualFold(v string) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldEqualFold(FieldReason, v))
}

// ReasonContainsFold applies the ContainsFold predicate on the "reason" field.
func ReasonContainsFold(v string) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldContainsFold(FieldReason, v))
}

// TypeEQ applies the EQ predicate on the "type" field.
func TypeEQ(v Type) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldEQ(FieldType, v))
}

// TypeNEQ applies the NEQ predicate on the "type" field.
func TypeNEQ(v Type) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldNEQ(FieldType, v))
}

// TypeIn applies the In predicate on the "type" field.
func TypeIn(vs ...Type) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldIn(FieldType, vs...))
}

// TypeNotIn applies the NotIn predicate on the "type" field.
func TypeNotIn(vs ...Type) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldNotIn(FieldType, vs...))
}

// StatusEQ applies the EQ predicate on the "status" field.
func StatusEQ(v Status) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldEQ(FieldStatus, v))
}

// StatusNEQ applies the NEQ predicate on the "status" field.
func StatusNEQ(v Status) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldNEQ(FieldStatus, v))
}

// StatusIn applies the In predicate on the "status" field.
func StatusIn(vs ...Status) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldIn(FieldStatus, vs...))
}

// StatusNotIn applies the NotIn predicate on the "status" field.
func StatusNotIn(vs ...Status) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldNotIn(FieldStatus, vs...))
}

// OrgIDEQ applies the EQ predicate on the "org_id" field.
func OrgIDEQ(v int) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldEQ(FieldOrgID, v))
}

// OrgIDNEQ applies the NEQ predicate on the "org_id" field.
func OrgIDNEQ(v int) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldNEQ(FieldOrgID, v))
}

// OrgIDIn applies the In predicate on the "org_id" field.
func OrgIDIn(vs ...int) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldIn(FieldOrgID, vs...))
}

// OrgIDNotIn applies the NotIn predicate on the "org_id" field.
func OrgIDNotIn(vs ...int) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldNotIn(FieldOrgID, vs...))
}

// EmployeeIDEQ applies the EQ predicate on the "employee_id" field.
func EmployeeIDEQ(v int) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldEQ(FieldEmployeeID, v))
}

// EmployeeIDNEQ applies the NEQ predicate on the "employee_id" field.
func EmployeeIDNEQ(v int) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldNEQ(FieldEmployeeID, v))
}

// EmployeeIDIn applies the In predicate on the "employee_id" field.
func EmployeeIDIn(vs ...int) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldIn(FieldEmployeeID, vs...))
}

// EmployeeIDNotIn applies the NotIn predicate on the "employee_id" field.
func EmployeeIDNotIn(vs ...int) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldNotIn(FieldEmployeeID, vs...))
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldEQ(FieldCreatedAt, v))
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldNEQ(FieldCreatedAt, v))
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldIn(FieldCreatedAt, vs...))
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...time.Time) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldNotIn(FieldCreatedAt, vs...))
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v time.Time) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldGT(FieldCreatedAt, v))
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldGTE(FieldCreatedAt, v))
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldLT(FieldCreatedAt, v))
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldLTE(FieldCreatedAt, v))
}

// UpdatedAtEQ applies the EQ predicate on the "updated_at" field.
func UpdatedAtEQ(v time.Time) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldEQ(FieldUpdatedAt, v))
}

// UpdatedAtNEQ applies the NEQ predicate on the "updated_at" field.
func UpdatedAtNEQ(v time.Time) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldNEQ(FieldUpdatedAt, v))
}

// UpdatedAtIn applies the In predicate on the "updated_at" field.
func UpdatedAtIn(vs ...time.Time) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldIn(FieldUpdatedAt, vs...))
}

// UpdatedAtNotIn applies the NotIn predicate on the "updated_at" field.
func UpdatedAtNotIn(vs ...time.Time) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldNotIn(FieldUpdatedAt, vs...))
}

// UpdatedAtGT applies the GT predicate on the "updated_at" field.
func UpdatedAtGT(v time.Time) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldGT(FieldUpdatedAt, v))
}

// UpdatedAtGTE applies the GTE predicate on the "updated_at" field.
func UpdatedAtGTE(v time.Time) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldGTE(FieldUpdatedAt, v))
}

// UpdatedAtLT applies the LT predicate on the "updated_at" field.
func UpdatedAtLT(v time.Time) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldLT(FieldUpdatedAt, v))
}

// UpdatedAtLTE applies the LTE predicate on the "updated_at" field.
func UpdatedAtLTE(v time.Time) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.FieldLTE(FieldUpdatedAt, v))
}

// HasLeaveApproves applies the HasEdge predicate on the "leave_approves" edge.
func HasLeaveApproves() predicate.LeaveRequest {
	return predicate.LeaveRequest(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, LeaveApprovesTable, LeaveApprovesColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasLeaveApprovesWith applies the HasEdge predicate on the "leave_approves" edge with a given conditions (other predicates).
func HasLeaveApprovesWith(preds ...predicate.LeaveApproval) predicate.LeaveRequest {
	return predicate.LeaveRequest(func(s *sql.Selector) {
		step := newLeaveApprovesStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasApplicant applies the HasEdge predicate on the "applicant" edge.
func HasApplicant() predicate.LeaveRequest {
	return predicate.LeaveRequest(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, ApplicantTable, ApplicantColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasApplicantWith applies the HasEdge predicate on the "applicant" edge with a given conditions (other predicates).
func HasApplicantWith(preds ...predicate.Employee) predicate.LeaveRequest {
	return predicate.LeaveRequest(func(s *sql.Selector) {
		step := newApplicantStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasOrganization applies the HasEdge predicate on the "organization" edge.
func HasOrganization() predicate.LeaveRequest {
	return predicate.LeaveRequest(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, OrganizationTable, OrganizationColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasOrganizationWith applies the HasEdge predicate on the "organization" edge with a given conditions (other predicates).
func HasOrganizationWith(preds ...predicate.Organization) predicate.LeaveRequest {
	return predicate.LeaveRequest(func(s *sql.Selector) {
		step := newOrganizationStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.LeaveRequest) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.LeaveRequest) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.LeaveRequest) predicate.LeaveRequest {
	return predicate.LeaveRequest(sql.NotPredicates(p))
}
