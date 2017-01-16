// This file was generated by counterfeiter
package fakes

import (
	"sync"

	"github.com/IBM-Bluemix/bluemix-cli-sdk/bluemix/token_refresher"
)

type FakeTokenRefresher struct {
	RefreshAuthTokenStub        func() (newToken string, refreshToken string, err error)
	refreshAuthTokenMutex       sync.RWMutex
	refreshAuthTokenArgsForCall []struct{}
	refreshAuthTokenReturns struct {
		result1 string
		result2 string
		result3 error
	}
}

func (fake *FakeTokenRefresher) RefreshAuthToken() (newToken string, refreshToken string, err error) {
	fake.refreshAuthTokenMutex.Lock()
	fake.refreshAuthTokenArgsForCall = append(fake.refreshAuthTokenArgsForCall, struct{}{})
	fake.refreshAuthTokenMutex.Unlock()
	if fake.RefreshAuthTokenStub != nil {
		return fake.RefreshAuthTokenStub()
	} else {
		return fake.refreshAuthTokenReturns.result1, fake.refreshAuthTokenReturns.result2, fake.refreshAuthTokenReturns.result3
	}
}

func (fake *FakeTokenRefresher) RefreshAuthTokenCallCount() int {
	fake.refreshAuthTokenMutex.RLock()
	defer fake.refreshAuthTokenMutex.RUnlock()
	return len(fake.refreshAuthTokenArgsForCall)
}

func (fake *FakeTokenRefresher) RefreshAuthTokenReturns(result1 string, result2 string, result3 error) {
	fake.RefreshAuthTokenStub = nil
	fake.refreshAuthTokenReturns = struct {
		result1 string
		result2 string
		result3 error
	}{result1, result2, result3}
}

var _ token_refresher.TokenRefresher = new(FakeTokenRefresher)
