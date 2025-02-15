package delivery

import (
	"context"
	"log/slog"
	"math"

	"github.com/Inspirate789/grpc-template/internal/models"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UseCase interface {
	HealthCheck(ctx context.Context) error
	CreateUser(ctx context.Context, name string) (id uint64, err error)
	UpdateUser(ctx context.Context, user models.User) (found bool, err error)
	DeleteUser(ctx context.Context, id uint64) error
	GetUser(ctx context.Context, id uint64) (user models.User, found bool, err error)
	GetUsers(ctx context.Context, limit, offset uint64) ([]models.User, uint64, error)
	GetUsersByEvent(ctx context.Context, eventID, limit, offset uint64) ([]models.User, uint64, error)
}

type Delivery struct {
	useCase UseCase
	logger  *slog.Logger
	UnimplementedUserServiceServer
}

func New(useCase UseCase, logger *slog.Logger) *Delivery {
	return &Delivery{
		useCase: useCase,
		logger:  logger,
	}
}

func (d *Delivery) Register(server grpc.ServiceRegistrar) {
	RegisterUserServiceServer(server, d)
}

func (d *Delivery) HealthCheck(ctx context.Context) error {
	return d.useCase.HealthCheck(ctx)
}

func (d *Delivery) CreateUser(ctx context.Context, request *CreateUserRequest) (*CreateUserResponse, error) {
	id, err := d.useCase.CreateUser(ctx, request.GetName())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &CreateUserResponse{Id: id}, nil
}

func (d *Delivery) UpdateUser(ctx context.Context, request *UpdateUserRequest) (*UpdateUserResponse, error) {
	user := models.User{
		ID:   request.GetUser().GetId(),
		Name: request.GetUser().GetName(),
	}

	found, err := d.useCase.UpdateUser(ctx, user)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	} else if !found {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &UpdateUserResponse{}, nil
}

func (d *Delivery) DeleteUser(ctx context.Context, request *DeleteUserRequest) (*DeleteUserResponse, error) {
	return &DeleteUserResponse{}, d.useCase.DeleteUser(ctx, request.GetId())
}

func (d *Delivery) GetUser(ctx context.Context, request *GetUserRequest) (*GetUserResponse, error) {
	user, found, err := d.useCase.GetUser(ctx, request.GetId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	} else if !found {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &GetUserResponse{
		User: &User{
			Id:   user.ID,
			Name: user.Name,
		},
	}, nil
}

func (d *Delivery) GetUsers(ctx context.Context, request *ListUsersRequest) (*ListUsersResponse, error) {
	eventID := request.GetEventId()
	offset := request.GetOffset()
	limit := request.GetLimit()
	if limit == 0 {
		limit = math.MaxInt32
	}

	var (
		users      []models.User
		totalCount uint64
		err        error
	)

	if request.EventId != nil {
		users, totalCount, err = d.useCase.GetUsersByEvent(ctx, eventID, limit, offset)
	} else {
		users, totalCount, err = d.useCase.GetUsers(ctx, limit, offset)
	}

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	dto := make([]*User, 0, len(users))
	for _, user := range users {
		dto = append(dto, &User{
			Id:   user.ID,
			Name: user.Name,
		})
	}

	return &ListUsersResponse{Users: dto, TotalCount: totalCount}, nil
}
