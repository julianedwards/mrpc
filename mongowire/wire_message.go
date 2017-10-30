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
	Collection() string
	Documents() []bson.Simple
}

type opMessagePayloadType0 struct {
	PayloadType uint8
	Document    bson.Simple
}

type opMessagePayloadType1 struct {
	PayloadType uint8
	Size        int32
	Identifer   string
	Documents   []bson.Simple
}

// TODO:
//   - implement section interface for payload type 0
//      - Type
//      - Name
//      - DB
//      - Collection
//      - Documents
//   - implement section interface for payload type 1
//      - Type
//      - Name
//      - DB
//      - Collection
//      - Documents
//   - finish implementation of parseMsgMessageBody
//   - implement message interface
//      - Header
//      - Serialize
//      - HasResponse
//      - Scope

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
