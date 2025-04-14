package api

import (
	"avito/internal/models"
	"avito/internal/models/core_errors"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

func (s *Server) PostPvz(w http.ResponseWriter, r *http.Request) {
	var body PVZ
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		slog.Error("json decode error", "error", err.Error())
		WriteBadRequest(w, "invalid request body")
		return
	}

	city := models.PvzCity(body.City)
	id, err := uuid.Parse(body.Id.String())
	if err != nil {
		slog.Error("uuid parse error", "error", err.Error())
		WriteBadRequest(w, "invalid request")
		return
	}
	pvz := &models.Pvz{
		ID:        id,
		City:      city,
		CreatedAt: *body.RegistrationDate,
	}
	pvz, err = s.pvzService.Create(r.Context(), pvz)

	if err != nil {
		slog.Error("create pvz error", "error", err.Error())
		switch {
		case errors.Is(err, core_errors.ErrAccessDenied):
			WriteForbidden(w, "access denied")
		default:
			WriteBadRequest(w, "invalid request")
		}
		return
	}

	resp := PVZ{
		Id:               &pvz.ID,
		City:             PVZCity(pvz.City),
		RegistrationDate: &pvz.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}

func (s *Server) GetPvz(w http.ResponseWriter, r *http.Request, params GetPvzParams) {
	filter := models.Filter{
		StartDate: params.StartDate,
		EndDate:   params.EndDate,
		Page:      *params.Page,
		Limit:     *params.Limit,
	}
	resp, err := s.pvzService.GetWithFilter(r.Context(), filter)
	if err != nil {
		slog.Error("get pvz error", "error", err.Error())
		switch {
		case errors.Is(err, core_errors.ErrAccessDenied):
			WriteForbidden(w, "access denied")
		default:
			WriteBadRequest(w, "invalid request")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(&resp)
}
