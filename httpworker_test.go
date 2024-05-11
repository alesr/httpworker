package httpworker

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNew(t *testing.T) {
	w := New(zap.NewNop())

	require.NotNil(t, w)

	assert.NotNil(t, w.logger)
	assert.NotNil(t, w.router)
	assert.NotNil(t, w.server)
	assert.Equal(t, defaultAddr, w.addr)
	assert.Equal(t, defaultAddr, w.server.Addr)
	assert.Equal(t, defaultShudownCtx, w.shutdownTimeoutCtx)

	// Check if the liveness check route was added.

	req := httptest.NewRequest(http.MethodGet, aliveCheckRoute, nil)
	rr := httptest.NewRecorder()
	w.router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestInit(t *testing.T) {
	logger := zap.NewNop()

	w := New(logger)

	err := w.Init(logger)
	require.NoError(t, err)

	assert.Equal(t, logger.Named("http_worker"), w.logger)
}

func TestRun(t *testing.T) {
	w := New(zap.NewNop())

	err := w.Init(zap.NewNop())
	require.NoError(t, err)

	go func() {
		err := w.Run()
		if err != nil && err != http.ErrServerClosed {
			t.Errorf("HTTP server error: %v", err)
		}
	}()
	defer w.Terminate()

	httpCli := &http.Client{}

	u, err := url.Parse(fmt.Sprintf("http://localhost%s%s", defaultAddr, aliveCheckRoute))
	require.NoError(t, err)

	resp, err := httpCli.Get(u.String())
	require.NoError(t, err)

	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestTerminate(t *testing.T) {
	w := New(zap.NewNop())

	err := w.Init(zap.NewNop())
	require.NoError(t, err)

	go func() {
		err := w.Run()
		if err != nil && err != http.ErrServerClosed {
			t.Errorf("HTTP server error: %v", err)
		}
	}()

	// Worker is running:

	httpCli := &http.Client{}

	u, err := url.Parse(fmt.Sprintf("http://localhost%s%s", defaultAddr, aliveCheckRoute))
	require.NoError(t, err)

	resp, err := httpCli.Get(u.String())
	require.NoError(t, err)

	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Now, terminate the worker:

	err = w.Terminate()
	assert.NoError(t, err)

	// Worker should be terminated:

	resp, err = httpCli.Get(u.String())
	assert.Error(t, err)

	assert.Nil(t, resp)
	assert.Error(t, err, err.Error())
}

func TestAlive(t *testing.T) {
	w := New(zap.NewNop())

	err := w.Init(zap.NewNop())

	require.NoError(t, err)

	go func() {
		err := w.Run()
		if err != nil && err != http.ErrServerClosed {
			t.Errorf("HTTP server error: %v", err)
		}
	}()
	defer w.Terminate()

	err = w.Alive()
	assert.NoError(t, err)
}
