// Code generated by ent, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// AppointmentHistoriesColumns holds the columns for the "appointment_histories" table.
	AppointmentHistoriesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "position_name", Type: field.TypeString},
		{Name: "joining_at", Type: field.TypeTime},
		{Name: "description", Type: field.TypeString, Nullable: true},
		{Name: "attachment_urls", Type: field.TypeJSON, Nullable: true},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "employee_id", Type: field.TypeInt},
	}
	// AppointmentHistoriesTable holds the schema information for the "appointment_histories" table.
	AppointmentHistoriesTable = &schema.Table{
		Name:       "appointment_histories",
		Columns:    AppointmentHistoriesColumns,
		PrimaryKey: []*schema.Column{AppointmentHistoriesColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "appointment_histories_employees_appointment_histories",
				Columns:    []*schema.Column{AppointmentHistoriesColumns[7]},
				RefColumns: []*schema.Column{EmployeesColumns[0]},
				OnDelete:   schema.NoAction,
			},
		},
	}
	// DepartmentsColumns holds the columns for the "departments" table.
	DepartmentsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "name", Type: field.TypeString},
		{Name: "code", Type: field.TypeString},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "org_id", Type: field.TypeInt},
	}
	// DepartmentsTable holds the schema information for the "departments" table.
	DepartmentsTable = &schema.Table{
		Name:       "departments",
		Columns:    DepartmentsColumns,
		PrimaryKey: []*schema.Column{DepartmentsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "departments_organizations_departments",
				Columns:    []*schema.Column{DepartmentsColumns[5]},
				RefColumns: []*schema.Column{OrganizationsColumns[0]},
				OnDelete:   schema.NoAction,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "department_code_org_id",
				Unique:  true,
				Columns: []*schema.Column{DepartmentsColumns[2], DepartmentsColumns[5]},
			},
		},
	}
	// EmployeesColumns holds the columns for the "employees" table.
	EmployeesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "user_id", Type: field.TypeString, Unique: true, Nullable: true},
		{Name: "code", Type: field.TypeString, Unique: true},
		{Name: "status", Type: field.TypeEnum, Enums: []string{"active", "inactive"}, Default: "active"},
		{Name: "joining_at", Type: field.TypeTime},
		{Name: "org_id", Type: field.TypeInt},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "position_id", Type: field.TypeInt},
	}
	// EmployeesTable holds the schema information for the "employees" table.
	EmployeesTable = &schema.Table{
		Name:       "employees",
		Columns:    EmployeesColumns,
		PrimaryKey: []*schema.Column{EmployeesColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "employees_positions_employees",
				Columns:    []*schema.Column{EmployeesColumns[8]},
				RefColumns: []*schema.Column{PositionsColumns[0]},
				OnDelete:   schema.NoAction,
			},
		},
	}
	// LabelsColumns holds the columns for the "labels" table.
	LabelsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "name", Type: field.TypeString},
		{Name: "description", Type: field.TypeString, Nullable: true},
		{Name: "color", Type: field.TypeString},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "org_id", Type: field.TypeInt, Nullable: true},
	}
	// LabelsTable holds the schema information for the "labels" table.
	LabelsTable = &schema.Table{
		Name:       "labels",
		Columns:    LabelsColumns,
		PrimaryKey: []*schema.Column{LabelsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "labels_organizations_labels",
				Columns:    []*schema.Column{LabelsColumns[6]},
				RefColumns: []*schema.Column{OrganizationsColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "label_org_id",
				Unique:  false,
				Columns: []*schema.Column{LabelsColumns[6]},
			},
		},
	}
	// LeaveApprovalsColumns holds the columns for the "leave_approvals" table.
	LeaveApprovalsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "comment", Type: field.TypeString, Nullable: true},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "reviewer_id", Type: field.TypeInt},
		{Name: "leave_request_id", Type: field.TypeInt},
	}
	// LeaveApprovalsTable holds the schema information for the "leave_approvals" table.
	LeaveApprovalsTable = &schema.Table{
		Name:       "leave_approvals",
		Columns:    LeaveApprovalsColumns,
		PrimaryKey: []*schema.Column{LeaveApprovalsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "leave_approvals_employees_leave_approves",
				Columns:    []*schema.Column{LeaveApprovalsColumns[4]},
				RefColumns: []*schema.Column{EmployeesColumns[0]},
				OnDelete:   schema.NoAction,
			},
			{
				Symbol:     "leave_approvals_leave_requests_leave_approves",
				Columns:    []*schema.Column{LeaveApprovalsColumns[5]},
				RefColumns: []*schema.Column{LeaveRequestsColumns[0]},
				OnDelete:   schema.NoAction,
			},
		},
	}
	// LeaveRequestsColumns holds the columns for the "leave_requests" table.
	LeaveRequestsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "total_days", Type: field.TypeFloat64},
		{Name: "start_at", Type: field.TypeTime},
		{Name: "end_at", Type: field.TypeTime},
		{Name: "reason", Type: field.TypeString, Nullable: true},
		{Name: "type", Type: field.TypeEnum, Enums: []string{"annual", "unpaid"}, Default: "annual"},
		{Name: "status", Type: field.TypeEnum, Enums: []string{"pending", "rejected", "approved"}, Default: "pending"},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "employee_id", Type: field.TypeInt},
		{Name: "org_id", Type: field.TypeInt},
	}
	// LeaveRequestsTable holds the schema information for the "leave_requests" table.
	LeaveRequestsTable = &schema.Table{
		Name:       "leave_requests",
		Columns:    LeaveRequestsColumns,
		PrimaryKey: []*schema.Column{LeaveRequestsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "leave_requests_employees_leave_requests",
				Columns:    []*schema.Column{LeaveRequestsColumns[9]},
				RefColumns: []*schema.Column{EmployeesColumns[0]},
				OnDelete:   schema.NoAction,
			},
			{
				Symbol:     "leave_requests_organizations_leave_requests",
				Columns:    []*schema.Column{LeaveRequestsColumns[10]},
				RefColumns: []*schema.Column{OrganizationsColumns[0]},
				OnDelete:   schema.NoAction,
			},
		},
	}
	// OrganizationsColumns holds the columns for the "organizations" table.
	OrganizationsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "name", Type: field.TypeString},
		{Name: "code", Type: field.TypeString, Unique: true},
		{Name: "logo_url", Type: field.TypeString, Nullable: true},
		{Name: "address", Type: field.TypeString, Nullable: true},
		{Name: "phone", Type: field.TypeString, Nullable: true},
		{Name: "email", Type: field.TypeString, Nullable: true},
		{Name: "website", Type: field.TypeString, Nullable: true},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "parent_id", Type: field.TypeInt, Nullable: true},
	}
	// OrganizationsTable holds the schema information for the "organizations" table.
	OrganizationsTable = &schema.Table{
		Name:       "organizations",
		Columns:    OrganizationsColumns,
		PrimaryKey: []*schema.Column{OrganizationsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "organizations_organizations_children",
				Columns:    []*schema.Column{OrganizationsColumns[10]},
				RefColumns: []*schema.Column{OrganizationsColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// PositionsColumns holds the columns for the "positions" table.
	PositionsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "name", Type: field.TypeString},
		{Name: "code", Type: field.TypeString},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "department_id", Type: field.TypeInt},
		{Name: "parent_id", Type: field.TypeInt, Nullable: true},
	}
	// PositionsTable holds the schema information for the "positions" table.
	PositionsTable = &schema.Table{
		Name:       "positions",
		Columns:    PositionsColumns,
		PrimaryKey: []*schema.Column{PositionsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "positions_departments_positions",
				Columns:    []*schema.Column{PositionsColumns[5]},
				RefColumns: []*schema.Column{DepartmentsColumns[0]},
				OnDelete:   schema.NoAction,
			},
			{
				Symbol:     "positions_positions_children",
				Columns:    []*schema.Column{PositionsColumns[6]},
				RefColumns: []*schema.Column{PositionsColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "position_department_id_code",
				Unique:  true,
				Columns: []*schema.Column{PositionsColumns[5], PositionsColumns[2]},
			},
		},
	}
	// ProjectsColumns holds the columns for the "projects" table.
	ProjectsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "name", Type: field.TypeString},
		{Name: "code", Type: field.TypeString},
		{Name: "description", Type: field.TypeString, Nullable: true},
		{Name: "start_at", Type: field.TypeTime, Nullable: true},
		{Name: "end_at", Type: field.TypeTime, Nullable: true},
		{Name: "process", Type: field.TypeInt, Nullable: true},
		{Name: "status", Type: field.TypeEnum, Enums: []string{"not_started", "in_progress", "completed"}, Default: "not_started"},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "creator_id", Type: field.TypeInt},
		{Name: "updater_id", Type: field.TypeInt},
		{Name: "org_id", Type: field.TypeInt},
	}
	// ProjectsTable holds the schema information for the "projects" table.
	ProjectsTable = &schema.Table{
		Name:       "projects",
		Columns:    ProjectsColumns,
		PrimaryKey: []*schema.Column{ProjectsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "projects_employees_created_projects",
				Columns:    []*schema.Column{ProjectsColumns[10]},
				RefColumns: []*schema.Column{EmployeesColumns[0]},
				OnDelete:   schema.NoAction,
			},
			{
				Symbol:     "projects_employees_updated_projects",
				Columns:    []*schema.Column{ProjectsColumns[11]},
				RefColumns: []*schema.Column{EmployeesColumns[0]},
				OnDelete:   schema.NoAction,
			},
			{
				Symbol:     "projects_organizations_projects",
				Columns:    []*schema.Column{ProjectsColumns[12]},
				RefColumns: []*schema.Column{OrganizationsColumns[0]},
				OnDelete:   schema.NoAction,
			},
		},
	}
	// TasksColumns holds the columns for the "tasks" table.
	TasksColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "name", Type: field.TypeString},
		{Name: "code", Type: field.TypeString, Unique: true},
		{Name: "description", Type: field.TypeString, Nullable: true},
		{Name: "process", Type: field.TypeInt, Default: 0},
		{Name: "status", Type: field.TypeEnum, Enums: []string{"not_received", "received", "in_progress", "completed", "cancelled"}, Default: "not_received"},
		{Name: "start_at", Type: field.TypeTime, Nullable: true},
		{Name: "due_date", Type: field.TypeTime, Nullable: true},
		{Name: "creator_id", Type: field.TypeInt},
		{Name: "updater_id", Type: field.TypeInt},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "type", Type: field.TypeEnum, Enums: []string{"task", "feature", "bug", "another"}, Default: "task"},
		{Name: "project_id", Type: field.TypeInt, Nullable: true},
	}
	// TasksTable holds the schema information for the "tasks" table.
	TasksTable = &schema.Table{
		Name:       "tasks",
		Columns:    TasksColumns,
		PrimaryKey: []*schema.Column{TasksColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "tasks_projects_tasks",
				Columns:    []*schema.Column{TasksColumns[13]},
				RefColumns: []*schema.Column{ProjectsColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// TaskReportsColumns holds the columns for the "task_reports" table.
	TaskReportsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "content", Type: field.TypeString, Nullable: true},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "reporter_id", Type: field.TypeInt},
		{Name: "task_id", Type: field.TypeInt},
	}
	// TaskReportsTable holds the schema information for the "task_reports" table.
	TaskReportsTable = &schema.Table{
		Name:       "task_reports",
		Columns:    TaskReportsColumns,
		PrimaryKey: []*schema.Column{TaskReportsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "task_reports_employees_task_reports",
				Columns:    []*schema.Column{TaskReportsColumns[4]},
				RefColumns: []*schema.Column{EmployeesColumns[0]},
				OnDelete:   schema.NoAction,
			},
			{
				Symbol:     "task_reports_tasks_reports",
				Columns:    []*schema.Column{TaskReportsColumns[5]},
				RefColumns: []*schema.Column{TasksColumns[0]},
				OnDelete:   schema.NoAction,
			},
		},
	}
	// ProjectMembersColumns holds the columns for the "project_members" table.
	ProjectMembersColumns = []*schema.Column{
		{Name: "project_id", Type: field.TypeInt},
		{Name: "employee_id", Type: field.TypeInt},
	}
	// ProjectMembersTable holds the schema information for the "project_members" table.
	ProjectMembersTable = &schema.Table{
		Name:       "project_members",
		Columns:    ProjectMembersColumns,
		PrimaryKey: []*schema.Column{ProjectMembersColumns[0], ProjectMembersColumns[1]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "project_members_project_id",
				Columns:    []*schema.Column{ProjectMembersColumns[0]},
				RefColumns: []*schema.Column{ProjectsColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "project_members_employee_id",
				Columns:    []*schema.Column{ProjectMembersColumns[1]},
				RefColumns: []*schema.Column{EmployeesColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// TaskLabelsColumns holds the columns for the "task_labels" table.
	TaskLabelsColumns = []*schema.Column{
		{Name: "task_id", Type: field.TypeInt},
		{Name: "label_id", Type: field.TypeInt},
	}
	// TaskLabelsTable holds the schema information for the "task_labels" table.
	TaskLabelsTable = &schema.Table{
		Name:       "task_labels",
		Columns:    TaskLabelsColumns,
		PrimaryKey: []*schema.Column{TaskLabelsColumns[0], TaskLabelsColumns[1]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "task_labels_task_id",
				Columns:    []*schema.Column{TaskLabelsColumns[0]},
				RefColumns: []*schema.Column{TasksColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "task_labels_label_id",
				Columns:    []*schema.Column{TaskLabelsColumns[1]},
				RefColumns: []*schema.Column{LabelsColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// TaskAssigneesColumns holds the columns for the "task_assignees" table.
	TaskAssigneesColumns = []*schema.Column{
		{Name: "task_id", Type: field.TypeInt},
		{Name: "employee_id", Type: field.TypeInt},
	}
	// TaskAssigneesTable holds the schema information for the "task_assignees" table.
	TaskAssigneesTable = &schema.Table{
		Name:       "task_assignees",
		Columns:    TaskAssigneesColumns,
		PrimaryKey: []*schema.Column{TaskAssigneesColumns[0], TaskAssigneesColumns[1]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "task_assignees_task_id",
				Columns:    []*schema.Column{TaskAssigneesColumns[0]},
				RefColumns: []*schema.Column{TasksColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "task_assignees_employee_id",
				Columns:    []*schema.Column{TaskAssigneesColumns[1]},
				RefColumns: []*schema.Column{EmployeesColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		AppointmentHistoriesTable,
		DepartmentsTable,
		EmployeesTable,
		LabelsTable,
		LeaveApprovalsTable,
		LeaveRequestsTable,
		OrganizationsTable,
		PositionsTable,
		ProjectsTable,
		TasksTable,
		TaskReportsTable,
		ProjectMembersTable,
		TaskLabelsTable,
		TaskAssigneesTable,
	}
)

func init() {
	AppointmentHistoriesTable.ForeignKeys[0].RefTable = EmployeesTable
	DepartmentsTable.ForeignKeys[0].RefTable = OrganizationsTable
	EmployeesTable.ForeignKeys[0].RefTable = PositionsTable
	LabelsTable.ForeignKeys[0].RefTable = OrganizationsTable
	LeaveApprovalsTable.ForeignKeys[0].RefTable = EmployeesTable
	LeaveApprovalsTable.ForeignKeys[1].RefTable = LeaveRequestsTable
	LeaveRequestsTable.ForeignKeys[0].RefTable = EmployeesTable
	LeaveRequestsTable.ForeignKeys[1].RefTable = OrganizationsTable
	OrganizationsTable.ForeignKeys[0].RefTable = OrganizationsTable
	PositionsTable.ForeignKeys[0].RefTable = DepartmentsTable
	PositionsTable.ForeignKeys[1].RefTable = PositionsTable
	ProjectsTable.ForeignKeys[0].RefTable = EmployeesTable
	ProjectsTable.ForeignKeys[1].RefTable = EmployeesTable
	ProjectsTable.ForeignKeys[2].RefTable = OrganizationsTable
	TasksTable.ForeignKeys[0].RefTable = ProjectsTable
	TaskReportsTable.ForeignKeys[0].RefTable = EmployeesTable
	TaskReportsTable.ForeignKeys[1].RefTable = TasksTable
	ProjectMembersTable.ForeignKeys[0].RefTable = ProjectsTable
	ProjectMembersTable.ForeignKeys[1].RefTable = EmployeesTable
	TaskLabelsTable.ForeignKeys[0].RefTable = TasksTable
	TaskLabelsTable.ForeignKeys[1].RefTable = LabelsTable
	TaskAssigneesTable.ForeignKeys[0].RefTable = TasksTable
	TaskAssigneesTable.ForeignKeys[1].RefTable = EmployeesTable
}
