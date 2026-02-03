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

		list, err := freshListHandler.Handle(ctx)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(list) != 0 {
			t.Errorf("expected empty list, got %d items", len(list))
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

		list, err := listHandler.Handle(ctx)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(list) < 2 {
			t.Errorf("expected at least 2 segments, got %d", len(list))
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

		list, err := freshListHandler.Handle(ctx)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(list) != 1 {
			t.Errorf("expected 1 segment, got %d", len(list))
		}

		if list[0].Name() != "still-exists" {
			t.Errorf("expected 'still-exists', got '%s'", list[0].Name())
		}
	})
}

func TestNewListSegmentsHandler_NilRepository(t *testing.T) {
	_, err := segments.NewListSegmentsHandler(nil)
	if err == nil {
		t.Error("expected error for nil repository")
	}
}
