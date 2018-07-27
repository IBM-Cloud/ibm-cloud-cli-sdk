package crn

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type CRNTestSuite struct {
	suite.Suite
}

func TestCRNSuite(t *testing.T) {
	suite.Run(t, new(CRNTestSuite))
}

func (suite *CRNTestSuite) TestParse() {
	crnString := "crn:v1:bluemix:public:test-service::a/account-guid::deployment:deployment-guid"
	crn, err := Parse(crnString)
	suite.NoError(err)
	suite.Equal("test-service", crn.ServiceName)
	suite.Equal("", crn.Region)
	suite.Equal(ScopeAccount, crn.ScopeType)
	suite.Equal("account-guid", crn.Scope)
	suite.Equal("a/account-guid", crn.ScopeSegment())
	suite.Equal(ResourceTypeDeployment, crn.ResourceType)
	suite.Equal("deployment-guid", crn.Resource)

	crnString = "crn:v1:bluemix:public:test-service:us-south:global::deployment:deployment-guid"
	crn, err = Parse(crnString)
	suite.NoError(err)
	suite.Equal("test-service", crn.ServiceName)
	suite.Equal("us-south", crn.Region)
	suite.Equal("", crn.ScopeType)
	suite.Equal("global", crn.Scope)
	suite.Equal("global", crn.ScopeSegment())
	suite.Equal(ResourceTypeDeployment, crn.ResourceType)
	suite.Equal("deployment-guid", crn.Resource)

	crnString = "crn:v1:bluemix:public:test-service:us-south:error::deployment:deployment-guid"
	_, err = Parse(crnString)
	suite.Error(err)
}
