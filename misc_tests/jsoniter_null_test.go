package misc_tests

import (
	"bytes"
	"io"
	"testing"

	"github.com/heskandari/jsoner"
	"github.com/stretchr/testify/require"
)

func Test_read_null(t *testing.T) {
	should := require.New(t)
	iter := jsoner.ParseString(jsoner.DefaultAPI(), `null`)
	should.True(iter.ReadNil())
	iter = jsoner.ParseString(jsoner.DefaultAPI(), `null`)
	should.Nil(iter.Read())
	iter = jsoner.ParseString(jsoner.DefaultAPI(), `navy`)
	iter.Read()
	should.True(iter.Error != nil && iter.Error != io.EOF)
	iter = jsoner.ParseString(jsoner.DefaultAPI(), `navy`)
	iter.ReadNil()
	should.True(iter.Error != nil && iter.Error != io.EOF)
}

func Test_write_null(t *testing.T) {
	should := require.New(t)
	buf := &bytes.Buffer{}
	stream := jsoner.NewStream(jsoner.DefaultAPI(), buf, 4096)
	stream.WriteNil()
	stream.Flush()
	should.Nil(stream.Error)
	should.Equal("null", buf.String())
}

func Test_decode_null_object_field(t *testing.T) {
	should := require.New(t)
	iter := jsoner.ParseString(jsoner.DefaultAPI(), `[null,"a"]`)
	iter.ReadArray()
	if iter.ReadObject() != "" {
		t.FailNow()
	}
	iter.ReadArray()
	if iter.ReadString() != "a" {
		t.FailNow()
	}
	type TestObject struct {
		Field string
	}
	objs := []TestObject{}
	should.Nil(jsoner.DefaultAPI().UnmarshalFromString("[null]", &objs))
	should.Len(objs, 1)
}

func Test_decode_null_array_element(t *testing.T) {
	should := require.New(t)
	iter := jsoner.ParseString(jsoner.DefaultAPI(), `[null,"a"]`)
	should.True(iter.ReadArray())
	should.True(iter.ReadNil())
	should.True(iter.ReadArray())
	should.Equal("a", iter.ReadString())
}

func Test_decode_null_string(t *testing.T) {
	should := require.New(t)
	iter := jsoner.ParseString(jsoner.DefaultAPI(), `[null,"a"]`)
	should.True(iter.ReadArray())
	should.Equal("", iter.ReadString())
	should.True(iter.ReadArray())
	should.Equal("a", iter.ReadString())
}

func Test_decode_null_skip(t *testing.T) {
	iter := jsoner.ParseString(jsoner.DefaultAPI(), `[null,"a"]`)
	iter.ReadArray()
	iter.Skip()
	iter.ReadArray()
	if iter.ReadString() != "a" {
		t.FailNow()
	}
}