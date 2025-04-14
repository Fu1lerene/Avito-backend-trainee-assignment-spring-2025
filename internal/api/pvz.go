package api

import (
	"avito/internal/models"
	core_errors "avito/internal/utils/errors"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"net/http"
)

func (s *Server) PostPvz(w http.ResponseWriter, r *http.Request) {
	var body PVZ
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		core_errors.Write(w, core_errors.ErrInvalidRequest, http.StatusBadRequest)
		return
	}

	city := models.PvzCity(body.City)
	id, err := uuid.Parse(body.Id.String())
	if err != nil {
		core_errors.Write(w, core_errors.ErrInvalidRequest, http.StatusBadRequest)
		return
	}
	pvz := &models.Pvz{
		ID:        id,
		City:      city,
		CreatedAt: *body.RegistrationDate,
	}
	pvz, err = s.pvzService.Create(r.Context(), pvz)

	if err != nil {
		switch {
		case errors.Is(err, core_errors.ErrAccessDenied):
			core_errors.Write(w, err, http.StatusForbidden)
		default:
			core_errors.Write(w, err, http.StatusBadRequest)
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
		switch {
		case errors.Is(err, core_errors.ErrAccessDenied):
			core_errors.Write(w, err, http.StatusForbidden)
		default:
			core_errors.Write(w, err, http.StatusBadRequest)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(&resp)
}
