package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/jcorrea-videoamp/GoPostgreCrud/project/models"
	"github.com/jcorrea-videoamp/GoPostgreCrud/project/proto"
	"github.com/jcorrea-videoamp/GoPostgreCrud/project/repository"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type OrderService interface {
	proto.OrderServiceServer
}

type orderService struct {
	repo repository.Repository
	proto.UnimplementedOrderServiceServer
}

func NewOrderService(r repository.Repository) (*orderService, error) {
	if r == nil {
		return nil, fmt.Errorf("repository instance not available or valid, unable to return service")
	}
	return &orderService{
		repo: r,
	}, nil
}

func (os *orderService) ListOrders(ctx context.Context, void *emptypb.Empty) (*proto.ListOrderResponse, error) {
	orders, err := os.repo.ListOrders()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch orders from database: %w", err)
	}
	respOrders := make([]*proto.Order, len(orders))
	for i, order := range orders {
		respOrders[i] = repositoryOrderToProtoOrder(order)
	}
	return &proto.ListOrderResponse{
		Orders: respOrders,
	}, nil
}

func (os *orderService) CreateOrder(ctx context.Context, req *proto.CreateRequest) (*proto.AcknowlegeResponse, error) {
	order := protoOrderToRepositoryOrder(req.Order)
	resp, err := os.repo.CreateOrder(order)
	if err != nil {
		return nil, fmt.Errorf("failed to insert order in database: %w", err)
	}
	return &proto.AcknowlegeResponse{
		Response: resp,
	}, nil
}

func (os *orderService) GetOrder(ctx context.Context, req *proto.GetRequest) (*proto.OrderResponse, error) {
	order, err := os.repo.GetOrder(strconv.Itoa(int(req.GetId())))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order from database: %w", err)
	}
	return &proto.OrderResponse{Order: repositoryOrderToProtoOrder(order)}, nil
}

func (os *orderService) DeleteOrder(ctx context.Context, req *proto.GetRequest) (*proto.AcknowlegeResponse, error) {
	resp, err := os.repo.DeleteOrder(strconv.Itoa(int(req.GetId())))
	if err != nil {
		return nil, fmt.Errorf("failed to remove order from database: %w", err)
	}
	return &proto.AcknowlegeResponse{
		Response: resp,
	}, nil
}

func (os *orderService) UpdateOrder(ctx context.Context, req *proto.CreateRequest) (*proto.AcknowlegeResponse, error) {
	orderID := strconv.Itoa(int(req.Order.Id))
	updateFields, err := mapUpdateableFields(req.Order)
	if err != nil {
		return nil, fmt.Errorf("failed to update order in database: %w", err)
	}
	resp, err := os.repo.UpdateOrder(updateFields, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to update order in database: %w", err)
	}
	return &proto.AcknowlegeResponse{
		Response: resp,
	}, nil
}

func mapUpdateableFields(order *proto.Order) (map[string]any, error) {
	fields := make(map[string]any)

	if strings.TrimSpace(order.Status) != "" {
		fields["status"] = order.Status
	}

	if strings.TrimSpace(order.Customer) != "" {
		fields["customer"] = order.Customer
	}

	if order.Quantity != 0 {
		fields["quantity"] = int(order.Quantity)
	}

	if order.Price != 0 {
		fields["price"] = float64(order.Price)
	}

	if len(fields) == 0 {
		return nil, fmt.Errorf("nothing to be updated in the order")
	}
	return fields, nil
}

func repositoryOrderToProtoOrder(order *models.Order) *proto.Order {
	return &proto.Order{
		Id:        int32(order.ID),
		Status:    order.Status,
		Customer:  order.Customer,
		Quantity:  int32(order.Quantity),
		Price:     float32(order.Price),
		CreatedAt: timestamppb.New(order.CreatedAt),
		UpdatedAt: timestamppb.New(order.UpdatedAt),
	}
}

func protoOrderToRepositoryOrder(order *proto.Order) *models.Order {
	return &models.Order{
		ID:       int(order.Id),
		Status:   order.Status,
		Customer: order.Customer,
		Quantity: int(order.Quantity),
		Price:    float64(order.Price),
		OrderHistory: models.OrderHistory{
			CreatedAt: order.CreatedAt.AsTime(),
			UpdatedAt: order.UpdatedAt.AsTime(),
		},
	}
}
