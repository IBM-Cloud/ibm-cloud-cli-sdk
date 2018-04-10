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
	table.Add("Username", context.Username())
	table.Add("Org", context.CurrentOrg().Name)
	table.Add("Space", context.CurrentSpace().Name)

	table.Add("CC endpoint", context.APIEndpoint())
	table.Add("UAA endpoint", context.UAAEndpoint())
	table.Add("IAM endpoint", context.IAMTokenEndpoint())
	table.Add("Doppler logging endpoint", context.DopplerEndpoint())

	table.Add("Color enabled", context.ColorEnabled())
	table.Add("HTTP timeout (second)", strconv.Itoa(context.HTTPTimeout()))
	table.Add("Trace", context.Trace())
	table.Add("Locale", context.Locale())

	table.Print()
}

func (p *PrintContext) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "PrintContext",
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
				Description: "Print Bluemix plugin context",
				Usage:       "bx context",
			},
		},
	}
}
