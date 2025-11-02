package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	inventoryv1 "github.com/evgeniyseleznev/bigproj/shared/pkg/proto/inventory/v1"
)

const grpcPort = 50051

type inventoryService struct {
	inventoryv1.UnimplementedInventoryServiceServer

	mu    sync.RWMutex
	parts map[string]*inventoryv1.Part
}

func (s *inventoryService) GetPart(ctx context.Context, req *inventoryv1.GetPartRequest) (*inventoryv1.GetPartResponse, error) {
	_ = ctx // context используется для соответствия интерфейсу gRPC
	s.mu.RLock()
	defer s.mu.RUnlock()

	part, ok := s.parts[req.GetUuid()]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "part with UUID %s not found", req.GetUuid())
	}

	return &inventoryv1.GetPartResponse{
		Part: part,
	}, nil
}

// matchesUUID проверяет, соответствует ли деталь фильтру по UUID
func matchesUUID(part *inventoryv1.Part, uuids []string) bool {
	if len(uuids) == 0 {
		return true
	}
	for _, uuid := range uuids {
		if part.GetUuid() == uuid {
			return true
		}
	}
	return false
}

// matchesName проверяет, соответствует ли деталь фильтру по имени
func matchesName(part *inventoryv1.Part, names []string) bool {
	if len(names) == 0 {
		return true
	}
	for _, name := range names {
		if part.GetName() == name {
			return true
		}
	}
	return false
}

// matchesCategory проверяет, соответствует ли деталь фильтру по категории
func matchesCategory(part *inventoryv1.Part, categories []inventoryv1.Category) bool {
	if len(categories) == 0 {
		return true
	}
	for _, cat := range categories {
		if part.GetCategory() == cat {
			return true
		}
	}
	return false
}

// matchesManufacturerCountry проверяет, соответствует ли деталь фильтру по стране производителя
func matchesManufacturerCountry(part *inventoryv1.Part, countries []string) bool {
	if len(countries) == 0 {
		return true
	}
	if part.GetManufacturer() == nil {
		return false
	}
	for _, country := range countries {
		if part.GetManufacturer().GetCountry() == country {
			return true
		}
	}
	return false
}

// matchesTags проверяет, соответствует ли деталь фильтру по тегам
func matchesTags(part *inventoryv1.Part, tags []string) bool {
	if len(tags) == 0 {
		return true
	}
	partTags := part.GetTags()
	for _, filterTag := range tags {
		for _, partTag := range partTags {
			if partTag == filterTag {
				return true
			}
		}
	}
	return false
}

func (s *inventoryService) ListParts(ctx context.Context, req *inventoryv1.ListPartsRequest) (*inventoryv1.ListPartsResponse, error) {
	_ = ctx // context используется для соответствия интерфейсу gRPC
	s.mu.RLock()
	defer s.mu.RUnlock()

	filter := req.GetFilter()

	// Если фильтр пустой или все поля пустые - возвращаем все детали
	if filter == nil || (len(filter.GetUuids()) == 0 && len(filter.GetNames()) == 0 &&
		len(filter.GetCategories()) == 0 && len(filter.GetManufacturerCountries()) == 0 &&
		len(filter.GetTags()) == 0) {
		parts := make([]*inventoryv1.Part, 0, len(s.parts))
		for _, part := range s.parts {
			parts = append(parts, part)
		}
		return &inventoryv1.ListPartsResponse{Parts: parts}, nil
	}

	var result []*inventoryv1.Part
	for _, part := range s.parts {
		if matchesUUID(part, filter.GetUuids()) &&
			matchesName(part, filter.GetNames()) &&
			matchesCategory(part, filter.GetCategories()) &&
			matchesManufacturerCountry(part, filter.GetManufacturerCountries()) &&
			matchesTags(part, filter.GetTags()) {
			result = append(result, part)
		}
	}

	return &inventoryv1.ListPartsResponse{Parts: result}, nil
}

