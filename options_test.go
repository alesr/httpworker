package httpworker

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWithRouter(t *testing.T) {
	t.Parallel()

	r := chi.NewRouter()

	endpoint := "/test"

	r.Get(endpoint, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})

	w := New(zap.NewNop(), WithRouter(r))

	req := httptest.NewRequest(http.MethodGet, endpoint, nil)
	rr := httptest.NewRecorder()

	w.router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusTeapot, rr.Code)
}

func TestWithAddress(t *testing.T) {
	t.Parallel()

	addr := ":4242"

	w := New(zap.NewNop(), WithAddress(addr))
	assert.Equal(t, addr, w.server.Addr)
}

func TestWithShutdownCtx(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	w := New(zap.NewNop(), WithShutdownCtx(ctx))

	assert.Equal(t, ctx, w.shutdownTimeoutCtx)
}
