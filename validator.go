package emailvalidator

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"
)

// ValidationState is used to describe the validation result
type ValidationState int

const (
	// ValidationStateNotChecked means the validation is not performed for this
	ValidationStateNotChecked ValidationState = iota // the zero value
	// ValidationStateTrue validated and result was true
	ValidationStateTrue
	// ValidationStateFalse validated and result was false
	ValidationStateFalse
)

// ValidationResult is the optional validation parts, which are not an error normally
type ValidationResult struct {
	FreeProvider ValidationState `json:"free_provider"`
	Disposable   ValidationState `json:"disposable"`
	MXValidation ValidationState `json:"mx_validation"`
	BlackList    ValidationState `json:"black_list"`
}

// Options internally used to handle the options, use OptionSetter to change the option
type Options struct {
	mxValidation        int
	mxValidationTimeout time.Duration
	mxForce             int
}

// OptionSetter is used to handle options in the file
type OptionSetter func(*Options) error

// MarshalJSON json transform for the value
func (v ValidationState) MarshalJSON() ([]byte, error) {
	switch v {
	case ValidationStateNotChecked:
		return []byte("null"), nil
	case ValidationStateFalse:
		return []byte("false"), nil
	case ValidationStateTrue:
		return []byte("true"), nil
	}

	return nil, fmt.Errorf("value %d not supported", v)
}

// CheckMX add the checking the me record to validation. if the force is active then it check the MX even if the
// disposable or free provider is detected.
func CheckMX(timeout time.Duration, forceCheck bool) OptionSetter {
	return func(opt *Options) error {
		if timeout < time.Microsecond {
			return errors.New("invalid timeout")
		}
		opt.mxValidation = 1
		opt.mxValidationTimeout = timeout
		if forceCheck {
			opt.mxForce = 1
		}
		return nil
	}
}

func isDisposable(domain string) bool {
	parts := strings.Split(domain, ".")
	if len(parts) > 2 {
		return wildDisposableDomain[strings.Join(parts[len(parts)-2:], ".")]
	}

	return disposableDomain[domain]
}

func isFreeProvider(domain string) bool {
	return freeProvider[domain]
}

func isValidTLD(in string) bool {
	return tlds[in]
}

func extractEmailParts(email string) (username string, domain string, tld string, err error) {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "", "", "", fmt.Errorf("invalid email address there is %d @", len(parts)-1)
	}

	username = parts[0]
	domain = parts[1]
	parts = strings.Split(domain, ".")
	if len(parts) < 2 {
		return "", "", "", errors.New("invalid email address there is no dot in the host name")
	}
	tld = parts[len(parts)-1]
	return
}

func validateMx(ctx context.Context, domain string) error {
	r := net.Resolver{}
	_, err := r.LookupMX(ctx, domain)
	if err != nil {
		// Based on RFC5321 if no MX record found, we should fallback to A or AAAA record check
		if _, err := r.LookupHost(ctx, domain); err != nil {
			return errors.New("lockup failed")
		}
	}
	return nil
}

func isBlackList(u string) bool {
	return blackList[strings.ToLower(u)]
}

// ValidateContext try to validate the email address, the context version, this context used for any
// extra validation used in the library (like MX validation)
func ValidateContext(ctx context.Context, address string, opts ...OptionSetter) (*ValidationResult, error) {
	opt := &Options{}
	for i := range opts {
		if err := opts[i](opt); err != nil {
			return nil, err
		}
	}

	username, domain, tld, err := extractEmailParts(address)
	if err != nil {
		return nil, err
	}

	/*
		In addition to restrictions on syntax, there is a length limit on
		email addresses.  That limit is a maximum of 64 characters (octets)
		in the "local part" (before the "@") and a maximum of 255 characters
		(octets) in the domain part (after the "@") for a total length of 320
		characters. However, there is a restriction in RFC 2821 on the length of an
		address in MAIL and RCPT commands of 256 characters.  Since addresses
		that do not fit in those fields are not normally useful, the upper
		limit on address lengths should normally be considered to be 256.
	*/
	/*
		https://www.rfc-editor.org/errata/eid1690
		I believe erratum ID 1003 is slightly wrong. RFC 2821 places a 256 character
		limit on the forward-path. But a path is defined as
		Path = "<" [ A-d-l ":" ] Mailbox ">"
		So the forward-path will contain at least a pair of angle brackets in addition to the Mailbox.
		This limits the Mailbox (i.e. the email address) to 254 characters.
	*/
	if len(address) > 254 {
		return nil, errors.New("maximum email address size is 254")
	}

	if len(username) > 64 {
		return nil, errors.New("maximum user name (before @) length is 64")
	}

	if !isValidTLD(tld) {
		return nil, fmt.Errorf("the %s is not valid tld", tld)
	}

	if err := isValidUserName(username, domain); err != nil {
		return nil, err
	}

	res := ValidationResult{
		Disposable:   ValidationStateFalse,
		FreeProvider: ValidationStateFalse,
		BlackList:    ValidationStateFalse,
	}

	var dispOrFree bool
	if isDisposable(domain) {
		res.Disposable = ValidationStateTrue
		dispOrFree = true
	}

	if isFreeProvider(domain) {
		res.FreeProvider = ValidationStateTrue
		dispOrFree = true
	}

	if isBlackList(username) {
		res.BlackList = ValidationStateTrue
	}

	mxCheck := opt.mxValidation == 1 && (!dispOrFree || opt.mxForce == 1)

	if mxCheck {
		res.MXValidation = ValidationStateTrue
		ctx, cancel := context.WithTimeout(ctx, opt.mxValidationTimeout)
		defer cancel()
		if err := validateMx(ctx, domain); err != nil {
			res.MXValidation = ValidationStateFalse
		}
	}

	return &res, nil
}

// Validate is for validating single email
func Validate(address string, opts ...OptionSetter) (*ValidationResult, error) {
	return ValidateContext(context.Background(), address, opts...)
}
