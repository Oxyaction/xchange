package main

import (
	"context"
	"log"

	pb "github.com/Oxyaction/xchange/rpc"
	"github.com/golang/protobuf/ptypes"
)

type OrderServer struct {
	orderRepository *orderRepository
	pb.UnimplementedOrderServer
}

func (s *OrderServer) CreateSellOrder(ctx context.Context, req *pb.CreateSellOrderRequest) (*pb.SellOrderReply, error) {
	order := s.orderRepository.Create(context.Background(), &orderDTO{
		amount:   int(req.GetAmount()),
		price:    int(req.GetPrice()),
		assetID:  req.GetAssetId(),
		sellerID: req.GetSellerId(),
	})

	ts, err := ptypes.TimestampProto(order.createdAt)

	if err != nil {
		log.Fatal(err)
	}

	return &pb.SellOrderReply{
		Id:        order.id,
		Amount:    int32(order.amount),
		Price:     int32(order.price),
		AssetId:   order.assetID,
		SellerId:  order.sellerID,
		CreatedAt: ts,
	}, nil
}
