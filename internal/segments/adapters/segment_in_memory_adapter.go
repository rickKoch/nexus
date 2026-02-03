package adapters

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/rickKoch/nexus/internal/segments/domain/segment"
)

// ErrSegmentNotFound is returned when a segment is not found.
var ErrSegmentNotFound = errors.New("segment not found")

// InMemorySegmentRepository is an in-memory implementation of segment.Repository.
type InMemorySegmentRepository struct {
	mu       sync.RWMutex
	segments map[int]*segment.Segment
	nextID   int
}

// NewInMemorySegmentRepository creates a new in-memory segment repository.
func NewInMemorySegmentRepository() *InMemorySegmentRepository {
	return &InMemorySegmentRepository{
		segments: make(map[int]*segment.Segment),
		nextID:   1,
	}
}

// List returns paginated non-deleted segments.
func (r *InMemorySegmentRepository) List(ctx context.Context, params segment.ListParams) (*segment.ListResult, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Collect all non-deleted segments
	all := make([]segment.Segment, 0, len(r.segments))
	for _, s := range r.segments {
		if !s.IsDeleted() {
			all = append(all, *s)
		}
	}

	// Sort by ID for consistent ordering
	for i := 0; i < len(all)-1; i++ {
		for j := i + 1; j < len(all); j++ {
			if all[i].ID() > all[j].ID() {
				all[i], all[j] = all[j], all[i]
			}
		}
	}

	totalCount := len(all)

	// Apply pagination
	start := (params.Page - 1) * params.PageSize
	if start >= totalCount {
		return &segment.ListResult{
			Segments:   []segment.Segment{},
			TotalCount: totalCount,
			Page:       params.Page,
			PageSize:   params.PageSize,
		}, nil
	}

	end := start + params.PageSize
	if end > totalCount {
		end = totalCount
	}

	return &segment.ListResult{
		Segments:   all[start:end],
		TotalCount: totalCount,
		Page:       params.Page,
		PageSize:   params.PageSize,
	}, nil
}

// Get returns a segment by ID.
func (r *InMemorySegmentRepository) Get(ctx context.Context, id int) (*segment.Segment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	s, ok := r.segments[id]
	if !ok || s.IsDeleted() {
		return nil, ErrSegmentNotFound
	}

	return s, nil
}

// Create stores a new segment and returns it with an assigned ID.
func (r *InMemorySegmentRepository) Create(ctx context.Context, s *segment.Segment) (*segment.Segment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := r.nextID
	r.nextID++

	now := time.Now()
	newSegment := segment.UnmarshalSegmentFromDatabase(
		id,
		s.Name(),
		s.TTLSeconds(),
		now,
		now,
		nil,
	)

	r.segments[id] = newSegment

	return newSegment, nil
}

// Update updates an existing segment.
func (r *InMemorySegmentRepository) Update(ctx context.Context, s *segment.Segment) (*segment.Segment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, ok := r.segments[s.ID()]
	if !ok || existing.IsDeleted() {
		return nil, ErrSegmentNotFound
	}

	updatedSegment := segment.UnmarshalSegmentFromDatabase(
		s.ID(),
		s.Name(),
		s.TTLSeconds(),
		existing.CreatedAt(),
		time.Now(),
		nil,
	)

	r.segments[s.ID()] = updatedSegment

	return updatedSegment, nil
}

// Delete soft-deletes a segment by ID.
func (r *InMemorySegmentRepository) Delete(ctx context.Context, id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	s, ok := r.segments[id]
	if !ok || s.IsDeleted() {
		return ErrSegmentNotFound
	}

	s.Delete(time.Now())

	return nil
}
