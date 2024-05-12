package examples

import (
	"context"
	"net/http"
	"time"

	"github.com/alesr/httpworker"
	"github.com/alesr/httpworker/router"
	"github.com/alesr/svc"
	"go.uber.org/zap"
)

func ExampleWorker() {
	r := router.New()

	// Add routes to the router.
	router.MountRouter(r, "/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, HTTP Worker!"))
	}))

	// Pass the router to the WithRouter option.
	routerOpt := httpworker.WithRouter(r)

	// Pass the address to the WithAddress option.
	addrOpt := httpworker.WithAddress(":8081")

	// Pass the shutdown context to the WithShutdownCtx option.

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	shutdownCtx := httpworker.WithShutdownCtx(ctx)

	// Instantiate the http worker.
	httpWorker := httpworker.New(zap.NewNop(), routerOpt, addrOpt, shutdownCtx)

	// Create a new service.
	svc, _ := svc.New("svcdemo", "0.0.1")

	// Add the http worker to the service.
	svc.AddWorker("http-worker", httpWorker)

	// Run the service.
	svc.Run()
}
