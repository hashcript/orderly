package items

import (
	"github.com/gin-gonic/gin"
)

type ItemSerializer struct {
	c         *gin.Context
	ItemModel ItemModel
}

func NewItemSerializer(c *gin.Context, Item ItemModel) *ItemSerializer {
	return &ItemSerializer{c: c, ItemModel: Item}
}

/*
	ID        uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string  `gorm:"column:name;type:varchar(255);not null" json:"name"`                        // Name of the item, used by the merchant to identify the item
	Unit      string  `gorm:"column:unit;type:varchar(255);not null" json:"unit"`                        // The unit of measurement for the item (e.g., piece, kg, liter)
	UnitPrice float32 `gorm:"column:unit_price;type:decimal(10,2);not null" json:"unit_price"`           // Price per unit of the item, formatted to two decimal places
	Status    string  `gorm:"column:status;type:varchar(20);default:'available';not null" json:"status"` // Status of the item (e.g., available, out of stock)
	Category  string  `gorm:"column:category;type:varchar(20);default:'any';not null" json:"category"`

*/

func (s *ItemSerializer) Response() map[string]interface{} {
	return map[string]interface{}{
		"id":        s.ItemModel.ID,
		"Name":      s.ItemModel.Name,
		"Unit":      s.ItemModel.Unit,
		"UnitPrice": s.ItemModel.UnitPrice,
		"status":    s.ItemModel.Status,
		"Category":  s.ItemModel.Category, // Added transferType

	}
}

type ItemSerializerList struct {
	c     *gin.Context
	Items []ItemModel
}

func NewItemSerializerList(c *gin.Context, Items []ItemModel) *ItemSerializerList {
	return &ItemSerializerList{c: c, Items: Items}
}

func (s *ItemSerializerList) Response() []map[string]interface{} {
	var responseList []map[string]interface{}
	for _, request := range s.Items {
		serializer := NewItemSerializer(s.c, request)
		responseList = append(responseList, serializer.Response())
	}
	return responseList
}
