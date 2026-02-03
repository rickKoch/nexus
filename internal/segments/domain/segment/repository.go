package segment

import "context"

type Repository interface {
	List(ctx context.Context) ([]Segment, error)
	Get(ctx context.Context, id int) (*Segment, error)
	Create(ctx context.Context, segment *Segment) (*Segment, error)
	Update(ctx context.Context, segment *Segment) (*Segment, error)
	Delete(ctx context.Context, id int) error
}
