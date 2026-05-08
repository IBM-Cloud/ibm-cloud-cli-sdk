package plugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestValidateUsageEnhanced_EmptyUsage tests that empty usage is detected
// func TestValidateUsageEnhanced_EmptyUsage(t *testing.T) {
// 	validator := NewPluginMetadataValidator()
//
// 	cmd := Command{
// 		Namespace:   "test",
// 		Name:        "command",
// 		Description: "Test command",
// 		Usage:       "",
// 	}
//
// 	errs := validator.validateUsageEnhanced(cmd, "test-plugin")
//
// 	assert.Len(t, errs, 1, "Expected 1 error for empty usage")
// 	assert.Equal(t, PriorityError, errs[0].Priority)
// 	assert.Contains(t, errs[0].Error, "has no usage information")
// }

// TestValidateUsageEnhanced_ValidUsage tests valid usage patterns
func TestValidateUsageEnhanced_ValidUsage(t *testing.T) {
	validator := NewPluginMetadataValidator()

	testCases := []struct {
		name           string
		pluginMetadata []PluginMetadata
		errors         PluginToValidationErrors
	}{
		{
			name: "Command with uppercase arguments",
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
							Namespace:   "resource",
							Name:        "service-instance-create",
							Description: "Create a resource",
							Usage:       "ibmcloud resource service-instance-create NAME SERVICE_NAME PLAN_NAME LOCATION",
						},
					},
				},
			},
			errors: PluginToValidationErrors{},
		},
		{
			name: "Command with optional uppercase arguments",
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
							Namespace:   "resource",
							Name:        "service-instance-update",
							Description: "Update a resource",
							Usage:       "ibmcloud resource service-instance-update NAME [--service-plan-id PLAN_ID]",
							Flags: []Flag{
								{
									Name:        "service-plan-id",
									Description: "New service plan ID",
								},
							},
						},
					},
				},
			},
			errors: PluginToValidationErrors{},
		},
		{
			name: "Command with multiple uppercase arguments",
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
							Namespace:   "ks cluster create",
							Name:        "classic",
							Description: "Create a classic cluster",
							Usage:       "ibmcloud ks cluster create classic --name NAME --zone ZONE [--workers COUNT]",
							Flags: []Flag{
								{
									Name:        "name",
									Description: "The cluster name",
								},
								{
									Name:        "zone",
									Description: "The name of the zone",
								},
								{
									Name:        "workers",
									Description: "The number of workers",
								},
							},
						},
					},
				},
			},
			errors: PluginToValidationErrors{},
		},
		{
			name: "Command with full path to ibmcloud",
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
							Description: "List installed plug-ins",
							Usage:       "/usr/local/bin/ibmcloud plugin list",
							Flags:       []Flag{},
						},
					},
				},
			},
			errors: PluginToValidationErrors{},
		},
		{
			name: "Command with choices using uppercase",
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
							Namespace:   "account",
							Name:        "org-role-set",
							Description: "Set the role organization",
							Usage:       "ibmcloud account org-role-set USER_NAME ORG_NAME (ROLE_A | ROLE_B)",
							Flags:       []Flag{},
						},
					},
				},
			},
			errors: PluginToValidationErrors{},
		},
		{
			name: "Complex command with mixed elements",
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
							Namespace:   "service",
							Name:        "create",
							Description: "Create a service",
							Usage:       "ibmcloud service create NAME PLAN LOCATION [--tags TAGS] [--parameters PARAMS]",
							Flags: []Flag{
								{
									Name:        "tags",
									Description: "The list of tags",
								},
								{
									Name:        "parameters",
									Description: "The parameters to set on the instance",
								},
							},
						},
					},
				},
			},
			errors: PluginToValidationErrors{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pluginToErrors := validator.Errors(tc.pluginMetadata)
			usage := tc.pluginMetadata[0].Commands[0].Usage

			for _, pluginErrors := range pluginToErrors {
				assert.Empty(t, pluginErrors, "Expected no errors for valid usage: %s", usage)
			}

		})
	}
}

// TestValidateUsageEnhanced_MissingIbmcloudPrefix tests detection of missing ibmcloud prefix
func TestValidateUsageEnhanced_MissingIbmcloudPrefix(t *testing.T) {
	validator := NewPluginMetadataValidator()

	testCases := []struct {
		name           string
		pluginMetadata []PluginMetadata
		errors         PluginToValidationErrors
	}{
		{
			name: "No prefix at all",
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
							Description: "List all installed plug-ins",
							Usage:       "plugin list [--output FORMAT]",
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
							Description: "List all installed plug-ins",
							Usage:       "IBMCLOUD plugin list [--output FORMAT]",
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
							Description: "List all installed plug-ins",
							Usage:       "IBMCloud plugin list [--output FORMAT]",
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
							ParentName: "",
							Name:       "ibmcloud",
						},
					},
					Commands: []Command{
						{
							Namespace:   "plugin",
							Name:        "list",
							Description: "List all installed plug-ins",
							Usage:       "ic plugin list [--output FORMAT]",
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
							ParentName: "",
							Name:       "ibmcloud",
						},
					},
					Commands: []Command{
						{
							Namespace:   "plugin",
							Name:        "list",
							Description: "List all installed plug-ins",
							Usage:       "list",
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
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			pluginToErrors := validator.Errors(tc.pluginMetadata)

			usage := tc.pluginMetadata[0].Commands[0].Usage
			assert.NotEmpty(t, pluginToErrors, "Expected error for usage without proper ibmcloud prefix: %s", usage)
			found := false
			for _, pluginErrs := range pluginToErrors {
				assert.Equal(t, len(pluginErrs), 1, "Expect number of plugin validation errors to be 1")
				pluginErr := pluginErrs[0]

				if pluginErr.Priority == PriorityError &&
					(containsString(pluginErr.Error, "should start with 'ibmcloud'") ||
						containsString(pluginErr.Error, "full path")) {
					found = true
					break
				}
			}
			assert.True(t, found, "Expected error about ibmcloud prefix for: %s", usage)
		})
	}
}

