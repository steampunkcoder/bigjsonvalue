// Package bigjsonvalue provides types to replace interface{} for decoding unknown JSON values
package bigjsonvalue

import (
	"encoding/json"
	"fmt"
	"math/big"
	//"strconv"
	"strings"
)

// BigJSONValue is wrapper around interface{} type to force
// json.Unmarshal() to decode integer values as big.Int instead of
// float64.  The problem with float64 is that it doesn't have enough
// precision to store exact values of large int64 and uint64 values.
// Instead of trying to unmarshal JSON into an interface{},
// unmarshal into a BigJSONValue instead.
//
// Compared to NatJSONValue, BigJSONValue uses big.Int and big.Float to
// store arbitrary-precision numbers, but is slower than NatJSONValue.
type BigJSONValue struct {
	proxy interface{}
}

// Kind returns the kind of BigJSONValue it is holding:
//
// Returns Bool if value is a bool.
//
// Returns String if value is a string.
//
// Returns BigInt if value is a big.Int.
//
// Returns BigFloat if value is a big.Float.
//
// Otherwise returns Nil.
func (bjv *BigJSONValue) Kind() Kind {
	switch bjv.proxy.(type) {
	case bool:
		return Bool
	case string:
		return String
	case big.Int:
		return BigInt
	case big.Float:
		return BigFloat
	default:
		return Nil
	}
}

// IsNil returns true if value is nil.
func (bjv *BigJSONValue) IsNil() bool {
	return (bjv.Kind() == Nil)
}

// IsBool returns true if value is a bool.
func (bjv *BigJSONValue) IsBool() bool {
	return (bjv.Kind() == Bool)
}

// IsString returns true if value is a string.
func (bjv *BigJSONValue) IsString() bool {
	return (bjv.Kind() == String)
}

// IsBigFloat returns true if value is a big.Float.
func (bjv *BigJSONValue) IsBigFloat() bool {
	return (bjv.Kind() == BigFloat)
}

// IsBigInt returns true if value is a big.Int.
func (bjv *BigJSONValue) IsBigInt() bool {
	return (bjv.Kind() == BigInt)
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

// String implements fmt.Stringer interface for BigJSONValue.
//
// Bool values return "true" or "false".
//
// String values return as-is (no surround double-quotes are added).
//
// Number values return with as much precision as possible.
//
// Nil values return "nil".
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
//
// The text "null" is decoded as a nil value.
//
// The text "true" and "false" are decoded as bool values.
//
// Text surrounded by double-quotes are decoded as string values.
//
// Number text containing period "." or the letters "e" or "E"
// are decoded as big.Float values.
//
// Otherwise, number text is decoded as big.Int values.
// Whether text is considered a number is based on http://json.org
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
		} else if strings.HasPrefix(text, "0") || strings.HasPrefix(text, "-0") {
			err = ErrInvalidJSON
		} else {
			var bigi big.Int
			//_, ok := bigi.SetString(text, 10)
			//if !ok {
			//	err = strconv.ErrSyntax
			//}
			err = bigi.UnmarshalJSON([]byte(text))
			bjv.proxy = bigi
		}
	} else {
		err = ErrInvalidJSON
	}
	return bjv, err
}

// UnmarshalJSON implements the json.Unmarshaler interface for BigJSONValue
func (bjv *BigJSONValue) UnmarshalJSON(text []byte) error {
	_, err := bjv.DecodeJSONValue(string(text))
	return err
}
