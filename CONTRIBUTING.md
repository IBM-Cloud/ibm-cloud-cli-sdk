# Contributing to Atom

:+1:First of all, thanks for your time to contribute!:tada:

The following is a set of guidelines for` contributing to IBM Cloud CLI SDK. If you have any suggestion or issue regarding IBM Cloud CLI, you can go to [ibm-cloud-cli-releases](https://github.com/IBM-Cloud/ibm-cloud-cli-release) and file issues there.

## Contribute Code

### Before You Submit PR

#### Code Style

We follow the offical [CodeReviewComments](https://github.com/golang/go/wiki/CodeReviewComments). Make sure you run [gofmt](https://golang.org/cmd/gofmt/) and [go vet](https://golang.org/cmd/vet/) to fix any major changes.

#### Unit Test

Make sure you have good unit test. Run `go test -cover $(go list ./...)`, and ensure coverage is above 80% for major packages (aka packages other than i18n, fakes, docs...).


#### Commit Message

Good commit message will greatly help review. We recommend [AngularJS spec of commit message](https://docs.google.com/document/d/1QrDFcIiPjSLDn3EL15IJygNPiHORgU1_OOAqWjiDU5Y/edit#heading=h.greljkmo14y0). You can use [commitzen](https://github.com/commitizen/cz-cli) to help you compose the commit.
