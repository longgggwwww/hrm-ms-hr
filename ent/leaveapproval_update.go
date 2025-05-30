// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/longgggwwww/hrm-ms-hr/ent/employee"
	"github.com/longgggwwww/hrm-ms-hr/ent/leaveapproval"
	"github.com/longgggwwww/hrm-ms-hr/ent/leaverequest"
	"github.com/longgggwwww/hrm-ms-hr/ent/predicate"
)

// LeaveApprovalUpdate is the builder for updating LeaveApproval entities.
type LeaveApprovalUpdate struct {
	config
	hooks    []Hook
	mutation *LeaveApprovalMutation
}

// Where appends a list predicates to the LeaveApprovalUpdate builder.
func (lau *LeaveApprovalUpdate) Where(ps ...predicate.LeaveApproval) *LeaveApprovalUpdate {
	lau.mutation.Where(ps...)
	return lau
}

// SetComment sets the "comment" field.
func (lau *LeaveApprovalUpdate) SetComment(s string) *LeaveApprovalUpdate {
	lau.mutation.SetComment(s)
	return lau
}

// SetNillableComment sets the "comment" field if the given value is not nil.
func (lau *LeaveApprovalUpdate) SetNillableComment(s *string) *LeaveApprovalUpdate {
	if s != nil {
		lau.SetComment(*s)
	}
	return lau
}

// ClearComment clears the value of the "comment" field.
func (lau *LeaveApprovalUpdate) ClearComment() *LeaveApprovalUpdate {
	lau.mutation.ClearComment()
	return lau
}

// SetLeaveRequestID sets the "leave_request_id" field.
func (lau *LeaveApprovalUpdate) SetLeaveRequestID(i int) *LeaveApprovalUpdate {
	lau.mutation.SetLeaveRequestID(i)
	return lau
}

// SetNillableLeaveRequestID sets the "leave_request_id" field if the given value is not nil.
func (lau *LeaveApprovalUpdate) SetNillableLeaveRequestID(i *int) *LeaveApprovalUpdate {
	if i != nil {
		lau.SetLeaveRequestID(*i)
	}
	return lau
}

// SetReviewerID sets the "reviewer_id" field.
func (lau *LeaveApprovalUpdate) SetReviewerID(i int) *LeaveApprovalUpdate {
	lau.mutation.SetReviewerID(i)
	return lau
}

// SetNillableReviewerID sets the "reviewer_id" field if the given value is not nil.
func (lau *LeaveApprovalUpdate) SetNillableReviewerID(i *int) *LeaveApprovalUpdate {
	if i != nil {
		lau.SetReviewerID(*i)
	}
	return lau
}

