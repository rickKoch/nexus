package adapters

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rickKoch/nexus/internal/segments/domain/segment"
)

// segmentRow represents a database row for a segment.
type segmentRow struct {
	ID         int        `db:"id"`
	Name       string     `db:"name"`
	TTLSeconds *int       `db:"ttl_seconds"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at"`
	DeletedAt  *time.Time `db:"deleted_at"`
}

// PostgreSQLSegmentRepository is a PostgreSQL implementation of segment.Repository.
type PostgreSQLSegmentRepository struct {
	db *sqlx.DB
}

// NewPostgreSQLSegmentRepository creates a new PostgreSQL segment repository.
func NewPostgreSQLSegmentRepository(db *sqlx.DB) *PostgreSQLSegmentRepository {
	return &PostgreSQLSegmentRepository{
		db: db,
	}
}

// List returns all non-deleted segments.
func (r *PostgreSQLSegmentRepository) List(ctx context.Context) ([]segment.Segment, error) {
	query := `
		SELECT id, name, ttl_seconds, created_at, updated_at, deleted_at
		FROM segments
		WHERE deleted_at IS NULL
		ORDER BY id
	`

	var rows []segmentRow
	if err := r.db.SelectContext(ctx, &rows, query); err != nil {
		return nil, err
	}

	segments := make([]segment.Segment, 0, len(rows))
	for _, row := range rows {
		s := segment.UnmarshalSegmentFromDatabase(
			row.ID, row.Name, row.TTLSeconds, row.CreatedAt, row.UpdatedAt, row.DeletedAt,
		)
		segments = append(segments, *s)
	}

	return segments, nil
}

// Get returns a segment by ID.
func (r *PostgreSQLSegmentRepository) Get(ctx context.Context, id int) (*segment.Segment, error) {
	query := `
		SELECT id, name, ttl_seconds, created_at, updated_at, deleted_at
		FROM segments
		WHERE id = $1 AND deleted_at IS NULL
	`

	var row segmentRow
	if err := r.db.GetContext(ctx, &row, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSegmentNotFound
		}
		return nil, err
	}

	return segment.UnmarshalSegmentFromDatabase(
		row.ID, row.Name, row.TTLSeconds, row.CreatedAt, row.UpdatedAt, row.DeletedAt,
	), nil
}

// Create stores a new segment and returns it with an assigned ID.
func (r *PostgreSQLSegmentRepository) Create(ctx context.Context, s *segment.Segment) (*segment.Segment, error) {
	query := `
		INSERT INTO segments (name, ttl_seconds, created_at, updated_at)
		VALUES (:name, :ttl_seconds, :created_at, :updated_at)
		RETURNING id, name, ttl_seconds, created_at, updated_at, deleted_at
	`

	now := time.Now()
	params := segmentRow{
		Name:       s.Name(),
		TTLSeconds: s.TTLSeconds(),
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	rows, err := r.db.NamedQueryContext(ctx, query, params)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var row segmentRow
	if rows.Next() {
		if err := rows.StructScan(&row); err != nil {
			return nil, err
		}
	}

	return segment.UnmarshalSegmentFromDatabase(
		row.ID, row.Name, row.TTLSeconds, row.CreatedAt, row.UpdatedAt, row.DeletedAt,
	), nil
}

// Update updates an existing segment.
func (r *PostgreSQLSegmentRepository) Update(ctx context.Context, s *segment.Segment) (*segment.Segment, error) {
	query := `
		UPDATE segments
		SET name = :name, ttl_seconds = :ttl_seconds, updated_at = :updated_at
		WHERE id = :id AND deleted_at IS NULL
		RETURNING id, name, ttl_seconds, created_at, updated_at, deleted_at
	`

	params := segmentRow{
		ID:         s.ID(),
		Name:       s.Name(),
		TTLSeconds: s.TTLSeconds(),
		UpdatedAt:  time.Now(),
	}

	rows, err := r.db.NamedQueryContext(ctx, query, params)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var row segmentRow
	if !rows.Next() {
		return nil, ErrSegmentNotFound
	}

	if err := rows.StructScan(&row); err != nil {
		return nil, err
	}

	return segment.UnmarshalSegmentFromDatabase(
		row.ID, row.Name, row.TTLSeconds, row.CreatedAt, row.UpdatedAt, row.DeletedAt,
	), nil
}

// Delete soft-deletes a segment by ID.
func (r *PostgreSQLSegmentRepository) Delete(ctx context.Context, id int) error {
	query := `
		UPDATE segments
		SET deleted_at = :deleted_at, updated_at = :updated_at
		WHERE id = :id AND deleted_at IS NULL
	`

	now := time.Now()
	params := map[string]interface{}{
		"id":         id,
		"deleted_at": now,
		"updated_at": now,
	}

	result, err := r.db.NamedExecContext(ctx, query, params)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrSegmentNotFound
	}

	return nil
}
