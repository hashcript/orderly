package items

import (
	"net/http"
	"strconv"

	"com.orderly/auth"
	"com.orderly/database"
	"github.com/gin-gonic/gin"
)

// Create function to set up the routes for Item
func Create(router *gin.RouterGroup) {
	// Apply the AuthMiddleware to all routes
	protected := router.Group("/")
	protected.Use(auth.AuthMiddleware())

	protected.GET("/read/:id", ReadSingleItem)
	protected.GET("/list", ListItems)
	protected.GET("/list/merchant/:merchant", ReadItemByMerchantID)
	protected.GET("/category/:status", ReadItemByCategory)
	protected.PUT("/update/:id", UpdateItemHandler)
	protected.DELETE("/delete/:id", DeleteItemHandler)
	protected.POST("/create", CreateItem)
}

// CreateItem handles the creation of a new withdrawal request.
func CreateItem(c *gin.Context) {
	// Create a new instance of the validator
	validator := NewItemModelValidator()

	// Bind the request data to the validator (this calls the Bind method)
	if err := validator.Bind(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	// Save the withdrawal request model to the database
	if err := SaveItem(&validator.ItemModel); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save withdrawal request", "details": err.Error()})
		return
	}

	// Respond with the created withdrawal request
	c.JSON(http.StatusOK, gin.H{"message": "Item request created successfully", "request": validator.ItemModel})
}

func ReadSingleItem(c *gin.Context) {
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Withdrawal Request ID"})
		return
	}

	Item, err := GetItemByID(uint(idUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, database.NewError("Item", err))
		return
	}

	serializer := NewItemSerializer(c, Item)
	c.JSON(http.StatusOK, gin.H{"Item": serializer.Response()})
}

// ListItem retrieves a list of all withdrawal requests.
func ListItems(c *gin.Context) {
	Item, err := GetAllItems()
	if err != nil {
		c.JSON(http.StatusInternalServerError, database.NewError("Item", err))
		return
	}

	// Serialize and respond with the withdrawal requests
	var response []map[string]interface{}
	for _, request := range Item {
		serializer := NewItemSerializer(c, request)
		response = append(response, serializer.Response())
	}

	c.JSON(http.StatusOK, gin.H{"Item": response})
}

// ReadItemByCategory retrieves withdrawal requests by their status.
func ReadItemByCategory(c *gin.Context) {
	status := c.Param("status")
	if status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status is required"})
		return
	}

	Item, err := GetItemsByCategory(status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve withdrawal requests", "details": err.Error()})
		return
	}

	// Serialize and respond with the withdrawal requests
	var response []map[string]interface{}
	for _, request := range Item {
		serializer := NewItemSerializer(c, request)
		response = append(response, serializer.Response())
	}

	c.JSON(http.StatusOK, gin.H{"Item": response})
}

// UpdateItemHandler updates a withdrawal request by ID.
// UpdateItemHandler updates a withdrawal request by ID.
func UpdateItemHandler(c *gin.Context) {
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Withdrawal Request ID"})
		return
	}

	Item, err := GetItemByID(uint(idUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, database.NewError("Item", err))
		return
	}

	// Extract fields to update
	type UpdateRequest struct {
		Name      string  `form:"name" json:"name" binding:"required"`
		Unit      string  `form:"unit" json:"unit" binding:"required"`
		UnitPrice float32 `form:"UnitPrice" json:"UnitPrice" binding:"required"`
		Status    string  `form:"status" json:"status" binding:"required"`
		Category  string  `form:"Category" json:"Category"`
	}
	var updateReq UpdateRequest

	// Parse JSON into the struct
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	// Log extracted data

	// Prepare data for update
	updatedData := map[string]interface{}{
		"status":     updateReq.Status,
		"name":       updateReq.Name,
		"unit_price": updateReq.UnitPrice,
		"unit":       updateReq.Unit,
		"category":   updateReq.Category,
	}

	// Update the withdrawal request
	if err := UpdateSingleItem(&Item, updatedData); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	// Serialize the updated request
	serializer := NewItemSerializer(c, Item)
	c.JSON(http.StatusOK, gin.H{"Item": serializer.Response()})
}

// DeleteItemHandler deletes a withdrawal request by ID.
func DeleteItemHandler(c *gin.Context) {
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Withdrawal Request ID"})
		return
	}

	Item, err := GetItemByID(uint(idUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, database.NewError("Item", err))
		return
	}

	if err := DeleteItem(Item.ID); err != nil {
		c.JSON(http.StatusUnprocessableEntity, database.NewError("database", err))
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "Withdrawal Request Deleted"})
}

func ReadItemByMerchantID(c *gin.Context) {
	merchant := c.Param("merchant")
	if merchant == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status is required"})
		return
	}

	Item, err := GetItemsByMerchantID(merchant)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve withdrawal requests", "details": err.Error()})
		return
	}

	// Serialize and respond with the withdrawal requests
	var response []map[string]interface{}
	for _, request := range Item {
		serializer := NewItemSerializer(c, request)
		response = append(response, serializer.Response())
	}

	c.JSON(http.StatusOK, gin.H{"Item": response})
}
