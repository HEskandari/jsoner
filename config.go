package jsoner

import (
	"encoding/json"
	"fmt"
	"github.com/modern-go/reflect2"
	"io"
	"reflect"
	"sync"
	"unsafe"
)

// Config customize how the API should behave.
// The API is created from Config by Froze.
type Config struct {
	IndentionStep                 int
	MarshalFloatWith6Digits       bool
	EscapeHTML                    bool
	SortMapKeys                   bool
	UseNumber                     bool
	DisallowUnknownFields         bool
	TagKey                        string
	OnlyTaggedField               bool
	ValidateJsonRawMessage        bool
	ObjectFieldMustBeSimpleString bool
	CaseSensitive                 bool
}

// API the public interface of this package.
// Primary Marshal and Unmarshal.
type API interface {
	IteratorPool
	StreamPool
	MarshalToString(v interface{}) (string, error)
	Marshal(v interface{}) ([]byte, error)
	MarshalIndent(v interface{}, prefix, indent string) ([]byte, error)
	UnmarshalFromString(str string, v interface{}) error
	Unmarshal(data []byte, v interface{}) error
	Get(data []byte, path ...interface{}) Any
	RegisterExtension(extension Extension)
	RegisterFieldEncoderFunc(typ string, field string, fun EncoderFunc, isEmptyFunc func(unsafe.Pointer) bool)
	RegisterFieldEncoder(typ string, field string, encoder ValEncoder)
	RegisterTypeEncoder(typ string, encoder ValEncoder)
	RegisterTypeEncoderFunc(typ string, fun EncoderFunc, isEmptyFunc func(unsafe.Pointer) bool)
	RegisterTypeDecoderFunc(typ string, fun DecoderFunc)
	RegisterTypeDecoder(typ string, decoder ValDecoder)
	RegisterFieldDecoderFunc(typ string, field string, fun DecoderFunc)
	RegisterFieldDecoder(typ string, field string, decoder ValDecoder)
	ClearExtensions()
	ClearDecoders()
	ClearEncoders()
	DecoderOf(typ reflect2.Type) ValDecoder
	EncoderOf(typ reflect2.Type) ValEncoder
}

// DefaultAPI returns the default API
func DefaultAPI() API {
	var defaultConfig = Config{
		EscapeHTML: true,
	}.Froze()
	return defaultConfig
}

// CompatibleAPI tries to be 100% compatible with standard library behavior
func CompatibleAPI() API {
	var compatibleWithStandardLibrary = Config{
		EscapeHTML:             true,
		SortMapKeys:            true,
		ValidateJsonRawMessage: true,
	}.Froze()
	return compatibleWithStandardLibrary
}

// FastestAPI marshals float with only 6 digits precision
func FastestAPI() API {
	var fastest = Config{
		EscapeHTML:                    false,
		MarshalFloatWith6Digits:       true, // will lose precession
		ObjectFieldMustBeSimpleString: true, // do not unescape object field
	}.Froze()
	return fastest
}

type frozenConfig struct {
	configBeforeFrozen            Config
	sortMapKeys                   bool
	indentionStep                 int
	objectFieldMustBeSimpleString bool
	onlyTaggedField               bool
	disallowUnknownFields         bool
	decoderCache                  *Map
	encoderCache                  *Map
	encoderExtension              Extension
	decoderExtension              Extension
	extraExtensions               []Extension
	streamPool                    *sync.Pool
	iteratorPool                  *sync.Pool
	caseSensitive                 bool
	typeDecoders                  map[string]ValDecoder
	fieldDecoders                 map[string]ValDecoder
	typeEncoders                  map[string]ValEncoder
	fieldEncoders                 map[string]ValEncoder
	locker                        *sync.Mutex
}

func (cfg *frozenConfig) initCache() {
	cfg.decoderCache = NewMap()
	cfg.encoderCache = NewMap()
}

func (cfg *frozenConfig) addDecoderToCache(cacheKey uintptr, decoder ValDecoder) {
	cfg.decoderCache.Store(cacheKey, decoder)
}

func (cfg *frozenConfig) addEncoderToCache(cacheKey uintptr, encoder ValEncoder) {
	cfg.encoderCache.Store(cacheKey, encoder)
}

func (cfg *frozenConfig) getDecoderFromCache(cacheKey uintptr) ValDecoder {
	decoder, found := cfg.decoderCache.Load(cacheKey)
	if found {
		return decoder.(ValDecoder)
	}
	return nil
}

func (cfg *frozenConfig) getEncoderFromCache(cacheKey uintptr) ValEncoder {
	encoder, found := cfg.encoderCache.Load(cacheKey)
	if found {
		return encoder.(ValEncoder)
	}
	return nil
}

