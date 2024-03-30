package rawv2_test

// Test vectors below are from Ruuvi's documentation at
// https://docs.ruuvi.com/communication/bluetooth-advertisements/data-format-5-rawv2

import (
	"encoding/hex"
	"errors"
	"testing"

	"github.com/susji/ruuvi/data"
	"github.com/susji/ruuvi/data/rawv2"
)

const (
	VECTOR_GOOD = "0512FC5394C37C0004FFFC040CAC364200CDCBB8334C884F"
	VECTOR_MAX  = "057FFFFFFEFFFE7FFF7FFF7FFFFFDEFEFFFECBB8334C884F"
	VECTOR_MIN  = "058001000000008001800180010000000000CBB8334C884F"
	VECTOR_BAD  = "058000FFFFFFFF800080008000FFFFFFFFFFFFFFFFFFFFFF"
)

func TestBasic(t *testing.T) {
	raw := make([]byte, len(VECTOR_GOOD)/2)
	n, _ := hex.Decode(raw, []byte(VECTOR_GOOD))
	if n != len(raw) {
		t.Fatal("unexpected n:", n)
	}
	p, err := rawv2.Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	if p.Type != data.VERSION_RAWV2 {
		t.Error(p.Type)
	}
	if !p.Temperature.Valid || !ae(p.Temperature.Value, 24.3) {
		t.Error(p.Temperature)
	}
	if !p.Pressure.Valid || p.Pressure.Value != 100044 {
		t.Error(p.Pressure)
	}
	if !p.Humidity.Valid || !ae(p.Humidity.Value, 53.49) {
		t.Error(p.Humidity)
	}
	if !p.AccelerationX.Valid || p.AccelerationX.Value != 4 {
		t.Error(p.AccelerationX)
	}
	if !p.AccelerationY.Valid || p.AccelerationY.Value != -4 {
		t.Error(p.AccelerationY)
	}
	if !p.AccelerationZ.Valid || p.AccelerationZ.Value != 1036 {
		t.Error(p.AccelerationZ)
	}
	if !p.TransmitPower.Valid || p.TransmitPower.Value != 4 {
		t.Error(p.TransmitPower)
	}
	if !p.BatteryVoltage.Valid || !ae(p.BatteryVoltage.Value, 2.977) {
		t.Error(p.BatteryVoltage)
	}
	if p.MAC.Value[0] != 0xCB || p.MAC.Value[1] != 0xB8 || p.MAC.Value[2] != 0x33 ||
		p.MAC.Value[3] != 0x4C || p.MAC.Value[4] != 0x88 || p.MAC.Value[5] != 0x4F {
		t.Error(p.MAC)
	}
}

func TestMaximum(t *testing.T) {
	raw := make([]byte, len(VECTOR_MAX)/2)
	n, _ := hex.Decode(raw, []byte(VECTOR_MAX))
	if n != len(raw) {
		t.Fatal("unexpected n:", n)
	}
	p, err := rawv2.Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	if p.Type != data.VERSION_RAWV2 {
		t.Error(p.Type)
	}
	if !p.Temperature.Valid || !ae(p.Temperature.Value, 163.835) {
		t.Error(p.Temperature)
	}
	if !p.Pressure.Valid || p.Pressure.Value != 115534 {
		t.Error(p.Pressure)
	}
	if !p.Humidity.Valid || !ae(p.Humidity.Value, 163.8350) {
		t.Error(p.Humidity)
	}
	if !p.AccelerationX.Valid || p.AccelerationX.Value != 32767 {
		t.Error(p.AccelerationX)
	}
	if !p.AccelerationY.Valid || p.AccelerationY.Value != 32767 {
		t.Error(p.AccelerationY)
	}
	if !p.AccelerationZ.Valid || p.AccelerationZ.Value != 32767 {
		t.Error(p.AccelerationZ)
	}
	if !p.TransmitPower.Valid || p.TransmitPower.Value != 20 {
		t.Error(p.TransmitPower)
	}
	if !p.BatteryVoltage.Valid || !ae(p.BatteryVoltage.Value, 3.646) {
		t.Error(p.BatteryVoltage)
	}
	if p.MAC.Value[0] != 0xCB || p.MAC.Value[1] != 0xB8 || p.MAC.Value[2] != 0x33 ||
		p.MAC.Value[3] != 0x4C || p.MAC.Value[4] != 0x88 || p.MAC.Value[5] != 0x4F {
		t.Error(p.MAC)
	}
}

