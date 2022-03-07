package test

import (
	"github.com/heskandari/jsoner"
	"testing"
)

type StringerType int
type StructWithStringer struct {
	MyVal StringerType
}

func (t StringerType) String() string {
	switch t {
	case 1:
		return "True"
	case 0:
		return "False"
	}
	return ""
}

func TestStringer(t *testing.T) {
	cfg := jsoner.Config{}.Froze()

	st := StructWithStringer{
		MyVal: StringerType(1),
	}
	b, err := cfg.Marshal(st)
	js := string(b)

	if err != nil {
		t.Fatalf("failed to marshal with jsoner: %v", err)
	}
	if js != "{\"MyVal\":\"True\"}" {
		t.Fatalf("failed to marshal Stringer with jsoner: %v", err)
	}
}