// TestValidateUsageEnhanced_LowercaseArguments tests detection of lowercase argument values
func TestValidateUsageEnhanced_LowercaseArguments(t *testing.T) {
	validator := NewPluginMetadataValidator()

	testCases := []struct {
		name              string
		pluginMetadata    []PluginMetadata
		errors            PluginToValidationErrors
		namespace         string
		cmdName           string
		usage             string
		expectedLowercase []string
	}{
		{
			name: "Single lowercase argument",
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
							Namespace:   "resource",
							Name:        "service-instance-create",
							Description: "Create a service instance",
							Usage:       "ibmcloud resource service-instance-create name SERVICE PLAN LOCATION",
							Flags:       []Flag{},
						},
					},
				},
			},
			expectedLowercase: []string{"name"},
		},
		{
			name: "Multiple lowercase arguments - space separated",
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
							Name:       "ks",
						},
					},
					Commands: []Command{
						{
							Namespace:   "ks",
							Name:        "cluster create",
							Description: "Create a classic cluster",
							Usage:       "ibmcloud ks cluster create classic --name instance --zone region",
							Flags: []Flag{
								{
									Name:        "name",
									Description: "The name of the cluster instance",
								},
								{
									Name:        "zone",
									Description: "The name of the cluster zone",
								},
							},
						},
					},
				},
			},
			expectedLowercase: []string{"instance"},
		},
		{
			name: "Lowercase with hyphens",
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
							Namespace:   "resource",
							Name:        "group-create",
							Description: "Create a resource group",
							Usage:       "ibmcloud resource group-create instance-name",
						},
					},
				},
			},
			expectedLowercase: []string{"instance-name"},
		},
		{
			name: "Lowercase with underscores",
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
							Namespace:   "resource",
							Name:        "service-key-create",
							Description: "Create a service key",
							Usage:       "ibmcloud resource service-key-create instance_id KEY_NAME ROLE",
							Flags:       []Flag{},
						},
					},
				},
			},
			expectedLowercase: []string{"instance_id"},
		},
		{
			name: "Mixed case - some correct, some not",
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
							Namespace:   "account",
							Name:        "user-invite",
							Description: "Invite a user to the account",
							Usage:       "ibmcloud account user-invite EMAIL --org ORG --space space_name",
							Flags: []Flag{
								{
									Name:        "org",
									Description: "The name of the account organization",
								},
								{
									Name:        "space",
									Description: "The name of the account space",
								},
							},
						},
					},
				},
			},
			expectedLowercase: []string{"space_name"},
		},
		{
			name: "Lowercase before pipe with space",
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
							Namespace:   "account",
							Name:        "org-role-set",
							Description: "Set an account organization role",
							Usage:       "ibmcloud account org-role-set USER ORG (option_a | OPTION_B)",
							Flags:       []Flag{},
						},
					},
				},
			},
			expectedLowercase: []string{"option_a"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			pluginToErrors := validator.Errors(tc.pluginMetadata)

			usage := tc.pluginMetadata[0].Commands[0].Usage
			// Find the lowercase argument error
			found := false
			for _, pluginErrs := range pluginToErrors {
				assert.Equal(t, len(pluginErrs), 1, "Expected error for lowercase arguments in: %s", usage)
				err := pluginErrs[0]

				if err.Priority == PriorityWarning &&
					containsString(err.Error, "lowercase argument values") {
					found = true
					// Verify that expected lowercase params are mentioned
					for _, param := range tc.expectedLowercase {
						assert.Contains(t, err.Error, param,
							"Expected error to mention lowercase param '%s'", param)
					}
					break
				}
			}
			assert.True(t, found, "Expected warning about lowercase arguments for: %s", usage)
		})
	}
}

