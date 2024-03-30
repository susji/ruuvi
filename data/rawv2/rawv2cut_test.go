package rawv2_test

// Test vectors below are from Ruuvi's documentation at
// https://docs.ruuvi.com/communication/bluetooth-advertisements/data-format-c5-rawv2

import (
	"encoding/hex"
	"errors"
	"testing"

	"github.com/susji/ruuvi/data"
	"github.com/susji/ruuvi/data/rawv2"
)

const (
	VECTOR_GOOD_CUT = "C512FC5394C37CAC364200CDCBB8334C884F"
	VECTOR_MAX_CUT  = "C57FFFFFFEFFFEFFDEFEFFFECBB8334C884F"
	VECTOR_MIN_CUT  = "C58001000000000000000000CBB8334C884F"
	VECTOR_BAD_CUT  = "C58000FFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"
)

func TestBasicCut(t *testing.T) {
	raw := make([]byte, len(VECTOR_GOOD_CUT)/2)
	n, _ := hex.Decode(raw, []byte(VECTOR_GOOD_CUT))
	if n != len(raw) {
		t.Fatal("unexpected n:", n)
	}
	p, err := rawv2.Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	if p.Type != data.VERSION_CUTRAWV2 {
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
	if p.AccelerationX.Valid {
		t.Error(p.AccelerationX)
	}
	if p.AccelerationY.Valid {
		t.Error(p.AccelerationY)
	}
	if p.AccelerationZ.Valid {
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

func TestMaximumCut(t *testing.T) {
	raw := make([]byte, len(VECTOR_MAX_CUT)/2)
	n, _ := hex.Decode(raw, []byte(VECTOR_MAX_CUT))
	if n != len(raw) {
		t.Fatal("unexpected n:", n)
	}
	p, err := rawv2.Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	if p.Type != data.VERSION_CUTRAWV2 {
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
	if p.AccelerationX.Valid {
		t.Error(p.AccelerationX)
	}
	if p.AccelerationY.Valid {
		t.Error(p.AccelerationY)
	}
	if p.AccelerationZ.Valid {
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

func TestMinimumCut(t *testing.T) {
	raw := make([]byte, len(VECTOR_MIN_CUT)/2)
	n, _ := hex.Decode(raw, []byte(VECTOR_MIN_CUT))
	if n != len(raw) {
		t.Fatal("unexpected n:", n)
	}
	p, err := rawv2.Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	if p.Type != data.VERSION_CUTRAWV2 {
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
	if p.AccelerationX.Valid {
		t.Error(p.AccelerationX)
	}
	if p.AccelerationY.Valid {
		t.Error(p.AccelerationY)
	}
	if p.AccelerationZ.Valid {
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

func TestInvalidCut(t *testing.T) {
	raw := make([]byte, len(VECTOR_BAD_CUT)/2)
	n, _ := hex.Decode(raw, []byte(VECTOR_BAD_CUT))
	if n != len(raw) {
		t.Fatal("unexpected n:", n)
	}
	p, err := rawv2.Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	if p.Type != data.VERSION_CUTRAWV2 {
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

func TestShortCut(t *testing.T) {
	vector := VECTOR_GOOD_CUT[:len(VECTOR_GOOD_CUT)-2]
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
