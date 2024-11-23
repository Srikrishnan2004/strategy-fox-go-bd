package controllers

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"os"
	"strategy-fox-go-bd/pkg/config"
	"time"
)

var ctx = context.Background()

// GetProducts fetches the product list from Shopify
func GetProducts(w http.ResponseWriter, r *http.Request) {

	redisKey := "products"

	// Check if the data is already in Redis
	cachedData, err := config.RedisClient.Get(ctx, redisKey).Result()
	if err == nil {
		// Data exists in Redis, return it
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cachedData))
		fmt.Println("Cache hit: Returning data from Redis")
		return
	}

	// Construct the Shopify URL
	shopifyURL := fmt.Sprintf(
		"https://%s:%s@%s/admin/api/2022-10/products.json",
		os.Getenv("SHOPIFY_API_KEY"),
		os.Getenv("SHOPIFY_ADMIN_API_PASS_TOKEN"),
		os.Getenv("SHOPIFY_STORE_NAME"),
	)

	// Make the GET request to Shopify
	resp, err := http.Get(shopifyURL)
	if err != nil {
		http.Error(w, "Failed to fetch product list", http.StatusInternalServerError)
		fmt.Println("Error fetching products:", err.Error())
		return
	}
	defer resp.Body.Close()

	// Read and parse the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading response", http.StatusInternalServerError)
		fmt.Println("Error reading response:", err.Error())
		return
	}

	err = config.RedisClient.Set(ctx, redisKey, body, 30*time.Second).Err()
	if err != nil {
		fmt.Println("Error caching data in Redis:", err.Error())
	}

	// Write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
	fmt.Println("Cache miss: Fetched data from Shopify API and cached it")
}

// GetProduct fetches a specific product by ID from Shopify
func GetProduct(w http.ResponseWriter, r *http.Request) {
	// Extract the product ID from the route parameters
	vars := mux.Vars(r)
	productID := vars["id"]

	redisKey := fmt.Sprintf("products/%s", productID)

	// Check if the data is already in Redis
	cachedData, err := config.RedisClient.Get(ctx, redisKey).Result()
	if err == nil {
		// Data exists in Redis, return it
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cachedData))
		fmt.Printf("Cache hit: Returning data for product ID %s from Redis\n", productID)
		return
	}

	// Construct the Shopify URL
	shopifyURL := fmt.Sprintf(
		"https://%s:%s@%s/admin/api/2022-10/products/%s.json",
		os.Getenv("SHOPIFY_API_KEY"),
		os.Getenv("SHOPIFY_ADMIN_API_PASS_TOKEN"),
		os.Getenv("SHOPIFY_STORE_NAME"),
		productID,
	)

	// Make the GET request to Shopify
	resp, err := http.Get(shopifyURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch product with ID %s", productID), http.StatusInternalServerError)
		fmt.Printf("Error fetching product with ID %s: %s\n", productID, err.Error())
		return
	}
	defer resp.Body.Close()

	// Read and parse the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading response", http.StatusInternalServerError)
		fmt.Printf("Error reading response for product ID %s: %s\n", productID, err.Error())
		return
	}

	err = config.RedisClient.Set(ctx, redisKey, body, 30*time.Second).Err()
	if err != nil {
		fmt.Printf("Error caching data for product ID %s in Redis: %s\n", productID, err.Error())
	}

	// Write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
	fmt.Printf("Cache miss: Fetched data for product ID %s from Shopify API and cached it\n", productID)
}
