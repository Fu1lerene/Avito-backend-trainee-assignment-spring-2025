package api

import (
	"avito/internal/models"
	"avito/internal/models/core_errors"
	"encoding/json"
	"errors"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"log/slog"
	"net/http"
)

func (s *Server) PostProducts(w http.ResponseWriter, r *http.Request) {
	var body PostProductsJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		slog.Error("json decode error", "error", err.Error())
		WriteBadRequest(w, "invalid request body")
		return
	}

	productType := models.ProductType(body.Type)
	product, err := s.productService.Add(r.Context(), productType, body.PvzId)
	if err != nil {
		slog.Error("add product error", "error", err.Error())
		switch {
		case errors.Is(err, core_errors.ErrAccessDenied):
			WriteForbidden(w, "access denied")
		default:
			WriteBadRequest(w, "invalid request or no active receptions")
		}
		return
	}

	response := Product{
		Id:          &product.ID,
		ReceptionId: product.ReceptionID,
		Type:        ProductType(product.Type),
		DateTime:    &product.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(response)
}

func (s *Server) PostPvzPvzIdDeleteLastProduct(w http.ResponseWriter, r *http.Request, pvzId openapi_types.UUID) {
	err := s.productService.DeleteLastProduct(r.Context(), pvzId)
	if err != nil {
		slog.Error("delete last product error", "error", err.Error())
		switch {
		case errors.Is(err, core_errors.ErrAccessDenied):
			WriteForbidden(w, "access denied")
		default:
			WriteBadRequest(w, "invalid request or no active receptions or no products to delete")
		}
		return
	}
}
