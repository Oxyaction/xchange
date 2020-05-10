package main

import (
	"context"

	pb "github.com/Oxyaction/xchange/rpc"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AccountServer struct {
	accountRepository *accountRepository
	assetRepository   *assetRepository
	pb.UnimplementedAccountServer
}

func (*AccountServer) ChangeBalance(ctx context.Context, req *pb.ChangeBalanceRequest) (*pb.AccountReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChangeBalance not implemented")
}

func (server *AccountServer) Create(ctx context.Context, req *pb.CreateRequest) (*pb.AccountReply, error) {
	account := server.accountRepository.Create(ctx)
	return &pb.AccountReply{
		Id:      account.Id,
		Balance: int32(account.Balance),
	}, nil
}

func (server *AccountServer) CreateAsset(ctx context.Context, req *pb.CreateAssetRequest) (*pb.Asset, error) {
	asset := server.assetRepository.Create(ctx)
	return &pb.Asset{
		Id:   asset.Id,
		Name: asset.Name,
	}, nil
}

func (server *AccountServer) ChangeAssetBalance(ctx context.Context, req *pb.AssetBalance) (*pb.AssetBalance, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChangeAssetBalance not implemented")
}

func (server *AccountServer) GetBalance(ctx context.Context, req *pb.GetBalanceRequest) (*pb.AccountReply, error) {
	account, err := server.accountRepository.Get(ctx, req.Id)

	if err != nil {
		if err == notFound {
			return nil, status.Errorf(codes.NotFound, "Entity not found")
		}
		return nil, err
	}

	return &pb.AccountReply{
		Id:      account.Id,
		Balance: int32(account.Balance),
	}, nil
}
