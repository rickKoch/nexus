package segments

import (
	"context"
	"errors"
	"fmt"

	"github.com/rickKoch/nexus/internal/segments/domain/segment"
)

// CreateSegment holds the required parameters for creating a segment.
type CreateSegment struct {
	Name       string
	TTLSeconds *int
}

// CreateSegmentHandler defines the interface for creating a segment.
type CreateSegmentHandler interface {
	Handle(ctx context.Context, props CreateSegment) (*segment.Segment, error)
}

type createSegmentHandler struct {
	segmentRepo segment.Repository
}

// NewCreateSegmentHandler creates a new CreateSegmentHandler.
func NewCreateSegmentHandler(segmentRepo segment.Repository) (CreateSegmentHandler, error) {
	if segmentRepo == nil {
		return createSegmentHandler{}, errors.New("segment repository is not provided")
	}

	return createSegmentHandler{segmentRepo}, nil
}

// Handle creates a new segment.
func (h createSegmentHandler) Handle(ctx context.Context, props CreateSegment) (*segment.Segment, error) {
	config := segment.SegmentConfig{
		Name:       props.Name,
		TTLSeconds: props.TTLSeconds,
	}

	factory, err := segment.NewFactory(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create segment factory: %w", err)
	}

	newSegment := factory.NewSegment()

	created, err := h.segmentRepo.Create(ctx, newSegment)
	if err != nil {
		return nil, fmt.Errorf("failed to create segment: %w", err)
	}

	return created, nil
}
