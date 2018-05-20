package bigjsonvalue

import (
	"encoding/json"
	"reflect"

	"testing"
)

func TestDefaultValueIsNil(t *testing.T) {
	bjvN := BigJSONValue{}
	if !bjvN.IsNil() {
		t.Errorf("bjvN is not nil BigJSONValue")
	}
	if bjvN.Kind() != reflect.Invalid {
		t.Errorf("bjvN.Kind() unexpectedly returned %s", bjvN.Kind())
	}

	bjvM := new(BigJSONValue)
	if !bjvM.IsNil() {
		t.Errorf("bjvM is not nil BigJSONValue")
	}
	if bjvM.Kind() != reflect.Invalid {
		t.Errorf("bjvM.Kind() unexpectedly returned %s", bjvM.Kind())
	}
}

func TestDecodeJSONValue(t *testing.T) {
	bjvN := BigJSONValue{}
	_, err := bjvN.DecodeJSONValue(`null`)
	if err != nil {
		t.Errorf("DecodeJSONValue of null string error: %s", err)
	}
	if bjvN.Value() != nil {
		t.Errorf("bjvN.Value() should be nil")
	}
	if bjvN.String() != "nil" {
		t.Errorf("DecodeJSONValue of null string failed: <%s>", bjvN.String())
	}
	if bjvN.Kind() != reflect.Invalid {
		t.Errorf("bjvN.Kind() unexpectedly returned %s", bjvN.Kind())
	}
	if !bjvN.IsNil() {
		t.Errorf("bjvN.IsNil() did not return true")
	}

	bjvB := BigJSONValue{}
	_, err = bjvB.DecodeJSONValue(`true`)
	if err != nil {
		t.Errorf("DecodeJSONValue of true bool string error: %s", err)
	}
	if bjvB.Value() == nil {
		t.Errorf("bjvB.Value() should not be nil")
	}
	if bjvB.String() != "true" {
		t.Errorf("DecodeJSONValue of true bool string failed: <%s>", bjvB.String())
	}
	if bjvB.Kind() != reflect.Bool {
		t.Errorf("bjvB.Kind() unexpectedly returned %s", bjvB.Kind())
	}
	if !bjvB.IsBool() {
		t.Errorf("bjvB.IsBool() did not return true")
	}
	_ = bjvB.Bool() // panics with runtime error if not a bool

	bjvC := BigJSONValue{}
	_, err = bjvC.DecodeJSONValue(`false`)
	if err != nil {
		t.Errorf("DecodeJSONValue of false bool string error: %s", err)
	}
	if bjvC.Value() == nil {
		t.Errorf("bjvC.Value() should not be nil")
	}
	if bjvC.String() != "false" {
		t.Errorf("DecodeJSONValue of false bool string failed: <%s>", bjvC.String())
	}
	if bjvC.Kind() != reflect.Bool {
		t.Errorf("bjvC.Kind() unexpectedly returned %s", bjvC.Kind())
	}
	if !bjvC.IsBool() {
		t.Errorf("bjvC.IsBool() did not return true")
	}
	_ = bjvB.Bool() // panics with runtime error if not a bool

	bjvS := BigJSONValue{}
	_, err = bjvS.DecodeJSONValue(`"a\\b\\c new\nline"`)
	if err != nil {
		t.Errorf("DecodeJSONValue of backslashed JSON string error: %s", err)
	}
	if bjvS.Value() == nil {
		t.Errorf("bjvS.Value() should not be nil")
	}
	if bjvS.String() != "a\\b\\c new\nline" {
		t.Errorf("DecodeJSONValue of backslashed JSON string failed: <%s>", bjvS.String())
	}
	if bjvS.Kind() != reflect.String {
		t.Errorf("bjvS.Kind() unexpectedly returned %s", bjvS.Kind())
	}
	if !bjvS.IsString() {
		t.Errorf("bjvS.IsString() did not return true")
	}

	bjvI := BigJSONValue{}
	_, err = bjvI.DecodeJSONValue(`987654321987654321`)
	if err != nil {
		t.Errorf("DecodeJSONValue of large positive integer string error: %s", err)
	}
	if bjvI.Value() == nil {
		t.Errorf("bjvI.Value() should not be nil")
	}
	if bjvI.String() != "987654321987654321" {
		t.Errorf("DecodeJSONValue of large positive integer string failed: <%s>", bjvI.String())
	}
	if bjvI.Kind() != reflect.Int {
		t.Errorf("bjvI.Kind() unexpectedly returned %s", bjvI.Kind())
	}
	if !bjvI.IsBigInt() {
		t.Errorf("bjvI.IsBigInt() did not return true")
	}
	_ = bjvI.BigInt() // panics with runtime error if not a big.Int

	bjvJ := BigJSONValue{}
	_, err = bjvJ.DecodeJSONValue(`-987654321987654321`)
	if err != nil {
		t.Errorf("DecodeJSONValue of large negative integer string error: %s", err)
	}
	if bjvJ.Value() == nil {
		t.Errorf("bjvJ.Value() should not be nil")
	}
	if bjvJ.String() != "-987654321987654321" {
		t.Errorf("DecodeJSONValue of large negative integer string failed: <%s>", bjvJ.String())
	}
	if bjvJ.Kind() != reflect.Int {
		t.Errorf("bjvJ.Kind() unexpectedly returned %s", bjvJ.Kind())
	}
	if !bjvJ.IsBigInt() {
		t.Errorf("bjvJ.IsBigInt() did not return true")
	}
	_ = bjvJ.BigInt() // panics with runtime error if not a big.Int

	bjvF := BigJSONValue{}
	_, err = bjvF.DecodeJSONValue(`-3.14`)
	if err != nil {
		t.Errorf("DecodeJSONValue of negative float string error: %s", err)
	}
	if bjvF.Value() == nil {
		t.Errorf("bjvF.Value() should not be nil")
	}
	if bjvF.String() != "-3.14" {
		t.Errorf("DecodeJSONValue of negative float string failed: <%s>", bjvF.String())
	}
	if bjvF.Kind() != reflect.Float64 {
		t.Errorf("bjvF.Kind() unexpectedly returned %s", bjvF.Kind())
	}
	if !bjvF.IsBigFloat() {
		t.Errorf("bjvF.IsBigFloat() did not return true")
	}
	_ = bjvF.BigFloat() // panics with runtime error if not a big.Float

	bjvG := BigJSONValue{}
	_, err = bjvG.DecodeJSONValue(`987654321987654321.987654321987654321`)
	if err != nil {
		t.Errorf("DecodeJSONValue of large positive float string error: %s", err)
	}
	if bjvG.Value() == nil {
		t.Errorf("bjvG.Value() should not be nil")
	}
	if bjvG.String() != "9.87654321987654321987654321987654321e+17" {
		t.Errorf("DecodeJSONValue of large positive float string failed: <%s>", bjvG.String())
	}
	if bjvG.Kind() != reflect.Float64 {
		t.Errorf("bjvG.Kind() unexpectedly returned %s", bjvG.Kind())
	}
	if !bjvG.IsBigFloat() {
		t.Errorf("bjvG.IsBigFloat() did not return true")
	}
	_ = bjvG.BigFloat() // panics with runtime error if not a big.Float

	bjvH := BigJSONValue{}
	_, err = bjvH.DecodeJSONValue(`-0.987654321987654321987654321987654321E-69`)
	if err != nil {
		t.Errorf("DecodeJSONValue of small negative float string error: %s", err)
	}
	if bjvH.Value() == nil {
		t.Errorf("bjvH.Value() should not be nil")
	}
	if bjvH.String() != "-9.87654321987654321987654321987654321e-70" {
		t.Errorf("DecodeJSONValue of small negative float string failed: <%s>", bjvH.String())
	}
	if bjvH.Kind() != reflect.Float64 {
		t.Errorf("bjvH.Kind() unexpectedly returned %s", bjvH.Kind())
	}
	if !bjvH.IsBigFloat() {
		t.Errorf("bjvH.IsBigFloat() did not return true")
	}
	_ = bjvH.BigFloat() // panics with runtime error if not a big.Float

	bjvX := BigJSONValue{}
	_, err = bjvX.DecodeJSONValue(`-987654321987654321987654321987654321e69`)
	if err != nil {
		t.Errorf("DecodeJSONValue of large negative float string error: %s", err)
	}
	if bjvX.Value() == nil {
		t.Errorf("bjvX.Value() should not be nil")
	}
	if bjvX.String() != "-9.87654321987654321987654321987654321e+104" {
		t.Errorf("DecodeJSONValue of large negative float string failed: <%s>", bjvX.String())
	}
	if bjvX.Kind() != reflect.Float64 {
		t.Errorf("bjvX.Kind() unexpectedly returned %s", bjvX.Kind())
	}
	if !bjvX.IsBigFloat() {
		t.Errorf("bjvX.IsBigFloat() did not return true")
	}
	_ = bjvX.BigFloat() // panics with runtime error if not a big.Float

	bjvE := BigJSONValue{}
	_, err = bjvE.DecodeJSONValue(`{ foo: bar }`)
	if err != ErrNotImplemented {
		t.Errorf("DecodeJSONValue of braces string unexpected error: %s", err)
	}

	_, err = bjvE.DecodeJSONValue(`[ foo bar ]`)
	if err != ErrNotImplemented {
		t.Errorf("DecodeJSONValue of sq brackets string unexpected error: %s", err)
	}

	_, err = bjvE.DecodeJSONValue(`-123,456`)
	if err != ErrInvalidJSON {
		t.Errorf("DecodeJSONValue of non-number string unexpected error: %s", err)
	}

	_, err = bjvE.DecodeJSONValue(`+123456`)
	if err != ErrInvalidJSON {
		t.Errorf("DecodeJSONValue of non-number string unexpected error: %s", err)
	}

	_, err = bjvE.DecodeJSONValue(`.123456`)
	if err != ErrInvalidJSON {
		t.Errorf("DecodeJSONValue of non-number string unexpected error: %s", err)
	}

	_, err = bjvE.DecodeJSONValue(`-123,456.789`)
	if err != ErrInvalidJSON {
		t.Errorf("DecodeJSONValue of non-number string unexpected error: %s", err)
	}

	_, err = bjvE.DecodeJSONValue(`+123456.789`)
	if err != ErrInvalidJSON {
		t.Errorf("DecodeJSONValue of non-number string unexpected error: %s", err)
	}

	_, err = bjvE.DecodeJSONValue(`-123,456e7`)
	if err != ErrInvalidJSON {
		t.Errorf("DecodeJSONValue of non-number string unexpected error: %s", err)
	}

	_, err = bjvE.DecodeJSONValue(`+123456e7`)
	if err != ErrInvalidJSON {
		t.Errorf("DecodeJSONValue of non-number string unexpected error: %s", err)
	}
}

type WalChangeRec struct {
	ColumnValues []BigJSONValue `json:"columnvalues"`
}

func TestUnmarshalJSON(t *testing.T) {
	jsonStr := `{ "columnvalues": [
		null,
		true,
		false,
		"a\\b\\c new\nline",
		3.14,
		-987654321987654321,
		-987654321987654321.987654321987654321
	] }`

	expectedKinds := []reflect.Kind{
		reflect.Invalid,
		reflect.Bool,
		reflect.Bool,
		reflect.String,
		reflect.Float64,
		reflect.Int,
		reflect.Float64,
	}

	var bigRec WalChangeRec
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
