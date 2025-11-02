package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	paymentv1 "github.com/evgeniyseleznev/bigproj/shared/pkg/proto/payment/v1"
)

const grpcPort = 50052

type paymentService struct {
	paymentv1.UnimplementedPaymentServiceServer
}

func (s *paymentService) PayOrder(ctx context.Context, req *paymentv1.PayOrderRequest) (*paymentv1.PayOrderResponse, error) {
	_ = ctx // context используется для соответствия интерфейсу gRPC
	_ = req // параметры могут использоваться в будущем (order_uuid, user_uuid, payment_method)

	// Генерируем UUID транзакции
	transactionUUID := uuid.New().String()

	// Логируем в консоль
	log.Printf("Оплата прошла успешно, transaction_uuid: %s", transactionUUID)

	// Возвращаем transaction_uuid
	return &paymentv1.PayOrderResponse{
		TransactionUuid: transactionUUID,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Создаем gRPC сервер
	s := grpc.NewServer()

	// Создаем сервис
	service := &paymentService{}
	paymentv1.RegisterPaymentServiceServer(s, service)

	// Включаем рефлексию для отладки
	reflection.Register(s)

	go func() {
		log.Printf("💰 Payment Service gRPC server listening on port %d", grpcPort)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down Payment Service...")
	s.GracefulStop()
	log.Println("✅ Server stopped")
}
