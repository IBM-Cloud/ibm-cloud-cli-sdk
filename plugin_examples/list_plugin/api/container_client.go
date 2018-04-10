package api

import (
	"fmt"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/common/rest"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin_examples/list_plugin/models"
)

type ContainerError struct {
	Code        string `json:"code"`
	StatusCode  string `json:"rc"`
	Description string `json:"description"`
	IncidentId  string `json:"incident_id"`
}

func (e ContainerError) Error() string {
	return fmt.Sprintf("Server error, status code: %s, error code: %s, incident id: %s, message: %s",
		e.StatusCode, e.Code, e.IncidentId, e.Description)
}

type ContainerClient interface {
	Containers(spaceId string) ([]models.Container, error)
	ContainersQuotaAndUsage(spaceId string) (models.ContainersQuotaAndUsage, error)
}

type containerClient struct {
	endpoint string
	client   *rest.Client
}

func NewContainerClient(endpoint string, client *rest.Client) ContainerClient {
	return &containerClient{
		endpoint: endpoint,
		client:   client,
	}
}

func (c *containerClient) Containers(spaceId string) ([]models.Container, error) {
	var containers []models.Container
	var apiErr ContainerError

	req := rest.GetRequest(c.endpoint+"/v3/containers/json").Add("X-Auth-Project-Id", spaceId).Query("all", "true")
	_, err := c.client.Do(req, &containers, &apiErr)
	return containers, c.handleError(err, apiErr)
}

func (c *containerClient) ContainersQuotaAndUsage(spaceId string) (models.ContainersQuotaAndUsage, error) {
	var quotaAndUsage models.ContainersQuotaAndUsage
	var apiErr ContainerError

	req := rest.GetRequest(c.endpoint+"/v3/containers/usage").Add("X-Auth-Project-Id", spaceId)
	_, err := c.client.Do(req, &quotaAndUsage, &apiErr)
	return quotaAndUsage, c.handleError(err, apiErr)
}

func (c *containerClient) handleError(httpErr error, apiErr ContainerError) error {
	if httpErr != nil {
		return httpErr
	}
	if apiErr != (ContainerError{}) {
		return apiErr
	}
	return nil
}
