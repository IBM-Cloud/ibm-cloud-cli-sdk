package api_test

import (
	"net/http"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/common/rest"
	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin_examples/list_plugin/api"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/ghttp"
)

var _ = Describe("ContainerClient", func() {
	var server *Server
	var client ContainerClient

	AfterEach(func() {
		server.Close()
	})

	BeforeEach(func() {
		server = NewServer()
		client = NewContainerClient(server.URL(), rest.NewClient())
	})

	Describe("Containers()", func() {
		Context("When server returns successfully", func() {
			BeforeEach(func() {
				response := `
                    [
                        {
                            "ContainerState": "Running",
                            "Created": 1484718255,
                            "Group": {
                                "Id": "group1-id", 
                                "Name": "group1"
                            }, 
                            "Id": "container1-id", 
                            "Name": "container1", 
                            "Image": "registry/image1", 
                            "Memory": 256,
                            "Started": 1484718271
                        }, 
                        {
                            "ContainerState": "Running", 
                            "Created": 1484718254,  
                            "Group": {
                                "Id": "group1-id", 
                                "Name": "group1" 
                            }, 
                            "Id": "container2-id", 
                            "Image": "registry/image1", 
                            "Memory": 256, 
                            "Name": "container2",
                            "Started": 1484718273
                        }, 
                        {
                            "ContainerState": "Shutdown", 
                            "Created": 1484718254,  
                            "Group": {}, 
                            "Id": "container3-id", 
                            "Image": "registry/image2", 
                            "Memory": 512, 
                            "Name": "container3",
                            "Started": 1484718273
                        }
                    ]
                `

				server.AppendHandlers(
					CombineHandlers(
						VerifyRequest(http.MethodGet, "/v3/containers/json", "all=true"),
						VerifyHeaderKV("X-Auth-Project-Id", "space-id"),
						RespondWith(http.StatusOK, response),
					),
				)
			})

			It("should return apps and services summary", func() {
				containers, err := client.Containers("space-id")

				Expect(err).NotTo(HaveOccurred())
				Expect(len(containers)).To(Equal(3))

				Expect(containers[0].Name).To(Equal("container1"))
				Expect(containers[0].Group.Name).To(Equal("group1"))
				Expect(containers[0].Memory).To(Equal(int64(256)))
				Expect(containers[0].Created).To(Equal(int64(1484718255)))
				Expect(containers[0].Image).To(Equal("registry/image1"))
				Expect(containers[0].State).To(Equal("Running"))

				Expect(containers[1].Name).To(Equal("container2"))
				Expect(containers[1].Group.Name).To(Equal("group1"))
				Expect(containers[1].Memory).To(Equal(int64(256)))
				Expect(containers[1].Created).To(Equal(int64(1484718254)))
				Expect(containers[1].Image).To(Equal("registry/image1"))
				Expect(containers[1].State).To(Equal("Running"))

				Expect(containers[2].Name).To(Equal("container3"))
				Expect(containers[2].Group.Name).To(Equal(""))
				Expect(containers[2].Memory).To(Equal(int64(512)))
				Expect(containers[2].Created).To(Equal(int64(1484718254)))
				Expect(containers[2].Image).To(Equal("registry/image2"))
				Expect(containers[2].State).To(Equal("Shutdown"))
			})
		})

		Context("When server return empty json array", func() {
			BeforeEach(func() {
				server.AppendHandlers(RespondWith(http.StatusOK, "[]"))
			})

			It("should return empty containers and no error", func() {
				containers, err := client.Containers("space-id")

				Expect(err).NotTo(HaveOccurred())
				Expect(containers).To(BeEmpty())
			})
		})

		Context("When server return errors", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					RespondWith(http.StatusBadRequest, `
                            {
                                "code": "the-error-code",
                                "description": "the-error-message",
                                "incident_id": "the-incident-id",
                                "rc": "400"
                            }`,
					),
				)
			})

			It("should return error", func() {
				containers, err := client.Containers("space-id")

				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(ContainerError{
					Code:        "the-error-code",
					StatusCode:  "400",
					Description: "the-error-message",
					IncidentId:  "the-incident-id",
				}))
				Expect(containers).To(BeEmpty())
			})
		})
	})

	Describe("ContainersQuotaAndUsage()", func() {
		Context("When server returns containers usage and quota in the given space", func() {
			BeforeEach(func() {
				response := `
                    {
                        "Limits": {
                            "containers": -1, 
                            "floating_ips": 2, 
                            "memory_MB": 2048, 
                            "networks": 5, 
                            "subnets": 5, 
                            "vcpu": -1
                        }, 
                        "Usage": {
                            "containers": 5, 
                            "custom_networks": 0, 
                            "file_share_count": 0, 
                            "floating_ips": 1, 
                            "floating_ips_bound": 1, 
                            "images": 0, 
                            "memory_MB": 1280, 
                            "networks": 1, 
                            "running": 2, 
                            "subnets": 1, 
                            "vcpu": 5
                        }
                    }
                `

				server.AppendHandlers(
					CombineHandlers(
						VerifyRequest(http.MethodGet, "/v3/containers/usage"),
						VerifyHeaderKV("X-Auth-Project-Id", "space-id"),
						RespondWith(http.StatusOK, response),
					),
				)
			})

			It("Should return the quota and usage", func() {
				quotaAndUsage, err := client.ContainersQuotaAndUsage("space-id")

				Expect(err).NotTo(HaveOccurred())

				Expect(quotaAndUsage.Usage.TotalInstances).To(Equal(5))
				Expect(quotaAndUsage.Usage.RunningInstances).To(Equal(2))
				Expect(quotaAndUsage.Usage.FloatingIpsCount).To(Equal(1))
				Expect(quotaAndUsage.Usage.BoundFloatingIpsCount).To(Equal(1))
				Expect(quotaAndUsage.Usage.MemoryInMB).To(Equal(int64(1280)))
				Expect(quotaAndUsage.Usage.CPUCount).To(Equal(5))

				Expect(quotaAndUsage.Limits.CPUCountLimit).To(Equal(-1))
				Expect(quotaAndUsage.Limits.MemoryLimitInMB).To(Equal(int64(2048)))
				Expect(quotaAndUsage.Limits.InstancesCountLimit).To(Equal(-1))
				Expect(quotaAndUsage.Limits.FloatingIpCountLimit).To(Equal(2))

			})
		})
	})
})
