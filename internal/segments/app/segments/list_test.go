package segments_test

import (
	"context"
	"testing"

	"github.com/rickKoch/nexus/internal/segments/adapters"
	"github.com/rickKoch/nexus/internal/segments/app/segments"
)

func TestListSegmentsHandler_Handle(t *testing.T) {
	repo := adapters.NewInMemorySegmentRepository()
	createHandler, _ := segments.NewCreateSegmentHandler(repo)
	listHandler, err := segments.NewListSegmentsHandler(repo)
	if err != nil {
		t.Fatalf("failed to create handler: %v", err)
	}

	ctx := context.Background()

	t.Run("returns empty list when no segments", func(t *testing.T) {
		// Use a fresh repository
		freshRepo := adapters.NewInMemorySegmentRepository()
		freshListHandler, _ := segments.NewListSegmentsHandler(freshRepo)

		result, err := freshListHandler.Handle(ctx, segments.ListSegments{})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(result.Segments) != 0 {
			t.Errorf("expected empty list, got %d items", len(result.Segments))
		}

		if result.TotalCount != 0 {
			t.Errorf("expected total count 0, got %d", result.TotalCount)
		}
	})

	t.Run("returns all created segments", func(t *testing.T) {
		// Create some segments
		_, err := createHandler.Handle(ctx, segments.CreateSegment{Name: "list-segment-1"})
		if err != nil {
			t.Fatalf("failed to create segment: %v", err)
		}
		_, err = createHandler.Handle(ctx, segments.CreateSegment{Name: "list-segment-2"})
		if err != nil {
			t.Fatalf("failed to create segment: %v", err)
		}

		result, err := listHandler.Handle(ctx, segments.ListSegments{})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(result.Segments) < 2 {
			t.Errorf("expected at least 2 segments, got %d", len(result.Segments))
		}
	})

	t.Run("does not return deleted segments", func(t *testing.T) {
		freshRepo := adapters.NewInMemorySegmentRepository()
		freshCreateHandler, _ := segments.NewCreateSegmentHandler(freshRepo)
		freshListHandler, _ := segments.NewListSegmentsHandler(freshRepo)
		deleteHandler, _ := segments.NewDeleteSegmentHandler(freshRepo)

		// Create and delete a segment
		created, _ := freshCreateHandler.Handle(ctx, segments.CreateSegment{Name: "to-be-deleted"})
		_ = deleteHandler.Handle(ctx, segments.DeleteSegment{ID: created.ID()})

		// Create another segment
		_, _ = freshCreateHandler.Handle(ctx, segments.CreateSegment{Name: "still-exists"})

		result, err := freshListHandler.Handle(ctx, segments.ListSegments{})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(result.Segments) != 1 {
			t.Errorf("expected 1 segment, got %d", len(result.Segments))
		}

		if result.Segments[0].Name() != "still-exists" {
			t.Errorf("expected 'still-exists', got '%s'", result.Segments[0].Name())
		}
	})

	t.Run("respects pagination parameters", func(t *testing.T) {
		freshRepo := adapters.NewInMemorySegmentRepository()
		freshCreateHandler, _ := segments.NewCreateSegmentHandler(freshRepo)
		freshListHandler, _ := segments.NewListSegmentsHandler(freshRepo)

		// Create 5 segments
		for i := 1; i <= 5; i++ {
			_, _ = freshCreateHandler.Handle(ctx, segments.CreateSegment{Name: "segment"})
		}

		// Request page 1 with page size 2
		result, err := freshListHandler.Handle(ctx, segments.ListSegments{Page: 1, PageSize: 2})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(result.Segments) != 2 {
			t.Errorf("expected 2 segments, got %d", len(result.Segments))
		}

		if result.TotalCount != 5 {
			t.Errorf("expected total count 5, got %d", result.TotalCount)
		}

		if result.Page != 1 {
			t.Errorf("expected page 1, got %d", result.Page)
		}

		if result.PageSize != 2 {
			t.Errorf("expected page size 2, got %d", result.PageSize)
		}

		if result.TotalPages != 3 {
			t.Errorf("expected total pages 3, got %d", result.TotalPages)
		}

		// Request page 3 with page size 2 (should have 1 item)
		result, err = freshListHandler.Handle(ctx, segments.ListSegments{Page: 3, PageSize: 2})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(result.Segments) != 1 {
			t.Errorf("expected 1 segment on last page, got %d", len(result.Segments))
		}
	})

	t.Run("applies default pagination when not specified", func(t *testing.T) {
		freshRepo := adapters.NewInMemorySegmentRepository()
		freshListHandler, _ := segments.NewListSegmentsHandler(freshRepo)

		result, err := freshListHandler.Handle(ctx, segments.ListSegments{})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if result.Page != segments.DefaultPage {
			t.Errorf("expected default page %d, got %d", segments.DefaultPage, result.Page)
		}

		if result.PageSize != segments.DefaultPageSize {
			t.Errorf("expected default page size %d, got %d", segments.DefaultPageSize, result.PageSize)
		}
	})

	t.Run("enforces max page size", func(t *testing.T) {
		freshRepo := adapters.NewInMemorySegmentRepository()
		freshListHandler, _ := segments.NewListSegmentsHandler(freshRepo)

		result, err := freshListHandler.Handle(ctx, segments.ListSegments{PageSize: 500})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if result.PageSize != segments.MaxPageSize {
			t.Errorf("expected max page size %d, got %d", segments.MaxPageSize, result.PageSize)
		}
	})
}

func TestNewListSegmentsHandler_NilRepository(t *testing.T) {
	_, err := segments.NewListSegmentsHandler(nil)
	if err == nil {
		t.Error("expected error for nil repository")
	}
}
