package orders

import (
	"fmt"

	"com.orderly/database"
	"com.orderly/items"
	"com.orderly/users"
	"github.com/gin-gonic/gin"
)

type OrderValidator struct {
	Order struct {
		UserID     uint    `form:"user_id" json:"user_id" binding:"required"`
		ItemID     uint    `form:"item_id" json:"item_id" binding:"required"`
		Quantity   uint    `form:"quantity" json:"quantity" binding:"required,min=1"`
		TotalPrice float32 `form:"total_price" json:"total_price"`
		Status     string  `form:"status" json:"status" binding:"omitempty,oneof=pending completed canceled"`
	} `json:"order"`
	OrderModel OrderModel `json:"-"`
}

// Bind maps input data to the `OrderModel`
// Bind maps input data to the `OrderModel`
// It also validates the existence of the user and item in the database.
func (v *OrderValidator) Bind(c *gin.Context) error {
	// Bind the request body to the OrderValidator struct
	err := database.Bind(c, v)
	if err != nil {
		return err
	}

	// Check if the User exists in the database
	var user users.UserModel
	if err := database.GetConnection().First(&user, v.Order.UserID).Error; err != nil {
		return fmt.Errorf("user with ID %d not found", v.Order.UserID)
	}

	// Check if the Item exists in the database
	var item items.ItemModel
	if err := database.GetConnection().First(&item, v.Order.ItemID).Error; err != nil {
		return fmt.Errorf("item with ID %d not found", v.Order.ItemID)
	}

	// Populate the OrderModel fields
	v.OrderModel.UserID = v.Order.UserID
	v.OrderModel.ItemID = v.Order.ItemID
	v.OrderModel.Quantity = v.Order.Quantity
	v.OrderModel.TotalPrice = v.Order.TotalPrice
	// Set default status if empty
	if v.Order.Status == "" {
		v.Order.Status = "pending"
	}
	v.OrderModel.Status = v.Order.Status

	return nil
}

// NewOrderValidator creates a new instance of the validator
func NewOrderValidator() OrderValidator {
	return OrderValidator{}
}

// NewOrderValidatorFillWith creates a validator prefilled with an existing order's data
func NewOrderValidatorFillWith(order OrderModel) OrderValidator {
	return OrderValidator{
		Order: struct {
			UserID     uint    `form:"user_id" json:"user_id" binding:"required"`
			ItemID     uint    `form:"item_id" json:"item_id" binding:"required"`
			Quantity   uint    `form:"quantity" json:"quantity" binding:"required,min=1"`
			TotalPrice float32 `form:"total_price" json:"total_price"`
			Status     string  `form:"status" json:"status" binding:"omitempty,oneof=pending completed canceled"`
		}{
			UserID:     order.UserID,
			ItemID:     order.ItemID,
			Quantity:   order.Quantity,
			TotalPrice: order.TotalPrice,
			Status:     order.Status,
		},
	}
}
