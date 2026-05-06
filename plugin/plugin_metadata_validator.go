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
	utf16Regexp        *regexp.Regexp
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
		utf16Regexp = regexp.MustCompile(`\\u[0-9a-fA-F]{4}`)
	})

	var (
		metadataErrs []PluginMetadataError
		pluginName   string
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
						validErr.CommandName = cmd.Namespace + " " + cmd.Name
					}
					validErr.Error = p.ParseValidationTagErrors(v)
					validErr.Priority = PriorityError // Default priority for struct validation errors
					metadataErrs = append(metadataErrs, validErr)
				}
			}
		}

		// Validate each command
		for _, cmd := range metadata.Commands {
			metadataErrs = append(metadataErrs, p.validateCommand(cmd, metadata.Name)...)
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

func (p pluginMetadataValidate) ParseValidationTagErrors(fieldErr validator.FieldError) string {
	switch fieldErr.Tag() {
	case "min":
		return i18n.T("{{.Field}} must contain at least {{.Param}} element", map[string]any{
			"Field": fieldErr.Field(),
			"Param": fieldErr.Param(),
		})
	case "ne":
		return i18n.T("{{.Field}} must not equal '{{.Param}}'", map[string]any{
			"Field": fieldErr.Field(),
			"Param": fieldErr.Param(),
		})
	case "required":
		return i18n.T("{{.Field}} is required", map[string]any{
			"Field": fieldErr.Field(),
		})
	case "usage":
		if err := validatePluginMetadataUsage(fmt.Sprint(fieldErr.Value())); err != nil {
			return err.Error()
		}
	case "mincliversion":
		if err := validateCliVersionMinimum(fieldErr.Value().(VersionType), fieldErr.Param()); err != nil {
			return fmt.Sprintf("%s. Remediation: Set MinCliVersion to %s or higher to ensure compatibility with supported CLI versions.", err.Error(), fieldErr.Param())
		}
	case "excludesall":
		return i18n.T("{{.Field}} contains the following forbidden characters: {{.Chars}}", map[string]any{
			"Field": fieldErr.Field(),
			"Chars": strings.Join(strings.Split(fieldErr.Param(), ""), " "),
		})
	default:
		return fieldErr.Error()
	}

	return ""
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

func validateCliVersionMinimum(versionType VersionType, minVersion string) error {
	minSemVer, minErr := semver.NewVersion(minVersion)
	if minErr != nil {
		return minErr
	}
	specifiedVersion, verErr := semver.NewVersion(versionType.String())
	if verErr != nil {
		return verErr
	}

	if specifiedVersion.LessThan(minSemVer) {
		return errors.New(i18n.T("MinCliVersion ({{.ProvidedMinVersion}}) is lower than the allowed minimum {{.AllowedMinimum}}", map[string]any{
			"ProvidedMinVersion": versionType.String(),
			"AllowedMinimum":     minVersion,
		}))
	}

	return nil
}

func validatePluginMetadataUsage(usageText string) error {

	// All arguments and flags MUST be explicitly stated (no placeholders)
	if placeholdersRegexp.MatchString(usageText) {
		return errors.New(i18n.T("Usage contains placeholder arguments/flags"))
	}

	// No basic usage text (Just COMMAND)
	parts := strings.SplitAfter(usageText, "COMMAND")
	if len(parts) == 2 && strings.TrimSpace(parts[1]) == "" {
		return errors.New(i18n.T("Usage does not have any usage text besides COMMAND"))
	}

	// All surrounded parenthesis and brackets MUST be closed
	unclosed, unclosedChar, indicies := UnclosedGroupings(usageText)
	if unclosed {
		return errors.New(i18n.T("Usage contains unclosed {{.UnclosedGroup}} between indicies {{.Indicies}}", map[string]any{
			"Indicies":      indicies,
			"UnclosedGroup": string(unclosedChar),
		}))
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

// validateCommand validates a single command against CLI standards
func (p pluginMetadataValidate) validateCommand(cmd Command, pluginName string) []PluginMetadataError {
	var errs []PluginMetadataError

	if cmd.Name == "" {
		return errs
	}

	// Validate command naming
	errs = append(errs, p.validateCommandNaming(cmd, pluginName)...)

	// Validate command depth
	errs = append(errs, p.validateCommandDepth(cmd, pluginName)...)

	// Validate description
	errs = append(errs, p.validateDescription(cmd, pluginName)...)

	// Validate usage
	errs = append(errs, p.validateUsageEnhanced(cmd, pluginName)...)

	// Validate flags
	errs = append(errs, p.validateFlags(cmd, pluginName)...)

	// Set command name for all errors
	for i := range errs {
		if errs[i].CommandName == "" {
			errs[i].CommandName = cmd.Namespace + " " + cmd.Name
		}
	}

	return errs
}

// validateCommandNaming validates command naming conventions
func (p pluginMetadataValidate) validateCommandNaming(cmd Command, pluginName string) []PluginMetadataError {
	var errs []PluginMetadataError
	cmdName := cmd.Name

	// Check for reserved names
	reservedNames := map[string]bool{
		"--version": true, "-v": true, "--help": true, "-h": true,
	}
	if reservedNames[cmdName] {
		errs = append(errs, PluginMetadataError{
			Namespace:   fmt.Sprintf("Command.%s", cmdName),
			Error:       fmt.Sprintf("Command '%s' uses a reserved flag name. These are handled by the CLI framework.", cmdName),
			Priority:    PriorityWarning,
			Remediation: "Remove this command; --version and --help are automatically provided.",
		})
		return errs
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
	for _, segment := range segments {
		if len(segment) < 2 {
			errs = append(errs, PluginMetadataError{
				Namespace:   fmt.Sprintf("Command.%s", cmdName),
				Error:       fmt.Sprintf("Command '%s' contains a segment '%s' that is less than 2 characters. Each word in a command should be at least 2 characters.", cmdName, segment),
				Priority:    PriorityError,
				Remediation: "Use more descriptive words with at least 2 characters for each segment.",
			})
			break
		}
	}

	// Check for uppercase letters
	if hasUppercase(cmdName) {
		errs = append(errs, PluginMetadataError{
			Namespace:   fmt.Sprintf("Command.%s", cmdName),
			Error:       fmt.Sprintf("Command '%s' contains uppercase letters. Use lower case words with hyphens or spaces.", cmdName),
			Priority:    PriorityError,
			Remediation: "Convert command name to lowercase with hyphens or spaces as separators.",
		})
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
		retrievalKeywords := []string{"get", "show", "display", "view", "retrieve", "return", "fetch"}
		for _, keyword := range retrievalKeywords {
			if strings.Contains(descLower, keyword) {
				isRetrievalCommand = true
				break
			}
		}
	}

	if !hasCommonVerb && !hasPluralForm && !isRetrievalCommand && len(words) > 1 && !isConfigPattern {
		errs = append(errs, PluginMetadataError{
			Namespace:   fmt.Sprintf("Command.%s", cmdName),
			Error:       fmt.Sprintf("Command '%s' does not use a common verb or plural form. Consider using verbs like: list, create, update, delete, show, get, set... or plural nouns like 'instances', 'services'", cmdName),
			Priority:    PriorityWarning,
			Remediation: "Use common verbs in command names such as list, create, update, delete, or use plural forms to indicate listing operations.",
		})
	}

	return errs
}

// validateCommandDepth validates command depth does not exceed maximum
func (p pluginMetadataValidate) validateCommandDepth(cmd Command, pluginName string) []PluginMetadataError {
	var errs []PluginMetadataError
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
		errs = append(errs, PluginMetadataError{
			Namespace: fmt.Sprintf("Command.%s", cmdName),
			Error: fmt.Sprintf("Command '%s' has %d levels, exceeding the maximum of %d. Deep command hierarchies are difficult for users to remember and discover.",
				cmdName, depth, maxCommandDepth),
			Priority: PriorityWarning,
			Remediation: fmt.Sprintf("Reduce command depth to %d or fewer levels. Options: (1) Flatten the hierarchy by combining levels, "+
				"(2) Use command flags/options instead of subcommands, (3) Reorganize the command structure to be more intuitive.", maxCommandDepth),
		})
	}

	return errs
}

// validateDescription validates command description
func (p pluginMetadataValidate) validateDescription(cmd Command, pluginName string) []PluginMetadataError {
	var errs []PluginMetadataError
	description := cmd.Description
	cmdName := cmd.Name

	if description == "" {
		errs = append(errs, PluginMetadataError{
			Namespace:   fmt.Sprintf("Command.%s.Description", cmdName),
			Error:       fmt.Sprintf("Command '%s' has no description. All commands must have a clear description.", cmdName),
			Priority:    PriorityError,
			Remediation: "Add a sentence without subject describing what the command does.",
		})
		return errs
	}

	// Check capitalization
	if len(description) > 0 && !unicode.IsUpper(rune(description[0])) {
		errs = append(errs, PluginMetadataError{
			Namespace:   fmt.Sprintf("Command.%s.Description", cmdName),
			Error:       fmt.Sprintf("Description for '%s' should start with a capital letter.", cmdName),
			Priority:    PriorityError,
			Remediation: "Capitalize the first letter of the description.",
		})
	}

	// Check for subject (anti-pattern)
	badStarts := []string{"this command", "plugin to", "commands to", "this plugin", "this is"}
	descLower := strings.ToLower(description)
	for _, badStart := range badStarts {
		if strings.HasPrefix(descLower, badStart) {
			errs = append(errs, PluginMetadataError{
				Namespace:   fmt.Sprintf("Command.%s.Description", cmdName),
				Error:       fmt.Sprintf("Description for '%s' starts with '%s'. Use a sentence without subject.", cmdName, badStart),
				Priority:    PriorityError,
				Remediation: fmt.Sprintf("Remove '%s' and start directly with the action. Example: 'List all instances' instead of 'This command lists all instances'.", badStart),
			})
			break
		}
	}

	// Check word count (guideline, not strict)
	wordCount := len(strings.Fields(description))
	if wordCount > maxDescriptionWordCount {
		errs = append(errs, PluginMetadataError{
			Namespace:   fmt.Sprintf("Command.%s.Description", cmdName),
			Error:       fmt.Sprintf("Description for '%s' has %d words. Consider limiting to less than %d words for better display.", cmdName, wordCount, maxDescriptionWordCount),
			Priority:    PriorityInfo,
			Remediation: "Shorten the description to be more concise.",
		})
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
			errs = append(errs, PluginMetadataError{
				Namespace:   fmt.Sprintf("Command.%s.Description", cmdName),
				Error:       fmt.Sprintf("Command '%s' uses plural form but description doesn't clearly indicate it returns a list or group of items.", cmdName),
				Priority:    PriorityWarning,
				Remediation: "Update description to include words like 'list', 'show', 'display', 'view', 'all', or 'multiple' to clarify it returns multiple items.",
			})
		}
	}

	return errs
}

// validateUsageEnhanced validates usage format with additional checks
func (p pluginMetadataValidate) validateUsageEnhanced(cmd Command, pluginName string) []PluginMetadataError {
	var errs []PluginMetadataError
	usage := cmd.Usage
	cmdName := cmd.Name

	if usage == "" {
		errs = append(errs, PluginMetadataError{
			Namespace:   fmt.Sprintf("Command.%s.Usage", cmdName),
			Error:       fmt.Sprintf("Command '%s' has no usage information.", cmdName),
			Priority:    PriorityError,
			Remediation: "Add usage text showing command syntax with parameters and options.",
		})
		return errs
	}

	// Check that usage starts with 'ibmcloud' or full path
	usageStripped := strings.TrimSpace(usage)
	if !strings.HasPrefix(usageStripped, "ibmcloud") && !strings.Contains(usageStripped[:min(50, len(usageStripped))], "/ibmcloud") {
		errs = append(errs, PluginMetadataError{
			Namespace:   fmt.Sprintf("Command.%s.Usage", cmdName),
			Error:       fmt.Sprintf("Command '%s' usage should start with 'ibmcloud' (lowercase) or the full path to the ibmcloud binary.", cmdName),
			Priority:    PriorityError,
			Remediation: "Start usage examples with 'ibmcloud' in lowercase (e.g., 'ibmcloud plugin-name command ...').",
		})
	}

	// Check for lowercase argument values (should be CAPS)
	// Look for lowercase words that appear to be user-input parameters
	// Filter out common words, command names, and words that are part of the command structure
	// Strategy: Find lowercase words including those in choice operators, but exclude flags and paths

	// Remove flag names (words after --) to avoid false positives
	// Flags like --name, --zone, --tags should not be flagged
	usageWithoutFlags := regexp.MustCompile(`--[a-z][a-z-]*`).ReplaceAllString(usage, "")

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
		"and": true, "or": true, "to": true, "from": true, "with": true, "ibmcloud": true, "ic": true,
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
		errs = append(errs, PluginMetadataError{
			Namespace: fmt.Sprintf("Command.%s.Usage", cmdName),
			Error: fmt.Sprintf("Command '%s' usage contains lowercase argument values: %s. User input values should be in CAPITAL letters.",
				cmdName, strings.Join(lowercaseParams, ", ")),
			Priority:    PriorityWarning,
			Remediation: "Convert argument values to CAPITAL letters (e.g., NAME, INSTANCE_ID, FORMAT).",
		})
	}

	return errs
}

// validateFlags validates command flags/options
func (p pluginMetadataValidate) validateFlags(cmd Command, pluginName string) []PluginMetadataError {
	var errs []PluginMetadataError
	cmdName := cmd.Name

	for _, flag := range cmd.Flags {
		flagName := flag.Name
		if flagName == "" {
			continue
		}

		// Check single letter flags use -
		if len(flagName) == 1 && !strings.HasPrefix(flagName, "-") {
			errs = append(errs, PluginMetadataError{
				Namespace:   fmt.Sprintf("Command.%s.Flag.%s", cmdName, flagName),
				Error:       fmt.Sprintf("Flag '%s' in command '%s' should use '-' prefix.", flagName, cmdName),
				Priority:    PriorityError,
				Remediation: fmt.Sprintf("Use '-%s' for single letter flags.", flagName),
			})
		}

		// Check multi-letter flags use --
		if len(flagName) > 1 && strings.HasPrefix(flagName, "-") && !strings.HasPrefix(flagName, "--") {
			errs = append(errs, PluginMetadataError{
				Namespace:   fmt.Sprintf("Command.%s.Flag.%s", cmdName, flagName),
				Error:       fmt.Sprintf("Flag '%s' in command '%s' should use '--' prefix.", flagName, cmdName),
				Priority:    PriorityError,
				Remediation: fmt.Sprintf("Use '--%s' for multi-letter flags.", strings.TrimPrefix(flagName, "-")),
			})
		}
	}

	return errs
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

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
