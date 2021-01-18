package tracestan

import (
	"bytes"
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/nats-io/stan.go"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"google.golang.org/protobuf/proto"
)

const component = "STAN"

type Client struct {
	conn   stan.Conn
	tracer opentracing.Tracer

	subs map[string][]stan.Subscription
	mu   sync.Mutex
}

// Handler can be: func(ctx context.Context, data proto.Message) error.
type Handler interface{}

func NewClient(config *Config, tracer opentracing.Tracer, opts ...stan.Option) (*Client, error) {
	opts = append(opts, config.addr())

	conn, err := stan.Connect(config.ClusterID, config.ClientID, opts...)
	if err != nil {
		return nil, fmt.Errorf("stan connect: %w", err)
	}

	client := &Client{
		conn:   conn,
		tracer: tracer,
		subs:   make(map[string][]stan.Subscription),
	}

	return client, err
}

func (c *Client) Subscribe(topic string, handler Handler, opts ...stan.SubscriptionOption) error {
	cb, err := c.buildHandler(topic, handler)
	if err != nil {
		return fmt.Errorf("build handler: %w", err)
	}

	sub, err := c.conn.Subscribe(topic, cb, opts...)
	if err != nil {
		return fmt.Errorf("subscribe: %w", err)
	}

	c.mu.Lock()
	c.subs[topic] = append(c.subs[topic], sub)
	c.mu.Unlock()

	return nil
}

func (c *Client) QueueSubscribe(topic, queue string, handler Handler, opts ...stan.SubscriptionOption) error {
	cb, err := c.buildHandler(topic, handler)
	if err != nil {
		return fmt.Errorf("build handler: %w", err)
	}

	sub, err := c.conn.QueueSubscribe(topic, queue, cb, opts...)
	if err != nil {
		return fmt.Errorf("subscribe: %w", err)
	}

	c.mu.Lock()
	c.subs[topic] = append(c.subs[topic], sub)
	c.mu.Unlock()

	return nil
}

func (c *Client) Publish(ctx context.Context, topic string, v interface{}) error {
	var parentCtx opentracing.SpanContext

	if parent := opentracing.SpanFromContext(ctx); parent != nil {
		parentCtx = parent.Context()
	}

	span := c.tracer.StartSpan(defaultNameFunc(topic), opentracing.ChildOf(parentCtx))
	defer span.Finish()

	ext.SpanKindProducer.Set(span)
	ext.MessageBusDestination.Set(span, topic)
	ext.Component.Set(span, component)

	buf := bytes.NewBuffer(nil)

	// We have no better place to record an error than the Span itself.
	if err := c.tracer.Inject(span.Context(), opentracing.Binary, buf); err != nil {
		span.LogFields(log.String("event", "Tracer.Inject() failed"), log.Error(err))
	}

	// Encode proto message.
	data, err := c.encode(v)
	if err != nil {
		return fmt.Errorf("encode proto message: %w", err)
	}

	// Write payload.
	if _, err := buf.Write(data); err != nil {
		return fmt.Errorf("write payload: %w", err)
	}

	return c.conn.Publish(topic, buf.Bytes())
}

func (c *Client) Unsubscribe(topic string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	buf := make([]error, 0)

	for _, sub := range c.subs[topic] {
		if err := sub.Unsubscribe(); err != nil {
			buf = append(buf, fmt.Errorf("unsubscribe %s: %w", topic, err))
		}

		delete(c.subs, topic)
	}

	if len(buf) > 0 {
		return fmt.Errorf("%w: %+v", ErrUnsubscribe, buf)
	}

	return nil
}

func (c *Client) UnsubscribeAll() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	buf := make([]error, 0)

	for topic, subs := range c.subs {
		for _, sub := range subs {
			if err := sub.Unsubscribe(); err != nil {
				buf = append(buf, fmt.Errorf("unsubscribe %s: %w", topic, err))
			}
		}

		delete(c.subs, topic)
	}

	if len(buf) > 0 {
		return fmt.Errorf("%w: %+v", ErrUnsubscribe, buf)
	}

	return nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) buildHandler(topic string, handler Handler) (stan.MsgHandler, error) {
	callback, argType, err := reflectCb(handler)
	if err != nil {
		return nil, err
	}

	return c.handler(topic, callback, argType), nil
}

