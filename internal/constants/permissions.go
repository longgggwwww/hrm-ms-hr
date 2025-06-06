package constants

// Employee permissions
const (
	EmployeeCreate = "employee:create"
	EmployeeRead   = "employee:read"
	EmployeeUpdate = "employee:update"
	EmployeeDelete = "employee:delete"
)

// Project permissions
const (
	ProjectCreate = "project:create"
	ProjectRead   = "project:read"
	ProjectUpdate = "project:update"
	ProjectDelete = "project:delete"
)

// Leave Request permissions
const (
	LeaveRequestReadAdmin      = "leave_request:read:admin"
	LeaveRequestReadEmployee   = "leave_request:read:employee"
	LeaveRequestApproveAdmin   = "leave_request:approve:admin"
	LeaveRequestRejectAdmin    = "leave_request:reject:admin"
	LeaveRequestCreateEmployee = "leave_request:create:employee"
)

// Task Report permissions
const (
	TaskReportCreate = "task_report:create"
	TaskReportUpdate = "task_report:update"
)

// Task permissions
const (
	TaskCreate = "task:create"
	TaskRead   = "task:read"
	TaskUpdate = "task:update"
	TaskDelete = "task:delete"
)

// Organization permissions
const (
	OrgCreate = "org:create"
	OrgRead   = "org:read"
	OrgUpdate = "org:update"
	OrgDelete = "org:delete"
)

// Department permissions
const (
	DepartmentCreate = "department:create"
	DepartmentRead   = "department:read"
	DepartmentUpdate = "department:update"
	DepartmentDelete = "department:delete"
)

// Position permissions
const (
	PositionCreate = "position:create"
	PositionRead   = "position:read"
	PositionUpdate = "position:update"
	PositionDelete = "position:delete"
)

// Label permissions
const (
	LabelCreate = "label:create"
	LabelRead   = "label:read"
	LabelUpdate = "label:update"
	LabelDelete = "label:delete"
)

// Permission groups for easier management
var (
	// Employee permission group
	EmployeePermissions = []string{
		EmployeeCreate,
		EmployeeRead,
		EmployeeUpdate,
		EmployeeDelete,
	}

	// Project permission group
	ProjectPermissions = []string{
		ProjectCreate,
		ProjectRead,
		ProjectUpdate,
		ProjectDelete,
	}

	// Leave Request permission group
	LeaveRequestPermissions = []string{
		LeaveRequestReadAdmin,
		LeaveRequestReadEmployee,
		LeaveRequestApproveAdmin,
		LeaveRequestRejectAdmin,
		LeaveRequestCreateEmployee,
	}

	// Task Report permission group
	TaskReportPermissions = []string{
		TaskReportCreate,
		TaskReportUpdate,
	}

	// Task permission group
	TaskPermissions = []string{
		TaskCreate,
		TaskRead,
		TaskUpdate,
		TaskDelete,
	}

	// Organization permission group
	OrgPermissions = []string{
		OrgCreate,
		OrgRead,
		OrgUpdate,
		OrgDelete,
	}

	// Department permission group
	DepartmentPermissions = []string{
		DepartmentCreate,
		DepartmentRead,
		DepartmentUpdate,
		DepartmentDelete,
	}

	// Position permission group
	PositionPermissions = []string{
		PositionCreate,
		PositionRead,
		PositionUpdate,
		PositionDelete,
	}

	// Label permission group
	LabelPermissions = []string{
		LabelCreate,
		LabelRead,
		LabelUpdate,
		LabelDelete,
	}

	// All permissions
	AllPermissions = []string{
		// Employee
		EmployeeCreate,
		EmployeeRead,
		EmployeeUpdate,
		EmployeeDelete,
		// Project
		ProjectCreate,
		ProjectRead,
		ProjectUpdate,
		ProjectDelete,
		// Leave Request
		LeaveRequestReadAdmin,
		LeaveRequestReadEmployee,
		LeaveRequestApproveAdmin,
		LeaveRequestRejectAdmin,
		LeaveRequestCreateEmployee,
		// Task Report
		TaskReportCreate,
		TaskReportUpdate,
		// Task
		TaskCreate,
		TaskRead,
		TaskUpdate,
		TaskDelete,
		// Organization
		OrgCreate,
		OrgRead,
		OrgUpdate,
		OrgDelete,
		// Department
		DepartmentCreate,
		DepartmentRead,
		DepartmentUpdate,
		DepartmentDelete,
		// Position
		PositionCreate,
		PositionRead,
		PositionUpdate,
		PositionDelete,
		// Label
		LabelCreate,
		LabelRead,
		LabelUpdate,
		LabelDelete,
	}
)

// Permission checker functions
func IsValidPermission(permission string) bool {
	for _, p := range AllPermissions {
		if p == permission {
			return true
		}
	}
	return false
}

func HasPermission(userPermissions []string, requiredPermission string) bool {
	for _, p := range userPermissions {
		if p == requiredPermission {
			return true
		}
	}
	return false
}

func HasAnyPermission(userPermissions []string, requiredPermissions []string) bool {
	for _, required := range requiredPermissions {
		if HasPermission(userPermissions, required) {
			return true
		}
	}
	return false
}

func HasAllPermissions(userPermissions []string, requiredPermissions []string) bool {
	for _, required := range requiredPermissions {
		if !HasPermission(userPermissions, required) {
			return false
		}
	}
	return true
}
