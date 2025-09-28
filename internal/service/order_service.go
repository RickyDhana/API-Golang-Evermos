package service

import (
	"errors"

	"project_satu/internal/config"
	"project_satu/internal/domain"
)

type OrderItemRequest struct {
	ProdukID uint `json:"produk_id"`
	Qty      int  `json:"qty"`
}

type CreateOrderRequest struct {
	Items []OrderItemRequest `json:"items"`
}

func CreateOrder(userID uint, req *CreateOrderRequest) (*domain.Order, error) {
	if len(req.Items) == 0 {
		return nil, errors.New("items cannot be empty")
	}

	order := &domain.Order{
		UserID: userID,
		Status: "pending",
	}


	if err := config.DB.Create(order).Error; err != nil {
		return nil, err
	}

	var total int
	for _, item := range req.Items {
		var produk domain.Produk
		if err := config.DB.First(&produk, item.ProdukID).Error; err != nil {
			return nil, err
		}

		subtotal := int(produk.Harga) * item.Qty
		total += subtotal

		orderItem := &domain.OrderItem{
			OrderID:  order.ID,
			ProdukID: produk.ID,
			Qty:      item.Qty,
			Subtotal: subtotal,
		}
		if err := config.DB.Create(orderItem).Error; err != nil {
			return nil, err
		}

		productLog := &domain.ProductLog{
			ProdukID: produk.ID,
			Nama:     produk.NamaProduk,
			Harga:    int(produk.Harga),
			Qty:      item.Qty,
			OrderID:  order.ID,
		}
		if err := config.DB.Create(productLog).Error; err != nil {
			return nil, err
		}
	}

	order.Total = total
	if err := config.DB.Save(order).Error; err != nil {
		return nil, err
	}

	return order, nil
}

func ListOrders(userID uint) ([]domain.Order, error) {
	var orders []domain.Order
	if err := config.DB.Where("user_id = ?", userID).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func GetOrderByID(userID uint, orderID uint) (*domain.Order, error) {
	var order domain.Order
	if err := config.DB.Where("id = ? AND user_id = ?", orderID, userID).First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}
