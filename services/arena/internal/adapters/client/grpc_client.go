package client

import (
	"context"
	pb "api/proto"
	"api/services/arena/internal/core/domain"
	"api/services/arena/internal/core/ports"
	"time"
)

type grpcClientAdapter struct {
	client pb.DuelistServiceClient
}

func NewGrpcClientAdapter(client pb.DuelistServiceClient) ports.CowboyProvider {
	return &grpcClientAdapter{client: client}
}

func (g *grpcClientAdapter) GetCowboy(id string) (*domain.Cowboy, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := g.client.GetCowboy(ctx, &pb.GetCowboyRequest{Id: id})
	if err != nil {
		return nil, err
	}

	return &domain.Cowboy{
		ID:       resp.Id,
		Name:     resp.Name,
		Health:   int(resp.Health),
		Damage:   int(resp.Damage),
		Speed:    int(resp.Speed),
		Accuracy: resp.Accuracy,
	}, nil
}