package client

import (
	"context"
	pb "api/proto"
	"api/services/arena/internal/core/ports"
	"api/services/arena/internal/core/domain/entity"
	"time"
)

type grpcClientAdapter struct {
	client pb.DuelistServiceClient
}

func NewGrpcClientAdapter(client pb.DuelistServiceClient) ports.CowboyProvider {
	return &grpcClientAdapter{client: client}
}

func (g *grpcClientAdapter) GetCowboy(id string) (*entity.Cowboy, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := g.client.GetCowboy(ctx, &pb.GetCowboyRequest{Id: id})
	if err != nil {
		return nil, err
	}

	return &entity.Cowboy{
		ID:       resp.Id,
		Name:     resp.Name,
		Health:   int(resp.Health),
		Damage:   int(resp.Damage),
		Speed:    int(resp.Speed),
		Accuracy: resp.Accuracy,
	}, nil
}