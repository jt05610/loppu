package modbus_test

import (
	mb "github.com/jt05610/loppu/hardware/modbus"
	"testing"
)

func TestCRC16(t *testing.T) {
	for _, tc := range []struct {
		Name   string
		Input  []byte
		Expect uint16
	}{
		{"Okay", []byte{0x01, 0x01, 0xfe, 0xed, 0xbe, 0xAD}, 0xCA2D},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			actual := mb.CRC16(tc.Input)
			if tc.Expect != actual {
				t.Fail()
			}
		})
	}
}
