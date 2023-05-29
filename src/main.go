package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/ddomd/feeder/internal/database"
	"github.com/ddomd/feeder/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main () {
	//Load environment variables from file
	err := godotenv.Load(); if err != nil {
		log.Fatal(err.Error())
	}

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("PORT environment var is not set!")
	}

	dbURL := os.Getenv("DATABASE_URL")

	if dbURL == "" {
		log.Fatal("DATABASE_URL environment var is not set!")
	}

	//Load DB
	db, err := sql.Open("postgres", dbURL); if err != nil {
		log.Fatal(err.Error())
	}

	dbQueries := database.New(db)

	//Setup config
	cfg := config.ApiConfig{dbQueries, port}

	//Router declarations
	router := chi.NewRouter()
	v1Router := chi.NewRouter()

	//CORS options
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	router.Mount("/v1", v1Router)

	//POST Routes
	v1Router.Post("/users", cfg.HandleCreateUser)
	v1Router.Post("/feeds", cfg.AuthMiddleware(cfg.HandleCreateFeed))

	//GET Routes
	v1Router.Get("/users", cfg.AuthMiddleware(cfg.HandleGetUserByAPIKey))
	v1Router.Get("/feeds", cfg.HandleGetAllFeeds)
	v1Router.Get("/err", config.HandleServerError)
	v1Router.Get("/ready", config.HandleServerReady)
	

	server := &http.Server{
		Addr: ":" + port,
		Handler: router,
	}

	err = server.ListenAndServe(); if err != nil {
		log.Fatal(err.Error())
	}
}