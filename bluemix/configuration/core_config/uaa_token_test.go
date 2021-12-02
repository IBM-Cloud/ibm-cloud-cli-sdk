package core_config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type uaaTokenTestCases struct {
	token     string
	name      string
	isExpired bool
}

var TestUAATokenHasExpiredTestCases = []uaaTokenTestCases{
	{
		token:     "",
		name:      "empty token",
		isExpired: true,
	},
	{
		token:     "eyJraWQiOiIyMDIxMTAxODA4MTkiLCJhbGciOiJSUzI1NiJ9.eyJpYW1faWQiOiJJQk1pZC02NjYwMDE1UktKIiwiaWQiOiJJQk1pZC02NjYwMDE1UktKIiwicmVhbG1pZCI6IklCTWlkIiwic2Vzc2lvbl9pZCI6IkMtMDBkNDIyYjAtYzcyZC00MzNmLWE0YmUtMzc2ZjkyZDEyNDliIiwianRpIjoiNzNmMzVmNGQtZmI2Ny00NTc3LThlNGMtNDE3YzA5MDYwNDU3IiwiaWRlbnRpZmllciI6IjY2NjAwMTVSS0oiLCJnaXZlbl9uYW1lIjoiTkFOQSIsImZhbWlseV9uYW1lIjoiQU1GTyIsIm5hbWUiOiJOQU5BIEFNRk8iLCJlbWFpbCI6Im5vYW1mb0BpYm0uY29tIiwic3ViIjoibm9hbWZvQGlibS5jb20iLCJhdXRobiI6eyJzdWIiOiJub2FtZm9AaWJtLmNvbSIsImlhbV9pZCI6IklCTWlkLTY2NjAwMTVSS0oiLCJuYW1lIjoiTkFOQSBBTUZPIiwiZ2l2ZW5fbmFtZSI6Ik5BTkEiLCJmYW1pbHlfbmFtZSI6IkFNRk8iLCJlbWFpbCI6Im5vYW1mb0BpYm0uY29tIn0sImFjY291bnQiOnsiYm91bmRhcnkiOiJnbG9iYWwiLCJ2YWxpZCI6dHJ1ZSwiYnNzIjoiMDY3OGUzOWY3ZWYxNDkyODk1OWM0YzFhOGY2YTdiYmYifSwiaWF0IjoxNjM1NDQyMDI3LCJleHAiOjE2MzU0NDI5MjcsImlzcyI6Imh0dHBzOi8vaWFtLmNsb3VkLmlibS5jb20vaWRlbnRpdHkiLCJncmFudF90eXBlIjoidXJuOmlibTpwYXJhbXM6b2F1dGg6Z3JhbnQtdHlwZTpwYXNzY29kZSIsInNjb3BlIjoiaWJtIG9wZW5pZCIsImNsaWVudF9pZCI6ImJ4IiwiYWNyIjozLCJhbXIiOlsidG90cCIsIm1mYSIsIm90cCIsInB3ZCJdfQ.RsBd371ACEKOlhkTJngqBVDCY90Z-MT-iYb1OiLA5OpLYPZunR0saHUzBLh2LxnV-Jo0oeitPBmIK38jDk8MCb-rZa3qYNB2qe0WgO50bCMLKgwhKqJwVM6jMMpg4vg6up8kH8Ftc61kivaa1GrJKmQkonnHrjgrLo5IB2yfkMEAbUAMPb_jcRfjEsSP44I-Vx3dYIVSZs8bIufkgmDbJjlMmdhRenh57iwtQ7uImFgK2d-qQ-7sWLvhfzj2VdBLRHPa-dWYlrVgOAMpk6SCMz8wh6LcDUx9LdNKHpxMGCXpGT_UUWvwYqBuLTI3nmkIWIb_Cqa6al7-gQKPTC00Fw",
		name:      "expired token",
		isExpired: true,
	},
	{
		token:     "ABCD.DEFG.HIGK",
		name:      "invalid token",
		isExpired: true,
	},
	{
		token:     "eyJhbGciOiJIUzI1NiIsImtpZCI6ImtleS0xIiwidHlwIjoiSldUIn0.eyJqdGkiOiJhZGFkNWYwYzQ0ZDI0ZDlkYmVhN2YyNGIzMDNmOWNhNyIsInN1YiI6IjY3ODdiMzM2LTAwNzUtNGYwYy1hZmZiLWUyOWZjMmVhZWI4OCIsInNjb3BlIjpbIm9wZW5pZCIsInVhYS51c2VyIiwiY2xvdWRfY29udHJvbGxlci5yZWFkIiwicGFzc3dvcmQud3JpdGUiLCJjbG91ZF9jb250cm9sbGVyLndyaXRlIl0sImNsaWVudF9pZCI6ImNmIiwiY2lkIjoiY2YiLCJhenAiOiJjZiIsImdyYW50X3R5cGUiOiJwYXNzd29yZCIsInVzZXJfaWQiOiI2Nzg3YjMzNi0wMDc1LTRmMGMtYWZmYi1lMjlmYzJlYWViODgiLCJvcmlnaW4iOiJ1YWEiLCJ1c2VyX25hbWUiOiJ3YW5nanVubEBjbi5pYm0uY29tIiwiZW1haWwiOiJ3YW5nanVubEBjbi5pYm0uY29tIiwiYXV0aF90aW1lIjoxNTE2MTczMjgxLCJpYXQiOjE4MTYxNzMyODEsImV4cCI6MjkyNjE3Njg4MCwiaXNzIjoiaHR0cHM6Ly91YWEubmcuYmx1ZW1peC5uZXQvb2F1dGgvdG9rZW4iLCJ6aWQiOiJ1YWEiLCJhdWQiOlsiY2xvdWRfY29udHJvbGxlciIsInBhc3N3b3JkIiwiY2YiLCJ1YWEiLCJvcGVuaWQiXX0.CPUPLMtDpGScEbK9pyD5pzOCWXUir7TX-gZSmFFSRLM",
		name:      "token expired in 2062",
		isExpired: false,
	},
}

func TestUAATokenHasExpired(t *testing.T) {
	for _, testCase := range TestUAATokenHasExpiredTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			tokenInfo := NewUAATokenInfo(testCase.token)
			assert.Equal(t, testCase.isExpired, tokenInfo.hasExpired())
		})
	}
}
