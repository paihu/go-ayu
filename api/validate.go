package main

import (
	"regexp"
)

var (
	email_re = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	url_re   = regexp.MustCompile(`^https?://[a-z]([a-z\.])*\.[a-z]{2,3}/.+$`)
)

func validateEmail(s string) bool {
	return email_re.Match([]byte(s))
}
func validateUrl(s string) bool {
	return url_re.Match([]byte(s))
}
