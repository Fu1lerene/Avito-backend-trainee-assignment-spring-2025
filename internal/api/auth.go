package api

import (
	"avito/internal/models"
	core_errors "avito/internal/utils/errors"
	"encoding/json"
	"errors"
	"net/http"
)

func (s *Server) PostDummyLogin(w http.ResponseWriter, r *http.Request) {
	var body PostDummyLoginJSONBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		core_errors.Write(w, core_errors.ErrInvalidRequest, http.StatusBadRequest)
		return
	}

	role := models.UserRole(body.Role)

	token, err := s.authService.DummyLogin(role)
	if err != nil {
		core_errors.Write(w, err, http.StatusBadRequest)
		return
	}

	_ = json.NewEncoder(w).Encode(token)
}

func (s *Server) PostRegister(w http.ResponseWriter, r *http.Request) {
	var body PostRegisterJSONBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		core_errors.Write(w, core_errors.ErrInvalidRequest, http.StatusBadRequest)
		return
	}

	role := models.UserRole(body.Role)
	_, err := s.authService.Register(r.Context(), role, body.Password, string(body.Email))
	if err != nil {
		core_errors.Write(w, err, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) PostLogin(w http.ResponseWriter, r *http.Request) {
	var body PostLoginJSONBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		core_errors.Write(w, core_errors.ErrInvalidRequest, http.StatusBadRequest)
		return
	}

	token, err := s.authService.Login(r.Context(), body.Password, string(body.Email))
	if err != nil {
		switch {
		case errors.Is(err, core_errors.ErrInvalidCredentials):
			core_errors.Write(w, err, http.StatusUnauthorized)
		default:
			core_errors.Write(w, err, http.StatusBadRequest)
		}
		return
	}

	_ = json.NewEncoder(w).Encode(token)
}
