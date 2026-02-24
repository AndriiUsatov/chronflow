package nats

import (
	"context"
	"log"
	"time"

	"github.com/AndriiUsatov/chronflow/internal/pb"
	"google.golang.org/protobuf/proto"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type taskNatsQueue struct {
	jetStream jetstream.JetStream
	consumer  jetstream.Consumer
	subject   string
}

func (queue taskNatsQueue) Close() {
	queue.jetStream.Conn().Close()
}

func connect(natsURL string) (*nats.Conn, error) {
	return nats.Connect(
		natsURL,
		nats.MaxReconnects(-1),
		nats.ReconnectWait(2*time.Second),
		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) {
			log.Printf("NATS disconnected: %v", err)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Printf("NATS Reconnected to %v", nc.ConnectedUrl())
		}),
	)
}

func NewQueue(ctx context.Context, natsURL, stream, subject, consumer string, retryFor time.Duration) (taskNatsQueue, error) {

	var conn *nats.Conn
	var err error

	globalTimer := time.NewTimer(retryFor)
	defer globalTimer.Stop()

	timer := time.NewTimer(0)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return taskNatsQueue{}, ctx.Err()
		case <-globalTimer.C:
			conn, err = connect(natsURL)
			if err == nil {
				goto Connected
			}
			return taskNatsQueue{}, err
		case <-timer.C:
			conn, err = connect(natsURL)
			if err == nil {
				goto Connected
			}
			timer.Reset(2 * time.Second)
		}

	}

Connected:

	js, err := jetstream.New(conn)

	if err != nil {
		return taskNatsQueue{}, err
	}

	strm, err := js.Stream(ctx, stream)
	if err != nil {
		return taskNatsQueue{}, err
	}

	consmr, err := strm.Consumer(ctx, consumer)
	if err != nil {
		return taskNatsQueue{}, err
	}

	return taskNatsQueue{
		jetStream: js,
		consumer:  consmr,
		subject:   subject,
	}, nil

}

func (queue taskNatsQueue) PublishTask(ctx context.Context, task *pb.ProtoTask) error {
	msg, err := proto.Marshal(task)
	if err != nil {
		return err
	}

	_, err = queue.jetStream.Publish(ctx, queue.subject, msg, jetstream.WithMsgID(task.Id))

	return err
}

func (queue taskNatsQueue) ConsumeTask(ctx context.Context) (*pb.ProtoTask, error) {

	msg, err := queue.consumer.Fetch(1)
	if err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case m := <-msg.Messages():
		if m == nil {
			return nil, context.DeadlineExceeded
		}

		res := pb.ProtoTask{}

		err = proto.Unmarshal(m.Data(), &res)

		if err != nil {
			log.Printf("Broken record: ID=`%s`. Terminating record.", m.Headers().Get(nats.MsgIdHdr))
			if errr := m.Term(); errr != nil {
				return nil, err
			}
			return nil, err
		}

		if err = m.Ack(); err != nil {
			return nil, err
		}

		return &res, nil
	}

}
