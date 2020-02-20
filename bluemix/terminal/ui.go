package terminal

import (
	"fmt"
	"io"
	"os"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/i18n"
)

// UI provides utilities to handle input and output streams
type UI interface {
	// Deprecated: this method could be removed in the future,
	// Use Verbose() if it's interactive message only
	// Or use Print() if it's command output
	// Say prints the formated message, the message will be suppressed in quiet mode
	Say(format string, args ...interface{})

	// Verbose prints message to StdErr, the message will be suppressed in quiet mode
	Verbose(format string, args ...interface{})

	// Warn prints the formated warning message, the message will be suppressed in quiet mode
	Warn(format string, args ...interface{})

	// Failed prints the formated failure message to StdErr, word `FAILED` will be suppressed in quiet mode.
	// But the message itself will not be.
	Failed(format string, args ...interface{})

	// Print will send the message to StdOut, the message will not be suppressed in quiet mode
	Print(format string, args ...interface{})

	// OK prints 'OK', the message will be suppressed in quiet mode
	Ok()

	// Prompt creates a single Prompt
	Prompt(message string, options *PromptOptions) *Prompt

	// ChoicePrompt creates a choice prompt
	ChoicesPrompt(message string, choices []string, options *PromptOptions) *Prompt

	// Ask asks for text answer
	// Deprecated: use Prompt instead
	Ask(format string, args ...interface{}) (answer string, err error)

	// AskForPassword asks for password
	// Deprecated: use Prompt instead
	AskForPassword(format string, args ...interface{}) (answer string, err error)

	// Confirm asks for user confirmation
	// Deprecated: use Prompt instead
	Confirm(format string, args ...interface{}) (bool, error)

	// ConfirmWithDefault asks for user confirmation. If user skipped, return
	// defaultBool Deprecated: use Prompt instead
	ConfirmWithDefault(defaultBool bool, format string, args ...interface{}) (bool, error)

	// SelectOne asks to select one from choices. It returns the selected index.
	// Deprecated: use ChoicesPrompt instead
	SelectOne(choices []string, format string, args ...interface{}) (int, error)

	// Table creates a table with the given headers
	Table(headers []string) Table

	// Writer returns writer of the terminal UI
	Writer() io.Writer

	// Enable or disable quiet mode. Contents passed to Say(), Warn(), Failed(), OK() will be ignored if under quiet mode.
	SetQuiet(bool)

	// Return whether quiet mode is enabled or not
	Quiet() bool
}

type terminalUI struct {
	In     io.Reader
	Out    io.Writer
	ErrOut io.Writer
	quiet  bool
}

// NewStdUI initialize a terminal UI with os.Stdin and os.Stdout
func NewStdUI() UI {
	return NewUI(os.Stdin, Output, ErrOutput)
}

// NewUI initialize a terminal UI with io.Reader and io.Writer
func NewUI(in io.Reader, out io.Writer, errOut io.Writer) UI {
	return &terminalUI{
		In:     in,
		Out:    out,
		ErrOut: errOut,
	}
}

func (ui *terminalUI) Say(format string, args ...interface{}) {
	if ui.quiet {
		return
	}

	ui.Print(format, args...)
}

func (ui *terminalUI) Verbose(format string, args ...interface{}) {
	if ui.quiet {
		return
	}
	ui.Error(format, args...)
}

func (ui *terminalUI) Warn(format string, args ...interface{}) {
	if ui.quiet {
		return
	}

	message := fmt.Sprintf(format, args...)
	ui.Error(WarningColor(message))
}

func (ui *terminalUI) Ok() {
	if ui.quiet {
		return
	}

	ui.Say(SuccessColor(T("OK")))
}

func (ui *terminalUI) Print(format string, args ...interface{}) {
	if args != nil {
		fmt.Fprintf(ui.Out, format+"\n", args...)
	} else {
		fmt.Fprint(ui.Out, format+"\n")
	}
}

func (ui *terminalUI) Error(format string, args ...interface{}) {
	if args != nil {
		fmt.Fprintf(ui.ErrOut, format+"\n", args...)
	} else {
		fmt.Fprint(ui.ErrOut, format+"\n")
	}
}

func (ui *terminalUI) Failed(format string, args ...interface{}) {
	ui.Verbose(FailureColor(T("FAILED")))
	ui.Error(format, args...)
	ui.Verbose("")
}

func (ui *terminalUI) Prompt(message string, options *PromptOptions) *Prompt {
	p := NewPrompt(message, options)
	p.Reader = ui.In
	p.Writer = ui.Out
	return p
}

func (ui *terminalUI) ChoicesPrompt(message string, choices []string, options *PromptOptions) *Prompt {
	p := NewChoicesPrompt(message, choices, options)
	p.Reader = ui.In
	p.Writer = ui.Out
	return p
}

func (ui *terminalUI) Ask(format string, args ...interface{}) (answer string, err error) {
	message := fmt.Sprintf(format, args...)
	err = ui.Prompt(message, &PromptOptions{HideDefault: true, NoLoop: true}).Resolve(&answer)
	return
}

func (ui *terminalUI) AskForPassword(format string, args ...interface{}) (passwd string, err error) {
	message := fmt.Sprintf(format, args...)
	err = ui.Prompt(message, &PromptOptions{HideInput: true, HideDefault: true, NoLoop: true}).Resolve(&passwd)
	return
}

func (ui *terminalUI) Confirm(format string, args ...interface{}) (yn bool, err error) {
	message := fmt.Sprintf(format, args...)
	err = ui.Prompt(message, &PromptOptions{HideDefault: true, NoLoop: true}).Resolve(&yn)
	return
}

func (ui *terminalUI) ConfirmWithDefault(defaultBool bool, format string, args ...interface{}) (yn bool, err error) {
	yn = defaultBool
	message := fmt.Sprintf(format, args...)
	err = ui.Prompt(message, &PromptOptions{HideDefault: true, NoLoop: true}).Resolve(&yn)
	return
}

func (ui *terminalUI) SelectOne(choices []string, format string, args ...interface{}) (int, error) {
	var selected string
	message := fmt.Sprintf(format, args...)

	err := ui.ChoicesPrompt(message, choices, &PromptOptions{HideDefault: true}).Resolve(&selected)
	if err != nil {
		return -1, err
	}

	for i, c := range choices {
		if selected == c {
			return i, nil
		}
	}

	return -1, nil
}

func (ui *terminalUI) Table(headers []string) Table {
	return NewTable(ui.Out, headers)
}

func (ui *terminalUI) Writer() io.Writer {
	return ui.Out
}

func (ui *terminalUI) SetQuiet(quiet bool) {
	ui.quiet = quiet
}

func (ui *terminalUI) Quiet() bool {
	return ui.quiet
}