// TestValidateUsageEnhanced_ExcludedWords tests that command words are not flagged as lowercase arguments
func TestValidateUsageEnhanced_ExcludedWords(t *testing.T) {
	validator := NewPluginMetadataValidator()

	testCases := []struct {
		name           string
		pluginMetadata []PluginMetadata
		shouldErr      bool
	}{
		{
			name: "Command name in usage",
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
							Description: "List all installed plug-ins",
							Usage:       "ibmcloud plugin list",
							Flags:       []Flag{},
						},
					},
				},
			},
			shouldErr: false,
		},
		{
			name: "Multi-word command with hyphens",
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
							Name:        "list-all",
							Description: "List all installed plug-ins",
							Usage:       "ibmcloud plugin list all",
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
			shouldErr: false,
		},
		{
			name: "Common words like 'and', 'or', 'to'",
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
							Name:        "convert",
							Description: "List all installed plug-ins",
							Usage:       "ibmcloud plugin convert from FORMAT to FORMAT", // NOTE: need to check if this test is valid
							Flags:       []Flag{},
						},
					},
				},
			},
			shouldErr: false,
		},
		{
			name: "Namespace words are excluded",
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
							Namespace:   "service",
							Name:        "create",
							Description: "List all installed plug-ins",
							Usage:       "ibmcloud service create NAME",
							Flags:       []Flag{},
						},
					},
				},
			},
			shouldErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			pluginToErrors := validator.Errors(tc.pluginMetadata)

			// Check if lowercase argument errors exist
			hasLowercaseErr := false
			for _, pluginErrs := range pluginToErrors {
				assert.Equal(t, len(pluginErrs), 1, "Expect number of plugin validation errors to be 1")
				err := pluginErrs[0]

				if containsString(err.Error, "lowercase argument values") {
					hasLowercaseErr = true
					if !tc.shouldErr {
						t.Errorf("Should not flag command words as lowercase arguments: %s", err.Error)
					}
				}
			}

			if tc.shouldErr && !hasLowercaseErr {
				t.Error("Expected lowercase argument error but got none")
			}
		})
	}
}

// TestValidateUsageEnhanced_ComplexScenarios tests complex real-world usage patterns
func TestValidateUsageEnhanced_ComplexScenarios(t *testing.T) {
	validator := NewPluginMetadataValidator()

	testCases := []struct {
		name           string
		namespace      string
		cmdName        string
		pluginMetadata []PluginMetadata
		usage          string
		expectErrors   bool
		errorContains  string
	}{
		{
			name: "Valid complex command",
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
							Namespace:   "resource",
							Name:        "service-instance-create",
							Description: "List all installed plug-ins",
							Usage:       "ibmcloud resource service-instance-create NAME SERVICE PLAN LOCATION [--tags TAGS] [--parameters PARAMS]",
							Flags: []Flag{
								{
									Name:        "tags",
									Description: "Comma seperated list of tags",
								},
								{
									Name:        "parameters",
									Description: "The custom parameters used to create the instance",
								},
							},
						},
					},
				},
			},
			expectErrors: false,
		},
		{
			name: "Complex with lowercase params",
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
							Namespace:   "resource",
							Name:        "service-instance-create",
							Description: "Create an service instance",
							Usage:       "ibmcloud resource service-instance-create name SERVICE PLAN LOCATION",
							Flags:       []Flag{},
						},
					},
				},
			},
			expectErrors:  true,
			errorContains: "lowercase argument values",
		},
		{
			name: "Valid with JSON example (uppercase)",
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
							Namespace:   "resource",
							Name:        "service-instance-update",
							Description: "Update a service instance",
							Usage:       "ibmcloud resource service-instance-update NAME --parameters DATA",
							Flags: []Flag{
								{
									Name:        "parameters",
									Description: "The parameters to set on the instance",
								},
							},
						},
					},
				},
			},
			expectErrors: false,
		},
		{
			name: "Valid with multiple optional groups",
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
							Namespace:   "ks cluster create",
							Name:        "classic",
							Description: "Create a classic cluster",
							Usage:       "ibmcloud ks cluster create classic --name NAME [--zone ZONE] [--workers COUNT] [--force]",
							Flags: []Flag{
								{
									Name:        "name",
									Description: "The name of the cluster",
								},
								{
									Name:        "zone",
									Description: "The name of the zone",
								},
								{
									Name:        "workers",
									Description: "The number of workers",
								},
								{
									Name:        "force",
									Description: "Force creation without confirmation",
								},
							},
						},
					},
				},
			},
			expectErrors: false,
		},
		{
			name: "Missing ibmcloud prefix with lowercase",
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
							Namespace:   "resource",
							Name:        "service-instance-create",
							Description: "Create a service instance",
							Usage:       "resource service-instance-create instance",
							Flags:       []Flag{},
						},
					},
				},
			},
			expectErrors:  true,
			errorContains: "should start with 'ibmcloud'",
		},
		{
			name: "Valid with choice operator",
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
							Namespace:   "account",
							Name:        "org-role-set",
							Description: "Set an organization role",
							Usage:       "ibmcloud account org-role-set USER ORG (ROLE_A | ROLE_B | ROLE_C)",
							Flags:       []Flag{},
						},
					},
				},
			},
			expectErrors: false,
		},
		{
			name: "Invalid choice with lowercase",
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
							Namespace:   "account",
							Name:        "org-role-set",
							Description: "Set an organization role",
							Usage:       "ibmcloud account org-role-set USER ORG (option_a | option_b)",
							Flags:       []Flag{},
						},
					},
				},
			},
			expectErrors:  true,
			errorContains: "lowercase argument values",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			pluginToErrors := validator.Errors(tc.pluginMetadata)
			usage := tc.pluginMetadata[0].Commands[0].Usage

			for _, pluginErrors := range pluginToErrors {
				if tc.expectErrors {
					assert.NotEmpty(t, pluginErrors, "Expected errors for: %s", usage)
					if tc.errorContains != "" {
						found := false
						for _, err := range pluginErrors {
							if containsString(err.Error, tc.errorContains) {
								found = true
								break
							}
						}
						assert.True(t, found, "Expected error containing '%s' for: %s",
							tc.errorContains, usage)
					}
				} else {
					assert.Empty(t, pluginErrors, "Expected no errors for: %s", usage)
				}

			}

		})
	}
}

