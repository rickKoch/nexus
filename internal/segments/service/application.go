package service

import (
	"context"

	"github.com/rickKoch/nexus/internal/segments/adapters"
	"github.com/rickKoch/nexus/internal/segments/app"
)

func NewApplication(ctx context.Context) (a app.Application, err error) {
	db, err := adapters.NewPostgreSQLConnection(adapters.DefaultPostgreSQLConfig())
	if err != nil {
		return a, err
	}

	segmentRepo := adapters.NewPostgreSQLSegmentRepository(db)

	seg, err := app.NewSegments(segmentRepo)
	if err != nil {
		return a, err
	}

	return app.Application{
		Segments: seg,
	}, nil
}
