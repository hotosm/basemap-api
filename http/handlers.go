package http

import (
	"os"
	"log/slog"
	"net/http"

	"github.com/protomaps/go-pmtiles/pmtiles"
	"github.com/tilezen/go-tilepacks/tilepack"
	"gitlab.com/spwoodcock/mb2osm/converter"
)

func getEndpoints(logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoContext(r.Context(), "msg", slog.String("handler", "getEndpoints"))

		jsonData := map[string]string{
			"generate basemaps": "http://localhost:8000/generate",
			"download mbtiles":  "http://localhost:8000/download/{id}?format=mbtiles",
			"download pmtiles":  "http://localhost:8000/download/{id}?format=pmtiles",
			"download osmand":   "http://localhost:8000/download/{id}?format=osmand",
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

		w.WriteHeader(http.StatusOK)
	})
}

func getDownloadBasemap(logger *slog.Logger) http.Handler {
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
			slog.String("handler", "getDownloadBasemap"),
			slog.String("format", format),
		)

		// Generate mbtiles
		// https://github.com/tilezen/go-tilepacks/blob/f75e8450bbe08f99964df15c86bdfb69852981e2/cmd/build/main.go#L161
		var jobCreator tilepack.JobGenerator
		var err error
		jobCreator, err = tilepack.NewFileTransportXYZJobGenerator(*fileTransportRoot, *urlTemplateStr, bounds, zooms, time.Duration(*requestTimeout)*time.Second, *invertedY, *ensureGzip)

		// Convert to pmtiles
		tempFile, err := os.CreateTemp("", "basemapId.tmp")
		if err != nil {
			logger.Error("Error generating pmtiles tempfile ", err)
		}
		// Clean up after test
		defer os.Remove(tempFile.Name())
	
		err = pmtiles.Convert(
			logger,
			"basemapId.mbtiles",
			"basemapId.pmtiles",
			true, // deduplicate,
			tempFile,
		)
		if err != nil {
			logger.Error("Error generating pmtiles file ", err)
		}

		// Convert to osmand
		err = converter.MbtilesToOsm(
			"basemapId.mbtiles",
			"basemapId.sqlitedb",
			80,   // jpegQuality
			true, // overwrite
		)
		if err != nil {
			logger.Error("Error generating osmand file ", err)
		}

		w.WriteHeader(http.StatusOK)
	})
}
