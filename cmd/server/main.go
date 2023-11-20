package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/reinaldosaraiva/go-api/configs"
	"github.com/reinaldosaraiva/go-api/internal/entity"
	"github.com/reinaldosaraiva/go-api/internal/infra/database"
	"github.com/reinaldosaraiva/go-api/internal/infra/webserver/handlers"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	config, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(postgres.Open(config.GetDBDSN()), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&entity.User{},&entity.Product{})
	ProductHandler := handlers.NewProductHandler(database.NewProduct(db))
	//Debug config values JWTExpireIn
	log.Println("Config JWTExpireIN: " + strconv.FormatUint(uint64(config.JWTExpireIn), 10))
	userHandler := handlers.NewUserHandler(database.NewUser(db),config.TokenAuth, config.JWTExpireIn)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Route("/products", func(r chi.Router) {	
		r.Use(jwtauth.Verifier(config.TokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Post("/", ProductHandler.CreateProduct)
		r.Get("/{id}", ProductHandler.GetProduct)
		r.Get("/", ProductHandler.GetProducts)
		r.Put("/{id}", ProductHandler.UpdateProduct)
		r.Delete("/{id}", ProductHandler.DeleteProduct)
	})
	r.Post("/users", userHandler.CreateUser)

	r.Post("/users/generate_token", userHandler.GetJWT)

	http.ListenAndServe(":"+config.WebServerPort, r)
}
