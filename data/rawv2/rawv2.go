// Package rawv2 parses and validates Ruuvi's Data formats 5 (RAWv2) and C5
// (Cut-RAWv2) as described in [RAWv2] and [Cut-RAWv2], respectively.
//
// Each value in a packet is considered usable if value-specific Valid is true.
//
// [RAWv2]: https://docs.ruuvi.com/communication/bluetooth-advertisements/data-format-5-rawv2
// [Cut-RAWv2]: https://docs.ruuvi.com/communication/bluetooth-advertisements/data-format-c5-rawv2
package rawv2

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

var (
	ErrorPacketTooSmall = errors.New("packet too small")
	ErrorPacketNotV2    = errors.New("packet has wrong version")
	ErrorBadMAC         = errors.New("not a valid MAC address")
)

const (
	VERSION_PLAIN = 0x05
	VERSION_CUT   = 0xC5
)

// Temperature is degrees Celsius.
type Temperature struct {
	Valid bool
	Value float32
}

// Humidity is relative air humidity in percentage.
type Humidity struct {
	Valid bool
	Value float32
}

// Pressure is Pascals.
type Pressure struct {
	Valid bool
	Value uint32
}

// Acceleration is Milli-G.
type Acceleration struct {
	Valid bool
	Value int16
}

// BatteryVoltage is Volts.
type BatteryVoltage struct {
	Valid bool
	Value float32
}

// TransmitPower is dBm.
type TransmitPower struct {
	Valid bool
	Value int16
}

// MovementCounter is raw count.
type MovementCounter struct {
	Valid bool
	Value uint8
}

// SequenceNumber is raw count.
type SequenceNumber struct {
	Valid bool
	Value uint16
}

// MAC is the device's MAC address.
type MAC struct {
	Value [6]byte
}

func (v MAC) String() string {
	return fmt.Sprintf(
		"%02x:%02x:%02x:%02x:%02x:%02x",
		v.Value[0], v.Value[1], v.Value[2], v.Value[3], v.Value[4], v.Value[5])
}

func (v MAC) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.String())
}

func (v *MAC) UnmarshalJSON(b []byte) error {
	var m0, m1, m2, m3, m4, m5 byte
	n, err := fmt.Sscanf(
		string(b),
		`"%02x:%02x:%02x:%02x:%02x:%02x"`,
		&m0, &m1, &m2, &m3, &m4, &m5)
	if n != 6 || err != nil {
		n, err = fmt.Sscanf(
			string(b),
			`"%02X:%02X:%02X:%02x:%02X:%02X"`,
			&m0, &m1, &m2, &m3, &m4, &m5)
		if n != 6 || err != nil {
			return fmt.Errorf("%w: %q: %v", ErrorBadMAC, b, err)
		}
	}
	v.Value[0] = m0
	v.Value[1] = m1
	v.Value[2] = m2
	v.Value[3] = m3
	v.Value[4] = m4
	v.Value[5] = m5
	return nil
}

type RuuviRawV2 struct {
	Type            uint8
	Timestamp       time.Time
	Temperature     Temperature
	Humidity        Humidity
	Pressure        Pressure
	AccelerationX   Acceleration
	AccelerationY   Acceleration
	AccelerationZ   Acceleration
	BatteryVoltage  BatteryVoltage
	TransmitPower   TransmitPower
	MovementCounter MovementCounter
	SequenceNumber  SequenceNumber
	MAC             MAC
}

// Parse parses Ruuvi's RAWv2 and Cut-RAWv2 formats. It expects to receive the
// raw payload as recorded from Ruuvi's Bluetooth Advertisements without the
// Manufacturer ID prefix of 0x0499. After a successful parse, it will return a
// pointer to a [RuuviRawV2]. Upon failing, it will return errors derived from
// [ErrorPacketTooSmall] or [ErrorPacketNotV2].
//
// Parse is optimistic so as long as it is given enough data, it will try
// parsing all the specified values.
func Parse(d []byte) (*RuuviRawV2, error) {
	return ParseWithTime(d, time.Now())
}

// ParseWithTime is as [Parse] except it permits passing an arbitrary timestamp
// to be included in the returned struct.
func ParseWithTime(d []byte, t time.Time) (*RuuviRawV2, error) {
	minlen := 0
	switch d[0] {
	case VERSION_PLAIN:
		minlen = 24
	case VERSION_CUT:
		minlen = 18
	default:
		return nil, fmt.Errorf("%w: %d", ErrorPacketNotV2, d[0])
	}
	if len(d) < minlen {
		return nil, fmt.Errorf("%w: %d octets", ErrorPacketTooSmall, len(d))
	}
	r := &RuuviRawV2{Timestamp: t, Type: uint8(d[0])}
	bo := binary.BigEndian
	s := 1
	r.Temperature = newTemperature(bo.Uint16(d[s:]))
	s += 2
	r.Humidity = newHumidity(bo.Uint16(d[s:]))
	s += 2
	r.Pressure = newPressure(bo.Uint16(d[s:]))
	s += 2
	if r.Type == VERSION_PLAIN {
		r.AccelerationX = newAcceleration(bo.Uint16(d[s:]))
		s += 2
		r.AccelerationY = newAcceleration(bo.Uint16(d[s:]))
		s += 2
		r.AccelerationZ = newAcceleration(bo.Uint16(d[s:]))
		s += 2
	}
	p := bo.Uint16(d[s:])
	r.BatteryVoltage = newBatteryVoltage(p)
	r.TransmitPower = newTransmitPower(p)
	s += 2
	r.MovementCounter = newMovementCounter(d[s])
	s += 1
	r.SequenceNumber = newSequenceNumber(bo.Uint16(d[s:]))
	s += 2
	r.MAC = newMAC(d[s : s+6])
	return r, nil
}

func newTemperature(raw uint16) Temperature {
	p := Temperature{Valid: true}
	if raw == 0x8000 {
		p.Valid = false
	}
	p.Value = float32(int16(raw)) * 0.005
	return p
}

func newHumidity(raw uint16) Humidity {
	p := Humidity{Valid: true}
	if raw == 0xffff {
		p.Valid = false
	}
	p.Value = float32(raw) * 0.0025
	return p
}

func newPressure(raw uint16) Pressure {
	p := Pressure{Valid: true}
	if raw == 0xffff {
		p.Valid = false
	}
	p.Value = uint32(raw) + uint32(50000)
	return p
}

func newAcceleration(raw uint16) Acceleration {
	p := Acceleration{Valid: true}
	if raw == 0x8000 {
		p.Valid = false
	}
	p.Value = int16(raw)
	return p
}

func newBatteryVoltage(raw uint16) BatteryVoltage {
	raw = raw >> 5
	p := BatteryVoltage{Valid: true}
	if raw == 0x7ff {
		p.Valid = false
	}
	p.Value = float32(raw)/1000 + 1.6
	return p
}

func newTransmitPower(raw uint16) TransmitPower {
	raw = (raw & 0b11111)
	p := TransmitPower{Valid: true}
	if raw == 0x1f {
		p.Valid = false
	}
	p.Value = int16(raw*2) - 40
	return p
}

func newMovementCounter(raw byte) MovementCounter {
	p := MovementCounter{Valid: true}
	if raw == 0xff {
		p.Valid = false
	}
	p.Value = uint8(raw)
	return p
}

func newSequenceNumber(raw uint16) SequenceNumber {
	p := SequenceNumber{Valid: true}
	if raw == 0xffff {
		p.Valid = false
	}
	p.Value = raw
	return p
}

func newMAC(raw []byte) MAC {
	p := MAC{}
	copy(p.Value[:], raw)
	return p
}
