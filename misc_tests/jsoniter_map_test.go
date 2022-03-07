package misc_tests

import (
	"math/big"
	"testing"

	"github.com/heskandari/jsoner"
	"github.com/stretchr/testify/require"
)

func Test_decode_TextMarshaler_key_map(t *testing.T) {
	should := require.New(t)
	var val map[*big.Float]string
	should.Nil(jsoner.DefaultAPI().UnmarshalFromString(`{"1":"2"}`, &val))
	str, err := jsoner.DefaultAPI().MarshalToString(val)
	should.Nil(err)
	should.Equal(`{"1":"2"}`, str)
}

func Test_map_eface_of_eface(t *testing.T) {
	should := require.New(t)
	json := jsoner.CompatibleAPI()
	output, err := json.MarshalToString(map[interface{}]interface{}{
		"1": 2,
		3:   "4",
	})
	should.NoError(err)
	should.Equal(`{"1":2,"3":"4"}`, output)
}

func Test_encode_nil_map(t *testing.T) {
	should := require.New(t)
	var nilMap map[string]string
	output, err := jsoner.DefaultAPI().MarshalToString(nilMap)
	should.NoError(err)
	should.Equal(`null`, output)
}