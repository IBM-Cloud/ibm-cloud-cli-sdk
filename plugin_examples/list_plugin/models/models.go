package models

type AppsAndServices struct {
	Apps     []App
	Services []ServiceInstance
}

type App struct {
	Name             string
	URLs             []string `json:"urls"`
	Memory           int64    `json:"memory"`
	TotalInstances   int      `json:"instances"`
	RunningInstances int      `json:"running_instances"`
	IsDiego          bool     `json:"diego"`
	State            string
}

type ServiceInstance struct {
	Name        string
	ServicePlan ServicePlan `json:"service_plan"`
}

type ServicePlan struct {
	Name            string
	ServiceOffering ServiceOffering `json:"service"`
}

type ServiceOffering struct {
	Label string
}

type OrgUsage struct {
	Org    string `json:"name"`
	Spaces []SpaceUsage
}

func (o OrgUsage) TotalMemoryUsed() int64 {
	var sum int64
	for _, s := range o.Spaces {
		sum += s.MemoryInDev + s.MemoryInProd
	}
	return sum
}

func (o OrgUsage) AppsCount() int {
	var sum int
	for _, s := range o.Spaces {
		sum += s.Apps
	}
	return sum
}

func (o OrgUsage) ServicesCount() int {
	var sum int
	for _, s := range o.Spaces {
		sum += s.Services
	}
	return sum
}

type SpaceUsage struct {
	Space        string `json:"name"`
	Apps         int    `json:"app_count"`
	Services     int    `json:"service_count"`
	MemoryInDev  int64  `json:"mem_dev_total"`
	MemoryInProd int64  `json:"mem_prod_total"`
}

type Container struct {
	Name    string
	Group   ContainerGroup
	Memory  int64
	Created int64
	Image   string
	State   string `json:"ContainerState"`
}

type ContainerGroup struct {
	Name string
}

type ContainersQuotaAndUsage struct {
	Limits ContainersQuota `json:"Limits"`
	Usage  ContainersUsage `json:"Usage"`
}

type ContainersQuota struct {
	InstancesCountLimit  int   `json:"containers"`
	CPUCountLimit        int   `json:"vcpu"`
	MemoryLimitInMB      int64 `json:"memory_MB"`
	FloatingIpCountLimit int   `json:"floating_ips"`
}

type ContainersUsage struct {
	TotalInstances        int   `json:"containers"`
	RunningInstances      int   `json:"running"`
	CPUCount              int   `json:"vcpu"`
	MemoryInMB            int64 `json:"memory_MB"`
	FloatingIpsCount      int   `json:"floating_ips"`
	BoundFloatingIpsCount int   `json:"floating_ips_bound"`
}