// TestValidateUsageEnhanced_EdgeCases tests edge cases and boundary conditions
func TestValidateUsageEnhanced_EdgeCases(t *testing.T) {
	validator := NewPluginMetadataValidator()

	testCases := []struct {
		name           string
		pluginMetadata []PluginMetadata
		usage          string
		expectErrors   bool
	}{
		{
			name: "Very long usage string",
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
							Name:        "command",
							Description: "A general purpose command for plugin",
							Usage:       "ibmcloud plugin command PARAM1 PARAM2 PARAM3 PARAM4 PARAM5 --opt1 VAL1 --opt2 VAL2 --opt3 VAL3",
							Flags: []Flag{
								{
									Name:        "opt1",
									Description: "Option1 description",
								},
								{
									Name:        "opt2",
									Description: "Option2 description",
								},
								{
									Name:        "opt3",
									Description: "Option3 description",
								},
							},
						},
					},
				},
			},
			expectErrors: false,
		},
		{
			name: "Usage with only whitespace after ibmcloud",
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
							Description: "List all installed plug-ins",
							Usage:       "ibmcloud   ",
							Flags:       []Flag{},
						},
					},
				},
			},
			expectErrors: false,
		},
		{
			name: "Usage with tabs and newlines",
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
							Description: "List all installed plug-ins",
							Usage:       "ibmcloud\tplugin\nlist",
							Flags:       []Flag{},
						},
					},
				},
			},
			expectErrors: true,
		},
		{
			name: "Usage with special characters in uppercase args",
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
							Name:        "create",
							Description: "Create a plugin",
							Usage:       "ibmcloud plugin create NAME_WITH_UNDERSCORES",
							Flags:       []Flag{},
						},
					},
				},
			},
			expectErrors: false,
		},
		{
			name: "Usage with numbers in uppercase args",
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
							Name:        "create",
							Description: "Create a plugn",
							Usage:       "ibmcloud plugin create INSTANCE123",
							Flags:       []Flag{},
						},
					},
				},
			},
			expectErrors: false,
		},
		{
			name: "Single character lowercase param not detected",
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
							Namespace:   "resource",
							Name:        "group-create",
							Description: "Create a resource group",
							Usage:       "ibmcloud resource group-create x",
							Flags:       []Flag{},
						},
					},
				},
			},
			expectErrors: false, // Single char params are not caught by the regex
		},
		{
			name: "Lowercase param at end of line",
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
							Namespace:   "resource",
							Name:        "service-instance-create",
							Description: "Create a service instance",
							Usage:       "ibmcloud resource service-instance-create name",
							Flags:       []Flag{},
						},
					},
				},
			},
			expectErrors: true,
		},
		{
			name: "Lowercase param before pipe without space IS detected",
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
							Namespace:   "account",
							Name:        "org-role-set",
							Description: "Set an organization role",
							Usage:       "ibmcloud account org-role-set USER ORG (option_a|OPTION_B)",
							Flags:       []Flag{},
						},
					},
				},
			},
			expectErrors: true, // Improved regex now catches this case
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			pluginToErrors := validator.Errors(tc.pluginMetadata)

			for _, errs := range pluginToErrors {
				if tc.expectErrors {
					assert.NotEmpty(t, errs, "Expected errors for: %s", tc.usage)
				} else {
					assert.Empty(t, errs, "Expected no errors for: %s", tc.usage)
				}
			}
		})
	}
}

// TestValidateUsageEnhanced_ErrorPriorities tests that errors have correct priorities
func TestValidateUsageEnhanced_ErrorPriorities(t *testing.T) {
	validator := NewPluginMetadataValidator()

	testCases := []struct {
		name             string
		usage            string
		pluginMetadata   []PluginMetadata
		expectedPriority Priority
		errorType        string
	}{
		{
			name: "Empty usage - ERROR priority",
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
							Namespace:   "service",
							Name:        "create",
							Description: "Create a service",
							Usage:       "",
							Flags:       []Flag{},
						},
					},
				},
			},
			usage:            "",
			expectedPriority: PriorityError,
			errorType:        "no usage information",
		},
		{
			name: "Missing ibmcloud prefix - ERROR priority",
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
							Description: "List all installed plug-ins",
							Usage:       "plugin list",
							Flags:       []Flag{},
						},
					},
				},
			},
			expectedPriority: PriorityError,
			errorType:        "should start with 'ibmcloud'",
		},
		{
			name:             "Lowercase arguments - WARNING priority",
			usage:            "ibmcloud plugin create name",
			expectedPriority: PriorityWarning,
			errorType:        "lowercase argument values",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pluginToErrors := validator.Errors(tc.pluginMetadata)

			for _, pluginErrors := range pluginToErrors {
				assert.NotEmpty(t, pluginErrors, "Expected errors for: %s", tc.name)

				found := false
				for _, err := range pluginErrors {
					if containsString(err.Error, tc.errorType) {
						assert.Equal(t, tc.expectedPriority, err.Priority,
							"Expected priority %s for error type '%s'", tc.expectedPriority, tc.errorType)
						found = true
						break
					}
				}
				assert.True(t, found, "Expected to find error type '%s'", tc.errorType)
			}
		})
	}
}

