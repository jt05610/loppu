package modbus

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/jt05610/loppu"
	"github.com/jt05610/loppu/hardware"
	"github.com/jt05610/loppu/yaml"
	"io"
	"net/http"
	"net/url"
	"time"
)

type ParamType string

type Param struct {
	Type        ParamType
	Description string `yaml:"desc"`
}

type Handler struct {
	Name        string              `yaml:"name"`
	Description string              `yaml:"desc"`
	Params      []map[string]*Param `yaml:"params,omitempty"`
}

type MetaData struct {
	Node    string    `yaml:"node"`
	Author  string    `yaml:"author"`
	Address MBAddress `yaml:"address"`
	Date    string    `yaml:"date"`
}

type MBusNode struct {
	MetaData    *MetaData             `yaml:"meta"`
	Tables      map[string][]*Handler `yaml:"tables"`
	Diag        []*Handler            `yaml:"diag,omitempty"`
	client      *Client
	rfLookup    map[string]map[string]func(uint16, uint16) *MBusPDU
	addrLookup  map[string]uint16
	paramLookup map[string]string
}

func (n *MBusNode) Proto(p ...hardware.Proto) hardware.Proto {
	if p != nil {
		n.client = p[0].(*Client)
	}
	return n.client
}

func (n *MBusNode) Meta() loppu.MetaData {
	return n.MetaData
}

var lookup = map[string]map[string]func(uint16, uint16) *MBusPDU{
	http.MethodGet: {
		"coils":             ReadCoils,
		"discrete_inputs":   ReadDiscreteInputs,
		"holding_registers": ReadHoldingRegisters,
		"input_registers":   ReadInputRegisters,
	},
	http.MethodPost: {
		"coils":             WriteCoil,
		"holding_registers": WriteRegister,
	},
}

func (n *MBusNode) Endpoints(baseURL string) []*loppu.Endpoint {
	res := make([]*loppu.Endpoint, 0)
	n.paramLookup = make(map[string]string)
	for name, handlers := range n.Tables {
		for _, h := range handlers {
			route, err := url.JoinPath(baseURL, n.MetaData.Node, h.Name)
			if err != nil {
				panic(err)
			}
			res = append(res, &loppu.Endpoint{
				Route:       route,
				Method:      http.MethodGet,
				Description: h.Description,
				Param:       nil,
				Func:        strcase.ToCamel(fmt.Sprintf("Get %s", h.Name)),
			})
			if name == "coils" || name == "holding_registers" {
				var param *loppu.EndpointParam
				if len(h.Params) > 0 {
					param = &loppu.EndpointParam{}
					for paramName, p := range h.Params[0] {
						param.Name = strcase.ToLowerCamel(paramName)
						param.Type = string(p.Type)
						param.Description = p.Description
						param.NameCap = strcase.ToCamel(paramName)
						param.Tag = fmt.Sprintf("`json:\"%s\"`", param.Name)
					}
					n.paramLookup[route] = param.Name
				} else {
					param = nil
				}
				res = append(res, &loppu.Endpoint{
					Route:       route,
					Method:      http.MethodPost,
					Description: h.Description,
					Param:       param,
					Func:        strcase.ToCamel(fmt.Sprintf("Post %s", h.Name)),
				})
			}
		}
	}
	return res
}

func (n *MBusNode) handlers() map[string]http.HandlerFunc {
	_ = n.Endpoints("/")
	res := make(map[string]http.HandlerFunc)
	n.rfLookup = make(map[string]map[string]func(uint16, uint16) *MBusPDU)
	n.rfLookup[http.MethodGet] = make(map[string]func(uint16, uint16) *MBusPDU)
	n.rfLookup[http.MethodPost] = make(map[string]func(uint16, uint16) *MBusPDU)
	n.addrLookup = make(map[string]uint16)
	for name, handlers := range n.Tables {
		for i, h := range handlers {
			endpoint := fmt.Sprintf("/%s/%s", n.MetaData.Node, h.Name)
			if reqFunc, ok := lookup[http.MethodGet][name]; ok {
				n.rfLookup[http.MethodGet][endpoint] = reqFunc
				n.addrLookup[endpoint] = uint16(i)
			} else {
				panic(errors.New("failed to find get request formatter"))
			}
			if name == "coils" || name == "holding_registers" {
				if rf, ok := lookup[http.MethodPost][name]; ok {
					n.rfLookup[http.MethodPost][endpoint] = rf
				} else {
					panic(errors.New("failed to find post request formatter"))
				}
			}
			res[endpoint] = func(w http.ResponseWriter, r *http.Request) {
				if reqFunc, ok := n.rfLookup[r.Method][r.RequestURI]; ok {
					if reqFunc == nil {
						http.Error(w, "internal server error", http.StatusInternalServerError)
					} else {
						var bytes []byte
						if r.Method == http.MethodGet {
							pack, err := n.client.Request(r.Context(),
								n.MetaData.Address, reqFunc(n.addrLookup[r.
									RequestURI], 1))
							pdu := pack.(*MBusPDU)
							if err != nil || res == nil {
								http.Error(w, err.Error(), http.StatusInternalServerError)
								return
							}
							if pdu != nil && pdu.FuncCode < ReadHoldingRegistersFC {
								bytes, err = json.Marshal(map[string]uint16{
									"result": uint16(pdu.Body[1]),
								})

							} else {
								if pdu.Body[1] == 2 {
									bytes, err = json.Marshal(map[string]uint16{
										"result": binary.BigEndian.Uint16(
											pdu.Body[2:]),
									})
								} else {
									rr := make([]uint16, 0)
									for i := uint8(0); i < pdu.Body[1]; i += 2 {
										rr = append(rr, binary.BigEndian.
											Uint16(pdu.Body[2+i:]))
									}
									bytes, err = json.Marshal(map[string][]uint16{
										"result": rr,
									})
								}
							}

						} else {
							dec := json.NewDecoder(r.Body)
							req := make(map[string]uint16)
							err := dec.Decode(&req)
							if err != nil {
								if _, found := n.paramLookup[r.RequestURI]; found {
									http.Error(w, "bad request", http.StatusBadRequest)
									return
								} else {
									if err.Error() != "EOF" {
										http.Error(w, "bad request", http.StatusBadRequest)
										return
									}
								}
							}
							var reqPDU *MBusPDU
							if len(req) > 0 {
								if param, found := n.paramLookup[r.RequestURI]; !found {
									http.Error(w,
										fmt.Sprintf("must include %s", param),
										http.StatusBadRequest)
								} else {
									if p, found := req[param]; found {
										reqPDU = reqFunc(n.addrLookup[r.RequestURI], p)
									} else {
										http.Error(w,
											fmt.Sprintf("must include %s", param),
											http.StatusBadRequest)
									}
								}
							} else {
								reqPDU = reqFunc(n.addrLookup[r.RequestURI], 1)
							}

							pack, err := n.client.Request(r.Context(),
								n.MetaData.Address, reqPDU)
							if err != nil {
								http.Error(w, "server error", http.StatusInternalServerError)
								return
							}
							pdu := pack.(*MBusPDU)
							if pdu.FuncCode != reqPDU.FuncCode {
								http.Error(w, "Modbus server error", http.StatusInternalServerError)
							} else {
								bytes = []byte("ok")
							}
						}
						w.WriteHeader(http.StatusOK)
						_, err := w.Write(bytes)
						if err != nil {
							panic(err)
						}
					}
				}
			}
		}
	}
	return res
}

