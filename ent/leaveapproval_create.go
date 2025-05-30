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
)

// LeaveApprovalCreate is the builder for creating a LeaveApproval entity.
type LeaveApprovalCreate struct {
	config
	mutation *LeaveApprovalMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetComment sets the "comment" field.
func (lac *LeaveApprovalCreate) SetComment(s string) *LeaveApprovalCreate {
	lac.mutation.SetComment(s)
	return lac
}

// SetNillableComment sets the "comment" field if the given value is not nil.
func (lac *LeaveApprovalCreate) SetNillableComment(s *string) *LeaveApprovalCreate {
	if s != nil {
		lac.SetComment(*s)
	}
	return lac
}

// SetLeaveRequestID sets the "leave_request_id" field.
func (lac *LeaveApprovalCreate) SetLeaveRequestID(i int) *LeaveApprovalCreate {
	lac.mutation.SetLeaveRequestID(i)
	return lac
}

// SetReviewerID sets the "reviewer_id" field.
func (lac *LeaveApprovalCreate) SetReviewerID(i int) *LeaveApprovalCreate {
	lac.mutation.SetReviewerID(i)
	return lac
}

// SetCreatedAt sets the "created_at" field.
func (lac *LeaveApprovalCreate) SetCreatedAt(t time.Time) *LeaveApprovalCreate {
	lac.mutation.SetCreatedAt(t)
	return lac
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (lac *LeaveApprovalCreate) SetNillableCreatedAt(t *time.Time) *LeaveApprovalCreate {
	if t != nil {
		lac.SetCreatedAt(*t)
	}
	return lac
}

// SetUpdatedAt sets the "updated_at" field.
func (lac *LeaveApprovalCreate) SetUpdatedAt(t time.Time) *LeaveApprovalCreate {
	lac.mutation.SetUpdatedAt(t)
	return lac
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (lac *LeaveApprovalCreate) SetNillableUpdatedAt(t *time.Time) *LeaveApprovalCreate {
	if t != nil {
		lac.SetUpdatedAt(*t)
	}
	return lac
}

// SetLeaveRequest sets the "leave_request" edge to the LeaveRequest entity.
func (lac *LeaveApprovalCreate) SetLeaveRequest(l *LeaveRequest) *LeaveApprovalCreate {
	return lac.SetLeaveRequestID(l.ID)
}

// SetReviewer sets the "reviewer" edge to the Employee entity.
func (lac *LeaveApprovalCreate) SetReviewer(e *Employee) *LeaveApprovalCreate {
	return lac.SetReviewerID(e.ID)
}

// Mutation returns the LeaveApprovalMutation object of the builder.
func (lac *LeaveApprovalCreate) Mutation() *LeaveApprovalMutation {
	return lac.mutation
}

// Save creates the LeaveApproval in the database.
func (lac *LeaveApprovalCreate) Save(ctx context.Context) (*LeaveApproval, error) {
	lac.defaults()
	return withHooks(ctx, lac.sqlSave, lac.mutation, lac.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (lac *LeaveApprovalCreate) SaveX(ctx context.Context) *LeaveApproval {
	v, err := lac.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (lac *LeaveApprovalCreate) Exec(ctx context.Context) error {
	_, err := lac.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (lac *LeaveApprovalCreate) ExecX(ctx context.Context) {
	if err := lac.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (lac *LeaveApprovalCreate) defaults() {
	if _, ok := lac.mutation.CreatedAt(); !ok {
		v := leaveapproval.DefaultCreatedAt()
		lac.mutation.SetCreatedAt(v)
	}
	if _, ok := lac.mutation.UpdatedAt(); !ok {
		v := leaveapproval.DefaultUpdatedAt()
		lac.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (lac *LeaveApprovalCreate) check() error {
	if _, ok := lac.mutation.LeaveRequestID(); !ok {
		return &ValidationError{Name: "leave_request_id", err: errors.New(`ent: missing required field "LeaveApproval.leave_request_id"`)}
	}
	if _, ok := lac.mutation.ReviewerID(); !ok {
		return &ValidationError{Name: "reviewer_id", err: errors.New(`ent: missing required field "LeaveApproval.reviewer_id"`)}
	}
	if _, ok := lac.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "LeaveApproval.created_at"`)}
	}
	if _, ok := lac.mutation.UpdatedAt(); !ok {
		return &ValidationError{Name: "updated_at", err: errors.New(`ent: missing required field "LeaveApproval.updated_at"`)}
	}
	if len(lac.mutation.LeaveRequestIDs()) == 0 {
		return &ValidationError{Name: "leave_request", err: errors.New(`ent: missing required edge "LeaveApproval.leave_request"`)}
	}
	if len(lac.mutation.ReviewerIDs()) == 0 {
		return &ValidationError{Name: "reviewer", err: errors.New(`ent: missing required edge "LeaveApproval.reviewer"`)}
	}
	return nil
}

func (lac *LeaveApprovalCreate) sqlSave(ctx context.Context) (*LeaveApproval, error) {
	if err := lac.check(); err != nil {
		return nil, err
	}
	_node, _spec := lac.createSpec()
	if err := sqlgraph.CreateNode(ctx, lac.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	lac.mutation.id = &_node.ID
	lac.mutation.done = true
	return _node, nil
}

func (lac *LeaveApprovalCreate) createSpec() (*LeaveApproval, *sqlgraph.CreateSpec) {
	var (
		_node = &LeaveApproval{config: lac.config}
		_spec = sqlgraph.NewCreateSpec(leaveapproval.Table, sqlgraph.NewFieldSpec(leaveapproval.FieldID, field.TypeInt))
	)
	_spec.OnConflict = lac.conflict
	if value, ok := lac.mutation.Comment(); ok {
		_spec.SetField(leaveapproval.FieldComment, field.TypeString, value)
		_node.Comment = &value
	}
	if value, ok := lac.mutation.CreatedAt(); ok {
		_spec.SetField(leaveapproval.FieldCreatedAt, field.TypeTime, value)
		_node.CreatedAt = value
	}
	if value, ok := lac.mutation.UpdatedAt(); ok {
		_spec.SetField(leaveapproval.FieldUpdatedAt, field.TypeTime, value)
		_node.UpdatedAt = value
	}
	if nodes := lac.mutation.LeaveRequestIDs(); len(nodes) > 0 {
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
		_node.LeaveRequestID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := lac.mutation.ReviewerIDs(); len(nodes) > 0 {
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
		_node.ReviewerID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.LeaveApproval.Create().
//		SetComment(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.LeaveApprovalUpsert) {
//			SetComment(v+v).
//		}).
//		Exec(ctx)
func (lac *LeaveApprovalCreate) OnConflict(opts ...sql.ConflictOption) *LeaveApprovalUpsertOne {
	lac.conflict = opts
	return &LeaveApprovalUpsertOne{
		create: lac,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.LeaveApproval.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (lac *LeaveApprovalCreate) OnConflictColumns(columns ...string) *LeaveApprovalUpsertOne {
	lac.conflict = append(lac.conflict, sql.ConflictColumns(columns...))
	return &LeaveApprovalUpsertOne{
		create: lac,
	}
}

type (
	// LeaveApprovalUpsertOne is the builder for "upsert"-ing
	//  one LeaveApproval node.
	LeaveApprovalUpsertOne struct {
		create *LeaveApprovalCreate
	}

	// LeaveApprovalUpsert is the "OnConflict" setter.
	LeaveApprovalUpsert struct {
		*sql.UpdateSet
	}
)

// SetComment sets the "comment" field.
func (u *LeaveApprovalUpsert) SetComment(v string) *LeaveApprovalUpsert {
	u.Set(leaveapproval.FieldComment, v)
	return u
}

// UpdateComment sets the "comment" field to the value that was provided on create.
func (u *LeaveApprovalUpsert) UpdateComment() *LeaveApprovalUpsert {
	u.SetExcluded(leaveapproval.FieldComment)
	return u
}

// ClearComment clears the value of the "comment" field.
func (u *LeaveApprovalUpsert) ClearComment() *LeaveApprovalUpsert {
	u.SetNull(leaveapproval.FieldComment)
	return u
}

// SetLeaveRequestID sets the "leave_request_id" field.
func (u *LeaveApprovalUpsert) SetLeaveRequestID(v int) *LeaveApprovalUpsert {
	u.Set(leaveapproval.FieldLeaveRequestID, v)
	return u
}

// UpdateLeaveRequestID sets the "leave_request_id" field to the value that was provided on create.
func (u *LeaveApprovalUpsert) UpdateLeaveRequestID() *LeaveApprovalUpsert {
	u.SetExcluded(leaveapproval.FieldLeaveRequestID)
	return u
}

// SetReviewerID sets the "reviewer_id" field.
func (u *LeaveApprovalUpsert) SetReviewerID(v int) *LeaveApprovalUpsert {
	u.Set(leaveapproval.FieldReviewerID, v)
	return u
}

// UpdateReviewerID sets the "reviewer_id" field to the value that was provided on create.
func (u *LeaveApprovalUpsert) UpdateReviewerID() *LeaveApprovalUpsert {
	u.SetExcluded(leaveapproval.FieldReviewerID)
	return u
}

// SetCreatedAt sets the "created_at" field.
func (u *LeaveApprovalUpsert) SetCreatedAt(v time.Time) *LeaveApprovalUpsert {
	u.Set(leaveapproval.FieldCreatedAt, v)
	return u
}

// UpdateCreatedAt sets the "created_at" field to the value that was provided on create.
func (u *LeaveApprovalUpsert) UpdateCreatedAt() *LeaveApprovalUpsert {
	u.SetExcluded(leaveapproval.FieldCreatedAt)
	return u
}

// SetUpdatedAt sets the "updated_at" field.
func (u *LeaveApprovalUpsert) SetUpdatedAt(v time.Time) *LeaveApprovalUpsert {
	u.Set(leaveapproval.FieldUpdatedAt, v)
	return u
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *LeaveApprovalUpsert) UpdateUpdatedAt() *LeaveApprovalUpsert {
	u.SetExcluded(leaveapproval.FieldUpdatedAt)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create.
// Using this option is equivalent to using:
//
//	client.LeaveApproval.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *LeaveApprovalUpsertOne) UpdateNewValues() *LeaveApprovalUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.LeaveApproval.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *LeaveApprovalUpsertOne) Ignore() *LeaveApprovalUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *LeaveApprovalUpsertOne) DoNothing() *LeaveApprovalUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the LeaveApprovalCreate.OnConflict
// documentation for more info.
func (u *LeaveApprovalUpsertOne) Update(set func(*LeaveApprovalUpsert)) *LeaveApprovalUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&LeaveApprovalUpsert{UpdateSet: update})
	}))
	return u
}

// SetComment sets the "comment" field.
func (u *LeaveApprovalUpsertOne) SetComment(v string) *LeaveApprovalUpsertOne {
	return u.Update(func(s *LeaveApprovalUpsert) {
		s.SetComment(v)
	})
}

// UpdateComment sets the "comment" field to the value that was provided on create.
func (u *LeaveApprovalUpsertOne) UpdateComment() *LeaveApprovalUpsertOne {
	return u.Update(func(s *LeaveApprovalUpsert) {
		s.UpdateComment()
	})
}

// ClearComment clears the value of the "comment" field.
func (u *LeaveApprovalUpsertOne) ClearComment() *LeaveApprovalUpsertOne {
	return u.Update(func(s *LeaveApprovalUpsert) {
		s.ClearComment()
	})
}

// SetLeaveRequestID sets the "leave_request_id" field.
func (u *LeaveApprovalUpsertOne) SetLeaveRequestID(v int) *LeaveApprovalUpsertOne {
	return u.Update(func(s *LeaveApprovalUpsert) {
		s.SetLeaveRequestID(v)
	})
}

// UpdateLeaveRequestID sets the "leave_request_id" field to the value that was provided on create.
func (u *LeaveApprovalUpsertOne) UpdateLeaveRequestID() *LeaveApprovalUpsertOne {
	return u.Update(func(s *LeaveApprovalUpsert) {
		s.UpdateLeaveRequestID()
	})
}

// SetReviewerID sets the "reviewer_id" field.
func (u *LeaveApprovalUpsertOne) SetReviewerID(v int) *LeaveApprovalUpsertOne {
	return u.Update(func(s *LeaveApprovalUpsert) {
		s.SetReviewerID(v)
	})
}

// UpdateReviewerID sets the "reviewer_id" field to the value that was provided on create.
func (u *LeaveApprovalUpsertOne) UpdateReviewerID() *LeaveApprovalUpsertOne {
	return u.Update(func(s *LeaveApprovalUpsert) {
		s.UpdateReviewerID()
	})
}

// SetCreatedAt sets the "created_at" field.
func (u *LeaveApprovalUpsertOne) SetCreatedAt(v time.Time) *LeaveApprovalUpsertOne {
	return u.Update(func(s *LeaveApprovalUpsert) {
		s.SetCreatedAt(v)
	})
}

// UpdateCreatedAt sets the "created_at" field to the value that was provided on create.
func (u *LeaveApprovalUpsertOne) UpdateCreatedAt() *LeaveApprovalUpsertOne {
	return u.Update(func(s *LeaveApprovalUpsert) {
		s.UpdateCreatedAt()
	})
}

// SetUpdatedAt sets the "updated_at" field.
func (u *LeaveApprovalUpsertOne) SetUpdatedAt(v time.Time) *LeaveApprovalUpsertOne {
	return u.Update(func(s *LeaveApprovalUpsert) {
		s.SetUpdatedAt(v)
	})
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *LeaveApprovalUpsertOne) UpdateUpdatedAt() *LeaveApprovalUpsertOne {
	return u.Update(func(s *LeaveApprovalUpsert) {
		s.UpdateUpdatedAt()
	})
}

// Exec executes the query.
func (u *LeaveApprovalUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for LeaveApprovalCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *LeaveApprovalUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *LeaveApprovalUpsertOne) ID(ctx context.Context) (id int, err error) {
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *LeaveApprovalUpsertOne) IDX(ctx context.Context) int {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// LeaveApprovalCreateBulk is the builder for creating many LeaveApproval entities in bulk.
type LeaveApprovalCreateBulk struct {
	config
	err      error
	builders []*LeaveApprovalCreate
	conflict []sql.ConflictOption
}

// Save creates the LeaveApproval entities in the database.
func (lacb *LeaveApprovalCreateBulk) Save(ctx context.Context) ([]*LeaveApproval, error) {
	if lacb.err != nil {
		return nil, lacb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(lacb.builders))
	nodes := make([]*LeaveApproval, len(lacb.builders))
	mutators := make([]Mutator, len(lacb.builders))
	for i := range lacb.builders {
		func(i int, root context.Context) {
			builder := lacb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*LeaveApprovalMutation)
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
					_, err = mutators[i+1].Mutate(root, lacb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = lacb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, lacb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, lacb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (lacb *LeaveApprovalCreateBulk) SaveX(ctx context.Context) []*LeaveApproval {
	v, err := lacb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (lacb *LeaveApprovalCreateBulk) Exec(ctx context.Context) error {
	_, err := lacb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (lacb *LeaveApprovalCreateBulk) ExecX(ctx context.Context) {
	if err := lacb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.LeaveApproval.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.LeaveApprovalUpsert) {
//			SetComment(v+v).
//		}).
//		Exec(ctx)
func (lacb *LeaveApprovalCreateBulk) OnConflict(opts ...sql.ConflictOption) *LeaveApprovalUpsertBulk {
	lacb.conflict = opts
	return &LeaveApprovalUpsertBulk{
		create: lacb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.LeaveApproval.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (lacb *LeaveApprovalCreateBulk) OnConflictColumns(columns ...string) *LeaveApprovalUpsertBulk {
	lacb.conflict = append(lacb.conflict, sql.ConflictColumns(columns...))
	return &LeaveApprovalUpsertBulk{
		create: lacb,
	}
}

// LeaveApprovalUpsertBulk is the builder for "upsert"-ing
// a bulk of LeaveApproval nodes.
type LeaveApprovalUpsertBulk struct {
	create *LeaveApprovalCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.LeaveApproval.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *LeaveApprovalUpsertBulk) UpdateNewValues() *LeaveApprovalUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.LeaveApproval.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *LeaveApprovalUpsertBulk) Ignore() *LeaveApprovalUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *LeaveApprovalUpsertBulk) DoNothing() *LeaveApprovalUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the LeaveApprovalCreateBulk.OnConflict
// documentation for more info.
func (u *LeaveApprovalUpsertBulk) Update(set func(*LeaveApprovalUpsert)) *LeaveApprovalUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&LeaveApprovalUpsert{UpdateSet: update})
	}))
	return u
}

// SetComment sets the "comment" field.
func (u *LeaveApprovalUpsertBulk) SetComment(v string) *LeaveApprovalUpsertBulk {
	return u.Update(func(s *LeaveApprovalUpsert) {
		s.SetComment(v)
	})
}

// UpdateComment sets the "comment" field to the value that was provided on create.
func (u *LeaveApprovalUpsertBulk) UpdateComment() *LeaveApprovalUpsertBulk {
	return u.Update(func(s *LeaveApprovalUpsert) {
		s.UpdateComment()
	})
}

// ClearComment clears the value of the "comment" field.
func (u *LeaveApprovalUpsertBulk) ClearComment() *LeaveApprovalUpsertBulk {
	return u.Update(func(s *LeaveApprovalUpsert) {
		s.ClearComment()
	})
}

// SetLeaveRequestID sets the "leave_request_id" field.
func (u *LeaveApprovalUpsertBulk) SetLeaveRequestID(v int) *LeaveApprovalUpsertBulk {
	return u.Update(func(s *LeaveApprovalUpsert) {
		s.SetLeaveRequestID(v)
	})
}

// UpdateLeaveRequestID sets the "leave_request_id" field to the value that was provided on create.
func (u *LeaveApprovalUpsertBulk) UpdateLeaveRequestID() *LeaveApprovalUpsertBulk {
	return u.Update(func(s *LeaveApprovalUpsert) {
		s.UpdateLeaveRequestID()
	})
}

// SetReviewerID sets the "reviewer_id" field.
func (u *LeaveApprovalUpsertBulk) SetReviewerID(v int) *LeaveApprovalUpsertBulk {
	return u.Update(func(s *LeaveApprovalUpsert) {
		s.SetReviewerID(v)
	})
}

// UpdateReviewerID sets the "reviewer_id" field to the value that was provided on create.
func (u *LeaveApprovalUpsertBulk) UpdateReviewerID() *LeaveApprovalUpsertBulk {
	return u.Update(func(s *LeaveApprovalUpsert) {
		s.UpdateReviewerID()
	})
}

// SetCreatedAt sets the "created_at" field.
func (u *LeaveApprovalUpsertBulk) SetCreatedAt(v time.Time) *LeaveApprovalUpsertBulk {
	return u.Update(func(s *LeaveApprovalUpsert) {
		s.SetCreatedAt(v)
	})
}

// UpdateCreatedAt sets the "created_at" field to the value that was provided on create.
func (u *LeaveApprovalUpsertBulk) UpdateCreatedAt() *LeaveApprovalUpsertBulk {
	return u.Update(func(s *LeaveApprovalUpsert) {
		s.UpdateCreatedAt()
	})
}

// SetUpdatedAt sets the "updated_at" field.
func (u *LeaveApprovalUpsertBulk) SetUpdatedAt(v time.Time) *LeaveApprovalUpsertBulk {
	return u.Update(func(s *LeaveApprovalUpsert) {
		s.SetUpdatedAt(v)
	})
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *LeaveApprovalUpsertBulk) UpdateUpdatedAt() *LeaveApprovalUpsertBulk {
	return u.Update(func(s *LeaveApprovalUpsert) {
		s.UpdateUpdatedAt()
	})
}

// Exec executes the query.
func (u *LeaveApprovalUpsertBulk) Exec(ctx context.Context) error {
	if u.create.err != nil {
		return u.create.err
	}
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the LeaveApprovalCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for LeaveApprovalCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *LeaveApprovalUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