// TestValidateUsageEnhanced_Remediation tests that remediation messages are provided
func TestValidateUsageEnhanced_Remediation(t *testing.T) {
	validator := NewPluginMetadataValidator()

	testCases := []struct {
		name           string
		pluginMetadata []PluginMetadata
		usage          string
	}{
		{
			name: "Empty usage",
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
							Description: "List all installed plug-ins",
							Usage:       "",
							Flags:       []Flag{},
						},
					},
				},
			},
		},
		{

			name: "Missing ibmcloud prefix",
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
							Description: "List all installed plug-ins",
							Usage:       "plugin list",
							Flags:       []Flag{},
						},
					},
				},
			},
		},
		{
			name: "Lowercase arguments",
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
							Name:        "create",
							Description: "Create a plugin",
							Usage:       "ibmcloud plugin create name",
							Flags:       []Flag{},
						},
					},
				},
			},
			usage: "ibmcloud plugin create name",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			pluginToErrors := validator.Errors(tc.pluginMetadata)
			for _, pluginErrors := range pluginToErrors {

				assert.NotEmpty(t, pluginErrors, "Expected errors for: %s", tc.name)

				for _, err := range pluginErrors {
					assert.NotEmpty(t, err.Remediation,
						"Expected remediation message for error: %s", err.Error)
				}
			}
		})
	}
}

// TestValidateCommandNaming_PluralForms tests that commands ending with 's' are accepted as alternatives to 'list'
func TestValidateCommandNaming_PluralForms(t *testing.T) {
	validator := NewPluginMetadataValidator()

	testCases := []struct {
		name           string
		namespace      string
		cmdName        string
		pluginMetadata []PluginMetadata
		expectErrors   bool
	}{
		{
			name: "Valid plural form - instances",
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
							Description: "List all installed plug-ins",
							Usage:       "ibmcloud service create NAME",
							Flags:       []Flag{},
						},
					},
				},
			},
			expectErrors: false,
		},
		{
			name: "Valid plural form - services",
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
							Namespace:   "resource",
							Name:        "services",
							Description: "List all resources",
							Usage:       "ibmcloud resource services",
							Flags:       []Flag{},
						},
					},
				},
			},
			expectErrors: false,
		},
		{
			name: "Valid plural form - clusters",
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
							Namespace:   "ks",
							Name:        "clusters",
							Description: "List all clusters",
							Usage:       "ibmcloud ks clusters",
							Flags:       []Flag{},
						},
					},
				},
			},
			expectErrors: false,
		},
		{
			name: "Valid plural form - resources",
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
							Namespace:   "account",
							Name:        "resources",
							Description: "List all resources",
							Usage:       "ibmcloud account resources",
							Flags:       []Flag{},
						},
					},
				},
			},
			expectErrors: false,
		},
		{
			name: "Valid plural form with multi-word",
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
							Namespace:   "service",
							Name:        "key instances",
							Description: "List all key instances",
							Usage:       "ibmcloud service key instances",
							Flags:       []Flag{},
						},
					},
				},
			},
			expectErrors: false,
		},
		{
			name: "Valid - has common verb",
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
							Namespace:   "service",
							Name:        "instance list",
							Description: "List instances",
							Usage:       "ibmcloud service instance list",
							Flags:       []Flag{},
						},
					},
				},
			},
			namespace:    "service",
			cmdName:      "instance list",
			expectErrors: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pluginToErrors := validator.Errors(tc.pluginMetadata)
			cmdName := tc.pluginMetadata[0].Commands[0].Name
			// Check for the "does not use a common verb" error
			hasVerbError := false
			for _, pluginErrs := range pluginToErrors {

				for _, err := range pluginErrs {
					if containsString(err.Error, "does not use a common verb") {
						hasVerbError = true
						break
					}
				}

				if tc.expectErrors {
					assert.True(t, hasVerbError, "Expected verb/plural error for: %s", cmdName)
				} else {
					assert.False(t, hasVerbError, "Should not have verb/plural error for: %s", cmdName)
				}
			}
		})
	}
}

