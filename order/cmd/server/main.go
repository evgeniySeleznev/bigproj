package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	inventoryv1 "github.com/evgeniyseleznev/bigproj/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/evgeniyseleznev/bigproj/shared/pkg/proto/payment/v1"
)

const (
	httpPort = "8080"
)

// OrderStatus represents order status
type OrderStatus string

const (
	OrderStatusPendingPayment OrderStatus = "PENDING_PAYMENT"
	OrderStatusPaid           OrderStatus = "PAID"
	OrderStatusCancelled      OrderStatus = "CANCELLED"
)

// Order represents an order
type Order struct {
	OrderUUID       string      `json:"order_uuid"`
	UserUUID        string      `json:"user_uuid"`
	PartUUIDs       []string    `json:"part_uuids"`
	TotalPrice      float64     `json:"total_price"`
	TransactionUUID *string     `json:"transaction_uuid,omitempty"`
	PaymentMethod   *string     `json:"payment_method,omitempty"`
	Status          OrderStatus `json:"status"`
}

// OrderStorage is thread-safe order storage
type OrderStorage struct {
	mu     sync.RWMutex
	orders map[string]*Order
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[string]*Order),
	}
}

func (s *OrderStorage) GetOrder(uuid string) (*Order, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.orders[uuid]
	if !ok {
		return nil, errors.New("not found")
	}
	return order, nil
}

func (s *OrderStorage) CreateOrder(order *Order) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.orders[order.OrderUUID] = order
}

func (s *OrderStorage) UpdateOrder(uuid string, updateFunc func(*Order)) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	order, ok := s.orders[uuid]
	if !ok {
		return errors.New("not found")
	}
	updateFunc(order)
	return nil
}

// OrderHandler handles order requests
type OrderHandler struct {
	storage         *OrderStorage
	inventoryClient inventoryv1.InventoryServiceClient
	paymentClient   paymentv1.PaymentServiceClient
	inventoryConn   *grpc.ClientConn
	paymentConn     *grpc.ClientConn
}

func NewOrderHandler() (*OrderHandler, error) {
	// Connect to Inventory Service
	inventoryConn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to inventory service: %w", err)
	}

	// Connect to Payment Service
	paymentConn, err := grpc.NewClient(
		"localhost:50052",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		if closeErr := inventoryConn.Close(); closeErr != nil {
			log.Printf("error closing inventory connection: %v", closeErr)
		}
		return nil, fmt.Errorf("failed to connect to payment service: %w", err)
	}

	return &OrderHandler{
		storage:         NewOrderStorage(),
		inventoryClient: inventoryv1.NewInventoryServiceClient(inventoryConn),
		paymentClient:   paymentv1.NewPaymentServiceClient(paymentConn),
		inventoryConn:   inventoryConn,
		paymentConn:     paymentConn,
	}, nil
}

func (h *OrderHandler) Close() {
	if h.inventoryConn != nil {
		if err := h.inventoryConn.Close(); err != nil {
			log.Printf("error closing inventory connection: %v", err)
		}
	}
	if h.paymentConn != nil {
		if err := h.paymentConn.Close(); err != nil {
			log.Printf("error closing payment connection: %v", err)
		}
	}
}

