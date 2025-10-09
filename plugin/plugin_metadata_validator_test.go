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
							Name:        "command",
							Description: "A sample command",
							Usage:       "command ABC (--DEF | --GHI)",
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
							Usage:       "command ABC (--DEF | --GHI)",
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
							Usage:       "command ABC (--DEF | --GHI)",
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
							Usage:       "command ABC (--DEF | --GHI)",
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
							Usage:       "command ABC (--DEF | --GHI)",
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
							Usage:       "command ABC (--DEF | --GHI)",
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
							Usage:     "n1 command ABC (--DEF | --GHI)",
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
						Error:       "Description is required",
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
						Namespace:   "PluginMetadata.Commands[0].Usage",
						Error:       "Usage is required",
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
						Namespace:   "PluginMetadata.Commands[0].Usage",
						Error:       "Usage is required",
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
							Usage:       "command [options...]",
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
							Usage:       "command COMMAND",
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
							Usage:       "command ([--DEF | ABD)",
							Description: "This is a command",
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
						Error:       "Usage contains unclosed [ between indicies [9 22]",
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
							Usage:       "command ([--DEF] | ABD)",
							Description: "This is a command",
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
							Usage:       "command ([--DEF] | ABD)",
							Description: "This is a command",
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
						Error:       "Name must not equal '--version'",
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
							Usage:       "command ABC (--DEF | --GHI)",
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
						Namespace:   "PluginMetadata.Commands[0].Name",
						Error:       "Name must not equal '--version'",
					},
					{
						CommandName: "n1 --version",
						Namespace:   "PluginMetadata.Commands[0].Description",
						Error:       "Description is required",
					},
					{

						CommandName: "n1 --version",
						Namespace:   "PluginMetadata.Commands[0].Usage",
						Error:       "Usage is required",
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
				assert.NotEmpty(t, d.errors[pluginName])
				for idx, err := range errs {
					assert.Equal(t, err.CommandName, d.errors[pluginName][idx].CommandName)
					assert.Equal(t, err.Error, d.errors[pluginName][idx].Error)
					assert.Equal(t, err.Namespace, d.errors[pluginName][idx].Namespace)
				}
			}
		})
	}
}
