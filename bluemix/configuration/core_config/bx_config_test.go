package core_config_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/configuration/core_config"
	"github.com/stretchr/testify/assert"
)

// test case 1: no last updated, no enabled
func TestNoLastUpdateAndNoEnabled(t *testing.T) {
	config := prepareConfigForCLI("", t)

	// check
	checkUsageStats(false, false, config, t)

	// enabled
	config.SetUsageStatsEnabled(true)
	checkUsageStats(true, true, config, t)

	// disabled
	config.SetUsageStatsEnabled(false)
	checkUsageStats(false, true, config, t)

	t.Cleanup(cleanupConfigFiles)
}

// test case 2: no last updated, enabled false
func TestNoLastUpdateAndDisabled(t *testing.T) {
	config := prepareConfigForCLI(`{"UsageStatsEnabled": false}`, t)

	// check
	checkUsageStats(false, false, config, t)

	// enabled
	config.SetUsageStatsEnabled(true)
	checkUsageStats(true, true, config, t)

	// disabled
	config.SetUsageStatsEnabled(false)
	checkUsageStats(false, true, config, t)

	t.Cleanup(cleanupConfigFiles)
}

// test case 3: no last updated, enabled true
func TestNoLastUpdateAndEnabled(t *testing.T) {
	config := prepareConfigForCLI(`{"UsageStatsEnabled": true}`, t)

	// check
	checkUsageStats(false, false, config, t)

	// write enabled
	config.SetUsageStatsEnabled(true)
	checkUsageStats(true, true, config, t)

	// disabled
	config.SetUsageStatsEnabled(false)
	checkUsageStats(false, true, config, t)

	t.Cleanup(cleanupConfigFiles)
}

// test case 4: zero update, no enabled
func TestZeroLastUpdateAndNoEnabled(t *testing.T) {
	config := prepareConfigForCLI(`{"UsageStatsEnabledLastUpdate": "0001-01-01T00:00:00Z"}`, t)

	// check
	checkUsageStats(false, false, config, t)

	// enabled
	config.SetUsageStatsEnabled(true)
	checkUsageStats(true, true, config, t)

	// disabled
	config.SetUsageStatsEnabled(false)
	checkUsageStats(false, true, config, t)

	t.Cleanup(cleanupConfigFiles)
}

// test case 5: zero updated, enabled false
func TestZeroLastUpdateAndDisabled(t *testing.T) {
	config := prepareConfigForCLI(`{"UsageStatsEnabledLastUpdate": "0001-01-01T00:00:00Z","UsageStatsEnabled": false}`, t)

	// check
	checkUsageStats(false, false, config, t)

	// enabled
	config.SetUsageStatsEnabled(true)
	checkUsageStats(true, true, config, t)

	// disabled
	config.SetUsageStatsEnabled(false)
	checkUsageStats(false, true, config, t)

	t.Cleanup(cleanupConfigFiles)
}

// test case 6: zero updated, enabled true
func TestZeroLastUpdateAndEnabled(t *testing.T) {
	config := prepareConfigForCLI(`{"UsageStatsEnabledLastUpdate": "0001-01-01T00:00:00Z","UsageStatsEnabled": true}`, t)

	// check
	checkUsageStats(false, false, config, t)

	// enabled
	config.SetUsageStatsEnabled(true)
	checkUsageStats(true, true, config, t)

	// disabled
	config.SetUsageStatsEnabled(false)
	checkUsageStats(false, true, config, t)

	t.Cleanup(cleanupConfigFiles)
}

// test case 7: updated, no enabled
func TestLastUpdateAndNoEnabled(t *testing.T) {
	config := prepareConfigForCLI(`{"UsageStatsEnabledLastUpdate": "2020-03-29T12:23:43.519017+08:00"}`, t)

	// check
	checkUsageStats(false, true, config, t)

	// enabled
	config.SetUsageStatsEnabled(true)
	checkUsageStats(true, true, config, t)

	// disabled
	config.SetUsageStatsEnabled(false)
	checkUsageStats(false, true, config, t)

	t.Cleanup(cleanupConfigFiles)
}

// test case 8: updated, enabled false
func TestLastUpdateAndDisabled(t *testing.T) {
	config := prepareConfigForCLI(`{"UsageStatsEnabledLastUpdate": "2020-03-29T12:23:43.519017+08:00","UsageStatsEnabled": false}`, t)

	// check
	checkUsageStats(false, true, config, t)

	// enabled
	config.SetUsageStatsEnabled(true)
	checkUsageStats(true, true, config, t)

	// disabled
	config.SetUsageStatsEnabled(false)
	checkUsageStats(false, true, config, t)

	t.Cleanup(cleanupConfigFiles)
}

// test case 9: updated, enabled true
func TestLastUpdateAndEnabled(t *testing.T) {
	config := prepareConfigForCLI(`{"UsageStatsEnabledLastUpdate": "2020-03-29T12:23:43.519017+08:00","UsageStatsEnabled": true}`, t)

	// check
	checkUsageStats(true, true, config, t)

	// disable
	config.SetUsageStatsEnabled(false)
	checkUsageStats(false, true, config, t)

	// enable
	config.SetUsageStatsEnabled(true)
	checkUsageStats(true, true, config, t)

	t.Cleanup(cleanupConfigFiles)
}

func checkUsageStats(enabled bool, timeStampExist bool, config core_config.Repository, t *testing.T) {
	assert.Equal(t, config.UsageStatsEnabled(), enabled)
	assert.Equal(t, config.UsageStatsEnabledLastUpdate().IsZero(), !timeStampExist)
}

func prepareConfigForCLI(cliConfigContent string, t *testing.T) core_config.Repository {
	ioutil.WriteFile("config.json", []byte(cliConfigContent), 0644)
	ioutil.WriteFile("cf_config.json", []byte(""), 0644)
	return core_config.NewCoreConfigFromPath("cf_config.json", "config.json", func(err error) {
		t.Fatal(err.Error())
	})
}

func cleanupConfigFiles() {
	os.Remove("config.json")
	os.Remove("cf_config.json")
}
