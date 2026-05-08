package plugin

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"unicode"

	"github.com/Masterminds/semver"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/i18n"
	"github.com/go-playground/validator/v10"
)

// Priority represents the severity level of a validation error
type Priority string

const (
	// PriorityError - High severity deviation from standards
	PriorityError Priority = "ERROR"
	// PriorityWarning - Recommended improvements - style and convention issues
	PriorityWarning Priority = "WARNING"
	// PriorityInfo - Nice-to-have suggestions - polish and enhancements
	PriorityInfo Priority = "INFO"
)

var (
	validate  *validator.Validate
	validErrs validator.ValidationErrors

	cmdIndxRegex       *regexp.Regexp
	placeholdersRegexp *regexp.Regexp
	idxRegex           *regexp.Regexp

	compileOnce sync.Once
)

// Common verbs that should be used in commands
var commonVerbs = map[string]bool{
	"list": true, "show": true, "create": true, "update": true, "delete": true, "remove": true,
	"add": true, "bind": true, "unbind": true, "get": true, "set": true, "unset": true,
	"enable": true, "disable": true, "start": true, "stop": true, "restart": true, "config": true, "download": true,
}

// Maximum command depth (excluding base CLI name like "ibmcloud")
const maxCommandDepth = 3

// lastWordIsPluralForm checks if the last word in a command name is a plural form
// Returns true if the last word ends with 's', is at least 3 characters long,
// and is not a common verb (to avoid false positives like "is", "as", "list", "set")
func lastWordIsPluralForm(cmdName string) bool {
	words := strings.Fields(strings.ReplaceAll(cmdName, "-", " "))
	if len(words) == 0 {
		return false
	}
	lastWord := words[len(words)-1]
	return len(lastWord) >= 3 && strings.HasSuffix(lastWord, "s") && !commonVerbs[lastWord]
}

// Maximum recommended word count for command descriptions
const maxDescriptionWordCount = 15

// MinimumRequiredCLIVersion is the minimum CLI version that plugins must support
// This should match the validation tag on PluginMetadata.MinCliVersion
const MinimumRequiredCLIVersion = "2.0.0"

// PluginToValidationErrors is a mapping between a plugin and one or more metadata validation errors
type PluginToValidationErrors map[string][]PluginMetadataError

type PluginMetadataError struct {
	Namespace   string   `json:"namespace"`
	CommandName string   `json:"command,omitempty"`
	Error       string   `json:"error"`
	Priority    Priority `json:"priority"`
	Remediation string   `json:"remediation,omitempty"`
}

type pluginMetadataValidate struct {
	validator *validator.Validate
}

type PluginMetadataValidate interface {
	Errors(metadatas []PluginMetadata) PluginToValidationErrors
}

// Validator returns an singleton of validator
func Validator() *validator.Validate {
	if validate == nil {
		validate = validator.New()
	}

	validate.RegisterValidation("usage", func(fl validator.FieldLevel) bool {
		return validatePluginMetadataUsage(fl.Field().String()) == nil
	})

	validate.RegisterValidation("mincliversion", func(fl validator.FieldLevel) bool {
		return validateCliVersionMinimum(fl.Field().Interface().(VersionType), fl.Param()) == nil
	})

	validate.RegisterStructValidation(validatePluginCommandStruct, Command{})

	return validate
}

func NewPluginMetadataValidator() *pluginMetadataValidate {
	return &pluginMetadataValidate{
		validator: Validator(),
	}
}

