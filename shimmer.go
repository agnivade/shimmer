// +build js,wasm

package shimmer

import (
	"bytes"
	"image"
	"syscall/js"
)

type Shimmer struct {
	buf                                   bytes.Buffer
	onImgLoadCb, brightnessCb, contrastCb js.Callback
	shutdownCb                            js.Callback
	sourceImg                             image.Image

	done chan struct{}
}

// New returns a new instance of shimmer
func New() *Shimmer {
	return &Shimmer{
		done: make(chan struct{}),
	}
}

// Start sets up all the callbacks and waits for the close signal
// to be sent from the browser.
func (s *Shimmer) Start() {
	// Setup callbacks
	s.setupOnImgLoadCb()
	js.Global.Get("document").
		Call("getElementById", "sourceImg").
		Call("addEventListener", "load", s.onImgLoadCb)

	s.setupBrightnessCb()
	js.Global.Get("document").
		Call("getElementById", "brightness").
		Call("addEventListener", "change", s.brightnessCb)

	s.setupContrastCb()
	js.Global.Get("document").
		Call("getElementById", "contrast").
		Call("addEventListener", "change", s.contrastCb)

	s.setupShutdownCb()
	js.Global.Get("document").
		Call("getElementById", "close").
		Call("addEventListener", "click", s.shutdownCb)

	<-s.done
	s.log("Shutting down app")
	s.onImgLoadCb.Close()
	s.brightnessCb.Close()
	s.contrastCb.Close()
	s.shutdownCb.Close()
}

// utility function to log a msg to the UI from inside a callback
func (s *Shimmer) log(msg string) {
	js.Global.Get("document").
		Call("getElementById", "status").
		Set("innerText", msg)
}

func (s *Shimmer) setupShutdownCb() {
	s.shutdownCb = js.NewEventCallback(js.PreventDefault, func(ev js.Value) {
		s.done <- struct{}{}
	})
}
