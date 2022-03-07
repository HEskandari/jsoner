package skip_tests

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/heskandari/jsoner"
	"github.com/stretchr/testify/require"
)

func Test_skip_number_in_array(t *testing.T) {
	should := require.New(t)
	iter := jsoner.ParseString(jsoner.DefaultAPI(), `[-0.12, "stream"]`)
	iter.ReadArray()
	iter.Skip()
	iter.ReadArray()
	should.Nil(iter.Error)
	should.Equal("stream", iter.ReadString())
}

func Test_skip_string_in_array(t *testing.T) {
	should := require.New(t)
	iter := jsoner.ParseString(jsoner.DefaultAPI(), `["hello", "stream"]`)
	iter.ReadArray()
	iter.Skip()
	iter.ReadArray()
	should.Nil(iter.Error)
	should.Equal("stream", iter.ReadString())
}

func Test_skip_null(t *testing.T) {
	iter := jsoner.ParseString(jsoner.DefaultAPI(), `[null , "stream"]`)
	iter.ReadArray()
	iter.Skip()
	iter.ReadArray()
	if iter.ReadString() != "stream" {
		t.FailNow()
	}
}

func Test_skip_true(t *testing.T) {
	iter := jsoner.ParseString(jsoner.DefaultAPI(), `[true , "stream"]`)
	iter.ReadArray()
	iter.Skip()
	iter.ReadArray()
	if iter.ReadString() != "stream" {
		t.FailNow()
	}
}

func Test_skip_false(t *testing.T) {
	iter := jsoner.ParseString(jsoner.DefaultAPI(), `[false , "stream"]`)
	iter.ReadArray()
	iter.Skip()
	iter.ReadArray()
	if iter.ReadString() != "stream" {
		t.FailNow()
	}
}

func Test_skip_array(t *testing.T) {
	iter := jsoner.ParseString(jsoner.DefaultAPI(), `[[1, [2, [3], 4]], "stream"]`)
	iter.ReadArray()
	iter.Skip()
	iter.ReadArray()
	if iter.ReadString() != "stream" {
		t.FailNow()
	}
}

func Test_skip_empty_array(t *testing.T) {
	iter := jsoner.ParseString(jsoner.DefaultAPI(), `[ [ ], "stream"]`)
	iter.ReadArray()
	iter.Skip()
	iter.ReadArray()
	if iter.ReadString() != "stream" {
		t.FailNow()
	}
}

func Test_skip_nested(t *testing.T) {
	iter := jsoner.ParseString(jsoner.DefaultAPI(), `[ {"a" : [{"stream": "c"}], "d": 102 }, "stream"]`)
	iter.ReadArray()
	iter.Skip()
	iter.ReadArray()
	if iter.ReadString() != "stream" {
		t.FailNow()
	}
}

func Test_skip_and_return_bytes(t *testing.T) {
	should := require.New(t)
	iter := jsoner.ParseString(jsoner.DefaultAPI(), `[ {"a" : [{"stream": "c"}], "d": 102 }, "stream"]`)
	iter.ReadArray()
	skipped := iter.SkipAndReturnBytes()
	should.Equal(`{"a" : [{"stream": "c"}], "d": 102 }`, string(skipped))
}

func Test_skip_and_return_bytes_with_reader(t *testing.T) {
	should := require.New(t)
	iter := jsoner.Parse(jsoner.DefaultAPI(), bytes.NewBufferString(`[ {"a" : [{"stream": "c"}], "d": 102 }, "stream"]`), 4)
	iter.ReadArray()
	skipped := iter.SkipAndReturnBytes()
	should.Equal(`{"a" : [{"stream": "c"}], "d": 102 }`, string(skipped))
}

func Test_append_skip_and_return_bytes_with_reader(t *testing.T) {
	should := require.New(t)
	iter := jsoner.Parse(jsoner.DefaultAPI(), bytes.NewBufferString(`[ {"a" : [{"stream": "c"}], "d": 102 }, "stream"]`), 4)
	iter.ReadArray()
	buf := make([]byte, 0, 1024)
	buf = iter.SkipAndAppendBytes(buf)
	should.Equal(`{"a" : [{"stream": "c"}], "d": 102 }`, string(buf))
}

func Test_skip_empty(t *testing.T) {
	should := require.New(t)
	should.NotNil(jsoner.DefaultAPI().Get([]byte("")).LastError())
}

type TestResp struct {
	Code uint64
}

func Benchmark_jsoner_skip(b *testing.B) {
	input := []byte(`
{
    "_shards":{
        "total" : 5,
        "successful" : 5,
        "failed" : 0
    },
    "hits":{
        "total" : 1,
        "hits" : [
            {
                "_index" : "twitter",
                "_type" : "tweet",
                "_id" : "1",
                "_source" : {
                    "user" : "kimchy",
                    "postDate" : "2009-11-15T14:12:12",
                    "message" : "trying out Elasticsearch"
                }
            }
        ]
    },
    "code": 200
}`)
	for n := 0; n < b.N; n++ {
		result := TestResp{}
		iter := jsoner.ParseBytes(jsoner.DefaultAPI(), input)
		for field := iter.ReadObject(); field != ""; field = iter.ReadObject() {
			switch field {
			case "code":
				result.Code = iter.ReadUint64()
			default:
				iter.Skip()
			}
		}
	}
}

func Benchmark_json_skip(b *testing.B) {
	input := []byte(`
{
    "_shards":{
        "total" : 5,
        "successful" : 5,
        "failed" : 0
    },
    "hits":{
        "total" : 1,
        "hits" : [
            {
                "_index" : "twitter",
                "_type" : "tweet",
                "_id" : "1",
                "_source" : {
                    "user" : "kimchy",
                    "postDate" : "2009-11-15T14:12:12",
                    "message" : "trying out Elasticsearch"
                }
            }
        ]
    },
    "code": 200
}`)
	for n := 0; n < b.N; n++ {
		result := TestResp{}
		json.Unmarshal(input, &result)
	}
}