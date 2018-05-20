package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"regexp"
	"strings"
)

// JSONNumRegexpPat defines the regexp pattern for matching JSON numbers
// (integers or floats) based on json.org
const JSONNumRegexpPat = `^-?\d+(\.\d+)?([eE][-+]?\d+)?$`

// This global regexp is (mostly) thread-safe according to
// https://golang.org/pkg/regexp/#Regexp
var jsonNumRegexp = regexp.MustCompile(JSONNumRegexpPat)

// ErrInvalidJSON defines the invalid JSON error
var ErrInvalidJSON = errors.New("invalid JSON")

// ErrNotImplemented defines the not-implemented error
var ErrNotImplemented = errors.New("not implemented")

// BigJSONValue is wrapper around interface{} type to force
// json.Unmarshal() to decode integer values as integers instead of
// float64.  The problem with float64 is that it doesn't have enough
// precision to store exact values of large absolute int64 values.
// So instead of trying to unmarshal JSON into an interface{},
// unmarshal into a BigJSONValue instead.
// The impetus for BigJSONValue is for decoding the JSON encoding of
// WAL Logical-Decoding changes emitted from Postgres wal2json plugin.
type BigJSONValue struct {
	proxy interface{}
}

// Kind returns the kind of BigJSONValue it is holding:
// Returns reflect.Bool if value is a bool.
// Returns reflect.String if value is a string.
// Returns reflect.Int if value is a big.Int.
// Returns reflect.Float64 if value is a big.Float.
// Otherwise returns reflect.Invalid (ie: nil value).
func (bjv *BigJSONValue) Kind() reflect.Kind {
	switch bjv.proxy.(type) {
	case bool:
		return reflect.Bool
	case string:
		return reflect.String
	case big.Int:
		return reflect.Int
	case big.Float:
		return reflect.Float64
	default:
		return reflect.Invalid
	}
}

// IsNil returns true if value is nil (invalid).
func (bjv *BigJSONValue) IsNil() bool {
	return (bjv.Kind() == reflect.Invalid)
}

// IsBool returns true if value is a bool.
func (bjv *BigJSONValue) IsBool() bool {
	return (bjv.Kind() == reflect.Bool)
}

// IsString returns true if value is a string.
func (bjv *BigJSONValue) IsString() bool {
	return (bjv.Kind() == reflect.String)
}

// IsBigFloat returns true if value is a big.Float.
func (bjv *BigJSONValue) IsBigFloat() bool {
	return (bjv.Kind() == reflect.Float64)
}

// IsBigInt returns true if value is a big.Int.
func (bjv *BigJSONValue) IsBigInt() bool {
	return (bjv.Kind() == reflect.Int)
}

// Value returns the underlying interface{} value that is being wrapped.
func (bjv *BigJSONValue) Value() interface{} {
	return bjv.proxy
}

// Bool returns the underlying bool value.
// Panics with runtime error if not a bool.
func (bjv *BigJSONValue) Bool() bool {
	return bjv.proxy.(bool)
}

// BigFloat returns the underlying big.Float value.
// Panics with runtime error if not a big.Float.
func (bjv *BigJSONValue) BigFloat() big.Float {
	return bjv.proxy.(big.Float)
}

// BigInt returns the underlying big.Int value.
// Panics with runtime error if not a big.Int.
func (bjv *BigJSONValue) BigInt() big.Int {
	return bjv.proxy.(big.Int)
}

// String implements fmt.Stringer interface for a BigJSONValue
// Bool values return "true" or "false".
// String values return as-is (no surround quotes are added).
// Number values return with as much precision as possible.
// Nil (invalid) values return "nil".
func (bjv *BigJSONValue) String() string {
	switch bjv.proxy.(type) {
	case bool:
		return fmt.Sprintf("%t", bjv.proxy.(bool))
	case string:
		return bjv.proxy.(string)
	case big.Int:
		bigi := bjv.proxy.(big.Int)
		return bigi.String()
	case big.Float:
		bigf := bjv.proxy.(big.Float)
		return bigf.Text('g', -1)
	default:
		return "nil"
	}
}

// DecodeJSONValue decodes a JSON value, and returns itself.
// Results are undefined if error is returned.
// The text "null" is decoded as a nil value.
// The text "true" and "false" are decoded as bool values.
// Text surrounded by double-quotes are decoded as string values.
// Number text containing period "." or the letters "e" or "E"
// are decoded as big.Float values.
// Otherwise, number text is decoded as big.Int values.
// Text is considered a number based on json.org.
func (bjv *BigJSONValue) DecodeJSONValue(text string) (*BigJSONValue, error) {
	var err error
	if text == "null" {
		bjv.proxy = nil
	} else if text == "true" {
		bjv.proxy = true
	} else if text == "false" {
		bjv.proxy = false
	} else if strings.HasPrefix(text, `{`) && strings.HasSuffix(text, `}`) {
		err = ErrNotImplemented
	} else if strings.HasPrefix(text, `[`) && strings.HasSuffix(text, `]`) {
		err = ErrNotImplemented
	} else if strings.HasPrefix(text, `"`) && strings.HasSuffix(text, `"`) {
		err = json.Unmarshal([]byte(text), &bjv.proxy)
	} else if jsonNumRegexp.MatchString(text) {
		if strings.ContainsAny(text, ".eE") {
			var bigf big.Float
			bigf.SetPrec(128)
			_, _, err = bigf.Parse(text, 10)
			bjv.proxy = bigf
		} else {
			var bigi big.Int
			_, ok := bigi.SetString(text, 10)
			bjv.proxy = bigi
			if !ok {
				err = fmt.Errorf("big.Int.SetString(%s) failed", text)
			}
		}
	} else {
		err = ErrInvalidJSON
	}
	return bjv, err
}

// UnmarshalJSON implements the json.Unmarshaler interface for a BigJSONValue
func (bjv *BigJSONValue) UnmarshalJSON(text []byte) error {
	//fmt.Printf("BigJSONValue.UnmarshalJSON(%s)\n", text)
	_, err := bjv.DecodeJSONValue(string(text))
	return err
}
