package bigjsonvalue

import (
	"encoding/json"
	"strconv"
	"testing"
)

type testRec struct {
	jsonStr   string
	bigString string
	bigKind   Kind
	bigIsNil  bool
	bigErr    error
	natString string
	natKind   Kind
	natIsNil  bool
	natErr    error
}

func TestBigDefaultValueIsNil(t *testing.T) {
	bjvN := BigJSONValue{}
	if !bjvN.IsNil() {
		t.Errorf("bjvN is not nil BigJSONValue")
	}
	if bjvN.Kind() != Nil {
		t.Errorf("bjvN.Kind() unexpectedly returned %s", bjvN.Kind())
	}

	bjvM := new(BigJSONValue)
	if !bjvM.IsNil() {
		t.Errorf("bjvM is not nil BigJSONValue")
	}
	if bjvM.Kind() != Nil {
		t.Errorf("bjvM.Kind() unexpectedly returned %s", bjvM.Kind())
	}
}

var testList = []testRec{
	{`null`,
		`nil`, Nil, true, nil,
		`nil`, Nil, true, nil},
	{`true`,
		`true`, Bool, false, nil,
		`true`, Bool, false, nil},
	{`false`,
		`false`, Bool, false, nil,
		`false`, Bool, false, nil},
	{`"a\\b\\c new\nline"`,
		"a\\b\\c new\nline", String, false, nil,
		"a\\b\\c new\nline", String, false, nil},
	{`18446744073709551616`, // math.MaxUint64 + 1
		"18446744073709551616", BigInt, false, nil,
		"18446744073709551615", Uint64, false, strconv.ErrRange},
	{`-9223372036854775809`, // math.MinInt64 - 1
		"-9223372036854775809", BigInt, false, nil,
		"-9223372036854775808", Int64, false, strconv.ErrRange},
	{`4.940656458412465441765687928682213723651e-324`, // math.SmallestNonzeroFloat64
		"4.94065645841246544176568792868221372365e-324", BigFloat, false, nil,
		"5e-324", Float64, false, nil},
	{`1.797693134862315708145274237317043567981e+308`, // math.MaxFloat64
		"1.79769313486231570814527423731704356798e+308", BigFloat, false, nil,
		"1.7976931348623157e+308", Float64, false, nil},
	{`1.7976931348623159e+308`, // more than math.MaxFloat64
		"1.7976931348623159e+308", BigFloat, false, nil,
		"+Inf", Float64, false, strconv.ErrRange},
	{`-1.7976931348623159e+308`, // less than -math.MaxFloat64
		"-1.7976931348623159e+308", BigFloat, false, nil,
		"-Inf", Float64, false, strconv.ErrRange},
	{`987654321987654321`,
		"987654321987654321", BigInt, false, nil,
		"987654321987654321", Uint64, false, nil},
	{`-987654321987654321`,
		"-987654321987654321", BigInt, false, nil,
		"-987654321987654321", Int64, false, nil},
	{`-3.14`,
		"-3.14", BigFloat, false, nil,
		"-3.14", Float64, false, nil},
	{`987654321987654321.987654321987654321`,
		"9.87654321987654321987654321987654321e+17", BigFloat, false, nil,
		"9.876543219876543e+17", Float64, false, nil},
	{`-0.987654321987654321987654321987654321E-69`,
		"-9.87654321987654321987654321987654321e-70", BigFloat, false, nil,
		"-9.876543219876543e-70", Float64, false, nil},
	{`-987654321987654321987654321987654321e69`,
		"-9.87654321987654321987654321987654321e+104", BigFloat, false, nil,
		"-9.876543219876543e+104", Float64, false, nil},
	{`{ "foo": "bar" }`,
		"nil", Nil, true, ErrNotImplemented,
		"nil", Nil, true, ErrNotImplemented},
	{`[ 123, 456 ]`,
		"nil", Nil, true, ErrNotImplemented,
		"nil", Nil, true, ErrNotImplemented},
	{`-123,456`,
		"nil", Nil, true, ErrInvalidJSON,
		"nil", Nil, true, ErrInvalidJSON},
	{`+123456`,
		"nil", Nil, true, ErrInvalidJSON,
		"nil", Nil, true, ErrInvalidJSON},
	{`.123456`,
		"nil", Nil, true, ErrInvalidJSON,
		"nil", Nil, true, ErrInvalidJSON},
	{`-123,456.789`,
		"nil", Nil, true, ErrInvalidJSON,
		"nil", Nil, true, ErrInvalidJSON},
	{`+123456.789`,
		"nil", Nil, true, ErrInvalidJSON,
		"nil", Nil, true, ErrInvalidJSON},
	{`-123,456e7`,
		"nil", Nil, true, ErrInvalidJSON,
		"nil", Nil, true, ErrInvalidJSON},
	{`+123456e7`,
		"nil", Nil, true, ErrInvalidJSON,
		"nil", Nil, true, ErrInvalidJSON},
	{`0123456`,
		"nil", Nil, true, ErrInvalidJSON,
		"nil", Nil, true, ErrInvalidJSON},
	{`-0123456`,
		"nil", Nil, true, ErrInvalidJSON,
		"nil", Nil, true, ErrInvalidJSON},
}