// Froze forge API from config
func (cfg Config) Froze() API {
	api := frozenConfig{
		sortMapKeys:                   cfg.SortMapKeys,
		indentionStep:                 cfg.IndentionStep,
		objectFieldMustBeSimpleString: cfg.ObjectFieldMustBeSimpleString,
		onlyTaggedField:               cfg.OnlyTaggedField,
		disallowUnknownFields:         cfg.DisallowUnknownFields,
		caseSensitive:                 cfg.CaseSensitive,
		locker:                        &sync.Mutex{},
	}
	api.streamPool = &sync.Pool{
		New: func() interface{} {
			return NewStream(&api, nil, 512)
		},
	}
	api.iteratorPool = &sync.Pool{
		New: func() interface{} {
			return NewIterator(&api)
		},
	}
	api.initCache()
	encoderExtension := EncoderExtension{}
	decoderExtension := DecoderExtension{}

	api.typeDecoders = map[string]ValDecoder{}
	api.fieldDecoders = map[string]ValDecoder{}
	api.typeEncoders = map[string]ValEncoder{}
	api.fieldEncoders = map[string]ValEncoder{}
	api.extraExtensions = []Extension{}

	if cfg.MarshalFloatWith6Digits {
		api.marshalFloatWith6Digits(encoderExtension)
	}
	if cfg.EscapeHTML {
		api.escapeHTML(encoderExtension)
	}
	if cfg.UseNumber {
		api.useNumber(decoderExtension)
	}
	if cfg.ValidateJsonRawMessage {
		api.validateJsonRawMessage(encoderExtension)
	}
	api.encoderExtension = encoderExtension
	api.decoderExtension = decoderExtension
	api.configBeforeFrozen = cfg
	return &api
}

type RawMessage []byte

func (cfg *frozenConfig) validateJsonRawMessage(extension EncoderExtension) {
	encoder := &funcEncoder{func(ptr unsafe.Pointer, stream *Stream) {
		rawMessage := *(*json.RawMessage)(ptr)
		iter := cfg.BorrowIterator([]byte(rawMessage))
		defer cfg.ReturnIterator(iter)
		iter.Read()
		if iter.Error != nil && iter.Error != io.EOF {
			stream.WriteRaw("null")
		} else {
			stream.WriteRaw(string(rawMessage))
		}
	}, func(ptr unsafe.Pointer) bool {
		return len(*((*json.RawMessage)(ptr))) == 0
	}}
	extension[reflect2.TypeOfPtr((*json.RawMessage)(nil)).Elem()] = encoder
	extension[reflect2.TypeOfPtr((*RawMessage)(nil)).Elem()] = encoder
}

func (cfg *frozenConfig) useNumber(extension DecoderExtension) {
	extension[reflect2.TypeOfPtr((*interface{})(nil)).Elem()] = &funcDecoder{func(ptr unsafe.Pointer, iter *Iterator) {
		exitingValue := *((*interface{})(ptr))
		if exitingValue != nil && reflect.TypeOf(exitingValue).Kind() == reflect.Ptr {
			iter.ReadVal(exitingValue)
			return
		}
		if iter.WhatIsNext() == NumberValue {
			*((*interface{})(ptr)) = json.Number(iter.readNumberAsString())
		} else {
			*((*interface{})(ptr)) = iter.Read()
		}
	}}
}
func (cfg *frozenConfig) getTagKey() string {
	tagKey := cfg.configBeforeFrozen.TagKey
	if tagKey == "" {
		return "json"
	}
	return tagKey
}

func (cfg *frozenConfig) RegisterExtension(extension Extension) {
	cfg.extraExtensions = append(cfg.extraExtensions, extension)
	copied := cfg.configBeforeFrozen
	cfg.configBeforeFrozen = copied
}

func (cfg *frozenConfig) RegisterTypeEncoderFunc(typ string, fun EncoderFunc, isEmptyFunc func(unsafe.Pointer) bool) {
	cfg.typeEncoders[typ] = &funcEncoder{fun, isEmptyFunc}
}

func (cfg *frozenConfig) RegisterTypeDecoderFunc(typ string, fun DecoderFunc) {
	cfg.typeDecoders[typ] = &funcDecoder{fun}
}

// RegisterTypeDecoder register TypeDecoder for a typ
func (cfg *frozenConfig) RegisterTypeDecoder(typ string, decoder ValDecoder) {
	cfg.typeDecoders[typ] = decoder
}

// RegisterFieldDecoderFunc register TypeDecoder for a struct field with function
func (cfg *frozenConfig) RegisterFieldDecoderFunc(typ string, field string, fun DecoderFunc) {
	cfg.RegisterFieldDecoder(typ, field, &funcDecoder{fun})
}

// RegisterFieldDecoder register TypeDecoder for a struct field
func (cfg *frozenConfig) RegisterFieldDecoder(typ string, field string, decoder ValDecoder) {
	cfg.fieldDecoders[fmt.Sprintf("%s/%s", typ, field)] = decoder
}
func (cfg *frozenConfig) RegisterFieldEncoderFunc(typ string, field string, fun EncoderFunc, isEmptyFunc func(unsafe.Pointer) bool) {
	cfg.RegisterFieldEncoder(typ, field, &funcEncoder{fun, isEmptyFunc})
}

func (cfg *frozenConfig) RegisterFieldEncoder(typ string, field string, encoder ValEncoder) {
	cfg.fieldEncoders[fmt.Sprintf("%s/%s", typ, field)] = encoder
}

