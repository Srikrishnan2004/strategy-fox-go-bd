package routes

import (
	"github.com/gorilla/mux"
	"strategy-fox-go-bd/pkg/controllers"
)

var ChatBotRoutes = func(router *mux.Router) {
	router.HandleFunc("/chat", controllers.HandleChat).Methods("POST")
}