func (s *inventoryService) initParts() {
	now := timestamppb.Now()

	s.parts = map[string]*inventoryv1.Part{
		"part-uuid-1": {
			Uuid:          "part-uuid-1",
			Name:          "Ion Engine Model X1",
			Description:   "High-efficiency ion engine for deep space missions",
			Price:         45000.00,
			StockQuantity: 5,
			Category:      inventoryv1.Category_CATEGORY_ENGINE,
			Dimensions: &inventoryv1.Dimensions{
				Length: 120.5,
				Width:  80.0,
				Height: 120.5,
				Weight: 85.0,
			},
			Manufacturer: &inventoryv1.Manufacturer{
				Name:    "SpaceTech Industries",
				Country: "USA",
				Website: "www.spacetech.com",
			},
			Tags:      []string{"engine", "ion", "electric"},
			Metadata:  make(map[string]*inventoryv1.Value),
			CreatedAt: now,
			UpdatedAt: now,
		},
		"part-uuid-2": {
			Uuid:          "part-uuid-2",
			Name:          "Liquid Hydrogen Tank 500L",
			Description:   "Storage tank for liquid hydrogen fuel",
			Price:         23000.00,
			StockQuantity: 12,
			Category:      inventoryv1.Category_CATEGORY_FUEL,
			Dimensions: &inventoryv1.Dimensions{
				Length: 200.0,
				Width:  150.0,
				Height: 150.0,
				Weight: 120.0,
			},
			Manufacturer: &inventoryv1.Manufacturer{
				Name:    "Hydrogen Systems Ltd",
				Country: "Germany",
				Website: "www.hydrogensystems.de",
			},
			Tags:      []string{"fuel", "hydrogen", "tank"},
			Metadata:  make(map[string]*inventoryv1.Value),
			CreatedAt: now,
			UpdatedAt: now,
		},
		"part-uuid-3": {
			Uuid:          "part-uuid-3",
			Name:          "Observation Window 50cm",
			Description:   "Reinforced observation porthole for crew quarters",
			Price:         8500.00,
			StockQuantity: 20,
			Category:      inventoryv1.Category_CATEGORY_PORTHOLE,
			Dimensions: &inventoryv1.Dimensions{
				Length: 50.0,
				Width:  50.0,
				Height: 10.0,
				Weight: 25.0,
			},
			Manufacturer: &inventoryv1.Manufacturer{
				Name:    "ClearView Space Windows",
				Country: "Russia",
				Website: "www.clearview.ru",
			},
			Tags:      []string{"window", "observation", "crew"},
			Metadata:  make(map[string]*inventoryv1.Value),
			CreatedAt: now,
			UpdatedAt: now,
		},
		"part-uuid-4": {
			Uuid:          "part-uuid-4",
			Name:          "Solar Wing Panel 4m",
			Description:   "Solar panel wing for extended missions",
			Price:         15000.00,
			StockQuantity: 8,
			Category:      inventoryv1.Category_CATEGORY_WING,
			Dimensions: &inventoryv1.Dimensions{
				Length: 400.0,
				Width:  200.0,
				Height: 5.0,
				Weight: 90.0,
			},
			Manufacturer: &inventoryv1.Manufacturer{
				Name:    "SolarWorks GmbH",
				Country: "Germany",
				Website: "www.solarworks.de",
			},
			Tags:      []string{"solar", "wing", "power"},
			Metadata:  make(map[string]*inventoryv1.Value),
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Создаем gRPC сервер
	s := grpc.NewServer()

	// Создаем сервис
	service := &inventoryService{}
	service.initParts()

	inventoryv1.RegisterInventoryServiceServer(s, service)

	// Включаем рефлексию для отладки
	reflection.Register(s)

	go func() {
		log.Printf("🚀 Inventory Service gRPC server listening on port %d", grpcPort)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down Inventory Service...")
	s.GracefulStop()
	log.Println("✅ Server stopped")
}
