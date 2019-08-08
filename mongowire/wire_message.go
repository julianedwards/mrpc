package mongowire

import (
	"github.com/mongodb/ftdc/bsonx"
	"github.com/pkg/errors"
	"github.com/tychoish/mongorpc/model"
)

type opMessageSection interface {
	Type() uint8
	Name() string
	DB() string
	Documents() []*bsonx.Document
	Serialize() []byte
}

type opMessagePayloadType0 struct {
	PayloadType uint8
	Document    *bsonx.Document
}

func (p *opMessagePayloadType0) Type() uint8                  { return 0 }
func (p *opMessagePayloadType0) Name() string                 { return "" }
func (p *opMessagePayloadType0) Documents() []*bsonx.Document { return []*bsonx.Document{p.Document} }
func (p *opMessagePayloadType0) Serialize() []byte            { return nil }
func (p *opMessagePayloadType0) DB() string {
	key, err := p.Document.LookupErr("$db")
	if err != nil {
		return ""
	}

	val, ok := key.StringValueOK()
	if !ok {
		return ""
	}

	return val
}

type opMessagePayloadType1 struct {
	PayloadType uint8
	Size        int32
	Identifer   string
	Payload     []*bsonx.Document
}

func (p *opMessagePayloadType1) Type() uint8                  { return 1 }
func (p *opMessagePayloadType1) Name() string                 { return p.Identifer }
func (p *opMessagePayloadType1) DB() string                   { return "" }
func (p *opMessagePayloadType1) Documents() []*bsonx.Document { return p.Payload }
func (p *opMessagePayloadType1) Serialize() []byte            { return nil }

func (m *opMessage) Header() MessageHeader { return m.header }
func (m *opMessage) HasResponse() bool     { return m.Flags > 1 }
func (m *opMessage) Scope() *OpScope {
	return &OpScope{
		Type: m.header.OpCode,
	}
}

func (m *opMessage) Serialize() []byte { return nil }

// TODO:
//   - finish implementation of parseMsgMessageBody
//   - implement message interface
//      - Serialize

func NewOpMessage(moreToCome bool, document *bsonx.Document, items ...model.SequenceItem) Message {
	msg := &opMessage{
		header: MessageHeader{
			OpCode:    OP_MSG,
			RequestID: 19,
		},
		Flags: 1,
		Items: []opMessageSection{
			&opMessagePayloadType0{
				PayloadType: 0,
				Document:    document,
			},
		},
	}

	for idx := range items {
		item := items[idx]
		it := &opMessagePayloadType1{
			PayloadType: 1,
			Identifer:   item.Identifer,
		}
		for _, i := range item.Documents {
			it.Size += int32(getDocSize(i))
		}
		msg.Items = append(msg.Items, it)
	}

	return msg
}

func (h *MessageHeader) parseMsgMessageBody(body []byte) (Message, error) {
	return nil, errors.New("op_message parsing not implemented")
}
