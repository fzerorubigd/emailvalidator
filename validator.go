package emailvalidator

import (
	"errors"
	"fmt"
	"strings"
)

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

func isValiUserName(u string) error {
	l := len(u) - 1
	for i, c := range u {
		switch c {
		case ' ', '\t', '\n':
			return errors.New("username contains white space")
		case '.':
			if i == 0 || i == l {
				return errors.New("dot can not begin or end username")
			}
		}
	}

	return nil
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

// Validate try to validate the email address
// TODO : add option support
func Validate(address string) (free bool, disposable bool, err error) {
	username, domain, tld, err := extractEmailParts(address)
	if err != nil {
		return false, false, err
	}

	if !isValidTLD(tld) {
		return false, false, fmt.Errorf("the %s is not valid tld", tld)
	}

	if err := isValiUserName(username); err != nil {
		return false, false, err
	}

	disposable = isDisposable(domain)
	if !disposable {
		free = isFreeProvider(domain)
	}

	return free, disposable, nil
}
