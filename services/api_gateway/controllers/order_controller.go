package controllers

import (
	"context"
	"encoding/json"
	"net/http"

	pb "github.com/Oxyaction/xchange/rpc"
)

type OrderController struct {
	Client pb.OrderClient
}

func (c *OrderController) CreateOrder(w http.ResponseWriter, r *http.Request) {
	reply, err := c.Client.CreateSellOrder(context.Background(), &pb.CreateSellOrderRequest{})
	if err != nil {
		DisplayAppError(w, err, "Create account error", 500)
		return
	}
	j, err := json.Marshal(reply)
	if err != nil {
		DisplayAppError(w, err, "JSON encoding error", 500)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	// Write the JSON data to the ResponseWriter
	w.Write(j)
}
