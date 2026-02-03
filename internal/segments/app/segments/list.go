package segments

import (
	"context"
	"errors"
	"fmt"

	"github.com/rickKoch/nexus/internal/segments/domain/segment"
)

// ListSegmentsHandler defines the interface for listing segments.
type ListSegmentsHandler interface {
	Handle(ctx context.Context) ([]segment.Segment, error)
}

type listSegmentsHandler struct {
	segmentRepo segment.Repository
}

// NewListSegmentsHandler creates a new ListSegmentsHandler.
func NewListSegmentsHandler(segmentRepo segment.Repository) (ListSegmentsHandler, error) {
	if segmentRepo == nil {
		return listSegmentsHandler{}, errors.New("segment repository is not provided")
	}

	return listSegmentsHandler{segmentRepo}, nil
}

// Handle returns all segments.
func (h listSegmentsHandler) Handle(ctx context.Context) ([]segment.Segment, error) {
	segments, err := h.segmentRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list segments: %w", err)
	}

	return segments, nil
}
