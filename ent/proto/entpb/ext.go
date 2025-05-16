package entpb

import (
	context "context"

	"github.com/longgggwwww/hrm-ms-hr/ent"
	"github.com/longgggwwww/hrm-ms-hr/ent/employee"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// ExtService implements ExtServiceServer.
type ExtService struct {
	client *ent.Client
	UnimplementedExtServiceServer
}

// NewExtService returns a new ExtService.
func NewExtService(client *ent.Client) *ExtService {
	return &ExtService{
		client: client,
	}
}
func (s *ExtService) GetBranchByUserId(ctx context.Context, req *GetBranchByUserIdRequest) (*Branch, error) {
	userID := req.GetUserId()

	employee, err := s.client.Employee.
		Query().
		Where(employee.UserID(userID)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	branchEntity, err := s.client.Branch.
		Get(ctx, employee.BranchID)
	if err != nil {
		return nil, err
	}

	return &Branch{
		Id:          branchEntity.ID[:],
		Name:        branchEntity.Name,
		Code:        branchEntity.Code,
		CompanyId:   branchEntity.CompanyID[:],
		Address:     wrapperspb.String(branchEntity.Address),
		ContactInfo: wrapperspb.String(branchEntity.ContactInfo),
		CreatedAt:   timestamppb.New(branchEntity.CreatedAt),
		UpdatedAt:   timestamppb.New(branchEntity.UpdatedAt),
	}, nil
}

func (s *ExtService) DeleteEmployeeByUserId(ctx context.Context, req *DeleteEmployeeByUserIdRequest) (*DeleteEmployeeByUserIdResponse, error) {
	userID := req.GetUserId()

	_, err := s.client.Employee.
		Delete().
		Where(employee.UserID(userID)).
		Exec(ctx)
	if err != nil {
		return &DeleteEmployeeByUserIdResponse{
			Success: false,
		}, nil
	}

	return &DeleteEmployeeByUserIdResponse{
		Success: true,
	}, nil
}
