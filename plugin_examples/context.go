package main

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
)

type PrintContext struct {
	ui terminal.UI
}

func main() {
	plugin.Start(NewPrintContext())
}

func NewPrintContext() *PrintContext {
	return &PrintContext{
		ui: terminal.NewStdUI(),
	}
}

func (p *PrintContext) Run(context plugin.PluginContext, args []string) {
	table := p.ui.Table([]string{"Name", "Value"})
	table.Add("API endpoint", context.APIEndpoint())
	table.Add("IAM endpoint", context.IAMEndpoint())
	table.Add("Username", context.UserEmail())

	cf := context.CF()
	table.Add("CC endpoint", cf.APIEndpoint())
	table.Add("UAA endpoint", cf.UAAEndpoint())
	table.Add("Doppler logging endpoint", cf.DopplerEndpoint())
	table.Add("Org", cf.CurrentOrganization().Name)
	table.Add("Space", cf.CurrentSpace().Name)

	table.Add("Color enabled", context.ColorEnabled())
	table.Add("HTTP timeout (second)", strconv.Itoa(context.HTTPTimeout()))
	table.Add("Trace", context.Trace())
	table.Add("Locale", context.Locale())

	table.Print()
}

func (p *PrintContext) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "context-sample",
		Version: plugin.VersionType{
			Major: 0,
			Minor: 1,
			Build: 0,
		},
		MinCliVersion: plugin.VersionType{
			Major: 0,
			Minor: 0,
			Build: 1,
		},
		Commands: []plugin.Command{
			{
				Name:        "context",
				Description: "Print IBM Cloud plugin context",
				Usage:       "ibmcloud context",
			},
		},
	}
}
