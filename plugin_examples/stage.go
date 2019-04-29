package main

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"

	"fmt"
)

type StageDemo struct{}

func main() {
	plugin.Start(new(StageDemo))
}

func (n *StageDemo) Run(context plugin.PluginContext, args []string) {
	switch args[0] {
	case "list":
		fmt.Println("Running command 'list'.")
	case "show":
		fmt.Println("Running command 'show'.")
	}
}

func (n *StageDemo) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "stage-demo",
		Version: plugin.VersionType{
			Major: 0,
			Minor: 0,
			Build: 1,
		},
		Namespaces: []plugin.Namespace{
			plugin.Namespace{
				Name:        "stage",
				Description: "Show example of stage annotation",
				Stage:       plugin.StageBeta,
			},
		},
		Commands: []plugin.Command{
			{
				Namespace:   "stage",
				Name:        "list",
				Description: "List resources",
				Usage:       "ibmcloud stage list",
				Stage:       plugin.StageDeprecated,
			},
			{
				Namespace:   "stage",
				Name:        "show",
				Description: "Show the details of a resource",
				Usage:       "ibmcloud stage show",
				Stage:       plugin.StageExperimental,
			},
		},
	}
}