// Errors returns a list of validation errors found for each plug-in metadata provided
func (p pluginMetadataValidate) Errors(metadatas []PluginMetadata) PluginToValidationErrors {

	// PERF: compile regular expression once since we are using this over many plugins
	compileOnce.Do(func() {
		placeholdersRegexp = regexp.MustCompile(`\[([aA]rguments|[oO]ptions)+\.{0,3}\]`)
		cmdIndxRegex = regexp.MustCompile(`Commands\[\d+\]`)
		idxRegex = regexp.MustCompile(`\d+`)
	})

	var (
		metadataErrs []PluginMetadataError
		pluginName   string
		cmdName      string
	)
	pluginToMetadataErrs := PluginToValidationErrors{}
	for _, metadata := range metadatas {
		metadataErrs = []PluginMetadataError{}

		// Validate plugin-level metadata
		metadataErrs = append(metadataErrs, p.validatePluginInfo(metadata)...)

		// Run struct validation
		if err := p.validator.Struct(metadata); err != nil {
			validErr := PluginMetadataError{}
			if errors.As(err, &validErrs) {
				for _, v := range validErrs {
					validErr.Namespace = v.StructNamespace()

					// attempt to include the CommandName if available
					cmdIdx := p.CommandIndexFromNamespace(v.StructNamespace())
					if cmdIdx != -1 {
						cmd := metadata.Commands[cmdIdx]
						cmdName = cmd.Name
						validErr.CommandName = cmd.Namespace + " " + cmd.Name
					}
					pluginErr := p.ParseValidationTagErrors(v, cmdName)
					validErr.Error = pluginErr.Error
					validErr.Priority = PriorityError // Default priority for struct validation errors
					if pluginErr.Priority != "" {
						validErr.Priority = pluginErr.Priority
					}
					validErr.Remediation = pluginErr.Remediation
					metadataErrs = append(metadataErrs, validErr)
				}
			}
		}

		// Only add to map if there are errors
		if len(metadataErrs) > 0 {
			// NOTE: Use "UNKNOWN" if the plugin name is not defined
			if metadata.Name != "" {
				pluginName = metadata.Name
			} else {
				pluginName = "UNKNOWN"
			}
			pluginToMetadataErrs[pluginName] = metadataErrs
		}
	}

	return pluginToMetadataErrs
}

