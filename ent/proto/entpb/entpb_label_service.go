// Code generated by protoc-gen-entgrpc. DO NOT EDIT.
package entpb

import (
	context "context"
	base64 "encoding/base64"
	entproto "entgo.io/contrib/entproto"
	runtime "entgo.io/contrib/entproto/runtime"
	sqlgraph "entgo.io/ent/dialect/sql/sqlgraph"
	fmt "fmt"
	ent "github.com/longgggwwww/hrm-ms-hr/ent"
	label "github.com/longgggwwww/hrm-ms-hr/ent/label"
	organization "github.com/longgggwwww/hrm-ms-hr/ent/organization"
	task "github.com/longgggwwww/hrm-ms-hr/ent/task"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
	strconv "strconv"
)

// LabelService implements LabelServiceServer
type LabelService struct {
	client *ent.Client
	UnimplementedLabelServiceServer
}

// NewLabelService returns a new LabelService
func NewLabelService(client *ent.Client) *LabelService {
	return &LabelService{
		client: client,
	}
}

// toProtoLabel transforms the ent type to the pb type
func toProtoLabel(e *ent.Label) (*Label, error) {
	v := &Label{}
	color := e.Color
	v.Color = color
	created_at := timestamppb.New(e.CreatedAt)
	v.CreatedAt = created_at
	description := wrapperspb.String(e.Description)
	v.Description = description
	id := int64(e.ID)
	v.Id = id
	name := e.Name
	v.Name = name
	organization := wrapperspb.Int64(int64(e.OrgID))
	v.OrgId = organization
	updated_at := timestamppb.New(e.UpdatedAt)
	v.UpdatedAt = updated_at
	if edg := e.Edges.Organization; edg != nil {
		id := int64(edg.ID)
		v.Organization = &Organization{
			Id: id,
		}
	}
	for _, edg := range e.Edges.Tasks {
		id := int64(edg.ID)
		v.Tasks = append(v.Tasks, &Task{
			Id: id,
		})
	}
	return v, nil
}

// toProtoLabelList transforms a list of ent type to a list of pb type
func toProtoLabelList(e []*ent.Label) ([]*Label, error) {
	var pbList []*Label
	for _, entEntity := range e {
		pbEntity, err := toProtoLabel(entEntity)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "internal error: %s", err)
		}
		pbList = append(pbList, pbEntity)
	}
	return pbList, nil
}

// Create implements LabelServiceServer.Create
func (svc *LabelService) Create(ctx context.Context, req *CreateLabelRequest) (*Label, error) {
	label := req.GetLabel()
	m, err := svc.createBuilder(label)
	if err != nil {
		return nil, err
	}
	res, err := m.Save(ctx)
	switch {
	case err == nil:
		proto, err := toProtoLabel(res)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "internal error: %s", err)
		}
		return proto, nil
	case sqlgraph.IsUniqueConstraintError(err):
		return nil, status.Errorf(codes.AlreadyExists, "already exists: %s", err)
	case ent.IsConstraintError(err):
		return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %s", err)
	default:
		return nil, status.Errorf(codes.Internal, "internal error: %s", err)
	}

}

// Get implements LabelServiceServer.Get
func (svc *LabelService) Get(ctx context.Context, req *GetLabelRequest) (*Label, error) {
	var (
		err error
		get *ent.Label
	)
	id := int(req.GetId())
	switch req.GetView() {
	case GetLabelRequest_VIEW_UNSPECIFIED, GetLabelRequest_BASIC:
		get, err = svc.client.Label.Get(ctx, id)
	case GetLabelRequest_WITH_EDGE_IDS:
		get, err = svc.client.Label.Query().
			Where(label.ID(id)).
			WithOrganization(func(query *ent.OrganizationQuery) {
				query.Select(organization.FieldID)
			}).
			WithTasks(func(query *ent.TaskQuery) {
				query.Select(task.FieldID)
			}).
			Only(ctx)
	default:
		return nil, status.Error(codes.InvalidArgument, "invalid argument: unknown view")
	}
	switch {
	case err == nil:
		return toProtoLabel(get)
	case ent.IsNotFound(err):
		return nil, status.Errorf(codes.NotFound, "not found: %s", err)
	default:
		return nil, status.Errorf(codes.Internal, "internal error: %s", err)
	}

}

// Update implements LabelServiceServer.Update
func (svc *LabelService) Update(ctx context.Context, req *UpdateLabelRequest) (*Label, error) {
	label := req.GetLabel()
	labelID := int(label.GetId())
	m := svc.client.Label.UpdateOneID(labelID)
	labelColor := label.GetColor()
	m.SetColor(labelColor)
	if label.GetDescription() != nil {
		labelDescription := label.GetDescription().GetValue()
		m.SetDescription(labelDescription)
	}
	labelName := label.GetName()
	m.SetName(labelName)
	if label.GetOrgId() != nil {
		labelOrgID := int(label.GetOrgId().GetValue())
		m.SetOrgID(labelOrgID)
	}
	labelUpdatedAt := runtime.ExtractTime(label.GetUpdatedAt())
	m.SetUpdatedAt(labelUpdatedAt)
	if label.GetOrganization() != nil {
		labelOrganization := int(label.GetOrganization().GetId())
		m.SetOrganizationID(labelOrganization)
	}
	for _, item := range label.GetTasks() {
		tasks := int(item.GetId())
		m.AddTaskIDs(tasks)
	}

	res, err := m.Save(ctx)
	switch {
	case err == nil:
		proto, err := toProtoLabel(res)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "internal error: %s", err)
		}
		return proto, nil
	case sqlgraph.IsUniqueConstraintError(err):
		return nil, status.Errorf(codes.AlreadyExists, "already exists: %s", err)
	case ent.IsConstraintError(err):
		return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %s", err)
	default:
		return nil, status.Errorf(codes.Internal, "internal error: %s", err)
	}

}

