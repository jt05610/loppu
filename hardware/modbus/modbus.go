package modbus

import (
	"encoding/binary"
	"github.com/jt05610/loppu"
)

type FuncCode byte

const (
	ReadCoilsFC            FuncCode = 0x01
	ReadDiscreteInputsFC   FuncCode = 0x02
	ReadHoldingRegistersFC FuncCode = 0x03
	ReadInputRegistersFC   FuncCode = 0x04
	WriteCoilFC            FuncCode = 0x05
	WriteHoldingRegisterFC FuncCode = 0x06
	DiagFC                 FuncCode = 0x08
)

type MBusPDU struct {
	FuncCode
	Body []byte
}

func (m *MBusPDU) CRC16() uint16 {
	return 0
}

func (m *MBusPDU) Addr() loppu.Addr {
	return nil
}

func (m *MBusPDU) Header() []byte {
	return []byte{byte(m.FuncCode)}
}

func (m *MBusPDU) Data() []byte {
	return m.Body
}

func (m *MBusPDU) Read(p []byte) (n int, err error) {
	p[0] = byte(m.FuncCode)
	for i := 1; i < len(p); i++ {
		p[i] = m.Body[i-1]
	}
	return len(p), nil
}

func (m *MBusPDU) Write(p []byte) (n int, err error) {
	m.FuncCode = FuncCode(p[0])
	if len(p) > 1 {
		m.Body = p[1:]
	}
	return len(p), nil

}

func ReadCoils(addr uint16, quantity uint16) *MBusPDU {
	m := &MBusPDU{
		FuncCode: ReadCoilsFC,
		Body:     make([]byte, 4),
	}
	binary.BigEndian.PutUint32(m.Body, (uint32(addr)<<16)+uint32(quantity))
	return m
}

func ReadDiscreteInputs(addr uint16, quantity uint16) *MBusPDU {
	m := &MBusPDU{
		FuncCode: ReadDiscreteInputsFC,
		Body:     make([]byte, 4),
	}
	binary.BigEndian.PutUint32(m.Body, (uint32(addr)<<16)+uint32(quantity))
	return m
}

func ReadHoldingRegisters(addr uint16, quantity uint16) *MBusPDU {
	m := &MBusPDU{
		FuncCode: ReadHoldingRegistersFC,
		Body:     make([]byte, 4),
	}
	binary.BigEndian.PutUint32(m.Body, (uint32(addr)<<16)+uint32(quantity))
	return m
}

func ReadInputRegisters(addr uint16, quantity uint16) *MBusPDU {
	m := &MBusPDU{
		FuncCode: ReadInputRegistersFC,
		Body:     make([]byte, 4),
	}
	binary.BigEndian.PutUint32(m.Body, (uint32(addr)<<16)+uint32(quantity))
	return m
}

func WriteCoil(addr uint16, value uint16) *MBusPDU {
	if value != 0 {
		value = 0xFF00
	}
	m := &MBusPDU{
		FuncCode: WriteCoilFC,
		Body:     make([]byte, 4),
	}
	binary.BigEndian.PutUint32(m.Body, (uint32(addr)<<16)+uint32(value))
	return m
}

func WriteRegister(addr uint16, value uint16) *MBusPDU {
	m := &MBusPDU{
		FuncCode: WriteHoldingRegisterFC,
		Body:     make([]byte, 4),
	}
	binary.BigEndian.PutUint32(m.Body, (uint32(addr)<<16)+uint32(value))
	return m
}

func Echo(data ...byte) *MBusPDU {
	return &MBusPDU{
		FuncCode: DiagFC,
		Body:     append([]byte{0x00, 0x00}, data...),
	}
}
