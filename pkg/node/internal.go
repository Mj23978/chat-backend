package gameNode

import (
	// "fmt"

	log "github.com/mj23978/chat-backend-x/logger/zerolog"
	test "github.com/mj23978/chat-backend/pkg"
	ep "github.com/mj23978/chat-backend/tests"
	"github.com/thoas/go-funk"
)

func (gn *gameNode) handleEmitters() {
	gp := gn.getPatternEvents("req:*")
	log.Debugf("handlingEmitters Started")

	go func() {
		for {
			msg := <-gp.On().Ch
			topic := msg.OriginalTopic

			switch topic {
			case test.EmitterRequestGetUser:
				log.Infof("Get User Arg - %v", msg.Args[0])
			}
		}
	}()
}

func (gn *gameNode) getPatternEvents(pattern string, middlewares ...func(*ep.Event)) *ep.Group {
	group := ep.NewGroup(30)
	r := funk.Map(gn.emitters, func(_ string, em *ep.Emitter) *ep.Listener {
		return em.On(pattern, middlewares...)
	}).([]*ep.Listener)
	group.Add(r...)
	log.Debugf("All Channels : %v\n", len(r))
	gn.groups[pattern] = group
	return group
}

func (gn *gameNode) newEmitter(name string) *ep.Emitter {
	em := ep.NewEmitter(2)
	// em.Use("*", Sync) // In Default Emit Uses a goroutine to send messages but with sync it dont
	gn.emitters[name] = em
	return em
}
