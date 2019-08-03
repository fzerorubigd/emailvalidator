package emailvalidator

import (
	"errors"
	"strings"
)

type domainRulesValidation func(u string) error

var (
	errInvalidChar   = errors.New("invalid character in username based on domain rules")
	errShortUserName = errors.New("short username based on domain rules")
)

var (
	domainRules = map[string]domainRulesValidation{
		"gmail.com": func(u string) error {
			parts := strings.SplitN(u, "+", 2)
			total := len(parts[0])
			l := total - 1
			for i, c := range parts[0] {
				switch c {
				case ' ', '\n', '\t', '&', '-', '_', '<', '>', ',':
					return errInvalidChar
				case '.':
					if i == 0 || i == l {
						return errInvalidChar
					}
					total--
				}
			}

			if total < 6 {
				return errShortUserName
			}

			if len(parts) < 2 {
				return nil
			}

			for _, c := range parts[1] {
				switch c {
				case ' ', '\n', '\t', '&', '-', '_', '<', '>', ',':
					return errInvalidChar
				}
			}

			return nil
		},
	}
)

func isValidUserName(u, domain string) error {
	if fn, ok := domainRules[domain]; ok {
		return fn(u)
	}

	l := len(u) - 1
	for i, c := range u {
		switch c {
		case ' ', '\t', '\n':
			return errInvalidChar
		case '.':
			if i == 0 || i == l {
				return errInvalidChar
			}
		}
	}

	return nil
}
