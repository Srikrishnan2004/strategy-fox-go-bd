package routes

import (
	"github.com/gorilla/mux"
	"strategy-fox-go-bd/pkg/controllers"
)

var ShopifyRoutes = func(router *mux.Router) {
	//router.HandleFunc("/v1/products", controllers.GetProducts).Methods("GET")
	router.HandleFunc("/v2/products", controllers.GetProductsGQ).Methods("GET")
	//router.HandleFunc("/products/{id}", controllers.GetProduct).Methods("GET")
	//router.HandleFunc("/products/{id}/model", controllers.GetModel).Methods("GET")
	router.HandleFunc("/v2/products/by-name/{name}", controllers.GetProductByNameGQ).Methods("GET")
	router.HandleFunc("/v2/products/by-id/{id}", controllers.GetProductByIdGQ).Methods("GET")
	router.HandleFunc("/product/metafield", controllers.UpdateMetafieldById).Methods("POST")

}
