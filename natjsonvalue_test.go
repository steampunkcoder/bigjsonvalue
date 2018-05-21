package bigjsonvalue

import (
	"encoding/json"
	"testing"
)

func TestNatDefaultValueIsNil(t *testing.T) {
	njvN := NatJSONValue{}
	if !njvN.IsNil() {
		t.Errorf("njvN is not nil NatJSONValue")
	}
	if njvN.Kind() != Nil {
		t.Errorf("njvN.Kind() unexpectedly returned %s", njvN.Kind())
	}

	njvM := new(NatJSONValue)
	if !njvM.IsNil() {
		t.Errorf("njvM is not nil NatJSONValue")
	}
	if njvM.Kind() != Nil {
		t.Errorf("njvM.Kind() unexpectedly returned %s", njvM.Kind())
	}
}

func TestNatDecodeJSONValue(t *testing.T) {
	for idx, rec := range testList {
		njv := NatJSONValue{}
		_, err := njv.DecodeJSONValue(rec.jsonStr)
		if err != rec.natErr {
			t.Errorf("%d: Unexpected err=%+v, natErr=%s, testRec=%+v\n", idx, err, rec.natErr, rec)
		}
		if njv.Kind() != rec.natKind {
			t.Errorf("%d: Unexpected Kind()=%s, testRec=%+v\n", idx, njv.Kind(), rec)
		}
		if njv.Kind() == Nil && !njv.IsNil() {
			t.Errorf("%d: Unexpected IsNil()=%t, testRec=%+v\n", idx, njv.IsNil(), rec)
		}
		if njv.Kind() == Bool && !njv.IsBool() {
			t.Errorf("%d: Unexpected IsBool()=%t, testRec=%+v\n", idx, njv.IsBool(), rec)
		}
		if njv.Kind() == String && !njv.IsString() {
			t.Errorf("%d: Unexpected IsString()=%t, testRec=%+v\n", idx, njv.IsString(), rec)
		}
		if njv.Kind() == Int64 && !njv.IsInt64() {
			t.Errorf("%d: Unexpected IsInt64()=%t, testRec=%+v\n", idx, njv.IsInt64(), rec)
		}
		if njv.Kind() == Uint64 && !njv.IsUint64() {
			t.Errorf("%d: Unexpected IsUint64()=%t, testRec=%+v\n", idx, njv.IsUint64(), rec)
		}
		if njv.Kind() == Float64 && !njv.IsFloat64() {
			t.Errorf("%d: Unexpected IsFloat64()=%t, testRec=%+v\n", idx, njv.IsFloat64(), rec)
		}
		if njv.IsNil() != rec.natIsNil {
			t.Errorf("%d: Unexpected IsNil()=%t, testRec=%+v\n", idx, njv.IsNil(), rec)
		}
		if njv.IsNil() && njv.Value() != nil {
			t.Errorf("%d: Unexpected non-nil Value()=%+v, testRec=%+v\n", idx, njv.Value(), rec)
		}
		if !njv.IsNil() && njv.Value() == nil {
			t.Errorf("%d: Unexpected nil Value()=%+v, testRec=%+v\n", idx, njv.Value(), rec)
		}
		if njv.IsBool() {
			njv.Bool() // panics with runtime error if not a bool
		}
		if njv.IsInt64() {
			njv.Int64() // panics with runtime error if not a int64
		}
		if njv.IsUint64() {
			njv.Uint64() // panics with runtime error if not a uint64
		}
		if njv.IsFloat64() {
			njv.Float64() // panics with runtime error if not a float64
		}
		if njv.String() != rec.natString {
			t.Errorf("%d: Unexpected String()=%s, testRec=%+v\n", idx, njv.String(), rec)
		}
	}
}

type natWalChangeRec struct {
	ColumnValues []NatJSONValue `json:"columnvalues"`
}

func TestNatUnmarshalJSON(t *testing.T) {
	jsonStr := `{ "columnvalues": [
		null,
		true,
		false,
		"a\\b\\c new\nline",
		987654321987654321,
		3.14,
		-987654321987654321,
		-987654321.987654321
	] }`

	expectedKinds := []Kind{
		Nil,
		Bool,
		Bool,
		String,
		Uint64,
		Float64,
		Int64,
		Float64,
	}

	var bigRec natWalChangeRec
	err := json.Unmarshal([]byte(jsonStr), &bigRec)
	if err != nil {
		t.Errorf("json.Unmarshal err=%s\n", err)
	} else {
		for idx, njv := range bigRec.ColumnValues {
			if njv.Kind() != expectedKinds[idx] {
				t.Errorf("json.Unmarshal NatJSONValue <%s> is of kind %s, does not match expectedKinds[%d]=%s",
					njv.String(), njv.Kind(), idx, expectedKinds[idx])
			}
		}
	}
}

func BenchmarkNatDecodeJSONNumbers(b *testing.B) {
	njv := NatJSONValue{}
	for n := 0; n < b.N; n++ {
		njv.DecodeJSONValue(`987654321987654321`)
		njv.DecodeJSONValue(`987654321.987654321`)
	}
}
