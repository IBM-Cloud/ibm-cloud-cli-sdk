package plugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestValidateUsageEnhanced_EmptyUsage tests that empty usage is detected
func TestValidateUsageEnhanced_EmptyUsage(t *testing.T) {
	validator := NewPluginMetadataValidator()

	cmd := Command{
		Namespace:   "test",
		Name:        "command",
		Description: "Test command",
		Usage:       "",
	}

	errs := validator.validateUsageEnhanced(cmd, "test-plugin")

	assert.Len(t, errs, 1, "Expected 1 error for empty usage")
	assert.Equal(t, PriorityError, errs[0].Priority)
	assert.Contains(t, errs[0].Error, "has no usage information")
}

// TestValidateUsageEnhanced_ValidUsage tests valid usage patterns
func TestValidateUsageEnhanced_ValidUsage(t *testing.T) {
	validator := NewPluginMetadataValidator()

	testCases := []struct {
		name      string
		namespace string
		cmdName   string
		usage     string
	}{
		{
			name:      "Simple command with ibmcloud prefix",
			namespace: "plugin",
			cmdName:   "list",
			usage:     "ibmcloud plugin list",
		},
		{
			name:      "Command with uppercase arguments",
			namespace: "resource",
			cmdName:   "service-instance-create",
			usage:     "ibmcloud resource service-instance-create NAME SERVICE_NAME PLAN_NAME LOCATION",
		},
		{
			name:      "Command with optional uppercase arguments",
			namespace: "resource",
			cmdName:   "service-instance-update",
			usage:     "ibmcloud resource service-instance-update NAME [--service-plan-id PLAN_ID]",
		},
		{
			name:      "Command with multiple uppercase arguments",
			namespace: "ks cluster create",
			cmdName:   "classic",
			usage:     "ibmcloud ks cluster create classic --name NAME --zone ZONE [--workers COUNT]",
		},
		{
			name:      "Command with full path to ibmcloud",
			namespace: "plugin",
			cmdName:   "list",
			usage:     "/usr/local/bin/ibmcloud plugin list",
		},
		{
			name:      "Command with choices using uppercase",
			namespace: "account",
			cmdName:   "org-role-set",
			usage:     "ibmcloud account org-role-set USER_NAME ORG_NAME (ROLE_A | ROLE_B)",
		},
		{
			name:      "Complex command with mixed elements",
			namespace: "service",
			cmdName:   "create",
			usage:     "ibmcloud service create NAME PLAN LOCATION [--tags TAGS] [--parameters PARAMS]",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := Command{
				Namespace:   tc.namespace,
				Name:        tc.cmdName,
				Description: "Test command",
				Usage:       tc.usage,
			}

			errs := validator.validateUsageEnhanced(cmd, "test-plugin")

			assert.Empty(t, errs, "Expected no errors for valid usage: %s", tc.usage)
		})
	}
}

// TestValidateUsageEnhanced_MissingIbmcloudPrefix tests detection of missing ibmcloud prefix
func TestValidateUsageEnhanced_MissingIbmcloudPrefix(t *testing.T) {
	validator := NewPluginMetadataValidator()

	testCases := []struct {
		name  string
		usage string
	}{
		{
			name:  "No prefix at all",
			usage: "plugin list",
		},
		{
			name:  "Wrong prefix - uppercase",
			usage: "IBMCLOUD plugin list",
		},
		{
			name:  "Wrong prefix - mixed case",
			usage: "IBMCloud plugin list",
		},
		{
			name:  "Wrong prefix - ic only",
			usage: "ic plugin list",
		},
		{
			name:  "Command name only",
			usage: "list",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := Command{
				Namespace:   "plugin",
				Name:        "list",
				Description: "Test command",
				Usage:       tc.usage,
			}

			errs := validator.validateUsageEnhanced(cmd, "test-plugin")

			assert.NotEmpty(t, errs, "Expected error for usage without proper ibmcloud prefix: %s", tc.usage)
			found := false
			for _, err := range errs {
				if err.Priority == PriorityError &&
					(containsString(err.Error, "should start with 'ibmcloud'") ||
						containsString(err.Error, "full path")) {
					found = true
					break
				}
			}
			assert.True(t, found, "Expected error about ibmcloud prefix for: %s", tc.usage)
		})
	}
}

