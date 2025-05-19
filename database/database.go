package database

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator"
	driver "github.com/go-sql-driver/mysql" // this is the MySQL driver
	"gorm.io/driver/mysql"                  // this is GORM's wrapper, keep name as-is
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() *gorm.DB {
	// Load Aiven CA certificate
	rootCertPool := x509.NewCertPool()
	pem, err := os.ReadFile("ca.pem") // Ensure this file exists
	if err != nil {
		log.Fatal("Failed to read CA cert:", err)
	}
	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		log.Fatal("Failed to append PEM.")
	}

	// Register custom TLS config with go-sql-driver/mysql
	err = driver.RegisterTLSConfig("custom", &tls.Config{
		RootCAs: rootCertPool,
	})
	if err != nil {
		log.Fatal("Failed to register TLS config:", err)
	}

	// Construct DSN using the custom TLS config
	dsn := "avnadmin:AVNS_z6fAW0ja8h5in-GeGez@tcp(mysql-365444d5-compliace.g.aivencloud.com:25391)/defaultdb?tls=custom&charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	DB = db
	return DB
}

// Function to get the database connection
func GetConnection() *gorm.DB {
	return DB
}

// Binder
func Bind(c *gin.Context, object interface{}) error {
	binder := binding.Default(c.Request.Method, c.ContentType())
	return c.ShouldBindWith(object, binder)
}

// Return customized error info
type CommonError struct {
	Errors map[string]interface{} `json:"errors"`
}

// Validators

func NewValidatorError(err error) CommonError {
	res := CommonError{}
	res.Errors = make(map[string]interface{})
	errs := err.(validator.ValidationErrors)
	for _, v := range errs {
		res1 := fmt.Sprintf("{%v: %v}", v.Tag(), v.Param())
		fmt.Println(res1)
		if v.Param() != "" {

			res.Errors[v.Field()] = fmt.Sprintf("{%v: %v}", v.Tag(), v.Param())
		} else {
			res.Errors[v.Field()] = fmt.Sprintf("{key: %v}", v.Tag())
		}
	}
	return res
}

// Wrapping error into an object
func NewError(key string, err error) CommonError {
	res := CommonError{}
	res.Errors = make(map[string]interface{})
	res.Errors[key] = err.Error()
	return res
}