// ParseValidationTagErrors translates validator.FieldError validation tags into user-friendly PluginMetadataError messages.
// This method serves as the central error message formatter for all plugin metadata validation failures,
// converting technical validation tags into actionable error messages with appropriate priority levels and remediation guidance.
//
// Parameters:
//   - fieldErr: The validation error from the go-playground/validator library containing the failed validation tag and context
//   - cmdName: The name of the command being validated (used for contextual error messages)
//
// Returns:
//   - PluginMetadataError: A structured error containing the error message, priority level, and remediation advice
//
// Supported Validation Tags:
//
//   - min: Validates minimum element count in collections (e.g., Commands, Namespaces must have at least 1 element)
//     Priority: ERROR
//
//   - required: Validates presence of mandatory fields (Name, Description, Usage)
//     Priority: ERROR
//     Special handling for Usage and Description fields with command-specific error messages
//
//   - uppercase: Validates that descriptions start with a capital letter
//     Priority: ERROR
//     Remediation: Capitalize the first letter
//
//   - cmdsegmin: Validates minimum character length for command name segments
//     Priority: ERROR
//     Format: "count|segment" where count is minimum length and segment is the offending word
//     Remediation: Use more descriptive words (at least 2 characters per segment)
//
//   - noreserves: Detects use of reserved command names (--version, --help, -v, -h)
//     Priority: WARNING
//     Remediation: Remove command as these are automatically provided by CLI framework
//
//   - maxnsdepth: Validates command depth does not exceed maximum (default: 3 levels)
//     Priority: WARNING
//     Remediation: Flatten hierarchy, use flags instead of subcommands, or reorganize structure
//
//   - usage: Validates command usage text meets requirements (proper prefix, no placeholders, closed brackets)
//     Priority: ERROR
//     Delegates to validatePluginMetadataUsage for detailed validation
//
//   - nosubject: Detects anti-pattern of starting descriptions with subjects like "This command", "Plugin to"
//     Priority: ERROR
//     Remediation: Remove subject and start directly with action verb
//
//   - mincliversion: Validates MinCliVersion meets minimum required version (2.0.0)
//     Priority: ERROR
//     Delegates to validateCliVersionMinimum for version comparison
//
//   - cmdwordcount: Validates description word count doesn't exceed recommended maximum (15 words)
//     Priority: INFO
//     Remediation: Shorten description for better display
//
//   - cleardesc: Validates that plural command names have descriptions indicating they return multiple items
//     Priority: WARNING
//     Remediation: Include words like 'list', 'show', 'display', 'all', or 'multiple'
//
//   - capargs: Validates that positional arguments in usage text are in CAPITAL letters
//     Priority: WARNING
//     Param contains comma-separated list of lowercase arguments found
//     Remediation: Convert to CAPITAL letters (e.g., NAME, INSTANCE_ID, FORMAT)
//
//   - noverb: Validates command names use common verbs or plural forms
//     Priority: ERROR (default, but can vary)
//     Remediation: Use common verbs (list, create, update, delete, show, get, set) or plural nouns
//
//   - excludesall: Validates field names don't contain forbidden characters
//     Priority: ERROR
//     Param contains the forbidden characters found
//
// Example Usage:
//
//	fieldErr := validator.FieldError{...}
//	pluginErr := validator.ParseValidationTagErrors(fieldErr, "service create")
//	// Returns: PluginMetadataError{
//	//   Error: "Command 'service create' has no usage information.",
//	//   Priority: PriorityError,
//	//   Remediation: "Add usage text showing command syntax with parameters and options."
//	// }
//
// Note: If no remediation message is explicitly set, it defaults to the error message.
// If no priority is set, it defaults to PriorityError.
func (p pluginMetadataValidate) ParseValidationTagErrors(fieldErr validator.FieldError, cmdName string) PluginMetadataError {
	var (
		errMsg         string
		remediationMsg string
		priority       Priority
	)
	switch fieldErr.Tag() {
	case "min":
		errMsg = i18n.T("{{.Field}} must contain at least {{.Param}} element", map[string]any{
			"Field": fieldErr.Field(),
			"Param": fieldErr.Param(),
		})
	case "required":
		if fieldErr.StructField() == "Usage" {
			remediationMsg = i18n.T("Add usage text showing command syntax with parameters and options.")
			errMsg = i18n.T("Command '{{.Name}}' has no usage information.", map[string]any{
				"Name": cmdName,
			})
		} else if fieldErr.StructField() == "Description" {
			if isCommandDescription(fieldErr.Namespace()) {
				remediationMsg = i18n.T("Add a sentence without subject describing what the command does.")
				errMsg = i18n.T("Command '{{.Name}}' has no description. All commands must have a clear description.", map[string]any{
					"Name": cmdName,
				})
			}
		} else {
			errMsg = i18n.T("{{.Field}} is required", map[string]any{
				"Field": fieldErr.StructField(),
			})
		}
	// Validate if property starts with a capital letter
	case "uppercase":

		errMsg = i18n.T("Description for '{{.Name}}' should start with a capital letter.", map[string]any{
			"Name": cmdName,
		})
		remediationMsg = i18n.T("Capitalize the first letter of the description.")

	// Validate if command name has at least X amount of segments
	case "cmdsegmin":
		params := strings.Split(fieldErr.Param(), "|")
		count := params[0]
		segment := params[1]

		errMsg = i18n.T("Command '{{.Name}}' contains a segment '{{.Segment}}' that is less than {{.Count}} characters. Each word in a command should be at least {{.Count}} characters.", map[string]any{
			"Name":    cmdName,
			"Segment": segment,
			"Count":   count,
		})
		remediationMsg = i18n.T("Use more descriptive words with at least {{.Count}} characters for each segment.", map[string]any{
			"Count": fieldErr.Param(),
		})
	case "noreserves":
		errMsg = i18n.T("Command '{{.Name}}' uses a reserved flag name. These are handled by the CLI framework.", map[string]any{
			"Name": cmdName,
		})
		priority = PriorityWarning
		remediationMsg = "Remove this command; --version and --help are automatically provided."
	case "maxnsdepth":

		errMsg = i18n.T("Command '{{.Name}}' has {{.Level}} levels, exceeding the maximum of {{.Level}}. Deep command hierarchies are difficult for users to remember and discover.",
			map[string]any{
				"Name":  cmdName,
				"Level": fieldErr.Param(),
			})

		priority = PriorityWarning
		remediationMsg = i18n.T("Reduce command depth to {{.Level}} or fewer levels. Options: (1) Flatten the hierarchy by combining levels, "+
			"(2) Use command flags/options instead of subcommands, (3) Reorganize the command structure to be more intuitive.", map[string]any{
			"Level": fieldErr.Param(),
		})
	case "usage":
		if err := validatePluginMetadataUsage(fmt.Sprint(fieldErr.Value())); err != nil {
			return *err
		}
	case "nosubject":
		errMsg = i18n.T("Description for '{{.Name}}' starts with '{{.Bad}}'. Use a sentence without subject.", map[string]any{
			"Name": cmdName,
			"Bad":  fieldErr.Param(),
		})
		remediationMsg = i18n.T("Remove '{{.Bad}}' and start directly with the action. Example: 'List all instances' instead of 'This command lists all instances'.", map[string]any{
			"Bad": fieldErr.Param(),
		})
	case "mincliversion":
		if err := validateCliVersionMinimum(fieldErr.Value().(VersionType), fieldErr.Param()); err != nil {
			return *err
		}
	case "excludesall":
		errMsg = i18n.T("{{.Field}} contains the following forbidden characters: {{.Chars}}", map[string]any{
			"Field": fieldErr.Field(),
			"Chars": strings.Join(strings.Split(fieldErr.Param(), ""), " "),
		})

	case "cmdwordcount":
		wordCount := len(strings.Fields(fieldErr.Value().(string)))
		errMsg = i18n.T("Description for '{{.Name}}' has {{.WordCount}} words. Consider limiting to less than {{.MaxWordCount}} words for better display.", map[string]any{
			"Name":         cmdName,
			"WordCount":    wordCount,
			"MaxWordCount": fieldErr.Param(),
		})
		priority = PriorityInfo
		remediationMsg = i18n.T("Shorten the description to be more concise.")

	case "cleardesc":
		errMsg = i18n.T("Command '{{.Name}}' uses plural form but description doesn't clearly indicate it returns a list or group of items.", map[string]any{
			"Name": cmdName,
		})
		priority = PriorityWarning
		remediationMsg = i18n.T("Update description to include words like 'list', 'show', 'display', 'view', 'all', or 'multiple' to clarify it returns multiple items.")

	case "capargs":
		errMsg = i18n.T("Command '{{.Name}}' usage contains lowercase argument values: {{.Args}}. User input values should be in CAPITAL letters.",
			map[string]any{
				"Name": cmdName,
				"Args": fieldErr.Param(),
			})

		priority = PriorityWarning
		remediationMsg = i18n.T("Convert positional arguments to CAPITAL letters (e.g., NAME, INSTANCE_ID, FORMAT).")
	case "noverb":
		errMsg = i18n.T("Command '{{.Name}}' does not use a common verb or plural form. Consider using verbs like: list, create, update, delete, show, get, set... or plural nouns like 'instances', 'services'", map[string]any{
			"Name": cmdName,
		})
		remediationMsg = i18n.T("Use common verbs in command names such as list, create, update, delete, or use plural forms to indicate listing operations.")
	}

	if remediationMsg == "" {
		remediationMsg = errMsg
	}

	if priority == "" {
		priority = PriorityError
	}

	return PluginMetadataError{
		Error:       errMsg,
		Remediation: remediationMsg,
		Priority:    priority,
	}
}