func (n *MBusNode) Register(srv *http.ServeMux) {
	for name, handler := range n.handlers() {
		srv.HandleFunc(name, handler)
	}
	for _, handler := range n.Diag {
		if handler.Name == "echo" {
			endpoint := fmt.Sprintf("/%s/%s", n.MetaData.Node, "echo")
			srv.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
				dec := json.NewDecoder(r.Body)
				body := make(map[string][]byte)
				err := dec.Decode(&body)
				if err != nil {
					http.Error(w, "bad request", http.StatusBadRequest)
				}
				req := Echo(body["message"]...)
				pack, err := n.client.Request(r.Context(),
					n.MetaData.Address, req)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				res := pack.(*MBusPDU)
				var bytes []byte
				if res.FuncCode != req.FuncCode {
					http.Error(w, "Modbus server error", http.StatusInternalServerError)
				} else {
					bytes = res.Body
				}
				w.WriteHeader(http.StatusOK)
				_, err = w.Write(bytes)
				if err != nil {
					panic(err)
				}
			})
		}
	}
}

func (n *MBusNode) Load(r io.Reader) error {
	l := yaml.NodeService[MBusNode]{}
	v, err := l.Load(r)
	n.MetaData = v.MetaData
	n.Tables = v.Tables
	return err
}

func (n *MBusNode) Flush(w io.Writer) error {
	l := yaml.NodeService[MBusNode]{}
	return l.Flush(w, n)
}

func NewMBusNode(name string, address byte) hardware.Node {
	return &MBusNode{
		MetaData: &MetaData{
			Node:    name,
			Author:  loppu.Username(),
			Address: MBAddress(address),
			Date:    time.Now().String(),
		},
		Tables: map[string][]*Handler{
			"discrete_inputs": {
				&Handler{
					Name:        "discrete_input_1",
					Description: "Discrete inputs are binary read only registers",
					Params:      nil,
				},
				&Handler{
					Name:        "discrete_input_2",
					Description: "Add a new one like this",
					Params:      nil,
				},
			},
			"coils": {
				&Handler{
					Name:        "coil_1",
					Description: "Coils are binary read/write registers. They are used to execute on/off functions on the target device",
					Params: []map[string]*Param{
						{
							"value": &Param{
								Type:        "int",
								Description: "Writing 0 will write 0 to the device. Writing anything else will write 1.",
							},
						},
					},
				},
				&Handler{
					Name:        "coil_2",
					Description: "They can also be parameter-less if you want",
					Params:      nil,
				},
			},
			"input_registers": {
				&Handler{
					Name:        "input_register_1",
					Description: "Input registers are 16-bit read only registers",
					Params:      nil,
				},
				&Handler{
					Name:        "input_register_2",
					Description: "Add a new one like this",
					Params:      nil,
				},
			},
			"holding_registers": {
				&Handler{
					Name:        "holding_register_1",
					Description: "Holding registers are 16-bit read/write registers. They are used to set variables on the device, and can execute a if desired",
					Params: []map[string]*Param{
						{
							"value": &Param{
								Type:        "int",
								Description: "Whatever you write will be converted to a uint16",
							},
						},
					},
				},
				&Handler{
					Name:        "holding_register_2",
					Description: "Add a new one like this",
					Params: []map[string]*Param{
						{
							"value": &Param{
								Type:        "int",
								Description: "You technically don't need to include a parameter but you probably should",
							},
						},
					},
				},
			},
		},
	}
}
