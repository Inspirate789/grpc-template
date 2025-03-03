package app

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/pkg/errors"
	slogfiber "github.com/samber/slog-fiber"
	"go.uber.org/multierr"
)

type HealthChecker interface {
	HealthCheck(ctx context.Context) error
}

type WebDelivery interface {
	HealthChecker
	AddHandlers(router fiber.Router)
}

type WebConfig struct {
	Host       string
	Port       string
	PathPrefix string
}

type WebApp struct {
	config WebConfig
	app    *fiber.App
	logger *slog.Logger
}

func newFiberError(msg string) fiber.Map {
	return fiber.Map{"message": msg}
}

func checkReadiness(appComponents ...HealthChecker) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var err error
		for _, component := range appComponents {
			err = multierr.Append(err, component.HealthCheck(ctx.Context()))
		}
		if err != nil {
			return ctx.Status(fiber.StatusServiceUnavailable).JSON(newFiberError(err.Error()))
		}

		return ctx.Status(fiber.StatusOK).SendString("healthy")
	}
}

func NewWebApp(
	config WebConfig,
	delivery []WebDelivery,
	auth fiber.Handler,
	logger *slog.Logger,
	appComponents ...HealthChecker,
) *WebApp {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			logger.Error(err.Error())
			msg := strings.SplitN(err.Error(), ":", 2)[0]

			var DNSError *net.DNSError
			if errors.As(err, &DNSError) {
				return ctx.Status(fiber.StatusServiceUnavailable).JSON(newFiberError(msg))
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(newFiberError(msg))
		},
	})

	app.Use(recover.New(recover.Config{EnableStackTrace: true}))
	app.Use(slogfiber.New(logger))
	app.Use(pprof.New())

	api := app.Group(config.PathPrefix)

	if auth != nil {
		api.Use(auth)
	}

	for _, d := range delivery {
		appComponents = append(appComponents, d)
		d.AddHandlers(api)
	}

	app.Get("/manage/health", checkReadiness(appComponents...))

	return &WebApp{
		config: config,
		app:    app,
		logger: logger,
	}
}

func (app *WebApp) Start() error {
	return errors.Wrap(app.app.Listen(app.config.Host+":"+app.config.Port), "start web app")
}

func (app *WebApp) Shutdown(ctx context.Context) error {
	return errors.Wrap(app.app.ShutdownWithContext(ctx), "stop web app")
}

func (app *WebApp) Test(req *http.Request, msTimeout ...int) (*http.Response, error) {
	return app.app.Test(req, msTimeout...)
}
