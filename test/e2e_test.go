package test

import (
	"context"
	"github.com/awakari/client-sdk-go/api"
	"github.com/awakari/client-sdk-go/model"
	"github.com/awakari/core/test/config"
	"github.com/awakari/core/test/data"
	"github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
	"testing"
	"time"
)

func Test_MessageDelivery(t *testing.T) {
	//
	cfg, err := config.NewConfigFromEnv()
	require.Nil(t, err)
	//
	var client api.Client
	client, err = api.
		NewClientBuilder().
		SubscriptionsUri(cfg.Uri.Subscriptions).
		WriterUri(cfg.Uri.Resolver).
		ReaderUri(cfg.Uri.Reader).
		Build()
	require.Nil(t, err)
	defer client.Close()
	//
	groupIdCtx := metadata.AppendToOutgoingContext(context.TODO(), "x-awakari-group-id", "test-group-0")
	//
	var subIds []string
	for _, subData := range data.Subs {
		var subId string
		subId, err = client.CreateSubscription(groupIdCtx, "test-user-0", subData)
		require.Nil(t, err)
		subIds = append(subIds, subId)
	}
	//
	//time.Sleep(100 * time.Second) // wait for the cond/sub cache entries expiration
	//
	var msgsWriter model.Writer[*pb.CloudEvent]
	msgsWriter, err = client.OpenMessagesWriter(groupIdCtx, "test-user-1")
	require.Nil(t, err)
	defer msgsWriter.Close()
	var msgCount uint32
	msgCount, err = msgsWriter.WriteBatch(data.Msgs)
	require.Equal(t, len(data.Msgs), int(msgCount))
	require.Nil(t, err)
	//
	cases := map[string]struct {
		subId string
		msgs  []*pb.CloudEvent
		err   error
	}{
		"disabled": {
			subId: subIds[0],
			msgs:  []*pb.CloudEvent{},
			err:   context.DeadlineExceeded,
		},
		"exact complete value match for a key": {
			subId: subIds[1],
			msgs: []*pb.CloudEvent{
				data.Msgs[4],
			},
		},
		"partial exact match": {
			subId: subIds[2],
			msgs: []*pb.CloudEvent{
				data.Msgs[0],
			},
		},
		"basic group condition with \"and\" logic and partial sub-conditions": {
			subId: subIds[3],
			msgs: []*pb.CloudEvent{
				data.Msgs[1],
			},
		},
		"basic group condition with \"or\" logic": {
			subId: subIds[4],
			msgs: []*pb.CloudEvent{
				data.Msgs[2],
				data.Msgs[3],
			},
		},
		"basic group condition with \"and\" logic and a negative sub-condition": {
			subId: subIds[5],
			msgs: []*pb.CloudEvent{
				data.Msgs[0],
				data.Msgs[1],
			},
		},
	}
	//
	for k, c := range cases {
		t.Run(k, func(t *testing.T) {
			//
			ctx, cancel := context.WithTimeout(context.TODO(), 50*time.Second)
			defer cancel()
			var msgsReader model.Reader[[]*pb.CloudEvent]
			msgsReader, err = client.OpenMessagesReader(ctx, "test-user-0", c.subId, 4)
			require.Nil(t, err)
			defer msgsReader.Close()
			//
			var msgs []*pb.CloudEvent
			var msgBatch []*pb.CloudEvent
			for {
				msgBatch, err = msgsReader.Read()
				assert.True(t, err == nil || err == context.DeadlineExceeded)
				msgs = append(msgs, msgBatch...)
				select {
				case <-ctx.Done():
					err = ctx.Err()
				default:
					continue
				}
				if err != nil {
					break
				}
			}
			assert.Equal(t, len(c.msgs), len(msgs), c.subId)
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
