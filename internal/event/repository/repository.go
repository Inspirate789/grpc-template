package repository

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/Inspirate789/grpc-template/internal/models"
	"github.com/Inspirate789/grpc-template/pkg/sqlxutils"
	"github.com/jmoiron/sqlx"
)

type SqlxRepository struct {
	db     *sqlx.DB
	logger *slog.Logger
}

func NewSqlx(db *sqlx.DB, logger *slog.Logger) *SqlxRepository {
	return &SqlxRepository{
		db:     db,
		logger: logger,
	}
}

func (r *SqlxRepository) HealthCheck(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

func (r *SqlxRepository) CreateEvent(ctx context.Context, name string, timestamp time.Time, userIDs []uint64) (uint64, error) {
	dto := EventDTO{
		ID:        0,
		Name:      name,
		Timestamp: timestamp.Format(TimestampLayout),
	}

	err := sqlxutils.RunTx(ctx, r.db, sql.LevelDefault, func(tx *sqlx.Tx) error {
		err := sqlxutils.NamedGet(ctx, tx, &dto.ID, insertEventQuery, dto)
		if err != nil {
			return err
		}

		for _, userID := range userIDs {
			_, err = sqlxutils.NamedExec(ctx, tx, insertEventUserQuery, EventUserDTO{
				UserID:  userID,
				EventID: dto.ID,
			})
			if err != nil {
				return err
			}
		}

		return nil
	})

	return dto.ID, err
}

func (*SqlxRepository) updateEventUsersTx(
	ctx context.Context,
	tx sqlx.ExtContext,
	eventID uint64,
	oldUsers, newUsers []uint64,
) error {
	old := make([]uint64, 0)

	for _, userID := range oldUsers {
		if !slices.Contains(newUsers, userID) {
			old = append(old, userID)
		}
	}

	for _, userID := range newUsers {
		if !slices.Contains(oldUsers, userID) {
			_, err := sqlxutils.NamedExec(ctx, tx, insertEventUserQuery, EventUserDTO{
				UserID:  userID,
				EventID: eventID,
			})
			if err != nil {
				return err
			}
		}
	}

	var err error
	if len(old) != 0 {
		arr := make([]string, 0, len(old))
		for _, userID := range old {
			arr = append(arr, strconv.FormatUint(userID, 10))
		}

		_, err = sqlxutils.Exec(ctx, tx, deleteEventUsersQuery, eventID, strings.Join(arr, ","))
	}

	return err
}

func (r *SqlxRepository) UpdateEvent(ctx context.Context, event models.Event) (found bool, err error) {
	var res sql.Result

	err = sqlxutils.RunTx(ctx, r.db, sql.LevelDefault, func(tx *sqlx.Tx) error {
		var existingEvent models.Event
		existingEvent, found, err = r.getEventTx(ctx, tx, event.ID)
		if err != nil {
			return err
		}

		if !slices.Equal(event.UserIDs, existingEvent.UserIDs) {
			err = r.updateEventUsersTx(ctx, tx, event.ID, existingEvent.UserIDs, event.UserIDs)
			if err != nil {
				return err
			}
		}

		dto := EventDTO{ID: event.ID, Name: event.Name, Timestamp: event.Timestamp.Format(TimestampLayout)}
		res, err = sqlxutils.NamedExec(ctx, tx, updateEventQuery, dto)

		return err
	})
	if err != nil {
		return false, err
	}

	rowsCount, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsCount != 0, nil
}

func (r *SqlxRepository) DeleteEvent(ctx context.Context, id uint64) error {
	_, err := sqlxutils.Exec(ctx, r.db, deleteEventQuery, id)

	return err
}

func (*SqlxRepository) getEventTx(ctx context.Context, tx sqlx.QueryerContext, id uint64) (models.Event, bool, error) {
	var dto EventWithUsersDTO
	err := sqlxutils.Get(ctx, tx, &dto, selectEventQuery, id)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Event{}, false, nil
	} else if err != nil {
		return models.Event{}, false, err
	}

	err = sqlxutils.Select(ctx, tx, &dto.UserIDs, selectUserIDsByEventQuery, id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return models.Event{}, false, err
	}

	model, err := dto.ToModel()

	return model, true, err
}

func (r *SqlxRepository) GetEvent(ctx context.Context, id uint64) (models.Event, bool, error) {
	var (
		event models.Event
		found bool
	)
	err := sqlxutils.RunTx(ctx, r.db, sql.LevelDefault, func(tx *sqlx.Tx) error {
		var txErr error
		event, found, txErr = r.getEventTx(ctx, tx, id)
		return txErr
	})

	return event, found, err
}

func (r *SqlxRepository) GetEvents(ctx context.Context, limit, offset uint64) ([]models.Event, uint64, error) {
	res := make(EventsDTO, 0)

	err := sqlxutils.RunTx(ctx, r.db, sql.LevelDefault, func(tx *sqlx.Tx) error {
		txErr := sqlxutils.Select(ctx, tx, &res, selectEventsQuery, limit, offset)
		if errors.Is(txErr, sql.ErrNoRows) {
			return nil
		} else if txErr != nil {
			return txErr
		}

		for i, dto := range res {
			txErr = sqlxutils.Select(ctx, tx, &res[i].UserIDs, selectUserIDsByEventQuery, dto.ID)
			if txErr != nil && !errors.Is(txErr, sql.ErrNoRows) {
				return txErr
			}
		}

		return nil
	})
	if err != nil {
		return nil, 0, err
	}

	return res.ToModel()
}

func (r *SqlxRepository) GetEventsByUser(ctx context.Context, userID, limit, offset uint64) ([]models.Event, uint64, error) {
	res := make(EventsDTO, 0)

	err := sqlxutils.RunTx(ctx, r.db, sql.LevelDefault, func(tx *sqlx.Tx) error {
		txErr := sqlxutils.Select(ctx, tx, &res, selectEventsByUserQuery, userID, limit, offset)
		if errors.Is(txErr, sql.ErrNoRows) {
			return nil
		} else if txErr != nil {
			return txErr
		}

		for i, dto := range res {
			txErr = sqlxutils.Select(ctx, tx, &res[i].UserIDs, selectUserIDsByEventQuery, dto.ID)
			if txErr != nil && !errors.Is(txErr, sql.ErrNoRows) {
				return txErr
			}
		}

		return nil
	})
	if err != nil {
		return nil, 0, err
	}

	return res.ToModel()
}
