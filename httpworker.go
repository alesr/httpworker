package httpworker

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/alesr/httpworker/router"
	"github.com/alesr/svc"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

const (
	aliveCheckRoute = "/livez"
	defaultAddr     = ":8080"
)

var (
	_                    svc.Worker = (*Worker)(nil)
	_                    svc.Aliver = (*Worker)(nil)
	defaultShudownCtx, _            = context.WithTimeout(context.Background(), 5*time.Second)
)

// Worker is a svc.Worker implementation that uses the go-chi/chi router.
type Worker struct {
	logger             *zap.Logger
	router             chi.Router
	addr               string
	server             *http.Server
	shutdownTimeoutCtx context.Context
}

// New creates a new Worker with the provided logger and options.
func New(logger *zap.Logger, opts ...Option) *Worker {
	w := Worker{
		logger:             logger,
		addr:               defaultAddr,
		router:             router.New(),
		shutdownTimeoutCtx: defaultShudownCtx,
	}

	w.setupLivenessCheckRoute()

	for _, opt := range opts {
		opt(&w)
	}

	w.server = &http.Server{
		Addr:    w.addr,
		Handler: w.router,
	}
	return &w
}

func (w *Worker) Init(logger *zap.Logger) error {
	w.logger = w.logger.Named("http_worker")
	return nil
}

func (w *Worker) Run() error {
	w.logger.Info("Starting chi server", zap.String("address", w.server.Addr))

	if err := w.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		w.logger.Error("failed to serve chi server", zap.Error(err))
		return err
	}
	return nil
}

func (w *Worker) Terminate() error {
	w.logger.Info("Terminating chi server")

	if err := w.server.Shutdown(w.shutdownTimeoutCtx); err != nil {
		w.logger.Error("failed to shutdown chi server", zap.Error(err))
		return err
	}
	return nil
}

// Alive checks if the server is running by making a request to the health check route.
func (w *Worker) Alive() error {
	client := &http.Client{Timeout: 5 * time.Second}
	url := "http://" + w.server.Addr + aliveCheckRoute

	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("server is not responding with status OK")
	}
	return nil
}

func (w *Worker) setupLivenessCheckRoute() {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
	router.MountRouter(w.router, aliveCheckRoute, http.HandlerFunc(handlerFunc))
}
