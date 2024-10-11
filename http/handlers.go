package http

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/google/uuid"

	"github.com/hotosm/basemap-api/tiles"
)

func getEndpoints(logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoContext(r.Context(), "msg", slog.String("handler", "getEndpoints"))

		jsonData := map[string]string{
			"generate basemaps": "http://localhost:8000/generate",
			"download mbtiles":  "http://localhost:8000/get-url/{id}?format=mbtiles",
			"download pmtiles":  "http://localhost:8000/get-url/{id}?format=pmtiles",
			"download osmand":   "http://localhost:8000/get-url/{id}?format=osmand",
			"heartbeat":         "http://localhost:8000/__lbheartbeat__",
		}
		err := encode(w, r, http.StatusOK, jsonData)
		if err != nil {
			logger.ErrorContext(r.Context(), "error encoding JSON", slog.String("error", err.Error()))
		}
	})
}

func getHeartbeat(logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

func postGenerateBasemaps(logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoContext(r.Context(), "generating basemaps", slog.String("handler", "postGenerateBasemaps"))

		basemapId := uuid.New()

		urlTemplate := r.URL.Query().Get("url-template")
		if urlTemplate == "" {
			logger.Error("url-template param was not provided")
			w.WriteHeader(http.StatusUnprocessableEntity)
		}

		tiles.GenerateMbTiles(logger, basemapId, urlTemplate)
		tiles.GeneratePmTiles(logger, basemapId)
		tiles.GenerateOsmAnd(logger, basemapId)
	})
}

func getBasemapS3URL(logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		basemapId := r.PathValue("id")
		format := r.URL.Query().Get("format")
		if format == "" {
			// Default mbtiles download if no format specfied
			format = "mbtiles"
		}

		logger.InfoContext(
			r.Context(),
			"downloading basemap",
			slog.String("id", basemapId),
			slog.String("handler", "getBasemapS3URL"),
			slog.String("format", format),
		)

		jsonData := map[string]string{
			"url": "BASEMAP_S3_ENDPOINT/BASEMAP_S3_BUCKET/BASEMAP_S3_PATH_PREFIX/basemapId.format",
		}
		err := encode(w, r, http.StatusOK, jsonData)
		if err != nil {
			logger.ErrorContext(r.Context(), "error encoding JSON", slog.String("error", err.Error()))
		}
	})
}
