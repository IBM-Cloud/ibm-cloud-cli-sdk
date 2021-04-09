package trace

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/i18n"
)

// Printer will write logs to specific output.
type Printer interface {
	// Print formats using the default formats for its operands and writes to specific output.
	Print(v ...interface{})
	// Printf formats according to a format specifier and writes to specific output.
	Printf(format string, v ...interface{})
	// Println formats using the default formats for its operands and writes to specific output.
	Println(v ...interface{})
}

// Closer will close the specific output
type Closer interface {
	Close() error
}

// PrinterCloser is the Printer which can be closed
type PrinterCloser interface {
	Printer
	Closer
}

// NullLogger will drop all inputs
type NullLogger struct{}

func (l *NullLogger) Print(v ...interface{})                 {}
func (l *NullLogger) Printf(format string, v ...interface{}) {}
func (l *NullLogger) Println(v ...interface{})               {}

type loggerImpl struct {
	*log.Logger
	c io.WriteCloser
}

func (loggerImpl *loggerImpl) Close() error {
	if loggerImpl.c != nil {
		return loggerImpl.c.Close()
	}
	return nil
}

func newLoggerImpl(out io.Writer, prefix string, flag int) *loggerImpl {
	l := log.New(out, prefix, flag)
	c, _ := out.(io.WriteCloser)
	return &loggerImpl{
		Logger: l,
		c:      c,
	}
}

// Logger is the default logger
var Logger Printer = NewLogger("")

// NewLogger returns a printer for the given trace setting.
func NewLogger(bluemixTrace string) Printer {
	switch strings.ToLower(bluemixTrace) {
	case "", "false":
		return new(NullLogger)
	case "true":
		return NewStdLogger()
	default:
		return NewFileLogger(bluemixTrace)
	}
}

// NewStdLogger creates a a printer that writes to StdOut.
func NewStdLogger() PrinterCloser {
	return newLoggerImpl(terminal.ErrOutput, "", 0)
}

// NewFileLogger creates a printer that writes to the given file path.
func NewFileLogger(path string) PrinterCloser {
	file, err := os.OpenFile(filepath.Clean(path), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0600)
	if err != nil {
		logger := NewStdLogger()
		logger.Printf(T("An error occurred when creating log file '{{.Path}}':\n{{.Error}}\n\n", map[string]interface{}{"Path": path, "Error": err.Error()}))
		return logger
	}
	return newLoggerImpl(file, "", 0)
}

var privateDataPlaceholder = "[PRIVATE DATA HIDDEN]"

// Sanitize returns a clean string with sensitive user data in the input
// replaced by PRIVATE_DATA_PLACEHOLDER.
func Sanitize(input string) string {
	re := regexp.MustCompile(`(?m)^(Authorization|X-Auth\S*): .*`)
	sanitized := re.ReplaceAllString(input, "$1: "+privateDataPlaceholder)

	re = regexp.MustCompile(`(?i)(password|token|apikey|passcode)=[^&]*(&|$)`)
	sanitized = re.ReplaceAllString(sanitized, "$1="+privateDataPlaceholder+"$2")

	re = regexp.MustCompile(`(?i)"([^"]*(password|token|apikey)[^"_]*)":\s*"[^\,]*"`)
	sanitized = re.ReplaceAllString(sanitized, fmt.Sprintf(`"$1":"%s"`, privateDataPlaceholder))

	return sanitized
}
