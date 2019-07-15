package main

import (
    "context"
    "io"

    "github.com/k0kubun/pp"
    "github.com/mongodb/grip"
    "github.com/tychoish/mongorpc"
    "github.com/tychoish/mongorpc/mongowire"
    "github.com/tychoish/mongorpc/bson"
		"github.com/mongodb/ftdc/bsonx"
)

func main() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    srv := mongorpc.NewService("localhost", 12345)

    // db.runCommand({whatever})
    if err := srv.RegisterOperation(&mongowire.OpScope{
        Type:    mongowire.OP_QUERY,
        Context: "test.$cmd",
    }, func(ctx context.Context, w io.Writer, msg mongowire.Message) {
        _, _ = pp.Print(msg.Header())
    }); err != nil {
        grip.Error(err)
        return
    }

    // db.isMaster()
    if err := srv.RegisterOperation(&mongowire.OpScope{
        Type:    mongowire.OP_COMMAND,
        Context: "admin",
        Command: "isMaster",
    }, func(ctx context.Context, w io.Writer, msg mongowire.Message) {
				requestHeader := msg.Header()
        _, _ = pp.Print(requestHeader)
				isMaster := bsonx.EC.Boolean("isMaster", true)

				doc := bsonx.NewDocument(isMaster)
				docBytes, err := doc.MarshalBSON()
				docSimple := bson.Simple{BSON: docBytes, Size: int32(len(docBytes))}
        _, _ = pp.Print(err)

				metadata := bsonx.NewDocument()
				metadataBytes, err := metadata.MarshalBSON()
				metadataSimple := bson.Simple{BSON: metadataBytes, Size: int32(len(metadataBytes))}
        _, _ = pp.Print(err)

				// outputDocs := bsonx.NewDocument()
				// outputDocsBytes, err := doc.MarshalBSON()
				// outputDocsSimple := bson.Simple{BSON: docBytes, Size: int32(len(outputDocsBytes))}

				responseHeader := mongowire.MessageHeader{ResponseTo: requestHeader.RequestID, RequestID: 324, OpCode: mongowire.OP_COMMAND_REPLY}
				commandReplyMessage := mongowire.CommandReplyMessage{CommandReplyHeader: responseHeader, CommandReply: docSimple, Metadata: metadataSimple}
				w.Write(commandReplyMessage.Serialize())
    }); err != nil {
        grip.Error(err)
        return
    }

    grip.Error(srv.Run(ctx))
}
