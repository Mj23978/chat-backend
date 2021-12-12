package ep

import (
	"reflect"
	"testing"
	"time"
)

func TestFlatBasic(t *testing.T) {
	ee := &Emitter{}
	go ee.Emit("test", nil)
	event := <-ee.On("test").Ch
	expect(t, len(event.Args), 1)
}

func TestFlatClose(t *testing.T) {
	ee := NewEmitter(0)
	wait := make(chan struct{})
	pipe := ee.On("test").Ch
	ee.On("test", Close)
	l := ee.Listeners("test")
	expect(t, len(l), 2)
	go func() {
		event := <-pipe
		expect(t, len(event.Args), 1)
		wait <- struct{}{}
	}()
	<-ee.Emit("test", "close")
	<-wait

	go func() {
		for range pipe {
		}
		wait <- struct{}{}
	}()
	l = ee.Listeners("test")

	expect(t, len(l), 1)
	ee.Off("test", pipe)

	<-wait
	expect(t, len(ee.Topics()), 0)
}

func TestBufferedBasic(t *testing.T) {
	ee := NewEmitter(1)
	// ee.Use("*", OrSkip)
	ch := make(chan struct{})
	pipe := ee.On("test").Ch
	go func() {
		event := <-pipe
		expect(t, len(event.Args), 2)
		ch <- struct{}{}
	}()
	<-ee.Emit("test", nil, true)
	<-ch
}

func TestOff(t *testing.T) {
	ee := NewEmitter(0)
	ee.On("test")
	ee.On("test")
	expect(t, len(ee.Topics()), 1)
	l := ee.Listeners("test")
	expect(t, len(l), 2)

	ee.Off("test")
	l = ee.Listeners("test")
	expect(t, len(l), 0)
	expect(t, len(ee.Topics()), 0)
}

func TestRange(t *testing.T) {
	ee := NewEmitter(0)
	c := 42
	go ee.Emit("test", "range", "it", c)
	for event := range ee.On("test", Close).Ch { // Close if channel is blocked
		expect(t, event.String(0), "range")
		expect(t, event.String(1), "it")
		expect(t, event.Int(2), c)
		// ee.Off("test")
		break
	}
	l := ee.Listeners("test")
	expect(t, len(l), 1)
	<-ee.Emit("test", "range", "it", 42)
	l = ee.Listeners("test")
	expect(t, len(l), 0)
}

func TestCloseOnBlock(t *testing.T) {
	ee := NewEmitter(0)

	ee.On("test0", Close)
	l := ee.Listeners("test0")
	expect(t, len(l), 1)
	expect(t, len(ee.Topics()), 1)
	<-ee.Emit("test0")
	l = ee.Listeners("test0")
	expect(t, len(l), 0)
	expect(t, len(ee.Topics()), 0)

	ee = NewEmitter(3)
	ee.Use("test*", Close)
	ee.On("test1")
	ee.On("test2")

	<-ee.Emit("test1")
	<-ee.Emit("test1")
	<-ee.Emit("test1")
	l = ee.Listeners("test1")
	expect(t, len(l), 1)
	expect(t, len(ee.Topics()), 2)
	<-ee.Emit("test1") // should raise blockedError
	// ^^^^ and remove the topic as well
	l = ee.Listeners("test1")
	expect(t, len(l), 0)
	expect(t, len(ee.Topics()), 1)
	<-ee.Emit("test2")
	<-ee.Emit("test2")
	<-ee.Emit("test2")
	expect(t, len(ee.Topics()), 1)
	l = ee.Listeners("test2")
	expect(t, len(l[0].Ch), 3)
	<-ee.Emit("test2") // should raise blockedError
	// ^^^^ and remove the topic as well
	l = ee.Listeners("test2")
	expect(t, len(l), 0)
	expect(t, len(ee.Topics()), 0)
}

// func TestInvalidPattern(t *testing.T) {
// 	ee := New(0)
// 	ee.On("test")
// 	list, err := ee.Listeners("\\")
// 	expect(t, len(list), 0)
// 	expect(t, err != nil, true)
// 	expect(t, err.Error(), "syntax error in pattern")
//
// 	err = ee.Off("\\")
// 	expect(t, err.Error(), "syntax error in pattern")
// 	err = <-ee.Emit("\\")
// 	expect(t, err.Error(), "syntax error in pattern")
// }

func TestOnOffAll(t *testing.T) {
	ee := NewEmitter(0)
	ee.On("*")
	l := ee.Listeners("test")
	expect(t, len(l), 1)

	ee.Off("*")
	l = ee.Listeners("test")
	expect(t, len(l), 0)
}

