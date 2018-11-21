package core_config

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"
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
	Expiray     jsonTime     `json:"exp"`
	Issuer      string       `json:"iss"`
	GrantType   string       `json:"grant_type"`
	Scope       string       `json:"scope"`
}

type AccountsInfo struct {
	AccountID    string `json:"bss"`
	IMSAccountID string `json:"ims"`
}

type jsonTime struct {
	time.Time
}

func (t jsonTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Unix())
}

func (t *jsonTime) UnmarshalJSON(bytes []byte) error {
	var sec int64
	if err := json.Unmarshal(bytes, &sec); err != nil {
		return err
	}
	t.Time = time.Unix(sec, 0)
	return nil
}

func NewIAMTokenInfo(token string) IAMTokenInfo {
	tokenJSON, err := decodeAccessToken(token)
	if err != nil {
		return IAMTokenInfo{}
	}

	var info IAMTokenInfo
	err = json.Unmarshal(tokenJSON, &info)
	if err != nil {
		return IAMTokenInfo{}
	}

	return info
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
