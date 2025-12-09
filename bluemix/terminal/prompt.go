package terminal

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"reflect"
	"strconv"
	"strings"

	"golang.org/x/crypto/ssh/terminal"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/i18n"
)

var (
	ErrInputEmpty          = errors.New("input is empty")
	ErrInputNotNumber      = errors.New("input is not a number")
	ErrInputNotFloatNumber = errors.New("input is not a floating point number")
	ErrInputNotBool        = errors.New("input must be 'y', 'n', 'yes' or 'no'")
	ErrInputOutOfRange     = errors.New("input is out of range")
)

type ValidateFunc func(string) error

type PromptOptions struct {
	Required     bool         // If true, user input is required
	HideInput    bool         // If true, user input is hide, typically when asking password. TODO: support mask later
	HideDefault  bool         // If true, hide default value in the prompt message
	NoLoop       bool         // if true, when input is invalid, return error instead of asking user for retry
	ValidateFunc ValidateFunc // customized input validation function
}

// Prompt represents a terminal prompt. Create Prompt with NewPrompt or
// NewChoicesPrompt
type Prompt struct {
	message string
	choices []string

	options PromptOptions

	Reader io.Reader
	Writer io.Writer
}

// NewPrompt returns a single prompt
func NewPrompt(message string, options *PromptOptions) *Prompt {
	p := &Prompt{
		message: message,
		Reader:  os.Stdin,
		Writer:  os.Stdout,
	}
	if options != nil {
		p.options = *options
	}
	return p
}

// NewChoicesPrompt returns a choice prompt
func NewChoicesPrompt(message string, choices []string, options *PromptOptions) *Prompt {
	p := NewPrompt(message, options)
	p.choices = choices
	return p
}

// Resolve reads user input and resolves it to the destination value
func (p *Prompt) Resolve(dest interface{}) error {
	if len(p.choices) > 0 {
		return p.resolveChoices(dest)
	}
	return p.resolveSingle(dest)
}

func (p *Prompt) resolveSingle(dest interface{}) error {
	prompt, err := p.singlePrompt(dest)
	if err != nil {
		return err
	}

	for {
		input, readErr := p.read(prompt)
		if readErr != nil {
			return fmt.Errorf("%s", T("Could not read from input: ") + readErr.Error())
		}

		var err error

		if input == "" {
			if p.options.Required {
				err = ErrInputEmpty
			}
		} else {
			if p.options.ValidateFunc != nil {
				err = p.options.ValidateFunc(input)
			}

			if err == nil {
				err = resolveValue(input, reflect.ValueOf(dest).Elem())
			}
		}

		if err != nil {
			if p.options.NoLoop {
				return err
			}

			var message string
			switch err {
			case ErrInputEmpty:
				message = T("Please enter value.")
			case ErrInputNotNumber:
				message = T("Please enter a valid number.")
			case ErrInputNotFloatNumber:
				message = T("Please enter a valid floating number.")
			case ErrInputNotBool:
				message = T("Please enter 'y', 'n', 'yes' or 'no'.")
			default:
				message = err.Error()
			}

			fmt.Fprintln(p.Writer, FailureColor(message))
			continue
		}

		return nil
	}
}

func (p *Prompt) singlePrompt(dest interface{}) (string, error) {
	err := checkDestination(dest)
	if err != nil {
		return "", fmt.Errorf("%s (%v)", p.message, err)
	}

	if p.options.HideDefault {
		return p.message, nil
	}

	e := reflect.ValueOf(dest).Elem()
	if e.Kind() == reflect.Interface {
		e = e.Elem()
	}

	var prompt string
	if e.Kind() == reflect.Bool {
		if p.options.Required {
			prompt = p.message + " [y/n]"
		} else {
			if e.Bool() == true {
				prompt = p.message + " [Y/n]"
			} else {
				prompt = p.message + " [y/N]"
			}
		}
	} else {
		if p.options.Required {
			prompt = p.message
		} else {
			prompt = fmt.Sprintf("%s (%v)", p.message, e)
		}
	}

	return prompt, nil
}

func checkDestination(dest interface{}) error {
	rv := reflect.ValueOf(dest)

	// check if pointer
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("invalid destination: non-pointer %T", dest)
	}

	// check if pointer to nil
	if rv.IsNil() {
		return fmt.Errorf("invalid destination: nil")
	}

	// check if pointer to supported type
	e := rv.Elem()
	if e.Kind() == reflect.Interface {
		if !isSupported(e.Elem().Kind()) {
			return fmt.Errorf("invalid destination: pointer to an interface of an unknown type %s", e.Elem().Type())
		}
	} else {
		if !isSupported(e.Kind()) {
			return fmt.Errorf("invalid destination: unknown type %T", dest)
		}
	}

	return nil
}

