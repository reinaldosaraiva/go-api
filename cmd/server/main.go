package main

import (
	"net/http"

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
	http.HandleFunc("/products", ProductHandler.CreateProduct)
	http.ListenAndServe(":"+config.WebServerPort, nil)
}
