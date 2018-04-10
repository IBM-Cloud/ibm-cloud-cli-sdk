package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin_examples/list_plugin/api"
	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin_examples/list_plugin/i18n"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin_examples/list_plugin/models"
)

type List struct {
	ui              terminal.UI
	cf              plugin.CFContext
	ccClient        api.CCClient
	containerClient api.ContainerClient
}

func NewList(
	ui terminal.UI,
	context plugin.PluginContext,
	ccClient api.CCClient,
	containerClient api.ContainerClient) *List {

	return &List{
		ui:              ui,
		cf:              context.CF(),
		ccClient:        ccClient,
		containerClient: containerClient,
	}
}

func (cmd *List) Run(args []string) error {
	err := checkTarget(cmd.cf)
	if err != nil {
		return err
	}

	orgId := cmd.cf.CurrentOrganization().GUID
	spaceId := cmd.cf.CurrentSpace().GUID

	summary, err := cmd.ccClient.AppsAndServices(spaceId)
	if err != nil {
		return fmt.Errorf(T("Unable to query apps and services of the target space:\n") + err.Error())
	}

	orgUsage, err := cmd.ccClient.OrgUsage(orgId)
	if err != nil {
		return fmt.Errorf(T("Unable to retrieve usage of the target org:\n") + err.Error())
	}

	cmd.printApps(summary.Apps, orgUsage)
	cmd.printServices(summary.Services, orgUsage)

	containers, err := cmd.containerClient.Containers(spaceId)
	if err != nil {
		return fmt.Errorf(T("Unable to query containers of the target space:\n") + err.Error())
	}

	containersQuotaAndUsage, err := cmd.containerClient.ContainersQuotaAndUsage(spaceId)
	if err != nil {
		return fmt.Errorf(T("Unable to retrieve containers' usage and quota of the target space:\n") + err.Error())
	}

	cmd.printContainers(containers, containersQuotaAndUsage)

	return nil
}

func (cmd *List) printApps(apps []models.App, orgUsage models.OrgUsage) {
	cmd.ui.Say(terminal.ColorizeBold(
		T("CloudFoundy Applications  {{.Used}}/{{.Limit}} used",
			map[string]interface{}{
				"Used":  formattedGB(orgUsage.TotalMemoryUsed()),
				"Limit": formattedGB(cmd.cf.CurrentOrganization().QuotaDefinition.InstanceMemoryLimitInMB)}), 33))

	table := cmd.ui.Table([]string{T("Name"), T("Routes"), T("Memory (MB)"), T("Instances"), T("State")})
	for _, a := range apps {
		table.Add(
			a.Name,
			strings.Join(a.URLs, "\n"),
			strconv.FormatInt(a.Memory, 10),
			fmt.Sprintf("%d/%d", a.RunningInstances, a.TotalInstances),
			a.State)
	}
	table.Print()

	cmd.ui.Say("")
}

func (cmd *List) printServices(services []models.ServiceInstance, orgUsage models.OrgUsage) {
	cmd.ui.Say(terminal.ColorizeBold(
		T("Services {{.Count}}/{{.Limit}} used", map[string]interface{}{
			"Count": orgUsage.ServicesCount(),
			"Limit": cmd.cf.CurrentOrganization().QuotaDefinition.ServicesLimit}), 33))

	table := cmd.ui.Table([]string{T("Name"), T("Service Offering"), T("Plan")})
	for _, s := range services {
		table.Add(
			s.Name,
			s.ServicePlan.ServiceOffering.Label,
			s.ServicePlan.Name)
	}
	table.Print()

	cmd.ui.Say("")
}

func (cmd *List) printContainers(containers []models.Container, quotaAndUsage models.ContainersQuotaAndUsage) {
	cmd.ui.Say(terminal.ColorizeBold(
		T("Containers  {{.MemoryUsed}}/{{.MemoryLimit}}  {{.IPCount}}/{{.IPLimit}} Public IPs Requested|{{.BoundIPCount}} Used",
			map[string]interface{}{
				"MemoryUsed":   formattedGB(quotaAndUsage.Usage.MemoryInMB),
				"MemoryLimit":  formattedGB(quotaAndUsage.Limits.MemoryLimitInMB),
				"IPCount":      quotaAndUsage.Usage.FloatingIpsCount,
				"IPLimit":      quotaAndUsage.Limits.FloatingIpCountLimit,
				"BoundIPCount": quotaAndUsage.Usage.BoundFloatingIpsCount}), 33))

	byName := make(map[string][]models.Container)
	for _, c := range containers {
		name := c.Group.Name
		if name == "" {
			name = c.Name
		}
		byName[name] = append(byName[name], c)
	}

	table := cmd.ui.Table([]string{T("Name"), T("Instances"), T("Image"), T("Created"), T("Status")})
	for name, containers := range byName {
		var image string
		parts := strings.Split(containers[0].Image, "/")
		if len(parts) > 0 {
			image = parts[len(parts)-1]
		}

		var createdStr string
		if len(containers) > 1 {
			createdStr = "--"
		} else {
			createdStr = fmt.Sprintf("%d", containers[0].Created)
		}

		status := containers[0].State
		if len(containers) > 1 {
			for i := 1; i < len(containers); i++ {
				if status != containers[i].State {
					status = "??" //TODO
					break
				}
			}
		}

		table.Add(
			name,
			strconv.Itoa(len(containers)),
			image,
			createdStr,
			status)
	}
	table.Print()

	cmd.ui.Say("")
}

func checkTarget(cf plugin.CFContext) error {
	if !cf.HasAPIEndpoint() {
		return fmt.Errorf(T("No CF API endpoint set. Use '{{.Command}}' to target a CloudFoundry environment.",
			map[string]interface{}{"Command": terminal.CommandColor("bx target --cf")}))
	}

	if !cf.IsLoggedIn() {
		return fmt.Errorf(T("Not logged in. Use '{{.Command}}' to log in.",
			map[string]interface{}{"Command": terminal.CommandColor("bx target --cf")}))
	}

	if !cf.HasTargetedSpace() {
		return fmt.Errorf(T("No space targeted. Use '{{.Command}}' to target an org and a space.",
			map[string]interface{}{"Command": terminal.CommandColor("bx target -o ORG -s SPACE")}))
	}

	return nil
}

func formattedGB(sizeInMB int64) string {
	return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%0.2f", float64(sizeInMB)/1024), "0"), ".") + " GB"
}
