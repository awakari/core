package test

import (
	"context"
	"fmt"
	"github.com/awakari/client-sdk-go/api"
	"github.com/awakari/client-sdk-go/model"
	"github.com/awakari/client-sdk-go/model/subscription"
	"github.com/awakari/client-sdk-go/model/subscription/condition"
	"github.com/awakari/core/test/config"
	"github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
	"github.com/google/uuid"
	"github.com/montanaflynn/stats"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
	"os/exec"
	"sync"
	"testing"
	"time"
)

const groupId = "perf-test-group-0"
const userId = "perf-test-user-0"

func Test_Perf_EndToEnd(t *testing.T) {
	//
	cfg, err := config.NewConfigFromEnv()
	require.Nil(t, err)
	//
	subCount := cfg.Test.Perf.E2e.SubCount
	batchSize := cfg.Test.Perf.E2e.BatchSize
	writeRate := cfg.Test.Perf.E2e.WriteRate
	duration := cfg.Test.Perf.E2e.Duration
	//
	var clients []api.Client
	for i := 0; i < subCount; i++ {
		client, err := api.
			NewClientBuilder().
			SubscriptionsUri(cfg.Uri.Subscriptions).
			WriterUri(cfg.Uri.Resolver).
			ReaderUri(cfg.Uri.Reader).
			Build()
		require.Nil(t, err)
		defer client.Close()
		clients = append(clients, client)
	}
	//
	ctxGroupId := metadata.AppendToOutgoingContext(context.TODO(), "x-awakari-group-id", groupId)
	//
	fmt.Printf("Create %d subscriptions...\n", subCount)
	var subIds []string
	for i := 0; i < subCount; i++ {
		subData := subscription.Data{
			Condition: condition.
				NewBuilder().
				AnyOfWords(fmt.Sprintf("term%d", i)).
				BuildTextCondition(),
			Description: fmt.Sprintf("perf-test-sub-%d", i),
			Enabled:     true,
		}
		var subId string
		subId, err = clients[i].CreateSubscription(ctxGroupId, userId, subData)
		require.Nil(t, err)
		subIds = append(subIds, subId)
	}
	defer func() {
		fmt.Printf("Delete %d subscriptions...\n", subCount)
		for _, subId := range subIds {
			err = clients[0].DeleteSubscription(ctxGroupId, userId, subId)
			require.Nil(t, err)
		}
	}()

	//
	groupIdAndTimeoutCtx, cancel := context.WithTimeout(ctxGroupId, duration)
	defer cancel()
	writeTsByEvtId := make(map[string]int64)

	fmt.Println("Write start...")
	var wg sync.WaitGroup
	ticker := time.NewTicker(time.Duration(float64(time.Second) * float64(batchSize) / (writeRate)))
	for n := 0; n < len(subIds); n++ {
		term := fmt.Sprintf("term%d", n)
		var writer model.Writer[*pb.CloudEvent]
		writer, err = clients[n].OpenMessagesWriter(ctxGroupId, userId)
		require.Nil(t, err)
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer writer.Close()
			done := false
			evtCount := 0
			for !done {
				select {
				case <-ticker.C:
					var evts []*pb.CloudEvent
					for i := 0; i < batchSize; i++ {
						evt := &pb.CloudEvent{
							Id:         uuid.NewString(),
							Attributes: map[string]*pb.CloudEventAttributeValue{},
							Data: &pb.CloudEvent_TextData{
								TextData: term,
							},
						}
						evts = append(evts, evt)
					}
					var ackCount uint32
					ackCount, err = writer.WriteBatch(evts)
					require.Nil(t, err)
					evtCount += int(ackCount)
					for i := 0; i < int(ackCount); i++ {
						writeTsByEvtId[evts[i].Id] = time.Now().UnixNano()
					}
				case <-groupIdAndTimeoutCtx.Done():
					done = true
					break
				}
			}
		}()
	}

	fmt.Printf("Read start...")
	readTsByEvtId := make(map[string]int64)
	ctx, cancel := context.WithTimeout(context.TODO(), duration)
	defer cancel()
	var mx sync.Mutex
	for n, subId := range subIds {
		client := clients[n]
		subIdCopy := subId
		wg.Add(1)
		go func() {
			defer wg.Done()
			var reader model.Reader[[]*pb.CloudEvent]
			reader, err = client.OpenMessagesReader(ctx, "test-user-0", subIdCopy, uint32(batchSize))
			require.Nil(t, err)
			var evts []*pb.CloudEvent
			done := false
			for !done {
				evts, err = reader.Read()
				if err == context.DeadlineExceeded {
					break
				}
				require.Nil(t, err)
				mx.Lock()
				for _, evt := range evts {
					readTsByEvtId[evt.Id] = time.Now().UnixNano()
				}
				mx.Unlock()
				select {
				case <-ctx.Done():
					done = true
					break
				default:
					continue
				}
			}
		}()
	}

	// output results
	wg.Wait()
	fmt.Printf("Written %d events, read %d\n", len(writeTsByEvtId), len(readTsByEvtId))
	var latencyData []float64
	var tsReadMin int64 = 0
	var tsReadMax int64 = 0
	for evtId, tsRead := range readTsByEvtId {
		tsWrite := writeTsByEvtId[evtId]
		latencyData = append(latencyData, float64(tsRead-tsWrite)/float64(time.Second))
		if tsReadMin == 0 || tsRead < tsReadMin {
			tsReadMin = tsRead
		}
		if tsRead > tsReadMax {
			tsReadMax = tsRead
		}
	}
	var lat50th, lat90th, lat99th float64
	lat50th, _ = stats.Percentile(latencyData, 0.5)
	lat90th, _ = stats.Percentile(latencyData, 0.9)
	lat99th, _ = stats.Percentile(latencyData, 0.99)
	fmt.Printf("Rate: %f, Latency: 50th=%f, 90th=%f, 99th=%f [s]\n", float64(len(readTsByEvtId))*float64(time.Second)/float64(tsReadMax-tsReadMin), lat50th, lat90th, lat99th)

	fmt.Println("wait until writer queue is consumed completely...")
	time.Sleep(5 * time.Minute)

	// clean matches
	matchesCleanOutput, err := exec.Command("grpcurl", "-plaintext", "-d", "{}", cfg.Uri.Matches, "awakari.matches.Service/Clean").Output()
	require.Nil(t, err)
	fmt.Printf("Matches cleaned: %s\n", string(matchesCleanOutput))

	// clean messages
	messagesCleanOutput, err := exec.Command("grpcurl", "-plaintext", "-d", "{}", cfg.Uri.Messages, "awakari.messages.cleaner.Service/Clean").Output()
	require.Nil(t, err)
	fmt.Printf("Messages cleaned: %s\n", string(messagesCleanOutput))
}
