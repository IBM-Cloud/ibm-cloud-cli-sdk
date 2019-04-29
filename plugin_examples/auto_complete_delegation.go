package main

import (
	"fmt"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
)

type AutoCompleteDelegationSample struct{}

func main() {
	plugin.Start(new(AutoCompleteDelegationSample))
}

func (pluginDemo *AutoCompleteDelegationSample) Run(context plugin.PluginContext, args []string) {
	switch args[0] {
	case "SendCompletion":
		// CLI invokes 'PATH_TO_PLUGIN_BINARY SendCompletion [NAMESPACE] ... [COMMAND] [ARGS]...'
		// for example: 'AutoCompleteDelegationSample SendCompletion autocomplete-sample set-role'

		subCommands := []string{"get-role", "set-role", "help"}
		var cmd string
		if len(args) > 2 {
			for _, c := range subCommands {
				if c == args[2] {
					cmd = c
					break
				}
			}
		}

		if cmd == "" {
			// completion for sub-commands
			fmt.Println(strings.Join(subCommands, "\n"))
			return
		}

		if cmd == "set-role" {
			// completion for command args
			fmt.Println("Viewer\nEditor\nOperator\nAdministrator")
		}
	default:
		fmt.Printf("Running command %q\n", args[0])
	}
}

func (pluginDemo *AutoCompleteDelegationSample) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "auto-complete-delegation-sample",
		Version: plugin.VersionType{
			Major: 0,
			Minor: 0,
			Build: 1,
		},
		DelegateBashCompletion: true, // enable command completion delegation
		Namespaces: []plugin.Namespace{
			{
				Name:        "autocomplete-sample",
				Description: "Demonstrate delegate command completion to plugin.",
			},
		},
		Commands: []plugin.Command{
			{
				Namespace:   "autocomplete-sample",
				Name:        "get-role",
				Description: "get user's role",
				Usage:       "ibmcloud autocomplete-sample get-role",
			},
			{
				Namespace:   "autocomplete-sample",
				Name:        "set-role",
				Description: "set a user role (Viewer, Editor, Operator or Administrator)",
				Usage:       "ibmcloud autocomplete-sample set-role",
			},
			{
				Namespace:   "autocomplete-sample",
				Name:        "help",
				Description: "show help",
				Usage:       "ibmcloud autocomplete-sample help",
			},
		},
	}
}
