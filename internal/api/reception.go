package api

import (
	"avito/internal/models/core_errors"
	"encoding/json"
	"errors"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"log/slog"
	"net/http"
)

func (s *Server) PostReceptions(w http.ResponseWriter, r *http.Request) {
	var body PostReceptionsJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		slog.Error("json decode error", "error", err.Error())
		WriteBadRequest(w, "invalid request body")
		return
	}

	reception, err := s.receptionService.Open(r.Context(), body.PvzId)
	if err != nil {
		slog.Error("open reception error", "error", err.Error())
		switch {
		case errors.Is(err, core_errors.ErrAccessDenied):
			WriteForbidden(w, "access denied")
		default:
			WriteBadRequest(w, "invalid request or there is an open reception")
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
	err := s.receptionService.Close(r.Context(), pvzId)
	if err != nil {
		slog.Error("close reception error", "error", err.Error())
		switch {
		case errors.Is(err, core_errors.ErrAccessDenied):
			WriteForbidden(w, "access denied")
		default:
			WriteBadRequest(w, "invalid request or reception is already closed")
		}
		return
	}
}
