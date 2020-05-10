package routers

import (
	pb "github.com/Oxyaction/xchange/rpc"

	"github.com/gorilla/mux"
)

func InitRoutes(c pb.AccountClient) *mux.Router {
	router := mux.NewRouter().StrictSlash(false)
	SetAccountRoutes(c, router)
	return router
}
