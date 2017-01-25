package api_test

import (
	"net/http"

	"github.com/IBM-Bluemix/bluemix-cli-sdk/common/rest"
	. "github.com/IBM-Bluemix/bluemix-cli-sdk/plugin_examples/list_plugin/api"
	"github.com/IBM-Bluemix/bluemix-cli-sdk/plugin_examples/list_plugin/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/ghttp"
)

var _ = Describe("CCClient", func() {
	var server *Server
	var client CCClient

	AfterEach(func() {
		server.Close()
	})

	BeforeEach(func() {
		server = NewServer()
		client = NewCCClient(server.URL(), rest.NewClient())
	})

	Describe("AppsAndServices()", func() {
		Context("When server returns successfully", func() {
			BeforeEach(func() {
				response := `
                  {  
                     "guid":"space-guid",
                     "name":"space-name",
                     "apps":[  
                        {  
                           "guid":"app1-guid",
                           "urls":[  
                              "app1.mybluemix.net"
                           ],
                           "routes":[  
                              {  
                                 "guid":"route-id",
                                 "host":"app1",
                                 "domain":{  
                                    "guid":"domain-id",
                                    "name":"mybluemix.net"
                                 }
                              }
                           ],
                           "service_count":1,
                           "service_names":[  
                              "serviceB-instance1"
                           ],
                           "running_instances":0,
                           "name":"app1",
                           "memory":256,
                           "instances":2,
                           "disk_quota":1024,
                           "state":"STOPPED"
                        }
                     ],
                     "services":[  
                        {  
                           "guid":"serviceA-instance1-id",
                           "name":"serviceA-instance1",
                           "bound_app_count":0,
                           "last_operation":{  
                              "type":"create",
                              "state":"succeeded",
                              "description":""
                           },
                           "dashboard_url":null,
                           "service_plan":{  
                              "guid":"serviceA-plan-guid",
                              "name":"serviceA-plan",
                              "service":{  
                                 "guid":"serviceA-guid",
                                 "label":"serviceA",
                                 "provider":"serviceA-provider",
                                 "version":"serviceA-version"
                              }
                           }
                        },
                        {  
                           "guid":"serviceB-instance1-id",
                           "name":"serviceB-instance1",
                           "bound_app_count":1,
                           "last_operation":{  
                              "type":"create",
                              "state":"succeeded",
                              "description":""
                           },
                           "dashboard_url":null,
                           "service_plan":{  
                              "guid":"serviceB-plan-guid",
                              "name":"serviceB-plan",
                              "service":{  
                                 "guid":"serviceB-guid",
                                 "label":"serviceB",
                                 "provider":"serviceB-provider",
                                 "version":"serviceB-version"
                              }
                           }
                        }
                     ]
                  }`

				server.AppendHandlers(
					CombineHandlers(
						VerifyRequest(http.MethodGet, "/v2/spaces/space-id/summary"),
						RespondWith(http.StatusOK, response),
					),
				)
			})

			It("should return apps and services summary", func() {
				summary, err := client.AppsAndServices("space-id")

				Expect(err).NotTo(HaveOccurred())
				Expect(len(summary.Apps)).To(Equal(1))
				Expect(summary.Apps[0].Name).To(Equal("app1"))
				Expect(summary.Apps[0].URLs).To(Equal([]string{"app1.mybluemix.net"}))
				Expect(summary.Apps[0].Memory).To(Equal(int64(256)))
				Expect(summary.Apps[0].TotalInstances).To(Equal(2))
				Expect(summary.Apps[0].RunningInstances).To(Equal(0))
				Expect(summary.Apps[0].State).To(Equal("STOPPED"))

				Expect(len(summary.Services)).To(Equal(2))
				Expect(summary.Services[0].Name).To(Equal("serviceA-instance1"))
				Expect(summary.Services[0].ServicePlan.Name).To(Equal("serviceA-plan"))
				Expect(summary.Services[0].ServicePlan.ServiceOffering.Label).To(Equal("serviceA"))
				Expect(summary.Services[1].Name).To(Equal("serviceB-instance1"))
				Expect(summary.Services[1].ServicePlan.Name).To(Equal("serviceB-plan"))
				Expect(summary.Services[1].ServicePlan.ServiceOffering.Label).To(Equal("serviceB"))
			})
		})

		Context("When server return recognized error", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					CombineHandlers(
						VerifyRequest(http.MethodGet, "/v2/spaces/space-id/summary"),
						RespondWith(http.StatusInternalServerError, `
                           {
                              "code":500,
                              "description":"the-error-message"
                           }`,
						),
					),
				)
			})

			It("should return cc api error", func() {
				summary, err := client.AppsAndServices("space-id")

				Expect(summary).To(Equal(models.AppsAndServices{}))
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(CCError{Code: 500, Description: "the-error-message"}))
			})
		})
	})

	Describe("OrgUsage()", func() {
		Context("When server returns the org summary", func() {
			BeforeEach(func() {
				response := `
               {
                  "guid": "org1-id",
                  "name": "org1",
                  "status": "active",
                  "spaces": [
                     {
                        "guid": "space1-id",
                        "name": "space1",
                        "service_count": 1,
                        "app_count": 6,
                        "mem_dev_total": 1152,
                        "mem_prod_total": 0
                     },
                    {
                        "guid": "space2-id",
                        "name": "space2",
                        "service_count": 1,
                        "app_count": 1,
                        "mem_dev_total": 128,
                        "mem_prod_total": 0
                     }
                  ]
               }
            `

				server.AppendHandlers(
					CombineHandlers(
						VerifyRequest(http.MethodGet, "/v2/organizations/org1-id/summary"),
						RespondWith(http.StatusOK, response),
					),
				)
			})

			It("Should returns the usage", func() {
				usage, err := client.OrgUsage("org1-id")

				Expect(err).NotTo(HaveOccurred())
				Expect(usage.Org).To(Equal("org1"))
				Expect(len(usage.Spaces)).To(Equal(2))
				Expect(usage.Spaces[0].Space).To(Equal("space1"))
				Expect(usage.Spaces[0].Apps).To(Equal(6))
				Expect(usage.Spaces[0].Services).To(Equal(1))
				Expect(usage.Spaces[0].MemoryInDev).To(Equal(int64(1152)))
				Expect(usage.Spaces[0].MemoryInProd).To(Equal(int64(0)))
			})
		})
	})
})
