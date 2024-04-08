package service

import (
	"context"
	"fmt"

	pb "github.com/767829413/advanced-go/api/orgSet"
)

type OrgConfService struct {
	pb.UnimplementedOrgConfServer
}

func NewOrgConfService() *OrgConfService {
	return &OrgConfService{}
}

func (s *OrgConfService) GetOrgConf(
	ctx context.Context,
	req *pb.GetOrgConfRequest,
) (*pb.GetOrgConfResponse, error) {
	fmt.Println("GetOrgConf")
	return &pb.GetOrgConfResponse{}, nil
}

func (s *OrgConfService) SetOrgConf(
	ctx context.Context,
	req *pb.SetOrgConfRequest,
) (*pb.SetOrgConfResponse, error) {
	fmt.Println("SetOrgConf")
	return &pb.SetOrgConfResponse{}, nil
}