func UnclosedGroupings(usageText string) (bool, rune, []int) {
	var stack []rune
	var startIdxStack []int
	var indicies = make([]int, 2)
	closingToOpeningGroupMap := map[rune]rune{
		')': '(',
		']': '[',
		'>': '<',
		'}': '{',
	}

	for idx, char := range usageText {
		switch char {
		// add to stack if an open bracket or parenthesis is found
		// keep track of starting index of last occurrence of an opening bracket/parenthesis
		case '(', '[', '{', '<':
			stack = append(stack, char)
			startIdxStack = append(startIdxStack, idx)

		case ')', ']', '}', '>':
			expectedOpen := closingToOpeningGroupMap[char]
			if len(stack) > 0 {
				actualOpen := stack[len(stack)-1]
				// remove from top of stack is the expected opening (ie. "]" should expect "[")
				if actualOpen == expectedOpen {
					stack = stack[:len(stack)-1]
					startIdxStack = startIdxStack[:len(startIdxStack)-1]
				} else {
					// otherwise there is an error
					// return true and start/end indicies
					indicies[1] = idx + 1
					if len(startIdxStack) > 0 {
						indicies[0] = startIdxStack[len(startIdxStack)-1]
					}
					return true, actualOpen, indicies
				}
			}
		}

	}

	// if there still is an opening groupings in the stack report as error
	if len(stack) > 0 {
		indicies[0] = startIdxStack[len(startIdxStack)-1]
		indicies[1] = len(usageText) - 1
		return true, stack[0], indicies
	}

	return false, ' ', indicies
}

