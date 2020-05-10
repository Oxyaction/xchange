package routers

import (
	pb "github.com/Oxyaction/xchange/rpc"
	"github.com/Oxyaction/xchange/services/api_gateway/controllers"

	"github.com/gorilla/mux"
)

func SetOrdertRoutes(c pb.OrderClient, router *mux.Router) {
	orderController := controllers.OrderController{c}
	router.HandleFunc("/order", orderController.CreateOrder).Methods("POST")
}
