package app

import (
	"avito/internal/api"
	"avito/internal/infrastructure/database"
	productrepo "avito/internal/repositories/product"
	pvzrepo "avito/internal/repositories/pvz"
	receptionrepo "avito/internal/repositories/reception"
	userrepo "avito/internal/repositories/user"
	"avito/internal/services/auth"
	"avito/internal/services/product"
	"avito/internal/services/pvz"
	"avito/internal/services/reception"
	"avito/internal/utils/jwt"
	"encoding/json"
	"github.com/flowchartsman/swaggerui"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

type app struct {
}

type App interface {
	Run()
}

func New() App {
	return &app{}
}

func (app *app) Run() {
	_ = godotenv.Load()
	jwt.InitSecret()
	db := database.MustInit()

	userRepo := userrepo.New(db)
	pvzRepo := pvzrepo.New(db)
	receptionRepo := receptionrepo.New(db)
	productRepo := productrepo.New(db)

	authService := auth.New(userRepo)
	productService := product.New(receptionRepo, productRepo)
	receptionService := reception.New(receptionRepo)
	pvzService := pvz.New(pvzRepo)

	server := api.NewServer(authService, productService, receptionService, pvzService)

	r := http.NewServeMux()
	h := api.HandlerFromMux(server, r)
	h = AuthMiddleware(h)

	sw, err := api.GetSwagger()
	if err != nil {
		log.Fatal("failed to get swagger spec: ", err)
	}
	spec, err := json.MarshalIndent(sw, "", "  ")
	if err != nil {
		log.Fatal("failed to marshal swagger spec: ", err)
	}

	r.Handle("/swagger/", http.StripPrefix("/swagger", swaggerui.Handler(spec)))

	s := &http.Server{
		Handler: h,
		Addr:    ":8080",
	}
	log.Fatal(s.ListenAndServe())
}
