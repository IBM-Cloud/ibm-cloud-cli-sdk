package core_config_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/authentication/vpc"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/configuration/core_config"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/models"
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

func TestHasTargetedProfile(t *testing.T) {
	config := prepareConfigForCLI(`{"UsageStatsEnabledLastUpdate": "2020-03-29T12:23:43.519017+08:00","UsageStatsEnabled": true}`, t)

	// check
	checkUsageStats(true, true, config, t)

	// verify profile is empty in config
	assert.False(t, config.HasTargetedProfile())

	// set profile without ID
	mockProfile := models.Profile{
		Name: "sample_name",
	}

	config.SetProfile(mockProfile)
	assert.False(t, config.HasTargetedProfile())

	mockProfile.ID = "mock_ID"

	config.SetProfile(mockProfile)
	assert.True(t, config.HasTargetedProfile())

	// validate profile
	parsedProfile := config.CurrentProfile()
	assert.Equal(t, mockProfile, parsedProfile)

	t.Cleanup(cleanupConfigFiles)
}

func TestHasTargetedComputeResource(t *testing.T) {
	config := prepareConfigForCLI(`{"UsageStatsEnabledLastUpdate": "2020-03-29T12:23:43.519017+08:00","UsageStatsEnabled": true}`, t)

	// check
	checkUsageStats(true, true, config, t)

	// verify profile is empty in config
	assert.False(t, config.HasTargetedProfile())

	// set profile without compute resource
	mockProfile := models.Profile{
		ID:   "mock_ID",
		Name: "sample_name",
	}

	config.SetProfile(mockProfile)
	assert.False(t, config.HasTargetedComputeResource())

	mockCR := models.Authn{
		Name: "mock_name",
		ID:   "my_cr",
	}
	mockProfile.ComputeResource = mockCR

	config.SetProfile(mockProfile)
	assert.True(t, config.HasTargetedComputeResource())

	// validate profile
	parsedProfile := config.CurrentProfile()
	assert.Equal(t, mockProfile, parsedProfile)

	t.Cleanup(cleanupConfigFiles)
}

func TestHasProfileWithUser(t *testing.T) {
	config := prepareConfigForCLI(`{"UsageStatsEnabledLastUpdate": "2020-03-29T12:23:43.519017+08:00","UsageStatsEnabled": true}`, t)

	// check
	checkUsageStats(true, true, config, t)

	// verify profile is empty in config
	assert.False(t, config.HasTargetedProfile())

	// set profile without compute resource
	mockProfile := models.Profile{
		ID:   "mock_ID",
		Name: "sample_name",
	}

	config.SetProfile(mockProfile)
	assert.False(t, config.HasTargetedComputeResource())

	mockCR := models.Authn{
		Name: "mock_name",
		ID:   "my_id",
	}
	mockProfile.User = mockCR

	config.SetProfile(mockProfile)
	assert.False(t, config.HasTargetedComputeResource())
	assert.False(t, config.IsLoggedInAsCRI())

	// validate profile
	parsedProfile := config.CurrentProfile()
	assert.Equal(t, mockProfile, parsedProfile)

	t.Cleanup(cleanupConfigFiles)
}

func TestCRIType(t *testing.T) {
	config := prepareConfigForCLI(`{"UsageStatsEnabledLastUpdate": "2020-03-29T12:23:43.519017+08:00","UsageStatsEnabled": true}`, t)

	// check
	assert.Empty(t, config.CRIType())
	config.SetCRIType("VPC")
	assert.Equal(t, "VPC", config.CRIType())

	t.Cleanup(cleanupConfigFiles)
}

func TestVPCCRITokenURL(t *testing.T) {
	config := prepareConfigForCLI(`{"UsageStatsEnabledLastUpdate": "2020-03-29T12:23:43.519017+08:00","UsageStatsEnabled": true}`, t)

	// check default value
	assert.Equal(t, vpc.DefaultServerEndpoint, config.VPCCRITokenURL())

	// overwrite with custom value and validate
	oldValue := bluemix.EnvCRVpcUrl.Get()
	bluemix.EnvCRVpcUrl.Set("https://ibm.com")
	assert.Equal(t, "https://ibm.com", config.VPCCRITokenURL())
	bluemix.EnvCRVpcUrl.Set(oldValue)

	t.Cleanup(cleanupConfigFiles)
}

