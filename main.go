package main

import (
	"log"
	"os"

	"github.com/Kontentski/burgersDb/handlers"
	"github.com/Kontentski/burgersDb/storage"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func main() {
	storage.Init()
	runMigrations()

	r := gin.Default()

	api := r.Group("api/")
	api.GET("/burgers", handlers.GetBurgers)
	api.GET("/burgers/:id", handlers.GetBurgerById)
	api.GET("/burgers/n=:name", handlers.GetBurgerByName)
	api.GET("/burgers/f=:name", handlers.GetBurgerByLetter)
	api.GET("/burgers/random", handlers.GetBurgerByRandom)
	api.GET("/burgers/randomten", handlers.GetBurgersByRandom)
	api.GET("/burgers/latest", handlers.GetLatestBurgers)
	api.POST("/burgers", handlers.CreateBurger)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}

	log.Fatal(r.Run(":" + port))
}

func runMigrations() {
	dsn := os.Getenv("DATABASE_URL")
	db, err := goose.OpenDBWithDriver("pgx", dsn)
	if err != nil {
		log.Fatalf("Could not open db: %v\n", err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("Could not set dialect: %v\n", err)
	}

	if err := goose.Up(db, "./migrations"); err != nil {
		log.Fatalf("Could not run up migrations: %v\n", err)
	}
}
