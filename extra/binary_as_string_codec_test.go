package extra

import (
	"github.com/heskandari/jsoner"
	"github.com/stretchr/testify/require"
	"testing"
)

var js jsoner.API

func init() {
	js = jsoner.DefaultAPI()
	js.RegisterExtension(&BinaryAsStringExtension{})
}

func TestBinaryAsStringCodec(t *testing.T) {
	t.Run("safe set", func(t *testing.T) {
		should := require.New(t)
		output, err := js.Marshal([]byte("hello"))
		should.NoError(err)
		should.Equal(`"hello"`, string(output))
		var val []byte
		should.NoError(js.Unmarshal(output, &val))
		should.Equal(`hello`, string(val))
	})
	t.Run("non safe set", func(t *testing.T) {
		should := require.New(t)
		output, err := js.Marshal([]byte{1, 2, 3, 23})
		should.NoError(err)
		should.Equal(`"\\x01\\x02\\x03\\x17"`, string(output))
		var val []byte
		should.NoError(js.Unmarshal(output, &val))
		should.Equal([]byte{1, 2, 3, 23}, val)
	})
}