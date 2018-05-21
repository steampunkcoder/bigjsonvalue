package bigjsonvalue

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// NatJSONValue is wrapper around interface{} type to force
// json.Unmarshal() to decode integer values as int64 or uint64 instead
// of float64.  The problem with float64 is that it doesn't have enough
// precision to store exact values of large int64 and uint64 values.
// Instead of trying to unmarshal JSON into an interface{},
// unmarshal into a NatJSONValue instead.
//
// Compared to BigJSONValue, NatJSONValue uses native Golang number types
// int64, uint64, and float64 to store numbers, so is faster than BigJSONValue.
type NatJSONValue struct {
	proxy interface{}
}

// Kind returns the kind of NatJSONValue it is holding:
//
// Returns Bool if value is a bool.
//
// Returns String if value is a string.
//
// Returns Int64 if value is a int64.
//
// Returns Uint64 if value is a uint64.
//
// Returns Float64 if value is a float64.
//
// Otherwise returns Nil.
func (bjv *NatJSONValue) Kind() Kind {
	switch bjv.proxy.(type) {
	case bool:
		return Bool
	case string:
		return String
	case int64:
		return Int64
	case uint64:
		return Uint64
	case float64:
		return Float64
	default:
		return Nil
	}
}

// IsNil returns true if value is nil.
func (bjv *NatJSONValue) IsNil() bool {
	return (bjv.Kind() == Nil)
}

// IsBool returns true if value is a bool.
func (bjv *NatJSONValue) IsBool() bool {
	return (bjv.Kind() == Bool)
}

// IsString returns true if value is a string.
func (bjv *NatJSONValue) IsString() bool {
	return (bjv.Kind() == String)
}

// IsFloat64 returns true if value is a float64.
func (bjv *NatJSONValue) IsFloat64() bool {
	return (bjv.Kind() == Float64)
}

// IsInt64 returns true if value is a int64.
func (bjv *NatJSONValue) IsInt64() bool {
	return (bjv.Kind() == Int64)
}

// IsUint64 returns true if value is a uint64.
func (bjv *NatJSONValue) IsUint64() bool {
	return (bjv.Kind() == Uint64)
}

// Value returns the underlying interface{} value that is being wrapped.
func (bjv *NatJSONValue) Value() interface{} {
	return bjv.proxy
}

// Bool returns the underlying bool value.
// Panics with runtime error if not a bool.
func (bjv *NatJSONValue) Bool() bool {
	return bjv.proxy.(bool)
}

// Float64 returns the underlying float64 value.
// Panics with runtime error if not a float64.
func (bjv *NatJSONValue) Float64() float64 {
	return bjv.proxy.(float64)
}

// Int64 returns the underlying int64 value.
// Panics with runtime error if not a int64.
func (bjv *NatJSONValue) Int64() int64 {
	return bjv.proxy.(int64)
}

// Uint64 returns the underlying uint64 value.
// Panics with runtime error if not a uint64.
func (bjv *NatJSONValue) Uint64() uint64 {
	return bjv.proxy.(uint64)
}

// String implements fmt.Stringer interface for NatJSONValue.
//
// Bool values return "true" or "false".
//
// String values return as-is (no surround double-quotes are added).
//
// Number values return with as much precision as possible.
//
// Nil values return "nil".
func (bjv *NatJSONValue) String() string {
	switch bjv.proxy.(type) {
	case bool:
		return fmt.Sprintf("%t", bjv.proxy.(bool))
	case string:
		return bjv.proxy.(string)
	case int64:
		return strconv.FormatInt(bjv.proxy.(int64), 10)
	case uint64:
		return strconv.FormatUint(bjv.proxy.(uint64), 10)
	case float64:
		return strconv.FormatFloat(bjv.proxy.(float64), 'g', -1, 64)
	default:
		return "nil"
	}
}

// DecodeJSONValue decodes a JSON value, and returns itself.
// Results are undefined if error is returned.
//
// The text "null" is decoded as a nil value.
//
// The text "true" and "false" are decoded as bool values.
//
// Text surrounded by double-quotes are decoded as string values.
//
// Number text containing period "." or the letters "e" or "E"
// are decoded as float64 values.
//
// Otherwise, number text is decoded as int64 for negative values
// or uint64 for positive values.
// Whether text is considered a number is based on http://json.org.
func (bjv *NatJSONValue) DecodeJSONValue(text string) (*NatJSONValue, error) {
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
			var f64 float64
			f64, err = strconv.ParseFloat(text, 64)
			bjv.proxy = f64
		} else if strings.HasPrefix(text, "0") || strings.HasPrefix(text, "-0") {
			err = ErrInvalidJSON
		} else if strings.HasPrefix(text, "-") {
			var i64 int64
			i64, err = strconv.ParseInt(text, 10, 64)
			bjv.proxy = i64
		} else {
			var u64 uint64
			u64, err = strconv.ParseUint(text, 10, 64)
			bjv.proxy = u64
		}

		if numErr, ok := err.(*strconv.NumError); ok {
			err = numErr.Err
		}
	} else {
		err = ErrInvalidJSON
	}
	return bjv, err
}

// UnmarshalJSON implements the json.Unmarshaler interface for NatJSONValue
func (bjv *NatJSONValue) UnmarshalJSON(text []byte) error {
	_, err := bjv.DecodeJSONValue(string(text))
	return err
}
