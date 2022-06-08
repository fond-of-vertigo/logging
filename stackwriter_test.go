package logger

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"
)

func TestStackWriter_Write(t *testing.T) {
	tests := []struct {
		name    string
		writer  StackWriter
		msg     string
		wantMsg string
		wantErr bool
	}{{
		name:    "Write empty string ",
		msg:     "",
		wantMsg: "",
	}, {
		name:    "Write one short line",
		msg:     "ABC",
		wantMsg: "ABC",
	}, {
		name:    "Write unicode text",
		msg:     "ABC äöü 🙂 abc",
		wantMsg: "ABC äöü 🙂 abc",
	}, {
		name:    "Write one line with max length",
		msg:     makeString(bufSize),
		wantMsg: makeString(bufSize),
	}, {
		name:    "Write one line with max+1 length",
		msg:     makeString(bufSize + 1),
		wantMsg: makeString(bufSize + 1),
	}, {
		name:    "Write one line with max*2 length",
		msg:     makeString(bufSize * 2),
		wantMsg: makeString(bufSize * 2),
	}, {
		name:    "Write one line with max*2+1 length",
		msg:     makeString((bufSize * 2) + 1),
		wantMsg: makeString((bufSize * 2) + 1),
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := bytes.NewBufferString("")
			tt.writer = MakeStackWriter(out)

			gotN, err := tt.writer.Write(tt.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotN != len(tt.msg) {
				t.Errorf("Write() gotN = %v, want %v", gotN, len(tt.msg))
			}

			err = tt.writer.Flush()
			if (err != nil) != tt.wantErr {
				t.Errorf("Flush() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if out.String() != tt.wantMsg {
				t.Errorf("out.String does not equal wantString:\nWant: %s\nGot.: %s", tt.wantMsg, out.String())
			}
		})
	}
}

func TestStackWriter_WriteEscaped(t *testing.T) {
	tests := []struct {
		name     string
		writer   StackWriter
		strValue string
		wantJSON string
		wantErr  bool
	}{{
		name:     "Marshal simple string value",
		strValue: `abcdefgh`,
	}, {
		name:     "Marshal unicode string value",
		strValue: `abcdefgh äöü 🙂 abc`,
	}, {
		name:     "Marshal control chars",
		strValue: string([]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31}),
	}, {
		name:     "Marshal special chars",
		strValue: "\\\"",
	}, {
		name:     "Marshal mixed text 1",
		strValue: "\"abc\\def\"",
	}, {
		name:     "Marshal mixed text 2",
		strValue: "\"a\"",
	}, {
		name:     "Marshal mixed text 3",
		strValue: "\"a\"a",
	}, {
		name:     "Marshal multi-line text",
		strValue: makeString(32) + "\n" + makeString(32),
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := bytes.NewBufferString("")
			tt.writer = MakeStackWriter(out)

			_, err := tt.writer.WriteEscaped(tt.strValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("WriteString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			err = tt.writer.Flush()
			if (err != nil) != tt.wantErr {
				t.Errorf("Flush() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			tt.wantJSON = mustMarshalJSONString(tt.strValue)
			if out.String() != tt.wantJSON {
				t.Errorf("out.String does not equal wantJSON:\nWant: \"%s\"\nGot.: \"%s\"", tt.wantJSON, out.String())
			}
		})
	}
}

func TestStackWriter_ZeroAlloc(t *testing.T) {
	longString := makeString(16 * 1024)
	allocs := testing.AllocsPerRun(1, func() {
		sw := MakeStackWriter(os.Stdout)
		n, err := sw.Write(longString)
		if err != nil {
			t.Errorf("WriteString error: %s", err)
		}
		if n != len(longString) {
			t.Errorf("wrote %d bytes, expected %d bytes", n, len(longString))
		}
	})

	if allocs > 0.0 {
		t.Errorf("Allocs detected! Want 0 allocs, got %f", allocs)
	}
}

func TestStackWriter_WritesEscaped_ZeroAlloc(t *testing.T) {
	longString := makeString(16 * 1024)
	allocs := testing.AllocsPerRun(1, func() {
		sw := MakeStackWriter(io.Discard)
		n, err := sw.WriteJSONString(longString)
		if err != nil {
			t.Errorf("WriteString error: %s", err)
		}
		if n != len(longString)+2 {
			t.Errorf("wrote %d bytes, expected %d bytes", n, len(longString)+2)
		}
	})

	if allocs > 0.0 {
		t.Errorf("Allocs detected! Want 0 allocs, got %f", allocs)
	}
}

func mustMarshalJSONString(str string) string {
	j, err := json.Marshal(str)
	if err != nil {
		panic(err)
	}
	return string(j[1 : len(j)-1])
}

func makeString(length int) string {
	var sb strings.Builder
	for i := 0; i < length; i++ {
		sb.WriteByte(byte('0' + (i % 10)))
	}
	return sb.String()
}
