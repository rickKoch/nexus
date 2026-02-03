package segments

import (
	"context"
	"errors"
	"fmt"

	"github.com/rickKoch/nexus/internal/segments/domain/segment"
)

const (
	DefaultPage     = 1
	DefaultPageSize = 20
	MaxPageSize     = 100
)

// ListSegments contains the pagination parameters for listing segments.
type ListSegments struct {
	Page     int
	PageSize int
}

// ListSegmentsResult contains the paginated list of segments.
type ListSegmentsResult struct {
	Segments   []segment.Segment
	TotalCount int
	Page       int
	PageSize   int
	TotalPages int
}

// ListSegmentsHandler defines the interface for listing segments.
type ListSegmentsHandler interface {
	Handle(ctx context.Context, cmd ListSegments) (*ListSegmentsResult, error)
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

// Handle returns paginated segments.
func (h listSegmentsHandler) Handle(ctx context.Context, cmd ListSegments) (*ListSegmentsResult, error) {
	// Apply defaults and constraints
	page := cmd.Page
	if page < 1 {
		page = DefaultPage
	}

	pageSize := cmd.PageSize
	if pageSize < 1 {
		pageSize = DefaultPageSize
	}
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}

	result, err := h.segmentRepo.List(ctx, segment.ListParams{
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list segments: %w", err)
	}

	totalPages := (result.TotalCount + pageSize - 1) / pageSize
	if totalPages < 1 {
		totalPages = 1
	}

	return &ListSegmentsResult{
		Segments:   result.Segments,
		TotalCount: result.TotalCount,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}
