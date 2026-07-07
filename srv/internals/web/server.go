package web

import (
	"context"
	"net/http"

	"github.com/yomiAdenaike01/support-queue/internals/config"
)

func Run(ctx context.Context, config *config.Config, handler http.Handler) {
	srv := &http.Server{
		Handler: handler,
		Addr:    config.GetEnvOrDefault("SERVER_ADDR", "0.0.0.0:2342"),
	}

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			panic(err)
		}
	}()

	<-ctx.Done()
}
