package entpb

import (
	context "context"
	"log"

	"github.com/longgggwwww/hrm-ms-hr/ent"
	employee "github.com/longgggwwww/hrm-ms-hr/ent/employee"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

// ExtService implements ExtServiceServer.
type ExtService struct {
	client *ent.Client
	UnimplementedExtServiceServer
}

func (s *ExtService) GetEmployeeByUserId(ctx context.Context, req *GetEmployeeByUserIdRequest) (*Employee, error) {
	log.Println("GetEmployeeByUserId called with UserId:", req.UserId)
	e, err := s.client.Employee.Query().Where(employee.UserID(req.UserId)).WithPosition(func(pq *ent.PositionQuery) {
		pq.WithDepartment(func(dq *ent.DepartmentQuery) {
			dq.WithOrganization()
		})
	}).Only(ctx)
	if err != nil {
		return nil, err
	}
	return toProtoEmployeeWithEdges(e)
}

func (s *ExtService) DeleteEmployeeByUserId(ctx context.Context, req *DeleteEmployeeByUserIdRequest) (*emptypb.Empty, error) {
	_, err := s.client.Employee.Delete().Where(employee.UserID(req.UserId)).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// NewExtService returns a new ExtService.
func NewExtService(client *ent.Client) *ExtService {
	return &ExtService{
		client: client,
	}
}

func toProtoEmployeeWithEdges(e *ent.Employee) (*Employee, error) {
	v := &Employee{}
	code := e.Code
	v.Code = code
	created_at := timestamppb.New(e.CreatedAt)
	v.CreatedAt = created_at
	id := int64(e.ID)
	v.Id = id
	joining_at := timestamppb.New(e.JoiningAt)
	v.JoiningAt = joining_at
	org_id := int64(e.OrgID)
	v.OrgId = org_id
	position := int64(e.PositionID)
	v.PositionId = position
	status := toProtoEmployee_Status(e.Status)
	v.Status = status
	updated_at := timestamppb.New(e.UpdatedAt)
	v.UpdatedAt = updated_at
	user_id := wrapperspb.String(e.UserID)
	v.UserId = user_id
	for _, edg := range e.Edges.AppointmentHistories {
		id := int64(edg.ID)
		v.AppointmentHistories = append(v.AppointmentHistories, &AppointmentHistory{
			Id: id,
		})
	}
	for _, edg := range e.Edges.AssignedTasks {
		id := int64(edg.ID)
		v.AssignedTasks = append(v.AssignedTasks, &Task{
			Id: id,
		})
	}
	for _, edg := range e.Edges.CreatedProjects {
		id := int64(edg.ID)
		v.CreatedProjects = append(v.CreatedProjects, &Project{
			Id: id,
		})
	}
	for _, edg := range e.Edges.LeaveApproves {
		id := int64(edg.ID)
		v.LeaveApproves = append(v.LeaveApproves, &LeaveApproval{
			Id: id,
		})
	}
	for _, edg := range e.Edges.LeaveRequests {
		id := int64(edg.ID)
		v.LeaveRequests = append(v.LeaveRequests, &LeaveRequest{
			Id: id,
		})
	}
	if edg := e.Edges.Position; edg != nil {
		id := int64(edg.ID)
		position := &Position{
			Id:   id,
			Name: edg.Name,
			Code: edg.Code,
		}

		// Include Department information
		if dept := edg.Edges.Department; dept != nil {
			deptId := int64(dept.ID)
			position.Department = &Department{
				Id:   deptId,
				Name: dept.Name,
				Code: dept.Code,
			}

			// Include Organization information
			if org := dept.Edges.Organization; org != nil {
				orgId := int64(org.ID)
				position.Department.Organization = &Organization{
					Id:   orgId,
					Name: org.Name,
					Code: org.Code,
				}
			}
		}

		v.Position = position
	}
	for _, edg := range e.Edges.Projects {
		id := int64(edg.ID)
		v.Projects = append(v.Projects, &Project{
			Id: id,
		})
	}
	for _, edg := range e.Edges.TaskReports {
		id := int64(edg.ID)
		v.TaskReports = append(v.TaskReports, &TaskReport{
			Id: id,
		})
	}
	for _, edg := range e.Edges.UpdatedProjects {
		id := int64(edg.ID)
		v.UpdatedProjects = append(v.UpdatedProjects, &Project{
			Id: id,
		})
	}
	return v, nil
}
