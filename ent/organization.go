// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/longgggwwww/hrm-ms-hr/ent/organization"
)

// Organization is the model entity for the Organization schema.
type Organization struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name"`
	// Code holds the value of the "code" field.
	Code string `json:"code"`
	// LogoURL holds the value of the "logo_url" field.
	LogoURL *string `json:"logo_url"`
	// Address holds the value of the "address" field.
	Address *string `json:"address"`
	// Phone holds the value of the "phone" field.
	Phone *string `json:"phone"`
	// Email holds the value of the "email" field.
	Email *string `json:"email"`
	// Website holds the value of the "website" field.
	Website *string `json:"website"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updated_at"`
	// ParentID holds the value of the "parent_id" field.
	ParentID *int `json:"parent_id"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the OrganizationQuery when eager-loading is set.
	Edges        OrganizationEdges `json:"edges"`
	selectValues sql.SelectValues
}

// OrganizationEdges holds the relations/edges for other nodes in the graph.
type OrganizationEdges struct {
	// Parent holds the value of the parent edge.
	Parent *Organization `json:"parent"`
	// Children holds the value of the children edge.
	Children []*Organization `json:"children"`
	// Departments holds the value of the departments edge.
	Departments []*Department `json:"departments"`
	// Projects holds the value of the projects edge.
	Projects []*Project `json:"projects"`
	// Labels holds the value of the labels edge.
	Labels []*Label `json:"labels"`
	// LeaveRequests holds the value of the leave_requests edge.
	LeaveRequests []*LeaveRequest `json:"leave_requests"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [6]bool
}

// ParentOrErr returns the Parent value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e OrganizationEdges) ParentOrErr() (*Organization, error) {
	if e.Parent != nil {
		return e.Parent, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: organization.Label}
	}
	return nil, &NotLoadedError{edge: "parent"}
}

// ChildrenOrErr returns the Children value or an error if the edge
// was not loaded in eager-loading.
func (e OrganizationEdges) ChildrenOrErr() ([]*Organization, error) {
	if e.loadedTypes[1] {
		return e.Children, nil
	}
	return nil, &NotLoadedError{edge: "children"}
}

// DepartmentsOrErr returns the Departments value or an error if the edge
// was not loaded in eager-loading.
func (e OrganizationEdges) DepartmentsOrErr() ([]*Department, error) {
	if e.loadedTypes[2] {
		return e.Departments, nil
	}
	return nil, &NotLoadedError{edge: "departments"}
}

// ProjectsOrErr returns the Projects value or an error if the edge
// was not loaded in eager-loading.
func (e OrganizationEdges) ProjectsOrErr() ([]*Project, error) {
	if e.loadedTypes[3] {
		return e.Projects, nil
	}
	return nil, &NotLoadedError{edge: "projects"}
}

// LabelsOrErr returns the Labels value or an error if the edge
// was not loaded in eager-loading.
func (e OrganizationEdges) LabelsOrErr() ([]*Label, error) {
	if e.loadedTypes[4] {
		return e.Labels, nil
	}
	return nil, &NotLoadedError{edge: "labels"}
}

// LeaveRequestsOrErr returns the LeaveRequests value or an error if the edge
// was not loaded in eager-loading.
func (e OrganizationEdges) LeaveRequestsOrErr() ([]*LeaveRequest, error) {
	if e.loadedTypes[5] {
		return e.LeaveRequests, nil
	}
	return nil, &NotLoadedError{edge: "leave_requests"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Organization) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case organization.FieldID, organization.FieldParentID:
			values[i] = new(sql.NullInt64)
		case organization.FieldName, organization.FieldCode, organization.FieldLogoURL, organization.FieldAddress, organization.FieldPhone, organization.FieldEmail, organization.FieldWebsite:
			values[i] = new(sql.NullString)
		case organization.FieldCreatedAt, organization.FieldUpdatedAt:
			values[i] = new(sql.NullTime)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Organization fields.
func (o *Organization) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case organization.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			o.ID = int(value.Int64)
		case organization.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				o.Name = value.String
			}
		case organization.FieldCode:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field code", values[i])
			} else if value.Valid {
				o.Code = value.String
			}
		case organization.FieldLogoURL:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field logo_url", values[i])
			} else if value.Valid {
				o.LogoURL = new(string)
				*o.LogoURL = value.String
			}
		case organization.FieldAddress:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field address", values[i])
			} else if value.Valid {
				o.Address = new(string)
				*o.Address = value.String
			}
		case organization.FieldPhone:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field phone", values[i])
			} else if value.Valid {
				o.Phone = new(string)
				*o.Phone = value.String
			}
		case organization.FieldEmail:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field email", values[i])
			} else if value.Valid {
				o.Email = new(string)
				*o.Email = value.String
			}
		case organization.FieldWebsite:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field website", values[i])
			} else if value.Valid {
				o.Website = new(string)
				*o.Website = value.String
			}
		case organization.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				o.CreatedAt = value.Time
			}
		case organization.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				o.UpdatedAt = value.Time
			}
		case organization.FieldParentID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field parent_id", values[i])
			} else if value.Valid {
				o.ParentID = new(int)
				*o.ParentID = int(value.Int64)
			}
		default:
			o.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Organization.
