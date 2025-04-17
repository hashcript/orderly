package orders

import (
	"log"

	"com.orderly/items"
	"github.com/gin-gonic/gin"
)

// OrderSerializer serializes a single OrderModel instance for API responses
type OrderSerializer struct {
	c          *gin.Context
	OrderModel OrderModel
}

// NewOrderSerializer creates a new instance of OrderSerializer
func NewOrderSerializer(c *gin.Context, order OrderModel) *OrderSerializer {
	return &OrderSerializer{
		c:          c,
		OrderModel: order,
	}
}

// Response returns a serialized representation of the OrderModel
func (s *OrderSerializer) Response() (map[string]interface{}, error) {
	// Serialize the associated user
	userData := map[string]interface{}{
		"id":    s.OrderModel.User.ID,
		"name":  s.OrderModel.User.Name,
		"email": s.OrderModel.User.Email,
		"phone": s.OrderModel.User.Phone,
	}

	// Serialize the associated item
	itemSerializer := items.NewItemSerializer(s.c, s.OrderModel.Item)
	itemResponse := itemSerializer.Response()

	return map[string]interface{}{
		"id":         s.OrderModel.ID,
		"user":       userData,
		"item":       itemResponse,
		"quantity":   s.OrderModel.Quantity,
		"totalPrice": s.OrderModel.TotalPrice,
		"status":     s.OrderModel.Status,
		"createdAt":  s.OrderModel.CreatedAt,
		"updatedAt":  s.OrderModel.UpdatedAt,
	}, nil
}

// OrderSerializerList serializes a list of OrderModel instances for API responses
type OrderSerializerList struct {
	c      *gin.Context
	Orders []OrderModel
}

// NewOrderSerializerList creates a new instance of OrderSerializerList
func NewOrderSerializerList(c *gin.Context, orders []OrderModel) *OrderSerializerList {
	return &OrderSerializerList{c: c, Orders: orders}
}

// Response returns a serialized representation of a list of OrderModels
func (s *OrderSerializerList) Response() ([]map[string]interface{}, error) {
	var responseList []map[string]interface{}
	for _, order := range s.Orders {
		orderSerializer := NewOrderSerializer(s.c, order)
		orderResponse, err := orderSerializer.Response()
		if err != nil {
			log.Printf("Error serializing order: %v", err)
			return nil, err // You could choose to handle the error differently, e.g., skip this order and continue
		}
		responseList = append(responseList, orderResponse)
	}
	return responseList, nil
}