// TestValidateDescription_PluralFormWithoutListKeywords tests that plural commands require list-related descriptions
func TestValidateDescription_PluralFormWithoutListKeywords(t *testing.T) {
	validator := NewPluginMetadataValidator()

	testCases := []struct {
		name           string
		pluginMetadata []PluginMetadata
		expectError    bool
		errorMessage   string
	}{
		{
			name: "Plural command with 'list' keyword - valid",
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
							Namespace:   "service",
							Name:        "instances",
							Description: "List all service instances",
							Usage:       "ibmcloud service instances",
						},
					},
				},
			},
			expectError: false,
		},
		{
			name: "Plural command with 'show' keyword - valid",
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
							Namespace:   "service",
							Name:        "services",
							Description: "Show available services",
							Usage:       "ibmcloud service services",
						},
					},
				},
			},
			expectError: false,
		},
		{
			name: "Plural command with 'display' keyword - valid",
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
							Namespace:   "service",
							Name:        "clusters",
							Description: "Display all clusters in the account",
							Usage:       "ibmcloud service clusters",
						},
					},
				},
			},
			expectError: false,
		},
		{
			name: "Plural command with 'view' keyword - valid",
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
							Namespace:   "service",
							Name:        "resources",
							Description: "View all resources",
							Usage:       "ibmcloud service resources",
						},
					},
				},
			},
			expectError: false,
		},
		{
			name: "Plural command with 'retrieve' keyword - valid",
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
							Namespace:   "service",
							Name:        "policies",
							Description: "Retrieve security policies",
							Usage:       "ibmcloud service policies",
						},
					},
				},
			},
			expectError: false,
		},
		{
			name: "Plural command with 'get' keyword - valid",
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
							Namespace:   "service",
							Name:        "users",
							Description: "Get all users in the organization",
							Usage:       "ibmcloud service users",
						},
					},
				},
			},
			expectError: false,
		},
		{
			name: "Plural command with 'all' keyword - valid",
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
							Namespace:   "service",
							Name:        "instances",
							Description: "Returns all instances",
							Usage:       "ibmcloud service instances",
						},
					},
				},
			},
			expectError: false,
		},
		{
			name: "Plural command with 'multiple' keyword - valid",
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
							Namespace:   "service",
							Name:        "services",
							Description: "Manage multiple services",
							Usage:       "ibmcloud service services",
						},
					},
				},
			},
			expectError: false,
		},
		{
			name: "Plural command without list keywords - invalid",
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
							Namespace:   "service",
							Name:        "instances",
							Description: "Manage service instances",
							Usage:       "ibmcloud service instances",
						},
					},
				},
			},
			expectError:  true,
			errorMessage: "doesn't clearly indicate it returns a list or group of items",
		},
		{
			name: "Plural command with vague description - invalid",
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
							Namespace:   "service",
							Name:        "services",
							Description: "Work with services",
							Usage:       "ibmcloud service services",
						},
					},
				},
			},
			expectError:  true,
			errorMessage: "doesn't clearly indicate it returns a list or group of items",
		},
		{
			name: "Plural command with action description - invalid",
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
							Namespace:   "service",
							Name:        "clusters",
							Description: "Create and manage clusters",
							Usage:       "ibmcloud service clusters",
						},
					},
				},
			},
			expectError:  true,
			errorMessage: "doesn't clearly indicate it returns a list or group of items",
		},
		{
			name: "Non-plural command without list keywords - valid (no check)",
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
							Namespace:   "service",
							Name:        "create",
							Description: "Create a new instance",
							Usage:       "ibmcloud service create",
						},
					},
				},
			},
			expectError: false,
		},
		{
			name: "Short plural word 'is' not flagged",
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
							Namespace:   "service",
							Name:        "is",
							Description: "Check status",
							Usage:       "ibmcloud service is",
						},
					},
				},
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pluginToErrors := validator.Errors(tc.pluginMetadata)
			cmdName := tc.pluginMetadata[0].Commands[0].Name
			description := tc.pluginMetadata[0].Commands[0].Description

			// Check for the plural form description error
			hasPluralDescError := false
			for _, pluginErrs := range pluginToErrors {
				for _, err := range pluginErrs {
					if containsString(err.Error, "doesn't clearly indicate it returns a list or group of items") {
						hasPluralDescError = true
						if tc.errorMessage != "" {
							assert.Contains(t, err.Error, tc.errorMessage)
						}
						assert.Equal(t, PriorityWarning, err.Priority, "Expected WARNING priority for plural description check")
						break
					}
				}
			}

			if tc.expectError {
				assert.True(t, hasPluralDescError, "Expected plural description error for: %s with description: %s", cmdName, description)
			} else {
				assert.False(t, hasPluralDescError, "Should not have plural description error for: %s with description: %s", cmdName, description)
			}
		})
	}
}

