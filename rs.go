package main

import (
	"errors"
	"fmt"
	"reflect"
	//  	"log"
)

// Command handlers are responsible
// for validating commands,
// both as a stand-alone set of data
// as well as in the context
// of the Command's aggregate (I know,
// lots of undefined terms here ...
// see the github wiki).
type CommandHandler interface {
	handle(c Command) (e []Event, err error)
}

type EventListener interface {
	apply(e Event) error
}

// Maps a command to its handler.
type HandlerRegistry map[reflect.Type]CommandHandler

type ListenerRegistry map[reflect.Type][]EventListener

// Registers event and command listeners.  Dispatches commands.
type messageDispatcher struct {
	registry HandlerRegistry
	listeners ListenerRegistry
}

func (md *messageDispatcher) SendCommand(c Command) ([]Event, error) {
	t := reflect.TypeOf(c)
	if h, ok := md.registry[t]; ok {
		return h.handle(c)
	}
	return nil, errors.New(fmt.Sprint("No handler registered for command ", t))
}

// Since events represent a thing that actually happened,
// a fact, having an event listener return an error
// is probably not the right thing to do.
// While errors can certainly occor,
// for example, email server or database is down
// or the file system is full,
// a better approach would be to stick the events
// in a durable queue so when the error condition clears
// the listener can successfully do it's thing.
func (md *messageDispatcher) PublishEvent(e Event) error {
	var err error;
	t := reflect.TypeOf(e)
	if a, ok := md.listeners[t]; ok {
		for _, listener := range a {
			err = listener.apply(e)
			if err != nil {
				return err;
			}
		}
	}
	return nil
}

func NewMessageDispatcher(hr HandlerRegistry, lr ListenerRegistry) (*messageDispatcher, error) {
	m := make(HandlerRegistry, len(hr))
	for commandtype, handler := range hr {
		m[commandtype] = handler
	}
	return &messageDispatcher{registry: m}, nil
}

func main() {}
