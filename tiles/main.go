package tiles

import (
	"os"
	"log/slog"

	"github.com/google/uuid"
	"github.com/protomaps/go-pmtiles/pmtiles"
	"github.com/tilezen/go-tilepacks/tilepack"
	"gitlab.com/spwoodcock/mb2osm/converter"
)

func GenerateMbTiles(logger *slog.Logger, basemapId uuid.UUID, urlTemplate string) {
	// Generate mbtiles
	// https://github.com/tilezen/go-tilepacks/blob/f75e8450bbe08f99964df15c86bdfb69852981e2/cmd/build/main.go#L161
	var jobCreator tilepack.JobGenerator
	var err error
	jobCreator, err = tilepack.NewFileTransportXYZJobGenerator(
		*fileTransportRoot,
		*urlTemplateStr,
		bounds,
		zooms,
		time.Duration(*requestTimeout)*time.Second,
		*invertedY,
		*ensureGzip,
	)
}

func GeneratePmTiles(logger *slog.Logger, basemapId uuid.UUID) {
	// Convert to pmtiles
	tempFile, err := os.CreateTemp("", "basemapId.tmp")
	if err != nil {
		logger.Error("Error generating pmtiles tempfile ", err)
	}
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
}

func GenerateOsmAnd(logger *slog.Logger, basemapId uuid.UUID) {
	// Convert to osmand
	err := converter.MbtilesToOsm(
		"basemapId.mbtiles",
		"basemapId.sqlitedb",
		80,   // jpegQuality
		true, // overwrite
	)
	if err != nil {
		logger.Error("Error generating osmand file ", err)
	}
}