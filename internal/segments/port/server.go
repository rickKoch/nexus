package port

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// (GET /segment/:id)
	GetSegment(w http.ResponseWriter, r *http.Request, params GetSegmentParams)

	// (GET /segment)
	ListSegments(w http.ResponseWriter, r *http.Request)

	// (POST /segment)
	CreateSegment(w http.ResponseWriter, r *http.Request)

	// (PUT /segment/:id)
	UpdateSegment(w http.ResponseWriter, r *http.Request, params UpdateSegmentParams)

	// (DELETE /segment/:id)
	DeleteSegment(w http.ResponseWriter, r *http.Request, params DeleteSegmentParams)
}

func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	if r == nil {
		r = chi.NewRouter()
	}
	wrapper := ServerInterfaceWrapper{
		Handler: si,
		ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		},
	}

	r.Group(func(r chi.Router) {
		r.Get("/segment", wrapper.ListSegments)
		r.Post("/segment", wrapper.CreateSegment)
		r.Get("/segment/{id}", wrapper.GetSegment)
		r.Put("/segment/{id}", wrapper.UpdateSegment)
		r.Delete("/segment/{id}", wrapper.DeleteSegment)
	})

	return r
}

type ServerInterfaceWrapper struct {
	Handler          ServerInterface
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

func (siw *ServerInterfaceWrapper) GetSegment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, err)
		return
	}

	params := GetSegmentParams{ID: id}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetSegment(w, r, params)
	})

	handler.ServeHTTP(w, r.WithContext(ctx))
}

func (siw *ServerInterfaceWrapper) ListSegments(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.ListSegments(w, r)
	})

	handler.ServeHTTP(w, r.WithContext(ctx))
}

func (siw *ServerInterfaceWrapper) CreateSegment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CreateSegment(w, r)
	})

	handler.ServeHTTP(w, r.WithContext(ctx))
}

func (siw *ServerInterfaceWrapper) UpdateSegment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, err)
		return
	}

	params := UpdateSegmentParams{ID: id}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.UpdateSegment(w, r, params)
	})

	handler.ServeHTTP(w, r.WithContext(ctx))
}

func (siw *ServerInterfaceWrapper) DeleteSegment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, err)
		return
	}

	params := DeleteSegmentParams{ID: id}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.DeleteSegment(w, r, params)
	})

	handler.ServeHTTP(w, r.WithContext(ctx))
}

type GetSegmentParams struct {
	ID int `json:"id"`
}

type UpdateSegmentParams struct {
	ID int `json:"id"`
}

type DeleteSegmentParams struct {
	ID int `json:"id"`
}
