package logger

import (
	"fmt"
	"io"
)

type StackWriter struct {
	w          io.Writer
	buf        [bufSize]byte
	bufDataLen int
}

const bufSize = 1024

func MakeStackWriter(w io.Writer) StackWriter {
	return StackWriter{
		w: w,
	}
}

func (sw *StackWriter) WriteEscaped(str string) (n int, err error) {
	var copyFrom int
	for i, c := range str {
		switch c {
		case 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
			17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
			'\\', '"':
			if copyFrom < i {
				nw, err := sw.Write(str[copyFrom:i])
				copyFrom += nw
				n += nw
				if err != nil {
					return n, err
				}
			}

			var sliceToWrite []byte
			switch c {
			case 0, 1, 2, 3, 4, 5, 6, 7, 8:
				sliceToWrite = []byte(`\u0000`)
				sliceToWrite[5] = byte('0' + c)
			case 11, 12, 14, 15:
				sliceToWrite = []byte(`\u0000`)
				sliceToWrite[5] = byte('a' + c - 10)
			case 16, 17, 18, 19, 20, 21, 22, 23, 24, 25:
				sliceToWrite = []byte(`\u0010`)
				sliceToWrite[5] = byte('0' + c - 16)
			case 26, 27, 28, 29, 30, 31:
				sliceToWrite = []byte(`\u0010`)
				sliceToWrite[5] = byte('a' + c - 26)
			case '\r':
				sliceToWrite = []byte(`\r`)
			case '\n':
				sliceToWrite = []byte(`\n`)
			case '\t':
				sliceToWrite = []byte(`\t`)
			case '\\':
				sliceToWrite = []byte(`\\`)
			case '"':
				sliceToWrite = []byte(`\"`)
			}

			nw, err := sw.Write(string(sliceToWrite))
			if err != nil {
				return n, err
			}
			if nw != len(sliceToWrite) {
				return n, fmt.Errorf("failed to write %d bytes, wrote %d bytes", len(sliceToWrite), nw)
			}
			copyFrom++
			n++
		}
	}

	if copyFrom < len(str) {
		nw, err := sw.Write(str[copyFrom:])
		n += nw
		if err != nil {
			return n, err
		}
	}

	return n, nil
}

func (sw *StackWriter) Write(s string) (n int, err error) {
	lenToWrite := len(s)
	for lenToWrite > 0 {
		endIndex := sw.bufDataLen + lenToWrite
		if endIndex > bufSize {
			copyLen := bufSize - sw.bufDataLen
			bytesCopied := copy(sw.buf[sw.bufDataLen:], s[:copyLen])
			sw.bufDataLen += copyLen
			n += bytesCopied
			if bytesCopied != copyLen {
				return n, fmt.Errorf("failed to copy %d chars", copyLen)
			}

			err := sw.Flush()
			if err != nil {
				return 0, fmt.Errorf("failed to write: %w", err)
			}

			s = s[copyLen:]
			lenToWrite = len(s)
		} else {
			bytesCopied := copy(sw.buf[sw.bufDataLen:], s)
			sw.bufDataLen += bytesCopied
			n += bytesCopied
			if bytesCopied != lenToWrite {
				return 0, fmt.Errorf("failed to copy %d chars", lenToWrite)
			}

			break
		}
	}

	return n, err
}

func (sw *StackWriter) Flush() error {
	if sw.bufDataLen == 0 {
		return nil
	}

	bytesToFlush := sw.bufDataLen
	bytesFlushed, err := sw.rawFlush()
	if err != nil {
		return err
	}
	if bytesFlushed != bytesToFlush {
		return fmt.Errorf("flushed only %d bytes, but buffer contained %d bytes", bytesFlushed, bytesToFlush)
	}

	return nil
}

func (sw *StackWriter) rawFlush() (n int, err error) {
	if sw.bufDataLen == 0 {
		return 0, nil
	}

	// We need a slice var for the unsafe pointer:
	slice := sw.buf[:sw.bufDataLen]

	// Reset the bufDataLen in any case, even if not all bytes were written.
	sw.bufDataLen = 0

	return sw.w.Write(noescape_bytearray(&slice))
}