// PostOrders creates a new order
func (h *OrderHandler) PostOrders(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserUUID  string   `json:"user_uuid"`
		PartUUIDs []string `json:"part_uuids"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if encodeErr := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}); encodeErr != nil {
			log.Printf("error encoding error response: %v", encodeErr)
		}
		return
	}

	// Валидация входных данных
	if req.UserUUID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if encodeErr := json.NewEncoder(w).Encode(map[string]string{"error": "user_uuid is required"}); encodeErr != nil {
			log.Printf("error encoding error response: %v", encodeErr)
		}
		return
	}

	if len(req.PartUUIDs) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if encodeErr := json.NewEncoder(w).Encode(map[string]string{"error": "part_uuids is required and cannot be empty"}); encodeErr != nil {
			log.Printf("error encoding error response: %v", encodeErr)
		}
		return
	}

	// Получаем детали из Inventory Service
	resp, err := h.inventoryClient.ListParts(r.Context(), &inventoryv1.ListPartsRequest{
		Filter: &inventoryv1.PartsFilter{
			Uuids: req.PartUUIDs,
		},
	})
	if err != nil {
		log.Printf("error calling InventoryService: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		if encodeErr := json.NewEncoder(w).Encode(map[string]string{"error": "Bad Gateway"}); encodeErr != nil {
			log.Printf("error encoding error response: %v", encodeErr)
		}
		return
	}

	// Проверяем, что все детали найдены
	partsFound := make(map[string]bool)
	for _, part := range resp.GetParts() {
		partsFound[part.GetUuid()] = true
	}

	for _, uuid := range req.PartUUIDs {
		if !partsFound[uuid] {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			if encodeErr := json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("part %s not found", uuid)}); encodeErr != nil {
				log.Printf("error encoding error response: %v", encodeErr)
			}
			return
		}
	}

	// Подсчитываем total_price
	var totalPrice float64
	for _, part := range resp.GetParts() {
		totalPrice += part.GetPrice()
	}

	// Создаем заказ
	order := &Order{
		OrderUUID:  uuid.New().String(),
		UserUUID:   req.UserUUID,
		PartUUIDs:  req.PartUUIDs,
		TotalPrice: totalPrice,
		Status:     OrderStatusPendingPayment,
	}

	h.storage.CreateOrder(order)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"uuid":        order.OrderUUID,
		"total_price": order.TotalPrice,
	}); err != nil {
		log.Printf("error encoding response: %v", err)
	}
}

// GetOrders gets order by UUID
func (h *OrderHandler) GetOrders(w http.ResponseWriter, _ *http.Request, orderUUID string) {
	order, err := h.storage.GetOrder(orderUUID)
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(order); err != nil {
		log.Printf("error encoding response: %v", err)
	}
}

// PostOrdersPay pays the order
func (h *OrderHandler) PostOrdersPay(w http.ResponseWriter, r *http.Request, orderUUID string) {
	var req struct {
		PaymentMethod string `json:"payment_method"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	order, err := h.storage.GetOrder(orderUUID)
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	// Вызываем Payment Service
	paymentMethod := paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	switch req.PaymentMethod {
	case "CARD", "PAYMENT_METHOD_CARD":
		paymentMethod = paymentv1.PaymentMethod_PAYMENT_METHOD_CARD
	case "SBP", "PAYMENT_METHOD_SBP":
		paymentMethod = paymentv1.PaymentMethod_PAYMENT_METHOD_SBP
	case "CREDIT_CARD", "PAYMENT_METHOD_CREDIT_CARD":
		paymentMethod = paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case "INVESTOR_MONEY", "PAYMENT_METHOD_INVESTOR_MONEY":
		paymentMethod = paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	}

	resp, err := h.paymentClient.PayOrder(r.Context(), &paymentv1.PayOrderRequest{
		OrderUuid:     orderUUID,
		UserUuid:      order.UserUUID,
		PaymentMethod: paymentMethod,
	})
	if err != nil {
		log.Printf("error calling PaymentService: %v", err)
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return
	}

	// Обновляем заказ
	transactionUUID := resp.GetTransactionUuid()
	if err := h.storage.UpdateOrder(orderUUID, func(o *Order) {
		o.Status = OrderStatusPaid
		o.TransactionUUID = &transactionUUID
		method := req.PaymentMethod
		o.PaymentMethod = &method
	}); err != nil {
		log.Printf("error updating order: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{
		"transaction_uuid": transactionUUID,
	}); err != nil {
		log.Printf("error encoding response: %v", err)
	}
}

// PostOrdersCancel cancels the order
func (h *OrderHandler) PostOrdersCancel(w http.ResponseWriter, _ *http.Request, orderUUID string) {
	order, err := h.storage.GetOrder(orderUUID)
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	if order.Status == OrderStatusPaid {
		http.Error(w, "Order already paid and cannot be cancelled", http.StatusConflict)
		return
	}

	err = h.storage.UpdateOrder(orderUUID, func(o *Order) {
		o.Status = OrderStatusCancelled
	})
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	handler, err := NewOrderHandler()
	if err != nil {
		log.Fatalf("failed to create handler: %v", err)
	}
	defer handler.Close()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	r.Post("/api/v1/orders", handler.PostOrders)
	r.Get("/api/v1/orders/{order_uuid}", func(w http.ResponseWriter, r *http.Request) {
		handler.GetOrders(w, r, chi.URLParam(r, "order_uuid"))
	})
	r.Post("/api/v1/orders/{order_uuid}/pay", func(w http.ResponseWriter, r *http.Request) {
		handler.PostOrdersPay(w, r, chi.URLParam(r, "order_uuid"))
	})
	r.Post("/api/v1/orders/{order_uuid}/cancel", func(w http.ResponseWriter, r *http.Request) {
		handler.PostOrdersCancel(w, r, chi.URLParam(r, "order_uuid"))
	})

	server := &http.Server{
		Addr:              net.JoinHostPort("0.0.0.0", httpPort),
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("📦 Order Service HTTP server listening on port %s", httpPort)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down Order Service...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	if err := server.Shutdown(ctx); err != nil {
		cancel()
		log.Printf("server shutdown error: %v", err)
		return
	}
	cancel()
	log.Println("✅ Server stopped")
}
