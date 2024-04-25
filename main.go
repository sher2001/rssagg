package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/sher2001/rss-aggregator/internal/database"

	_ "github.com/lib/pq"
	// underscore in the import url is to let compiler know that don't conside this as unused
	// import even if its not used
	// in Go we need to have drivers as import
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	// If you call Load without any args it will default to loading .env in the current path.
	// You can otherwise tell it which files to load (there can be more than one) like:
	godotenv.Load(".env")

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("port is not specified")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("db URL is not specified")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Unable to connect to database: ", err)
	}

	apiCFG := apiConfig{
		DB: database.New(conn),
	}

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	// v1Router.HandleFunc("/healthz", handlerReadiness)
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handlerError)

	v1Router.Post("/users", apiCFG.handlerCreateUser)
	v1Router.Get("/users", apiCFG.middlewareAuth(apiCFG.handlerGetUserByAPIKey))

	v1Router.Post("/feeds", apiCFG.middlewareAuth(apiCFG.handlerCreateFeed))
	v1Router.Get("/feeds", apiCFG.handlerGetFeeds)

	v1Router.Post("/feed_follows", apiCFG.middlewareAuth(apiCFG.handlerCreateFeedFollow))
	v1Router.Get("/feed_follows", apiCFG.middlewareAuth(apiCFG.handlerGetFeedFollows))
	v1Router.Delete("/feed_follows/{feedFollowId}", apiCFG.middlewareAuth(apiCFG.handlerDeleteFeedFollow))

	router.Mount("/v1", v1Router)

	server := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	fmt.Printf("Server is starting at : %v", port)
	err = server.ListenAndServe()
	// from line 32, code will just listen and serve user requests if anything goes wrong it goes further in this file
	if err != nil {
		log.Fatal(err)
	}
}
