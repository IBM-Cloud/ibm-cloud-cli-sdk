package main

import (
	"fmt"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
)

type NamespaceDemo struct{}

func main() {
	plugin.Start(new(NamespaceDemo))
}

func (n *NamespaceDemo) Run(context plugin.PluginContext, args []string) {
	switch args[0] {
	case "list":
		fmt.Println("Running command 'list'.")
	case "show":
		fmt.Println("Running command 'show'.")
	case "delete":
		fmt.Println("Running command 'delete'.")
	}
}

func (n *NamespaceDemo) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "namespace-sample",
		Version: plugin.VersionType{
			Major: 0,
			Minor: 0,
			Build: 1,
		},
		MinCliVersion: plugin.VersionType{
			Major: 0,
			Minor: 0,
			Build: 1,
		},
		Namespaces: []plugin.Namespace{
			plugin.Namespace{
				Name:        "ns",
				Description: "Demonstrate namespace.",
			},
		},
		Commands: []plugin.Command{
			{
				Namespace:   "ns",
				Name:        "list",
				Description: "List resources.",
				Usage:       "ibmcloud ns list",
			},
			{
				Namespace:   "ns",
				Name:        "show",
				Description: "Show details of a resource.",
				Usage:       "ibmcloud ns show",
			},
			{
				Namespace:   "ns",
				Name:        "delete",
				Description: "Delete a resource.",
				Usage:       "ibmcloud ns delete",
			},
		},
	}
}
