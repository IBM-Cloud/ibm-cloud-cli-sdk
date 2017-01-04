package main

import (
	"fmt"

	"github.com/IBM-Bluemix/bluemix-cli-sdk/plugin"
)

type NamespaceDemo struct{}

func main() {
	plugin.Start(new(NamespaceDemo))
}

func (n *NamespaceDemo) Run(context plugin.PluginContext, args []string) {
	switch args[0] {
	case "cmd-1":
		fmt.Println("Running command 'cmd-1'.")
	case "cmd-2":
		fmt.Println("Running command 'cmd-2'.")
	case "cmd-3":
		fmt.Println("Running command 'cmd-3'.")
	}
}

func (n *NamespaceDemo) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "NamespaceDemo",
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
				Name:        "demo",
				Description: "Plugin demonstration",
			},
		},
		Commands: []plugin.Command{
			{
				Namespace:   "demo",
				Name:        "cmd-1",
				Description: "Help text of cmd-1",
				Usage:       "bx demo cmd-1",
			},
			{
				Namespace:   "demo",
				Name:        "cmd-2",
				Description: "Help text of cmd-2",
				Usage:       "bx demo cmd-2",
			},
			{
				Namespace:   "demo",
				Name:        "cmd-3",
				Description: "Help text of cmd-3",
				Usage:       "bx demo cmd-3",
			},
		},
	}
}
