package terminal

import (
	"os"
	"regexp"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix"
	"github.com/fatih/color"
	colorable "github.com/mattn/go-colorable"
)

var (
	Output    = colorable.NewColorableStdout()
	ErrOutput = colorable.NewColorableStderr()

	TerminalSupportsColors = isTerminal()
	UserAskedForColors     = ""
)

func init() {
	InitColorSupport()
}

func InitColorSupport() {
	color.NoColor = !ColorsEnabled()
}

func ColorsEnabled() bool {
	return userDidNotDisableColor() &&
		(userEnabledColors() || TerminalSupportsColors)
}

func userEnabledColors() bool {
	return UserAskedForColors == "true" || bluemix.EnvColor.Get() == "true"
}

func userDidNotDisableColor() bool {
	colorEnv := bluemix.EnvColor.Get()
	return colorEnv != "false" && (UserAskedForColors != "false" || colorEnv == "true")
}

func Colorize(message string, color *color.Color) string {
	return color.Sprintf(message)
}

var decolorizerRegex = regexp.MustCompile(`\x1B\[([0-9]{1,2}(;[0-9]{1,2})?)?[m|K]`)

func Decolorize(message string) string {
	return string(decolorizerRegex.ReplaceAll([]byte(message), []byte("")))
}

func HeaderColor(message string) string {
	return Colorize(message, color.New(color.Bold))
}

func CommandColor(message string) string {
	return Colorize(message, color.New(color.FgYellow, color.Bold))
}

func StoppedColor(message string) string {
	return Colorize(message, color.New(color.FgWhite, color.Bold))
}

func AdvisoryColor(message string) string {
	return Colorize(message, color.New(color.FgYellow, color.Bold))
}

func CrashedColor(message string) string {
	return Colorize(message, color.New(color.FgRed, color.Bold))
}

func FailureColor(message string) string {
	return Colorize(message, color.New(color.FgRed, color.Bold))
}

func SuccessColor(message string) string {
	return Colorize(message, color.New(color.FgGreen, color.Bold))
}

func EntityNameColor(message string) string {
	return Colorize(message, color.New(color.FgCyan, color.Bold))
}

func PromptColor(message string) string {
	return Colorize(message, color.New(color.FgCyan, color.Bold))
}

func TableContentHeaderColor(message string) string {
	return Colorize(message, color.New(color.FgCyan, color.Bold))
}

func WarningColor(message string) string {
	return Colorize(message, color.New(color.FgMagenta, color.Bold))
}

func LogStdoutColor(message string) string {
	return Colorize(message, color.New(color.FgWhite, color.Bold))
}

func LogStderrColor(message string) string {
	return Colorize(message, color.New(color.FgRed, color.Bold))
}

func LogHealthHeaderColor(message string) string {
	return Colorize(message, color.New(color.FgWhite, color.Bold))
}

func LogAppHeaderColor(message string) string {
	return Colorize(message, color.New(color.FgYellow, color.Bold))
}

func LogSysHeaderColor(message string) string {
	return Colorize(message, color.New(color.FgCyan, color.Bold))
}

func isTerminal() bool {
	return terminal.IsTerminal(int(os.Stdout.Fd()))
}
