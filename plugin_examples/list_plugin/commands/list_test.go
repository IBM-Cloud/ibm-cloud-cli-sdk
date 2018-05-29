package commands_test

import (
	sdkmodels "github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/models"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin/pluginfakes"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin_examples/list_plugin/api/fakes"
	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin_examples/list_plugin/commands"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin_examples/list_plugin/models"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ListCommand", func() {
	var ui *terminal.FakeUI
	var cf *pluginfakes.FakeCFContext
	var ccClient *fakes.FakeCCClient
	var containerClient *fakes.FakeContainerClient
	var cmd *List
	var err error

	BeforeEach(func() {
		ui = terminal.NewFakeUI()
		cf = new(pluginfakes.FakeCFContext)
		cf.HasAPIEndpointReturns(true)
		cf.IsLoggedInReturns(true)
		cf.HasTargetedSpaceReturns(true)

		context := new(pluginfakes.FakePluginContext)
		context.CFReturns(cf)

		ccClient = new(fakes.FakeCCClient)
		containerClient = new(fakes.FakeContainerClient)
		cmd = NewList(ui, context, ccClient, containerClient)

	})

	JustBeforeEach(func() {
		err = cmd.Run([]string{})
	})

	Context("When API endpoint not set", func() {
		BeforeEach(func() {
			cf.HasAPIEndpointReturns(false)
		})

		It("Should fail", func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("No API endpoint set"))
		})
	})

	Context("When user not logged in", func() {
		BeforeEach(func() {
			cf.IsLoggedInReturns(false)
		})

		It("Should fail", func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Not logged in"))
		})
	})

	Context("When user not target a space", func() {
		BeforeEach(func() {
			cf.HasTargetedSpaceReturns(false)
		})

		It("Should fail", func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("No space targeted"))
		})
	})

	Context("When user is logged in and a space is target", func() {
		BeforeEach(func() {
			cf.CurrentOrganizationReturns(sdkmodels.OrganizationFields{
				QuotaDefinition: sdkmodels.QuotaFields{
					InstanceMemoryLimitInMB: 2048,
					ServicesLimit:           10,
				},
			})

			ccClient.AppsAndServicesReturns(models.AppsAndServices{
				Apps: []models.App{
					{
						Name:             "app1",
						URLs:             []string{"https://app1.example.com"},
						Memory:           int64(512),
						TotalInstances:   2,
						RunningInstances: 1,
						IsDiego:          true,
						State:            "STARTED",
					},
					{
						Name:             "app2",
						URLs:             []string{"https://app2.example.com", "https://app2.another.com"},
						Memory:           int64(256),
						TotalInstances:   1,
						RunningInstances: 0,
						IsDiego:          false,
						State:            "STOPPED",
					},
				},
				Services: []models.ServiceInstance{
					{
						Name: "service1-instance1",
						ServicePlan: models.ServicePlan{
							Name: "plan1",
							ServiceOffering: models.ServiceOffering{
								Label: "service1",
							},
						},
					},
				},
			}, nil)

			ccClient.OrgUsageReturns(models.OrgUsage{
				Org: "org1",
				Spaces: []models.SpaceUsage{
					{
						Space:        "space1",
						Apps:         2,
						Services:     1,
						MemoryInDev:  int64(1028),
						MemoryInProd: int64(512),
					},
					{
						Space:        "space2",
						Apps:         1,
						Services:     1,
						MemoryInDev:  int64(256),
						MemoryInProd: int64(0),
					},
				},
			}, nil)

			containerClient.ContainersReturns([]models.Container{
				{
					Name: "container1",
					Group: models.ContainerGroup{
						Name: "group1",
					},
					Memory:  int64(256),
					Created: 1484718271,
					Image:   "registry/image1",
					State:   "Running",
				},
				{
					Name: "container2",
					Group: models.ContainerGroup{
						Name: "group1",
					},
					Memory:  int64(256),
					Created: 1484718271,
					Image:   "registry/image1",
					State:   "Running",
				},
				{
					Name:    "container3",
					Memory:  int64(512),
					Created: 1484718271,
					Image:   "registry/image2",
					State:   "Shutdown",
				},
			}, nil)

			containerClient.ContainersQuotaAndUsageReturns(models.ContainersQuotaAndUsage{
				Limits: models.ContainersQuota{
					InstancesCountLimit:  10,
					CPUCountLimit:        -1,
					MemoryLimitInMB:      20480,
					FloatingIpCountLimit: 2,
				},
				Usage: models.ContainersUsage{
					TotalInstances:        3,
					RunningInstances:      1,
					CPUCount:              6,
					MemoryInMB:            1024,
					FloatingIpsCount:      2,
					BoundFloatingIpsCount: 1,
				},
			}, nil)
		})

		It("List apps", func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(ui.Outputs()).To(ContainSubstring("CloudFoundy Applications  1.75 GB/2 GB used"))
			Expect(ui.Outputs()).To(ContainSubstring("Name   Routes                     Memory (MB)   Instances   State"))
			Expect(ui.Outputs()).To(ContainSubstring("app1   https://app1.example.com   512           1/2         STARTED"))
			Expect(ui.Outputs()).To(ContainSubstring("app2   https://app2.example.com   256           0/1         STOPPED"))
			Expect(ui.Outputs()).To(ContainSubstring("       https://app2.another.com"))
		})

		It("List services", func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(ui.Outputs()).To(ContainSubstring("Services 2/10 used"))
			Expect(ui.Outputs()).To(ContainSubstring("Name                 Service Offering   Plan"))
			Expect(ui.Outputs()).To(ContainSubstring("service1-instance1   service1           plan1"))
		})

		It("List containers", func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(ui.Outputs()).To(ContainSubstring("Containers  1 GB/20 GB  2/2 Public IPs Requested|1 Used"))
			Expect(ui.Outputs()).To(ContainSubstring("Name         Instances   Image    Created      Status"))
			Expect(ui.Outputs()).To(ContainSubstring("group1       2           image1   --           Running"))
			Expect(ui.Outputs()).To(ContainSubstring("container3   1           image2   1484718271   Shutdown"))
		})
	})
})
