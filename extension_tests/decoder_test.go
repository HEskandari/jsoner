package test

import (
	"bytes"
	"fmt"
	"github.com/heskandari/jsoner"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
	"time"
	"unsafe"
)

func Test_customize_type_decoder(t *testing.T) {
	api := jsoner.DefaultAPI()
	api.RegisterTypeDecoderFunc("time.Time", func(ptr unsafe.Pointer, iter *jsoner.Iterator) {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", iter.ReadString(), time.UTC)
		if err != nil {
			iter.Error = err
			return
		}
		*((*time.Time)(ptr)) = t
	})
	val := time.Time{}
	err := api.Unmarshal([]byte(`"2016-12-05 08:43:28"`), &val)
	if err != nil {
		t.Fatal(err)
	}
	year, month, day := val.Date()
	if year != 2016 || month != 12 || day != 5 {
		t.Fatal(val)
	}
}

func Test_customize_byte_array_encoder(t *testing.T) {
	should := require.New(t)
	api := jsoner.DefaultAPI()
	api.RegisterTypeEncoderFunc("[]uint8", func(ptr unsafe.Pointer, stream *jsoner.Stream) {
		t := *((*[]byte)(ptr))
		stream.WriteString(string(t))
	}, nil)

	val := []byte("abc")
	str, err := api.MarshalToString(val)
	should.Nil(err)
	should.Equal(`"abc"`, str)
}

type CustomEncoderAttachmentTestStruct struct {
	Value int32 `json:"value"`
}

type CustomEncoderAttachmentTestStructEncoder struct{}

func (c *CustomEncoderAttachmentTestStructEncoder) Encode(ptr unsafe.Pointer, stream *jsoner.Stream) {
	attachVal, ok := stream.Attachment.(int)
	stream.WriteRaw(`"`)
	stream.WriteRaw(fmt.Sprintf("%t %d", ok, attachVal))
	stream.WriteRaw(`"`)
}

func (c *CustomEncoderAttachmentTestStructEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return false
}

func Test_custom_encoder_attachment(t *testing.T) {

	config := jsoner.Config{SortMapKeys: true}.Froze()
	config.RegisterTypeEncoder("test.CustomEncoderAttachmentTestStruct", &CustomEncoderAttachmentTestStructEncoder{})
	expectedValue := 17
	should := require.New(t)
	buf := &bytes.Buffer{}
	stream := jsoner.NewStream(config, buf, 4096)
	stream.Attachment = expectedValue
	val := map[string]CustomEncoderAttachmentTestStruct{"a": {}}
	stream.WriteVal(val)
	stream.Flush()
	should.Nil(stream.Error)
	should.Equal("{\"a\":\"true 17\"}", buf.String())
}

func Test_customize_field_decoder(t *testing.T) {
	t.Skip("decoder not being picked up")
	type Tom struct {
		field1 string
	}
	api := jsoner.DefaultAPI()
	api.RegisterFieldDecoderFunc("jsoner.Tom", "field1", func(ptr unsafe.Pointer, iter *jsoner.Iterator) {
		*((*string)(ptr)) = strconv.Itoa(iter.ReadInt())
	})

	tom := Tom{}
	err := api.Unmarshal([]byte(`{"field1": 100}`), &tom)
	if err != nil {
		t.Fatal(err)
	}

	should := require.New(t)
	should.Equal(100, tom.field1)
}

func Test_recursive_empty_interface_customization(t *testing.T) {
	var obj interface{}
	api := jsoner.DefaultAPI()
	api.RegisterTypeDecoderFunc("interface {}", func(ptr unsafe.Pointer, iter *jsoner.Iterator) {
		switch iter.WhatIsNext() {
		case jsoner.NumberValue:
			*(*interface{})(ptr) = iter.ReadInt64()
		default:
			*(*interface{})(ptr) = iter.Read()
		}
	})
	should := require.New(t)
	api.Unmarshal([]byte("[100]"), &obj)
	should.Equal([]interface{}{int64(100)}, obj)
}

type MyInterface interface {
	Hello() string
}

type MyString string

func (ms MyString) Hello() string {
	return string(ms)
}

func Test_read_custom_interface(t *testing.T) {
	t.Skip()
	should := require.New(t)
	var val MyInterface
	api := jsoner.DefaultAPI()
	api.RegisterTypeDecoderFunc("jsoner.MyInterface", func(ptr unsafe.Pointer, iter *jsoner.Iterator) {
		*((*MyInterface)(ptr)) = MyString(iter.ReadString())
	})
	err := api.UnmarshalFromString(`"hello"`, &val)
	should.Nil(err)
	should.Equal("hello", val.Hello())
}

const flow1 = `
{"A":"hello"}
{"A":"hello"}
{"A":"hello"}
{"A":"hello"}
{"A":"hello"}`

const flow2 = `
{"A":"hello"}
{"A":"hello"}
{"A":"hello"}
{"A":"hello"}
{"A":"hello"}
`

type (
	Type1 struct {
		A string
	}

	Type2 struct {
		A string
	}
)

func (t *Type2) UnmarshalJSON(data []byte) error {
	return nil
}

func (t *Type2) MarshalJSON() ([]byte, error) {
	return nil, nil
}