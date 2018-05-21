package bigjsonvalue

import (
	"errors"
	"regexp"
)

// Errors
var (
	// ErrInvalidJSON defines the invalid JSON error
	ErrInvalidJSON = errors.New("invalid JSON")

	// ErrNotImplemented defines the not-implemented error
	ErrNotImplemented = errors.New("not implemented")
)

// Constants
const (
	// JSONNumRegexpPat defines the regexp pattern for matching JSON numbers
	// (integers or floats) based on http://json.org
	JSONNumRegexpPat = `^-?\d+(\.\d+)?([eE][-+]?\d+)?$`
)

// This global regexp is (mostly) thread-safe according to
// https://golang.org/pkg/regexp/#Regexp
var jsonNumRegexp = regexp.MustCompile(JSONNumRegexpPat)
