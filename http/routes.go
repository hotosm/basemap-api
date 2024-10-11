package http

import (
	"log/slog"
	"net/http"
)

func addRoutes(
	mux *http.ServeMux,
	logger *slog.Logger,
	config Config,
) {
	mux.Handle("/", getEndpoints(logger))
	mux.Handle("/__lbheartbeat__", getHeartbeat(logger))
	mux.Handle("/generate/", postGenerateBasemaps(logger))
	mux.Handle("/get-url/{id}", getBasemapS3URL(logger))
}
