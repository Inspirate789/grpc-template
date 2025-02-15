package usecase

import (
	"context"
	"log/slog"

	"github.com/Inspirate789/grpc-template/internal/models"
)

type Repository interface {
	HealthCheck(ctx context.Context) error
	CreateUser(ctx context.Context, name string) (id uint64, err error)
	UpdateUser(ctx context.Context, user models.User) (found bool, err error)
	DeleteUser(ctx context.Context, id uint64) error
	GetUser(ctx context.Context, id uint64) (user models.User, found bool, err error)
	GetUsers(ctx context.Context, limit, offset uint64) ([]models.User, uint64, error)
	GetUsersByEvent(ctx context.Context, eventID, limit, offset uint64) ([]models.User, uint64, error)
}

type UseCase struct {
	repository Repository
	logger     *slog.Logger
}

func New(repository Repository, logger *slog.Logger) *UseCase {
	return &UseCase{
		repository: repository,
		logger:     logger,
	}
}

func (u *UseCase) HealthCheck(ctx context.Context) error {
	return u.repository.HealthCheck(ctx)
}

func (u *UseCase) CreateUser(ctx context.Context, name string) (id uint64, err error) {
	return u.repository.CreateUser(ctx, name)
}

func (u *UseCase) UpdateUser(ctx context.Context, user models.User) (found bool, err error) {
	return u.repository.UpdateUser(ctx, user)
}

func (u *UseCase) DeleteUser(ctx context.Context, id uint64) error {
	return u.repository.DeleteUser(ctx, id)
}

func (u *UseCase) GetUser(ctx context.Context, id uint64) (user models.User, found bool, err error) {
	return u.repository.GetUser(ctx, id)
}

func (u *UseCase) GetUsers(ctx context.Context, limit, offset uint64) ([]models.User, uint64, error) {
	return u.repository.GetUsers(ctx, limit, offset)
}

func (u *UseCase) GetUsersByEvent(ctx context.Context, eventID, limit, offset uint64) ([]models.User, uint64, error) {
	return u.repository.GetUsersByEvent(ctx, eventID, limit, offset)
}
