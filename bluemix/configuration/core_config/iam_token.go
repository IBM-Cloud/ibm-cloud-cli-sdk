package core_config

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/common/types"
)

const (
	SubjectTypeServiceID      = "ServiceId"
	SubjectTypeTrustedProfile = "Profile"
)

const expiryDelta = 10 * time.Second

type IAMTokenInfo struct {
	IAMID       string       `json:"iam_id"`
	ID          string       `json:"id"`
	RealmID     string       `json:"realmid"`
	SessionID   string       `json:"session_id"`
	Identifier  string       `json:"identifier"`
	Firstname   string       `json:"given_name"`
	Lastname    string       `json:"family_name"`
	Fullname    string       `json:"name"`
	UserEmail   string       `json:"email"`
	Accounts    AccountsInfo `json:"account"`
	Subject     string       `json:"sub"`
	SubjectType string       `json:"sub_type"`
	Issuer      string       `json:"iss"`
	GrantType   string       `json:"grant_type"`
	Scope       string       `json:"scope"`
	Authn       Authn        `json:"authn"`
	Expiry      time.Time
	IssueAt     time.Time
}

type AccountsInfo struct {
	AccountID    string `json:"bss"`
	IMSAccountID string `json:"ims"`
	Valid        bool   `json:"valid"`
}

type Authn struct {
	Subject   string `json:"sub"`
	IAMID     string `json:"iam_id"`
	Name      string `json:"name"`
	Firstname string `json:"given_name"`
	Lastname  string `json:"family_name"`
	Email     string `json:"email"`
}

func NewIAMTokenInfo(token string) IAMTokenInfo {
	tokenJSON, err := DecodeAccessToken(token)
	if err != nil {
		return IAMTokenInfo{}
	}

	var t struct {
		IAMTokenInfo
		Expiry  types.UnixTime `json:"exp"`
		IssueAt types.UnixTime `json:"iat"`
	}
	err = json.Unmarshal(tokenJSON, &t)
	if err != nil {
		return IAMTokenInfo{}
	}

	ret := t.IAMTokenInfo
	ret.Expiry = t.Expiry.Time()
	ret.IssueAt = t.IssueAt.Time()
	return ret
}

// DecodeAccessToken will decode an access token string into a raw JSON.
// The encoded string is expected to be in three parts separated by a period.
// This method does not validate the contents of the parts
func DecodeAccessToken(token string) (tokenJSON []byte, err error) {
	encodedParts := strings.Split(token, ".")

	if len(encodedParts) < 3 {
		return
	}

	encodedTokenJSON := encodedParts[1]
	return base64.RawURLEncoding.DecodeString(encodedTokenJSON)
}

func (t IAMTokenInfo) exists() bool {
	// token without an ID is invalid
	return t.ID != ""
}

func (t IAMTokenInfo) hasExpired() bool {
	if !t.exists() {
		return true
	}
	if t.Expiry.IsZero() {
		return false
	}
	return t.Expiry.Before(time.Now().Add(expiryDelta))
}