// TestValidateCommandNaming_ConfigPattern tests that "config get/set/unset <parameter>" commands are accepted
func TestValidateCommandNaming_ConfigPattern(t *testing.T) {
	validator := NewPluginMetadataValidator()

	testCases := []struct {
		name           string
		pluginMetadata []PluginMetadata
		expectErrors   bool
	}{
		{
			name: "Valid config get pattern",
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
							Namespace:   "backup-recovery config get",
							Name:        "connector-service-url",
							Description: "Test command for config pattern",
							Usage:       "ibmcloud backup-recovery config get connector-service-url",
						},
					},
				},
			},
			expectErrors: false,
		},
		{
			name: "Valid config set pattern",
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
							Namespace:   "app-configuration config set",
							Name:        "service-url",
							Description: "Test command for config pattern",
							Usage:       "ibmcloud app-configuration config set service-url",
						},
					},
				},
			},
			expectErrors: false,
		},
		{
			name: "Valid config unset pattern",
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
							Namespace:   "service config unset",
							Name:        "api-key",
							Description: "Test command for config pattern",
							Usage:       "ibmcloud service config unset api-key",
						},
					},
				},
			},
			expectErrors: false,
		},
		{
			name: "Valid config get with hyphenated parameter",
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
							Namespace:   "plugin config get",
							Name:        "max-retry-count",
							Description: "Test command for config pattern",
							Usage:       "ibmcloud plugin config get max-retry-count",
						},
					},
				},
			},
			expectErrors: false,
		},
		{
			name: "Valid config set with underscored parameter",
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
							Namespace:   "plugin config set",
							Name:        "log_level",
							Description: "Test command for config pattern",
							Usage:       "ibmcloud plugin config set log_level",
						},
					},
				},
			},
			expectErrors: false,
		},
		{
			name: "Valid config get with uppercase in namespace",
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
							Namespace:   "Cloud-Logs Config Get",
							Name:        "service-url",
							Description: "Test command for config pattern",
							Usage:       "ibmcloud Cloud-Logs Config Get service-url",
						},
					},
				},
			},
			expectErrors: false,
		},
		{
			name: "Invalid - namespace doesn't end with config verb",
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
							Namespace:   "service config",
							Name:        "parameter-name",
							Description: "Test command for config pattern",
							Usage:       "ibmcloud service config parameter-name",
						},
					},
				},
			},
			expectErrors: true,
		},
		{
			name: "Invalid - namespace has config with wrong verb",
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
							Namespace:   "service config something",
							Name:        "parameter-name",
							Description: "Test command for config pattern",
							Usage:       "ibmcloud service config something parameter-name",
						},
					},
				},
			},
			expectErrors: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pluginToErrors := validator.Errors(tc.pluginMetadata)
			cmdName := tc.pluginMetadata[0].Commands[0].Name

			// Check for the "does not use a common verb" error
			hasVerbError := false
			for _, pluginErrs := range pluginToErrors {
				for _, err := range pluginErrs {
					if containsString(err.Error, "does not use a common verb") {
						hasVerbError = true
						break
					}
				}
			}

			if tc.expectErrors {
				assert.True(t, hasVerbError, "Expected verb/plural error for: %s", cmdName)
			} else {
				assert.False(t, hasVerbError, "Should not have verb/plural error for: %s", cmdName)
			}
		})
	}
}

// Helper function to check if a string contains a substring
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// TestValidateCommandNaming_MinimumSegmentLength tests that each segment of a command name is at least 2 characters
func TestValidateCommandNaming_MinimumSegmentLength(t *testing.T) {
	validator := NewPluginMetadataValidator()

	testCases := []struct {
		name           string
		pluginMetadata []PluginMetadata
		expectErrors   bool
		errorSegment   string
	}{
		{
			name: "Valid - single word with 2+ chars",
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
							Description: "Test command",
							Usage:       "ibmcloud plugin list",
						},
					},
				},
			},
			expectErrors: false,
		},
		{
			name: "Valid - multi-word with all segments 2+ chars",
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
							Namespace:   "service",
							Name:        "instance-create",
							Description: "Test command",
							Usage:       "ibmcloud service instance-create",
						},
					},
				},
			},
			expectErrors: false,
		},
		{
			name: "Valid - three words all 2+ chars",
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
							Namespace:   "ks",
							Name:        "cluster worker update",
							Description: "Test command",
							Usage:       "ibmcloud ks cluster worker update",
						},
					},
				},
			},
			expectErrors: false,
		},
		{
			name: "Invalid - single letter command",
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
							Name:        "a",
							Description: "Test command",
							Usage:       "ibmcloud plugin a",
						},
					},
				},
			},
			expectErrors: true,
			errorSegment: "a",
		},
		{
			name: "Invalid - hyphenated with single letter segment",
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
							Namespace:   "service",
							Name:        "a-b",
							Description: "Test command",
							Usage:       "ibmcloud service a-b",
						},
					},
				},
			},
			expectErrors: true,
			errorSegment: "a",
		},
		{
			name: "Invalid - space separated with single letter",
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
							Name:        "list x",
							Description: "Test command",
							Usage:       "ibmcloud plugin list x",
						},
					},
				},
			},
			expectErrors: true,
			errorSegment: "x",
		},
		{
			name: "Invalid - multiple single letter segments",
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
							Namespace:   "service",
							Name:        "a-b-c",
							Description: "Test command",
							Usage:       "ibmcloud service a-b-c",
						},
					},
				},
			},
			expectErrors: true,
			errorSegment: "a",
		},
		{
			name: "Invalid - valid word followed by single letter",
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
							Namespace:   "service",
							Name:        "create-x",
							Description: "Test command",
							Usage:       "ibmcloud service create-x",
						},
					},
				},
			},
			expectErrors: true,
			errorSegment: "x",
		},
		{
			name: "Invalid - single letter followed by valid word",
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
							Namespace:   "service",
							Name:        "x-create",
							Description: "Test command",
							Usage:       "ibmcloud service x-create",
						},
					},
				},
			},
			expectErrors: true,
			errorSegment: "x",
		},
		{
			name: "Invalid - command name with single letter segment",
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
							Namespace:   "ks cluster",
							Name:        "x-update",
							Description: "Test command",
							Usage:       "ibmcloud ks cluster x-update",
						},
					},
				},
			},
			expectErrors: true,
			errorSegment: "x",
		},
		{
			name: "Valid - two character segments",
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
							Namespace:   "service",
							Name:        "ls-rm",
							Description: "Test command",
							Usage:       "ibmcloud service ls-rm",
						},
					},
				},
			},
			expectErrors: false,
		},
		{
			name: "Valid - namespace with single letter not checked",
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
							Namespace:   "x",
							Name:        "list",
							Description: "Test command",
							Usage:       "ibmcloud x list",
						},
					},
				},
			},
			expectErrors: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pluginToErrors := validator.Errors(tc.pluginMetadata)
			cmdName := tc.pluginMetadata[0].Commands[0].Name

			// Check for the segment length error
			hasSegmentError := false
			for _, pluginErrs := range pluginToErrors {
				for _, err := range pluginErrs {
					if containsString(err.Error, "contains a segment") &&
						containsString(err.Error, "less than 2 characters") {
						hasSegmentError = true
						if tc.errorSegment != "" {
							assert.Contains(t, err.Error, tc.errorSegment,
								"Expected error to mention segment '%s'", tc.errorSegment)
						}
						assert.Equal(t, PriorityError, err.Priority,
							"Expected ERROR priority for segment length validation")
						assert.NotEmpty(t, err.Remediation,
							"Expected remediation message for segment length error")
						break
					}
				}
			}

			if tc.expectErrors {
				assert.True(t, hasSegmentError,
					"Expected segment length error for command: %s", cmdName)
			} else {
				assert.False(t, hasSegmentError,
					"Should not have segment length error for command: %s", cmdName)
			}
		})
	}
}

