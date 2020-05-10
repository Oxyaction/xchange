package routers

import (
	pb "github.com/Oxyaction/xchange/rpc"
	"github.com/Oxyaction/xchange/services/api_gateway/controllers"

	"github.com/gorilla/mux"
)

func SetAccountRoutes(c pb.AccountClient, router *mux.Router) {
	accountController := controllers.AccountController{c}
	router.HandleFunc("/account", accountController.CreateAccount).Methods("POST")
	router.HandleFunc("/account/{id}", accountController.GetAccount).Methods("GET")
}
