package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"strategy-fox-go-bd/pkg/routes"
)

func main() {
	envPath, err := filepath.Abs("../../.env")
	if err != nil {
		log.Fatalf("Error determining .env path: %v", err)
	}

	// Load the .env file
	err = godotenv.Load(envPath)
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		log.Fatal("ALLOWED_ORIGINS is not defined in the .env file")
	}

	origins := strings.Split(allowedOrigins, ",")

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins: origins,
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})

	router := mux.NewRouter()

	shopifyRouter := router.PathPrefix("/api/shopify").Subrouter()
	chatbotRouter := router.PathPrefix("/api/chatbot").Subrouter()

	routes.ShopifyRoutes(shopifyRouter)
	routes.ChatBotRoutes(chatbotRouter)

	handler := corsMiddleware.Handler(router)

	log.Printf("Server is running on port %s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