// TestValidateUsageEnhanced_LowercaseArguments tests detection of lowercase argument values
func TestValidateUsageEnhanced_LowercaseArguments(t *testing.T) {
	validator := NewPluginMetadataValidator()

	testCases := []struct {
		name              string
		namespace         string
		cmdName           string
		usage             string
		expectedLowercase []string
	}{
		{
			name:              "Single lowercase argument",
			namespace:         "resource",
			cmdName:           "service-instance-create",
			usage:             "ibmcloud resource service-instance-create name SERVICE PLAN LOCATION",
			expectedLowercase: []string{"name"},
		},
		{
			name:              "Multiple lowercase arguments - space separated",
			namespace:         "ks",
			cmdName:           "cluster create",
			usage:             "ibmcloud ks cluster create classic --name instance --zone region",
			expectedLowercase: []string{"instance"},
		},
		{
			name:              "Lowercase with hyphens",
			namespace:         "resource",
			cmdName:           "group-create",
			usage:             "ibmcloud resource group-create instance-name",
			expectedLowercase: []string{"instance-name"},
		},
		{
			name:              "Lowercase with underscores",
			namespace:         "resource",
			cmdName:           "service-key-create",
			usage:             "ibmcloud resource service-key-create instance_id KEY_NAME ROLE",
			expectedLowercase: []string{"instance_id"},
		},
		{
			name:              "Mixed case - some correct, some not",
			namespace:         "account",
			cmdName:           "user-invite",
			usage:             "ibmcloud account user-invite EMAIL --org ORG --space space_name",
			expectedLowercase: []string{"space_name"},
		},
		{
			name:              "Lowercase before pipe with space",
			namespace:         "account",
			cmdName:           "org-role-set",
			usage:             "ibmcloud account org-role-set USER ORG (option_a | OPTION_B)",
			expectedLowercase: []string{"option_a"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := Command{
				Namespace:   tc.namespace,
				Name:        tc.cmdName,
				Description: "Test command",
				Usage:       tc.usage,
			}

			errs := validator.validateUsageEnhanced(cmd, "test-plugin")

			assert.NotEmpty(t, errs, "Expected error for lowercase arguments in: %s", tc.usage)

			// Find the lowercase argument error
			found := false
			for _, err := range errs {
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
			assert.True(t, found, "Expected warning about lowercase arguments for: %s", tc.usage)
		})
	}
}

// TestValidateUsageEnhanced_ExcludedWords tests that command words are not flagged as lowercase arguments
func TestValidateUsageEnhanced_ExcludedWords(t *testing.T) {
	validator := NewPluginMetadataValidator()

	testCases := []struct {
		name      string
		namespace string
		cmdName   string
		usage     string
		shouldErr bool
	}{
		{
			name:      "Command name in usage",
			namespace: "plugin",
			cmdName:   "list",
			usage:     "ibmcloud plugin list",
			shouldErr: false,
		},
		{
			name:      "Multi-word command with hyphens",
			namespace: "plugin",
			cmdName:   "list-all",
			usage:     "ibmcloud plugin list all",
			shouldErr: false,
		},
		{
			name:      "Common words like 'and', 'or', 'to'",
			namespace: "plugin",
			cmdName:   "convert",
			usage:     "ibmcloud plugin convert from FORMAT to FORMAT",
			shouldErr: false,
		},
		{
			name:      "Namespace words are excluded",
			namespace: "service",
			cmdName:   "create",
			usage:     "ibmcloud service create NAME",
			shouldErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := Command{
				Namespace:   tc.namespace,
				Name:        tc.cmdName,
				Description: "Test command",
				Usage:       tc.usage,
			}

			errs := validator.validateUsageEnhanced(cmd, "test-plugin")

			// Check if lowercase argument errors exist
			hasLowercaseErr := false
			for _, err := range errs {
				if containsString(err.Error, "lowercase argument values") {
					hasLowercaseErr = true
					if !tc.shouldErr {
						t.Errorf("Should not flag command words as lowercase arguments: %s", err.Error)
					}
				}
			}

			if tc.shouldErr && !hasLowercaseErr {
				t.Errorf("Expected lowercase argument error but got none")
			}
		})
	}
}