func TestBigDecodeJSONValue(t *testing.T) {
	for idx, rec := range testList {
		bjv := BigJSONValue{}
		_, err := bjv.DecodeJSONValue(rec.jsonStr)
		if err != rec.bigErr {
			t.Errorf("%d: Unexpected err=%+v, bigErr=%s, testRec=%+v\n", idx, err, rec.bigErr, rec)
		}
		if bjv.Kind() != rec.bigKind {
			t.Errorf("%d: Unexpected Kind()=%s, testRec=%+v\n", idx, bjv.Kind(), rec)
		}
		if bjv.Kind() == Nil && !bjv.IsNil() {
			t.Errorf("%d: Unexpected IsNil()=%t, testRec=%+v\n", idx, bjv.IsNil(), rec)
		}
		if bjv.Kind() == Bool && !bjv.IsBool() {
			t.Errorf("%d: Unexpected IsBool()=%t, testRec=%+v\n", idx, bjv.IsBool(), rec)
		}
		if bjv.Kind() == String && !bjv.IsString() {
			t.Errorf("%d: Unexpected IsString()=%t, testRec=%+v\n", idx, bjv.IsString(), rec)
		}
		if bjv.Kind() == BigInt && !bjv.IsBigInt() {
			t.Errorf("%d: Unexpected IsBigInt()=%t, testRec=%+v\n", idx, bjv.IsBigInt(), rec)
		}
		if bjv.Kind() == BigFloat && !bjv.IsBigFloat() {
			t.Errorf("%d: Unexpected IsBigFloat()=%t, testRec=%+v\n", idx, bjv.IsBigFloat(), rec)
		}
		if bjv.IsNil() != rec.bigIsNil {
			t.Errorf("%d: Unexpected IsNil()=%t, testRec=%+v\n", idx, bjv.IsNil(), rec)
		}
		if bjv.IsNil() && bjv.Value() != nil {
			t.Errorf("%d: Unexpected non-nil Value()=%+v, testRec=%+v\n", idx, bjv.Value(), rec)
		}
		if !bjv.IsNil() && bjv.Value() == nil {
			t.Errorf("%d: Unexpected nil Value()=%+v, testRec=%+v\n", idx, bjv.Value(), rec)
		}
		if bjv.IsBool() {
			bjv.Bool() // panics with runtime error if not a bool
		}
		if bjv.IsBigInt() {
			bjv.BigInt() // panics with runtime error if not a big.Int
		}
		if bjv.IsBigFloat() {
			bjv.BigFloat() // panics with runtime error if not a big.Float
		}
		if bjv.String() != rec.bigString {
			t.Errorf("%d: Unexpected String()=%s, testRec=%+v\n", idx, bjv.String(), rec)
		}
	}
}

type bigWalChangeRec struct {
	ColumnValues []BigJSONValue `json:"columnvalues"`
}

func TestBigUnmarshalJSON(t *testing.T) {
	jsonStr := `{ "columnvalues": [
		null,
		true,
		false,
		"a\\b\\c new\nline",
		3.14,
		-987654321987654321,
		-987654321987654321.987654321987654321
	] }`

	expectedKinds := []Kind{
		Nil,
		Bool,
		Bool,
		String,
		BigFloat,
		BigInt,
		BigFloat,
	}

	var bigRec bigWalChangeRec
	err := json.Unmarshal([]byte(jsonStr), &bigRec)
	if err != nil {
		t.Errorf("json.Unmarshal err=%s\n", err)
	} else {
		for idx, bjv := range bigRec.ColumnValues {
			if bjv.Kind() != expectedKinds[idx] {
				t.Errorf("json.Unmarshal BigJSONValue <%s> is of kind %s, does not match expectedKinds[%d]=%s",
					bjv.String(), bjv.Kind(), idx, expectedKinds[idx])
			}
		}
	}
}

func BenchmarkBigDecodeJSONNumbers(b *testing.B) {
	bjv := BigJSONValue{}
	for n := 0; n < b.N; n++ {
		bjv.DecodeJSONValue(`987654321987654321`)
		bjv.DecodeJSONValue(`987654321.987654321`)
	}
}
