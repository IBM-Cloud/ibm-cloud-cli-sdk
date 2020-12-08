package endpoints_test

import (
	"fmt"
	"testing"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/endpoints"
	"github.com/stretchr/testify/assert"
)

func TestEndpointUnknownService(t *testing.T) {
	_, err := Endpoint(Service("unknown"), "cloud.ibm.com", "us-south", false)
	assert.Error(t, err, "an error is expected")
}

func TestEndpointEmptyCloudDomain(t *testing.T) {
	_, err := Endpoint(AccountManagement, "", "", false)
	assert.Error(t, err, "an error is expected")
}

func TestEndpointPublic(t *testing.T) {
	cloudDomain := "cloud.ibm.com"
	region := ""
	private := false

	endpoints := map[Service]string{
		GlobalSearch:       fmt.Sprintf("https://api.global-search-tagging.%s", cloudDomain),
		GlobalTagging:      fmt.Sprintf("https://tags.global-search-tagging.%s", cloudDomain),
		AccountManagement:  fmt.Sprintf("https://accounts.%s", cloudDomain),
		UserManagement:     fmt.Sprintf("https://user-management.%s", cloudDomain),
		Billing:            fmt.Sprintf("https://billing.%s", cloudDomain),
		Enterprise:         fmt.Sprintf("https://enterprise.%s", cloudDomain),
		ResourceController: fmt.Sprintf("https://resource-controller.%s", cloudDomain),
		ResourceCatalog:    fmt.Sprintf("https://globalcatalog.%s", cloudDomain),
	}
	for svc, expected := range endpoints {
		actual, err := Endpoint(svc, cloudDomain, region, private)
		assert.NoError(t, err, "public endpoint of service '%s'", svc)
		assert.Equal(t, expected, actual, "public endpoint of service '%s'", svc)
	}
}

func TestEndpointPrivate(t *testing.T) {
	cloudDomain := "cloud.ibm.com"
	region := "us-south"
	private := true

	endpoints := map[Service]string{
		GlobalSearch:       fmt.Sprintf("https://api.private.%s.global-search-tagging.%s", region, cloudDomain),
		GlobalTagging:      fmt.Sprintf("https://tags.private.%s.global-search-tagging.%s", region, cloudDomain),
		AccountManagement:  fmt.Sprintf("https://private.%s.accounts.%s", region, cloudDomain),
		UserManagement:     fmt.Sprintf("https://private.%s.user-management.%s", region, cloudDomain),
		Billing:            fmt.Sprintf("https://private.%s.billing.%s", region, cloudDomain),
		Enterprise:         fmt.Sprintf("https://private.%s.enterprise.%s", region, cloudDomain),
		ResourceController: fmt.Sprintf("https://private.%s.resource-controller.%s", region, cloudDomain),
		ResourceCatalog:    fmt.Sprintf("https://private.%s.globalcatalog.%s", region, cloudDomain),
	}
	for svc, expected := range endpoints {
		actual, err := Endpoint(svc, cloudDomain, region, private)
		assert.NoError(t, err, "private endpoint of service '%s'", svc)
		assert.Equal(t, expected, actual, "private endpoint of service '%s'", svc)
	}
}

func TestEndpointPrivateNoRegion(t *testing.T) {
	_, err := Endpoint(AccountManagement, "cloud.ibm.com", "", true)
	assert.Error(t, err, "an error is expected")
}