func TestIsLoggedInAsProfile(t *testing.T) {
	config := prepareConfigForCLI(`{"UsageStatsEnabledLastUpdate": "2020-03-29T12:23:43.519017+08:00","UsageStatsEnabled": true}`, t)
	testIAMCRTokenData := "Bearer eyJraWQiOiIyMDE3MTAzMC0wMDowMDowMCIsImFsZyI6IkhTMjU2In0.eyJpYW1faWQiOiJpYW0tUHJvZmlsZS05NDQ5N2QwZC0yYWMzLTQxYmYtYTk5My1hNDlkMWIxNDYyN2MiLCJpZCI6IklCTWlkLXRlc3QiLCJyZWFsbWlkIjoiaWFtIiwianRpIjoiMDRkMjBiMjUtZWUyZC00MDBmLTg2MjMtOGNkODA3MGI1NDY4IiwiaWRlbnRpZmllciI6IlByb2ZpbGUtOTQ0OTdkMGQtMmFjMy00MWJmLWE5OTMtYTQ5ZDFiMTQ2MjdjIiwibmFtZSI6Ik15IFByb2ZpbGUiLCJzdWIiOiJQcm9maWxlLTk0NDk3ZDBkLTJhYzMtNDFiZi1hOTkzLWE0OWQxYjE0NjI3YyIsInN1Yl90eXBlIjoiUHJvZmlsZSIsImF1dGhuIjp7InN1YiI6ImNybjp2MTpzdGFnaW5nOnB1YmxpYzppYW0taWRlbnRpdHk6OmEvMThlMzAyMDc0OWNlNDc0NGIwYjQ3MjQ2NmQ2MWZkYjQ6OmNvbXB1dGVyZXNvdXJjZTpGYWtlLUNvbXB1dGUtUmVzb3VyY2UiLCJpYW1faWQiOiJjcm4tY3JuOnYxOnN0YWdpbmc6cHVibGljOmlhbS1pZGVudGl0eTo6YS8xOGUzMDIwNzQ5Y2U0NzQ0YjBiNDcyNDY2ZDYxZmRiNDo6Y29tcHV0ZXJlc291cmNlOkZha2UtQ29tcHV0ZS1SZXNvdXJjZSIsIm5hbWUiOiJteV9jb21wdXRlX3Jlc291cmNlIn0sImFjY291bnQiOnsiYm91bmRhcnkiOiJnbG9iYWwiLCJ2YWxpZCI6dHJ1ZSwiYnNzIjoiZmFrZV9ic3MifSwiaWF0IjoxNjI5OTI5NDYzLCJleHAiOjgwMjk5MzMwNjMsImlzcyI6Imh0dHBzOi8vaWFtLmNsb3VkLmlibS5jb20vaWRlbnRpdHkiLCJncmFudF90eXBlIjoidXJuOmlibTpwYXJhbXM6b2F1dGg6Z3JhbnQtdHlwZTpjci10b2tlbiIsInNjb3BlIjoiaWJtIG9wZW5pZCIsImNsaWVudF9pZCI6ImJ4In0.ACeIK_8Wi0QmgQ19w4J2OA0OKgC4zb6M6PuGuPTEY_E"

	// check
	checkUsageStats(true, true, config, t)

	assert.Empty(t, config.IAMToken())
	config.SetIAMToken(testIAMCRTokenData)
	assert.True(t, config.IsLoggedInAsProfile())
	assert.False(t, config.IsLoggedInAsCRI())

	t.Cleanup(cleanupConfigFiles)
}

