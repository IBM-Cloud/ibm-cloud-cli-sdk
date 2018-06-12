package terminal

import (
	"fmt"
	"io"
	"os"

	"github.com/mattn/go-colorable"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/i18n"
)

type UI interface {
	// Say prints the formated message
	Say(format string, args ...interface{})

	// Warn prints the formated warning message
	Warn(format string, args ...interface{})

	// Failed prints the formated failure message
	Failed(format string, args ...interface{})

	// OK prints 'OK'
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
}

type terminalUI struct {
	In  io.Reader
	Out io.Writer
}

// NewStdUI initialize a terminal UI with os.Stdin and os.Stdout
func NewStdUI() UI {
	return NewUI(os.Stdin, colorable.NewColorableStdout())
}

// NewUI initialize a terminal UI with io.Reader and io.Writer
func NewUI(in io.Reader, out io.Writer) UI {
	return &terminalUI{
		In:  in,
		Out: out,
	}
}

func (ui *terminalUI) Say(format string, args ...interface{}) {
	if args != nil {
		fmt.Fprintf(ui.Out, format+"\n", args...)
	} else {
		fmt.Fprint(ui.Out, format+"\n")
	}
}

func (ui *terminalUI) Warn(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	ui.Say(WarningColor(message))
}

func (ui *terminalUI) Ok() {
	ui.Say(SuccessColor(T("OK")))
}

func (ui *terminalUI) Failed(format string, args ...interface{}) {
	ui.Say(FailureColor(T("FAILED")))
	ui.Say(format, args...)
	ui.Say("")
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
