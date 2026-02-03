package segments_test

import (
	"context"
	"testing"

	"github.com/rickKoch/nexus/internal/segments/adapters"
	"github.com/rickKoch/nexus/internal/segments/app/segments"
)

func TestCreateSegmentHandler_Handle(t *testing.T) {
	repo := adapters.NewInMemorySegmentRepository()
	handler, err := segments.NewCreateSegmentHandler(repo)
	if err != nil {
		t.Fatalf("failed to create handler: %v", err)
	}

	ctx := context.Background()

	t.Run("creates segment successfully", func(t *testing.T) {
		ttl := 3600
		created, err := handler.Handle(ctx, segments.CreateSegment{
			Name:       "test-segment",
			TTLSeconds: &ttl,
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if created.ID() == 0 {
			t.Error("expected non-zero ID")
		}
		if created.Name() != "test-segment" {
			t.Errorf("expected name 'test-segment', got '%s'", created.Name())
		}
		if created.TTLSeconds() == nil || *created.TTLSeconds() != 3600 {
			t.Errorf("expected TTL 3600, got %v", created.TTLSeconds())
		}
	})

	t.Run("creates segment without TTL", func(t *testing.T) {
		created, err := handler.Handle(ctx, segments.CreateSegment{
			Name: "no-ttl-segment",
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if created.Name() != "no-ttl-segment" {
			t.Errorf("expected name 'no-ttl-segment', got '%s'", created.Name())
		}
		if created.TTLSeconds() != nil {
			t.Errorf("expected nil TTL, got %v", created.TTLSeconds())
		}
	})

	t.Run("fails with empty name", func(t *testing.T) {
		_, err := handler.Handle(ctx, segments.CreateSegment{
			Name: "",
		})
		if err == nil {
			t.Error("expected error for empty name")
		}
	})

	t.Run("fails with invalid TTL", func(t *testing.T) {
		invalidTTL := -1
		_, err := handler.Handle(ctx, segments.CreateSegment{
			Name:       "invalid-ttl",
			TTLSeconds: &invalidTTL,
		})
		if err == nil {
			t.Error("expected error for invalid TTL")
		}
	})
}

func TestNewCreateSegmentHandler_NilRepository(t *testing.T) {
	_, err := segments.NewCreateSegmentHandler(nil)
	if err == nil {
		t.Error("expected error for nil repository")
	}
}
