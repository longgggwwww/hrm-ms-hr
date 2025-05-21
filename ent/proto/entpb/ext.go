package entpb

import (
	context "context"

	"github.com/longgggwwww/hrm-ms-hr/ent"
	employee "github.com/longgggwwww/hrm-ms-hr/ent/employee"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// ExtService implements ExtServiceServer.
type ExtService struct {
	client *ent.Client
	UnimplementedHRExtServiceServer
}

func (s *ExtService) GetEmployeeByUserId(ctx context.Context, req *GetEmployeeByUserIdRequest) (*Employee, error) {
	e, err := s.client.Employee.Query().Where(employee.UserID(req.UserId)).WithPosition().Only(ctx)
	if err != nil {
		return nil, err
	}
	return toProtoEmployee(e)
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
