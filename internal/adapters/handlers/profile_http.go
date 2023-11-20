package handlers

import (
	"net/http"
	"shedstat/internal/core/domain"
	"shedstat/internal/core/services"
	"strconv"

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

type ProfileCreateBody struct {
	ShedevrumID string `json:"shedevrum_id"`
	Link        string `json:"link"`
}

func (h *ProfileHTTPHandler) Setup(r *chi.Mux) {
	r.Route("/profile", func(r chi.Router) {
		r.Post("/", h.Create)
		r.Get("/{id}", h.Get)
	})
}

func (h *ProfileHTTPHandler) Create(w http.ResponseWriter, r *http.Request) {
	var body ProfileCreateBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err := h.svc.Create(r.Context(), &domain.ProfileEnity{
		ShedevrumID: body.ShedevrumID,
		Link:        body.Link,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *ProfileHTTPHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idn, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	profile, err := h.svc.Get(r.Context(), idn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(profile)
}
