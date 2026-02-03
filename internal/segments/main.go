package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rickKoch/nexus/internal/segments/port"
	"github.com/rickKoch/nexus/internal/segments/service"
	"github.com/rickKoch/nexus/pkg/server"
	"github.com/rickKoch/nexus/pkg/signals"
	"github.com/sirupsen/logrus"
)

func main() {
	ctx := signals.Context()
	application, err := service.NewApplication(ctx)
	if err != nil {
		logrus.WithError(err).Panic("Failed to initialize application")
	}

	server.RunHTTPServer(func(router chi.Router) http.Handler {
		return port.HandlerFromMux(port.NewHttpServer(application), router)
	})
}
