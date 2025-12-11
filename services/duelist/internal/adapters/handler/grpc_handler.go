package handler

import (
	pb "api/proto" // Import generated proto
	"api/services/duelist/internal/core/domain"
	"api/services/duelist/internal/core/ports"
	"context"
)

type GrpcHandler struct {
	pb.UnimplementedDuelistServiceServer
	service ports.DuelistService
}

func NewGrpcHandler(service ports.DuelistService) *GrpcHandler {
	return &GrpcHandler{service: service}
}

func (h *GrpcHandler) CreateCowboy(ctx context.Context, req *pb.CreateCowboyRequest) (*pb.CowboyResponse, error) {
	domainCowboy := &domain.Cowboy{
		ID:       req.Id,
		Name:     req.Name,
		Health:   int(req.Health),
		Damage:   int(req.Damage),
		Speed:    int(req.Speed),
		Accuracy: req.Accuracy,
	}

	created, err := h.service.Create(domainCowboy)
	if err != nil {
		return nil, err
	}

	return h.toProto(created), nil
}

func (h *GrpcHandler) GetCowboy(ctx context.Context, req *pb.GetCowboyRequest) (*pb.CowboyResponse, error) {
	cowboy, err := h.service.Get(req.Id)
	if err != nil {
		return nil, err
	}
	return h.toProto(cowboy), nil
}

func (h *GrpcHandler) toProto(c *domain.Cowboy) *pb.CowboyResponse {
	return &pb.CowboyResponse{
		Id:       c.ID,
		Name:     c.Name,
		Health:   int32(c.Health),
		Damage:   int32(c.Damage),
		Speed:    int32(c.Speed),
		Accuracy: c.Accuracy,
	}
}
