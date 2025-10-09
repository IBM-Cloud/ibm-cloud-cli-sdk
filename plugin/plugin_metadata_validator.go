package plugin

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/Masterminds/semver"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/i18n"
	"github.com/go-playground/validator/v10"
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

// PluginToValidationErrors is a mapping between a plugin and one or more metadata validation errors
type PluginToValidationErrors map[string][]PluginMetadataError

type PluginMetadataError struct {
	Namespace   string `json:"namespace"`
	CommandName string `json:"command,omitempty"`
	Error       string `json:"error"`
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
		pluginName string
	)
	pluginToMetadataErrs := PluginToValidationErrors{}
	for _, metadata := range metadatas {
		if err := p.validator.Struct(metadata); err != nil {
			validErr := PluginMetadataError{}
			if errors.As(err, &validErrs) {
				metadataErrs = []PluginMetadataError{}
				for _, v := range validErrs {
					validErr.Namespace = v.StructNamespace()

					// attempt to include the CommandName if available
					cmdIdx := p.CommandIndexFromNamespace(v.StructNamespace())
					if cmdIdx != -1 {
						cmd := metadata.Commands[cmdIdx]
						validErr.CommandName = cmd.Namespace + " " + cmd.Name

					}
					validErr.Error = p.ParseValidationTagErrors(v)
					metadataErrs = append(metadataErrs, validErr)
				}
			}

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
			return err.Error()
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
