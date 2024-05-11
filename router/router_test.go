package router

import (
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestNewRouter(t *testing.T) {
	t.Parallel()

	t.Run("should create a new chi router", func(t *testing.T) {
		router := New()

		assert.NotEmpty(t, router)
		assert.Implements(t, (*chi.Router)(nil), router)
	})

	t.Run("should add logger and recoverer middleware", func(t *testing.T) {
		router := New()
		assert.Len(t, router.Middlewares(), 2)
	})
}

func TestGroupRouter(t *testing.T) {
	t.Parallel()

	t.Skip("to be implemented")
}

func TestUseMiddleware(t *testing.T) {
	t.Parallel()

	t.Skip("to be implemented")
}

func TestMountRouter(t *testing.T) {
	t.Parallel()

	t.Skip("to be implemented")
}
