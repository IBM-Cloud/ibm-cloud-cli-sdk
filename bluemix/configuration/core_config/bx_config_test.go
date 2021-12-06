package core_config_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
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

func TestIsLoggedIn(t *testing.T) {
	config := prepareConfigForCLI(`{"UsageStatsEnabledLastUpdate": "2021-11-29T12:23:43.519017+08:00","UsageStatsEnabled": true}`, t)
	refresh := "eyJhbGciOiJydCJ9.eyJzZXNzaW9uX2lkIjoiQy02YzUzN2U3My0wMjAxLTQ0YmUtODZlZC0zMDY0YTUwMTMwNGYiLCJpYW1faWQiOiJJQk1pZC02NjYwMDE1UktKIiwiYWNjb3VudF9pZCI6IjYwOWMxNGI4NjhmNTQ1OWM5YmZkOGJhOGI4OTZiYmE5In0._WLlp2iYEpA9PTStswtoXPj5GhiXwrEB0EhEVIF-SKu34qL3-gmfjIeN7RJGT--nwvmuuDfcVIdhA2j7MCmjPqk8YUMVmvR38YY7pA81sk7ynBOeJSg_D2_QCDeH-p1waMZlmedvzhVhSJAqhsorQYOR5GMDmz-kiwOwbQ0ewBiX7Bkc5S_49spJG7T5qsBwLjd8EXjCFGGWg4MS-QXn1SnVKBJZ82VyRL8lSTrMCC2DGseA29ptLOhJJldKVjDrfgLAMDILge08Rbz5NZplwkBRLYT7bvMEaO1cj_we6Ya1DPEO70rYkkIJQ-UtYyOUMEVw4th1LKHkKFNYK4oR_it5DpX-w4jTi49yNbV8ragDtmfUQy1dQ7Vxv5Xc7IsF8htFboxtqYwqRi32M0821ftoGYbxZRX2W6BsijmaBUpE_iaBV39ulFJrvW4Uf9fp-GvjIkZo3iJZN8syrr6LQ_RicB73s3rhZ1tIA7i4w7lapSMAgGH2xqufh6Ca31YBN7Karj0Cy6CxX_2P4aLPlAgZL4qJK3gOz56h4hWqmLhurjf5bn0uPXznPAMAoWpim2MwvSSH3EPRbxUyqGusUe-AXhcY_bTjfWibQJ59wL0q3s6gNwkn1y0RKEdlGoc1ofkJU9XfJWc_HCrZ0Lkq5PjQDqk3WkhkxxVGuD92Ha8LNlff-u_hqpBT86nV7D9r9das0vl4etJA6QK8FJMJ5prRNUpKmrVWpMJrSCQr5o8FxvB94yC0H6FcvDTdIl2slVtOO5wtOkaqtcHLm2trpQO6T3nyJQMWd9rUalJjWP32Gpv5yBSEHKC_toJk050oVemk8fis2B5_Qlb86b5vJsnAMqhBC4SttBztlPf0vIDlReMle_HHRnNXUl58Eaady9jgaGqD2b-dU20wjX8p6SLtfMuYbN0hWhYl2_hP2jfzY4zCtQtO3nGNYKeHy_H678NK_WTjvvqvLighORbbi0KINDEVWqepxLxBnzXywWrSMki6CF6016nUh49DMpEpSeziriuxYllg_j2Uwcl396ecuGxkuN9iziHt9qKhnjS4oyjRo8pE3lGw6oc5DMBgxZcpgPJ8QHGlti1rs8vxdCxl2XxFu-Rsckj9jKjTeULKRV4voQKmiFoFELbgodp91fanimBVCLdBygLhm0v04o5tw6J6NkLnu4GPlFvASTLTG9PydIdSnWvP8sXR0T8PS_4l4ot96N2jqd-UKcaQgk5A-fo3NVfTELMz9CvdFYH1r33GtQvO2BtWmlJxU24wge5KtHVwmXe1BFZOxwE_nfV2AMxY-SplWh8lkjM34vJ7YUw8kmvjdSB7nFjCAyxmhPA" // pragma: allowlist secret
	expiredToken := "eyJpYW1faWQiOiJJQk1pZC02NjYwMDE1UktKIiwiaWQiOiJJQk1pZC02NjYwMDE1UktKIiwicmVhbG1pZCI6IklCTWlkIiwic2Vzc2lvbl9pZCI6IkMtMDBkNDIyYjAtYzcyZC00MzNmLWE0YmUtMzc2ZjkyZDEyNDliIiwianRpIjoiNzNmMzVmNGQtZmI2Ny00NTc3LThlNGMtNDE3YzA5MDYwNDU3IiwiaWRlbnRpZmllciI6IjY2NjAwMTVSS0oiLCJnaXZlbl9uYW1lIjoiTkFOQSIsImZhbWlseV9uYW1lIjoiQU1GTyIsIm5hbWUiOiJOQU5BIEFNRk8iLCJlbWFpbCI6Im5vYW1mb0BpYm0uY29tIiwic3ViIjoibm9hbWZvQGlibS5jb20iLCJhdXRobiI6eyJzdWIiOiJub2FtZm9AaWJtLmNvbSIsImlhbV9pZCI6IklCTWlkLTY2NjAwMTVSS0oiLCJuYW1lIjoiTkFOQSBBTUZPIiwiZ2l2ZW5fbmFtZSI6Ik5BTkEiLCJmYW1pbHlfbmFtZSI6IkFNRk8iLCJlbWFpbCI6Im5vYW1mb0BpYm0uY29tIn0sImFjY291bnQiOnsiYm91bmRhcnkiOiJnbG9iYWwiLCJ2YWxpZCI6dHJ1ZSwiYnNzIjoiMDY3OGUzOWY3ZWYxNDkyODk1OWM0YzFhOGY2YTdiYmYifSwiaWF0IjoxNjM1NDQyMDI3LCJleHAiOjE2MzU0NDI5MjcsImlzcyI6Imh0dHBzOi8vaWFtLmNsb3VkLmlibS5jb20vaWRlbnRpdHkiLCJncmFudF90eXBlIjoidXJuOmlibTpwYXJhbXM6b2F1dGg6Z3JhbnQtdHlwZTpwYXNzY29kZSIsInNjb3BlIjoiaWJtIG9wZW5pZCIsImNsaWVudF9pZCI6ImJ4IiwiYWNyIjozLCJhbXIiOlsidG90cCIsIm1mYSIsIm90cCIsInB3ZCJdfQ.RsBd371ACEKOlhkTJngqBVDCY90Z-MT-iYb1OiLA5OpLYPZunR0saHUzBLh2LxnV-Jo0oeitPBmIK38jDk8MCb-rZa3qYNB2qe0WgO50bCMLKgwhKqJwVM6jMMpg4vg6up8kH8Ftc61kivaa1GrJKmQkonnHrjgrLo5IB2yfkMEAbUAMPb_jcRfjEsSP44I-Vx3dYIVSZs8bIufkgmDbJjlMmdhRenh57iwtQ7uImFgK2d-qQ-7sWLvhfzj2VdBLRHPa-dWYlrVgOAMpk6SCMz8wh6LcDUx9LdNKHpxMGCXpGT_UUWvwYqBuLTI3nmkIWIb_Cqa6al7-gQKPTC00Fw"
	newToken := "eyJraWQiOiIyMDE3MTAzMC0wMDowMDowMCIsImFsZyI6IkhTMjU2In0.eyJpYW1faWQiOiJpYW0tUHJvZmlsZS05NDQ5N2QwZC0yYWMzLTQxYmYtYTk5My1hNDlkMWIxNDYyN2MiLCJpZCI6IklCTWlkLXRlc3QiLCJyZWFsbWlkIjoiaWFtIiwianRpIjoiMDRkMjBiMjUtZWUyZC00MDBmLTg2MjMtOGNkODA3MGI1NDY4IiwiaWRlbnRpZmllciI6IlByb2ZpbGUtOTQ0OTdkMGQtMmFjMy00MWJmLWE5OTMtYTQ5ZDFiMTQ2MjdjIiwibmFtZSI6Ik15IFByb2ZpbGUiLCJzdWIiOiJQcm9maWxlLTk0NDk3ZDBkLTJhYzMtNDFiZi1hOTkzLWE0OWQxYjE0NjI3YyIsInN1Yl90eXBlIjoiUHJvZmlsZSIsImF1dGhuIjp7InN1YiI6ImNybjp2MTpzdGFnaW5nOnB1YmxpYzppYW0taWRlbnRpdHk6OmEvMThlMzAyMDc0OWNlNDc0NGIwYjQ3MjQ2NmQ2MWZkYjQ6OmNvbXB1dGVyZXNvdXJjZTpGYWtlLUNvbXB1dGUtUmVzb3VyY2UiLCJpYW1faWQiOiJjcm4tY3JuOnYxOnN0YWdpbmc6cHVibGljOmlhbS1pZGVudGl0eTo6YS8xOGUzMDIwNzQ5Y2U0NzQ0YjBiNDcyNDY2ZDYxZmRiNDo6Y29tcHV0ZXJlc291cmNlOkZha2UtQ29tcHV0ZS1SZXNvdXJjZSIsIm5hbWUiOiJteV9jb21wdXRlX3Jlc291cmNlIn0sImFjY291bnQiOnsiYm91bmRhcnkiOiJnbG9iYWwiLCJ2YWxpZCI6dHJ1ZSwiYnNzIjoiZmFrZV9ic3MifSwiaWF0IjoxNjI5OTI5NDYzLCJleHAiOjgwMjk5MzMwNjMsImlzcyI6Imh0dHBzOi8vaWFtLmNsb3VkLmlibS5jb20vaWRlbnRpdHkiLCJncmFudF90eXBlIjoidXJuOmlibTpwYXJhbXM6b2F1dGg6Z3JhbnQtdHlwZTpjci10b2tlbiIsInNjb3BlIjoiaWJtIG9wZW5pZCIsImNsaWVudF9pZCI6ImJ4In0.ACeIK_8Wi0QmgQ19w4J2OA0OKgC4zb6M6PuGuPTEY_E"
	tests := []struct {
		name       string
		token      string
		newToken   string
		refresh    string
		newRefresh string
		isLoggedIn bool
	}{
		{
			name:       "token is expired and refresh token is present",
			token:      expiredToken,
			newToken:   "bearer " + newToken, // on refresh the bearer header is append to the token
			refresh:    refresh,
			newRefresh: refresh,
			isLoggedIn: true,
		},
		{
			name:       "token is expired and refresh token is NOT present",
			token:      expiredToken,
			newToken:   expiredToken,
			refresh:    "",
			newRefresh: "",
			isLoggedIn: false,
		},
		{
			name:       "token is not expired",
			token:      newToken,
			newToken:   newToken,
			refresh:    refresh,
			newRefresh: refresh,
			isLoggedIn: true,
		},
		{
			name:       "token is not expired and refresh is NOT present",
			token:      newToken,
			newToken:   newToken,
			refresh:    "",
			newRefresh: "",
			isLoggedIn: true,
		},
	}

	IAMEndpoints := models.Endpoints{}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				fmt.Fprintf(w, "{\"access_token\": \"%s\", \"refresh_token\": \"%s\", \"token_type\": \"bearer\"}", newToken, refresh)
			}))
			defer ts.Close()

			IAMEndpoints.PublicEndpoint = ts.URL
			IAMEndpoints.PrivateEndpoint = ts.URL
			IAMEndpoints.PrivateVPCEndpoint = ts.URL

			config.SetIAMToken(test.token)
			config.SetIAMRefreshToken(test.refresh)
			config.SetIAMEndpoints(IAMEndpoints)
			assert.Equal(t, test.isLoggedIn, config.IsLoggedIn())
			assert.Equal(t, test.newToken, config.IAMToken())
			assert.Equal(t, test.newRefresh, config.IAMRefreshToken())
			t.Cleanup(cleanupConfigFiles)
		})
	}
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
