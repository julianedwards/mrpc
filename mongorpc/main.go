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
				// requestHeader := msg.Header()
				isOk := bsonx.EC.Int32("ok", 1)

				doc := bsonx.NewDocument(isOk)
				docBytes, err := doc.MarshalBSON()
				docSimple := bson.Simple{BSON: docBytes, Size: int32(len(docBytes))}
        _, _ = pp.Print(err)

				newReply := mongowire.NewReply(0, 0, 0, 1, []bson.Simple{docSimple}) 
				w.Write(newReply.Serialize())
    }); err != nil {
        grip.Error(err)
        return
    }

    // db.whatsmyuri()
    if err := srv.RegisterOperation(&mongowire.OpScope{
        Type:    mongowire.OP_COMMAND,
        Context: "admin",
        Command: "whatsmyuri",
    }, func(ctx context.Context, w io.Writer, msg mongowire.Message) {
				// requestHeader := msg.Header()
				myUri := bsonx.EC.String("you", "localhost:12345")

				doc := bsonx.NewDocument(myUri)
				docBytes, err := doc.MarshalBSON()
				docSimple := bson.Simple{BSON: docBytes, Size: int32(len(docBytes))}
        _, _ = pp.Print(err)

				newReply := mongowire.NewReply(0, 0, 0, 1, []bson.Simple{docSimple}) 
				w.Write(newReply.Serialize())
    }); err != nil {
        grip.Error(err)
        return
    }

    grip.Error(srv.Run(ctx))
}
