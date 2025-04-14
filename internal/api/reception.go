package api

import (
	core_errors "avito/internal/utils/errors"
	"encoding/json"
	"errors"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"net/http"
)

func (s *Server) PostReceptions(w http.ResponseWriter, r *http.Request) {
	var body PostReceptionsJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		core_errors.Write(w, core_errors.ErrInvalidRequest, http.StatusBadRequest)
		return
	}

	reception, err := s.receptionService.Open(r.Context(), body.PvzId)
	if err != nil {
		switch {
		case errors.Is(err, core_errors.ErrAccessDenied):
			core_errors.Write(w, err, http.StatusForbidden)
		default:
			core_errors.Write(w, err, http.StatusBadRequest)
		}
		return
	}

	resp := Reception{
		Id:       &reception.ID,
		PvzId:    reception.PvzID,
		Status:   ReceptionStatus(reception.Status),
		DateTime: reception.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}

func (s *Server) PostPvzPvzIdCloseLastReception(w http.ResponseWriter, r *http.Request, pvzId openapi_types.UUID) {
	_, err := s.receptionService.Close(r.Context(), pvzId)
	if err != nil {
		switch {
		case errors.Is(err, core_errors.ErrAccessDenied):
			core_errors.Write(w, err, http.StatusForbidden)
		default:
			core_errors.Write(w, err, http.StatusBadRequest)
		}
		return
	}
}
