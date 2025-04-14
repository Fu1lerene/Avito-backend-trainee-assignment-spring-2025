package api

import (
	"avito/internal/models"
	core_errors "avito/internal/utils/errors"
	"encoding/json"
	"errors"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"net/http"
)

func (s *Server) PostProducts(w http.ResponseWriter, r *http.Request) {
	var body PostProductsJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		core_errors.Write(w, core_errors.ErrInvalidRequest, http.StatusBadRequest)
		return
	}

	productType := models.ProductType(body.Type)
	product, err := s.productService.Add(r.Context(), productType, body.PvzId)
	if err != nil {
		switch {
		case errors.Is(err, core_errors.ErrAccessDenied):
			core_errors.Write(w, err, http.StatusForbidden)
		default:
			core_errors.Write(w, err, http.StatusBadRequest)
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
	_, err := s.productService.DeleteLastProduct(r.Context(), pvzId)
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
