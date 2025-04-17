package items

import (
	"com.orderly/database"
	"github.com/gin-gonic/gin"
)


type ItemModelValidator struct {
	Item struct {
		Name      string  `form:"name" json:"name" binding:"required"`
		Unit      string  `form:"unit" json:"unit" binding:"required"`
		UnitPrice float32 `form:"UnitPrice" json:"UnitPrice" binding:"required"`
		Status    string  `form:"status" json:"status" binding:"required"`
		Category  string  `form:"Category" json:"Category"`
	} `json:"Item"`
	ItemModel ItemModel `json:"-"`
}

// Bind binds the Item data to the ItemModelValidator.
func (v *ItemModelValidator) Bind(c *gin.Context) error {
	err := database.Bind(c, v)
	if err != nil {
		return err
	}

	v.ItemModel.Name = v.Item.Name
	v.ItemModel.Unit = v.Item.Unit
	v.ItemModel.UnitPrice = v.Item.UnitPrice
	v.ItemModel.Status = v.Item.Status
	v.ItemModel.Category = v.Item.Category

	// Check values after assignment

	return nil
}

// NewItemModelValidator creates a new ItemModelValidator instance.
func NewItemModelValidator() ItemModelValidator {
	return ItemModelValidator{}
}

// NewItemModelValidatorFillWith creates a new ItemModelValidator instance and fills it with Item model data.
func NewItemModelValidatorFillWith(ItemModel ItemModel) ItemModelValidator {
	return ItemModelValidator{
		Item: struct {
			Name      string  `form:"name" json:"name" binding:"required"`
			Unit      string  `form:"unit" json:"unit" binding:"required"`
			UnitPrice float32 `form:"UnitPrice" json:"UnitPrice" binding:"required"`
			Status    string  `form:"status" json:"status" binding:"required"`
			Category  string  `form:"Category" json:"Category"`
		}{
			Name:      ItemModel.Name,
			Unit:      ItemModel.Unit,
			UnitPrice: ItemModel.UnitPrice,
			Status:    ItemModel.Status,
			Category:  ItemModel.Category,
		},
	}
}
