package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"os"
	"strategy-fox-go-bd/pkg/config"
	"time"
)

var ctx = context.Background()

func GetProducts(w http.ResponseWriter, r *http.Request) {

	redisKey := "products"

	cachedData, err := config.RedisClient.Get(ctx, redisKey).Result()
	if err == nil {
		// Data exists in Redis, return it
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cachedData))
		fmt.Println("Cache hit: Returning data from Redis")
		return
	}

	shopifyURL := fmt.Sprintf(
		"https://%s:%s@%s/admin/api/2022-10/products.json",
		os.Getenv("SHOPIFY_API_KEY"),
		os.Getenv("SHOPIFY_ADMIN_API_PASS_TOKEN"),
		os.Getenv("SHOPIFY_STORE_NAME"),
	)

	resp, err := http.Get(shopifyURL)
	if err != nil {
		http.Error(w, "Failed to fetch product list", http.StatusInternalServerError)
		fmt.Println("Error fetching products:", err.Error())
		return
	}
	defer resp.Body.Close()

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

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
	fmt.Println("Cache miss: Fetched data from Shopify API and cached it")
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]

	redisKey := fmt.Sprintf("products/%s", productID)

	cachedData, err := config.RedisClient.Get(ctx, redisKey).Result()
	if err == nil {
		// Data exists in Redis, return it
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cachedData))
		fmt.Printf("Cache hit: Returning data for product ID %s from Redis\n", productID)
		return
	}

	shopifyURL := fmt.Sprintf(
		"https://%s:%s@%s/admin/api/2022-10/products/%s.json",
		os.Getenv("SHOPIFY_API_KEY"),
		os.Getenv("SHOPIFY_ADMIN_API_PASS_TOKEN"),
		os.Getenv("SHOPIFY_STORE_NAME"),
		productID,
	)

	resp, err := http.Get(shopifyURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch product with ID %s", productID), http.StatusInternalServerError)
		fmt.Printf("Error fetching product with ID %s: %s\n", productID, err.Error())
		return
	}
	defer resp.Body.Close()

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

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
	fmt.Printf("Cache miss: Fetched data for product ID %s from Shopify API and cached it\n", productID)
}

func GetModel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]

	query := fmt.Sprintf(`
		{
			product(id: "gid://shopify/Product/%s") {
				id
				title
				media(first: 5) {
					edges {
						node {
							mediaContentType
							alt
							... on Model3d {
								id
								sources {
									url
									format
									mimeType
								}
							}
						}
					}
				}
			}
		}
	`, productID)

	apiURL := fmt.Sprintf("https://%s/admin/api/2024-10/graphql.json", os.Getenv("SHOPIFY_STORE_NAME"))
	accessToken := os.Getenv("SHOPIFY_ADMIN_API_PASS_TOKEN")

	if accessToken == "" {
		http.Error(w, "Shopify access token not set", http.StatusInternalServerError)
	}

	responseBody, err := executeGraphQLRequest(apiURL, accessToken, query)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error executing GraphQL request: %v", err), http.StatusInternalServerError)
		return
	}

	// Write the response back to the client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseBody)
}

func executeGraphQLRequest(url, accessToken, query string) ([]byte, error) {
	// Prepare request body
	requestBody, err := json.Marshal(map[string]string{"query": query})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	// Create the HTTP POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Shopify-Access-Token", accessToken)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GraphQL query failed with status: %d, response: %s", resp.StatusCode, body)
	}

	return body, nil
}
