package api

import (
	"avito/internal/models"
	"avito/internal/models/core_errors"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

func (s *Server) PostDummyLogin(w http.ResponseWriter, r *http.Request) {
	var body PostDummyLoginJSONBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		slog.Error("json decode error", "error", err.Error())
		WriteBadRequest(w, "invalid request body")
		return
	}

	role := models.UserRole(body.Role)

	token, err := s.authService.DummyLogin(role)
	if err != nil {
		slog.Error("dummylogin error", "error", err.Error())
		WriteBadRequest(w, "invalid request")
		return
	}

	_ = json.NewEncoder(w).Encode(token)
}

func (s *Server) PostRegister(w http.ResponseWriter, r *http.Request) {
	var body PostRegisterJSONBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		slog.Error("json decode error", "error", err.Error())
		WriteBadRequest(w, "invalid request body")
		return
	}

	role := models.UserRole(body.Role)
	err := s.authService.Register(r.Context(), role, body.Password, string(body.Email))
	if err != nil {
		slog.Error("register error", "error", err.Error())
		WriteBadRequest(w, "invalid request")
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) PostLogin(w http.ResponseWriter, r *http.Request) {
	var body PostLoginJSONBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		slog.Error("json decode error", "error", err.Error())
		WriteBadRequest(w, "invalid request body")
		return
	}

	token, err := s.authService.Login(r.Context(), body.Password, string(body.Email))
	if err != nil {
		slog.Error("login error", "error", err.Error())
		switch {
		case errors.Is(err, core_errors.ErrInvalidCredentials):
			WriteUnauthorized(w, "invalid credentials")
		default:
			WriteBadRequest(w, "invalid request")
		}
		return
	}

	_ = json.NewEncoder(w).Encode(token)
}
