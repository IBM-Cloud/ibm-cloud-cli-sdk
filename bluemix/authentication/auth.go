package authentication

import (
	"net/url"
)

type GrantType string

func (g GrantType) String() string {
	return string(g)
}

type ResponseType string

func (r ResponseType) String() string {
	return string(r)
}

type TokenRequest struct {
	grantType     GrantType
	params        url.Values
	responseTypes []ResponseType
}

func NewTokenRequest(grantType GrantType) *TokenRequest {
	return &TokenRequest{
		grantType: grantType,
		params:    make(url.Values),
	}
}

func (r *TokenRequest) GrantType() GrantType {
	return r.grantType
}

func (r *TokenRequest) ResponseTypes() []ResponseType {
	return r.responseTypes
}

func (r *TokenRequest) GetTokenParam(key string) string {
	return r.params.Get(key)
}

func (r *TokenRequest) SetResponseType(responseTypes ...ResponseType) {
	r.responseTypes = responseTypes
}

func (r *TokenRequest) SetTokenParam(key, value string) {
	r.params.Set(key, value)
}

func (r *TokenRequest) WithOption(opt TokenOption) *TokenRequest {
	opt(r)
	return r
}

func (r *TokenRequest) SetValue(v url.Values) {
	if r == nil {
		panic("authentication: token request is nil")
	}

	v.Set("grant_type", r.grantType.String())

	var responseTypeStr string
	for i, t := range r.responseTypes {
		if i > 0 {
			responseTypeStr += ","
		}
		responseTypeStr += t.String()
	}
	v.Set("response_type", responseTypeStr)

	for k, ss := range r.params {
		for _, s := range ss {
			v.Set(k, s)
		}
	}
}

type TokenOption func(r *TokenRequest)

func SetResponseType(responseTypes ...ResponseType) TokenOption {
	return func(r *TokenRequest) {
		r.SetResponseType(responseTypes...)
	}
}

func SetTokenParam(key, value string) TokenOption {
	return func(r *TokenRequest) {
		r.SetTokenParam(key, value)
	}
}
