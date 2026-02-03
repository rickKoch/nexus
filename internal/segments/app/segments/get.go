package segments

import (
	"context"
	"errors"
	"fmt"

	"github.com/rickKoch/nexus/internal/segments/domain/segment"
)

type GetSegment struct {
	ID int
}

type GetSegmentHandler interface {
	Handle(ctx context.Context, props GetSegment) (*segment.Segment, error)
}

type getSegmentHandler struct {
	segmentRepo segment.Repository
}

func NewGetSegmentHandler(segmentRepo segment.Repository) (GetSegmentHandler, error) {
	if segmentRepo == nil {
		return getSegmentHandler{}, errors.New("segment repository is not provided")
	}

	return getSegmentHandler{segmentRepo}, nil
}

func (gs getSegmentHandler) Handle(ctx context.Context, props GetSegment) (*segment.Segment, error) {
	segment, err := gs.segmentRepo.Get(ctx, props.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get segment '%d': %w", props.ID, err)
	}

	return segment, nil
}
