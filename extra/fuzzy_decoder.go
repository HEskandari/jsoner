package extra

import (
	"encoding/json"
	"io"
	"math"
	"reflect"
	"strings"
	"unsafe"

	"github.com/heskandari/jsoner"
	"github.com/modern-go/reflect2"
)

const maxUint = ^uint(0)
const maxInt = int(maxUint >> 1)
const minInt = -maxInt - 1

// RegisterFuzzyDecoders decode input from PHP with tolerance.
// It will handle string/number auto conversation, and treat empty [] as empty struct.
func RegisterFuzzyDecoders() {
	api := jsoner.DefaultAPI()
	api.RegisterExtension(&tolerateEmptyArrayExtension{})
	api.RegisterTypeDecoder("string", &fuzzyStringDecoder{})
	api.RegisterTypeDecoder("float32", &fuzzyFloat32Decoder{})
	api.RegisterTypeDecoder("float64", &fuzzyFloat64Decoder{})
	api.RegisterTypeDecoder("int", &fuzzyIntegerDecoder{func(isFloat bool, ptr unsafe.Pointer, iter *jsoner.Iterator) {
		if isFloat {
			val := iter.ReadFloat64()
			if val > float64(maxInt) || val < float64(minInt) {
				iter.ReportError("fuzzy decode int", "exceed range")
				return
			}
			*((*int)(ptr)) = int(val)
		} else {
			*((*int)(ptr)) = iter.ReadInt()
		}
	}})
	api.RegisterTypeDecoder("uint", &fuzzyIntegerDecoder{func(isFloat bool, ptr unsafe.Pointer, iter *jsoner.Iterator) {
		if isFloat {
			val := iter.ReadFloat64()
			if val > float64(maxUint) || val < 0 {
				iter.ReportError("fuzzy decode uint", "exceed range")
				return
			}
			*((*uint)(ptr)) = uint(val)
		} else {
			*((*uint)(ptr)) = iter.ReadUint()
		}
	}})
	api.RegisterTypeDecoder("int8", &fuzzyIntegerDecoder{func(isFloat bool, ptr unsafe.Pointer, iter *jsoner.Iterator) {
		if isFloat {
			val := iter.ReadFloat64()
			if val > float64(math.MaxInt8) || val < float64(math.MinInt8) {
				iter.ReportError("fuzzy decode int8", "exceed range")
				return
			}
			*((*int8)(ptr)) = int8(val)
		} else {
			*((*int8)(ptr)) = iter.ReadInt8()
		}
	}})
	api.RegisterTypeDecoder("uint8", &fuzzyIntegerDecoder{func(isFloat bool, ptr unsafe.Pointer, iter *jsoner.Iterator) {
		if isFloat {
			val := iter.ReadFloat64()
			if val > float64(math.MaxUint8) || val < 0 {
				iter.ReportError("fuzzy decode uint8", "exceed range")
				return
			}
			*((*uint8)(ptr)) = uint8(val)
		} else {
			*((*uint8)(ptr)) = iter.ReadUint8()
		}
	}})
	api.RegisterTypeDecoder("int16", &fuzzyIntegerDecoder{func(isFloat bool, ptr unsafe.Pointer, iter *jsoner.Iterator) {
		if isFloat {
			val := iter.ReadFloat64()
			if val > float64(math.MaxInt16) || val < float64(math.MinInt16) {
				iter.ReportError("fuzzy decode int16", "exceed range")
				return
			}
			*((*int16)(ptr)) = int16(val)
		} else {
			*((*int16)(ptr)) = iter.ReadInt16()
		}
	}})
	api.RegisterTypeDecoder("uint16", &fuzzyIntegerDecoder{func(isFloat bool, ptr unsafe.Pointer, iter *jsoner.Iterator) {
		if isFloat {
			val := iter.ReadFloat64()
			if val > float64(math.MaxUint16) || val < 0 {
				iter.ReportError("fuzzy decode uint16", "exceed range")
				return
			}
			*((*uint16)(ptr)) = uint16(val)
		} else {
			*((*uint16)(ptr)) = iter.ReadUint16()
		}
	}})
	api.RegisterTypeDecoder("int32", &fuzzyIntegerDecoder{func(isFloat bool, ptr unsafe.Pointer, iter *jsoner.Iterator) {
		if isFloat {
			val := iter.ReadFloat64()
			if val > float64(math.MaxInt32) || val < float64(math.MinInt32) {
				iter.ReportError("fuzzy decode int32", "exceed range")
				return
			}
			*((*int32)(ptr)) = int32(val)
		} else {
			*((*int32)(ptr)) = iter.ReadInt32()
		}
	}})
	api.RegisterTypeDecoder("uint32", &fuzzyIntegerDecoder{func(isFloat bool, ptr unsafe.Pointer, iter *jsoner.Iterator) {
		if isFloat {
			val := iter.ReadFloat64()
			if val > float64(math.MaxUint32) || val < 0 {
				iter.ReportError("fuzzy decode uint32", "exceed range")
				return
			}
			*((*uint32)(ptr)) = uint32(val)
		} else {
			*((*uint32)(ptr)) = iter.ReadUint32()
		}
	}})
	api.RegisterTypeDecoder("int64", &fuzzyIntegerDecoder{func(isFloat bool, ptr unsafe.Pointer, iter *jsoner.Iterator) {
		if isFloat {
			val := iter.ReadFloat64()
			if val > float64(math.MaxInt64) || val < float64(math.MinInt64) {
				iter.ReportError("fuzzy decode int64", "exceed range")
				return
			}
			*((*int64)(ptr)) = int64(val)
		} else {
			*((*int64)(ptr)) = iter.ReadInt64()
		}
	}})
	api.RegisterTypeDecoder("uint64", &fuzzyIntegerDecoder{func(isFloat bool, ptr unsafe.Pointer, iter *jsoner.Iterator) {
		if isFloat {
			val := iter.ReadFloat64()
			if val > float64(math.MaxUint64) || val < 0 {
				iter.ReportError("fuzzy decode uint64", "exceed range")
				return
			}
			*((*uint64)(ptr)) = uint64(val)
		} else {
			*((*uint64)(ptr)) = iter.ReadUint64()
		}
	}})
}

