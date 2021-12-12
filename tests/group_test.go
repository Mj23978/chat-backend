package ep

import (
	"fmt"
	"testing"
	"time"
)

func TestGroupInternals(t *testing.T) {
	g := &Group{}
	g.init()
	expect(t, g.isInit, true)
	g.stopIfListen()
	expect(t, len(g.stop), 0)
	g.listen()
	expect(t, g.isListen, true)
	g.stopIfListen()
}

func TestGroupBasic(t *testing.T) {
	g := &Group{Cap: 5}

	e := NewEmitter(5)
	e2 := NewEmitter(5)
	e3 := NewEmitter(5)

	e.Use("*", Sync)
	e2.Use("*", Sync)
	e3.Use("*", Sync)

	g.Add(
		e.On("*"),
		e2.On("*"),
		e3.On("*"),
	)

	pip := g.On()
	pipe := pip.Ch

	expect(t, len(pipe), 0)

	<-e.Emit("*", 1)
	<-e.Emit("*", 2)
	<-e2.Emit("*", 3)
	<-e3.Emit("*", 4)
	<-e3.Emit("*", 5)

	// departure/arrival order
	expect(t, (<-pipe).Int(0), 1)
	expect(t, (<-pipe).Int(0), 2)
	expect(t, (<-pipe).Int(0), 3)
	expect(t, (<-pipe).Int(0), 4)
	expect(t, (<-pipe).Int(0), 5)

	g.Off(pip.ID)

	_, ok := <-pipe
	expect(t, ok, false)

}

func TestGroupFlushOnOff(t *testing.T) {
	g := NewGroup(10)
	g.On()
	expect(t, len(g.listeners), 1)
	g.Flush()
	expect(t, len(g.listeners), 0)
	g.On()
	g.On()
	expect(t, len(g.listeners), 2)
	g.Off()
	expect(t, len(g.listeners), 0)
}

func TestGrouEmitterOff(t *testing.T) {
	g := NewGroup(10)
	expect(t, len(g.listCases()), 1)
	em1 := NewEmitter(3)
	el1 := em1.On("*")
	g.Add(el1)
	expect(t, len(g.listCases()), 2)
	em2 := NewEmitter(3)
	el2 := em2.On("*")
	g.Add(el2)
	expect(t, len(g.listCases()), 3)
	gl := g.On()
	go func() {
		for {
			select {
			case t := <-gl.Ch:
				fmt.Printf("Get Message %v from %v\n", t.OriginalTopic, t.From)
			}
		}
	}()
	em2.Emit("First")
	em1.Emit("Second")
	em2.Emit("Third")
	time.Sleep(time.Second)
}

func TestGrouEmitterRemove(t *testing.T) {
	g := NewGroup(2)
	em1 := NewEmitter(1)
	el1 := em1.On("*")
	g.Add(el1)
	expect(t, len(g.listCases()), 2)
	expect(t, len(em1.Listeners("*")), 1)
	expect(t, em1.Listeners("*")[0].ID, el1.ID)
	g.Remove(em1, el1.ID)
	expect(t, len(em1.Listeners("*")), 0)
	expect(t, len(g.listCases()), 1)
	time.Sleep(time.Millisecond * 1500)
}

func TestGrouCaseOrder(t *testing.T) {
	g := NewGroup(10)
	em := NewEmitter(1)
	els := []*Listener{}
	for i := 0; i < 5; i++ {
		el := em.On(fmt.Sprintf("%v:*", i))
		els = append(els, el)
		g.Add(el)
	}
	expect(t, len(g.listCases()), 6)
	expect(t, len(els), 5)
	cn := g.On().Ch
	go func() {
		for {
			select {
			case t := <-cn:
				fmt.Printf("Get Message %v\n", t)
			}
		}
	}()
	em.Emit("1:9", 1)
	em.Emit("*", 2)
	time.Sleep(time.Second)
}

func TestGroupOrder(t *testing.T) {
	g := &Group{Cap: 5}

	e := new(Emitter)
	e2 := new(Emitter)
	e3 := new(Emitter)

	e.Use("*", Sync)
	e2.Use("*", Sync)
	e3.Use("*", Sync)

	g.Add(
		e.On("*"),
		e2.On("*"),
		e3.On("*"),
	)

	pip := g.On()
	pipe := pip.Ch

	expect(t, len(pipe), 0)

	<-e.Emit("*", 1)
	<-e.Emit("*", 2)
	<-e2.Emit("*", 3)
	<-e3.Emit("*", 4)
	<-e3.Emit("*", 5)

	// departure/arrival order
	expect(t, (<-pipe).Int(0), 1)
	expect(t, (<-pipe).Int(0), 2)
	expect(t, (<-pipe).Int(0), 3)
	expect(t, (<-pipe).Int(0), 4)
	expect(t, (<-pipe).Int(0), 5)

	g.Off(pip.ID)

	_, ok := <-pipe
	expect(t, ok, false)

}