func validateCliVersionMinimum(versionType VersionType, minVersion string) *PluginMetadataError {
	minSemVer, minErr := semver.NewVersion(minVersion)
	if minErr != nil {
		return &PluginMetadataError{
			Error:    minErr.Error(),
			Priority: PriorityError,
		}
	}
	specifiedVersion, verErr := semver.NewVersion(versionType.String())
	if verErr != nil {
		return &PluginMetadataError{
			Error:    verErr.Error(),
			Priority: PriorityError,
		}
	}

	if specifiedVersion.LessThan(minSemVer) {
		return &PluginMetadataError{
			Error: i18n.T("MinCliVersion ({{.ProvidedMinVersion}}) is lower than the allowed minimum {{.AllowedMinimum}}", map[string]any{
				"ProvidedMinVersion": versionType.String(),
				"AllowedMinimum":     minVersion,
			}),
			Remediation: i18n.T("Set MinCliVersion to {{.Version}} or higher to ensure compatibility with supported CLI versions.", minVersion),
			Priority:    PriorityError,
		}
	}

	return nil
}

func validatePluginMetadataUsage(usageText string) *PluginMetadataError {

	// (1) All arguments and flags MUST be explicitly stated (no placeholders)
	if placeholdersRegexp.MatchString(usageText) {
		return &PluginMetadataError{
			Error:       i18n.T("Usage contains placeholder arguments/flags"),
			Priority:    PriorityError,
			Remediation: i18n.T("Remove placeholders from command usage text"),
		}

	}

	// (2) Check that usage starts with 'ibmcloud' or full path
	usageStripped := strings.TrimSpace(usageText)
	if !strings.HasPrefix(usageStripped, "ibmcloud") && !strings.Contains(usageStripped[:min(50, len(usageStripped))], "/ibmcloud") {
		return &PluginMetadataError{
			Error: i18n.T("Usage should start with '{{.Command}}' (lowercase) or the full path to the {{.Command}} binary.",
				map[string]any{
					"Command": "ibmcloud",
				}),
			Priority: PriorityError,
			Remediation: i18n.T("Start usage examples with '{{.Command}}' in lowercase (e.g., '{{.CommandPlugin}}...').",
				map[string]any{
					"Command":       "ibmcloud",
					"CommandPlugin": "ibmcloud plugin-name command",
				}),
		}
	}

	// (3) No basic usage text (Just COMMAND)
	parts := strings.SplitAfter(usageText, "COMMAND")
	if len(parts) == 2 && strings.TrimSpace(parts[1]) == "" {
		return &PluginMetadataError{
			Error:    i18n.T("Usage does not have any usage text besides COMMAND"),
			Priority: PriorityError,
		}

	}

	// (4) All surrounded parenthesis and brackets MUST be closed
	unclosed, unclosedChar, indicies := UnclosedGroupings(usageText)
	if unclosed {
		return &PluginMetadataError{
			Error: i18n.T("Usage contains unclosed {{.UnclosedGroup}} between indicies {{.Indicies}}", map[string]any{
				"Indicies":      indicies,
				"UnclosedGroup": string(unclosedChar),
			}),
			Priority: PriorityError,
		}
	}

	return nil
}

// CommandIndexFromNamespace will attempt to return the index from the namespace of the Command (eg. PluginMetadata.Commands[8].Usage)
// otherwise return -1
func (p pluginMetadataValidate) CommandIndexFromNamespace(namespace string) int {

	if !cmdIndxRegex.MatchString(namespace) {
		return -1
	}

	matchStr := cmdIndxRegex.FindString(namespace)
	matchIndx := idxRegex.FindString(matchStr)

	if matchIndx == "" {
		return -1
	}

	idx, err := strconv.Atoi(matchIndx)
	if err != nil {
		return -1
	}

	return idx
}

// validatePluginInfo validates plugin-level information
func (p pluginMetadataValidate) validatePluginInfo(metadata PluginMetadata) []PluginMetadataError {
	var errs []PluginMetadataError

	// Check plugin name for uppercase
	if metadata.Name != "" && hasUppercase(metadata.Name) {
		errs = append(errs, PluginMetadataError{
			Namespace:   "PluginMetadata.Name",
			Error:       fmt.Sprintf("Plugin name '%s' contains uppercase letters. Use lowercase with hyphens.", metadata.Name),
			Priority:    PriorityError,
			Remediation: "Convert plugin name to lowercase with hyphens (e.g., 'my-service').",
		})
	}

	return errs
}

