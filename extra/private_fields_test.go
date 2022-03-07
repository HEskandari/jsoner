package extra

import (
	"github.com/heskandari/jsoner"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_private_fields(t *testing.T) {
	type TestObject struct {
		field1 string
	}
	js := jsoner.DefaultAPI()
	SupportPrivateFields(js)
	should := require.New(t)
	obj := TestObject{}
	should.Nil(js.UnmarshalFromString(`{"field1":"Hello"}`, &obj))
	should.Equal("Hello", obj.field1)
}