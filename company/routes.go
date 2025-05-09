package company

import (
	"net/http"
	"strconv"

	"com.orderly/auth"
	"com.orderly/database"
	"github.com/gin-gonic/gin"
)

// Create sets up the company-related routes
func Create(router *gin.RouterGroup) {
	protected := router.Group("/")
	protected.Use(auth.AuthMiddleware())
	protected.POST("/create", createCompany)
	protected.GET("/read/:id", ReadSingleCompany)
	protected.GET("/list", ListCompanies)
	protected.PUT("/update/:id", UpdateCompanyHandler)
	protected.DELETE("/delete/:id", DeleteCompanyHandler)
}

// CreateCompany handles creating a new company
func createCompany(c *gin.Context) {
	validator := NewCompanyModelValidator()

	if err := validator.Bind(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	if err := CreateCompany(&validator.CompanyModel); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save company", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Company created successfully", "company": validator.CompanyModel})
}

// ReadSingleCompany fetches a company by ID
func ReadSingleCompany(c *gin.Context) {
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Company ID"})
		return
	}

	company, err := GetCompanyByID(uint(idUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, database.NewError("Company", err))
		return
	}

	serializer := NewCompanySerializer(c, company)
	c.JSON(http.StatusOK, gin.H{"company": serializer.Response()})
}

// ListCompanies retrieves all companies
func ListCompanies(c *gin.Context) {
	companies, err := GetAllCompanies()
	if err != nil {
		c.JSON(http.StatusInternalServerError, database.NewError("Company", err))
		return
	}

	serializer := NewCompanySerializerList(c, companies)
	c.JSON(http.StatusOK, gin.H{"companies": serializer.Response()})
}

// UpdateCompanyHandler updates a company by ID
func UpdateCompanyHandler(c *gin.Context) {
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Company ID"})
		return
	}

	company, err := GetCompanyByID(uint(idUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, database.NewError("Company", err))
		return
	}

	var updateReq struct {
		Name     string `form:"name" json:"name" binding:"required"`
		Industry string `form:"industry" json:"industry" binding:"required"`
		Status   string `form:"status" json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	updatedData := map[string]interface{}{
		"name":     updateReq.Name,
		"industry": updateReq.Industry,
		"status":   updateReq.Status,
	}

	if err := UpdateSingleCompany(&company, updatedData); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	serializer := NewCompanySerializer(c, company)
	c.JSON(http.StatusOK, gin.H{"company": serializer.Response()})
}

// DeleteCompanyHandler deletes a company by ID
func DeleteCompanyHandler(c *gin.Context) {
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Company ID"})
		return
	}

	company, err := GetCompanyByID(uint(idUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, database.NewError("Company", err))
		return
	}

	if err := DeleteCompany(company.ID); err != nil {
		c.JSON(http.StatusUnprocessableEntity, database.NewError("database", err))
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "Company deleted successfully"})
}
