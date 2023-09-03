package twitch

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nicklaw5/helix"
)

type TwitchService struct {
	client *helix.Client
}

func NewTwitchService(clientID string, clientSecret string) *TwitchService {
	client, err := helix.NewClient(&helix.Options{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	})
	if err != nil {
		log.Fatal("Failed to create Twitch client:", err)
	}

	resp, err := client.RequestAppAccessToken([]string{"user:read:email"})

	if err != nil {
		log.Fatal("Failed to create Twitch client:", err)
	}

	client.SetAppAccessToken(resp.Data.AccessToken)

	return &TwitchService{
		client: client,
	}
}

func (s *TwitchService) GetTopStreams(count int) ([]helix.Stream, error) {
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	streamsResponse, err := s.client.GetStreams(&helix.StreamsParams{
		First: count,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get streams: %w", err)
	}

	return streamsResponse.Data.Streams, nil
}

func (s *TwitchService) GetStreamsByCategory(categoryIDs []string) ([]helix.Stream, error) {
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	streamsResponse, err := s.client.GetStreams(&helix.StreamsParams{
		GameIDs: categoryIDs,
	})

	if err != nil {
		log.Printf("Failed to get streams: %v\n", err)
	}

	return streamsResponse.Data.Streams, nil
}
