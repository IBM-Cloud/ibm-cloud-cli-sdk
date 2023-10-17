package motd_test

import (
	"fmt"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/api"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Versioning utils", func() {
	Describe("SemverConstraint", func() {

		type constraintChecks struct {
			prints      string
			satisfies   []string
			unsatisfies []string
		}

		Context("When constraint is a specific version", func() {
			for version, checks := range map[string]constraintChecks{
				"0.1.1": {
					satisfies:   []string{"0.1.1"},
					unsatisfies: []string{"1.2.3", "999.0.0"},
				},
				"1.2.3": {
					satisfies:   []string{"1.2.3"},
					unsatisfies: []string{"0.1.1", "999.0.0"},
				},
				"999.0.0": {
					satisfies:   []string{"999.0.0"},
					unsatisfies: []string{"0.1.1", "1.2.3"},
				},
			} {
				It(version, func() {
					version, checks := version, checks // gotcha

					sv, err := NewSemverConstraint(version)
					Expect(err).ToNot(HaveOccurred())
					Expect(sv).ToNot(BeNil())
					Expect(sv.IsRange()).To(BeFalse(), "should return a version SemverConstraint")
					Expect(sv.String()).To(Equal(version), "should print itself")

					Expect(sv.Satisfied(checks.satisfies[0])).To(BeTrue(), "should be satisfied by itself")
					for _, v := range checks.unsatisfies {
						Expect(sv.Satisfied(v)).To(BeFalse(), "should be unsatisfied by other version '%s'")
					}

				})
			}
		})

		Context("When constraint is a range", func() {
			for constraint, checks := range map[string]constraintChecks{
				/* catch all range */
				ConstraintAllVersions: {
					satisfies: []string{"0.1.1", "1.2.3", "999.0.0"},
				},

				/* tilde ranges */
				"~1": {
					satisfies:   []string{"1.2.3", "1.0.0", "1.9999.9999"},
					unsatisfies: []string{"0.9999.9999", "2.0.0", "12.1.2"},
				},
				"~1.2": {
					satisfies:   []string{"1.2.0", "1.2.3", "1.2.9999"},
					unsatisfies: []string{"0.9999.9999", "1.3.0", "1.1.9999", "2.0.0", "12.1.2"},
				},
				"~1.2.3": {
					satisfies:   []string{"1.2.3", "1.2.9999"},
					unsatisfies: []string{"0.9999.9999", "1.3.0", "1.1.9999", "2.0.0", "12.1.2"},
				},

				/* caret ranges */
				"^1": {
					satisfies:   []string{"1.2.3", "1.0.0", "1.9999.9999"},
					unsatisfies: []string{"0.9999.9999", "2.0.0", "12.1.2"},
				},
				"^1.2": {
					satisfies:   []string{"1.2.0", "1.2.3", "1.2.9999", "1.3.0"},
					unsatisfies: []string{"0.9999.9999", "1.1.9999", "2.0.0", "12.1.2"},
				},
				"^1.2.3": {
					satisfies:   []string{"1.2.3", "1.2.9999", "1.3.0"},
					unsatisfies: []string{"0.9999.9999", "1.1.9999", "2.0.0", "12.1.2"},
				},

				/* comparison ranges */
				">=1.0.0, <2.0.0": {
					satisfies:   []string{"1.2.3", "1.0.0", "1.9999.9999"},
					unsatisfies: []string{"0.9999.9999", "2.0.0", "12.1.2"},
				},
				">=1.2.0, <2.0.0": {
					satisfies:   []string{"1.2.0", "1.2.3", "1.2.9999", "1.3.0"},
					unsatisfies: []string{"0.9999.9999", "1.1.9999", "2.0.0", "12.1.2"},
				},
				">=1.2.3, <2.0.0": {
					satisfies:   []string{"1.2.3", "1.2.9999", "1.3.0"},
					unsatisfies: []string{"0.9999.9999", "1.1.9999", "2.0.0", "12.1.2"},
				},
				">=1.2.3, <=2.0.0": {
					satisfies:   []string{"1.2.3", "1.2.9999", "1.3.0", "2.0.0"},
					unsatisfies: []string{"0.9999.9999", "1.1.9999", "2.0.1", "12.1.2"},
				},

				/* dash ranges */
				"1 - 1.9999.9999": {
					satisfies:   []string{"1.2.3", "1.0.0", "1.9999.9999"},
					unsatisfies: []string{"0.9999.9999", "2.0.0", "12.1.2"},
				},
				"1.2 - 1.9999.9999": {
					satisfies:   []string{"1.2.0", "1.2.3", "1.2.9999", "1.3.0"},
					unsatisfies: []string{"0.9999.9999", "1.1.9999", "2.0.0", "12.1.2"},
				},
				"1.2.3 - 1.9999.9999": {
					satisfies:   []string{"1.2.3", "1.2.9999", "1.3.0"},
					unsatisfies: []string{"0.9999.9999", "1.1.9999", "2.0.0", "12.1.2"},
				},
				"1.2.3 - 2": {
					satisfies:   []string{"1.2.3", "1.2.9999", "1.3.0", "2.0.0", "2.9999.9999"},
					unsatisfies: []string{"0.9999.9999", "1.1.9999", "12.1.2"},
				},

				/* wildcard ranges */
				"1.x": {
					satisfies:   []string{"1.2.3", "1.0.0", "1.9999.9999"},
					unsatisfies: []string{"0.9999.9999", "2.0.0", "12.1.2"},
				},
				"1.2.X": {
					satisfies:   []string{"1.2.0", "1.2.3", "1.2.9999"},
					unsatisfies: []string{"0.9999.9999", "1.1.9999", "1.3.0", "2.0.0", "12.1.2"},
				},

				/* hybrid */
				">=1.2.*": {
					satisfies:   []string{"1.2.3", "1.2.9999", "1.3.0", "2.0.0", "2.0.1", "12.1.2"},
					unsatisfies: []string{"0.9999.9999", "1.1.9999"},
				},

				/* coercion */
				"1": {
					prints:      "1.x",
					satisfies:   []string{"1.2.3", "1.0.0", "1.9999.9999"},
					unsatisfies: []string{"0.9999.9999", "2.0.0", "12.1.2"},
				},
				"1.2": {
					prints:      "1.2.x",
					satisfies:   []string{"1.2.0", "1.2.3", "1.2.9999"},
					unsatisfies: []string{"0.9999.9999", "1.1.9999", "1.3.0", "2.0.0", "12.1.2"},
				},
				"v1": {
					prints:      "1.x",
					satisfies:   []string{"1.2.3", "1.0.0", "1.9999.9999"},
					unsatisfies: []string{"0.9999.9999", "2.0.0", "12.1.2"},
				},
				"v1.2": {
					prints:      "1.2.x",
					satisfies:   []string{"1.2.0", "1.2.3", "1.2.9999"},
					unsatisfies: []string{"0.9999.9999", "1.1.9999", "1.3.0", "2.0.0", "12.1.2"},
				},
			} {
				It(constraint, func(constraint string, checks constraintChecks) func() {
					return func() {
						sv, err := NewSemverConstraint(constraint)
						Expect(err).ToNot(HaveOccurred())
						Expect(sv).ToNot(BeNil())

						fmt.Printf("constraint is range: %T\n", sv)
						Expect(sv.IsRange()).To(BeTrue(), "should return a range SemverConstraint")

						if checks.prints == "" {
							checks.prints = constraint
						}
						Expect(sv.String()).To(Equal(checks.prints), "should print itself")

						for _, v := range checks.satisfies {
							Expect(sv.Satisfied(v)).To(BeTrue(), "should be satisfied by '%s'", v)
						}

						for _, v := range checks.unsatisfies {
							Expect(sv.Satisfied(v)).To(BeFalse(), "should not be satisfied by '%s'", v)
						}
					}
				}(constraint, checks))
			}
		})

		Context("Invalid constraint", func() {
			for _, constraint := range []string{
				"not-a-version", "1.2.3.4", "a.b.c",
			} {
				It("Returns an error", func() {
					constraint := constraint // gotcha
					sv, err := NewSemverConstraint(constraint)
					Expect(err).To(HaveOccurred())
					Expect(sv).To(BeNil())
				})
			}
		})

	})
})
