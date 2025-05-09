package company

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

type CompanySerializer struct {
	c            *gin.Context
	CompanyModel Company
}

func NewCompanySerializer(c *gin.Context, company Company) *CompanySerializer {
	return &CompanySerializer{c: c, CompanyModel: company}
}

func (s *CompanySerializer) Response() map[string]interface{} {
	var targetMarkets, countriesOfOperation, previouslyUsedGateways interface{}

	// Decode JSON fields (they're stored as datatypes.JSON)
	json.Unmarshal(s.CompanyModel.TargetMarkets, &targetMarkets)
	json.Unmarshal(s.CompanyModel.CountriesOfOperation, &countriesOfOperation)
	json.Unmarshal(s.CompanyModel.PreviouslyUsedPaymentGateways, &previouslyUsedGateways)

	return map[string]interface{}{
		"id":                               s.CompanyModel.ID,
		"company_name":                     s.CompanyModel.CompanyName,
		"merchant_url":                     s.CompanyModel.MerchantURL,
		"date_of_incorporation":            s.CompanyModel.DateOfIncorporation,
		"company_reg_number":               s.CompanyModel.CompanyRegNumber,
		"country_of_incorporation":         s.CompanyModel.CountryOfIncorporation,
		"company_email":                    s.CompanyModel.CompanyEmail,
		"contact_person_full_name":         s.CompanyModel.ContactPersonFullName,
		"contact_person_telephone":         s.CompanyModel.ContactPersonTelephoneNumber,
		"contact_person_email":             s.CompanyModel.ContactPersonEmail,
		"business_description":             s.CompanyModel.BusinessDescription,
		"company_source_of_funds":          s.CompanyModel.CompanySourceOfFunds,
		"business_relationship_purpose":    s.CompanyModel.BusinessRelationshipPurpose,
		"is_licensed":                      s.CompanyModel.IsCompanySubjectToLicensing,
		"bank_name":                        s.CompanyModel.BankName,
		"swift_code":                       s.CompanyModel.SwiftCode,
		"target_markets":                   targetMarkets,
		"countries_of_operation":           countriesOfOperation,
		"previously_used_payment_gateways": previouslyUsedGateways,
		"created_at":                       s.CompanyModel.CreatedAt,
		"updated_at":                       s.CompanyModel.UpdatedAt,
	}
}
