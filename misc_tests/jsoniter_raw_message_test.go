package misc_tests

import (
	"encoding/json"
	"github.com/heskandari/jsoner"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_jsoner_RawMessage(t *testing.T) {
	should := require.New(t)
	var data jsoner.RawMessage
	should.Nil(jsoner.DefaultAPI().Unmarshal([]byte(`[1,2,3]`), &data))
	should.Equal(`[1,2,3]`, string(data))
	str, err := jsoner.DefaultAPI().MarshalToString(data)
	should.Nil(err)
	should.Equal(`[1,2,3]`, str)
}

func Test_encode_map_of_jsoner_raw_message(t *testing.T) {
	should := require.New(t)
	type RawMap map[string]*jsoner.RawMessage
	value := jsoner.RawMessage("[]")
	rawMap := RawMap{"hello": &value}
	output, err := jsoner.DefaultAPI().MarshalToString(rawMap)
	should.Nil(err)
	should.Equal(`{"hello":[]}`, output)
}

func Test_marshal_invalid_json_raw_message(t *testing.T) {
	type A struct {
		Raw json.RawMessage `json:"raw"`
	}
	message := []byte(`{}`)

	a := A{}
	should := require.New(t)
	should.Nil(jsoner.CompatibleAPI().Unmarshal(message, &a))
	aout, aouterr := jsoner.CompatibleAPI().Marshal(&a)
	should.Equal(`{"raw":null}`, string(aout))
	should.Nil(aouterr)
}

func Test_marshal_nil_json_raw_message(t *testing.T) {
	type A struct {
		Nil1 jsoner.RawMessage `json:"raw1"`
		Nil2 json.RawMessage   `json:"raw2"`
	}

	a := A{}
	should := require.New(t)
	aout, aouterr := jsoner.DefaultAPI().Marshal(&a)
	should.Equal(`{"raw1":null,"raw2":null}`, string(aout))
	should.Nil(aouterr)

	a.Nil1 = []byte(`Any`)
	a.Nil2 = []byte(`Any`)
	should.Nil(jsoner.DefaultAPI().Unmarshal(aout, &a))
	should.Nil(a.Nil1)
	should.Nil(a.Nil2)
}