package terminal

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	term "github.com/IBM-Bluemix/bluemix-cli-sdk/bluemix/terminal"
)

type choicesPrompt struct {
	Message string
	Choices []string
}

func ChoicesPrompt(message string, choices ...string) choicesPrompt {
	return choicesPrompt{
		Message: message,
		Choices: choices,
	}
}

type FakeUI struct {
	Prompts         []string
	PasswordPrompts []string
	ChoicesPrompts  []choicesPrompt
	WarnOutputs     []string

	inputs  bytes.Buffer
	outputs bytes.Buffer
}

func NewFakeUI() *FakeUI {
	return &FakeUI{}
}

func (ui *FakeUI) Say(template string, args ...interface{}) {
	message := fmt.Sprintf(template, args...)
	fmt.Fprintln(&ui.outputs, message)
}

func (ui *FakeUI) Ok() {
	ui.Say("OK")
}

func (ui *FakeUI) Failed(template string, args ...interface{}) {
	message := fmt.Sprintf(template, args...)
	ui.Say("FAILED")
	ui.Say(message)
}

func (ui *FakeUI) Warn(template string, args ...interface{}) {
	message := fmt.Sprintf(template, args...)
	ui.WarnOutputs = append(ui.WarnOutputs, message)

	ui.Say(template, args...)
}

func (ui *FakeUI) Prompt(message string, options *term.PromptOptions) *term.Prompt {
	if options.HideInput {
		ui.PasswordPrompts = append(ui.PasswordPrompts, message)
	} else {
		ui.Prompts = append(ui.Prompts, message)
	}

	if ui.inputs.Len() == 0 {
		panic("No input provided to Fake UI for prompt: " + message)
	}

	p := term.NewPrompt(message, options)
	p.Reader = &ui.inputs
	p.Writer = &ui.outputs
	return p
}

func (ui *FakeUI) ChoicesPrompt(message string, choices []string, options *term.PromptOptions) *term.Prompt {
	ui.ChoicesPrompts = append(ui.ChoicesPrompts, ChoicesPrompt(message, choices...))

	if ui.inputs.Len() == 0 {
		panic(fmt.Sprintf("No input provided to Fake UI for choices prompt: %s [%s]",
			message, strings.Join(choices, ", ")))
	}

	p := term.NewChoicesPrompt(message, choices, options)
	p.Reader = &ui.inputs
	p.Writer = &ui.outputs
	return p
}

func (ui *FakeUI) Ask(template string, args ...interface{}) (string, error) {
	message := fmt.Sprintf(template, args...)

	var answer string
	err := ui.Prompt(message,
		&term.PromptOptions{
			HideDefault: true,
			NoLoop:      true,
		},
	).Resolve(&answer)

	return answer, err
}

func (ui *FakeUI) AskForPassword(template string, args ...interface{}) (string, error) {
	message := fmt.Sprintf(template, args...)

	var passwd string
	err := ui.Prompt(
		message,
		&term.PromptOptions{
			HideDefault: true,
			HideInput:   true,
			NoLoop:      true,
		},
	).Resolve(&passwd)

	return passwd, err
}

func (ui *FakeUI) Confirm(template string, args ...interface{}) (bool, error) {
	return ui.ConfirmWithDefault(false, template, args...)
}

func (ui *FakeUI) ConfirmWithDefault(defaultBool bool, template string, args ...interface{}) (bool, error) {
	message := fmt.Sprintf(template, args...)

	var yn = defaultBool
	err := ui.Prompt(
		message,
		&term.PromptOptions{
			HideDefault: true,
			NoLoop:      true,
		},
	).Resolve(&yn)
	return yn, err
}

func (ui *FakeUI) SelectOne(choices []string, template string, args ...interface{}) (int, error) {
	message := fmt.Sprintf(template, args...)

	var selected string
	err := ui.ChoicesPrompt(
		message,
		choices,
		&term.PromptOptions{
			HideDefault: true,
		},
	).Resolve(&selected)

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

func (ui *FakeUI) Table(headers []string) term.Table {
	return term.NewTable(&ui.outputs, headers)
}

func (ui *FakeUI) Inputs(lines ...string) {
	for _, line := range lines {
		ui.inputs.WriteString(line + "\n")
	}
}

func (ui *FakeUI) Outputs() string {
	return ui.outputs.String()
}

func (ui *FakeUI) Writer() io.Writer {
	return &ui.outputs
}
