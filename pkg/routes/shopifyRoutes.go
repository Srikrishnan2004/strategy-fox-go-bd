package routes

import (
	"github.com/gorilla/mux"
	"strategy-fox-go-bd/pkg/controllers"
)

var ShopifyRoutes = func(router *mux.Router) {
	router.HandleFunc("/products", controllers.GetProducts).Methods("GET")
	router.HandleFunc("/products/{id}", controllers.GetProduct).Methods("GET")
	router.HandleFunc("/products/{id}/model", controllers.GetModel).Methods("GET")
}
