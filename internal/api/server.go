package api

import (
	"avito/internal/services/auth"
	"avito/internal/services/product"
	"avito/internal/services/pvz"
	"avito/internal/services/reception"
)

var _ ServerInterface = (*Server)(nil)

type Server struct {
	authService      auth.Service
	productService   product.Service
	receptionService reception.Service
	pvzService       pvz.Service
}

func NewServer(
	authService auth.Service,
	productService product.Service,
	receptionService reception.Service,
	pvzService pvz.Service) *Server {
	return &Server{
		authService:      authService,
		productService:   productService,
		receptionService: receptionService,
		pvzService:       pvzService,
	}
}
