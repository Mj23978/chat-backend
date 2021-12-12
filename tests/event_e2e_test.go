package ep

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/thoas/go-funk"
)

var (
	wg   *sync.WaitGroup
	tn   *testNode
	ta   *testAPI
	ts   *testStore
	tb   *testBroker
	done chan bool
)

type testNode struct {
	emitters map[string]*Emitter
	groups   map[string]*Group
}

type testAPI struct {
	emitter *Emitter
}

type testStore struct {
	emitter *Emitter
}

type testBroker struct {
	emitter *Emitter
}

func init() {
	wg = new(sync.WaitGroup)
	done = make(chan bool, 1)
	tn = &testNode{
		emitters: make(map[string]*Emitter),
		groups:   make(map[string]*Group),
	}
	ta = &testAPI{
		emitter: tn.newEmitter("api"),
	}
	tb = &testBroker{
		emitter: tn.newEmitter("broker"),
	}
	ts = &testStore{
		emitter: tn.newEmitter("store"),
	}
	tn.groups["req:*"] = tn.getPatternEvents("req:*")
}

func (t *testNode) newEmitter(name string) *Emitter {
	em := NewEmitter(2)
	// em.Use("*", Sync) // In Default Emit Uses a goroutine to send messages but with sync it dont
	t.emitters[name] = em
	return em
}

func (t *testNode) getPatternEvents(pattern string, middlewares ...func(*Event)) *Group {
	group := NewGroup(30)
	r := funk.Map(t.emitters, func(_ string, em *Emitter) *Listener {
		return em.On(pattern, middlewares...)
	}).([]*Listener)
	group.Add(r...)
	fmt.Printf("All Channels : %v\n", len(r))
	return group
}

func TestServiceNode(t *testing.T) {
	go func() {
		for {
			select {
			case res := <-tn.groups["req:*"].On().Ch:
				fmt.Printf("Event %v\n", res)
			}
		}
	}()
	<-tn.emitters["api"].Emit("req:store:get", "Something to Get from Store")
	<-tn.emitters["store"].Emit("req:api:get", "Test for Buffer")
	<-tn.emitters["api"].Emit("req:broker:PubUserJoin", "Publish user joined in a broker for listeners")
	time.Sleep(time.Second * 3)
}
