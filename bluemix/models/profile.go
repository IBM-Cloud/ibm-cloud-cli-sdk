package models

type Profile struct {
	ID              string
	Name            string
	ComputeResource Authn
	User            Authn
}
