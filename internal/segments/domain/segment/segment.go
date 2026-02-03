package segment

import (
	"errors"
	"fmt"
	"time"
)

var (
	// ErrNameRequired is returned when the segment name is empty.
	ErrNameRequired = errors.New("segment name is required")
	// ErrNameTooLong is returned when the segment name exceeds the maximum length.
	ErrNameTooLong = errors.New("segment name must be 255 characters or less")
	// ErrInvalidTTL is returned when the TTL is not positive.
	ErrInvalidTTL = errors.New("TTL must be a positive number")
)

// Segment represents a segment entity in the domain.
type Segment struct {
	id int

	name       string
	ttlSeconds *int

	createdAt time.Time
	updatedAt time.Time
	deletedAt *time.Time
}

// CreateSegment holds the required parameters for creating a new Segment.
type SegmentConfig struct {
	Name       string
	TTLSeconds *int
}

// Validate checks if the CreateSegment fields are valid.
func (c SegmentConfig) Validate() error {
	if c.Name == "" {
		return ErrNameRequired
	}

	if len(c.Name) > 255 {
		return ErrNameTooLong
	}

	if c.TTLSeconds != nil && *c.TTLSeconds <= 0 {
		return ErrInvalidTTL
	}

	return nil
}

type Factory struct {
	sc SegmentConfig
}

func NewFactory(sc SegmentConfig) (Factory, error) {
	if err := sc.Validate(); err != nil {
		return Factory{}, fmt.Errorf("invalid segment config: %w", err)
	}

	return Factory{sc}, nil
}

// NewSegment creates a new Segment with the factory's configuration.
func (f Factory) NewSegment() *Segment {
	now := time.Now()
	return &Segment{
		name:       f.sc.Name,
		ttlSeconds: f.sc.TTLSeconds,
		createdAt:  now,
		updatedAt:  now,
	}
}

// UnmarshalSegmentFromDatabase reconstructs a Segment from database fields.
func UnmarshalSegmentFromDatabase(
	id int,
	name string,
	ttlSeconds *int,
	createdAt time.Time,
	updatedAt time.Time,
	deletedAt *time.Time,
) *Segment {
	return &Segment{
		id:         id,
		name:       name,
		ttlSeconds: ttlSeconds,
		createdAt:  createdAt,
		updatedAt:  updatedAt,
		deletedAt:  deletedAt,
	}
}

// ID returns the segment's unique identifier.
func (s *Segment) ID() int { return s.id }

// Name returns the segment's name.
func (s *Segment) Name() string { return s.name }

// TTLSeconds returns the segment's TTL in seconds.
func (s *Segment) TTLSeconds() *int { return s.ttlSeconds }

// CreatedAt returns when the segment was created.
func (s *Segment) CreatedAt() time.Time { return s.createdAt }

// UpdatedAt returns when the segment was last updated.
func (s *Segment) UpdatedAt() time.Time { return s.updatedAt }

// DeletedAt returns when the segment was deleted, or nil if not deleted.
func (s *Segment) DeletedAt() *time.Time { return s.deletedAt }

// Delete marks the segment as deleted.
func (s *Segment) Delete(deletedAt time.Time) {
	s.deletedAt = &deletedAt
	s.updatedAt = deletedAt
}

// Update modifies the segment's mutable fields.
func (s *Segment) Update(name string, ttlSeconds *int, updatedAt time.Time) {
	s.name = name
	s.ttlSeconds = ttlSeconds
	s.updatedAt = updatedAt
}

// IsDeleted returns true if the segment has been deleted.
func (s *Segment) IsDeleted() bool { return s.deletedAt != nil }
