package test

import (
	"context"
	"errors"
	"github.com/awakari/core/api/grpc/messages"
	"github.com/awakari/core/api/grpc/subscriptions"
	"github.com/awakari/core/api/grpc/writer"
	"github.com/awakari/core/test/config"
	"github.com/awakari/core/test/data"
	"github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"io"
	"testing"
	"time"
)

func Test_MessageDelivery(t *testing.T) {
	//
	cfg, err := config.NewConfigFromEnv()
	require.Nil(t, err)
	connSubs, err := grpc.Dial(cfg.Uri.Subscriptions, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.Nil(t, err)
	defer connSubs.Close()
	clientSubs := subscriptions.NewServiceClient(connSubs)
	connWriter, err := grpc.Dial(cfg.Uri.Writer, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.Nil(t, err)
	defer connWriter.Close()
	clientWriter := writer.NewServiceClient(connWriter)
	//
	mdUser := []string{
		"X-Awakari-Group-Id", "test-group0",
		"X-Awakari-User-Id", "test-user0",
	}
	ctxUser := metadata.AppendToOutgoingContext(context.TODO(), mdUser...)
	//
	var subIds []string
	for _, sub := range data.Subs {
		var resp *subscriptions.CreateResponse
		resp, err = clientSubs.Create(ctxUser, &sub)
		require.Nil(t, err)
		subIds = append(subIds, resp.Id)
	}
	//
	var streamMsgsWrite writer.Service_SubmitMessagesClient
	streamMsgsWrite, err = clientWriter.SubmitMessages(ctxUser)
	require.Nil(t, err)
	err = streamMsgsWrite.Send(&data.Msgs)
	require.Nil(t, err)
	var respWrite *writer.SubmitMessagesResponse
	respWrite, err = streamMsgsWrite.Recv()
	require.Nil(t, err)
	require.Equal(t, len(data.Msgs.Msgs), int(respWrite.AckCount))
	err = streamMsgsWrite.CloseSend()
	require.Nil(t, err)
	//
	cases := map[string]struct {
		subId string
		msgs  []*pb.CloudEvent
	}{
		"disabled": {
			subId: subIds[0],
			msgs:  []*pb.CloudEvent{},
		},
		"exact complete value match for a key": {
			subId: subIds[1],
			msgs: []*pb.CloudEvent{
				data.Msgs.Msgs[4],
			},
		},
		"partial exact match": {
			subId: subIds[2],
			msgs: []*pb.CloudEvent{
				data.Msgs.Msgs[0],
			},
		},
		"basic group condition with \"and\" logic and partial sub-conditions": {
			subId: subIds[3],
			msgs: []*pb.CloudEvent{
				data.Msgs.Msgs[1],
			},
		},
		"basic group condition with \"or\" logic": {
			subId: subIds[4],
			msgs: []*pb.CloudEvent{
				data.Msgs.Msgs[2],
				data.Msgs.Msgs[3],
			},
		},
		"basic group condition with \"and\" logic and a negative sub-condition": {
			subId: subIds[5],
			msgs: []*pb.CloudEvent{
				data.Msgs.Msgs[0],
				data.Msgs.Msgs[1],
			},
		},
		"single symbol wildcard": {
			subId: subIds[6],
			msgs: []*pb.CloudEvent{
				data.Msgs.Msgs[0],
			},
		},
		"multiple symbol wildcard": {
			subId: subIds[7],
			msgs: []*pb.CloudEvent{
				data.Msgs.Msgs[1],
				data.Msgs.Msgs[3],
			},
		},
	}
	//
	for k, c := range cases {
		t.Run(k, func(t *testing.T) {
			var connMsgs *grpc.ClientConn
			connMsgs, err = grpc.Dial(cfg.Uri.Messages, grpc.WithTransportCredentials(insecure.NewCredentials()))
			require.Nil(t, err)
			defer connMsgs.Close()
			clientMsgs := messages.NewServiceClient(connMsgs)
			ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Second)
			defer cancel()
			var streamRcv messages.Service_ReceiveClient
			streamRcv, err = clientMsgs.Receive(ctx)
			require.Nil(t, err)
			err = streamRcv.Send(&messages.ReceiveRequest{Command: &messages.ReceiveRequest_Start{Start: &messages.ReceiveCommandStart{
				SubId: c.subId,
			}}})
			require.Nil(t, err)
			defer streamRcv.CloseSend()
			var msgs []*pb.CloudEvent
			var msg *pb.CloudEvent
			for {
				msg, err = streamRcv.Recv()
				if err == io.EOF {
					break
				}
				if errors.Is(err, status.Error(codes.DeadlineExceeded, "context deadline exceeded")) {
					break
				}
				require.Nil(t, err, err)
				err = streamRcv.Send(&messages.ReceiveRequest{Command: &messages.ReceiveRequest_Ack{Ack: &messages.ReceiveCommandAck{
					Ack: true,
				}}})
				require.Nil(t, err, err)
				msgs = append(msgs, msg)
			}
			assert.Equal(t, len(c.msgs), len(msgs))
			var msgFound bool
			for _, msgWant := range c.msgs {
				for _, msgGot := range msgs {
					msgFound = msgWant.Id == msgGot.Id
					if !msgFound {
						continue
					}
					msgFound = msgWant.GetTextData() == msgGot.GetTextData()
					if !msgFound {
						continue
					}
					for attrKey, attrVal := range msgWant.Attributes {
						switch attrValT := attrVal.Attr.(type) {
						case *pb.CloudEventAttributeValue_CeBoolean:
							msgFound = attrValT.CeBoolean == msgGot.Attributes[attrKey].GetCeBoolean()
						case *pb.CloudEventAttributeValue_CeInteger:
							msgFound = attrValT.CeInteger == msgGot.Attributes[attrKey].GetCeInteger()
						case *pb.CloudEventAttributeValue_CeString:
							msgFound = attrValT.CeString == msgGot.Attributes[attrKey].GetCeString()
						case *pb.CloudEventAttributeValue_CeUri:
							msgFound = attrValT.CeUri == msgGot.Attributes[attrKey].GetCeUri()
						case *pb.CloudEventAttributeValue_CeUriRef:
							msgFound = attrValT.CeUriRef == msgGot.Attributes[attrKey].GetCeUriRef()
						case *pb.CloudEventAttributeValue_CeTimestamp:
							msgFound = attrValT.CeTimestamp.AsTime() == msgGot.Attributes[attrKey].GetCeTimestamp().AsTime()
						}
						// no need to check other attributes if any not equals
						if !msgFound {
							break
						}
					}
					if msgFound {
						break
					}
				}
				assert.Truef(t, msgFound, "msg id = %s was not found among the received messages: %+v", msgWant.Id, msgs)
			}
		})
	}
}
