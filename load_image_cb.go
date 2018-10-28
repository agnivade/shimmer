// +build js,wasm

package shimmer

import (
	"bytes"
	"image"
	"reflect"
	"syscall/js"
	"unsafe"
)

const (
	jpegPrefix = "data:image/jpeg;base64,"
	pngPrefix  = "data:image/png;base64,"
)

func (s *Shimmer) setupOnImgLoadCb() {
	s.onImgLoadCb = js.NewCallback(func(args []js.Value) {
		reader := bytes.NewReader(s.buf2)
		var err error
		s.sourceImg, _, err = image.Decode(reader)
		if err != nil {
			s.log(err.Error())
			return
		}
		s.log("Ready for operations")

		// reset brightness and contrast sliders
		js.Global().Get("document").
			Call("getElementById", "brightness").
			Set("value", 0)

		js.Global().Get("document").
			Call("getElementById", "contrast").
			Set("value", 0)
	})
}

func (s *Shimmer) setupInitMemCb() {
	// The length of the image array buffer is passed.
	// Then the buf slice is initialized to that length.
	// And a pointer to that slice is passed back to the browser.
	s.initMemCb = js.NewCallback(func(i []js.Value) {
		length := i[0].Int()
		s.console.Call("log", "length:", length)
		s.buf2 = make([]uint8, length)
		hdr := (*reflect.SliceHeader)(unsafe.Pointer(&s.buf2))
		ptr := uintptr(unsafe.Pointer(hdr.Data))
		s.console.Call("log", "ptr:", ptr)
		js.Global().Call("gotMem", ptr)
	})
}