// TestValidateUsageEnhanced_ComplexScenarios tests complex real-world usage patterns
func TestValidateUsageEnhanced_ComplexScenarios(t *testing.T) {
	validator := NewPluginMetadataValidator()

	testCases := []struct {
		name          string
		namespace     string
		cmdName       string
		usage         string
		expectErrors  bool
		errorContains string
	}{
		{
			name:         "Valid complex command",
			namespace:    "resource",
			cmdName:      "service-instance-create",
			usage:        "ibmcloud resource service-instance-create NAME SERVICE PLAN LOCATION [--tags TAGS] [--parameters PARAMS]",
			expectErrors: false,
		},
		{
			name:          "Complex with lowercase params",
			namespace:     "resource",
			cmdName:       "service-instance-create",
			usage:         "ibmcloud resource service-instance-create name SERVICE PLAN LOCATION",
			expectErrors:  true,
			errorContains: "lowercase argument values",
		},
		{
			name:         "Valid with JSON example (uppercase)",
			namespace:    "resource",
			cmdName:      "service-instance-update",
			usage:        "ibmcloud resource service-instance-update NAME --parameters DATA",
			expectErrors: false,
		},
		{
			name:         "Valid with multiple optional groups",
			namespace:    "ks cluster create",
			cmdName:      "classic",
			usage:        "ibmcloud ks cluster create classic --name NAME [--zone ZONE] [--workers COUNT] [--force]",
			expectErrors: false,
		},
		{
			name:          "Missing ibmcloud prefix with lowercase",
			namespace:     "resource",
			cmdName:       "service-instance-create",
			usage:         "resource service-instance-create instance",
			expectErrors:  true,
			errorContains: "should start with 'ibmcloud'",
		},
		{
			name:         "Valid with choice operator",
			namespace:    "account",
			cmdName:      "org-role-set",
			usage:        "ibmcloud account org-role-set USER ORG (ROLE_A | ROLE_B | ROLE_C)",
			expectErrors: false,
		},
		{
			name:          "Invalid choice with lowercase",
			namespace:     "account",
			cmdName:       "org-role-set",
			usage:         "ibmcloud account org-role-set USER ORG (option_a | option_b)",
			expectErrors:  true,
			errorContains: "lowercase argument values",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := Command{
				Namespace:   tc.namespace,
				Name:        tc.cmdName,
				Description: "Test command",
				Usage:       tc.usage,
			}

			errs := validator.validateUsageEnhanced(cmd, "test-plugin")

			if tc.expectErrors {
				assert.NotEmpty(t, errs, "Expected errors for: %s", tc.usage)
				if tc.errorContains != "" {
					found := false
					for _, err := range errs {
						if containsString(err.Error, tc.errorContains) {
							found = true
							break
						}
					}
					assert.True(t, found, "Expected error containing '%s' for: %s",
						tc.errorContains, tc.usage)
				}
			} else {
				assert.Empty(t, errs, "Expected no errors for: %s", tc.usage)
			}
		})
	}
}

