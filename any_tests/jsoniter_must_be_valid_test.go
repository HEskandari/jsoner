package any_tests

import (
	"testing"

	"github.com/heskandari/jsoner"
	"github.com/stretchr/testify/require"
)

// if must be valid is useless, just drop this test
func Test_must_be_valid(t *testing.T) {
	should := require.New(t)
	api := jsoner.DefaultAPI()
	any := api.Get([]byte("123"))
	should.Equal(any.MustBeValid().ToInt(), 123)

	any = jsoner.Wrap(int8(10))
	should.Equal(any.MustBeValid().ToInt(), 10)

	any = jsoner.Wrap(int16(10))
	should.Equal(any.MustBeValid().ToInt(), 10)

	any = jsoner.Wrap(int32(10))
	should.Equal(any.MustBeValid().ToInt(), 10)

	any = jsoner.Wrap(int64(10))
	should.Equal(any.MustBeValid().ToInt(), 10)

	any = jsoner.Wrap(uint(10))
	should.Equal(any.MustBeValid().ToInt(), 10)

	any = jsoner.Wrap(uint8(10))
	should.Equal(any.MustBeValid().ToInt(), 10)

	any = jsoner.Wrap(uint16(10))
	should.Equal(any.MustBeValid().ToInt(), 10)

	any = jsoner.Wrap(uint32(10))
	should.Equal(any.MustBeValid().ToInt(), 10)

	any = jsoner.Wrap(uint64(10))
	should.Equal(any.MustBeValid().ToInt(), 10)

	any = jsoner.Wrap(float32(10))
	should.Equal(any.MustBeValid().ToFloat64(), float64(10))

	any = jsoner.Wrap(float64(10))
	should.Equal(any.MustBeValid().ToFloat64(), float64(10))

	any = jsoner.Wrap(true)
	should.Equal(any.MustBeValid().ToFloat64(), float64(1))

	any = jsoner.Wrap(false)
	should.Equal(any.MustBeValid().ToFloat64(), float64(0))

	any = jsoner.Wrap(nil)
	should.Equal(any.MustBeValid().ToFloat64(), float64(0))

	any = jsoner.Wrap(struct{ age int }{age: 1})
	should.Equal(any.MustBeValid().ToFloat64(), float64(0))

	any = jsoner.Wrap(map[string]interface{}{"abc": 1})
	should.Equal(any.MustBeValid().ToFloat64(), float64(0))

	any = jsoner.Wrap("abc")
	should.Equal(any.MustBeValid().ToFloat64(), float64(0))

	any = jsoner.Wrap([]int{})
	should.Equal(any.MustBeValid().ToFloat64(), float64(0))

	any = jsoner.Wrap([]int{1, 2})
	should.Equal(any.MustBeValid().ToFloat64(), float64(1))
}