package models

import (
	"strings"
)

const regionSeparator = ":"

type Region struct {
	ID   string
	Name string
	Type string
}

func (r Region) Customer() string {
	return customer(strings.Split(r.ID, regionSeparator))
}

func customer(parts []string) string {
	return parts[0]
}

func (r Region) Deployment() string {
	return deployment(strings.Split(r.ID, regionSeparator))
}

func deployment(parts []string) string {
	if len(parts) < 2 {
		return ""
	}
	return parts[1]
}

func (r Region) CloudName() string {
	splits := strings.Split(r.ID, ":")

	customer := customer(splits)
	if customer != "ibm" {
		return customer
	}

	deployment := deployment(splits)
	switch {
	case deployment == "yp":
		return "bluemix"
	case strings.HasPrefix(deployment, "ys"):
		return "staging"
	default:
		return ""
	}
}

func (r Region) CloudType() string {
	return r.Type
}
