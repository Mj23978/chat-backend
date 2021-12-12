package ep

import (
	"path"
	"sync"

	"github.com/jinzhu/copier"
	"github.com/thoas/go-funk"
)

// Flag used to describe what behavior
// do you expect.
type Flag int

const (
	// FlagReset only to clear previously defined flags.
	// Example:
	// ee.Use("*", Reset) // clears flags for this pattern
	FlagReset Flag = 0
	// FlagOnce indicates to remove the listener after first sending.
	FlagOnce Flag = 1 << iota
	// FlagVoid indicates to skip sending.
	FlagVoid
	// FlagSkip indicates to skip sending if channel is blocked.
	FlagSkip
	// FlagClose indicates to drop listener if channel is blocked.
	FlagClose
	// FlagSync indicates to send an event synchronously.
	FlagSync
)

// Middlewares.

// Reset middleware resets flags
func Reset(e *Event) { e.Flags = FlagReset }

// Once middleware sets FlagOnce flag for an event
func Once(e *Event) { e.Flags = e.Flags | FlagOnce }

// Void middleware sets FlagVoid flag for an event
func Void(e *Event) { e.Flags = e.Flags | FlagVoid }

// Skip middleware sets FlagSkip flag for an event
func Skip(e *Event) { e.Flags = e.Flags | FlagSkip }

// Close middleware sets FlagClose flag for an event
func Close(e *Event) { e.Flags = e.Flags | FlagClose }

// Sync middleware sets FlagSync flag for an event
func Sync(e *Event) { e.Flags = e.Flags | FlagSync }

// NewEmitter returns just created Emitter struct. Capacity argument
// will be used to create channels with given capacity
func NewEmitter(capacity uint) *Emitter {
	return &Emitter{
		Cap:         capacity,
		ID:          randomString(15),
		listeners:   make(map[string][]*Listener),
		middlewares: make(map[string][]func(*Event)),
		isInit:      true,
	}
}

// Emitter is a struct that allows to emit, receive
// event, close receiver channel, get info
// about topics and listeners
type Emitter struct {
	Cap         uint
	ID          string
	mu          sync.Mutex
	listeners   map[string][]*Listener
	isInit      bool
	middlewares map[string][]func(*Event)
}

func newListener(capacity uint, EmitterID string, middlewares ...func(*Event)) *Listener {
	return &Listener{
		Ch:           make(chan Event, capacity),
		ID: 					randomString(15),
		EmitterID: 		EmitterID,
		middlewares:  middlewares,
	}
}

type Listener struct {
	Ch          chan Event
	ID          string
	EmitterID   string
	middlewares []func(*Event)
}

func (e *Emitter) init() {
	if !e.isInit {
		e.listeners = make(map[string][]*Listener)
		e.middlewares = make(map[string][]func(*Event))
		e.isInit = true
	}
}

// Use registers middlewares for the pattern.
func (e *Emitter) Use(pattern string, middlewares ...func(*Event)) {
	e.mu.Lock()
	e.init()
	defer e.mu.Unlock()

	e.middlewares[pattern] = middlewares
	if len(e.middlewares[pattern]) == 0 {
		delete(e.middlewares, pattern)
	}
}

// On returns a channel that will receive events. As optional second
// argument it takes middlewares.
func (e *Emitter) On(topic string, middlewares ...func(*Event)) *Listener {
	e.mu.Lock()
	e.init()
	l := newListener(e.Cap, e.ID, middlewares...)
	if listeners, ok := e.listeners[topic]; ok {
		e.listeners[topic] = append(listeners, l)
	} else {
		e.listeners[topic] = []*Listener{l}
	}
	e.mu.Unlock()
	return l
}

// Once works exactly like On(see above) but with `Once` as the first middleware.
func (e *Emitter) Once(topic string, middlewares ...func(*Event)) *Listener {
	return e.On(topic, append(middlewares, Once)...)
}

// Off unsubscribes all listeners which were covered by
// topic, it can be pattern as well.
func (e *Emitter) Off(topic string, channels ...<-chan Event) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.init()
	match, _ := e.matched(topic)

	for _, _topic := range match {
		if listeners, ok := e.listeners[_topic]; ok {

			if len(channels) == 0 {
				for i := len(listeners) - 1; i >= 0; i-- {
					close(listeners[i].Ch)
					listeners = drop(listeners, i)
				}
			} else {
				for chi := range channels {
					curr := channels[chi]
					for i := len(listeners) - 1; i >= 0; i-- {
						if curr == listeners[i].Ch {
							close(listeners[i].Ch)
							listeners = drop(listeners, i)
						}
					}
				}
			}
			e.listeners[_topic] = listeners
		}
		if len(e.listeners[_topic]) == 0 {
			delete(e.listeners, _topic)
		}
	}
}

// Unsubscribe close listner with given id from a given
// pattern
func (e *Emitter) Unsubscribe(topic string, ids ...string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.init()
	match, _ := e.matched(topic)

	for _, _topic := range match {
		if listeners, ok := e.listeners[_topic]; ok {
			for i := len(listeners) - 1; i >= 0; i-- {
				if funk.Contains(ids, listeners[i].ID) {
					close(listeners[i].Ch)
					listeners = drop(listeners, i)
				}
			}
			e.listeners[_topic] = listeners
		}
		if len(e.listeners[_topic]) == 0 {
			delete(e.listeners, _topic)
		}
	}
}

