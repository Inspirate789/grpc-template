package repository

import (
	"time"

	"github.com/Inspirate789/grpc-template/internal/models"
)

const TimestampLayout = time.RFC3339

type EventDTO struct {
	ID        uint64 `db:"id"`
	Name      string `db:"name"`
	Timestamp string `db:"timestamp"`
}

type EventUserDTO struct {
	UserID  uint64 `db:"user_id"`
	EventID uint64 `db:"event_id"`
}

type EventWithUsersDTO struct {
	EventDTO
	UserIDs []uint64 `db:"user_ids"`
}

func (dto EventWithUsersDTO) ToModel() (models.Event, error) {
	timestamp, err := time.Parse(TimestampLayout, dto.Timestamp)
	if err != nil {
		return models.Event{}, err
	}

	return models.Event{
		ID:        dto.ID,
		Name:      dto.Name,
		Timestamp: timestamp,
		UserIDs:   dto.UserIDs,
	}, nil
}

type CountedEventDTO struct {
	EventWithUsersDTO
	TotalCount uint64 `db:"total_count"`
}

func (dto CountedEventDTO) ToModel() (models.Event, error) {
	timestamp, err := time.Parse(TimestampLayout, dto.Timestamp)
	if err != nil {
		return models.Event{}, err
	}

	return models.Event{
		ID:        dto.ID,
		Name:      dto.Name,
		Timestamp: timestamp,
		UserIDs:   dto.UserIDs,
	}, nil
}

type EventsDTO []CountedEventDTO

func (dto EventsDTO) ToModel() ([]models.Event, uint64, error) {
	res := make([]models.Event, 0, len(dto))

	for _, event := range dto {
		model, err := event.ToModel()
		if err != nil {
			return nil, 0, err
		}

		res = append(res, model)
	}

	if len(dto) != 0 {
		return res, dto[0].TotalCount, nil
	}

	return res, 0, nil
}