func TestIsLoggedInWithServiceID(t *testing.T) {
	// Setup
	config := prepareConfigForCLI(`{"UsageStatsEnabledLastUpdate": "2020-03-29T12:23:43.519017+08:00","UsageStatsEnabled": true}`, t)
	testIAMCRTokenData := "eyJraWQiOiIyMDE3MTAzMC0wMDowMDowMCIsImFsZyI6IkhTMjU2In0.eyJpYW1faWQiOiJpYW0tUHJvZmlsZS05NDQ5N2QwZC0yYWMzLTQxYmYtYTk5My1hNDlkMWIxNDYyN2MiLCJpZCI6IklCTWlkLXRlc3QiLCJyZWFsbWlkIjoiaWFtIiwianRpIjoiMDRkMjBiMjUtZWUyZC00MDBmLTg2MjMtOGNkODA3MGI1NDY4IiwiaWRlbnRpZmllciI6IlByb2ZpbGUtOTQ0OTdkMGQtMmFjMy00MWJmLWE5OTMtYTQ5ZDFiMTQ2MjdjIiwibmFtZSI6Ik15IFNlcnZpY2UiLCJzdWIiOiJQcm9maWxlLTk0NDk3ZDBkLTJhYzMtNDFiZi1hOTkzLWE0OWQxYjE0NjI3YyIsInN1Yl90eXBlIjoiU2VydmljZUlkIiwiYXV0aG4iOnsic3ViIjoiY3JuOnYxOnN0YWdpbmc6cHVibGljOmlhbS1pZGVudGl0eTo6YS8xOGUzMDIwNzQ5Y2U0NzQ0YjBiNDcyNDY2ZDYxZmRiNDo6Y29tcHV0ZXJlc291cmNlOkZha2UtQ29tcHV0ZS1SZXNvdXJjZSIsImlhbV9pZCI6ImNybi1jcm46djE6c3RhZ2luZzpwdWJsaWM6aWFtLWlkZW50aXR5OjphLzE4ZTMwMjA3NDljZTQ3NDRiMGI0NzI0NjZkNjFmZGI0Ojpjb21wdXRlcmVzb3VyY2U6RmFrZS1Db21wdXRlLVJlc291cmNlIiwibmFtZSI6Im15X2NvbXB1dGVfcmVzb3VyY2UifSwiYWNjb3VudCI6eyJib3VuZGFyeSI6Imdsb2JhbCIsInZhbGlkIjp0cnVlLCJic3MiOiJmYWtlX2JzcyJ9LCJpYXQiOjE2Mjk5Mjk0NjMsImV4cCI6MjYyOTkzMzA2MywiaXNzIjoiaHR0cHM6Ly9pYW0uY2xvdWQuaWJtLmNvbS9pZGVudGl0eSIsImdyYW50X3R5cGUiOiJ1cm46aWJtOnBhcmFtczpvYXV0aDpncmFudC10eXBlOmNyLXRva2VuIiwic2NvcGUiOiJpYm0gb3BlbmlkIiwiY2xpZW50X2lkIjoiYngifQ.n3WF2O2KMmhW0nBDN2CooOwgcSDGI2qK858BzaHB6YI"
	config.SetIAMToken(testIAMCRTokenData)

	// Assertions
	assert.True(t, config.IsLoggedInWithServiceID())
}

func TestIsLoggedInAsCRI(t *testing.T) {
	config := prepareConfigForCLI(`{"UsageStatsEnabledLastUpdate": "2020-03-29T12:23:43.519017+08:00","UsageStatsEnabled": true}`, t)

	assert.False(t, config.IsLoggedInAsCRI())
	config.SetIsLoggedInAsCRI(true)
	assert.True(t, config.IsLoggedInAsCRI())

	t.Cleanup(cleanupConfigFiles)
}

func TestClearSession(t *testing.T) {
	config := prepareConfigForCLI(`{"UsageStatsEnabledLastUpdate": "2020-03-29T12:23:43.519017+08:00","UsageStatsEnabled": true}`, t)

	// check
	checkUsageStats(true, true, config, t)

	// verify profile is empty in config
	assert.False(t, config.HasTargetedProfile())

	// set profile
	mockProfile := models.Profile{
		ID:   "mock_ID",
		Name: "sample_name",
	}

	config.SetProfile(mockProfile)
	assert.True(t, config.HasTargetedProfile())

	// clear session
	config.ClearSession()
	assert.False(t, config.HasTargetedProfile())

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
