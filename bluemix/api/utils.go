package api

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/i18n"
	"github.com/Masterminds/semver"
)

const (
	ConstraintAllVersions = "*"
)

var coercableSemver = regexp.MustCompile(`^\d+(\.\d+)?$`)

type SemverConstraintInvalidError struct {
	Constraint string
	Err        error
}

func (e SemverConstraintInvalidError) Error() string {
	return i18n.T("Version constraint {{.Constraint}} is invalid:\n",
		map[string]interface{}{"Constraint": e.Constraint}) + e.Err.Error()
}

type SemverConstraint interface {
	Satisfied(string) bool
	IsRange() bool

	fmt.Stringer
}

func NewSemverConstraint(versionOrRange string) (SemverConstraint, error) {
	versionOrRange = strings.TrimPrefix(versionOrRange, "v")
	versionOrRange = coerce(versionOrRange)

	if _, err := semver.NewVersion(versionOrRange); err == nil {
		return semverVersion(versionOrRange), nil
	}

	constraints, err := semver.NewConstraint(versionOrRange)
	if err != nil {
		return nil, SemverConstraintInvalidError{Constraint: versionOrRange, Err: err}
	}

	return semverRange{repr: versionOrRange, constraints: constraints}, nil
}

type semverVersion string

func (v semverVersion) Satisfied(version string) bool {
	return strings.EqualFold(string(v), version)
}

func (v semverVersion) IsRange() bool {
	return false
}

func (v semverVersion) String() string {
	return string(v)
}

type semverRange struct {
	repr        string // user-provided string representation
	constraints *semver.Constraints
}

func (r semverRange) Satisfied(version string) bool {
	sv, err := semver.NewVersion(version)
	if err != nil {
		return false
	}

	return r.constraints.Check(sv)
}

func (r semverRange) IsRange() bool {
	return true
}

func (r semverRange) String() string {
	return r.repr
}

// coerce takes an incomplete semver range (e.g. '1' or '1.2') and turns them into a valid constraint. github.com/mastermind/semver's
// default behavior will fill any a missing minor/patch with 0's, so we bypass that to create ranges; e.g.
//
//	'1' -> '1.x'
//	'1.2' -> '1.2.x'
func coerce(semverRange string) string {
	if !coercableSemver.MatchString(semverRange) {
		return semverRange
	}
	return semverRange + ".x"
}
