package api

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/dupreehkuda/balance-microservice/internal/actions"
	"github.com/dupreehkuda/balance-microservice/internal/configuration"
	"github.com/dupreehkuda/balance-microservice/internal/handlers"
	i "github.com/dupreehkuda/balance-microservice/internal/interfaces"
	"github.com/dupreehkuda/balance-microservice/internal/logger"
	"github.com/dupreehkuda/balance-microservice/internal/middleware"
	"github.com/dupreehkuda/balance-microservice/internal/storage"
)

type api struct {
	handlers i.Handlers
	logger   *zap.Logger
	config   *configuration.Config
	mw       i.Middleware
}

func NewByConfig() *api {
	log := logger.InitializeLogger()

	cfg := configuration.New(log)

	store := storage.New(cfg.DatabasePath, log)
	store.CreateSchema()

	act := actions.New(store, log)

	handle := handlers.New(store, act, log)
	mware := middleware.New(log)

	return &api{
		handlers: handle,
		logger:   log,
		config:   cfg,
		mw:       mware,
	}
}

// Run runs the service
func (a api) Run() {
	serv := &http.Server{Addr: a.config.Address, Handler: a.service()}

	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		shutdownCtx, cancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer cancel()

		go func() {
			<-shutdownCtx.Done()
			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				a.logger.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		err := serv.Shutdown(shutdownCtx)
		if err != nil {
			a.logger.Fatal("Error shutting down", zap.Error(err))
		}
		a.logger.Info("Server shut down", zap.String("port", a.config.Address))
		serverStopCtx()
	}()

	a.logger.Info("Server started", zap.String("port", a.config.Address))
	err := serv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		a.logger.Fatal("Cant start server", zap.Error(err))
	}

	<-serverCtx.Done()
}

func (a api) service() http.Handler {
	r := chi.NewRouter()

	r.Use(a.mw.CheckCompression)
	r.Use(a.mw.WriteCompressed)

	r.Route("/api", func(r chi.Router) {
		r.Route("/balance", func(r chi.Router) {
			r.Route("/get/{account_id}", func(r chi.Router) {
				r.Use(a.mw.ParamCtx)
				r.Get("/", a.handlers.GetBalance)
			})

			r.Post("/add", a.handlers.AddFunds)
			r.Post("/transfer", a.handlers.TransferFunds)
		})

		r.Route("/order", func(r chi.Router) {
			r.Post("/reserve", a.handlers.ReserveFunds)
			r.Post("/withdraw", a.handlers.WithdrawBalance)
			r.Post("/cancel", a.handlers.CancelReserve)
		})

		r.Route("/accounting", func(r chi.Router) {
			r.Route("/get/{report_id}", func(r chi.Router) {
				r.Use(a.mw.ParamCtx)
				r.Get("/", a.handlers.GetReport)
			})

			r.Post("/add", a.handlers.GetReportLink)
		})
	})

	return r
}
