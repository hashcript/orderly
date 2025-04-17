package items

import (
	"errors"
	"fmt"

	"com.orderly/database"
)

type ItemModel struct {
	ID        uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string  `gorm:"column:name;type:varchar(255);not null" json:"name"`                        // Name of the item, used by the merchant to identify the item
	Unit      string  `gorm:"column:unit;type:varchar(255);not null" json:"unit"`                        // The unit of measurement for the item (e.g., piece, kg, liter)
	UnitPrice float32 `gorm:"column:unit_price;type:decimal(10,2);not null" json:"unit_price"`           // Price per unit of the item, formatted to two decimal places
	Status    string  `gorm:"column:status;type:varchar(20);default:'available';not null" json:"status"` // Status of the item (e.g., available, out of stock)
	Category  string  `gorm:"column:category;type:varchar(20);default:'any';not null" json:"category"`   // Status of the item (e.g., available, out of stock), category, breakfast,lumch etc

}

// Migrate the schema to the database if needed
func AutoMigrate() {
	db := database.GetConnection()
	db.AutoMigrate(&ItemModel{})
}

// SaveItem creates a new withdrawal request in the database.
func SaveItem(data *ItemModel) error {
	fmt.Println(data)
	db := database.GetConnection()
	return db.Create(data).Error
}

// GetItemByID retrieves a withdrawal request by its ID.
func GetItemByID(id uint) (ItemModel, error) {
	db := database.GetConnection()
	var request ItemModel
	err := db.First(&request, id).Error
	return request, err
}

// GetAllItems retrieves all withdrawal requests from the database.
func GetAllItems() ([]ItemModel, error) {
	db := database.GetConnection()
	var requests []ItemModel
	err := db.Find(&requests).Error
	return requests, err
}

// GetItemsByUnitPrice retrieves withdrawal requests by their UnitPrice.
func GetItemsByCategory(Category string) ([]ItemModel, error) {
	db := database.GetConnection()
	var requests []ItemModel
	err := db.Where("category = ?", Category).Find(&requests).Error
	return requests, err
}

// UpdateItem updates an existing withdrawal request by its ID.
func UpdateItem(id uint, updatedData map[string]interface{}) error {
	db := database.GetConnection()
	return db.Model(&ItemModel{}).Where("id = ?", id).Updates(updatedData).Error
}

// fancy update method
func UpdateSingleItem(model *ItemModel, data interface{}) error {
	db := database.GetConnection()

	// fmt.Printf("Model Before Update: %+v\n", model)
	// fmt.Printf("Data to Update: %+v\n", data)

	result := db.Model(model).Updates(data)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

// DeleteItem deletes a withdrawal request by its ID.
func DeleteItem(id uint) error {
	db := database.GetConnection()
	return db.Where("id = ?", id).Delete(&ItemModel{}).Error
}

// get Item by recipinet //mwexhNR
func GetItemsByMerchantID(Name string) ([]ItemModel, error) {
	db := database.GetConnection()
	var requests []ItemModel
	err := db.Where("Name = ?", Name).Find(&requests).Error
	return requests, err
}
