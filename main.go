package main

import (
	"log"
	"os"
	"time"

	"github.com/Kontentski/burgersDb/handlers"
	"github.com/Kontentski/burgersDb/middleware"
	"github.com/Kontentski/burgersDb/storage"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	storage.Init()
	storage.RunMigrations()

	middleware.InitializeCache(5*time.Minute, 10*time.Minute)

	r := gin.Default()

	r.Use(middleware.RateLimiter())
	r.Use(middleware.CacheResponse())

	r.Static("/api/assets", "./public/assets")

	api := r.Group("api/")
	api.GET("/burgers", handlers.GetBurgers)
	api.GET("/burgers/:id", handlers.GetBurgerById)
	api.GET("/burgers/n=:name", handlers.GetBurgerByName)
	api.GET("/burgers/f=:name", handlers.GetBurgerByLetter)
	api.GET("/burgers/random", handlers.GetBurgerByRandom)
	api.GET("/burgers/randomten", handlers.GetBurgersByRandom)
	api.GET("/burgers/latest", handlers.GetLatestBurgers)
	api.GET("/ingredients/:name", handlers.GetIngredientByName)
	api.GET("/burgers/i=:name", handlers.GetBurgersByIngredientName)
	api.GET("/burgers/ingredients", handlers.GetBurgersByIngredients)
	api.GET("/burgers/vegan", handlers.GetVeganBurgers)
	api.GET("/burgers/nonvegan", handlers.GetNonVeganBurgers)
	api.POST("/burgers", handlers.CreateBurger)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}

	log.Fatal(r.Run(":" + port))
}
