package services

import pb "github.com/Oxyaction/xchange/rpc"

type Services struct {
	OrderService *OrderService
}

func InitServices(accountClient *pb.AccountClient, orderClient *pb.OrderClient) *Services {
	return &Services{
		OrderService: &OrderService{accountClient, orderClient},
	}
}