func isSupported(k reflect.Kind) bool {
	switch k {
	case reflect.String:
	case reflect.Int:
	case reflect.Int8:
	case reflect.Int16:
	case reflect.Int32:
	case reflect.Int64:
	case reflect.Uint:
	case reflect.Uint8:
	case reflect.Uint16:
	case reflect.Uint32:
	case reflect.Uint64:
	case reflect.Float32:
	case reflect.Float64:
	case reflect.Bool:
	default:
		return false
	}
	return true
}

func (p *Prompt) resolveChoices(dest interface{}) error {
	prompt, err := p.choicesPrompt(dest)
	if err != nil {
		return err
	}

	for {
		input, readErr := p.read(prompt)
		if readErr != nil {
			return fmt.Errorf("%s", T("Could not read from input: ") + readErr.Error())
		}

		var selectedNum int
		var err error

		if input == "" {
			if p.options.Required {
				err = ErrInputEmpty
			}
		} else {
			selectedNum, err = strconv.Atoi(input)
			if err != nil {
				err = ErrInputNotNumber
			} else if selectedNum < 1 || selectedNum > len(p.choices) {
				err = ErrInputOutOfRange
			}
		}

		if err != nil {
			if p.options.NoLoop {
				return err
			}

			fmt.Fprintln(p.Writer, FailureColor(T("Please enter a number between 1 to {{.Count}}.", map[string]interface{}{"Count": len(p.choices)})))
			continue
		}

		if selectedNum > 0 {
			reflect.ValueOf(dest).Elem().SetString(p.choices[selectedNum-1])
		}

		return nil
	}
}

func (p *Prompt) choicesPrompt(dest interface{}) (string, error) {
	if _, ok := dest.(*string); !ok {
		return "", fmt.Errorf("%s (unsupported destination type: %T)", p.message, dest)
	}

	prompt := p.message
	defaultChoice := -1
	for i := 0; i < len(p.choices); i++ {
		prompt += "\n" + fmt.Sprintf("%d. %s", i+1, p.choices[i])

		if p.choices[i] == *(dest.(*string)) {
			defaultChoice = i
		}
	}
	prompt += T("\nEnter a number")

	if p.options.HideDefault {
		return prompt, nil
	}

	if !p.options.Required {
		if defaultChoice >= 0 {
			prompt += fmt.Sprintf(" (%d)", defaultChoice+1)
		} else {
			prompt += " ()"
		}
	}

	return prompt, nil
}

func (p *Prompt) read(prompt string) (string, error) {
	fmt.Fprintf(p.Writer, "%s%s ", prompt, PromptColor(">"))

	f, ok := p.Reader.(*os.File)
	isTerminal := ok && terminal.IsTerminal(int(f.Fd()))

	var input string
	var err error
	if p.options.HideInput {
		if isTerminal {
			input, err = readPassword(int(f.Fd()))
		} else {
			input, err = readLine(p.Reader)
		}
		fmt.Fprintln(p.Writer)
	} else {
		if isTerminal {
			input, err = readLine(p.Reader)
		} else {
			input, err = readLine(p.Reader)
			fmt.Fprintln(p.Writer, input)
		}
	}
	return input, err
}

func readPassword(fd int) (string, error) {
	oldState, err := terminal.GetState(fd)
	if err != nil {
		return "", err
	}

	// Catch interrupt signal to restore terminal's initial state
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)
	defer signal.Stop(c)

	go func() {
		<-c
		err := terminal.Restore(fd, oldState)
		if err != nil {
			fmt.Println(err.Error())
		}
		os.Exit(2)
	}()

	passwd, err := terminal.ReadPassword(fd)
	return string(passwd), err
}

func readLine(r io.Reader) (string, error) {
	var line string

	buf := make([]byte, 1)
	for {
		n, err := r.Read(buf)

		if err != nil {
			return "", err
		}

		if n == 0 || buf[0] == '\n' {
			break
		}

		line += string(buf[0])
	}

	return strings.TrimSpace(line), nil
}

func resolveValue(s string, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Interface:
		nv := reflect.New(v.Elem().Type()).Elem()
		err := resolveValue(s, nv)
		if err != nil {
			return err
		}
		v.Set(nv)
	case reflect.String:
		v.SetString(s)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return ErrInputNotNumber
		}
		v.SetInt(val)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return ErrInputNotNumber
		}
		v.SetUint(val)
	case reflect.Float32, reflect.Float64:
		val, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return ErrInputNotFloatNumber
		}
		v.SetFloat(val)
	case reflect.Bool:
		switch strings.ToLower(s) {
		case "y", "yes":
			v.SetBool(true)
		case "n", "no":
			v.SetBool(false)
		default:
			return ErrInputNotBool
		}
	default:
		return fmt.Errorf("Unable to set value of unknown type '%T'", v)
	}

	return nil
}
