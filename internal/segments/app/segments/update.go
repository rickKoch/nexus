package segments

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rickKoch/nexus/internal/segments/domain/segment"
)

// UpdateSegment holds the required parameters for updating a segment.
type UpdateSegment struct {
	ID         int
	Name       string
	TTLSeconds *int
}

// UpdateSegmentHandler defines the interface for updating a segment.
type UpdateSegmentHandler interface {
	Handle(ctx context.Context, props UpdateSegment) (*segment.Segment, error)
}

type updateSegmentHandler struct {
	segmentRepo segment.Repository
}

// NewUpdateSegmentHandler creates a new UpdateSegmentHandler.
func NewUpdateSegmentHandler(segmentRepo segment.Repository) (UpdateSegmentHandler, error) {
	if segmentRepo == nil {
		return updateSegmentHandler{}, errors.New("segment repository is not provided")
	}

	return updateSegmentHandler{segmentRepo}, nil
}

// Handle updates an existing segment.
func (h updateSegmentHandler) Handle(ctx context.Context, props UpdateSegment) (*segment.Segment, error) {
	// Validate input
	config := segment.SegmentConfig{
		Name:       props.Name,
		TTLSeconds: props.TTLSeconds,
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid segment: %w", err)
	}

	// Check if segment exists
	existing, err := h.segmentRepo.Get(ctx, props.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get segment '%d': %w", props.ID, err)
	}

	// Update the segment
	existing.Update(props.Name, props.TTLSeconds, time.Now())

	updated, err := h.segmentRepo.Update(ctx, existing)
	if err != nil {
		return nil, fmt.Errorf("failed to update segment '%d': %w", props.ID, err)
	}

	return updated, nil
}
