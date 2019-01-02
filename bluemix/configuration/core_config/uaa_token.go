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
}

func NewUAATokenInfo(token string) UAATokenInfo {
	tokenJSON, err := decodeAccessToken(token)
	if err != nil {
		return UAATokenInfo{}
	}

	var t struct {
		UAATokenInfo
		Expiry types.UnixTime `json:"exp"`
	}
	err = json.Unmarshal(tokenJSON, &t)
	if err != nil {
		return UAATokenInfo{}
	}

	ret := t.UAATokenInfo
	ret.Expiry = t.Expiry.Time()
	return ret
}
