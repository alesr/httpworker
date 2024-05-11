package httpworker

import (
	"context"

	"github.com/go-chi/chi/v5"
)

// Option is a functional option for configuring the Worker.
type Option func(w *Worker)

// WithRouter sets the chi router for the Worker.
func WithRouter(router chi.Router) Option {
	return func(w *Worker) {
		w.router = router
	}
}

// WithAddress sets the address for the Worker. Default is ":8080".
func WithAddress(addr string) Option {
	return func(w *Worker) {
		w.addr = addr
	}
}

// WithShutdownCtx sets the context for the Worker's shutdown timeout.
func WithShutdownCtx(ctx context.Context) Option {
	return func(w *Worker) {
		w.shutdownTimeoutCtx = ctx
	}
}
