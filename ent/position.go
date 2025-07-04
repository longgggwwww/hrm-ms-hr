// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/longgggwwww/hrm-ms-hr/ent/department"
	"github.com/longgggwwww/hrm-ms-hr/ent/position"
)

// Position is the model entity for the Position schema.
type Position struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name"`
	// Code holds the value of the "code" field.
	Code string `json:"code"`
	// DepartmentID holds the value of the "department_id" field.
	DepartmentID int `json:"department_id"`
	// ParentID holds the value of the "parent_id" field.
	ParentID int `json:"parent_id"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updated_at"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the PositionQuery when eager-loading is set.
	Edges        PositionEdges `json:"edges"`
	selectValues sql.SelectValues
}

// PositionEdges holds the relations/edges for other nodes in the graph.
type PositionEdges struct {
	// Employees holds the value of the employees edge.
	Employees []*Employee `json:"employees"`
	// Department holds the value of the department edge.
	Department *Department `json:"department"`
	// Children holds the value of the children edge.
	Children []*Position `json:"children"`
	// Parent holds the value of the parent edge.
	Parent *Position `json:"parent"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [4]bool
}

// EmployeesOrErr returns the Employees value or an error if the edge
// was not loaded in eager-loading.
func (e PositionEdges) EmployeesOrErr() ([]*Employee, error) {
	if e.loadedTypes[0] {
		return e.Employees, nil
	}
	return nil, &NotLoadedError{edge: "employees"}
}

// DepartmentOrErr returns the Department value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e PositionEdges) DepartmentOrErr() (*Department, error) {
	if e.Department != nil {
		return e.Department, nil
	} else if e.loadedTypes[1] {
		return nil, &NotFoundError{label: department.Label}
	}
	return nil, &NotLoadedError{edge: "department"}
}

// ChildrenOrErr returns the Children value or an error if the edge
// was not loaded in eager-loading.
func (e PositionEdges) ChildrenOrErr() ([]*Position, error) {
	if e.loadedTypes[2] {
		return e.Children, nil
	}
	return nil, &NotLoadedError{edge: "children"}
}

// ParentOrErr returns the Parent value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e PositionEdges) ParentOrErr() (*Position, error) {
	if e.Parent != nil {
		return e.Parent, nil
	} else if e.loadedTypes[3] {
		return nil, &NotFoundError{label: position.Label}
	}
	return nil, &NotLoadedError{edge: "parent"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Position) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case position.FieldID, position.FieldDepartmentID, position.FieldParentID:
			values[i] = new(sql.NullInt64)
		case position.FieldName, position.FieldCode:
			values[i] = new(sql.NullString)
		case position.FieldCreatedAt, position.FieldUpdatedAt:
			values[i] = new(sql.NullTime)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Position fields.
func (po *Position) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case position.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			po.ID = int(value.Int64)
		case position.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				po.Name = value.String
			}
		case position.FieldCode:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field code", values[i])
			} else if value.Valid {
				po.Code = value.String
			}
		case position.FieldDepartmentID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field department_id", values[i])
			} else if value.Valid {
				po.DepartmentID = int(value.Int64)
			}
		case position.FieldParentID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field parent_id", values[i])
			} else if value.Valid {
				po.ParentID = int(value.Int64)
			}
		case position.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				po.CreatedAt = value.Time
			}
		case position.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				po.UpdatedAt = value.Time
			}
		default:
			po.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Position.
// This includes values selected through modifiers, order, etc.
func (po *Position) Value(name string) (ent.Value, error) {
	return po.selectValues.Get(name)
}

// QueryEmployees queries the "employees" edge of the Position entity.
func (po *Position) QueryEmployees() *EmployeeQuery {
	return NewPositionClient(po.config).QueryEmployees(po)
}

// QueryDepartment queries the "department" edge of the Position entity.
func (po *Position) QueryDepartment() *DepartmentQuery {
	return NewPositionClient(po.config).QueryDepartment(po)
}

// QueryChildren queries the "children" edge of the Position entity.
func (po *Position) QueryChildren() *PositionQuery {
	return NewPositionClient(po.config).QueryChildren(po)
}

// QueryParent queries the "parent" edge of the Position entity.
func (po *Position) QueryParent() *PositionQuery {
	return NewPositionClient(po.config).QueryParent(po)
}

// Update returns a builder for updating this Position.
// Note that you need to call Position.Unwrap() before calling this method if this Position
// was returned from a transaction, and the transaction was committed or rolled back.
func (po *Position) Update() *PositionUpdateOne {
	return NewPositionClient(po.config).UpdateOne(po)
}

// Unwrap unwraps the Position entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (po *Position) Unwrap() *Position {
	_tx, ok := po.config.driver.(*txDriver)
	if !ok {
		panic("ent: Position is not a transactional entity")
	}
	po.config.driver = _tx.drv
	return po
}

// String implements the fmt.Stringer.
func (po *Position) String() string {
	var builder strings.Builder
	builder.WriteString("Position(")
	builder.WriteString(fmt.Sprintf("id=%v, ", po.ID))
	builder.WriteString("name=")
	builder.WriteString(po.Name)
	builder.WriteString(", ")
	builder.WriteString("code=")
	builder.WriteString(po.Code)
	builder.WriteString(", ")
	builder.WriteString("department_id=")
	builder.WriteString(fmt.Sprintf("%v", po.DepartmentID))
	builder.WriteString(", ")
	builder.WriteString("parent_id=")
	builder.WriteString(fmt.Sprintf("%v", po.ParentID))
	builder.WriteString(", ")
	builder.WriteString("created_at=")
	builder.WriteString(po.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(po.UpdatedAt.Format(time.ANSIC))
	builder.WriteByte(')')
	return builder.String()
}

// Positions is a parsable slice of Position.
type Positions []*Position
