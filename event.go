package gaarx

import (
	"context"
	"github.com/rs/zerolog/log"
)

type (
	Event struct {
		name  string
		ctx   context.Context
		in    chan interface{}
		out   []chan interface{}
		debug bool
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
				log.Debug().
					Str("Event Name", e.name).
					Interface("Received Data", data).
					Int("Recipients", len(e.out)).
					Msg("Received data")
			}
			if len(e.out) > 0 {
				for _, c := range e.out {
					c <- data
					if e.debug {
						log.Debug().Str("Event Name", e.name).Msgf("Event was send to recipients")
					}
				}
			}
		}
	}
}

// Dispatch notice all listeners with data
func (e *Event) Dispatch(data ...interface{}) {
	for _, d := range data {
		e.in <- d
	}
}

// Listen return channel to receive data from dispatch
func (e *Event) Listen() <-chan interface{} {
	c := make(chan interface{}, 100)
	e.out = append(e.out, c)
	return c
}

// Close finalize event
func (e *Event) Close() {
	for _, c := range e.out {
		if _, ok := <-c; ok {
			close(c)
		}
	}
}
