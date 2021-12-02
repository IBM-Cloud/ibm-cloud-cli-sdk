package core_config

import (
	"encoding/json"
	"time"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/common/types"
)

type UAATokenInfo struct {
	Username string `json:"user_name"`
	Email    string `json:"email"`
	UserGUID string `json:"user_id"`
	Expiry   time.Time
	IssueAt  time.Time
}

func NewUAATokenInfo(token string) UAATokenInfo {
	tokenJSON, err := DecodeAccessToken(token)
	if err != nil {
		return UAATokenInfo{}
	}

	var t struct {
		UAATokenInfo
		Expiry  types.UnixTime `json:"exp"`
		IssueAt types.UnixTime `json:"iat"`
	}
	err = json.Unmarshal(tokenJSON, &t)
	if err != nil {
		return UAATokenInfo{}
	}

	ret := t.UAATokenInfo
	ret.Expiry = t.Expiry.Time()
	ret.IssueAt = t.IssueAt.Time()
	return ret
}

func (t UAATokenInfo) isValid() bool {
	// UAA token without an UserGUID is invalid
	return t.UserGUID != ""
}

func (t UAATokenInfo) hasExpired() bool {
	if !t.isValid() {
		return true
	}
	if t.Expiry.IsZero() {
		return false
	}
	return t.Expiry.Before(time.Now().Add(expiryDelta))
}
