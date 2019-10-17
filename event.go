package gaarx

import (
	"context"
	"github.com/sirupsen/logrus"
)

type (
	Event struct {
		name  string
		ctx   context.Context
		in    chan interface{}
		out   []chan interface{}
		debug bool
		log   *logrus.Logger
	}
)

func newEvent(name string, ctx context.Context) *Event {
	return &Event{
		name: name,
		ctx:  ctx,
		in:   make(chan interface{}, 100),
		out:  make([]chan interface{}, 0),
	}
}

func (e *Event) iterate() {
	for {
		select {
		case <-e.ctx.Done():
			e.Close()
			break
		case data := <-e.in:
			if e.debug {
				e.log.
					WithFields(map[string]interface{}{
						"Event Name":    e.name,
						"Received Data": data,
						"Recipients":    len(e.out),
					}).
					Debug("Received data")
			}
			if len(e.out) > 0 {
				for _, c := range e.out {
					c <- data
					if e.debug {
						e.log.WithField("Event Name", e.name).Debugf("Event was send to recipients")
					}
				}
			}
		}
	}
}

func (e *Event) Dispatch(data ...interface{}) {
	for _, d := range data {
		e.in <- d
	}
}

func (e *Event) Listen() <-chan interface{} {
	c := make(chan interface{}, 100)
	e.out = append(e.out, c)
	return c
}

func (e *Event) Close() {
	for _, c := range e.out {
		if _, ok := <-c; ok {
			close(c)
		}
	}
}
