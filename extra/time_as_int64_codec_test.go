package extra

import (
	"github.com/heskandari/jsoner"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_time_as_int64(t *testing.T) {
	should := require.New(t)
	api := jsoner.DefaultAPI()
	RegisterTimeAsInt64Codec(api, time.Nanosecond)
	output, err := api.Marshal(time.Unix(1497952257, 1002))
	should.Nil(err)
	should.Equal("1497952257000001002", string(output))
	var val time.Time
	should.Nil(api.Unmarshal(output, &val))
	should.Equal(int64(1497952257000001002), val.UnixNano())
}

func Test_time_as_int64_keep_microsecond(t *testing.T) {
	should := require.New(t)
	api := jsoner.DefaultAPI()
	RegisterTimeAsInt64Codec(api, time.Microsecond)
	output, err := api.Marshal(time.Unix(1, 1002))
	should.Nil(err)
	should.Equal("1000001", string(output))
	var val time.Time
	should.Nil(api.Unmarshal(output, &val))
	should.Equal(int64(1000001000), val.UnixNano())
}