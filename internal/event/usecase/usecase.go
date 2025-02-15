package usecase

import (
	"context"
	"log/slog"
	"time"

	"github.com/Inspirate789/grpc-template/internal/models"
)

type Repository interface {
	HealthCheck(ctx context.Context) error
	CreateEvent(ctx context.Context, name string, timestamp time.Time, userIDs []uint64) (id uint64, err error)
	UpdateEvent(ctx context.Context, event models.Event) (found bool, err error)
	DeleteEvent(ctx context.Context, id uint64) error
	GetEvent(ctx context.Context, id uint64) (event models.Event, found bool, err error)
	GetEvents(ctx context.Context, limit, offset uint64) ([]models.Event, uint64, error)
	GetEventsByUser(ctx context.Context, userID, limit, offset uint64) ([]models.Event, uint64, error)
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

func (u *UseCase) CreateEvent(ctx context.Context, name string, timestamp time.Time, userIDs []uint64) (id uint64, err error) {
	return u.repository.CreateEvent(ctx, name, timestamp, userIDs)
}

func (u *UseCase) UpdateEvent(ctx context.Context, event models.Event) (found bool, err error) {
	return u.repository.UpdateEvent(ctx, event)
}

func (u *UseCase) DeleteEvent(ctx context.Context, id uint64) error {
	return u.repository.DeleteEvent(ctx, id)
}

func (u *UseCase) GetEvent(ctx context.Context, id uint64) (event models.Event, found bool, err error) {
	return u.repository.GetEvent(ctx, id)
}

func (u *UseCase) GetEvents(ctx context.Context, limit, offset uint64) ([]models.Event, uint64, error) {
	return u.repository.GetEvents(ctx, limit, offset)
}

func (u *UseCase) GetEventsByUser(ctx context.Context, userID, limit, offset uint64) ([]models.Event, uint64, error) {
	return u.repository.GetEventsByUser(ctx, userID, limit, offset)
}
