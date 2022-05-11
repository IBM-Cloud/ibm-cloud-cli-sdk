package plugin

import (
	"fmt"
	"testing"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers"
	"github.com/stretchr/testify/assert"
)

func TestConvertCobraFlagsToPluginFlags(t *testing.T) {
	assert := assert.New(t)
	cmd := testhelpers.GenerateCobraCommand()
	outputFlag := cmd.Flag("output")
	quietFlag := cmd.Flag("quiet")
	deprecateFlag := cmd.Flag("outputJSON")

	flags := ConvertCobraFlagsToPluginFlags(cmd)

	assert.Equal(3, len(flags))

	// NOTE: flags are sorted in lexicographical order
	assert.Equal(outputFlag.Usage, flags[0].Description)
	assert.True(flags[0].HasValue)
	assert.Equal(outputFlag.Hidden, flags[0].Hidden)
	assert.Equal(outputFlag.Name, flags[0].Name)

	assert.Equal(deprecateFlag.Usage, flags[1].Description)
	assert.False(flags[1].HasValue)
	assert.Equal(deprecateFlag.Hidden, flags[1].Hidden)
	assert.Equal(deprecateFlag.Name, flags[1].Name)

	assert.Equal(quietFlag.Usage, flags[2].Description)
	assert.False(flags[2].HasValue)
	assert.Equal(quietFlag.Hidden, flags[2].Hidden)
	assert.Equal(fmt.Sprintf("%s,%s", quietFlag.Shorthand, quietFlag.Name), flags[2].Name)
}
