// Package emailvalidator validates email, in a little more different way. it checks for disposable email providers and
// also checks the top level domain for more information
// This package is based on information in the https://github.com/ivolo/disposable-email-domains (MIT License) for disposable domain,
// And the data in https://github.com/daveearley/Email-Validation-Tool (MIT? License) for the free email providers.
package emailvalidator

//go:generate go run generate.go
