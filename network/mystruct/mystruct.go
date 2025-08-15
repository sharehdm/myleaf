package mystruct

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/sharehdm/myleaf/chanrpc"
	"github.com/sharehdm/myleaf/log"
)

type Processor struct {
	msgInfo map[int]*MsgInfo
}

type MsgInfo struct {
	msgType   int
	msgRouter *chanrpc.Server
}

type MsgHandler func([]interface{})

type MsgRaw struct {
	msgType    int
	msgRawData []byte
}

func NewProcessor() *Processor {
	p := new(Processor)
	p.msgInfo = make(map[int]*MsgInfo)
	return p
}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *Processor) Register(msgtype int) int {
	i := new(MsgInfo)
	i.msgType = msgtype
	p.msgInfo[msgtype] = i
	return msgtype
}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *Processor) SetRouter(msgtype int, msgRouter *chanrpc.Server) {
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
	msgtype := (int)(binary.LittleEndian.Uint16(data))
	_, ok := p.msgInfo[msgtype]
	if !ok {
		return nil, fmt.Errorf("message %v not registered", msgtype)
	}
	return MsgRaw{msgtype, data}, nil
}

// goroutine safe
func (p *Processor) Marshal(mid uint16, sid uint16, msg interface{}) ([][]byte, error) {
	topbys := make([]byte, 4)
	binary.LittleEndian.PutUint16(topbys, mid)
	binary.LittleEndian.PutUint16(topbys[2:], sid)
	if msg != nil {
		bufs := new(bytes.Buffer)
		if err := binary.Write(bufs, binary.LittleEndian, msg); err != nil {
			fmt.Println("err: ", err)
		}
		return [][]byte{topbys, bufs.Bytes()}, nil
	} else {
		return [][]byte{topbys}, nil
	}
}
