package services

import pb "github.com/Oxyaction/xchange/rpc"

type OrderService struct {
	AccountClient *pb.AccountClient
	OrderClient   *pb.OrderClient
}

type SellOrder struct {
	Amount int
	Price  int
}

func (os *OrderService) CreateSellOrder() {
	
}