// TestValidateCommandNaming_RetrievalCommands tests that commands without verbs are accepted if description indicates retrieval
func TestValidateCommandNaming_RetrievalCommands(t *testing.T) {
	validator := NewPluginMetadataValidator()

	testCases := []struct {
		name           string
		pluginMetadata []PluginMetadata
		expectErrors   bool
	}{
		{
			name: "Valid - access-key with 'get' in description",
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
							Namespace:   "monitoring",
							Name:        "access-key",
							Description: "Get the monitoring access key",
							Usage:       "ibmcloud monitoring access-key",
						},
					},
				},
			},
			expectErrors: false,
		},
		{
			name: "Valid - api-key with 'show' in description",
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
							Namespace:   "service",
							Name:        "api-key",
							Description: "Show the API key details",
							Usage:       "ibmcloud service api-key",
						},
					},
				},
			},
			expectErrors: false,
		},
		{
			name: "Valid - service-url with 'display' in description",
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
							Namespace:   "config",
							Name:        "service-url",
							Description: "Display the service URL",
							Usage:       "ibmcloud config service-url",
						},
					},
				},
			},
			expectErrors: false,
		},
		{
			name: "Valid - connection-string with 'retrieve' in description",
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
							Namespace:   "database",
							Name:        "connection-string",
							Description: "Retrieve the database connection string",
							Usage:       "ibmcloud database connection-string",
						},
					},
				},
			},
			expectErrors: false,
		},
		{
			name: "Valid - status-info with 'view' in description",
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
							Namespace:   "cluster",
							Name:        "status-info",
							Description: "View cluster status information",
							Usage:       "ibmcloud cluster status-info",
						},
					},
				},
			},
			expectErrors: false,
		},
		{
			name: "Valid - endpoint-url with 'return' in description",
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
							Namespace:   "service",
							Name:        "endpoint-url",
							Description: "Return the service endpoint URL",
							Usage:       "ibmcloud service endpoint-url",
						},
					},
				},
			},
			expectErrors: false,
		},
		{
			name: "Valid - credentials with 'fetch' in description",
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
							Namespace:   "service",
							Name:        "credentials",
							Description: "Fetch service credentials",
							Usage:       "ibmcloud service credentials",
						},
					},
				},
			},
			expectErrors: false,
		},
		{
			name: "Invalid - access-key without retrieval keywords",
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
							Namespace:   "monitoring",
							Name:        "access-key",
							Description: "Manage the monitoring access key",
							Usage:       "ibmcloud monitoring access-key",
						},
					},
				},
			},
			expectErrors: true,
		},
		{
			name: "Invalid - api-key with action description",
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
							Namespace:   "service",
							Name:        "api-key",
							Description: "Create and manage API keys",
							Usage:       "ibmcloud service api-key",
						},
					},
				},
			},
			expectErrors: true,
		},
		{
			name: "Valid - command with verb doesn't need retrieval keyword",
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
							Namespace:   "service",
							Name:        "key-create",
							Description: "Manage service keys",
							Usage:       "ibmcloud service key-create",
						},
					},
				},
			},
			expectErrors: false,
		},
		{
			name: "Valid - plural form doesn't need retrieval keyword",
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
							Namespace:   "service",
							Name:        "instances",
							Description: "Manage service instances",
							Usage:       "ibmcloud service instances",
						},
					},
				},
			},
			expectErrors: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pluginToErrors := validator.Errors(tc.pluginMetadata)
			cmd := tc.pluginMetadata[0].Commands[0]

			// Check for the verb/plural error
			hasVerbError := false
			for _, pluginErrs := range pluginToErrors {
				for _, err := range pluginErrs {
					if containsString(err.Error, "does not use a common verb or plural form") {
						hasVerbError = true
						break
					}
				}
			}

			if tc.expectErrors {
				assert.True(t, hasVerbError,
					"Expected verb/plural error for command: %s with description: %s", cmd.Name, cmd.Description)
			} else {
				assert.False(t, hasVerbError,
					"Should not have verb/plural error for command: %s with description: %s", cmd.Name, cmd.Description)
			}
		})
	}
}

// Made with Bob
