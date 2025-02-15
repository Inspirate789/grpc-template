package migrations

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
)

func up(migrator *migrate.Migrate, version uint, logger *slog.Logger) error {
	err := migrator.Up()
	switch {
	case err != nil && !errors.Is(err, migrate.ErrNoChange):
		rollbackErr := migrator.Force(int(version))
		if rollbackErr != nil {
			logger.Warn(err.Error() + "; rollback")
			logger.Warn("rollback migration")
			return rollbackErr
		}

		logger.Warn("rollback migration")

		return err
	case errors.Is(err, migrate.ErrNoChange):
		logger.Info("current database migration version is up to date")
	default:
		var dirty bool

		version, dirty, err = migrator.Version()
		if err != nil {
			return err
		} else if dirty {
			return fmt.Errorf("database migration version %d is dirty", version)
		}

		logger.Info(fmt.Sprintf("migrated to version %d", version))
	}

	return nil
}

func Do(driverName, migrationsPath string, dbInstance database.Driver, logger *slog.Logger) error {
	migrator, err := migrate.NewWithDatabaseInstance("file://"+migrationsPath, driverName, dbInstance)
	if err != nil {
		return err
	}

	version, dirty, err := migrator.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		return err
	} else if dirty {
		return fmt.Errorf("current database migration version %d is dirty", version)
	}

	logger.Info(fmt.Sprintf("current database migration version is %d; migrate up", version))

	return up(migrator, version, logger)
}
