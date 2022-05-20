package testhelpers

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

// GenerateCobraCommand will create a cobra command with basic flags used for testing
func GenerateCobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "example",
	}

	cmd.Flags().StringP("output", "", "", "Specify output format, only 'JSON' is supported.")
	cmd.Flags().BoolP("quiet", "q", false, "Suppress verbose output")
	cmd.Flags().BoolP("outputJSON", "", false, "Output data into JSON format")

	// NOTE: Added #nosec tag since flag is not attached to a real command
	cmd.Flags().MarkDeprecated("outputJSON", "outputJSON deprecated use --output instead") // #nosec
	return cmd
}

type mockStdoutFile struct {
	File *os.File
}

// CreateMockStdout will create a temp file used for mocking stdout for testing
func CreateMockStdout() *mockStdoutFile {
	f, err := os.CreateTemp("", "cli_sdk_mock_stdout")
	if err != nil {
		panic(fmt.Errorf("failed to create tmp file for mocking stdout: %v", err.Error()))
	}
	return &mockStdoutFile{
		File: f,
	}
}

// Read will open the temp mock stdout file and return contents as a string
func (m *mockStdoutFile) Read() string {
	out, err := ioutil.ReadFile(m.File.Name())
	if err != nil {
		panic(fmt.Errorf("failed to read stdout file: %v", err.Error()))
	}
	return string(out)
}
