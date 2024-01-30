package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/ricardoraposo/api-again/config"
	"github.com/ricardoraposo/api-again/internal/database"
	"github.com/ricardoraposo/api-again/internal/entity"
	"github.com/ricardoraposo/api-again/internal/handlers"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	c, err := config.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&entity.Product{}, &entity.User{})

	productDB := database.NewProduct(db)
	userDB := database.NewUser(db)

	productHandler := handlers.NewProductHandler(productDB)
	userHandler := handlers.NewUserHandler(userDB, c.TokenAuth, c.JwtExpiresIn)

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Route("/products", func(r chi.Router) {
        r.Use(jwtauth.Verifier(c.TokenAuth))
        r.Use(jwtauth.Authenticator)

		r.Post("/", productHandler.CreateProduct)
		r.Get("/", productHandler.GetProducts)
		r.Get("/{id}", productHandler.GetProduct)
		r.Put("/{id}", productHandler.UpdateProduct)
		r.Delete("/{id}", productHandler.DeleteProduct)
	})

	r.Post("/users", userHandler.CreateUser)
	r.Post("/users/login", userHandler.GetJWT)

	http.ListenAndServe(":8000", r)
}
