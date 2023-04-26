package modbus

import (
	"errors"
	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"
	"go.uber.org/zap"
)

type SerialOpt struct {
	Port     string          `yaml:"port"`
	VID      string          `yaml:"VID"`
	PID      string          `yaml:"PID"`
	Baud     int             `yaml:"baud"`
	Parity   serial.Parity   `yaml:"parity"`
	DataBits int             `yaml:"dataBits"`
	StopBits serial.StopBits `yaml:"stopBits"`
}

type Serial struct {
	port serial.Port
	log  *zap.Logger
	opt  *SerialOpt
}

var DefaultSerial = &SerialOpt{
	Port:     "",
	VID:      "1A86",
	PID:      "7523",
	Baud:     19200,
	Parity:   serial.NoParity,
	DataBits: 8,
	StopBits: serial.TwoStopBits,
}

var NoPort = errors.New("no port given to NewSerial")

func NewSerial(opt *SerialOpt, log *zap.Logger) (*Serial, error) {
	s := &Serial{log: log}
	var err error
	if opt.Port == "" {
		if opt.PID == "" && opt.VID == "" {
			return nil, NoPort
		}
		ports, err := enumerator.GetDetailedPortsList()
		if err != nil {
			panic(err)
		}
		vCheck := len(opt.PID) > 0
		for _, port := range ports {
			if port.PID == opt.PID {
				opt.Port = port.Name
				if vCheck {
					if port.VID == opt.VID {
						opt.Port = port.Name
						s.log.Info("found port", zap.String("port", port.Name))
						break
					}
				} else {
					s.log.Info("found port", zap.String("port", port.Name))
					break
				}
			}
		}
	}
	if len(opt.Port) == 0 {
		panic(errors.New("no port found"))
	}
	s.opt = opt

	return s, err
}

func (s *Serial) Open() error {
	mode := &serial.Mode{
		BaudRate: s.opt.Baud,
		Parity:   s.opt.Parity,
		DataBits: s.opt.DataBits,
		StopBits: s.opt.StopBits,
	}
	var err error
	s.port, err = serial.Open(s.opt.Port, mode)
	if err == nil {
		s.log.Info("connected port", zap.String("port", s.opt.Port))
	}
	return err
}

func (s *Serial) Close() {
	err := s.port.Close()
	if err != nil {
		panic(err)
	}
}

func (s *Serial) Read(p []byte) (n int, err error) {
	n, err = s.port.Read(p)
	s.log.Info("received", zap.ByteString("pdu", p[:n]))
	return n, err
}

func (s *Serial) Write(p []byte) (n int, err error) {
	s.log.Info("sending", zap.ByteString("pdu", p))
	n, err = s.port.Write(p)
	return n, err
}
