package emailvalidator

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"
)

type Options struct {
	mxValidation        int
	mxValidationTimeout time.Duration
}

// OptionSetter is used to handle options in the file
type OptionSetter func(*Options) error

// CheckMX is for checking the timeout (only if the domain is not in free providers and disposable providers)
func CheckMX(timeout time.Duration) OptionSetter {
	return func(opt *Options) error {
		if timeout < time.Microsecond {
			return errors.New("invalid timeout")
		}
		opt.mxValidation = 1
		opt.mxValidationTimeout = timeout
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

// ValidateContext try to validate the email address, the context version, this context used for any
// extra validation used in teh library (like MX validation)
func ValidateContext(ctx context.Context, address string, opts ...OptionSetter) (free bool, disposable bool, err error) {
	opt := &Options{}
	for i := range opts {
		if err := opts[i](opt); err != nil {
			return false, false, err
		}
	}

	username, domain, tld, err := extractEmailParts(address)
	if err != nil {
		return false, false, err
	}

	if !isValidTLD(tld) {
		return false, false, fmt.Errorf("the %s is not valid tld", tld)
	}

	if err := isValidUserName(username, domain); err != nil {
		return false, false, err
	}

	disposable = isDisposable(domain)
	if !disposable {
		free = isFreeProvider(domain)
	}

	// Lets accept free domain list and disposable list as resolved domains
	if !disposable && !free && opt.mxValidation != 0 {
		ctx, cancel := context.WithTimeout(ctx, opt.mxValidationTimeout)
		defer cancel()
		if err := validateMx(ctx, domain); err != nil {
			return false, false, err
		}
	}

	return free, disposable, nil
}

// Validate is for validating single email
func Validate(address string, opts ...OptionSetter) (free bool, disposable bool, err error) {
	return ValidateContext(context.Background(), address, opts...)
}
