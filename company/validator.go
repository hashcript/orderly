package company

import (
	"encoding/json"
	"errors"

	"com.orderly/database"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

// CompanyModelValidator handles validation for incoming company data
type CompanyModelValidator struct {
	Company struct {
		CompanyName                   string          `json:"company_name" binding:"required"`
		MerchantURL                   string          `json:"merchant_url" binding:"required"`
		DateOfIncorporation           string          `json:"date_of_incorporation" binding:"required"`
		CompanyRegNumber              string          `json:"company_reg_number" binding:"required"`
		CountryOfIncorporation        string          `json:"country_of_incorporation" binding:"required"`
		CompanyEmail                  string          `json:"company_email" binding:"required,email"`
		ContactPersonFullName         string          `json:"contact_person_full_name" binding:"required"`
		ContactPersonTelephoneNumber  string          `json:"contact_person_telephone" binding:"required"`
		ContactPersonEmail            string          `json:"contact_person_email" binding:"required,email"`
		BusinessDescription           string          `json:"business_description" binding:"required"`
		CompanySourceOfFunds          string          `json:"company_source_of_funds" binding:"required"`
		BusinessRelationshipPurpose   string          `json:"business_relationship_purpose" binding:"required"`
		IsCompanySubjectToLicensing   bool            `json:"is_licensed" binding:"required"`
		BankName                      string          `json:"bank_name" binding:"required"`
		SwiftCode                     string          `json:"swift_code" binding:"required"`
		TargetMarkets                 []string        `json:"target_markets" binding:"required"`
		CountriesOfOperation          []string        `json:"countries_of_operation" binding:"required"`
		PreviouslyUsedPaymentGateways []string        `json:"previously_used_payment_gateways"`
	} `json:"company"`
	CompanyModel Company `json:"-"`
}

// Bind binds and validates the request body to CompanyModel
func (v *CompanyModelValidator) Bind(c *gin.Context) error {
	err := database.Bind(c, v)
	if err != nil {
		return err
	}

	// Convert slices to JSON
	targetMarketsJSON, err := json.Marshal(v.Company.TargetMarkets)
	if err != nil {
		return errors.New("invalid target_markets format")
	}

	countriesJSON, err := json.Marshal(v.Company.CountriesOfOperation)
	if err != nil {
		return errors.New("invalid countries_of_operation format")
	}

	paymentGatewaysJSON, err := json.Marshal(v.Company.PreviouslyUsedPaymentGateways)
	if err != nil {
		return errors.New("invalid previously_used_payment_gateways format")
	}

	v.CompanyModel = Company{
		CompanyName:                   v.Company.CompanyName,
		MerchantURL:                   v.Company.MerchantURL,
		DateOfIncorporation:           v.Company.DateOfIncorporation,
		CompanyRegNumber:              v.Company.CompanyRegNumber,
		CountryOfIncorporation:        v.Company.CountryOfIncorporation,
		CompanyEmail:                  v.Company.CompanyEmail,
		ContactPersonFullName:         v.Company.ContactPersonFullName,
		ContactPersonTelephoneNumber:  v.Company.ContactPersonTelephoneNumber,
		ContactPersonEmail:            v.Company.ContactPersonEmail,
		BusinessDescription:           v.Company.BusinessDescription,
		CompanySourceOfFunds:          v.Company.CompanySourceOfFunds,
		BusinessRelationshipPurpose:   v.Company.BusinessRelationshipPurpose,
		IsCompanySubjectToLicensing:   v.Company.IsCompanySubjectToLicensing,
		BankName:                      v.Company.BankName,
		SwiftCode:                     v.Company.SwiftCode,
		TargetMarkets:                 datatypes.JSON(targetMarketsJSON),
		CountriesOfOperation:          datatypes.JSON(countriesJSON),
		PreviouslyUsedPaymentGateways: datatypes.JSON(paymentGatewaysJSON),
	}

	return nil
}

// NewCompanyModelValidator creates a new CompanyModelValidator
func NewCompanyModelValidator() CompanyModelValidator {
	return CompanyModelValidator{}
}

// NewCompanyModelValidatorFillWith creates a validator filled with existing data
func NewCompanyModelValidatorFillWith(model Company) CompanyModelValidator {
	var targetMarkets, countries, gateways []string
	json.Unmarshal(model.TargetMarkets, &targetMarkets)
	json.Unmarshal(model.CountriesOfOperation, &countries)
	json.Unmarshal(model.PreviouslyUsedPaymentGateways, &gateways)

	return CompanyModelValidator{
		Company: struct {
			CompanyName                   string   `json:"company_name" binding:"required"`
			MerchantURL                   string   `json:"merchant_url" binding:"required"`
			DateOfIncorporation           string   `json:"date_of_incorporation" binding:"required"`
			CompanyRegNumber              string   `json:"company_reg_number" binding:"required"`
			CountryOfIncorporation        string   `json:"country_of_incorporation" binding:"required"`
			CompanyEmail                  string   `json:"company_email" binding:"required,email"`
			ContactPersonFullName         string   `json:"contact_person_full_name" binding:"required"`
			ContactPersonTelephoneNumber  string   `json:"contact_person_telephone" binding:"required"`
			ContactPersonEmail            string   `json:"contact_person_email" binding:"required,email"`
			BusinessDescription           string   `json:"business_description" binding:"required"`
			CompanySourceOfFunds          string   `json:"company_source_of_funds" binding:"required"`
			BusinessRelationshipPurpose   string   `json:"business_relationship_purpose" binding:"required"`
			IsCompanySubjectToLicensing   bool     `json:"is_licensed" binding:"required"`
			BankName                      string   `json:"bank_name" binding:"required"`
			SwiftCode                     string   `json:"swift_code" binding:"required"`
			TargetMarkets                 []string `json:"target_markets" binding:"required"`
			CountriesOfOperation          []string `json:"countries_of_operation" binding:"required"`
			PreviouslyUsedPaymentGateways []string `json:"previously_used_payment_gateways"`
		}{
			CompanyName:                   model.CompanyName,
			MerchantURL:                   model.MerchantURL,
			DateOfIncorporation:           model.DateOfIncorporation,
			CompanyRegNumber:              model.CompanyRegNumber,
			CountryOfIncorporation:        model.CountryOfIncorporation,
			CompanyEmail:                  model.CompanyEmail,
			ContactPersonFullName:         model.ContactPersonFullName,
			ContactPersonTelephoneNumber:  model.ContactPersonTelephoneNumber,
			ContactPersonEmail:            model.ContactPersonEmail,
			BusinessDescription:           model.BusinessDescription,
			CompanySourceOfFunds:          model.CompanySourceOfFunds,
			BusinessRelationshipPurpose:   model.BusinessRelationshipPurpose,
			IsCompanySubjectToLicensing:   model.IsCompanySubjectToLicensing,
			BankName:                      model.BankName,
			SwiftCode:                     model.SwiftCode,
			TargetMarkets:                 targetMarkets,
			CountriesOfOperation:          countries,
			PreviouslyUsedPaymentGateways: gateways,
		},
	}
}