type tolerateEmptyArrayExtension struct {
	jsoner.DummyExtension
}

func (extension *tolerateEmptyArrayExtension) DecorateDecoder(typ reflect2.Type, decoder jsoner.ValDecoder) jsoner.ValDecoder {
	if typ.Kind() == reflect.Struct || typ.Kind() == reflect.Map {
		return &tolerateEmptyArrayDecoder{decoder}
	}
	return decoder
}

type tolerateEmptyArrayDecoder struct {
	valDecoder jsoner.ValDecoder
}

func (decoder *tolerateEmptyArrayDecoder) Decode(ptr unsafe.Pointer, iter *jsoner.Iterator) {
	if iter.WhatIsNext() == jsoner.ArrayValue {
		iter.Skip()
		newIter := iter.Pool().BorrowIterator([]byte("{}"))
		defer iter.Pool().ReturnIterator(newIter)
		decoder.valDecoder.Decode(ptr, newIter)
	} else {
		decoder.valDecoder.Decode(ptr, iter)
	}
}

type fuzzyStringDecoder struct {
}

func (decoder *fuzzyStringDecoder) Decode(ptr unsafe.Pointer, iter *jsoner.Iterator) {
	valueType := iter.WhatIsNext()
	switch valueType {
	case jsoner.NumberValue:
		var number json.Number
		iter.ReadVal(&number)
		*((*string)(ptr)) = string(number)
	case jsoner.StringValue:
		*((*string)(ptr)) = iter.ReadString()
	case jsoner.NilValue:
		iter.Skip()
		*((*string)(ptr)) = ""
	default:
		iter.ReportError("fuzzyStringDecoder", "not number or string")
	}
}

type fuzzyIntegerDecoder struct {
	fun func(isFloat bool, ptr unsafe.Pointer, iter *jsoner.Iterator)
}

func (decoder *fuzzyIntegerDecoder) Decode(ptr unsafe.Pointer, iter *jsoner.Iterator) {
	valueType := iter.WhatIsNext()
	var str string
	switch valueType {
	case jsoner.NumberValue:
		var number json.Number
		iter.ReadVal(&number)
		str = string(number)
	case jsoner.StringValue:
		str = iter.ReadString()
	case jsoner.BoolValue:
		if iter.ReadBool() {
			str = "1"
		} else {
			str = "0"
		}
	case jsoner.NilValue:
		iter.Skip()
		str = "0"
	default:
		iter.ReportError("fuzzyIntegerDecoder", "not number or string")
	}
	if len(str) == 0 {
		str = "0"
	}
	newIter := iter.Pool().BorrowIterator([]byte(str))
	defer iter.Pool().ReturnIterator(newIter)
	isFloat := strings.IndexByte(str, '.') != -1
	decoder.fun(isFloat, ptr, newIter)
	if newIter.Error != nil && newIter.Error != io.EOF {
		iter.Error = newIter.Error
	}
}

type fuzzyFloat32Decoder struct {
}

func (decoder *fuzzyFloat32Decoder) Decode(ptr unsafe.Pointer, iter *jsoner.Iterator) {
	valueType := iter.WhatIsNext()
	var str string
	switch valueType {
	case jsoner.NumberValue:
		*((*float32)(ptr)) = iter.ReadFloat32()
	case jsoner.StringValue:
		str = iter.ReadString()
		newIter := iter.Pool().BorrowIterator([]byte(str))
		defer iter.Pool().ReturnIterator(newIter)
		*((*float32)(ptr)) = newIter.ReadFloat32()
		if newIter.Error != nil && newIter.Error != io.EOF {
			iter.Error = newIter.Error
		}
	case jsoner.BoolValue:
		// support bool to float32
		if iter.ReadBool() {
			*((*float32)(ptr)) = 1
		} else {
			*((*float32)(ptr)) = 0
		}
	case jsoner.NilValue:
		iter.Skip()
		*((*float32)(ptr)) = 0
	default:
		iter.ReportError("fuzzyFloat32Decoder", "not number or string")
	}
}

type fuzzyFloat64Decoder struct {
}

func (decoder *fuzzyFloat64Decoder) Decode(ptr unsafe.Pointer, iter *jsoner.Iterator) {
	valueType := iter.WhatIsNext()
	var str string
	switch valueType {
	case jsoner.NumberValue:
		*((*float64)(ptr)) = iter.ReadFloat64()
	case jsoner.StringValue:
		str = iter.ReadString()
		newIter := iter.Pool().BorrowIterator([]byte(str))
		defer iter.Pool().ReturnIterator(newIter)
		*((*float64)(ptr)) = newIter.ReadFloat64()
		if newIter.Error != nil && newIter.Error != io.EOF {
			iter.Error = newIter.Error
		}
	case jsoner.BoolValue:
		// support bool to float64
		if iter.ReadBool() {
			*((*float64)(ptr)) = 1
		} else {
			*((*float64)(ptr)) = 0
		}
	case jsoner.NilValue:
		iter.Skip()
		*((*float64)(ptr)) = 0
	default:
		iter.ReportError("fuzzyFloat64Decoder", "not number or string")
	}
}