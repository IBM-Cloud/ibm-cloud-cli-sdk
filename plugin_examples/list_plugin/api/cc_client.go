package api

import (
	"fmt"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/common/rest"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin_examples/list_plugin/models"
)

type CCError struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
}

func (c CCError) Error() string {
	return fmt.Sprintf("Error response from server. Status code: %d; description: %s.", c.Code, c.Description)
}

type CCClient interface {
	AppsAndServices(spaceId string) (models.AppsAndServices, error)
	OrgUsage(orgId string) (models.OrgUsage, error)
}

type ccClient struct {
	endpoint string
	client   *rest.Client
}

func NewCCClient(endpoint string, client *rest.Client) CCClient {
	return &ccClient{
		endpoint: endpoint,
		client:   client,
	}
}

func (c *ccClient) AppsAndServices(spaceId string) (models.AppsAndServices, error) {
	var summary models.AppsAndServices
	var apiErr CCError

	req := rest.GetRequest(c.endpoint + fmt.Sprintf("/v2/spaces/%s/summary", spaceId))
	_, err := c.client.Do(req, &summary, &apiErr)
	return summary, c.handleError(err, apiErr)
}

func (c *ccClient) OrgUsage(orgId string) (models.OrgUsage, error) {
	var summary models.OrgUsage
	var apiErr CCError

	req := rest.GetRequest(c.endpoint + fmt.Sprintf("/v2/organizations/%s/summary", orgId))
	_, err := c.client.Do(req, &summary, &apiErr)
	return summary, c.handleError(err, apiErr)
}

func (c *ccClient) handleError(httpErr error, apiErr CCError) error {
	if httpErr != nil {
		return httpErr
	}
	if apiErr != (CCError{}) {
		return apiErr
	}
	return nil
}
