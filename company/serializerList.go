package company


import (
	"github.com/gin-gonic/gin"
)

type CompanySerializerList struct {
	c        *gin.Context
	Companies []Company
}

func NewCompanySerializerList(c *gin.Context, companies []Company) *CompanySerializerList {
	return &CompanySerializerList{c: c, Companies: companies}
}

func (s *CompanySerializerList) Response() []map[string]interface{} {
	var responseList []map[string]interface{}
	for _, company := range s.Companies {
		serializer := NewCompanySerializer(s.c, company)
		responseList = append(responseList, serializer.Response())
	}
	return responseList
}
