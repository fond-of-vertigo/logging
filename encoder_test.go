package logger

import (
	"bytes"
	"testing"
)

func Test_encodeValue(t *testing.T) {
	tests := []struct {
		name       string
		value      interface{}
		wantN      int
		wantErr    bool
		wantString string
	}{{
		name:       "Encode string",
		value:      "test",
		wantString: `"test"`,
	}, {
		name:       "Encode float32",
		value:      float32(1.25),
		wantString: `1.250000`,
	}, {
		name:       "Encode float64",
		value:      float64(1.25),
		wantString: `1.250000`,
	}, {
		name:       "Encode int",
		value:      1,
		wantString: `1`,
	}, {
		name:       "Encode int64",
		value:      9223372036854775807,
		wantString: `9223372036854775807`,
	}, {
		name:       "Encode uint",
		value:      uint(4294967295),
		wantString: `4294967295`,
	}, {
		name:       "Encode uint64",
		value:      uint64(18446744073709551615),
		wantString: `18446744073709551615`,
	}, {
		name:       "Encode bool",
		value:      true,
		wantString: `true`,
	}, {
		name:       "Encode nil",
		value:      nil,
		wantString: `null`,
	}, {
		name:       "Encode fmt.Stringer",
		value:      &testStringer{value: "testFmtStringer"},
		wantString: `"testFmtStringer"`,
	}, {
		name:       "Encode with custom encoder func",
		value:      &testCustomEncoder{value: "testCustomerEncoder"},
		wantString: `"testCustomerEncoder"`,
	}, {
		name: "Encode custom type",
		value: testUnknownType{
			A: 123,
			B: "Some value",
		},
		wantString: `{"A":123,"B":"Some value"}`,
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := bytes.NewBufferString("")
			sw := MakeStackWriter(out)

			if tt.wantN == 0 {
				tt.wantN = len(tt.wantString)
			}

			bytesWritten, err := encodeValue(&sw, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("encodeValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			sw.Flush()
			outString := out.String()
			if tt.wantString != outString {
				t.Errorf("Written string does not match.\nWant: %s\nGot : %s", tt.wantString, outString)
			}

			if bytesWritten != tt.wantN {
				t.Errorf("encodeValue() did not write %d bytes, %d bytes were written", tt.wantN, bytesWritten)
				return
			}
		})
	}
}

type testStringer struct {
	value string
}

func (ts *testStringer) String() string {
	return ts.value
}

type testCustomEncoder struct {
	value string
}

func (tce *testCustomEncoder) WriteJSONValue(sw *StackWriter) (n int, err error) {
	return sw.WriteJSONString(tce.value)
}

type testUnknownType struct {
	A int
	B string
}
