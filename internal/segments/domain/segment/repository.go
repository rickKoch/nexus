package segment

import "context"

// ListParams contains pagination parameters for listing segments.
type ListParams struct {
	Page     int
	PageSize int
}

// ListResult contains the paginated list of segments and total count.
type ListResult struct {
	Segments   []Segment
	TotalCount int
	Page       int
	PageSize   int
}

type Repository interface {
	List(ctx context.Context, params ListParams) (*ListResult, error)
	Get(ctx context.Context, id int) (*Segment, error)
	Create(ctx context.Context, segment *Segment) (*Segment, error)
	Update(ctx context.Context, segment *Segment) (*Segment, error)
	Delete(ctx context.Context, id int) error
}
