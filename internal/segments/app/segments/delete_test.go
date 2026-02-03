package segments_test

import (
	"context"
	"testing"

	"github.com/rickKoch/nexus/internal/segments/adapters"
	"github.com/rickKoch/nexus/internal/segments/app/segments"
)

func TestDeleteSegmentHandler_Handle(t *testing.T) {
	repo := adapters.NewInMemorySegmentRepository()
	createHandler, _ := segments.NewCreateSegmentHandler(repo)
	getHandler, _ := segments.NewGetSegmentHandler(repo)
	deleteHandler, err := segments.NewDeleteSegmentHandler(repo)
	if err != nil {
		t.Fatalf("failed to create handler: %v", err)
	}

	ctx := context.Background()

	t.Run("deletes existing segment", func(t *testing.T) {
		// Create a segment first
		created, err := createHandler.Handle(ctx, segments.CreateSegment{
			Name: "to-delete",
		})
		if err != nil {
			t.Fatalf("failed to create segment: %v", err)
		}

		// Delete the segment
		err = deleteHandler.Handle(ctx, segments.DeleteSegment{ID: created.ID()})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Verify it's deleted (should not be retrievable)
		_, err = getHandler.Handle(ctx, segments.GetSegment{ID: created.ID()})
		if err == nil {
			t.Error("expected error when getting deleted segment")
		}
	})

	t.Run("fails for non-existent segment", func(t *testing.T) {
		err := deleteHandler.Handle(ctx, segments.DeleteSegment{ID: 99999})
		if err == nil {
			t.Error("expected error for non-existent segment")
		}
	})

	t.Run("fails when deleting already deleted segment", func(t *testing.T) {
		// Create and delete a segment
		created, _ := createHandler.Handle(ctx, segments.CreateSegment{
			Name: "double-delete",
		})
		_ = deleteHandler.Handle(ctx, segments.DeleteSegment{ID: created.ID()})

		// Try to delete again
		err := deleteHandler.Handle(ctx, segments.DeleteSegment{ID: created.ID()})
		if err == nil {
			t.Error("expected error when deleting already deleted segment")
		}
	})
}

func TestNewDeleteSegmentHandler_NilRepository(t *testing.T) {
	_, err := segments.NewDeleteSegmentHandler(nil)
	if err == nil {
		t.Error("expected error for nil repository")
	}
}
