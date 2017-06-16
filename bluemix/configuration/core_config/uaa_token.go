package core_config

import "encoding/json"

type UAATokenInfo struct {
	Username string `json:"user_name"`
	Email    string `json:"email"`
	UserGUID string `json:"user_id"`
}

func NewUAATokenInfo(token string) UAATokenInfo {
	tokenJSON, err := decodeAccessToken(token)
	if err != nil {
		return UAATokenInfo{}
	}

	var info UAATokenInfo
	err = json.Unmarshal(tokenJSON, &info)
	if err != nil {
		return UAATokenInfo{}
	}

	return info
}
