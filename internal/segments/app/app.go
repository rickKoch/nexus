package app

import (
	"github.com/rickKoch/nexus/internal/segments/app/segments"
	"github.com/rickKoch/nexus/internal/segments/domain/segment"
)

type Application struct {
	Segments Segments
}

type Segments struct {
	GetSegment    segments.GetSegmentHandler
	ListSegments  segments.ListSegmentsHandler
	CreateSegment segments.CreateSegmentHandler
	UpdateSegment segments.UpdateSegmentHandler
	DeleteSegment segments.DeleteSegmentHandler
}

func NewSegments(repo segment.Repository) (Segments, error) {
	seg := Segments{}
	getHandler, err := segments.NewGetSegmentHandler(repo)
	if err != nil {
		return seg, err
	}

	listHandler, err := segments.NewListSegmentsHandler(repo)
	if err != nil {
		return seg, err
	}

	createHandler, err := segments.NewCreateSegmentHandler(repo)
	if err != nil {
		return seg, err
	}

	updateHandler, err := segments.NewUpdateSegmentHandler(repo)
	if err != nil {
		return seg, err
	}

	deleteHandler, err := segments.NewDeleteSegmentHandler(repo)
	if err != nil {
		return seg, err
	}

	return Segments{
		GetSegment:    getHandler,
		ListSegments:  listHandler,
		CreateSegment: createHandler,
		UpdateSegment: updateHandler,
		DeleteSegment: deleteHandler,
	}, nil
}
