package ep

import (
	"reflect"
	"sync"

	"github.com/thoas/go-funk"
)

// Group marges given subscribed channels into
// on subscribed channel
type Group struct {
	// Cap is capacity to create new channel
	Cap uint

	mu        sync.Mutex
	listeners []*Listener
	isInit    bool
	ID        string

	stop chan struct{}
	done chan struct{}

	cmu   sync.Mutex
	cases []*subcase

	lmu      sync.Mutex
	isListen bool
}

// subcase .
type subcase struct {
	Case        reflect.SelectCase
	emitterID     string
	listenerID  string
}

func newSubcase(ei, li string, Case reflect.SelectCase) *subcase {
	return &subcase{
		Case: Case,
		emitterID: ei,
		listenerID: li,
	}
}

// NewGroup returns new Emitter Group with a given capacity
// and perform init on it
func NewGroup(cap uint) *Group {
	res := &Group{Cap: cap, ID: randomString(15)}
	res.init()
	return res
}

// Flush reset the group to the initial state.
// All references will dropped.
func (g *Group) Flush() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.stopIfListen()
	close(g.stop)
	close(g.done)
	g.isInit = false
	g.init()
}

// Add adds channels which were already subscribed to
// some events.
func (g *Group) Add(listeners ...*Listener) {
	g.mu.Lock()
	defer g.listen()
	defer g.mu.Unlock()
	g.init()

	g.stopIfListen()

	g.cmu.Lock()
	for _, l := range listeners {
		ele := funk.Find(g.cases, func(ca *subcase) bool {
			return ca.listenerID == l.ID
		})
		if ele == nil {
			newCase := newSubcase(l.EmitterID, l.ID, reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(l.Ch),
			})
			g.cases = append(g.cases, newCase)
		}
	}
	g.cmu.Unlock()
}

// On returns subscribed channel.
func (g *Group) On() *Listener {
	g.mu.Lock()
	defer g.listen()
	defer g.mu.Unlock()
	g.init()

	g.stopIfListen()

	l := newListener(g.Cap, g.ID)
	g.listeners = append(g.listeners, l)
	return l
}

// Off unsubscribed given channels if any or unsubscribed all
// channels in other case
func (g *Group) Off(ids ...string) {
	g.mu.Lock()
	defer g.listen()
	defer g.mu.Unlock()
	g.init()

	g.stopIfListen()

	if len(ids) != 0 {
		for _, id := range ids {
			i := -1
		Listeners:
			for in := range g.listeners {
				if g.listeners[in].ID == id {
					i = in
					break Listeners
				}
			}
			if i != -1 {
				l := g.listeners[i]
				g.listeners = append(g.listeners[:i], g.listeners[i+1:]...)
				close(l.Ch)
			}
		}
		} else {
			for _, l := range g.listeners {
				close(l.Ch)
		}
		g.listeners = make([]*Listener, 0)
	}
}

// Remove unsubscribed from any subscribed channel 
func (g *Group) Remove(emitter *Emitter, ids ...string) {
	g.mu.Lock()
	defer g.listen()
	defer g.mu.Unlock()
	g.init()
	
	g.stopIfListen()
	
	if len(ids) != 0 {
		for _, id := range ids {
			if i, ca := g.findCase(id); ca != nil {
				emitter.Unsubscribe("*", id)
				g.cases = append(g.cases[:i], g.cases[i+1:]...)
			}
		}
	}
}

func (g *Group) stopIfListen() bool {
	g.lmu.Lock()
	defer g.lmu.Unlock()

	if !g.isListen {
		return false
	}

	g.stop <- struct{}{}
	g.isListen = false
	return true
}

func (g *Group) listen() {
	g.lmu.Lock()
	defer g.lmu.Unlock()
	g.cmu.Lock()
	g.isListen = true

	go func() {
		// unlock cases and isListen flag when func is exit
		defer g.cmu.Unlock()

		for {
			cases := g.listCases()
			i, val, isOpened := reflect.Select(cases)

			// exit if listening is stopped
			if i == 0 {
				return
			}

			if !isOpened && len(g.cases) > i {
				// remove this case
				g.cases = append(g.cases[:i], g.cases[i+1:]...)
			}
			var e Event
			if e2, ok := val.Interface().(Event); !ok {
				continue
			} else {
				e = e2
			}
			// use unblocked mode
			e.Flags = e.Flags | FlagSkip
			// send events to all listeners
			g.mu.Lock()
			for index := range g.listeners {
				l := g.listeners[index]
				pushEvent(g.done, l.Ch, &e)
			}
			g.mu.Unlock()
		}
	}()
}

func (g *Group) listCases() []reflect.SelectCase {
	res := []reflect.SelectCase{}
	for _, v := range g.cases {
		res = append(res, v.Case)
	}
	return res
}

func (g *Group) findCase(id string) (int, *subcase) {
	key, ele := funk.FindKey(g.cases, func(ca *subcase) bool {
		return ca.listenerID == id
	})
	return key.(int), ele.(*subcase)
}

func (g *Group) init() {
	if g.isInit {
		return
	}
	g.stop = make(chan struct{})
	g.done = make(chan struct{})
	g.cases = []*subcase{newSubcase("none", "stop", reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(g.stop),
		})}
	g.listeners = make([]*Listener, 0)
	g.isInit = true
}
