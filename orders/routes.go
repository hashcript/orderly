package orders

import (
	"net/http"
	"strconv"

	"com.orderly/auth"
	"com.orderly/database"
	"github.com/gin-gonic/gin"
)

// Create function to set up the routes for Order
func Create(router *gin.RouterGroup) {
	// Apply the AuthMiddleware to all routes
	protected := router.Group("/")
	protected.Use(auth.AuthMiddleware())

	protected.POST("/create", CreateOrderHandler)
	protected.GET("/read/:id", ReadSingleOrderHandler)
	protected.GET("/list", ListOrdersHandler)
	protected.PUT("/update/:id", UpdateOrderHandler)
	protected.DELETE("/delete/:id", DeleteOrderHandler)
}

// CreateOrderHandler handles the creation of a new order
func CreateOrderHandler(c *gin.Context) {
	// Create a new instance of the validator
	validator := NewOrderValidator()

	// Bind the request data to the validator (this calls the Bind method)
	if err := validator.Bind(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	// Save the order model to the database
	if err := SaveOrder(&validator.OrderModel); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save order", "details": err.Error()})
		return
	}

	// Respond with the created order
	c.JSON(http.StatusOK, gin.H{"message": "Order created successfully", "order": validator.OrderModel})
}

// ReadSingleOrderHandler retrieves a single order by its ID
func ReadSingleOrderHandler(c *gin.Context) {
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Order ID"})
		return
	}

	order, err := GetOrderByID(uint(idUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, database.NewError("Order", err))
		return
	}

	serializer := NewOrderSerializer(c, order)
	response, err := serializer.Response()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize order", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"order": response})
}

// ListOrdersHandler retrieves a list of all orders
func ListOrdersHandler(c *gin.Context) {
	orders, err := GetAllOrders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, database.NewError("Order", err))
		return
	}

	// Serialize and respond with the orders
	var response []map[string]interface{}
	for _, order := range orders {
		serializer := NewOrderSerializer(c, order)
		resp, err := serializer.Response()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize order", "details": err.Error()})
			return
		}
		response = append(response, resp)
	}

	c.JSON(http.StatusOK, gin.H{"orders": response})
}

// UpdateOrderHandler updates an order by its ID
// UpdateOrderHandler updates an order by its ID
func UpdateOrderHandler(c *gin.Context) {
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Order ID"})
		return
	}

	// Retrieve order from DB
	order, err := GetOrderByID(uint(idUint))
	if err != nil {
		// if err == database.ErrNotFound {
		// 	c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		// 	return
		// }
		c.JSON(http.StatusInternalServerError, database.NewError("Order", err))
		return
	}

	// Extract fields to update
	var updateReq struct {
		Quantity   uint    `json:"quantity" binding:"required"`
		TotalPrice float32 `json:"total_price" binding:"required"`
		Status     string  `json:"status" binding:"required"`
	}

	// Bind the JSON to the struct
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	// Apply updates to the order model
	order.Quantity = updateReq.Quantity
	order.TotalPrice = updateReq.TotalPrice
	order.Status = updateReq.Status

	// Save the updated order in the database
	updates := map[string]interface{}{
		"quantity":    order.Quantity,
		"total_price": order.TotalPrice,
		"status":      order.Status,
	}
	if err := UpdateOrder(order.ID, updates); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	// Serialize and respond with the updated order
	serializer := NewOrderSerializer(c, order)
	response, err := serializer.Response()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize order", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"order": response})
}

// DeleteOrderHandler deletes an order by its ID
func DeleteOrderHandler(c *gin.Context) {
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Order ID"})
		return
	}

	order, err := GetOrderByID(uint(idUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, database.NewError("Order", err))
		return
	}

	if err := DeleteOrder(order.ID); err != nil {
		c.JSON(http.StatusUnprocessableEntity, database.NewError("database", err))
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "Order deleted successfully"})
}
