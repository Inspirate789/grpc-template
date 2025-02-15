package repository

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

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

func (r *SqlxRepository) CreateUser(ctx context.Context, name string) (id uint64, err error) {
	dto := UserDTO{
		ID:   0,
		Name: name,
	}

	err = sqlxutils.NamedGet(ctx, r.db, &dto.ID, insertUserQuery, dto)
	if err != nil {
		return 0, err
	}

	return dto.ID, nil
}

func (r *SqlxRepository) UpdateUser(ctx context.Context, user models.User) (found bool, err error) {
	dto := UserDTO{ID: user.ID, Name: user.Name}
	res, err := sqlxutils.NamedExec(ctx, r.db, updateUserQuery, dto)
	if err != nil {
		return false, err
	}

	rowsCount, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsCount != 0, nil
}

func (r *SqlxRepository) DeleteUser(ctx context.Context, id uint64) error {
	_, err := sqlxutils.Exec(ctx, r.db, deleteUserQuery, id)

	return err
}

func (r *SqlxRepository) GetUser(ctx context.Context, id uint64) (user models.User, found bool, err error) {
	var dto UserDTO

	err = sqlxutils.Get(ctx, r.db, &dto, selectUserQuery, id)
	if errors.Is(err, sql.ErrNoRows) {
		return models.User{}, false, nil
	} else if err != nil {
		return models.User{}, false, err
	}

	return dto.ToModel(), true, nil
}

func (r *SqlxRepository) GetUsers(ctx context.Context, limit, offset uint64) ([]models.User, uint64, error) {
	res := make(UsersDTO, 0)

	err := sqlxutils.Select(ctx, r.db, &res, selectUsersQuery, limit, offset)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, 0, nil
	} else if err != nil {
		return nil, 0, err
	}

	model, totalCount := res.ToModel()

	return model, totalCount, nil
}

func (r *SqlxRepository) GetUsersByEvent(ctx context.Context, eventID, limit, offset uint64) ([]models.User, uint64, error) {
	res := make(UsersDTO, 0)

	err := sqlxutils.Select(ctx, r.db, &res, selectUsersByEventQuery, eventID, limit, offset)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, 0, nil
	} else if err != nil {
		return nil, 0, err
	}

	model, totalCount := res.ToModel()

	return model, totalCount, nil
}
