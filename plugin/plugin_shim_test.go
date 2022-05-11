package plugin

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

type cobraTestPlugin struct {
	cmd      *cobra.Command
	metadata string
}

type urfaveTestPlugin struct {
	metadata string
}

func marshalMetadata(meta PluginMetadata) string {
	json, err := json.Marshal(meta)
	if err != nil {
		panic(fmt.Errorf("could not marshal metadata: %v", err.Error()))
	}

	return string(json)
}

var pluginMetadata = PluginMetadata{
	Name: "test",
	Commands: []Command{
		{
			Name:        "list",
			Description: "List your apps, containers and services in the target space.",
			Usage:       "ibmcloud list",
		},
	},
}

func (p *cobraTestPlugin) GetMetadata() PluginMetadata {
	p.cmd = testhelpers.GenerateCobraCommand()
	pluginMetadata.Commands[0].Flags = ConvertCobraFlagsToPluginFlags(p.cmd)
	p.metadata = marshalMetadata(fillMetadata(pluginMetadata))

	return pluginMetadata
}

func (p *urfaveTestPlugin) GetMetadata() PluginMetadata {
	pluginMetadata.Commands[0].Flags = []Flag{
		{
			Name:        "output",
			Description: "Specify output format, only 'JSON' is supported.",
			Hidden:      false,
			HasValue:    true,
		},
	}
	p.metadata = marshalMetadata(fillMetadata(pluginMetadata))

	return pluginMetadata
}

func (p *cobraTestPlugin) Run(context PluginContext, args []string)  {}
func (p *urfaveTestPlugin) Run(context PluginContext, args []string) {}

func TestStartWithArgsWithCobraCommand(t *testing.T) {
	// TODO(me): Consider moving this to a method
	orgStdout := os.Stdout
	stdoutMock := testhelpers.CreateMockStdout()
	stdoutFile := stdoutMock.File

	// cleanup mock
	defer func() {
		os.Stdout = orgStdout
		os.RemoveAll(stdoutFile.Name())
		stdoutFile.Close()
	}()

	// mock stdout with empty file
	os.Stdout = stdoutFile

	cmd := []string{"SendMetadata"}
	pl := &cobraTestPlugin{}

	StartWithArgs(pl, cmd)

	stdoutMockOut := stdoutMock.Read()

	assert.Equal(t, pl.metadata, string(stdoutMockOut))
}

func TestStartWithArgsWithUrfaveCommand(t *testing.T) {
	orgStdout := os.Stdout
	stdoutMock := testhelpers.CreateMockStdout()
	stdoutFile := stdoutMock.File

	// cleanup mock
	defer func() {
		os.Stdout = orgStdout
		os.RemoveAll(stdoutFile.Name())
		stdoutFile.Close()
	}()

	// mock stdout with empty file
	os.Stdout = stdoutFile

	cmd := []string{"SendMetadata"}
	pl := &urfaveTestPlugin{}

	StartWithArgs(pl, cmd)

	stdoutMockOut := stdoutMock.Read()

	assert.Equal(t, pl.metadata, string(stdoutMockOut))

}
