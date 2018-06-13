package terminal

import (
	"fmt"
	"os"
	"regexp"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix"
)

type Color uint

const (
	red    Color = 31
	green        = 32
	yellow       = 33
	//	blue          = 34
	magenta = 35
	cyan    = 36
	white   = 37

	defaultFGColor = 39
)

var (
	colorize               func(message string, color Color, bold int) string
	TerminalSupportsColors = isTerminal()
	UserAskedForColors     = ""
)

func init() {
	InitColorSupport()
}

func InitColorSupport() {
	if ColorsEnabled() {
		colorize = func(message string, color Color, bold int) string {
			return fmt.Sprintf("\033[%d;%dm%s\033[0m", bold, color, message)
		}
	} else {
		colorize = func(message string, _ Color, _ int) string {
			return message
		}
	}
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

func Colorize(message string, color Color) string {
	return colorize(message, color, 0)
}

func ColorizeBold(message string, color Color) string {
	return colorize(message, color, 1)
}

var decolorizerRegex = regexp.MustCompile(`\x1B\[([0-9]{1,2}(;[0-9]{1,2})?)?[m|K]`)

func Decolorize(message string) string {
	return string(decolorizerRegex.ReplaceAll([]byte(message), []byte("")))
}

func HeaderColor(message string) string {
	return ColorizeBold(message, defaultFGColor)
}

func CommandColor(message string) string {
	return ColorizeBold(message, yellow)
}

func StoppedColor(message string) string {
	return ColorizeBold(message, white)
}

func AdvisoryColor(message string) string {
	return ColorizeBold(message, yellow)
}

func CrashedColor(message string) string {
	return ColorizeBold(message, red)
}

func FailureColor(message string) string {
	return ColorizeBold(message, red)
}

func SuccessColor(message string) string {
	return ColorizeBold(message, green)
}

func EntityNameColor(message string) string {
	return ColorizeBold(message, cyan)
}

func PromptColor(message string) string {
	return ColorizeBold(message, cyan)
}

func TableContentHeaderColor(message string) string {
	return ColorizeBold(message, cyan)
}

func WarningColor(message string) string {
	return ColorizeBold(message, magenta)
}

func LogStdoutColor(message string) string {
	return Colorize(message, white)
}

func LogStderrColor(message string) string {
	return Colorize(message, red)
}

func LogHealthHeaderColor(message string) string {
	return Colorize(message, white)
}

func LogAppHeaderColor(message string) string {
	return ColorizeBold(message, yellow)
}

func LogSysHeaderColor(message string) string {
	return ColorizeBold(message, cyan)
}

func isTerminal() bool {
	return terminal.IsTerminal(int(os.Stdout.Fd()))
}
