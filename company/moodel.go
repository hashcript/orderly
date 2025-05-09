package company

import (
	"errors"

	"com.orderly/database"
	"gorm.io/gorm"
	"gorm.io/datatypes"
)

type Company struct {
	gorm.Model
	CompanyName                   string         `gorm:"type:varchar(255);not null" json:"company_name"`
	MerchantURL                   string         `gorm:"type:varchar(255);not null" json:"merchant_url"`
	DateOfIncorporation           string         `gorm:"type:varchar(100);not null" json:"date_of_incorporation"`
	CompanyRegNumber              string         `gorm:"type:varchar(100);not null" json:"company_reg_number"`
	CountryOfIncorporation        string         `gorm:"type:varchar(100);not null" json:"country_of_incorporation"`
	CompanyEmail                  string         `gorm:"type:varchar(100);not null" json:"company_email"`
	ContactPersonFullName         string         `gorm:"type:varchar(100);not null" json:"contact_person_full_name"`
	ContactPersonTelephoneNumber  string         `gorm:"type:varchar(20);not null" json:"contact_person_telephone"`
	ContactPersonEmail            string         `gorm:"type:varchar(100);not null" json:"contact_person_email"`
	BusinessDescription           string         `gorm:"type:varchar(255);not null" json:"business_description"`
	CompanySourceOfFunds          string         `gorm:"type:varchar(255);not null" json:"company_source_of_funds"`
	BusinessRelationshipPurpose   string         `gorm:"type:varchar(255);not null" json:"business_relationship_purpose"`
	IsCompanySubjectToLicensing   bool           `gorm:"not null" json:"is_licensed"`
	BankName                      string         `gorm:"type:varchar(100);not null" json:"bank_name"`
	SwiftCode                     string         `gorm:"type:varchar(50);not null" json:"swift_code"`
	TargetMarkets                 datatypes.JSON `gorm:"type:json" json:"target_markets"` // e.g. ["EU", "ASIA"]
	CountriesOfOperation          datatypes.JSON `gorm:"type:json" json:"countries_of_operation"`
	PreviouslyUsedPaymentGateways datatypes.JSON `gorm:"type:json" json:"previously_used_payment_gateways"`
}

// AutoMigrate migrates the company schema
func AutoMigrate() {
	db := database.GetConnection()
	db.AutoMigrate(&Company{})
}

// CreateCompany inserts a new company record
func CreateCompany(data *Company) error {
	db := database.GetConnection()
	return db.Create(data).Error
}

// GetCompanyByID returns a single company by ID
func GetCompanyByID(id uint) (Company, error) {
	db := database.GetConnection()
	var company Company
	err := db.First(&company, id).Error
	return company, err
}

// GetAllCompanies returns all company records
func GetAllCompanies() ([]Company, error) {
	db := database.GetConnection()
	var companies []Company
	err := db.Find(&companies).Error
	return companies, err
}

// UpdateCompany updates a company record by ID
func UpdateCompany(id uint, updatedData map[string]interface{}) error {
	db := database.GetConnection()
	return db.Model(&Company{}).Where("id = ?", id).Updates(updatedData).Error
}

// Fancy update method
func UpdateSingleCompany(model *Company, data interface{}) error {
	db := database.GetConnection()
	result := db.Model(model).Updates(data)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

// DeleteCompany deletes a company by ID
func DeleteCompany(id uint) error {
	db := database.GetConnection()
	return db.Delete(&Company{}, id).Error
}
