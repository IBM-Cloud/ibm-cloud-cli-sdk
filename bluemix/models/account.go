package models

type LinkedIMSAccount struct {
	GUID string
}

type Account struct {
	GUID       string
	Name       string
	Owner      string
	IMSAccount LinkedIMSAccount
}
