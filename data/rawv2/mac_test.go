package rawv2_test

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"github.com/susji/ruuvi/data/rawv2"
)

func TestMacEncodeJson(t *testing.T) {
	m := rawv2.MAC{Value: [6]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0xF6}}
	raw, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	enc := string(raw)
	if enc != `"01:02:03:04:05:f6"` {
		t.Error(enc)
	}
}

func TestMacDecodeJsonLower(t *testing.T) {
	m := `"06:05:04:a3:02:01"`
	var dec rawv2.MAC
	err := json.Unmarshal([]byte(m), &dec)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(dec.Value, [6]byte{0x06, 0x05, 0x04, 0xA3, 0x02, 0x01}) {
		t.Error(dec)
	}
}

func TestMacDecodeJsonUpper(t *testing.T) {
	m := `"06:05:04:F3:02:01"`
	var dec rawv2.MAC
	err := json.Unmarshal([]byte(m), &dec)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(dec.Value, [6]byte{0x06, 0x05, 0x04, 0xF3, 0x02, 0x01}) {
		t.Error(dec)
	}
}

func TestMacDecodeJsonBad(t *testing.T) {
	m := `"06:05:04:F3:02:XY"`
	var dec rawv2.MAC
	err := json.Unmarshal([]byte(m), &dec)
	if !errors.Is(err, rawv2.ErrorBadMAC) {
		t.Fatal(err)
	}
}
