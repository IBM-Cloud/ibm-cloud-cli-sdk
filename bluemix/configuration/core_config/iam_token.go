package core_config

import (
	"encoding/base64"
	"encoding/json"
	"strings"
)

type IAMTokenInfo struct {
	Identifier string       `json:"identifier"`
	IAMId      string       `json:"iam_id"`
	UserName   string       `json:"name"`
	UserEmail  string       `json:"email"`
	Accounts   AccountsInfo `json:"account"`
}

type AccountsInfo struct {
	AccountID    string `json:"bss"`
	IMSAccountID string `json:"ims"`
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
	tokenParts := strings.Split(token, " ")

	if len(tokenParts) < 2 {
		return
	}

	encodedParts := strings.Split(tokenParts[1], ".")

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