func validatePluginCommandStruct(sl validator.StructLevel) {
	validateCommandNaming(sl)
	validateCommandDepth(sl)
	validateDescription(sl)
	validatePositionalArguments(sl)
}

func validatePositionalArguments(sl validator.StructLevel) {
	cmd := sl.Current().Interface().(Command)
	if cmd.Name == "" {
		sl.ReportError(cmd.Name, "name", "Name", "required", "")
		return
	}

	usageText := cmd.Usage
	// Check for lowercase positional argument values (should be CAPS)
	// Look for lowercase words that appear to be user-input parameters
	// Filter out common words, command names, and words that are part of the command structure
	// Strategy: Find lowercase words including those in choice operators, but exclude flags and paths

	// Remove flag names (words after --) to avoid false positives
	// Flags like --name, --zone, --tags or -a, -A should not be flagged
	usageWithoutFlags := regexp.MustCompile(`(-[a-zA-Z]{1}\s+)|(--[a-zA-Z][a-zA-Z-]*)`).ReplaceAllString(usageText, "")

	// Remove file paths to avoid flagging path components like /usr/local/bin
	usageWithoutPaths := regexp.MustCompile(`/[a-z/]+`).ReplaceAllString(usageWithoutFlags, "")

	// Match lowercase words (2+ chars, may contain hyphens/underscores)
	// Use word boundaries to properly match words
	// This will match lowercase words even inside choice operators like (option_a | OPTION_B)
	paramPattern := regexp.MustCompile(`\b([a-z][a-z_-]+)\b`)
	matches := paramPattern.FindAllStringSubmatch(usageWithoutPaths, -1)
	var lowercaseParams []string

	// Build list of words to exclude: common words + words from the command name itself
	excludeWords := map[string]bool{
		"and": true, "or": true, "to": true, "from": true, "with": true, "ibmcloud": true, "ic": true, "bx": true,
		"options": true, "arguments": true,
	}

	// Add words from the command name to exclusion list (these are command/subcommand names)
	// Split on spaces but keep hyphenated words intact (e.g., "service-instance-create" stays as one word)
	cmdWords := strings.Fields(cmd.Namespace + " " + cmd.Name)
	for _, word := range cmdWords {
		excludeWords[strings.ToLower(word)] = true
		// Also exclude individual parts of hyphenated words (e.g., "list" and "all" from "list-all")
		hyphenParts := strings.Split(word, "-")
		for _, part := range hyphenParts {
			if part != "" {
				excludeWords[strings.ToLower(part)] = true
			}
		}
	}

	for _, match := range matches {
		if len(match) > 1 {
			word := match[1]
			if !excludeWords[word] {
				lowercaseParams = append(lowercaseParams, word)
			}
		}
	}

	if len(lowercaseParams) > 0 {
		sl.ReportError(cmd.Usage, "usage", "Usage", "capargs", strings.Join(lowercaseParams, ", "))
	}
}

// validateDescription validates command description
func validateDescription(sl validator.StructLevel) {
	cmd := sl.Current().Interface().(Command)

	description := cmd.Description

	// Check capitalization
	if len(description) > 0 && !unicode.IsUpper(rune(description[0])) {
		sl.ReportError(cmd.Description, "description", "Description", "uppercase", "")
	}

	// Check for subject (anti-pattern)
	badStarts := []string{"this command", "plugin to", "commands to", "this plugin", "this is"}
	descLower := strings.ToLower(description)
	for _, badStart := range badStarts {
		if strings.HasPrefix(descLower, badStart) {
			sl.ReportError(cmd.Description, "description", "Description", "nosubject", badStart)
		}
	}

	// Check word count (guideline, not strict)
	wordCount := len(strings.Fields(description))
	if wordCount > maxDescriptionWordCount {
		sl.ReportError(cmd.Description, "description", "Description", "cmdwordcount", strconv.Itoa(maxDescriptionWordCount))
	}

	// Check if command uses plural form but description doesn't indicate listing
	// Only check the LAST word - if it ends with 's' and is not a verb, it's likely a plural noun
	if lastWordIsPluralForm(cmd.Name) {
		// Check if description contains list-related keywords
		listKeywords := []string{"list", "show", "display", "view", "retrieve", "get", "all", "multiple"}
		hasListKeyword := false
		descLower := strings.ToLower(description)

		for _, keyword := range listKeywords {
			if strings.Contains(descLower, keyword) {
				hasListKeyword = true
				break
			}
		}

		if !hasListKeyword {
			sl.ReportError(cmd.Description, "description", "Description", "cleardesc", "")
		}
	}
}