// Delete implements LabelServiceServer.Delete
func (svc *LabelService) Delete(ctx context.Context, req *DeleteLabelRequest) (*emptypb.Empty, error) {
	var err error
	id := int(req.GetId())
	err = svc.client.Label.DeleteOneID(id).Exec(ctx)
	switch {
	case err == nil:
		return &emptypb.Empty{}, nil
	case ent.IsNotFound(err):
		return nil, status.Errorf(codes.NotFound, "not found: %s", err)
	default:
		return nil, status.Errorf(codes.Internal, "internal error: %s", err)
	}

}

// List implements LabelServiceServer.List
func (svc *LabelService) List(ctx context.Context, req *ListLabelRequest) (*ListLabelResponse, error) {
	var (
		err      error
		entList  []*ent.Label
		pageSize int
	)
	pageSize = int(req.GetPageSize())
	switch {
	case pageSize < 0:
		return nil, status.Errorf(codes.InvalidArgument, "page size cannot be less than zero")
	case pageSize == 0 || pageSize > entproto.MaxPageSize:
		pageSize = entproto.MaxPageSize
	}
	listQuery := svc.client.Label.Query().
		Order(ent.Desc(label.FieldID)).
		Limit(pageSize + 1)
	if req.GetPageToken() != "" {
		bytes, err := base64.StdEncoding.DecodeString(req.PageToken)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "page token is invalid")
		}
		token, err := strconv.ParseInt(string(bytes), 10, 32)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "page token is invalid")
		}
		pageToken := int(token)
		listQuery = listQuery.
			Where(label.IDLTE(pageToken))
	}
	switch req.GetView() {
	case ListLabelRequest_VIEW_UNSPECIFIED, ListLabelRequest_BASIC:
		entList, err = listQuery.All(ctx)
	case ListLabelRequest_WITH_EDGE_IDS:
		entList, err = listQuery.
			WithOrganization(func(query *ent.OrganizationQuery) {
				query.Select(organization.FieldID)
			}).
			WithTasks(func(query *ent.TaskQuery) {
				query.Select(task.FieldID)
			}).
			All(ctx)
	}
	switch {
	case err == nil:
		var nextPageToken string
		if len(entList) == pageSize+1 {
			nextPageToken = base64.StdEncoding.EncodeToString(
				[]byte(fmt.Sprintf("%v", entList[len(entList)-1].ID)))
			entList = entList[:len(entList)-1]
		}
		protoList, err := toProtoLabelList(entList)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "internal error: %s", err)
		}
		return &ListLabelResponse{
			LabelList:     protoList,
			NextPageToken: nextPageToken,
		}, nil
	default:
		return nil, status.Errorf(codes.Internal, "internal error: %s", err)
	}

}

// BatchCreate implements LabelServiceServer.BatchCreate
func (svc *LabelService) BatchCreate(ctx context.Context, req *BatchCreateLabelsRequest) (*BatchCreateLabelsResponse, error) {
	requests := req.GetRequests()
	if len(requests) > entproto.MaxBatchCreateSize {
		return nil, status.Errorf(codes.InvalidArgument, "batch size cannot be greater than %d", entproto.MaxBatchCreateSize)
	}
	bulk := make([]*ent.LabelCreate, len(requests))
	for i, req := range requests {
		label := req.GetLabel()
		var err error
		bulk[i], err = svc.createBuilder(label)
		if err != nil {
			return nil, err
		}
	}
	res, err := svc.client.Label.CreateBulk(bulk...).Save(ctx)
	switch {
	case err == nil:
		protoList, err := toProtoLabelList(res)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "internal error: %s", err)
		}
		return &BatchCreateLabelsResponse{
			Labels: protoList,
		}, nil
	case sqlgraph.IsUniqueConstraintError(err):
		return nil, status.Errorf(codes.AlreadyExists, "already exists: %s", err)
	case ent.IsConstraintError(err):
		return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %s", err)
	default:
		return nil, status.Errorf(codes.Internal, "internal error: %s", err)
	}

}

func (svc *LabelService) createBuilder(label *Label) (*ent.LabelCreate, error) {
	m := svc.client.Label.Create()
	labelColor := label.GetColor()
	m.SetColor(labelColor)
	labelCreatedAt := runtime.ExtractTime(label.GetCreatedAt())
	m.SetCreatedAt(labelCreatedAt)
	if label.GetDescription() != nil {
		labelDescription := label.GetDescription().GetValue()
		m.SetDescription(labelDescription)
	}
	labelName := label.GetName()
	m.SetName(labelName)
	if label.GetOrgId() != nil {
		labelOrgID := int(label.GetOrgId().GetValue())
		m.SetOrgID(labelOrgID)
	}
	labelUpdatedAt := runtime.ExtractTime(label.GetUpdatedAt())
	m.SetUpdatedAt(labelUpdatedAt)
	if label.GetOrganization() != nil {
		labelOrganization := int(label.GetOrganization().GetId())
		m.SetOrganizationID(labelOrganization)
	}
	for _, item := range label.GetTasks() {
		tasks := int(item.GetId())
		m.AddTaskIDs(tasks)
	}
	return m, nil
}