func (cfg *frozenConfig) RegisterTypeEncoder(typ string, encoder ValEncoder) {
	cfg.typeEncoders[typ] = encoder
}

type lossyFloat32Encoder struct {
}

func (encoder *lossyFloat32Encoder) Encode(ptr unsafe.Pointer, stream *Stream) {
	stream.WriteFloat32Lossy(*((*float32)(ptr)))
}

func (encoder *lossyFloat32Encoder) IsEmpty(ptr unsafe.Pointer) bool {
	return *((*float32)(ptr)) == 0
}

type lossyFloat64Encoder struct {
}

func (encoder *lossyFloat64Encoder) Encode(ptr unsafe.Pointer, stream *Stream) {
	stream.WriteFloat64Lossy(*((*float64)(ptr)))
}

func (encoder *lossyFloat64Encoder) IsEmpty(ptr unsafe.Pointer) bool {
	return *((*float64)(ptr)) == 0
}

// EnableLossyFloatMarshalling keeps 10**(-6) precision
// for float variables for better performance.
func (cfg *frozenConfig) marshalFloatWith6Digits(extension EncoderExtension) {
	// for better performance
	extension[reflect2.TypeOfPtr((*float32)(nil)).Elem()] = &lossyFloat32Encoder{}
	extension[reflect2.TypeOfPtr((*float64)(nil)).Elem()] = &lossyFloat64Encoder{}
}

type htmlEscapedStringEncoder struct {
}

func (encoder *htmlEscapedStringEncoder) Encode(ptr unsafe.Pointer, stream *Stream) {
	str := *((*string)(ptr))
	stream.WriteStringWithHTMLEscaped(str)
}

func (encoder *htmlEscapedStringEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return *((*string)(ptr)) == ""
}

func (cfg *frozenConfig) escapeHTML(encoderExtension EncoderExtension) {
	encoderExtension[reflect2.TypeOfPtr((*string)(nil)).Elem()] = &htmlEscapedStringEncoder{}
}

func (cfg *frozenConfig) ClearDecoders() {
	cfg.typeDecoders = map[string]ValDecoder{}
	cfg.fieldDecoders = map[string]ValDecoder{}
	*cfg = *(cfg.configBeforeFrozen.Froze().(*frozenConfig))
}

func (cfg *frozenConfig) ClearExtensions() {
	cfg.extraExtensions = []Extension{}
	*cfg = *(cfg.configBeforeFrozen.Froze().(*frozenConfig))
}

func (cfg *frozenConfig) ClearEncoders() {
	cfg.typeEncoders = map[string]ValEncoder{}
	cfg.fieldEncoders = map[string]ValEncoder{}
	*cfg = *(cfg.configBeforeFrozen.Froze().(*frozenConfig))
}

func (cfg *frozenConfig) MarshalToString(v interface{}) (string, error) {
	stream := cfg.BorrowStream(nil)
	defer cfg.ReturnStream(stream)
	stream.WriteVal(v)
	if stream.Error != nil {
		return "", stream.Error
	}
	return string(stream.Buffer()), nil
}

func (cfg *frozenConfig) Marshal(v interface{}) ([]byte, error) {
	return cfg._marshal(v, 0)
}

func (cfg *frozenConfig) _marshal(v interface{}, indent int) ([]byte, error) {
	cfg.indentionStep = indent
	stream := cfg.BorrowStream(nil)
	defer cfg.ReturnStream(stream)
	stream.WriteVal(v)
	if stream.Error != nil {
		return nil, stream.Error
	}
	result := stream.Buffer()
	copied := make([]byte, len(result))
	copy(copied, result)
	return copied, nil
}

func (cfg *frozenConfig) MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	if prefix != "" {
		panic("prefix is not supported")
	}
	for _, r := range indent {
		if r != ' ' {
			panic("indent can only be space")
		}
	}
	//newCfg := cfg.configBeforeFrozen
	//newCfg.IndentionStep = len(indent)
	return cfg._marshal(v, len(indent))
}

func (cfg *frozenConfig) UnmarshalFromString(str string, v interface{}) error {
	data := []byte(str)
	iter := cfg.BorrowIterator(data)
	defer cfg.ReturnIterator(iter)
	iter.ReadVal(v)
	c := iter.nextToken()
	if c == 0 {
		if iter.Error == io.EOF {
			return nil
		}
		return iter.Error
	}
	iter.ReportError("Unmarshal", "there are bytes left after unmarshal")
	return iter.Error
}

func (cfg *frozenConfig) Get(data []byte, path ...interface{}) Any {
	iter := cfg.BorrowIterator(data)
	defer cfg.ReturnIterator(iter)
	return locatePath(iter, path)
}

func (cfg *frozenConfig) Unmarshal(data []byte, v interface{}) error {
	iter := cfg.BorrowIterator(data)
	defer cfg.ReturnIterator(iter)
	iter.ReadVal(v)
	c := iter.nextToken()
	if c == 0 {
		if iter.Error == io.EOF {
			return nil
		}
		return iter.Error
	}
	iter.ReportError("Unmarshal", "there are bytes left after unmarshal")
	return iter.Error
}