// Listeners returns slice of listeners which were covered by
// topic(it can be pattern) and error if pattern is invalid.
func (e *Emitter) Listeners(topic string) []*Listener {
	e.mu.Lock()
	e.init()
	defer e.mu.Unlock()
	var acc []*Listener
	match, _ := e.matched(topic)

	for _, _topic := range match {
		list := e.listeners[_topic]
		for i := range e.listeners[_topic] {
			acc = append(acc, list[i])
		}
	}

	return acc
}

// Topics returns all existing topics.
func (e *Emitter) Topics() []string {
	e.mu.Lock()
	e.init()
	defer e.mu.Unlock()
	acc := make([]string, len(e.listeners))
	i := 0
	for k := range e.listeners {
		acc[i] = k
		i++
	}
	return acc
}

// Emit emits an event with the rest arguments to all
// listeners which were covered by topic(it can be pattern).
func (e *Emitter) Emit(topic string, args ...interface{}) chan struct{} {
	e.mu.Lock()
	e.init()
	done := make(chan struct{}, 1)

	match, _ := e.matched(topic)

	var wg sync.WaitGroup
	var haveToWait bool
	for _, _topic := range match {
		listeners := e.listeners[_topic]
		event := Event{
			Topic:         _topic,
			From:          e.ID,
			OriginalTopic: topic,
			Args:          args,
		}

		applyMiddlewares(&event, e.getMiddlewares(_topic))

		// whole topic is skipping
		// if (event.Flags | FlagVoid) == event.Flags {
		// 	continue
		// }

	Loop:
		for i := len(listeners) - 1; i >= 0; i-- {
			lstnr := listeners[i]
			evn := Event{}
			if err := copier.Copy(evn, event); err != nil {
				evn = event
			}
			applyMiddlewares(&evn, lstnr.middlewares)

			if (evn.Flags | FlagVoid) == evn.Flags {
				continue Loop
			}

			if (evn.Flags | FlagSync) == evn.Flags {
				_, remove, _ := pushEvent(done, lstnr.Ch, &evn)
				if remove {
					defer e.Off(event.Topic, lstnr.Ch)
				}
			} else {
				wg.Add(1)
				haveToWait = true
				go func(lstnr *Listener, event *Event) {
					e.mu.Lock()
					_, remove, _ := pushEvent(done, lstnr.Ch, event)
					if remove {
						defer e.Off(event.Topic, lstnr.Ch)
					}
					wg.Done()
					e.mu.Unlock()
				}(lstnr, &evn)
			}
		}

	}
	if haveToWait {
		go func(done chan struct{}) {
			defer func() { recover() }()
			wg.Wait()
			close(done)
		}(done)
	} else {
		close(done)
	}

	e.mu.Unlock()
	return done
}

func pushEvent(
	done chan struct{},
	lstnr chan Event,
	event *Event,
) (success, remove bool, err error) {
	// unwind the flags
	isOnce := (event.Flags | FlagOnce) == event.Flags
	isSkip := (event.Flags | FlagSkip) == event.Flags
	isClose := (event.Flags | FlagClose) == event.Flags

	sent, canceled := send(
		done,
		lstnr,
		*event,
		!(isSkip || isClose),
	)
	success = sent

	if !sent && !canceled {
		remove = isClose
		// if not sent
	} else if !canceled {
		// if event was sent successfully
		remove = isOnce
	}
	return
}

func (e *Emitter) getMiddlewares(topic string) []func(*Event) {
	var acc []func(*Event)
	for pattern, v := range e.middlewares {
		if match, _ := path.Match(pattern, topic); match {
			acc = append(acc, v...)
		} else if match, _ := path.Match(topic, pattern); match {
			acc = append(acc, v...)
		}
	}
	return acc
}

func applyMiddlewares(e *Event, fns []func(*Event)) {
	for i := range fns {
		fns[i](e)
	}
}

func (e *Emitter) matched(topic string) ([]string, error) {
	acc := []string{}
	var err error
	for k := range e.listeners {
		if matched, err := path.Match(topic, k); err != nil {
			return []string{}, err
		} else if matched {
			acc = append(acc, k)
		} else {
			if matched, _ := path.Match(k, topic); matched {
				acc = append(acc, k)
			}
		}
	}
	return acc, err
}

func drop(l []*Listener, i int) []*Listener {
	return append(l[:i], l[i+1:]...)
}

func send(
	done chan struct{},
	ch chan Event,
	e Event, wait bool,
) (sent, canceled bool) {

	defer func() {
		if r := recover(); r != nil {
			canceled = false
			sent = false
		}
	}()

	if !wait {
		select {
		case <-done:
			break
		case ch <- e:
			sent = true
			return
		default:
			return
		}

	} else {
		select {
		case <-done:
			break
		case ch <- e:
			sent = true
			return
		}

	}
	canceled = true
	return
}