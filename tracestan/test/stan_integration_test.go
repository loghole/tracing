// +build integration

package test

import (
	"context"
	"testing"
	"time"

	"github.com/nats-io/stan.go"
	"github.com/stretchr/testify/assert"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	"github.com/loghole/tracing"
	"github.com/loghole/tracing/tracestan"
)

func TestIntegration(t *testing.T) {
	tracer, err := tracing.NewTracer(tracing.DefaultConfiguration("stan-subscribe", "127.0.0.1:6831"))
	if err != nil {
		t.Error(err)

		return
	}

	defer tracer.Close()

	client, err := tracestan.NewClient(&tracestan.Config{
		Addr:      "127.0.0.1:4222",
		ClientID:  "test_subscribe",
		ClusterID: "test-cluster",
	}, tracer)
	if err != nil {
		t.Error(err)

		return
	}

	defer client.Close()

	t.Run("PublishSubscribeTest", PublishSubscribeTest(client))
	t.Run("PublishQueueSubscribeTest", PublishQueueSubscribeTest(client))
}

func PublishSubscribeTest(client *tracestan.Client) func(t *testing.T) {
	return func(t *testing.T) {
		var (
			received = make([]string, 0, 0)
			expected = []string{"1", "2", "3", "4", "5"}
			topic    = time.Now().Format("200601021504059999999990700")
		)

		err := client.Subscribe(topic, func(ctx context.Context, data *errdetails.DebugInfo) error {
			received = append(received, data.Detail)

			return nil
		}, stan.SetManualAckMode(), stan.StartWithLastReceived())
		if err != nil {
			t.Error(err)

			return
		}

		for _, val := range expected {
			err = client.Publish(context.TODO(), topic, &errdetails.DebugInfo{Detail: val})
			if err != nil {
				t.Error(err)

				return
			}
		}

		time.Sleep(time.Second)

		assert.Equal(t, expected, received)

		if err := client.Unsubscribe(topic); err != nil {
			t.Error(err)
		}
	}
}

func PublishQueueSubscribeTest(client *tracestan.Client) func(t *testing.T) {
	return func(t *testing.T) {
		var (
			received = make([]string, 0, 0)
			expected = []string{"1", "2", "3", "4", "5"}
			topic    = time.Now().Format("200601021504059999999990700")
		)

		err := client.QueueSubscribe(topic, time.Now().String(), func(ctx context.Context, data *errdetails.DebugInfo) error {
			received = append(received, data.Detail)

			return nil
		}, stan.SetManualAckMode(), stan.StartWithLastReceived())
		if err != nil {
			t.Error(err)

			return
		}

		for _, val := range expected {
			err = client.Publish(context.TODO(), topic, &errdetails.DebugInfo{Detail: val})
			if err != nil {
				t.Error(err)

				return
			}
		}

		time.Sleep(time.Second)

		assert.Equal(t, expected, received)

		if err := client.UnsubscribeAll(); err != nil {
			t.Error(err)
		}
	}
}
