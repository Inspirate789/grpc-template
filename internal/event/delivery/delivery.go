package delivery

import (
	"context"
	"log/slog"
	"math"
	"time"

	"github.com/Inspirate789/grpc-template/internal/models"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UseCase interface {
	HealthCheck(ctx context.Context) error
	CreateEvent(ctx context.Context, name string, timestamp time.Time, userIDs []uint64) (id uint64, err error)
	UpdateEvent(ctx context.Context, event models.Event) (found bool, err error)
	DeleteEvent(ctx context.Context, id uint64) error
	GetEvent(ctx context.Context, id uint64) (event models.Event, found bool, err error)
	GetEvents(ctx context.Context, limit, offset uint64) ([]models.Event, uint64, error)
	GetEventsByUser(ctx context.Context, userID, limit, offset uint64) ([]models.Event, uint64, error)
}

type Delivery struct {
	useCase UseCase
	logger  *slog.Logger
	UnimplementedEventServiceServer
}

func New(useCase UseCase, logger *slog.Logger) *Delivery {
	return &Delivery{
		useCase: useCase,
		logger:  logger,
	}
}

func (d *Delivery) Register(server grpc.ServiceRegistrar) {
	RegisterEventServiceServer(server, d)
}

func (d *Delivery) HealthCheck(ctx context.Context) error {
	return d.useCase.HealthCheck(ctx)
}

func (d *Delivery) CreateEvent(ctx context.Context, request *CreateEventRequest) (*CreateEventResponse, error) {
	id, err := d.useCase.CreateEvent(ctx, request.GetName(), request.GetTimestamp().AsTime(), request.GetUserIds())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &CreateEventResponse{Id: id}, nil
}

func (d *Delivery) UpdateEvent(ctx context.Context, request *UpdateEventRequest) (*UpdateEventResponse, error) {
	event := models.Event{
		ID:        request.GetEvent().GetId(),
		Name:      request.GetEvent().GetName(),
		Timestamp: request.GetEvent().GetTimestamp().AsTime(),
		UserIDs:   request.GetEvent().GetUserIds(),
	}

	found, err := d.useCase.UpdateEvent(ctx, event)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	} else if !found {
		return nil, status.Error(codes.NotFound, "event not found")
	}

	return &UpdateEventResponse{}, nil
}

func (d *Delivery) DeleteEvent(ctx context.Context, request *DeleteEventRequest) (*DeleteEventResponse, error) {
	return &DeleteEventResponse{}, d.useCase.DeleteEvent(ctx, request.GetId())
}

func (d *Delivery) GetEvent(ctx context.Context, request *GetEventRequest) (*GetEventResponse, error) {
	event, found, err := d.useCase.GetEvent(ctx, request.GetId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	} else if !found {
		return nil, status.Error(codes.NotFound, "event not found")
	}

	return &GetEventResponse{
		Event: &Event{
			Id:        event.ID,
			Name:      event.Name,
			Timestamp: timestamppb.New(event.Timestamp),
			UserIds:   event.UserIDs,
		},
	}, nil
}

func (d *Delivery) GetEvents(ctx context.Context, request *ListEventsRequest) (*ListEventsResponse, error) {
	userID := request.GetUserId()
	offset := request.GetOffset()
	limit := request.GetLimit()
	if limit == 0 {
		limit = math.MaxInt32
	}

	var (
		events     []models.Event
		totalCount uint64
		err        error
	)

	if request.UserId != nil {
		events, totalCount, err = d.useCase.GetEventsByUser(ctx, userID, limit, offset)
	} else {
		events, totalCount, err = d.useCase.GetEvents(ctx, limit, offset)
	}

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	dto := make([]*Event, 0, len(events))
	for _, event := range events {
		dto = append(dto, &Event{
			Id:        event.ID,
			Name:      event.Name,
			Timestamp: timestamppb.New(event.Timestamp),
			UserIds:   event.UserIDs,
		})
	}

	return &ListEventsResponse{Events: dto, TotalCount: totalCount}, nil
}