func TestOrSkipOnce(t *testing.T) {
	ee := NewEmitter(0)
	pipe := ee.On("test", Skip, Once).Ch
	<-ee.Emit("test")
	l := ee.Listeners("test")
	expect(t, len(l), 1)
	go ee.Emit("test")
	<-pipe
	l = ee.Listeners("test")
	expect(t, len(l), 0)
}

func TestVoid(t *testing.T) {
	ee := NewEmitter(0)
	expect(t, len(ee.middlewares), 0)
	ee.Use("*", Void)
	expect(t, len(ee.middlewares), 1)
	ch := make(chan struct{})
	pipe := ee.On("test").Ch
	go func() {
		select {
		case <-pipe:
		default:
			ch <- struct{}{}
		}
	}()
	go ee.Emit("test")
	<-ch
	ee.Use("*")
	ee.Off("*", pipe)
	expect(t, len(ee.middlewares), 0)
	l := ee.Listeners("*")
	expect(t, len(l), 0)
	ee.On("test", Void)
	// unblocked, sending will be skipped
	<-ee.Emit("test")
}

func TestOnceClose(t *testing.T) {
	ee := NewEmitter(0)
	ee.On("test", Close, Once)
	// unblocked, the listener will be
	// closed after first attempt
	<-ee.Emit("test")
}

func TestCancellation(t *testing.T) {
	ee := NewEmitter(0)
	pipe := ee.On("test", Once).Ch
	ch := make(chan struct{})
	go func() {
		done := ee.Emit("test", 1)
		select {
		case <-done:
			expect(t, "cancellation success", "cancellation failure")
		case <-time.After(1e5):
			done <- struct{}{}
			ch <- struct{}{}
		}
	}()

	<-ch

	go ee.Emit("test", 2)
	l := ee.Listeners("*")
	expect(t, len(l), 1)
	e := <-pipe
	expect(t, e.Int(0), 2)
	expect(t, e.Flags, e.Flags|FlagOnce)
}

func TestSyncCancellation(t *testing.T) {
	ee := NewEmitter(0)
	pipe := ee.On("test", Once, Skip).Ch
	close(ee.Emit("test"))
	select {
	case e := <-pipe:
		expect(t, e, nil)
	default:
	}
}

func TestBackwardPattern(t *testing.T) {
	ee := NewEmitter(0)
	ee.Use("test", Close)
	go ee.Emit("test")
	e := <-ee.On("*", Once).Ch
	expect(t, e.OriginalTopic, "test")
	expect(t, e.Topic, "*")
	expect(t, e.Flags, e.Flags|FlagClose)
	expect(t, e.Flags, e.Flags|FlagOnce)
}

func TestResetMiddleware(t *testing.T) {
	ee := NewEmitter(0)
	ee.Use("*", Void, Reset)
	go ee.Emit("test")
	<-ee.On("test").Ch
}

func TestMiddleware(t *testing.T) {
	ee := NewEmitter(10)
	pipe := ee.On("test", func(e *Event) {
		if e.Int(0)%3 != 0 {
			e.Flags = e.Flags | FlagVoid
		}
	}).Ch
	pipe2 := ee.On("test").Ch

	for i := 0; i < 10; i++ {
		<-ee.Emit("test", i)
	}
	expect(t, len(pipe), 4)
	expect(t, len(pipe2), 10)
}

func TestSync(t *testing.T) {
	ee := NewEmitter(1)
	ee.Use("*", Sync)
	pipe := ee.On("test").Ch
	pipe2 := ee.On("test", Once).Ch
	_, isOpened := <-ee.Emit("test", 42)
	expect(t, len(pipe), 1)
	expect(t, len(pipe2), 1)

	expect(t, isOpened, false)

	e, isOpened := <-pipe2
	expect(t, e.Int(0), 42)
	expect(t, isOpened, true)
	_, isOpened = <-pipe2
	expect(t, isOpened, false)

	// void
	ee = NewEmitter(0)
	ee.Once("*", Void)
	ee.On("test:void", Void, func(e *Event) {})
	<-ee.Emit("test:void")
}

func TestCallbackOnlyUsage(t *testing.T) {
	ee := NewEmitter(0)
	ee.Use("*", Void)
	var called bool

	ee.On("call", func(e *Event) {
		called = e.Bool(0)
	})
	ee.Emit("call", true)
	expect(t, called, true)
}

func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func expectRev(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}