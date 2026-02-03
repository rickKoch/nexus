package port

import (
	"encoding/json"
	"net/http"

	"github.com/rickKoch/nexus/internal/segments/app"
	"github.com/rickKoch/nexus/internal/segments/app/segments"
	"github.com/rickKoch/nexus/internal/segments/domain/segment"
)

type HttpServer struct {
	app app.Application
}

func NewHttpServer(application app.Application) HttpServer {
	return HttpServer{
		app: application,
	}
}

// GetSegment handles GET /segment/:id
func (h HttpServer) GetSegment(w http.ResponseWriter, r *http.Request, params GetSegmentParams) {
	seg, err := h.app.Segments.GetSegment.Handle(r.Context(), segments.GetSegment{ID: params.ID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	render(w, http.StatusOK, toSegmentResponse(seg))
}

// ListSegments handles GET /segment
func (h HttpServer) ListSegments(w http.ResponseWriter, r *http.Request) {
	segs, err := h.app.Segments.ListSegments.Handle(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := make([]SegmentResponse, 0, len(segs))
	for i := range segs {
		response = append(response, toSegmentResponse(&segs[i]))
	}

	render(w, http.StatusOK, response)
}

// CreateSegment handles POST /segment
func (h HttpServer) CreateSegment(w http.ResponseWriter, r *http.Request) {
	var req CreateSegmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	seg, err := h.app.Segments.CreateSegment.Handle(r.Context(), segments.CreateSegment{
		Name:       req.Name,
		TTLSeconds: req.TTLSeconds,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	render(w, http.StatusCreated, toSegmentResponse(seg))
}

// UpdateSegment handles PUT /segment/:id
func (h HttpServer) UpdateSegment(w http.ResponseWriter, r *http.Request, params UpdateSegmentParams) {
	var req UpdateSegmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	seg, err := h.app.Segments.UpdateSegment.Handle(r.Context(), segments.UpdateSegment{
		ID:         params.ID,
		Name:       req.Name,
		TTLSeconds: req.TTLSeconds,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	render(w, http.StatusOK, toSegmentResponse(seg))
}

// DeleteSegment handles DELETE /segment/:id
func (h HttpServer) DeleteSegment(w http.ResponseWriter, r *http.Request, params DeleteSegmentParams) {
	err := h.app.Segments.DeleteSegment.Handle(r.Context(), segments.DeleteSegment{ID: params.ID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Request/Response types
type CreateSegmentRequest struct {
	Name       string `json:"name"`
	TTLSeconds *int   `json:"ttl_seconds,omitempty"`
}

type UpdateSegmentRequest struct {
	Name       string `json:"name"`
	TTLSeconds *int   `json:"ttl_seconds,omitempty"`
}

type SegmentResponse struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	TTLSeconds *int   `json:"ttl_seconds,omitempty"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

func render(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func toSegmentResponse(s *segment.Segment) SegmentResponse {
	return SegmentResponse{
		ID:         s.ID(),
		Name:       s.Name(),
		TTLSeconds: s.TTLSeconds(),
		CreatedAt:  s.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  s.UpdatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}
}
