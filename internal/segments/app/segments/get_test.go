package segments_test

import (
	"context"
	"testing"

	"github.com/rickKoch/nexus/internal/segments/adapters"
	"github.com/rickKoch/nexus/internal/segments/app/segments"
)

func TestGetSegmentHandler_Handle(t *testing.T) {
	repo := adapters.NewInMemorySegmentRepository()
	createHandler, _ := segments.NewCreateSegmentHandler(repo)
	getHandler, err := segments.NewGetSegmentHandler(repo)
	if err != nil {
		t.Fatalf("failed to create handler: %v", err)
	}

	ctx := context.Background()

	t.Run("gets existing segment", func(t *testing.T) {
		// Create a segment first
		ttl := 3600
		created, err := createHandler.Handle(ctx, segments.CreateSegment{
			Name:       "get-test-segment",
			TTLSeconds: &ttl,
		})
		if err != nil {
			t.Fatalf("failed to create segment: %v", err)
		}

		// Get the segment
		got, err := getHandler.Handle(ctx, segments.GetSegment{ID: created.ID()})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if got.ID() != created.ID() {
			t.Errorf("expected ID %d, got %d", created.ID(), got.ID())
		}
		if got.Name() != "get-test-segment" {
			t.Errorf("expected name 'get-test-segment', got '%s'", got.Name())
		}
	})

	t.Run("fails for non-existent segment", func(t *testing.T) {
		_, err := getHandler.Handle(ctx, segments.GetSegment{ID: 99999})
		if err == nil {
			t.Error("expected error for non-existent segment")
		}
	})
}

func TestNewGetSegmentHandler_NilRepository(t *testing.T) {
	_, err := segments.NewGetSegmentHandler(nil)
	if err == nil {
		t.Error("expected error for nil repository")
	}
}
