package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/reinaldosaraiva/go-api/configs"
	_ "github.com/reinaldosaraiva/go-api/docs"
	"github.com/reinaldosaraiva/go-api/internal/entity"
	"github.com/reinaldosaraiva/go-api/internal/infra/database"
	"github.com/reinaldosaraiva/go-api/internal/infra/webserver/handlers"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title	API Go
// @version	1.0
// @description	API Go with authentication JWT and CRUD
// @termsOfService	http://swagger.io/terms
// @contact.name	Reinaldo Saraiva
// @contact.email	reinaldo.saraiva@gmail.com
// license.name		Bearware
// host				localhost:8000
// @BasePath	/
// @securityDefinitions.apikey	ApiKeyAuth
// @in	header
// @name	Authorization
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
	userHandler := handlers.NewUserHandler(database.NewUser(db))
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.WithValue("jwt", config.TokenAuth))
	r.Use(middleware.WithValue("jwtExpiresIn", config.JWTExpireIn))
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
	r.Get("/users", userHandler.GetUserByEmail)

	r.Post("/users/generate_token", userHandler.GetJWT)

	r.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8000/docs/swagger/doc.json")))

	http.ListenAndServe(":"+config.WebServerPort, r)
}