// TestValidateUsageEnhanced_EdgeCases tests edge cases and boundary conditions
func TestValidateUsageEnhanced_EdgeCases(t *testing.T) {
	validator := NewPluginMetadataValidator()

	testCases := []struct {
		name         string
		namespace    string
		cmdName      string
		usage        string
		expectErrors bool
	}{
		{
			name:         "Very long usage string",
			namespace:    "plugin",
			cmdName:      "command",
			usage:        "ibmcloud plugin command PARAM1 PARAM2 PARAM3 PARAM4 PARAM5 --opt1 VAL1 --opt2 VAL2 --opt3 VAL3",
			expectErrors: false,
		},
		{
			name:         "Usage with only whitespace after ibmcloud",
			namespace:    "plugin",
			cmdName:      "list",
			usage:        "ibmcloud   ",
			expectErrors: false,
		},
		{
			name:         "Usage with tabs and newlines",
			namespace:    "plugin",
			cmdName:      "list",
			usage:        "ibmcloud\tplugin\nlist",
			expectErrors: false,
		},
		{
			name:         "Usage with special characters in uppercase args",
			namespace:    "plugin",
			cmdName:      "create",
			usage:        "ibmcloud plugin create NAME_WITH_UNDERSCORES",
			expectErrors: false,
		},
		{
			name:         "Usage with numbers in uppercase args",
			namespace:    "plugin",
			cmdName:      "create",
			usage:        "ibmcloud plugin create INSTANCE123",
			expectErrors: false,
		},
		{
			name:         "Single character lowercase param not detected",
			namespace:    "resource",
			cmdName:      "group-create",
			usage:        "ibmcloud resource group-create x",
			expectErrors: false, // Single char params are not caught by the regex
		},
		{
			name:         "Lowercase param at end of line",
			namespace:    "resource",
			cmdName:      "service-instance-create",
			usage:        "ibmcloud resource service-instance-create name",
			expectErrors: true,
		},
		{
			name:         "Lowercase param before pipe without space IS detected",
			namespace:    "account",
			cmdName:      "org-role-set",
			usage:        "ibmcloud account org-role-set USER ORG (option_a|OPTION_B)",
			expectErrors: true, // Improved regex now catches this case
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := Command{
				Namespace:   tc.namespace,
				Name:        tc.cmdName,
				Description: "Test command",
				Usage:       tc.usage,
			}

			errs := validator.validateUsageEnhanced(cmd, "test-plugin")

			if tc.expectErrors {
				assert.NotEmpty(t, errs, "Expected errors for: %s", tc.usage)
			} else {
				assert.Empty(t, errs, "Expected no errors for: %s", tc.usage)
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
		expectedPriority Priority
		errorType        string
	}{
		{
			name:             "Empty usage - ERROR priority",
			usage:            "",
			expectedPriority: PriorityError,
			errorType:        "no usage information",
		},
		{
			name:             "Missing ibmcloud prefix - ERROR priority",
			usage:            "plugin list",
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
			cmd := Command{
				Namespace:   "plugin",
				Name:        "command",
				Description: "Test command",
				Usage:       tc.usage,
			}

			errs := validator.validateUsageEnhanced(cmd, "test-plugin")

			assert.NotEmpty(t, errs, "Expected errors for: %s", tc.name)

			found := false
			for _, err := range errs {
				if containsString(err.Error, tc.errorType) {
					assert.Equal(t, tc.expectedPriority, err.Priority,
						"Expected priority %s for error type '%s'", tc.expectedPriority, tc.errorType)
					found = true
					break
				}
			}
			assert.True(t, found, "Expected to find error type '%s'", tc.errorType)
		})
	}
}

// TestValidateUsageEnhanced_Remediation tests that remediation messages are provided
func TestValidateUsageEnhanced_Remediation(t *testing.T) {
	validator := NewPluginMetadataValidator()

	testCases := []struct {
		name  string
		usage string
	}{
		{
			name:  "Empty usage",
			usage: "",
		},
		{
			name:  "Missing ibmcloud prefix",
			usage: "plugin list",
		},
		{
			name:  "Lowercase arguments",
			usage: "ibmcloud plugin create name",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := Command{
				Namespace:   "plugin",
				Name:        "command",
				Description: "Test command",
				Usage:       tc.usage,
			}

			errs := validator.validateUsageEnhanced(cmd, "test-plugin")

			assert.NotEmpty(t, errs, "Expected errors for: %s", tc.name)

			for _, err := range errs {
				assert.NotEmpty(t, err.Remediation,
					"Expected remediation message for error: %s", err.Error)
			}
		})
	}
}

// TestValidateCommandNaming_PluralForms tests that commands ending with 's' are accepted as alternatives to 'list'
func TestValidateCommandNaming_PluralForms(t *testing.T) {
	validator := NewPluginMetadataValidator()

	testCases := []struct {
		name         string
		namespace    string
		cmdName      string
		expectErrors bool
	}{
		{
			name:         "Valid plural form - instances",
			namespace:    "service",
			cmdName:      "instances",
			expectErrors: false,
		},
		{
			name:         "Valid plural form - services",
			namespace:    "resource",
			cmdName:      "services",
			expectErrors: false,
		},
		{
			name:         "Valid plural form - clusters",
			namespace:    "ks",
			cmdName:      "clusters",
			expectErrors: false,
		},
		{
			name:         "Valid plural form - resources",
			namespace:    "account",
			cmdName:      "resources",
			expectErrors: false,
		},
		{
			name:         "Valid plural form with multi-word",
			namespace:    "service",
			cmdName:      "key instances",
			expectErrors: false,
		},
		{
			name:         "Valid - has common verb",
			namespace:    "service",
			cmdName:      "instance list",
			expectErrors: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := Command{
				Namespace:   tc.namespace,
				Name:        tc.cmdName,
				Description: "Test command",
				Usage:       "ibmcloud " + tc.namespace + " " + tc.cmdName,
			}

			errs := validator.validateCommandNaming(cmd, "test-plugin")

			// Check for the "does not use a common verb" error
			hasVerbError := false
			for _, err := range errs {
				if containsString(err.Error, "does not use a common verb") {
					hasVerbError = true
					break
				}
			}

			if tc.expectErrors {
				assert.True(t, hasVerbError, "Expected verb/plural error for: %s", tc.cmdName)
			} else {
				assert.False(t, hasVerbError, "Should not have verb/plural error for: %s", tc.cmdName)
			}
		})
	}
}

