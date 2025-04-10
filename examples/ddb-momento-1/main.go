package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/momento"
)

type CachedLayout struct {
	SessionId  string `json:"sessionId"`
	Layout     string `json:"layout"`
	LastUpdate int64  `json:"lastUpdate"`
	LayoutHash string `json:"layoutHash"`
}

// Service Base service that simulates non-caching logic
type Service struct{}

func (s *Service) GetCachedSortedLayout(ctx context.Context, pageReference, userID, deviceID string) (*CachedLayout, error) {
	fmt.Println("[Service] Fetching layout from base implementation")
	return &CachedLayout{
		SessionId:  userID,
		Layout:     fmt.Sprintf("layout-%s", pageReference),
		LastUpdate: time.Now().Unix(),
		LayoutHash: "hash-123-base",
	}, nil
}

// MomentoService Momento-enhanced service that wraps base Service
type MomentoService struct {
	Service
	MomentoClient momento.CacheClient
	CacheName     string
	TTL           time.Duration
}

func (m *MomentoService) GetCachedSortedLayout(ctx context.Context, pageReference, userID, deviceID string) (*CachedLayout, error) {
	fmt.Println("[MomentoService] Fetching layout with caching logic")
	return &CachedLayout{
		SessionId:  userID,
		Layout:     fmt.Sprintf("layout-%s", pageReference),
		LastUpdate: time.Now().Unix(),
		LayoutHash: "hash-123-momento",
	}, nil
}

func main() {
	ctx := context.Background()

	credentialProvider, err := auth.NewEnvMomentoTokenProvider("MOMENTO_API_KEY")
	if err != nil {
		log.Fatalf("failed to load momento credentials: %v", err)
	}

	client, err := momento.NewCacheClient(
		config.LaptopLatest(),
		credentialProvider,
		5*time.Second,
	)
	if err != nil {
		log.Fatalf("failed to create momento client: %v", err)
	}

	cacheName := "layout-cache"
	momentoService := &MomentoService{
		Service:       Service{},
		MomentoClient: client,
		CacheName:     cacheName,
		TTL:           10 * time.Second,
	}

	// First call - from momento service
	layout1, err := momentoService.GetCachedSortedLayout(ctx, "home-momento", "user-momento", "deviceABC-momento")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Returned layout (1st call): %+v\n\n", layout1)

	baseService := &Service{}
	// Second call - from base service
	layout3, err := baseService.GetCachedSortedLayout(ctx, "home-base", "user123-base", "deviceABC-base")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Returned layout (base service): %+v\n", layout3)
}
