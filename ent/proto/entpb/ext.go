package entpb

import (
	context "context"

	"github.com/longgggwwww/hrm-ms-hr/ent"
	"github.com/longgggwwww/hrm-ms-hr/ent/employee"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
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

func (s *ExtService) GetBranchByUserId(ctx context.Context, req *emptypb.Empty) (*Branch, error) {
	// Retrieve userID from the request
	userID := req.String()

	// Query the employee by user ID
	employee, err := s.client.Employee.
		Query().
		Where(employee.UserID(userID)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	// Query the branch by branch ID from the employee
	branchEntity, err := s.client.Branch.
		Get(ctx, employee.BranchID)
	if err != nil {
		return nil, err
	}

	// Convert the branch entity to the protobuf Branch message
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
