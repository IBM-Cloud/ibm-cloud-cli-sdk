package main

import (
	"fmt"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
)

type HelloWorldPlugin struct{}

func main() {
	plugin.Start(new(HelloWorldPlugin))
}

func (pluginDemo *HelloWorldPlugin) Run(context plugin.PluginContext, args []string) {
	fmt.Println("Hi, this is my first plugin for Bluemix")
}

func (pluginDemo *HelloWorldPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "SayHello",
		Version: plugin.VersionType{
			Major: 0,
			Minor: 0,
			Build: 1,
		},
		Commands: []plugin.Command{
			{
				Name:        "hello",
				Alias:       "hi",
				Description: "say hello to Bluemix.",
				Usage:       "bluemix hello",
			},
		},
	}
}