// SetCreatedAt sets the "created_at" field.
func (lau *LeaveApprovalUpdate) SetCreatedAt(t time.Time) *LeaveApprovalUpdate {
	lau.mutation.SetCreatedAt(t)
	return lau
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (lau *LeaveApprovalUpdate) SetNillableCreatedAt(t *time.Time) *LeaveApprovalUpdate {
	if t != nil {
		lau.SetCreatedAt(*t)
	}
	return lau
}

// SetUpdatedAt sets the "updated_at" field.
func (lau *LeaveApprovalUpdate) SetUpdatedAt(t time.Time) *LeaveApprovalUpdate {
	lau.mutation.SetUpdatedAt(t)
	return lau
}

// SetLeaveRequest sets the "leave_request" edge to the LeaveRequest entity.
func (lau *LeaveApprovalUpdate) SetLeaveRequest(l *LeaveRequest) *LeaveApprovalUpdate {
	return lau.SetLeaveRequestID(l.ID)
}

// SetReviewer sets the "reviewer" edge to the Employee entity.
func (lau *LeaveApprovalUpdate) SetReviewer(e *Employee) *LeaveApprovalUpdate {
	return lau.SetReviewerID(e.ID)
}

// Mutation returns the LeaveApprovalMutation object of the builder.
func (lau *LeaveApprovalUpdate) Mutation() *LeaveApprovalMutation {
	return lau.mutation
}

// ClearLeaveRequest clears the "leave_request" edge to the LeaveRequest entity.
func (lau *LeaveApprovalUpdate) ClearLeaveRequest() *LeaveApprovalUpdate {
	lau.mutation.ClearLeaveRequest()
	return lau
}

// ClearReviewer clears the "reviewer" edge to the Employee entity.
func (lau *LeaveApprovalUpdate) ClearReviewer() *LeaveApprovalUpdate {
	lau.mutation.ClearReviewer()
	return lau
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (lau *LeaveApprovalUpdate) Save(ctx context.Context) (int, error) {
	lau.defaults()
	return withHooks(ctx, lau.sqlSave, lau.mutation, lau.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (lau *LeaveApprovalUpdate) SaveX(ctx context.Context) int {
	affected, err := lau.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (lau *LeaveApprovalUpdate) Exec(ctx context.Context) error {
	_, err := lau.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (lau *LeaveApprovalUpdate) ExecX(ctx context.Context) {
	if err := lau.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (lau *LeaveApprovalUpdate) defaults() {
	if _, ok := lau.mutation.UpdatedAt(); !ok {
		v := leaveapproval.UpdateDefaultUpdatedAt()
		lau.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (lau *LeaveApprovalUpdate) check() error {
	if lau.mutation.LeaveRequestCleared() && len(lau.mutation.LeaveRequestIDs()) > 0 {
		return errors.New(`ent: clearing a required unique edge "LeaveApproval.leave_request"`)
	}
	if lau.mutation.ReviewerCleared() && len(lau.mutation.ReviewerIDs()) > 0 {
		return errors.New(`ent: clearing a required unique edge "LeaveApproval.reviewer"`)
	}
	return nil
}

func (lau *LeaveApprovalUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := lau.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(leaveapproval.Table, leaveapproval.Columns, sqlgraph.NewFieldSpec(leaveapproval.FieldID, field.TypeInt))
	if ps := lau.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := lau.mutation.Comment(); ok {
		_spec.SetField(leaveapproval.FieldComment, field.TypeString, value)
	}
	if lau.mutation.CommentCleared() {
		_spec.ClearField(leaveapproval.FieldComment, field.TypeString)
	}
	if value, ok := lau.mutation.CreatedAt(); ok {
		_spec.SetField(leaveapproval.FieldCreatedAt, field.TypeTime, value)
	}
	if value, ok := lau.mutation.UpdatedAt(); ok {
		_spec.SetField(leaveapproval.FieldUpdatedAt, field.TypeTime, value)
	}
	if lau.mutation.LeaveRequestCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   leaveapproval.LeaveRequestTable,
			Columns: []string{leaveapproval.LeaveRequestColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(leaverequest.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := lau.mutation.LeaveRequestIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   leaveapproval.LeaveRequestTable,
			Columns: []string{leaveapproval.LeaveRequestColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(leaverequest.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if lau.mutation.ReviewerCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   leaveapproval.ReviewerTable,
			Columns: []string{leaveapproval.ReviewerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(employee.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := lau.mutation.ReviewerIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   leaveapproval.ReviewerTable,
			Columns: []string{leaveapproval.ReviewerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(employee.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, lau.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{leaveapproval.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	lau.mutation.done = true
	return n, nil
}

// LeaveApprovalUpdateOne is the builder for updating a single LeaveApproval entity.
type LeaveApprovalUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *LeaveApprovalMutation
}

// SetComment sets the "comment" field.
func (lauo *LeaveApprovalUpdateOne) SetComment(s string) *LeaveApprovalUpdateOne {
	lauo.mutation.SetComment(s)
	return lauo
}

// SetNillableComment sets the "comment" field if the given value is not nil.
func (lauo *LeaveApprovalUpdateOne) SetNillableComment(s *string) *LeaveApprovalUpdateOne {
	if s != nil {
		lauo.SetComment(*s)
	}
	return lauo
}

// ClearComment clears the value of the "comment" field.
func (lauo *LeaveApprovalUpdateOne) ClearComment() *LeaveApprovalUpdateOne {
	lauo.mutation.ClearComment()
	return lauo
}

// SetLeaveRequestID sets the "leave_request_id" field.
func (lauo *LeaveApprovalUpdateOne) SetLeaveRequestID(i int) *LeaveApprovalUpdateOne {
	lauo.mutation.SetLeaveRequestID(i)
	return lauo
}

// SetNillableLeaveRequestID sets the "leave_request_id" field if the given value is not nil.
func (lauo *LeaveApprovalUpdateOne) SetNillableLeaveRequestID(i *int) *LeaveApprovalUpdateOne {
	if i != nil {
		lauo.SetLeaveRequestID(*i)
	}
	return lauo
}

// SetReviewerID sets the "reviewer_id" field.
func (lauo *LeaveApprovalUpdateOne) SetReviewerID(i int) *LeaveApprovalUpdateOne {
	lauo.mutation.SetReviewerID(i)
	return lauo
}

// SetNillableReviewerID sets the "reviewer_id" field if the given value is not nil.
func (lauo *LeaveApprovalUpdateOne) SetNillableReviewerID(i *int) *LeaveApprovalUpdateOne {
	if i != nil {
		lauo.SetReviewerID(*i)
	}
	return lauo
}

// SetCreatedAt sets the "created_at" field.
func (lauo *LeaveApprovalUpdateOne) SetCreatedAt(t time.Time) *LeaveApprovalUpdateOne {
	lauo.mutation.SetCreatedAt(t)
	return lauo
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (lauo *LeaveApprovalUpdateOne) SetNillableCreatedAt(t *time.Time) *LeaveApprovalUpdateOne {
	if t != nil {
		lauo.SetCreatedAt(*t)
	}
	return lauo
}

// SetUpdatedAt sets the "updated_at" field.
func (lauo *LeaveApprovalUpdateOne) SetUpdatedAt(t time.Time) *LeaveApprovalUpdateOne {
	lauo.mutation.SetUpdatedAt(t)
	return lauo
}

// SetLeaveRequest sets the "leave_request" edge to the LeaveRequest entity.
func (lauo *LeaveApprovalUpdateOne) SetLeaveRequest(l *LeaveRequest) *LeaveApprovalUpdateOne {
	return lauo.SetLeaveRequestID(l.ID)
}

// SetReviewer sets the "reviewer" edge to the Employee entity.
func (lauo *LeaveApprovalUpdateOne) SetReviewer(e *Employee) *LeaveApprovalUpdateOne {
	return lauo.SetReviewerID(e.ID)
}

// Mutation returns the LeaveApprovalMutation object of the builder.
func (lauo *LeaveApprovalUpdateOne) Mutation() *LeaveApprovalMutation {
	return lauo.mutation
}

// ClearLeaveRequest clears the "leave_request" edge to the LeaveRequest entity.
func (lauo *LeaveApprovalUpdateOne) ClearLeaveRequest() *LeaveApprovalUpdateOne {
	lauo.mutation.ClearLeaveRequest()
	return lauo
}

// ClearReviewer clears the "reviewer" edge to the Employee entity.
func (lauo *LeaveApprovalUpdateOne) ClearReviewer() *LeaveApprovalUpdateOne {
	lauo.mutation.ClearReviewer()
	return lauo
}

// Where appends a list predicates to the LeaveApprovalUpdate builder.
func (lauo *LeaveApprovalUpdateOne) Where(ps ...predicate.LeaveApproval) *LeaveApprovalUpdateOne {
	lauo.mutation.Where(ps...)
	return lauo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (lauo *LeaveApprovalUpdateOne) Select(field string, fields ...string) *LeaveApprovalUpdateOne {
	lauo.fields = append([]string{field}, fields...)
	return lauo
}

// Save executes the query and returns the updated LeaveApproval entity.
func (lauo *LeaveApprovalUpdateOne) Save(ctx context.Context) (*LeaveApproval, error) {
	lauo.defaults()
	return withHooks(ctx, lauo.sqlSave, lauo.mutation, lauo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (lauo *LeaveApprovalUpdateOne) SaveX(ctx context.Context) *LeaveApproval {
	node, err := lauo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (lauo *LeaveApprovalUpdateOne) Exec(ctx context.Context) error {
	_, err := lauo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (lauo *LeaveApprovalUpdateOne) ExecX(ctx context.Context) {
	if err := lauo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (lauo *LeaveApprovalUpdateOne) defaults() {
	if _, ok := lauo.mutation.UpdatedAt(); !ok {
		v := leaveapproval.UpdateDefaultUpdatedAt()
		lauo.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (lauo *LeaveApprovalUpdateOne) check() error {
	if lauo.mutation.LeaveRequestCleared() && len(lauo.mutation.LeaveRequestIDs()) > 0 {
		return errors.New(`ent: clearing a required unique edge "LeaveApproval.leave_request"`)
	}
	if lauo.mutation.ReviewerCleared() && len(lauo.mutation.ReviewerIDs()) > 0 {
		return errors.New(`ent: clearing a required unique edge "LeaveApproval.reviewer"`)
	}
	return nil
}

func (lauo *LeaveApprovalUpdateOne) sqlSave(ctx context.Context) (_node *LeaveApproval, err error) {
	if err := lauo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(leaveapproval.Table, leaveapproval.Columns, sqlgraph.NewFieldSpec(leaveapproval.FieldID, field.TypeInt))
	id, ok := lauo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "LeaveApproval.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := lauo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, leaveapproval.FieldID)
		for _, f := range fields {
			if !leaveapproval.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != leaveapproval.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := lauo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := lauo.mutation.Comment(); ok {
		_spec.SetField(leaveapproval.FieldComment, field.TypeString, value)
	}
	if lauo.mutation.CommentCleared() {
		_spec.ClearField(leaveapproval.FieldComment, field.TypeString)
	}
	if value, ok := lauo.mutation.CreatedAt(); ok {
		_spec.SetField(leaveapproval.FieldCreatedAt, field.TypeTime, value)
	}
	if value, ok := lauo.mutation.UpdatedAt(); ok {
		_spec.SetField(leaveapproval.FieldUpdatedAt, field.TypeTime, value)
	}
	if lauo.mutation.LeaveRequestCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   leaveapproval.LeaveRequestTable,
			Columns: []string{leaveapproval.LeaveRequestColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(leaverequest.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := lauo.mutation.LeaveRequestIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   leaveapproval.LeaveRequestTable,
			Columns: []string{leaveapproval.LeaveRequestColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(leaverequest.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if lauo.mutation.ReviewerCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   leaveapproval.ReviewerTable,
			Columns: []string{leaveapproval.ReviewerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(employee.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := lauo.mutation.ReviewerIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   leaveapproval.ReviewerTable,
			Columns: []string{leaveapproval.ReviewerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(employee.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &LeaveApproval{config: lauo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, lauo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{leaveapproval.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	lauo.mutation.done = true
	return _node, nil
}