func (c *Client) handler(topic string, callback reflect.Value, argType reflect.Type) stan.MsgHandler {
	return func(msg *stan.Msg) {
		// Read trace carrier.
		buf := bytes.NewBuffer(msg.Data)

		spanContext, _ := c.tracer.Extract(opentracing.Binary, buf)

		span := c.tracer.StartSpan(defaultNameFunc(topic), opentracing.FollowsFrom(spanContext))
		defer span.Finish()

		ext.SpanKindConsumer.Set(span)
		ext.MessageBusDestination.Set(span, msg.Subject)
		ext.Component.Set(span, component)

		ctx := opentracing.ContextWithSpan(context.Background(), span)

		// Decode bytes to proto message.
		var oPtr reflect.Value

		if argType.Kind() != reflect.Ptr {
			oPtr = reflect.New(argType)
		} else {
			oPtr = reflect.New(argType.Elem())
		}

		if err := c.decode(buf.Bytes(), oPtr.Interface()); err != nil {
			ext.Error.Set(span, true)
			ext.LogError(span, err)

			return
		}

		if argType.Kind() != reflect.Ptr {
			oPtr = reflect.Indirect(oPtr)
		}

		// Run handler with args.
		result := callback.Call([]reflect.Value{reflect.ValueOf(ctx), oPtr})
		if errVal := result[0]; !errVal.IsNil() {
			ext.Error.Set(span, true)
			ext.LogError(span, errVal.Interface().(error))

			return
		}

		if err := msg.Ack(); err != nil {
			ext.Error.Set(span, true)
			ext.LogError(span, err)
		}
	}
}

func (c *Client) encode(v interface{}) ([]byte, error) {
	if v == nil {
		return nil, ErrEmptyProtoMsg
	}

	msg, ok := v.(proto.Message)
	if !ok {
		return nil, ErrInvalidProtoMsgEncode
	}

	return proto.Marshal(msg)
}

func (c *Client) decode(data []byte, v interface{}) error {
	msg, ok := v.(proto.Message)
	if !ok {
		return ErrInvalidProtoMsgDecode
	}

	return proto.Unmarshal(data, msg)
}

func defaultNameFunc(method string) string {
	return "STAN " + method
}

func reflectCb(cb Handler) (callback reflect.Value, argType reflect.Type, err error) {
	const cbArgs = 2

	cbType := reflect.TypeOf(cb)
	if cbType.Kind() != reflect.Func {
		return reflect.Value{}, nil, fmt.Errorf("%w: handler has to be a func", ErrInvalidArgs)
	}

	numArgs := cbType.NumIn()
	if numArgs != cbArgs {
		return reflect.Value{}, nil, fmt.Errorf("%w: num of arguments should be 2", ErrInvalidArgs)
	}

	if !cbType.In(0).Implements(reflect.TypeOf((*context.Context)(nil)).Elem()) {
		return reflect.Value{}, nil, fmt.Errorf("%w: first arg must be context", ErrInvalidArgs)
	}

	if !cbType.In(1).Implements(reflect.TypeOf((*proto.Message)(nil)).Elem()) {
		return reflect.Value{}, nil, fmt.Errorf("%w: second arg must implement proto.Message", ErrInvalidArgs)
	}

	if cbType.NumOut() != 1 {
		return reflect.Value{}, nil, fmt.Errorf("%w: num of return args should be 1", ErrInvalidResult)
	}

	if !cbType.Out(0).Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		return reflect.Value{}, nil, fmt.Errorf("%w: return type should be 'error'", ErrInvalidResult)
	}

	return reflect.ValueOf(cb), cbType.In(numArgs - 1), nil
}
