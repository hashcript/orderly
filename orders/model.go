package orders

import (
	"fmt"

	"com.orderly/database"
	"com.orderly/items"
	"com.orderly/users"

	"gorm.io/gorm"
)

type OrderModel struct {
	gorm.Model
	ID         uint            `gorm:"primary_key"`
	UserID     uint            `gorm:"not null" json:"user_id"`                  // Foreign key to User
	User       users.UserModel `gorm:"foreignKey:UserID" json:"user"`            // Relationship with UserModel
	ItemID     uint            `gorm:"not null" json:"item_id"`                  // Foreign key to Item
	Item       items.ItemModel `gorm:"foreignKey:ItemID" json:"item"`            // Relationship with ItemModel
	Quantity   uint            `gorm:"not null" json:"quantity"`                 // Quantity of the item ordered
	TotalPrice float32         `gorm:"not null" json:"total_price"`              // Total price for the order
	Status     string          `gorm:"default:'pending';not null" json:"status"` // Order status (pending, completed, canceled)
}

// AutoMigrate ensures the schema is updated in the database
func AutoMigrate() {
	db := database.GetConnection()
	db.AutoMigrate(&OrderModel{})
}

// SaveOrder saves a new order to the database
// SaveOrder saves a new order to the database
func SaveOrder(order *OrderModel) error {
	db := database.GetConnection()

	// Check if the associated User and Item exist before saving
	var user users.UserModel
	if err := db.First(&user, order.UserID).Error; err != nil {
		return fmt.Errorf("User with ID %d not found", order.UserID)
	}

	var item items.ItemModel
	if err := db.First(&item, order.ItemID).Error; err != nil {
		return fmt.Errorf("Item with ID %d not found", order.ItemID)
	}

	return db.Create(order).Error
}

// GetOrderByID retrieves a single order by its ID
func GetOrderByID(id uint) (OrderModel, error) {
	db := database.GetConnection()
	var order OrderModel
	err := db.Preload("User").Preload("Item").First(&order, id).Error
	return order, err
}

// GetAllOrders retrieves all orders from the database
func GetAllOrders() ([]OrderModel, error) {
	db := database.GetConnection()
	var orders []OrderModel
	err := db.Preload("User").Preload("Item").Find(&orders).Error
	return orders, err
}

// UpdateOrder updates an order by its ID
func UpdateOrder(id uint, updatedData map[string]interface{}) error {
	db := database.GetConnection()
	return db.Model(&OrderModel{}).Where("id = ?", id).Updates(updatedData).Error
}

// DeleteOrder deletes an order by its ID
func DeleteOrder(id uint) error {
	db := database.GetConnection()
	return db.Where("id = ?", id).Delete(&OrderModel{}).Error
}