// TestValidateDescription_PluralFormWithoutListKeywords tests that plural commands require list-related descriptions
func TestValidateDescription_PluralFormWithoutListKeywords(t *testing.T) {
	validator := NewPluginMetadataValidator()

	testCases := []struct {
		name         string
		cmdName      string
		description  string
		expectError  bool
		errorMessage string
	}{
		{
			name:        "Plural command with 'list' keyword - valid",
			cmdName:     "instances",
			description: "List all service instances",
			expectError: false,
		},
		{
			name:        "Plural command with 'show' keyword - valid",
			cmdName:     "services",
			description: "Show available services",
			expectError: false,
		},
		{
			name:        "Plural command with 'display' keyword - valid",
			cmdName:     "clusters",
			description: "Display all clusters in the account",
			expectError: false,
		},
		{
			name:        "Plural command with 'view' keyword - valid",
			cmdName:     "resources",
			description: "View all resources",
			expectError: false,
		},
		{
			name:        "Plural command with 'retrieve' keyword - valid",
			cmdName:     "policies",
			description: "Retrieve security policies",
			expectError: false,
		},
		{
			name:        "Plural command with 'get' keyword - valid",
			cmdName:     "users",
			description: "Get all users in the organization",
			expectError: false,
		},
		{
			name:        "Plural command with 'all' keyword - valid",
			cmdName:     "instances",
			description: "Returns all instances",
			expectError: false,
		},
		{
			name:        "Plural command with 'multiple' keyword - valid",
			cmdName:     "services",
			description: "Manage multiple services",
			expectError: false,
		},
		{
			name:         "Plural command without list keywords - invalid",
			cmdName:      "instances",
			description:  "Manage service instances",
			expectError:  true,
			errorMessage: "doesn't clearly indicate it returns a list or group of items",
		},
		{
			name:         "Plural command with vague description - invalid",
			cmdName:      "services",
			description:  "Work with services",
			expectError:  true,
			errorMessage: "doesn't clearly indicate it returns a list or group of items",
		},
		{
			name:         "Plural command with action description - invalid",
			cmdName:      "clusters",
			description:  "Create and manage clusters",
			expectError:  true,
			errorMessage: "doesn't clearly indicate it returns a list or group of items",
		},
		{
			name:        "Non-plural command without list keywords - valid (no check)",
			cmdName:     "create",
			description: "Create a new instance",
			expectError: false,
		},
		{
			name:        "Short plural word 'is' not flagged",
			cmdName:     "is",
			description: "Check status",
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := Command{
				Namespace:   "service",
				Name:        tc.cmdName,
				Description: tc.description,
				Usage:       "ibmcloud service " + tc.cmdName,
			}

			errs := validator.validateDescription(cmd, "test-plugin")

			// Check for the plural form description error
			hasPluralDescError := false
			for _, err := range errs {
				if containsString(err.Error, "doesn't clearly indicate it returns a list or group of items") {
					hasPluralDescError = true
					if tc.errorMessage != "" {
						assert.Contains(t, err.Error, tc.errorMessage)
					}
					assert.Equal(t, PriorityWarning, err.Priority, "Expected WARNING priority for plural description check")
					break
				}
			}

			if tc.expectError {
				assert.True(t, hasPluralDescError, "Expected plural description error for: %s with description: %s", tc.cmdName, tc.description)
			} else {
				assert.False(t, hasPluralDescError, "Should not have plural description error for: %s with description: %s", tc.cmdName, tc.description)
			}
		})
	}
}

