package core_config

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/common/types"
)

const (
	SubjectTypeServiceID = "ServiceId"
)

type IAMTokenInfo struct {
	IAMID       string       `json:"iam_id"`
	ID          string       `json:"id"`
	RealmID     string       `json:"realmid"`
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
	Expiry      time.Time
	IssueAt     time.Time
}

type AccountsInfo struct {
	AccountID    string `json:"bss"`
	IMSAccountID string `json:"ims"`
	Valid        bool   `json:"valid"`
}

func NewIAMTokenInfo(token string) IAMTokenInfo {
	tokenJSON, err := decodeAccessToken(token)
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

func decodeAccessToken(token string) (tokenJSON []byte, err error) {
	encodedParts := strings.Split(token, ".")

	if len(encodedParts) < 3 {
		return
	}

	encodedTokenJSON := encodedParts[1]
	return base64Decode(encodedTokenJSON)
}

func base64Decode(encodedData string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(restorePadding(encodedData))
}

func restorePadding(seg string) string {
	switch len(seg) % 4 {
	case 2:
		seg = seg + "=="
	case 3:
		seg = seg + "="
	}
	return seg
}
