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
func (njv *NatJSONValue) Kind() Kind {
	switch njv.proxy.(type) {
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
func (njv *NatJSONValue) IsNil() bool {
	return (njv.Kind() == Nil)
}

// IsBool returns true if value is a bool.
func (njv *NatJSONValue) IsBool() bool {
	return (njv.Kind() == Bool)
}

// IsString returns true if value is a string.
func (njv *NatJSONValue) IsString() bool {
	return (njv.Kind() == String)
}

// IsFloat64 returns true if value is a float64.
func (njv *NatJSONValue) IsFloat64() bool {
	return (njv.Kind() == Float64)
}

// IsInt64 returns true if value is a int64.
func (njv *NatJSONValue) IsInt64() bool {
	return (njv.Kind() == Int64)
}

// IsUint64 returns true if value is a uint64.
func (njv *NatJSONValue) IsUint64() bool {
	return (njv.Kind() == Uint64)
}

// Value returns the underlying interface{} value that is being wrapped.
func (njv *NatJSONValue) Value() interface{} {
	return njv.proxy
}

// Bool returns the underlying bool value.
// Panics with runtime error if not a bool.
func (njv *NatJSONValue) Bool() bool {
	return njv.proxy.(bool)
}

// Float64 returns the underlying float64 value.
// Panics with runtime error if not a float64.
func (njv *NatJSONValue) Float64() float64 {
	return njv.proxy.(float64)
}

// Int64 returns the underlying int64 value.
// Panics with runtime error if not a int64.
func (njv *NatJSONValue) Int64() int64 {
	return njv.proxy.(int64)
}

// Uint64 returns the underlying uint64 value.
// Panics with runtime error if not a uint64.
func (njv *NatJSONValue) Uint64() uint64 {
	return njv.proxy.(uint64)
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
func (njv *NatJSONValue) String() string {
	switch njv.proxy.(type) {
	case bool:
		return fmt.Sprintf("%t", njv.proxy.(bool))
	case string:
		return njv.proxy.(string)
	case int64:
		return strconv.FormatInt(njv.proxy.(int64), 10)
	case uint64:
		return strconv.FormatUint(njv.proxy.(uint64), 10)
	case float64:
		return strconv.FormatFloat(njv.proxy.(float64), 'g', -1, 64)
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
// Whether text is considered a number is based on http://json.org
func (njv *NatJSONValue) DecodeJSONValue(text string) (*NatJSONValue, error) {
	var err error
	if text == "null" {
		njv.proxy = nil
	} else if text == "true" {
		njv.proxy = true
	} else if text == "false" {
		njv.proxy = false
	} else if strings.HasPrefix(text, `{`) && strings.HasSuffix(text, `}`) {
		err = ErrNotImplemented
	} else if strings.HasPrefix(text, `[`) && strings.HasSuffix(text, `]`) {
		err = ErrNotImplemented
	} else if strings.HasPrefix(text, `"`) && strings.HasSuffix(text, `"`) {
		err = json.Unmarshal([]byte(text), &njv.proxy)
	} else if jsonNumRegexp.MatchString(text) {
		if strings.ContainsAny(text, ".eE") {
			var f64 float64
			f64, err = strconv.ParseFloat(text, 64)
			njv.proxy = f64
		} else if strings.HasPrefix(text, "0") || strings.HasPrefix(text, "-0") {
			err = ErrInvalidJSON
		} else if strings.HasPrefix(text, "-") {
			var i64 int64
			i64, err = strconv.ParseInt(text, 10, 64)
			njv.proxy = i64
		} else {
			var u64 uint64
			u64, err = strconv.ParseUint(text, 10, 64)
			njv.proxy = u64
		}

		if numErr, ok := err.(*strconv.NumError); ok {
			err = numErr.Err
		}
	} else {
		err = ErrInvalidJSON
	}
	return njv, err
}

// UnmarshalJSON implements the json.Unmarshaler interface for NatJSONValue
func (njv *NatJSONValue) UnmarshalJSON(text []byte) error {
	_, err := njv.DecodeJSONValue(string(text))
	return err
}
