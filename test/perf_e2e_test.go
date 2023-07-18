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
	var client api.Client
	client, err = api.
		NewClientBuilder().
		SubscriptionsUri(cfg.Uri.Subscriptions).
		WriterUri(cfg.Uri.Writer).
		ReaderUri(cfg.Uri.Reader).
		Build()
	require.Nil(t, err)
	defer client.Close()
	//
	groupIdCtx := metadata.AppendToOutgoingContext(context.TODO(), "x-awakari-group-id", groupId)
	//
	cases := map[string]struct {
		subCount  int
		writeRate float64
		batchSize int
		duration  time.Duration
	}{
		// subCount = 1
		"subCount = 1, writeRate = 2": {
			subCount:  1,
			writeRate: 2,
			batchSize: 1,
			duration:  300 * time.Second,
		},
		"subCount = 1, writeRate = 5": {
			subCount:  1,
			writeRate: 5,
			batchSize: 1,
			duration:  200 * time.Second,
		},
		"subCount = 1, writeRate = 10": {
			subCount:  1,
			writeRate: 10,
			batchSize: 2,
			duration:  200 * time.Second,
		},
		"subCount = 1, writeRate = 20": {
			subCount:  1,
			writeRate: 20,
			batchSize: 16,
			duration:  200 * time.Second,
		},
		"subCount = 1, writeRate = 50": {
			subCount:  1,
			writeRate: 50,
			batchSize: 16,
			duration:  200 * time.Second,
		},
		"subCount = 1, writeRate = 100": {
			subCount:  1,
			writeRate: 100,
			batchSize: 16,
			duration:  200 * time.Second,
		},
		"subCount = 1, writeRate = 200": {
			subCount:  1,
			writeRate: 200,
			batchSize: 16,
			duration:  200 * time.Second,
		},
		"subCount = 1, writeRate = 500": {
			subCount:  1,
			writeRate: 500,
			batchSize: 16,
			duration:  200 * time.Second,
		},
		// subCount = 10
		"subCount = 10, writeRate = 2": {
			subCount:  10,
			writeRate: 2,
			batchSize: 1,
			duration:  300 * time.Second,
		},
		"subCount = 10, writeRate = 5": {
			subCount:  10,
			writeRate: 5,
			batchSize: 1,
			duration:  200 * time.Second,
		},
		"subCount = 10, writeRate = 10": {
			subCount:  10,
			writeRate: 10,
			batchSize: 2,
			duration:  200 * time.Second,
		},
		"subCount = 10, writeRate = 20": {
			subCount:  10,
			writeRate: 20,
			batchSize: 16,
			duration:  200 * time.Second,
		},
		"subCount = 10, writeRate = 50": {
			subCount:  10,
			writeRate: 50,
			batchSize: 16,
			duration:  200 * time.Second,
		},
		"subCount = 10, writeRate = 100": {
			subCount:  10,
			writeRate: 100,
			batchSize: 16,
			duration:  200 * time.Second,
		},
		"subCount = 10, writeRate = 200": {
			subCount:  10,
			writeRate: 200,
			batchSize: 16,
			duration:  200 * time.Second,
		},
		"subCount = 10, writeRate = 500": {
			subCount:  10,
			writeRate: 500,
			batchSize: 16,
			duration:  200 * time.Second,
		},
		// subCount = 100
		"subCount = 100, writeRate = 2": {
			subCount:  100,
			writeRate: 2,
			batchSize: 1,
			duration:  300 * time.Second,
		},
		"subCount = 100, writeRate = 5": {
			subCount:  100,
			writeRate: 5,
			batchSize: 1,
			duration:  200 * time.Second,
		},
		"subCount = 100, writeRate = 10": {
			subCount:  100,
			writeRate: 10,
			batchSize: 2,
			duration:  200 * time.Second,
		},
		"subCount = 100, writeRate = 20": {
			subCount:  100,
			writeRate: 20,
			batchSize: 16,
			duration:  200 * time.Second,
		},
		"subCount = 100, writeRate = 50": {
			subCount:  100,
			writeRate: 50,
			batchSize: 16,
			duration:  200 * time.Second,
		},
		"subCount = 100, writeRate = 100": {
			subCount:  100,
			writeRate: 100,
			batchSize: 16,
			duration:  200 * time.Second,
		},
		"subCount = 100, writeRate = 200": {
			subCount:  100,
			writeRate: 200,
			batchSize: 16,
			duration:  200 * time.Second,
		},
		"subCount = 100, writeRate = 500": {
			subCount:  100,
			writeRate: 500,
			batchSize: 16,
			duration:  200 * time.Second,
		},
		// subCount = 1000
		"subCount = 1000, writeRate = 2": {
			subCount:  1000,
			writeRate: 2,
			batchSize: 1,
			duration:  300 * time.Second,
		},
		"subCount = 1000, writeRate = 5": {
			subCount:  1000,
			writeRate: 5,
			batchSize: 1,
			duration:  200 * time.Second,
		},
		"subCount = 1000, writeRate = 10": {
			subCount:  1000,
			writeRate: 10,
			batchSize: 2,
			duration:  200 * time.Second,
		},
		"subCount = 1000, writeRate = 20": {
			subCount:  1000,
			writeRate: 20,
			batchSize: 16,
			duration:  200 * time.Second,
		},
		"subCount = 1000, writeRate = 50": {
			subCount:  1000,
			writeRate: 50,
			batchSize: 16,
			duration:  200 * time.Second,
		},
		"subCount = 1000, writeRate = 100": {
			subCount:  1000,
			writeRate: 100,
			batchSize: 16,
			duration:  200 * time.Second,
		},
		"subCount = 1000, writeRate = 200": {
			subCount:  1000,
			writeRate: 200,
			batchSize: 16,
			duration:  200 * time.Second,
		},
		"subCount = 1000, writeRate = 500": {
			subCount:  1000,
			writeRate: 500,
			batchSize: 16,
			duration:  200 * time.Second,
		},
	}
	//
	for k, c := range cases {
		t.Run(k, func(t *testing.T) {

			t.Log("Create subscriptions...")
			var subIds []string
			for i := 0; i < c.subCount; i++ {
				subData := subscription.Data{
					Condition: condition.
						NewBuilder().
						MatchText(fmt.Sprintf("yohoho term%d", i)).
						BuildTextCondition(),
					Description: fmt.Sprintf("perf-test-sub-%d", i),
					Enabled:     true,
				}
				var subId string
				subId, err = client.CreateSubscription(groupIdCtx, userId, subData)
				require.Nil(t, err)
				subIds = append(subIds, subId)
			}
			defer func() {
				for _, subId := range subIds {
					err = client.DeleteSubscription(groupIdCtx, userId, subId)
					require.Nil(t, err)
				}
			}()

			//
			groupIdAndTimeoutCtx, cancel := context.WithTimeout(groupIdCtx, c.duration)
			defer cancel()
			writeTsByEvtId := make(map[string]int64)

			t.Log("Write start...")
			var writer model.Writer[*pb.CloudEvent]
			writer, err = client.OpenMessagesWriter(groupIdCtx, userId)
			require.Nil(t, err)
			defer writer.Close()
			var wg sync.WaitGroup
			ticker := time.NewTicker(time.Duration(float64(time.Second) * float64(c.batchSize) / (c.writeRate)))
			wg.Add(1)
			go func() {
				defer wg.Done()
				done := false
				evtCount := 0
				for !done {
					select {
					case <-ticker.C:
						var evts []*pb.CloudEvent
						for i := 0; i < c.batchSize; i++ {
							evt := &pb.CloudEvent{
								Id:         uuid.NewString(),
								Attributes: map[string]*pb.CloudEventAttributeValue{},
								Data: &pb.CloudEvent_TextData{
									// distribute events evenly across all subscriptions and always send to the last one
									TextData: fmt.Sprintf("term%d yohoho", (evtCount+i)%c.subCount),
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
				topOutput, err := exec.Command("kubectl", "top", "node").Output()
				require.Nil(t, err)
				t.Log(fmt.Sprintf("\n%s", string(topOutput)))
			}()

			t.Log("Read start...")
			readTsByEvtId := make(map[string]int64)
			ctx, cancel := context.WithTimeout(context.TODO(), c.duration)
			defer cancel()
			var mx sync.Mutex
			var reader model.Reader[[]*pb.CloudEvent]
			reader, err = client.OpenMessagesReader(ctx, "test-user-0", subIds[c.subCount-1], uint32(c.batchSize))
			require.Nil(t, err)
			var tsReadFirst int64 = 0
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
					if tsReadFirst == 0 {
						tsReadFirst = time.Now().Unix()
					}
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
			timeRead := time.Now().Unix() - tsReadFirst

			// output results
			wg.Wait()
			fmt.Printf("Written %d events, read %d\n", len(writeTsByEvtId), len(readTsByEvtId))
			var latencyData []float64
			for evtId, tsRead := range readTsByEvtId {
				tsWrite := writeTsByEvtId[evtId]
				latencyData = append(latencyData, float64(tsRead-tsWrite)/float64(time.Second))
			}
			var lat50th, lat90th, lat99th float64
			lat50th, _ = stats.Percentile(latencyData, 0.5)
			lat90th, _ = stats.Percentile(latencyData, 0.9)
			lat99th, _ = stats.Percentile(latencyData, 0.99)
			fmt.Printf("Rate: %f, Latency: 50th=%f, 90th=%f, 99th=%f [s]\n", float64(len(readTsByEvtId))/float64(timeRead), lat50th, lat90th, lat99th)

			t.Log("wait until writer queue is consumed completely...")
			time.Sleep(5 * time.Minute)

			// clean matches
			matchesCleanOutput, err := exec.Command("grpcurl", "-plaintext", "-d", "{}", "localhost:50054", "awakari.matches.Service/Clean").Output()
			require.Nil(t, err)
			t.Log(fmt.Sprintf("Matches cleaned: %s", string(matchesCleanOutput)))

			// clean messages
			messagesCleanOutput, err := exec.Command("grpcurl", "-plaintext", "-d", "{}", "localhost:50055", "awakari.messages.cleaner.Service/Clean").Output()
			require.Nil(t, err)
			t.Log(fmt.Sprintf("Messages cleaned: %s", string(messagesCleanOutput)))
		})
	}
}

func Test_Perf_MaxRate_WriteRead(t *testing.T) {
	//
	cfg, err := config.NewConfigFromEnv()
	require.Nil(t, err)
	//
	var client api.Client
	client, err = api.
		NewClientBuilder().
		SubscriptionsUri(cfg.Uri.Subscriptions).
		WriterUri(cfg.Uri.Writer).
		ReaderUri(cfg.Uri.Reader).
		Build()
	require.Nil(t, err)
	defer client.Close()
	//
	groupIdCtx := metadata.AppendToOutgoingContext(context.TODO(), "x-awakari-group-id", groupId)
	//
	cases := map[string]struct {
		subCount  int
		writeRate float64
		batchSize int
		duration  time.Duration
	}{
		// subCount = 1
		"subCount = 1, writeRate = 2": {
			subCount:  1,
			writeRate: 2,
			batchSize: 1,
			duration:  300 * time.Second,
		},
		"subCount = 1, writeRate = 5": {
			subCount:  1,
			writeRate: 5,
			batchSize: 1,
			duration:  200 * time.Second,
		},
		"subCount = 1, writeRate = 10": {
			subCount:  1,
			writeRate: 10,
			batchSize: 2,
			duration:  200 * time.Second,
		},
		"subCount = 1, writeRate = 20": {
			subCount:  1,
			writeRate: 20,
			batchSize: 16,
			duration:  200 * time.Second,
		},
		"subCount = 1, writeRate = 50": {
			subCount:  1,
			writeRate: 50,
			batchSize: 16,
			duration:  200 * time.Second,
		},
		"subCount = 1, writeRate = 100": {
			subCount:  1,
			writeRate: 100,
			batchSize: 16,
			duration:  200 * time.Second,
		},
		"subCount = 1, writeRate = 200": {
			subCount:  1,
			writeRate: 200,
			batchSize: 16,
			duration:  200 * time.Second,
		},
		"subCount = 1, writeRate = 500": {
			subCount:  1,
			writeRate: 500,
			batchSize: 16,
			duration:  200 * time.Second,
		},
		// subCount = 10
		"subCount = 10, writeRate = 2": {
			subCount:  10,
			writeRate: 2,
			batchSize: 1,
			duration:  300 * time.Second,
		},
		"subCount = 10, writeRate = 5": {
			subCount:  10,
			writeRate: 5,
			batchSize: 1,
			duration:  200 * time.Second,
		},
		"subCount = 10, writeRate = 10": {
			subCount:  10,
			writeRate: 10,
			batchSize: 2,
			duration:  200 * time.Second,
		},
		"subCount = 10, writeRate = 20": {
			subCount:  10,
			writeRate: 20,
			batchSize: 16,
			duration:  200 * time.Second,
		},
		"subCount = 10, writeRate = 50": {
			subCount:  10,
			writeRate: 50,
			batchSize: 16,
			duration:  200 * time.Second,
		},
		"subCount = 10, writeRate = 100": {
			subCount:  10,
			writeRate: 100,
			batchSize: 16,
			duration:  200 * time.Second,
		},
		"subCount = 10, writeRate = 200": {
			subCount:  10,
			writeRate: 200,
			batchSize: 16,
			duration:  200 * time.Second,
		},
		"subCount = 10, writeRate = 500": {
			subCount:  10,
			writeRate: 500,
			batchSize: 16,
			duration:  200 * time.Second,
		},
		// subCount = 100
		"subCount = 100, writeRate = 2": {
			subCount:  100,
			writeRate: 2,
			batchSize: 1,
			duration:  300 * time.Second,
		},
		"subCount = 100, writeRate = 5": {
			subCount:  100,
			writeRate: 5,
			batchSize: 1,
			duration:  200 * time.Second,
		},
		"subCount = 100, writeRate = 10": {
			subCount:  100,
			writeRate: 10,
			batchSize: 2,
			duration:  200 * time.Second,
		},
		"subCount = 100, writeRate = 20": {
			subCount:  100,
			writeRate: 20,
			batchSize: 16,
			duration:  200 * time.Second,
		},
		"subCount = 100, writeRate = 50": {
			subCount:  100,
			writeRate: 50,
			batchSize: 16,
			duration:  200 * time.Second,
		},
		"subCount = 100, writeRate = 100": {
			subCount:  100,
			writeRate: 100,
			batchSize: 16,
			duration:  200 * time.Second,
		},
		"subCount = 100, writeRate = 200": {
			subCount:  100,
			writeRate: 200,
			batchSize: 16,
			duration:  200 * time.Second,
		},
		"subCount = 100, writeRate = 500": {
			subCount:  100,
			writeRate: 500,
			batchSize: 16,
			duration:  200 * time.Second,
		},
		// subCount = 1000
		"subCount = 1000, writeRate = 2": {
			subCount:  1000,
			writeRate: 2,
			batchSize: 1,
			duration:  300 * time.Second,
		},
		"subCount = 1000, writeRate = 5": {
			subCount:  1000,
			writeRate: 5,
			batchSize: 1,
			duration:  200 * time.Second,
		},
		"subCount = 1000, writeRate = 10": {
			subCount:  1000,
			writeRate: 10,
			batchSize: 2,
			duration:  200 * time.Second,
		},
		"subCount = 1000, writeRate = 20": {
			subCount:  1000,
			writeRate: 20,
			batchSize: 16,
			duration:  200 * time.Second,
		},
		"subCount = 1000, writeRate = 50": {
			subCount:  1000,
			writeRate: 50,
			batchSize: 16,
			duration:  200 * time.Second,
		},
		"subCount = 1000, writeRate = 100": {
			subCount:  1000,
			writeRate: 100,
			batchSize: 16,
			duration:  200 * time.Second,
		},
		"subCount = 1000, writeRate = 200": {
			subCount:  1000,
			writeRate: 200,
			batchSize: 16,
			duration:  200 * time.Second,
		},
		"subCount = 1000, writeRate = 500": {
			subCount:  1000,
			writeRate: 500,
			batchSize: 16,
			duration:  200 * time.Second,
		},
	}
	//
	for k, c := range cases {
		t.Run(k, func(t *testing.T) {

			t.Log("Create subscriptions...")
			var subIds []string
			for i := 0; i < c.subCount; i++ {
				subData := subscription.Data{
					Condition: condition.
						NewBuilder().
						MatchText(fmt.Sprintf("yohoho term%d", i)).
						BuildTextCondition(),
					Description: fmt.Sprintf("perf-test-sub-%d", i),
					Enabled:     true,
				}
				var subId string
				subId, err = client.CreateSubscription(groupIdCtx, userId, subData)
				require.Nil(t, err)
				subIds = append(subIds, subId)
			}
			defer func() {
				for _, subId := range subIds {
					err = client.DeleteSubscription(groupIdCtx, userId, subId)
					require.Nil(t, err)
				}
			}()

			//
			groupIdAndTimeoutCtx, cancel := context.WithTimeout(groupIdCtx, c.duration)
			defer cancel()
			writeTsByEvtId := make(map[string]int64)

			t.Log("Write start...")
			var writer model.Writer[*pb.CloudEvent]
			writer, err = client.OpenMessagesWriter(groupIdCtx, userId)
			require.Nil(t, err)
			defer writer.Close()
			var wg sync.WaitGroup
			ticker := time.NewTicker(time.Duration(float64(time.Second) * float64(c.batchSize) / (c.writeRate)))
			wg.Add(1)
			go func() {
				defer wg.Done()
				done := false
				evtCount := 0
				for !done {
					select {
					case <-ticker.C:
						var evts []*pb.CloudEvent
						for i := 0; i < c.batchSize; i++ {
							evt := &pb.CloudEvent{
								Id:         uuid.NewString(),
								Attributes: map[string]*pb.CloudEventAttributeValue{},
								Data: &pb.CloudEvent_TextData{
									// distribute events evenly across all subscriptions and always send to the last one
									TextData: fmt.Sprintf("term%d yohoho", (evtCount+i)%c.subCount),
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
				topOutput, err := exec.Command("kubectl", "top", "node").Output()
				require.Nil(t, err)
				t.Log(fmt.Sprintf("\n%s", string(topOutput)))
			}()

			t.Log("Read start...")
			readTsByEvtId := make(map[string]int64)
			ctx, cancel := context.WithTimeout(context.TODO(), c.duration)
			defer cancel()
			var mx sync.Mutex
			var reader model.Reader[[]*pb.CloudEvent]
			reader, err = client.OpenMessagesReader(ctx, "test-user-0", subIds[c.subCount-1], uint32(c.batchSize))
			require.Nil(t, err)
			var tsReadFirst int64 = 0
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
					if tsReadFirst == 0 {
						tsReadFirst = time.Now().Unix()
					}
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
			timeRead := time.Now().Unix() - tsReadFirst

			// output results
			wg.Wait()
			fmt.Printf("Written %d events, read %d\n", len(writeTsByEvtId), len(readTsByEvtId))
			var latencyData []float64
			for evtId, tsRead := range readTsByEvtId {
				tsWrite := writeTsByEvtId[evtId]
				latencyData = append(latencyData, float64(tsRead-tsWrite)/float64(time.Second))
			}
			var lat50th, lat90th, lat99th float64
			lat50th, _ = stats.Percentile(latencyData, 0.5)
			lat90th, _ = stats.Percentile(latencyData, 0.9)
			lat99th, _ = stats.Percentile(latencyData, 0.99)
			fmt.Printf("Rate: %f, Latency: 50th=%f, 90th=%f, 99th=%f [s]\n", float64(len(readTsByEvtId))/float64(timeRead), lat50th, lat90th, lat99th)

			t.Log("wait until writer queue is consumed completely...")
			time.Sleep(5 * time.Minute)

			// clean matches
			matchesCleanOutput, err := exec.Command("grpcurl", "-plaintext", "-d", "{}", "localhost:50054", "awakari.matches.Service/Clean").Output()
			require.Nil(t, err)
			t.Log(fmt.Sprintf("Matches cleaned: %s", string(matchesCleanOutput)))

			// clean messages
			messagesCleanOutput, err := exec.Command("grpcurl", "-plaintext", "-d", "{}", "localhost:50055", "awakari.messages.cleaner.Service/Clean").Output()
			require.Nil(t, err)
			t.Log(fmt.Sprintf("Messages cleaned: %s", string(messagesCleanOutput)))
		})
	}
}
