package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"shedstat/internal/core/domain"
	"shedstat/internal/core/services"

	"github.com/go-chi/chi/v5"
	"github.com/goccy/go-json"
)

type ProfileHTTPHandler struct {
	svc *services.ProfileService
}

func NewProfileHTTPHandler(svc *services.ProfileService) *ProfileHTTPHandler {
	return &ProfileHTTPHandler{
		svc: svc,
	}
}

func (h *ProfileHTTPHandler) Setup(r *chi.Mux) {
	r.Route("/profile", func(r chi.Router) {
		r.Post("/", h.Create)
		r.Get("/{id}", h.GetByShedevrumID)
		r.Get("/{id}/metrics", h.GetMetrics)
	})
}

type ProfileCreateBody struct {
	ShedevrumID string `json:"shedevrum_id"`
}

func (h *ProfileHTTPHandler) Create(w http.ResponseWriter, r *http.Request) {
	var body ProfileCreateBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err := h.svc.Create(r.Context(), body.ShedevrumID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *ProfileHTTPHandler) GetByShedevrumID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	profile, err := h.svc.GetByShedevrumID(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(profile)
}

func (h *ProfileHTTPHandler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	metrics, err := h.svc.GetMetrics(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(metrics)
}

func (h *ProfileHTTPHandler) GetTop(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	qFilter := query.Get("filter")

	var filter domain.ProfileMetrics_GetTopFilter
	filter.Scan(qFilter)
}
