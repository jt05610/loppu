package modbus

import (
	"encoding/binary"
	"github.com/jt05610/loppu"
	"github.com/jt05610/loppu/hardware"
)

type MBAddress byte

func (a MBAddress) Bytes() []byte {
	return []byte{byte(a)}
}

func (a MBAddress) Byte() byte {
	return byte(a)
}

func (a MBAddress) String() string {
	return string(a)
}

func NewMBAddress(v byte) loppu.Addr {
	return MBAddress(v)
}

type SerialPDU struct {
	Address loppu.Addr
	PDU     hardware.Packet
	CRC     uint16
}

func (s *SerialPDU) CRC16() uint16 {
	return s.CRC
}

func (s *SerialPDU) Addr() loppu.Addr {
	return s.Address
}

func (s *SerialPDU) Header() []byte {
	return append(s.Address.Bytes(), s.PDU.Header()...)
}

func (s *SerialPDU) Data() []byte {
	return s.PDU.Data()
}

func (s *SerialPDU) Read(p []byte) (n int, err error) {
	p[0] = s.Address.Byte()
	n, err = s.PDU.Read(p[1 : len(p)-2])
	binary.LittleEndian.PutUint16(p[len(p)-2:], s.CRC)
	return n + 3, err
}

func (s *SerialPDU) Write(p []byte) (n int, err error) {
	s.Address = NewMBAddress(p[0])
	s.PDU = &MBusPDU{}
	n, err = s.PDU.Write(p[1 : len(p)-2])
	s.CRC = binary.LittleEndian.Uint16(p[len(p)-2:])
	return n + 3, err
}

func NewSerialPDU(addr byte, pdu hardware.Packet) hardware.Packet {
	tmp := make([]byte, len(pdu.Data())+2)
	tmp[0] = addr
	_, _ = pdu.Read(tmp[1:])
	return &SerialPDU{
		Address: NewMBAddress(addr),
		PDU:     pdu,
		CRC:     CRC16(tmp),
	}
}
