package models

type ResourceGroup struct {
	GUID    string
	Name    string
	State   string
	Default bool
	QuotaID string
}
