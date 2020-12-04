package endpoints

import (
	"fmt"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/models"
)

type Service string

const (
	GlobalSearch       Service = "global-search"
	GlobalTagging      Service = "global-tagging"
	AccountManagement  Service = "account-management"
	UserManagement     Service = "user-management"
	Billing            Service = "billing"
	Enterprise         Service = "enterprise"
	ResourceController Service = "resource-controller"
	ResourceCatalog    Service = "global-catalog"
)

func (s Service) String() string {
	return string(s)
}

var endpointsMapping = map[Service]models.Endpoints{
	GlobalSearch: models.Endpoints{
		PublicEndpoint:  "https://api.global-search-tagging.<cloud_domain>",
		PrivateEndpoint: "https://api.private.<region>.global-search-tagging.<cloud_domain>",
	},
	GlobalTagging: models.Endpoints{
		PublicEndpoint:  "https://tags.global-search-tagging.<cloud_domain>",
		PrivateEndpoint: "https://tags.private.<region>.global-search-tagging.<cloud_domain>",
	},
	AccountManagement: models.Endpoints{
		PublicEndpoint:  "https://accounts.<cloud_domain>",
		PrivateEndpoint: "https://private.<region>.accounts.<cloud_domain>",
	},
	UserManagement: models.Endpoints{
		PublicEndpoint:  "https://user-management.<cloud_domain>",
		PrivateEndpoint: "https://private.<region>.user-management.<cloud_domain>",
	},
	Billing: models.Endpoints{
		PublicEndpoint:  "https://billing.<cloud_domain>",
		PrivateEndpoint: "https://private.<region>.billing.<cloud_domain>",
	},
	Enterprise: models.Endpoints{
		PublicEndpoint:  "https://enterprise.<cloud_domain>",
		PrivateEndpoint: "https://private.<region>.enterprise.<cloud_domain>",
	},
	ResourceController: models.Endpoints{
		PublicEndpoint:  "https://resource-controller.<cloud_domain>",
		PrivateEndpoint: "https://private.<region>.resource-controller.<cloud_domain>",
	},
	ResourceCatalog: models.Endpoints{
		PublicEndpoint:  "https://globalcatalog.<cloud_domain>",
		PrivateEndpoint: "https://private.<region>.globalcatalog.<cloud_domain>",
	},
}

func Endpoint(svc Service, cloudDomain, region string, private bool) (string, error) {
	var endpoint string
	if endpoints, found := endpointsMapping[svc]; found {
		if private {
			endpoint = endpoints.PrivateEndpoint
		} else {
			endpoint = endpoints.PublicEndpoint
		}
	}
	if endpoint == "" {
		return "", fmt.Errorf("the endpoint of service '%s' was unknown", svc)
	}

	// replace <cloud_domain>
	if cloudDomain == "" {
		return "", fmt.Errorf("the cloud domain is empty")
	}
	endpoint = strings.ReplaceAll(endpoint, "<cloud_domain>", cloudDomain)

	// replace <region>
	if region != "" {
		endpoint = strings.ReplaceAll(endpoint, "<region>", region)
	} else if strings.Index(endpoint, "<region>") >= 0 {
		return "", fmt.Errorf("region is required to get the endpoint of service '%s'", svc)
	}

	return endpoint, nil
}
