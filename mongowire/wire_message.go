package mongowire

import (
	"github.com/pkg/errors"
	"github.com/tychoish/mongorpc/bson"
	"github.com/tychoish/mongorpc/model"
)

type opMessageSection interface {
	Type() uint8
	Name() string
	DB() string
	Documents() []bson.Simple
	Serialize() []byte
}

type opMessagePayloadType0 struct {
	PayloadType uint8
	Document    bson.Simple
}

func (p *opMessagePayloadType0) Type() uint8              { return 0 }
func (p *opMessagePayloadType0) Name() string             { return "" }
func (p *opMessagePayloadType0) Documents() []bson.Simple { return []bson.Simple{p.Document} }
func (p *opMessagePayloadType0) DB() string {
	m, err := sec.Document.ToBSONM()
	if err != nil {
		return ""
	}

	return m["$db"]
}

type opMessagePayloadType1 struct {
	PayloadType uint8
	Size        int32
	Identifer   string
	Documents   []bson.Simple
}

func (p *opMessagePayloadType1) Type() uint8              { return 1 }
func (p *opMessagePayloadType1) Name() string             { return p.Identifer }
func (p *opMessagePayloadType1) DB() string               { return "" }
func (p *opMessagePayloadType1) Documents() []bson.Simple { return p.Documents }

func (m *opMessage) Header() MessageHeader { return m.header }
func (m *opMessage) HasResponse() bool     { return m.Flags > 1 }
func (m *opMessage) Scope() *OpScope {
	s := &OpScope{
		Type: m.header.OpCode,
	}

	for _, sec := range m.Items {
		if sec.Type() == 0 {
			s.Context = sec.DB()
			if s.Context == "" {
				continue
			}
			break
		}

		if sec.Type() == 1 {
			s.Command = sec.Name()
		}
	}

	return nil
}

func (m *opMessage) Serialize() []byte {

}

// TODO:
//   - finish implementation of parseMsgMessageBody
//   - implement message interface
//      - Serialize

func NewOpMessage(moreToCome bool, document bson.Simple, items ...model.SequenceItem) Message {
	msg := &opMessage{
		header: MessageHeader{
			OpCode:    OP_MSG,
			RequestID: 19,
		},
		Flags: 1,
		Items: []opMessageSection{
			opMessagePayloadType0{
				PayloadType: 0,
				Document:    document,
			},
		},
	}

	for idx := range items {
		item := items[idx]
		it := opMessagePayloadType1{
			PayloadType: 1,
			Identifer:   item.Identifier,
		}
		for _, i := range item.Ducments {
			it.Size += i.Size
		}
		msg.Items = append(msg.Items, it)
	}

	return msg
}

func (h *MessageHeader) parseMsgMessageBody(body []byte) (Message, error) {
	return nil, errors.New("op_message parsing not implemented")
}