// This includes values selected through modifiers, order, etc.
func (o *Organization) Value(name string) (ent.Value, error) {
	return o.selectValues.Get(name)
}

// QueryParent queries the "parent" edge of the Organization entity.
func (o *Organization) QueryParent() *OrganizationQuery {
	return NewOrganizationClient(o.config).QueryParent(o)
}

// QueryChildren queries the "children" edge of the Organization entity.
func (o *Organization) QueryChildren() *OrganizationQuery {
	return NewOrganizationClient(o.config).QueryChildren(o)
}

// QueryDepartments queries the "departments" edge of the Organization entity.
func (o *Organization) QueryDepartments() *DepartmentQuery {
	return NewOrganizationClient(o.config).QueryDepartments(o)
}

// QueryProjects queries the "projects" edge of the Organization entity.
func (o *Organization) QueryProjects() *ProjectQuery {
	return NewOrganizationClient(o.config).QueryProjects(o)
}

// QueryLabels queries the "labels" edge of the Organization entity.
func (o *Organization) QueryLabels() *LabelQuery {
	return NewOrganizationClient(o.config).QueryLabels(o)
}

// QueryLeaveRequests queries the "leave_requests" edge of the Organization entity.
func (o *Organization) QueryLeaveRequests() *LeaveRequestQuery {
	return NewOrganizationClient(o.config).QueryLeaveRequests(o)
}

// Update returns a builder for updating this Organization.
// Note that you need to call Organization.Unwrap() before calling this method if this Organization
// was returned from a transaction, and the transaction was committed or rolled back.
func (o *Organization) Update() *OrganizationUpdateOne {
	return NewOrganizationClient(o.config).UpdateOne(o)
}

// Unwrap unwraps the Organization entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (o *Organization) Unwrap() *Organization {
	_tx, ok := o.config.driver.(*txDriver)
	if !ok {
		panic("ent: Organization is not a transactional entity")
	}
	o.config.driver = _tx.drv
	return o
}

// String implements the fmt.Stringer.
func (o *Organization) String() string {
	var builder strings.Builder
	builder.WriteString("Organization(")
	builder.WriteString(fmt.Sprintf("id=%v, ", o.ID))
	builder.WriteString("name=")
	builder.WriteString(o.Name)
	builder.WriteString(", ")
	builder.WriteString("code=")
	builder.WriteString(o.Code)
	builder.WriteString(", ")
	if v := o.LogoURL; v != nil {
		builder.WriteString("logo_url=")
		builder.WriteString(*v)
	}
	builder.WriteString(", ")
	if v := o.Address; v != nil {
		builder.WriteString("address=")
		builder.WriteString(*v)
	}
	builder.WriteString(", ")
	if v := o.Phone; v != nil {
		builder.WriteString("phone=")
		builder.WriteString(*v)
	}
	builder.WriteString(", ")
	if v := o.Email; v != nil {
		builder.WriteString("email=")
		builder.WriteString(*v)
	}
	builder.WriteString(", ")
	if v := o.Website; v != nil {
		builder.WriteString("website=")
		builder.WriteString(*v)
	}
	builder.WriteString(", ")
	builder.WriteString("created_at=")
	builder.WriteString(o.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(o.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	if v := o.ParentID; v != nil {
		builder.WriteString("parent_id=")
		builder.WriteString(fmt.Sprintf("%v", *v))
	}
	builder.WriteByte(')')
	return builder.String()
}

// Organizations is a parsable slice of Organization.
type Organizations []*Organization
