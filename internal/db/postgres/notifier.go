package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/lib/pq"
)

type pgNotificationListener struct {
	listener *pq.Listener
	event    string
}

func handleError(ev pq.ListenerEventType, err error) {
	if err != nil {
		log.Println(err)
	}
}

func (listener pgNotificationListener) Close() {
	listener.listener.Close()
}

func (notifier pgNotificationListener) goListen(ctx context.Context, listener *pq.Listener, channel chan<- struct{}) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-listener.Notify:
			select {
			case channel <- struct{}{}:
			default:
			}
		}

	}
}

func (notifier pgNotificationListener) Listen(ctx context.Context) (<-chan struct{}, error) {
	res := make(chan struct{}, 1)

	go notifier.goListen(ctx, notifier.listener, res)

	return res, nil
}

func initialListen(ctx context.Context, listener *pq.Listener, event string, retryFor time.Duration) error {
	globalTimer := time.NewTimer(retryFor)
	defer globalTimer.Stop()

	timer := time.NewTimer(0)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-globalTimer.C:
			return listener.Listen(event)
		case <-timer.C:
			err := listener.Listen(event)
			if err == nil {
				return nil
			}
			timer.Reset(5 * time.Second)
		}
	}
}

func NewPGListener(ctx context.Context, url string, port int, username string, password string, database string, sslMode string, event string, retryFor time.Duration) (pgNotificationListener, error) {
	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", url, port, username, password, database, sslMode)
	listener := pq.NewListener(conn, 10*time.Second, time.Minute, handleError)

	if err := initialListen(ctx, listener, event, retryFor); err != nil {
		return pgNotificationListener{}, err
	}

	return pgNotificationListener{
		listener: listener,
		event:    event,
	}, nil
}
