package plugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnclosedGroupings(t *testing.T) {

	testData := []struct {
		usageText         string
		indicies          []int
		invalid           bool
		actualOpeningChar rune
	}{
		{
			usageText:         "command [ABC [DEF]",
			invalid:           true,
			actualOpeningChar: '[',
			indicies:          []int{8, 17},
		},
		{
			usageText: "command ([ABC] | [DEF])",
			invalid:   false,
		},
		{
			usageText:         "command --data [{\"a\": 1, \"b\": { \"d\": \"id\"}, \"e\": 2 }, {\"a\": 2, \"b\": { \"d\": \"id\", \"e\": 3 }]'",
			invalid:           true,
			actualOpeningChar: '{',
			indicies:          []int{54, 90},
		},
	}

	for _, d := range testData {
		isInvalid, char, indicies := UnclosedGroupings(d.usageText)

		assert.Equal(t, d.invalid, isInvalid, "expected usage text to be invalid")
		if isInvalid {
			assert.Equal(t, d.actualOpeningChar, char)
			assert.Equal(t, d.indicies[0], indicies[0])
			assert.Equal(t, d.indicies[1], indicies[1])
		}
	}
}

func TestErrors(t *testing.T) {
	validator := NewPluginMetadataValidator()

	testData := []struct {
		name           string
		pluginMetadata []PluginMetadata
		errors         PluginToValidationErrors
	}{
		{
			name: "No errors",
			pluginMetadata: []PluginMetadata{
				{
					Name: "plugin",
					Version: VersionType{
						Major: 1,
						Minor: 0,
						Build: 0,
					},
					MinCliVersion: VersionType{
						Major: 2,
						Minor: 0,
						Build: 0,
					},
					Namespaces: []Namespace{
						{
							ParentName: "namespace",
							Name:       "command",
						},
					},
					Commands: []Command{
						{
							Namespace:   "plugin",
							Name:        "list",
							Description: "List all resources",
							Usage:       "ibmcloud plugin list [--output FORMAT]",
							Flags: []Flag{
								{
									Name:        "output",
									Description: "Specify output format",
								},
							},
						},
					},
				},
			},
			errors: PluginToValidationErrors{},
		},
		{
			name: "No errors - with full path to ibmcloud",
			pluginMetadata: []PluginMetadata{
				{
					Name: "plugin",
					Version: VersionType{
						Major: 1,
						Minor: 0,
						Build: 0,
					},
					MinCliVersion: VersionType{
						Major: 2,
						Minor: 0,
						Build: 0,
					},
					Namespaces: []Namespace{
						{
							ParentName: "namespace",
							Name:       "command",
						},
					},
					Commands: []Command{
						{
							Namespace:   "plugin",
							Name:        "list",
							Description: "List all resources",
							Usage:       "/usr/local/bin/ibmcloud plugin list",
							Flags:       []Flag{},
						},
					},
				},
			},
			errors: PluginToValidationErrors{},
		},
		{
			name: "No errors - with uppercase arguments",
			pluginMetadata: []PluginMetadata{
				{
					Name: "plugin",
					Version: VersionType{
						Major: 1,
						Minor: 0,
						Build: 0,
					},
					MinCliVersion: VersionType{
						Major: 2,
						Minor: 0,
						Build: 0,
					},
					Namespaces: []Namespace{
						{
							ParentName: "namespace",
							Name:       "command",
						},
					},
					Commands: []Command{
						{
							Namespace:   "resource",
							Name:        "service-instance-create",
							Description: "Create a service instance",
							Usage:       "ibmcloud resource service-instance-create NAME SERVICE_NAME PLAN_NAME LOCATION",
							Flags:       []Flag{},
						},
					},
				},
			},
			errors: PluginToValidationErrors{},
		},
		{
			name: "No errors - with multiple uppercase arguments",
			pluginMetadata: []PluginMetadata{
				{
					Name: "ks",
					Version: VersionType{
						Major: 1,
						Minor: 0,
						Build: 0,
					},
					MinCliVersion: VersionType{
						Major: 2,
						Minor: 0,
						Build: 0,
					},
					Namespaces: []Namespace{
						{
							ParentName: "",
							Name:       "ks",
						},
					},
					Commands: []Command{
						{
							Namespace:   "ks cluster create",
							Name:        "classic",
							Description: "Create a cluster with worker nodes on classic infrastructure.",
							Usage:       "ibmcloud ks cluster create classic --name NAME --zone ZONE [--workers COUNT]",
							Flags:       []Flag{},
						},
					},
				},
			},
			errors: PluginToValidationErrors{},
		},
		{
			name: "No errors - with proper positional arguments",
			pluginMetadata: []PluginMetadata{
				{
					Name: "plugin",
					Version: VersionType{
						Major: 1,
						Minor: 0,
						Build: 0,
					},
					MinCliVersion: VersionType{
						Major: 2,
						Minor: 0,
						Build: 0,
					},
					Namespaces: []Namespace{
						{
							ParentName: "namespace",
							Name:       "command",
						},
					},
					Commands: []Command{
						{
							Namespace:   "account",
							Name:        "org-role-set",
							Description: "List all resources",
							Usage:       "ibmcloud account org-role-set USER_NAME ORG_NAME (ROLE_A | ROLE_B)",
							Flags:       []Flag{},
						},
					},
				},
			},
			errors: PluginToValidationErrors{},
		},
		{
			name: "No errors - with arguments and flags",
			pluginMetadata: []PluginMetadata{
				{
					Name: "plugin",
					Version: VersionType{
						Major: 1,
						Minor: 0,
						Build: 0,
					},
					MinCliVersion: VersionType{
						Major: 2,
						Minor: 0,
						Build: 0,
					},
					Namespaces: []Namespace{
						{
							ParentName: "namespace",
							Name:       "command",
						},
					},
					Commands: []Command{
						{
							Namespace:   "service",
							Name:        "create",
							Description: "Createe a service",
							Usage:       "ibmcloud service create NAME PLAN LOCATION [--tags TAGS] [--parameters PARAMS]",
							Flags: []Flag{
								{
									Name:        "tags",
									Description: "Specify tags",
								},
								{
									Name:        "parameters",
									Description: "Specify parameters",
								},
							},
						},
					},
				},
			},
			errors: PluginToValidationErrors{},
		},
		{
			name: "No plugin name",
			pluginMetadata: []PluginMetadata{
				{
					Name: "",
					Version: VersionType{
						Major: 1,
						Minor: 0,
						Build: 0,
					},
					MinCliVersion: VersionType{
						Major: 2,
						Minor: 0,
						Build: 0,
					},
					Namespaces: []Namespace{
						{
							ParentName: "namespace",
							Name:       "command",
						},
					},
					Commands: []Command{
						{
							Namespace:   "plugin",
							Name:        "command",
							Description: "A sample command",
							Usage:       "ibmcloud plugin command ABC (--DEF | --GHI)",
							Flags: []Flag{
								{
									Name:        "DEF",
									Description: "Description for flag DEF",
								},
								{
									Name:        "GHI",
									Description: "Description for flag GHI",
								},
							},
						},
					},
				},
			},
			errors: PluginToValidationErrors{
				"UNKNOWN": []PluginMetadataError{
					{
						Namespace: "PluginMetadata.Name",
						Error:     "Name is required",
						Priority:  PriorityError,
					},
				},
			},
		},
		{
			name: "No commands",
			pluginMetadata: []PluginMetadata{
				{
					Name: "plugin",
					Version: VersionType{
						Major: 1,
						Minor: 0,
						Build: 0,
					},
					MinCliVersion: VersionType{
						Major: 2,
						Minor: 0,
						Build: 0,
					},
					Namespaces: []Namespace{
						{
							ParentName: "namespace",
							Name:       "command",
						},
					},
					Commands: []Command{},
				},
			},
			errors: PluginToValidationErrors{
				"plugin": []PluginMetadataError{
					{
						Namespace: "PluginMetadata.Commands",
						Error:     "Commands must contain at least 1 element",
						Priority:  PriorityError,
					},
				},
			},
		},
		{
			name: "No namespaces",
			pluginMetadata: []PluginMetadata{
				{
					Name: "plugin",
					Version: VersionType{
						Major: 1,
						Minor: 0,
						Build: 0,
					},
					MinCliVersion: VersionType{
						Major: 2,
						Minor: 0,
						Build: 0,
					},
					Namespaces: []Namespace{},
					Commands: []Command{
						{
							Namespace:   "plugin",
							Name:        "command",
							Description: "A sample command",
							Usage:       "ibmcloud plugin command ABC (--DEF | --GHI)",
							Flags: []Flag{
								{
									Name:        "DEF",
									Description: "Description for flag DEF",
								},
								{
									Name:        "GHI",
									Description: "Description for flag GHI",
								},
							},
						},
					},
				},
			},
			errors: PluginToValidationErrors{
				"plugin": []PluginMetadataError{
					{
						Namespace: "PluginMetadata.Namespaces",
						Error:     "Namespaces must contain at least 1 element",
						Priority:  PriorityError,
					},
				},
			},
		},
		{
			name: "No name on namespace",
			pluginMetadata: []PluginMetadata{
				{
					Name: "plugin",
					Version: VersionType{
						Major: 1,
						Minor: 0,
						Build: 0,
					},
					MinCliVersion: VersionType{
						Major: 2,
						Minor: 0,
						Build: 0,
					},
					Namespaces: []Namespace{
						{
							ParentName: "plugin",
						},
					},
					Commands: []Command{
						{
							Namespace:   "plugin",
							Name:        "command",
							Description: "A sample command",
							Usage:       "ibmcloud plugin command ABC (--DEF | --GHI)",
							Flags: []Flag{
								{
									Name:        "DEF",
									Description: "Description for flag DEF",
								},
								{
									Name:        "GHI",
									Description: "Description for flag GHI",
								},
							},
						},
					},
				},
			},
			errors: PluginToValidationErrors{
				"plugin": []PluginMetadataError{
					{
						Namespace: "PluginMetadata.Namespaces[0].Name",
						Error:     "Name is required",
						Priority:  PriorityError,
					},
				},
			},
		},
		{
			name: "Command is missing a namespace",
			pluginMetadata: []PluginMetadata{
				{
					Name: "plugin",
					Version: VersionType{
						Major: 1,
						Minor: 0,
						Build: 0,
					},
					MinCliVersion: VersionType{
						Major: 2,
						Minor: 0,
						Build: 0,
					},
					Namespaces: []Namespace{
						{
							ParentName: "plugin",
							Name:       "n1",
						},
					},
					Commands: []Command{
						{
							Name:        "command",
							Description: "A sample command",
							Usage:       "ibmcloud command ABC (--DEF | --GHI)",
							Flags: []Flag{
								{
									Name:        "DEF",
									Description: "Description for flag DEF",
								},
								{
									Name:        "GHI",
									Description: "Description for flag GHI",
								},
							},
						},
					},
				},
			},
			errors: PluginToValidationErrors{
				"plugin": []PluginMetadataError{
					{
						CommandName: " command",
						Namespace:   "PluginMetadata.Commands[0].Namespace",
						Error:       "Namespace is required",
						Priority:    PriorityError,
					},
				},
			},
		},
		{
			name: "Command is missing a name",
			pluginMetadata: []PluginMetadata{
				{
					Name: "plugin",
					Version: VersionType{
						Major: 1,
						Minor: 0,
						Build: 0,
					},
					MinCliVersion: VersionType{
						Major: 2,
						Minor: 0,
						Build: 0,
					},
					Namespaces: []Namespace{
						{
							ParentName: "plugin",
							Name:       "n1",
						},
					},
					Commands: []Command{
						{
							Namespace:   "n1",
							Description: "A sample command",
							Usage:       "ibmcloud n1 command ABC (--DEF | --GHI)",
							Flags: []Flag{
								{
									Name:        "DEF",
									Description: "Description for flag DEF",
								},
								{
									Name:        "GHI",
									Description: "Description for flag GHI",
								},
							},
						},
					},
				},
			},
			errors: PluginToValidationErrors{
				"plugin": []PluginMetadataError{
					{
						CommandName: "n1 ",
						Namespace:   "PluginMetadata.Commands[0].Name",
						Error:       "Name is required",
						Priority:    PriorityError,
					},
				},
			},
		},
		{
			name: "Command is missing a description",
			pluginMetadata: []PluginMetadata{
				{
					Name: "plugin",
					Version: VersionType{
						Major: 1,
						Minor: 0,
						Build: 0,
					},
					MinCliVersion: VersionType{
						Major: 2,
						Minor: 0,
						Build: 0,
					},
					Namespaces: []Namespace{
						{
							ParentName: "plugin",
							Name:       "n1",
						},
					},
					Commands: []Command{
						{
							Namespace: "n1",
							Name:      "command",
							Usage:     "ibmcloud n1 command ABC (--DEF | --GHI)",
							Flags: []Flag{
								{
									Name:        "DEF",
									Description: "Description for flag DEF",
								},
								{
									Name:        "GHI",
									Description: "Description for flag GHI",
								},
							},
						},
					},
				},
			},
			errors: PluginToValidationErrors{
				"plugin": []PluginMetadataError{
					{
						CommandName: "n1 command",
						Namespace:   "PluginMetadata.Commands[0].Description",
						Error:       "Command 'command' has no description. All commands must have a clear description.",
						Priority:    PriorityError,
						Remediation: "Add a sentence without subject describing what the command does.",
					},
				},
			},
		},
		{
			name: "Command is missing usage text",
			pluginMetadata: []PluginMetadata{
				{
					Name: "plugin",
					Version: VersionType{
						Major: 1,
						Minor: 0,
						Build: 0,
					},
					MinCliVersion: VersionType{
						Major: 2,
						Minor: 0,
						Build: 0,
					},
					Namespaces: []Namespace{
						{
							ParentName: "plugin",
							Name:       "n1",
						},
					},
					Commands: []Command{
						{
							Namespace:   "n1",
							Name:        "command",
							Description: "This is a command",
							Flags: []Flag{
								{
									Name:        "DEF",
									Description: "Description for flag DEF",
								},
								{
									Name:        "GHI",
									Description: "Description for flag GHI",
								},
							},
						},
					},
				},
			},
			errors: PluginToValidationErrors{
				"plugin": []PluginMetadataError{
					{
						CommandName: "n1 command",
						Namespace:   "PluginMetadata.Commands[0].Description",
						Error:       "Description for 'command' starts with 'this is'. Use a sentence without subject.",
						Priority:    PriorityError,
						Remediation: "Remove 'this is' and start directly with the action. Example: 'List all instances' instead of 'This command lists all instances'.",
					},
					{
						CommandName: "n1 command",
						Namespace:   "PluginMetadata.Commands[0].Usage",
						Error:       "Command 'command' has no usage information.",
						Priority:    PriorityError,
						Remediation: "Add usage text showing command syntax with parameters and options.",
					},
				},
			},
		},
		{
			name: "Command is missing usage text",
			pluginMetadata: []PluginMetadata{
				{
					Name: "plugin",
					Version: VersionType{
						Major: 1,
						Minor: 0,
						Build: 0,
					},
					MinCliVersion: VersionType{
						Major: 2,
						Minor: 0,
						Build: 0,
					},
					Namespaces: []Namespace{
						{
							ParentName: "plugin",
							Name:       "n1",
						},
					},
					Commands: []Command{
						{
							Namespace:   "n1",
							Name:        "command",
							Description: "This is a command",
							Flags: []Flag{
								{
									Name:        "DEF",
									Description: "Description for flag DEF",
								},
								{
									Name:        "GHI",
									Description: "Description for flag GHI",
								},
							},
						},
					},
				},
			},
			errors: PluginToValidationErrors{
				"plugin": []PluginMetadataError{
					{
						CommandName: "n1 command",
						Namespace:   "PluginMetadata.Commands[0].Description",
						Error:       "Description for 'command' starts with 'this is'. Use a sentence without subject.",
						Priority:    PriorityError,
						Remediation: "Remove 'this is' and start directly with the action. Example: 'List all instances' instead of 'This command lists all instances'.",
					},
					{
						CommandName: "n1 command",
						Namespace:   "PluginMetadata.Commands[0].Usage",
						Error:       "Command 'command' has no usage information.",
						Priority:    PriorityError,
						Remediation: "Add usage text showing command syntax with parameters and options.",
					},
				},
			},
		},
		{
			name: "Command's usage text contains placeholders",
			pluginMetadata: []PluginMetadata{
				{
					Name: "plugin",
					Version: VersionType{
						Major: 1,
						Minor: 0,
						Build: 0,
					},
					MinCliVersion: VersionType{
						Major: 2,
						Minor: 0,
						Build: 0,
					},
					Namespaces: []Namespace{
						{
							ParentName: "plugin",
							Name:       "n1",
						},
					},
					Commands: []Command{
						{
							Namespace:   "n1",
							Name:        "command",
							Usage:       "ibmcloud n1 command [options...]",
							Description: "This is a command",
							Flags: []Flag{
								{
									Name:        "DEF",
									Description: "Description for flag DEF",
								},
								{
									Name:        "GHI",
									Description: "Description for flag GHI",
								},
							},
						},
					},
				},
			},
			errors: PluginToValidationErrors{
				"plugin": []PluginMetadataError{
					{
						CommandName: "n1 command",
						Namespace:   "PluginMetadata.Commands[0].Usage",
						Error:       "Usage contains placeholder arguments/flags",
						Priority:    PriorityError,
					},
					{
						CommandName: "n1 command",
						Namespace:   "PluginMetadata.Commands[0].Description",
						Error:       "Description for 'command' starts with 'this is'. Use a sentence without subject.",
						Priority:    PriorityError,
						Remediation: "Remove 'this is' and start directly with the action. Example: 'List all instances' instead of 'This command lists all instances'.",
					},
				},
			},
		},
		{
			name: "Command's usage just COMMAND",
			pluginMetadata: []PluginMetadata{
				{
					Name: "plugin",
					Version: VersionType{
						Major: 1,
						Minor: 0,
						Build: 0,
					},
					MinCliVersion: VersionType{
						Major: 2,
						Minor: 0,
						Build: 0,
					},
					Namespaces: []Namespace{
						{
							ParentName: "plugin",
							Name:       "n1",
						},
					},
					Commands: []Command{
						{
							Namespace:   "n1",
							Name:        "command",
							Usage:       "ibmcloud n1 command COMMAND",
							Description: "This is a command",
							Flags: []Flag{
								{
									Name:        "DEF",
									Description: "Description for flag DEF",
								},
								{
									Name:        "GHI",
									Description: "Description for flag GHI",
								},
							},
						},
					},
				},
			},
			errors: PluginToValidationErrors{
				"plugin": []PluginMetadataError{
					{
						CommandName: "n1 command",
						Namespace:   "PluginMetadata.Commands[0].Usage",
						Error:       "Usage does not have any usage text besides COMMAND",
						Priority:    PriorityError,
					},
					{
						CommandName: "n1 command",
						Namespace:   "PluginMetadata.Commands[0].Description",
						Error:       "Description for 'command' starts with 'this is'. Use a sentence without subject.",
						Priority:    PriorityError,
						Remediation: "Remove 'this is' and start directly with the action. Example: 'List all instances' instead of 'This command lists all instances'.",
					},
				},
			},
		},
		{
			name: "Missing prefix in usage text",
			pluginMetadata: []PluginMetadata{
				{
					Name: "plugin",
					Version: VersionType{
						Major: 1,
						Minor: 0,
						Build: 0,
					},
					MinCliVersion: VersionType{
						Major: 2,
						Minor: 0,
						Build: 0,
					},
					Namespaces: []Namespace{
						{
							ParentName: "plugin",
							Name:       "n1",
						},
					},
					Commands: []Command{
						{
							Namespace:   "n1",
							Name:        "command",
							Usage:       "n1 command COMMAND",
							Description: "This is a command",
							Flags: []Flag{
								{
									Name:        "DEF",
									Description: "Description for flag DEF",
								},
								{
									Name:        "GHI",
									Description: "Description for flag GHI",
								},
							},
						},
					},
				},
			},
			errors: PluginToValidationErrors{
				"plugin": []PluginMetadataError{
					{
						CommandName: "n1 command",
						Namespace:   "PluginMetadata.Commands[0].Usage",
						Error:       "Usage should start with 'ibmcloud' (lowercase) or the full path to the ibmcloud binary.",
						Priority:    PriorityError,
						Remediation: "Start usage examples with 'ibmcloud' in lowercase (e.g., 'ibmcloud plugin-name command...').",
					},
					{
						CommandName: "n1 command",
						Namespace:   "PluginMetadata.Commands[0].Description",
						Error:       "Description for 'command' starts with 'this is'. Use a sentence without subject.",
						Priority:    PriorityError,
						Remediation: "Remove 'this is' and start directly with the action. Example: 'List all instances' instead of 'This command lists all instances'.",
					},
				},
			},
		},
		{
			name: "Wrong prefix - uppercase",
			pluginMetadata: []PluginMetadata{
				{
					Name: "plugin",
					Version: VersionType{
						Major: 1,
						Minor: 0,
						Build: 0,
					},
					MinCliVersion: VersionType{
						Major: 2,
						Minor: 0,
						Build: 0,
					},
					Namespaces: []Namespace{
						{
							ParentName: "",
							Name:       "ibmcloud",
						},
					},
					Commands: []Command{
						{
							Namespace:   "plugin",
							Name:        "list",
							Usage:       "IBMCLOUD plugin list",
							Description: "List all installed plug-ins",
							Flags:       []Flag{},
						},
					},
				},
			},
			errors: PluginToValidationErrors{
				"plugin": []PluginMetadataError{
					{
						CommandName: "plugin list",
						Namespace:   "PluginMetadata.Commands[0].Usage",
						Error:       "Usage should start with 'ibmcloud' (lowercase) or the full path to the ibmcloud binary.",
						Priority:    PriorityError,
						Remediation: "Start usage examples with 'ibmcloud' in lowercase (e.g., 'ibmcloud plugin-name command...').",
					},
				},
			},
		},
		{
			name: "Wrong prefix - mixed case",
			pluginMetadata: []PluginMetadata{
				{
					Name: "plugin",
					Version: VersionType{
						Major: 1,
						Minor: 0,
						Build: 0,
					},
					MinCliVersion: VersionType{
						Major: 2,
						Minor: 0,
						Build: 0,
					},
					Namespaces: []Namespace{
						{
							ParentName: "",
							Name:       "ibmcloud",
						},
					},
					Commands: []Command{
						{

							Namespace:   "plugin",
							Name:        "list",
							Usage:       "IBMCloud plugin list",
							Description: "List all installed plug-ins",
							Flags:       []Flag{},
						},
					},
				},
			},
			errors: PluginToValidationErrors{
				"plugin": []PluginMetadataError{
					{
						CommandName: "plugin list",
						Namespace:   "PluginMetadata.Commands[0].Usage",
						Error:       "Usage should start with 'ibmcloud' (lowercase) or the full path to the ibmcloud binary.",
						Priority:    PriorityError,
						Remediation: "Start usage examples with 'ibmcloud' in lowercase (e.g., 'ibmcloud plugin-name command').",
					},
				},
			},
		},
		{
			name: "Wrong prefix - ic only",
			pluginMetadata: []PluginMetadata{
				{
					Name: "plugin",
					Version: VersionType{
						Major: 1,
						Minor: 0,
						Build: 0,
					},
					MinCliVersion: VersionType{
						Major: 2,
						Minor: 0,
						Build: 0,
					},
					Namespaces: []Namespace{
						{
							ParentName: "plugin",
							Name:       "n1",
						},
					},
					Commands: []Command{
						{

							Namespace:   "plugin",
							Name:        "list",
							Usage:       "ic plugin list",
							Description: "List all installed plug-ins",
							Flags:       []Flag{},
						},
					},
				},
			},
			errors: PluginToValidationErrors{
				"plugin": []PluginMetadataError{
					{
						CommandName: "plugin list",
						Namespace:   "PluginMetadata.Commands[0].Usage",
						Error:       "Usage should start with 'ibmcloud' (lowercase) or the full path to the ibmcloud binary.",
						Priority:    PriorityError,
						Remediation: "Start usage examples with 'ibmcloud' in lowercase (e.g., 'ibmcloud plugin-name command').",
					},
				},
			},
		},
		{
			name: "Command name only",
			pluginMetadata: []PluginMetadata{
				{
					Name: "plugin",
					Version: VersionType{
						Major: 1,
						Minor: 0,
						Build: 0,
					},
					MinCliVersion: VersionType{
						Major: 2,
						Minor: 0,
						Build: 0,
					},
					Namespaces: []Namespace{
						{
							ParentName: "plugin",
							Name:       "n1",
						},
					},
					Commands: []Command{
						{

							Namespace:   "plugin",
							Name:        "list",
							Usage:       "list",
							Description: "List all installed plug-ins",
							Flags:       []Flag{},
						},
					},
				},
			},
			errors: PluginToValidationErrors{
				"plugin": []PluginMetadataError{
					{
						CommandName: "plugin list",
						Namespace:   "PluginMetadata.Commands[0].Usage",
						Error:       "Usage should start with 'ibmcloud' (lowercase) or the full path to the ibmcloud binary.",
						Priority:    PriorityError,
						Remediation: "Start usage examples with 'ibmcloud' in lowercase (e.g., 'ibmcloud plugin-name command').",
					},
				},
			},
		},
		{
			name: "Command's usage contains unclosed brackets",
			pluginMetadata: []PluginMetadata{
				{
					Name: "plugin",
					Version: VersionType{
						Major: 1,
						Minor: 0,
						Build: 0,
					},
					MinCliVersion: VersionType{
						Major: 2,
						Minor: 0,
						Build: 0,
					},
					Namespaces: []Namespace{
						{
							ParentName: "plugin",
							Name:       "n1",
						},
					},
					Commands: []Command{
						{
							Namespace:   "n1",
							Name:        "command",
							Usage:       "ibmcloud n1 command ([--DEF | ABD)",
							Description: "List all the things",
							Flags: []Flag{
								{
									Name:        "DEF",
									Description: "Description for flag DEF",
								},
							},
						},
					},
				},
			},
			errors: PluginToValidationErrors{
				"plugin": []PluginMetadataError{
					{
						CommandName: "n1 command",
						Namespace:   "PluginMetadata.Commands[0].Usage",
						Error:       "Usage contains unclosed [ between indicies [21 34]",
						Priority:    PriorityError,
					},
				},
			},
		},
		{
			name: "Command flag name contains UTF-16 characters",
			pluginMetadata: []PluginMetadata{
				{
					Name: "plugin",
					Version: VersionType{
						Major: 1,
						Minor: 0,
						Build: 0,
					},
					MinCliVersion: VersionType{
						Major: 2,
						Minor: 0,
						Build: 0,
					},
					Namespaces: []Namespace{
						{
							ParentName: "plugin",
							Name:       "n1",
						},
					},
					Commands: []Command{
						{
							Namespace:   "n1",
							Name:        "command",
							Usage:       "ibmcloud n1 command ([--DEF] | ABD)",
							Description: "List all installed things",
							Flags: []Flag{
								{
									Name:        "\u003cinvalid Value\u003e",
									Description: "Description for flag DEF",
								},
							},
						},
					},
				},
			},
			errors: PluginToValidationErrors{
				"plugin": []PluginMetadataError{
					{
						CommandName: "n1 command",
						Namespace:   "PluginMetadata.Commands[0].Flags[0].Name",
						Error:       "Name contains the following forbidden characters: < >",
						Priority:    PriorityError,
					},
				},
			},
		},
		{
			name: "Command name is --version",
			pluginMetadata: []PluginMetadata{
				{
					Name: "plugin",
					Version: VersionType{
						Major: 1,
						Minor: 0,
						Build: 0,
					},
					MinCliVersion: VersionType{
						Major: 2,
						Minor: 0,
						Build: 0,
					},
					Namespaces: []Namespace{
						{
							ParentName: "plugin",
							Name:       "n1",
						},
					},
					Commands: []Command{
						{
							Namespace:   "n1",
							Name:        "--version",
							Usage:       "ibmcloud n1 --version ([--DEF] | ABD)",
							Description: "List all the things",
							Flags: []Flag{
								{
									Name:        "DEF",
									Description: "Description for flag DEF",
								},
							},
						},
					},
				},
			},
			errors: PluginToValidationErrors{
				"plugin": []PluginMetadataError{
					{
						CommandName: "n1 --version",
						Namespace:   "PluginMetadata.Commands[0].Name",
						Error:       "Command '--version' uses a reserved flag name. These are handled by the CLI framework.",
						Priority:    PriorityWarning,
						Remediation: "Remove this command; --version and --help are automatically provided.",
					},
				},
			},
		},
		{
			name: "Invalid minimum CLI version",
			pluginMetadata: []PluginMetadata{
				{
					Name: "plugin",
					Version: VersionType{
						Major: 1,
						Minor: 0,
						Build: 0,
					},
					MinCliVersion: VersionType{
						Major: 0,
						Minor: 0,
						Build: 0,
					},
					Namespaces: []Namespace{
						{
							ParentName: "namespace",
							Name:       "command",
						},
					},
					Commands: []Command{
						{
							Namespace:   "plugin",
							Name:        "command",
							Description: "A sample command",
							Usage:       "ibmcloud plugin command ABC (--DEF | --GHI)",
							Flags: []Flag{
								{
									Name:        "DEF",
									Description: "Description for flag DEF",
								},
								{
									Name:        "GHI",
									Description: "Description for flag GHI",
								},
							},
						},
					},
				},
			},
			errors: PluginToValidationErrors{
				"plugin": []PluginMetadataError{
					{
						CommandName: "",
						Namespace:   "PluginMetadata.MinCliVersion",
						Error:       "MinCliVersion (0.0.0) is lower than the allowed minimum 2.0.0",
						Remediation: "Set MinCliVersion to 2.0.0 or higher to ensure compatibility with supported CLI versions.",
						Priority:    PriorityError,
					},
				},
			},
		},
		{
			name: "Multiple errors",
			pluginMetadata: []PluginMetadata{
				{
					Name: "plugin",
					Version: VersionType{
						Major: 1,
						Minor: 0,
						Build: 0,
					},
					MinCliVersion: VersionType{
						Major: 2,
						Minor: 0,
						Build: 0,
					},
					Namespaces: []Namespace{
						{
							ParentName: "namespace",
							Name:       "command",
						},
					},
					Commands: []Command{
						{

							Namespace: "n1",
							Name:      "--version",
							Alias:     "",
							Aliases: []string{
								"-v",
							},
						},
					},
				},
			},
			errors: PluginToValidationErrors{
				"plugin": []PluginMetadataError{
					{
						CommandName: "n1 --version",
						Namespace:   "PluginMetadata.Commands[0].Description",
						Error:       "Command '--version' has no description. All commands must have a clear description.",
						Remediation: "Add a sentence without subject describing what the command does.",
						Priority:    PriorityError,
					},
					{
						CommandName: "n1 --version",
						Namespace:   "PluginMetadata.Commands[0].Usage",
						Error:       "Command '--version' has no usage information.",
						Remediation: "Add usage text showing command syntax with parameters and options.",
						Priority:    PriorityError,
					},
					{
						CommandName: "n1 --version",
						Namespace:   "PluginMetadata.Commands[0].Name",
						Error:       "Command '--version' uses a reserved flag name. These are handled by the CLI framework.",
						Priority:    PriorityWarning,
						Remediation: "Remove this command; --version and --help are automatically provided.",
					},
				},
			},
		},
		{
			name: "Command exceeds maximum depth",
			pluginMetadata: []PluginMetadata{
				{
					Name: "plugin",
					Version: VersionType{
						Major: 1,
						Minor: 0,
						Build: 0,
					},
					MinCliVersion: VersionType{
						Major: 2,
						Minor: 0,
						Build: 0,
					},
					Namespaces: []Namespace{
						{
							ParentName: "plugin",
							Name:       "mynamespace",
						},
					},
					Commands: []Command{
						{
							Namespace:   "mynamespace",
							Name:        "level1 level2 level3 level4 create",
							Description: "Create a resource",
							Usage:       "ibmcloud mynamespace level1 level2 level3 level4 create RESOURCE",
						},
					},
				},
			},
			errors: PluginToValidationErrors{
				"plugin": []PluginMetadataError{
					{
						CommandName: "mynamespace level1 level2 level3 level4 create",
						Namespace:   "PluginMetadata.Commands[0].Name",
						Error:       "Command 'level1 level2 level3 level4 create' has 3 levels, exceeding the maximum of 3. Deep command hierarchies are difficult for users to remember and discover.",
						Priority:    PriorityWarning,
						Remediation: "Reduce command depth to 3 or fewer levels. Options: (1) Flatten the hierarchy by combining levels, (2) Use command flags/options instead of subcommands, (3) Reorganize the command structure to be more intuitive.",
					},
				},
			},
		},
		{
			name: "Command with long description (Info priority)",
			pluginMetadata: []PluginMetadata{
				{
					Name: "plugin",
					Version: VersionType{
						Major: 1,
						Minor: 0,
						Build: 0,
					},
					MinCliVersion: VersionType{
						Major: 2,
						Minor: 0,
						Build: 0,
					},
					Namespaces: []Namespace{
						{
							ParentName: "plugin",
							Name:       "n1",
						},
					},
					Commands: []Command{
						{
							Namespace:   "n1",
							Name:        "command",
							Description: "A very long description that contains more than fifteen words which should trigger an info level validation error message",
							Usage:       "ibmcloud n1 command [OPTIONS]",
						},
					},
				},
			},
			errors: PluginToValidationErrors{
				"plugin": []PluginMetadataError{
					{
						CommandName: "n1 command",
						Namespace:   "PluginMetadata.Commands[0].Description",
						Error:       "Description for 'command' has 19 words. Consider limiting to less than 15 words for better display.",
						Priority:    PriorityInfo,
						Remediation: "Shorten the description to be more concise.",
					},
				},
			},
		},
	}

	for _, d := range testData {
		t.Run(d.name, func(t *testing.T) {
			actualErrors := validator.Errors(d.pluginMetadata)

			assert.Equal(t, len(d.errors), len(actualErrors))

			for pluginName, errs := range actualErrors {
				expectedErrs, exists := d.errors[pluginName]
				assert.True(t, exists, "Unexpected plugin '%s' in errors", pluginName)
				assert.Equal(t, len(expectedErrs), len(errs), "Error count mismatch for plugin '%s'", pluginName)

				// Check that each expected error is present (order-independent)
				for _, expectedErr := range expectedErrs {
					found := false
					for _, actualErr := range errs {
						if expectedErr.CommandName == actualErr.CommandName &&
							expectedErr.Error == actualErr.Error &&
							expectedErr.Namespace == actualErr.Namespace {
							found = true
							break
						}
					}
					assert.True(t, found, "Expected error not found: %s - %s", expectedErr.Namespace, expectedErr.Error)
				}

				// Original order-dependent check for backwards compatibility (kept for reference)
				for idx, err := range errs {
					if idx < len(expectedErrs) {
						_ = err
						_ = expectedErrs[idx]
						assert.Equal(t, expectedErrs[idx].Priority, err.Priority, "Priority mismatch at index %d", idx)
					}
				}
			}
		})
	}
}