func TestMinimum(t *testing.T) {
	raw := make([]byte, len(VECTOR_MIN)/2)
	n, _ := hex.Decode(raw, []byte(VECTOR_MIN))
	if n != len(raw) {
		t.Fatal("unexpected n:", n)
	}
	p, err := rawv2.Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	if p.Type != data.VERSION_RAWV2 {
		t.Error(p.Type)
	}
	if !p.Temperature.Valid || !ae(p.Temperature.Value, -163.835) {
		t.Error(p.Temperature)
	}
	if !p.Pressure.Valid || p.Pressure.Value != 50000 {
		t.Error(p.Pressure)
	}
	if !p.Humidity.Valid || !ae(p.Humidity.Value, 0) {
		t.Error(p.Humidity)
	}
	if !p.AccelerationX.Valid || p.AccelerationX.Value != -32767 {
		t.Error(p.AccelerationX)
	}
	if !p.AccelerationY.Valid || p.AccelerationY.Value != -32767 {
		t.Error(p.AccelerationY)
	}
	if !p.AccelerationZ.Valid || p.AccelerationZ.Value != -32767 {
		t.Error(p.AccelerationZ)
	}
	if !p.TransmitPower.Valid || p.TransmitPower.Value != -40 {
		t.Error(p.TransmitPower)
	}
	if !p.BatteryVoltage.Valid || !ae(p.BatteryVoltage.Value, 1.6) {
		t.Error(p.BatteryVoltage)
	}
	if p.MAC.Value[0] != 0xCB || p.MAC.Value[1] != 0xB8 || p.MAC.Value[2] != 0x33 ||
		p.MAC.Value[3] != 0x4C || p.MAC.Value[4] != 0x88 || p.MAC.Value[5] != 0x4F {
		t.Error(p.MAC)
	}
}

func TestInvalid(t *testing.T) {
	raw := make([]byte, len(VECTOR_BAD)/2)
	n, _ := hex.Decode(raw, []byte(VECTOR_BAD))
	if n != len(raw) {
		t.Fatal("unexpected n:", n)
	}
	p, err := rawv2.Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	if p.Type != data.VERSION_RAWV2 {
		t.Error(p.Type)
	}
	if p.Temperature.Valid {
		t.Error(p.Temperature)
	}
	if p.Humidity.Valid {
		t.Error(p.Humidity)
	}
	if p.Pressure.Valid {
		t.Error(p.Pressure)
	}
	if p.AccelerationX.Valid {
		t.Error(p.AccelerationX)
	}
	if p.AccelerationY.Valid {
		t.Error(p.AccelerationY)
	}
	if p.AccelerationZ.Valid {
		t.Error(p.AccelerationZ)
	}
	if p.TransmitPower.Valid {
		t.Error(p.TransmitPower)
	}
	if p.BatteryVoltage.Valid {
		t.Error(p.BatteryVoltage)
	}
}

func TestShort(t *testing.T) {
	vector := VECTOR_GOOD[:len(VECTOR_GOOD)-2]
	raw := make([]byte, len(vector)/2)
	n, _ := hex.Decode(raw, []byte(vector))
	if n != len(raw) {
		t.Fatal("unexpected n:", n)
	}
	_, err := rawv2.Parse(raw)
	if !errors.Is(err, rawv2.ErrorPacketTooSmall) {
		t.Fatal(err)
	}
}
