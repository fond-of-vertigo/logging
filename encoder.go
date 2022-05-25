package logger

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// JSONValueWriter can be implemented to write values types that are not directly supported
// by the logger lib.
type JSONValueWriter interface {
	// WriteJSONValue must be implemented to write a JSON value. The implementation is in
	// full control how to write the JSON value. Call sw.WriteJSONString for a full string value
	// or sw.WriteEscaped for partial string content that is added by multiple calls.
	// Call sw.Write for numbers or boolean values.
	WriteJSONValue(sw *StackWriter) (n int, err error)
}

func encodeKey(sw *StackWriter, key interface{}) (n int, err error) {
	switch k := key.(type) {
	case string:
		return sw.WriteJSONString(noescape_string(&k))
	case fmt.Stringer:
		return sw.WriteJSONString(noescape_stringer(&k).String())
	default:
		return sw.WriteJSONString(fmt.Sprintf("INVALID_KEY_%v", noescape_interface(&k)))
	}
}

func encodeValue(sw *StackWriter, value interface{}) (n int, err error) {
	switch v := value.(type) {
	case string:
		return sw.WriteJSONString(noescape_string(&v))
	case float32:
		return sw.Write(string(strconv.AppendFloat(nil, float64(v), 'f', 6, 32)))
	case float64:
		return sw.Write(string(strconv.AppendFloat(nil, v, 'f', 6, 64)))
	case int:
		return sw.Write(strconv.Itoa(v))
	case int32:
		return sw.Write(strconv.Itoa(int(v)))
	case int64:
		return sw.Write(strconv.FormatInt(v, 10))
	case uint:
		return sw.Write(strconv.FormatUint(uint64(v), 10))
	case uint64:
		return sw.Write(strconv.FormatUint(v, 10))
	case bool:
		return sw.Write(strconv.FormatBool(v))
	case fmt.Stringer:
		return sw.WriteJSONString(noescape_stringer(&v).String())
	case JSONValueWriter:
		return v.WriteJSONValue(noescape_stackwriterptr(sw))
	default:
		jsonString, err := json.Marshal(noescape_interface(&v))
		if err != nil {
			return 0, err
		}
		return sw.Write(string(jsonString))
	}
}