// validateCommandDepth validates command depth does not exceed maximum
func validateCommandDepth(sl validator.StructLevel) {
	cmd := sl.Current().Interface().(Command)
	cmdName := cmd.Name

	// Count depth by splitting on spaces
	parts := strings.Fields(cmdName)

	// Remove base CLI name if present
	if len(parts) > 1 {
		firstPart := strings.ToLower(parts[0])
		if firstPart == "ibmcloud" || firstPart == "ic" {
			parts = parts[1:]
		}
	}

	depth := len(parts)

	if depth > maxCommandDepth {
		sl.ReportError(cmd.Name, "name", "Name", "maxnsdepth", strconv.Itoa(maxCommandDepth))
	}

}

// validateCommandNaming validates command naming conventions
func validateCommandNaming(sl validator.StructLevel) {
	cmd := sl.Current().Interface().(Command)
	cmdName := cmd.Name

	// Check for reserved names
	reservedNames := map[string]bool{
		"--version": true, "-v": true, "--help": true, "-h": true,
	}
	if reservedNames[cmdName] {
		sl.ReportError(cmd.Name, "name", "Name", "noreserves", "")
	}

	// Check minimum length
	parts := strings.Fields(cmdName)
	actualCmd := cmdName
	if len(parts) > 1 {
		actualCmd = strings.Join(parts[1:], " ")
	}

	// Split by hyphens and spaces to get individual segments
	segments := strings.FieldsFunc(actualCmd, func(r rune) bool {
		return r == '-' || r == ' '
	})

	// Check each segment is at least 2 characters
	minSegCount := 2
	for _, segment := range segments {
		if len(segment) < minSegCount {
			sl.ReportError(cmd.Name, "Name", "name", "cmdsegmin", strconv.Itoa(minSegCount)+"|"+segment)
		}
	}

	// Check for common verbs or plural forms (ending with 's')
	words := strings.Fields(strings.ReplaceAll(actualCmd, "-", " "))
	hasCommonVerb := false

	// Special case: Skip validation for "config get/set/unset <parameter>" pattern
	// These commands follow a specific pattern where the parameter name comes after the verb
	// Check if the namespace ends with "config get", "config set", or "config unset"
	isConfigPattern := false
	namespace := strings.ToLower(cmd.Namespace)
	if strings.HasSuffix(namespace, "config get") ||
		strings.HasSuffix(namespace, "config set") ||
		strings.HasSuffix(namespace, "config unset") {
		isConfigPattern = true
	}

	for _, word := range words {
		if commonVerbs[word] {
			hasCommonVerb = true
			break
		}
	}

	// Check if the last word is a plural form
	hasPluralForm := lastWordIsPluralForm(actualCmd)

	// Check if description indicates this is a retrieval/display command
	// Commands like "access-key" are valid if they're getting/showing information
	isRetrievalCommand := false
	if cmd.Description != "" {
		descLower := strings.ToLower(cmd.Description)
		retrievalKeywords := []string{"get", "show", "display", "view", "retrieve", "return", "fetch", "invite"}
		for _, keyword := range retrievalKeywords {
			if strings.Contains(descLower, keyword) {
				isRetrievalCommand = true
				break
			}
		}
	}

	if !hasCommonVerb && !hasPluralForm && !isRetrievalCommand && len(words) > 1 && !isConfigPattern {
		sl.ReportError(cmd.Name, "name", "Name", "noverb", "")
	}

}

// hasUppercase checks if a string contains any uppercase letters
func hasUppercase(s string) bool {
	for _, r := range s {
		if unicode.IsUpper(r) {
			return true
		}
	}
	return false
}

func isCommandDescription(namespace string) bool {
	cmdDescRegex := regexp.MustCompile(`Commands(\[\d+\])*\.Description`)

	return cmdDescRegex.MatchString(namespace)
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
