package mystruct

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/sharehdm/myleaf/chanrpc"
	"github.com/sharehdm/myleaf/log"
)

type Processor struct {
	msgInfo map[uint16]*MsgInfo
}

type MsgInfo struct {
	msgType   uint16
	msgRouter *chanrpc.Server
}

type MsgHandler func([]interface{})

type MsgRaw struct {
	msgType    uint16
	msgRawData []byte
}

func NewProcessor() *Processor {
	p := new(Processor)
	p.msgInfo = make(map[uint16]*MsgInfo)
	return p
}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *Processor) Register(msgtype uint16) uint16 {
	i := new(MsgInfo)
	i.msgType = msgtype
	p.msgInfo[msgtype] = i
	return msgtype
}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *Processor) SetRouter(msgtype uint16, msgRouter *chanrpc.Server) {
	i, ok := p.msgInfo[msgtype]
	if !ok {
		log.Fatal("message %v not registered", msgtype)
	}
	i.msgRouter = msgRouter
}

// goroutine safe
func (p *Processor) Route(msg interface{}, userData interface{}) error {
	if msgRaw, ok := msg.(MsgRaw); ok {
		i, ok := p.msgInfo[msgRaw.msgType]
		if !ok {
			return fmt.Errorf("message %v not registered", msgRaw.msgType)
		}
		if i.msgRouter != nil {
			i.msgRouter.Go(msgRaw.msgType, msgRaw.msgRawData, userData)
		}
		return nil
	}
	return nil
}

// goroutine safe
func (p *Processor) Unmarshal(data []byte) (interface{}, error) {
	msgtype := binary.LittleEndian.Uint16(data)
	_, ok := p.msgInfo[msgtype]
	if !ok {
		return nil, fmt.Errorf("message %v not registered", msgtype)
	}
	return MsgRaw{msgtype, data}, nil
}

// goroutine safe
func (p *Processor) Marshal(msg interface{}) ([][]byte, error) {
	bufs := new(bytes.Buffer)
	if err := binary.Write(bufs, binary.LittleEndian, msg); err != nil {
		fmt.Println("err: ", err)
	}
	return [][]byte{bufs.Bytes()}, nil
}
