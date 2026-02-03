package segments_test

import (
	"context"
	"testing"

	"github.com/rickKoch/nexus/internal/segments/adapters"
	"github.com/rickKoch/nexus/internal/segments/app/segments"
)

func TestUpdateSegmentHandler_Handle(t *testing.T) {
	repo := adapters.NewInMemorySegmentRepository()
	createHandler, _ := segments.NewCreateSegmentHandler(repo)
	updateHandler, err := segments.NewUpdateSegmentHandler(repo)
	if err != nil {
		t.Fatalf("failed to create handler: %v", err)
	}

	ctx := context.Background()

	t.Run("updates segment name", func(t *testing.T) {
		// Create a segment first
		created, err := createHandler.Handle(ctx, segments.CreateSegment{
			Name: "original-name",
		})
		if err != nil {
			t.Fatalf("failed to create segment: %v", err)
		}

		// Update the segment
		updated, err := updateHandler.Handle(ctx, segments.UpdateSegment{
			ID:   created.ID(),
			Name: "updated-name",
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if updated.Name() != "updated-name" {
			t.Errorf("expected name 'updated-name', got '%s'", updated.Name())
		}
	})

	t.Run("updates segment TTL", func(t *testing.T) {
		// Create a segment without TTL
		created, err := createHandler.Handle(ctx, segments.CreateSegment{
			Name: "ttl-update-test",
		})
		if err != nil {
			t.Fatalf("failed to create segment: %v", err)
		}

		// Update with TTL
		newTTL := 7200
		updated, err := updateHandler.Handle(ctx, segments.UpdateSegment{
			ID:         created.ID(),
			Name:       "ttl-update-test",
			TTLSeconds: &newTTL,
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if updated.TTLSeconds() == nil || *updated.TTLSeconds() != 7200 {
			t.Errorf("expected TTL 7200, got %v", updated.TTLSeconds())
		}
	})

	t.Run("fails for non-existent segment", func(t *testing.T) {
		_, err := updateHandler.Handle(ctx, segments.UpdateSegment{
			ID:   99999,
			Name: "non-existent",
		})
		if err == nil {
			t.Error("expected error for non-existent segment")
		}
	})

	t.Run("fails with empty name", func(t *testing.T) {
		created, _ := createHandler.Handle(ctx, segments.CreateSegment{
			Name: "will-fail-update",
		})

		_, err := updateHandler.Handle(ctx, segments.UpdateSegment{
			ID:   created.ID(),
			Name: "",
		})
		if err == nil {
			t.Error("expected error for empty name")
		}
	})

	t.Run("fails with invalid TTL", func(t *testing.T) {
		created, _ := createHandler.Handle(ctx, segments.CreateSegment{
			Name: "invalid-ttl-update",
		})

		invalidTTL := -100
		_, err := updateHandler.Handle(ctx, segments.UpdateSegment{
			ID:         created.ID(),
			Name:       "invalid-ttl-update",
			TTLSeconds: &invalidTTL,
		})
		if err == nil {
			t.Error("expected error for invalid TTL")
		}
	})
}

func TestNewUpdateSegmentHandler_NilRepository(t *testing.T) {
	_, err := segments.NewUpdateSegmentHandler(nil)
	if err == nil {
		t.Error("expected error for nil repository")
	}
}
