package segments

import (
	"context"
	"errors"
	"fmt"

	"github.com/rickKoch/nexus/internal/segments/domain/segment"
)

// DeleteSegment holds the required parameters for deleting a segment.
type DeleteSegment struct {
	ID int
}

// DeleteSegmentHandler defines the interface for deleting a segment.
type DeleteSegmentHandler interface {
	Handle(ctx context.Context, props DeleteSegment) error
}

type deleteSegmentHandler struct {
	segmentRepo segment.Repository
}

// NewDeleteSegmentHandler creates a new DeleteSegmentHandler.
func NewDeleteSegmentHandler(segmentRepo segment.Repository) (DeleteSegmentHandler, error) {
	if segmentRepo == nil {
		return deleteSegmentHandler{}, errors.New("segment repository is not provided")
	}

	return deleteSegmentHandler{segmentRepo}, nil
}

// Handle soft-deletes a segment.
func (h deleteSegmentHandler) Handle(ctx context.Context, props DeleteSegment) error {
	if err := h.segmentRepo.Delete(ctx, props.ID); err != nil {
		return fmt.Errorf("failed to delete segment '%d': %w", props.ID, err)
	}

	return nil
}