// TestValidateCommandNaming_ConfigPattern tests that "config get/set/unset <parameter>" commands are accepted
func TestValidateCommandNaming_ConfigPattern(t *testing.T) {
	validator := NewPluginMetadataValidator()

	testCases := []struct {
		name         string
		namespace    string
		cmdName      string
		expectErrors bool
	}{
		{
			name:         "Valid config get pattern",
			namespace:    "backup-recovery config get",
			cmdName:      "connector-service-url",
			expectErrors: false,
		},
		{
			name:         "Valid config set pattern",
			namespace:    "app-configuration config set",
			cmdName:      "service-url",
			expectErrors: false,
		},
		{
			name:         "Valid config unset pattern",
			namespace:    "service config unset",
			cmdName:      "api-key",
			expectErrors: false,
		},
		{
			name:         "Valid config get with hyphenated parameter",
			namespace:    "plugin config get",
			cmdName:      "max-retry-count",
			expectErrors: false,
		},
		{
			name:         "Valid config set with underscored parameter",
			namespace:    "plugin config set",
			cmdName:      "log_level",
			expectErrors: false,
		},
		{
			name:         "Valid config get with uppercase in namespace",
			namespace:    "Cloud-Logs Config Get",
			cmdName:      "service-url",
			expectErrors: false,
		},
		{
			name:         "Invalid - namespace doesn't end with config verb",
			namespace:    "service config",
			cmdName:      "parameter-name",
			expectErrors: true,
		},
		{
			name:         "Invalid - namespace has config with wrong verb",
			namespace:    "service config something",
			cmdName:      "parameter-name",
			expectErrors: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := Command{
				Namespace:   tc.namespace,
				Name:        tc.cmdName,
				Description: "Test command for config pattern",
				Usage:       "ibmcloud " + tc.namespace + " " + tc.cmdName,
			}

			errs := validator.validateCommandNaming(cmd, "test-plugin")

			// Check for the "does not use a common verb" error
			hasVerbError := false
			for _, err := range errs {
				if containsString(err.Error, "does not use a common verb") {
					hasVerbError = true
					break
				}
			}

			if tc.expectErrors {
				assert.True(t, hasVerbError, "Expected verb/plural error for: %s", tc.cmdName)
			} else {
				assert.False(t, hasVerbError, "Should not have verb/plural error for: %s", tc.cmdName)
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
		name         string
		namespace    string
		cmdName      string
		expectErrors bool
		errorSegment string
	}{
		{
			name:         "Valid - single word with 2+ chars",
			namespace:    "plugin",
			cmdName:      "list",
			expectErrors: false,
		},
		{
			name:         "Valid - multi-word with all segments 2+ chars",
			namespace:    "service",
			cmdName:      "instance-create",
			expectErrors: false,
		},
		{
			name:         "Valid - three words all 2+ chars",
			namespace:    "ks",
			cmdName:      "cluster worker update",
			expectErrors: false,
		},
		{
			name:         "Invalid - single letter command",
			namespace:    "plugin",
			cmdName:      "a",
			expectErrors: true,
			errorSegment: "a",
		},
		{
			name:         "Invalid - hyphenated with single letter segment",
			namespace:    "service",
			cmdName:      "a-b",
			expectErrors: true,
			errorSegment: "a",
		},
		{
			name:         "Invalid - space separated with single letter",
			namespace:    "plugin",
			cmdName:      "list x",
			expectErrors: true,
			errorSegment: "x",
		},
		{
			name:         "Invalid - multiple single letter segments",
			namespace:    "service",
			cmdName:      "a-b-c",
			expectErrors: true,
			errorSegment: "a",
		},
		{
			name:         "Invalid - valid word followed by single letter",
			namespace:    "service",
			cmdName:      "create-x",
			expectErrors: true,
			errorSegment: "x",
		},
		{
			name:         "Invalid - single letter followed by valid word",
			namespace:    "service",
			cmdName:      "x-create",
			expectErrors: true,
			errorSegment: "x",
		},
		{
			name:         "Invalid - command name with single letter segment",
			namespace:    "ks cluster",
			cmdName:      "x-update",
			expectErrors: true,
			errorSegment: "x",
		},
		{
			name:         "Valid - two character segments",
			namespace:    "service",
			cmdName:      "ls-rm",
			expectErrors: false,
		},
		{
			name:         "Valid - namespace with single letter not checked",
			namespace:    "x",
			cmdName:      "list",
			expectErrors: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := Command{
				Namespace:   tc.namespace,
				Name:        tc.cmdName,
				Description: "Test command",
				Usage:       "ibmcloud " + tc.namespace + " " + tc.cmdName,
			}

			errs := validator.validateCommandNaming(cmd, "test-plugin")

			// Check for the segment length error
			hasSegmentError := false
			for _, err := range errs {
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

			if tc.expectErrors {
				assert.True(t, hasSegmentError,
					"Expected segment length error for command: %s", tc.cmdName)
			} else {
				assert.False(t, hasSegmentError,
					"Should not have segment length error for command: %s", tc.cmdName)
			}
		})
	}
}

// TestValidateCommandNaming_RetrievalCommands tests that commands without verbs are accepted if description indicates retrieval
func TestValidateCommandNaming_RetrievalCommands(t *testing.T) {
	validator := NewPluginMetadataValidator()

	testCases := []struct {
		name         string
		namespace    string
		cmdName      string
		description  string
		expectErrors bool
	}{
		{
			name:         "Valid - access-key with 'get' in description",
			namespace:    "monitoring",
			cmdName:      "access-key",
			description:  "Get the monitoring access key",
			expectErrors: false,
		},
		{
			name:         "Valid - api-key with 'show' in description",
			namespace:    "service",
			cmdName:      "api-key",
			description:  "Show the API key details",
			expectErrors: false,
		},
		{
			name:         "Valid - service-url with 'display' in description",
			namespace:    "config",
			cmdName:      "service-url",
			description:  "Display the service URL",
			expectErrors: false,
		},
		{
			name:         "Valid - connection-string with 'retrieve' in description",
			namespace:    "database",
			cmdName:      "connection-string",
			description:  "Retrieve the database connection string",
			expectErrors: false,
		},
		{
			name:         "Valid - status-info with 'view' in description",
			namespace:    "cluster",
			cmdName:      "status-info",
			description:  "View cluster status information",
			expectErrors: false,
		},
		{
			name:         "Valid - endpoint-url with 'return' in description",
			namespace:    "service",
			cmdName:      "endpoint-url",
			description:  "Return the service endpoint URL",
			expectErrors: false,
		},
		{
			name:         "Valid - credentials with 'fetch' in description",
			namespace:    "service",
			cmdName:      "credentials",
			description:  "Fetch service credentials",
			expectErrors: false,
		},
		{
			name:         "Invalid - access-key without retrieval keywords",
			namespace:    "monitoring",
			cmdName:      "access-key",
			description:  "Manage the monitoring access key",
			expectErrors: true,
		},
		{
			name:         "Invalid - api-key with action description",
			namespace:    "service",
			cmdName:      "api-key",
			description:  "Create and manage API keys",
			expectErrors: true,
		},
		{
			name:         "Valid - command with verb doesn't need retrieval keyword",
			namespace:    "service",
			cmdName:      "key-create",
			description:  "Manage service keys",
			expectErrors: false,
		},
		{
			name:         "Valid - plural form doesn't need retrieval keyword",
			namespace:    "service",
			cmdName:      "instances",
			description:  "Manage service instances",
			expectErrors: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := Command{
				Namespace:   tc.namespace,
				Name:        tc.cmdName,
				Description: tc.description,
				Usage:       "ibmcloud " + tc.namespace + " " + tc.cmdName,
			}

			errs := validator.validateCommandNaming(cmd, "test-plugin")

			// Check for the verb/plural error
			hasVerbError := false
			for _, err := range errs {
				if containsString(err.Error, "does not use a common verb or plural form") {
					hasVerbError = true
					break
				}
			}

			if tc.expectErrors {
				assert.True(t, hasVerbError,
					"Expected verb/plural error for command: %s with description: %s", tc.cmdName, tc.description)
			} else {
				assert.False(t, hasVerbError,
					"Should not have verb/plural error for command: %s with description: %s", tc.cmdName, tc.description)
			}
		})
	}
}

// Made with Bob
