# Email Validation in Go

[![Build Status](https://travis-ci.org/fzerorubigd/emailvalidator.svg)](https://travis-ci.org/fzerorubigd/emailvalidator)
[![Coverage Status](https://coveralls.io/repos/github/fzerorubigd/emailvalidator/badge.svg?branch=master)](https://coveralls.io/github/fzerorubigd/emailvalidator?branch=master)
[![GoDoc](https://godoc.org/github.com/fzerorubigd/emailvalidator?status.svg)](https://godoc.org/github.com/fzerorubigd/emailvalidator)
[![Go Report Card](https://goreportcard.com/badge/github.com/fzerorubigd/emailvalidator/die-github-cache-die)](https://goreportcard.com/report/github.com/fzerorubigd/emailvalidator)


package emailvalidator is an Email validating library in go, the idea is not using a regular expression to validate the email, but to check several data sources (all embedded in the library) to validate the email address and also to identify if the email is a disposable email provider or is a free email provider.

This library is in the Alpha stage and I have plan to extend it.

This package is based on information in the https://github.com/ivolo/disposable-email-domains (MIT License) for disposable domain,
and the data in https://github.com/daveearley/Email-Validation-Tool (MIT? License) for the free email providers. also the valid tlds are from https://data.iana.org/TLD/tlds-alpha-by-domain.txt


