package service

import (
	"context"
	"fmt"

	pb "github.com/767829413/advanced-go/api/wxopen"
)

type WxopenService struct {
	pb.UnimplementedWxopenServer
}

func NewWxopenService() *WxopenService {
	return &WxopenService{}
}

func (s *WxopenService) GetAccessToken(
	ctx context.Context,
	req *pb.GetAccessTokenRequest,
) (*pb.GetAccessTokenResponse, error) {
	fmt.Println("GetAccessToken")
	return &pb.GetAccessTokenResponse{}, nil
}

func (s *WxopenService) LoginQrCodeCreate(
	ctx context.Context,
	req *pb.LoginQrCodeCreateRequest,
) (*pb.LoginQrCodeCreateResponse, error) {
	fmt.Println("LoginQrCodeCreate")
	return &pb.LoginQrCodeCreateResponse{}, nil
}

func (s *WxopenService) GetWxUserInfoByCode(
	ctx context.Context,
	req *pb.GetWxUserInfoByCodeRequest,
) (*pb.GetWxUserInfoByCodeResponse, error) {
	fmt.Println("GetWxUserInfoByCode")
	return &pb.GetWxUserInfoByCodeResponse{}, nil
}

func (s *WxopenService) LoginUrlCreate(
	ctx context.Context,
	req *pb.LoginUrlCreateRequest,
) (*pb.LoginUrlCreateResponse, error) {
	fmt.Println("LoginUrlCreate")
	return &pb.LoginUrlCreateResponse{}, nil
}
