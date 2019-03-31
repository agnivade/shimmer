// +build js,wasm

package shimmer

import (
	"bytes"
	"image"
	"reflect"
	"syscall/js"
	"unsafe"
)

func (s *Shimmer) setupOnImgLoadCb() {
	s.onImgLoadCb = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		reader := bytes.NewReader(s.inBuf)
		var err error
		s.sourceImg, _, err = image.Decode(reader)
		if err != nil {
			s.log(err.Error())
			return nil
		}
		s.log("Ready for operations")

		// reset brightness and contrast sliders
		js.Global().Get("document").
			Call("getElementById", "brightness").
			Set("value", 0)

		js.Global().Get("document").
			Call("getElementById", "contrast").
			Set("value", 0)
		return nil
	})
}

func (s *Shimmer) setupInitMemCb() {
	// The length of the image array buffer is passed.
	// Then the buf slice is initialized to that length.
	// And a pointer to that slice is passed back to the browser.
	s.initMemCb = js.FuncOf(func(this js.Value, i []js.Value) interface{} {
		length := i[0].Int()
		s.console.Call("log", "length:", length)
		s.inBuf = make([]uint8, length)
		hdr := (*reflect.SliceHeader)(unsafe.Pointer(&s.inBuf))
		ptr := uintptr(unsafe.Pointer(hdr.Data))
		s.console.Call("log", "ptr:", ptr)
		js.Global().Call("gotMem", ptr)
		return nil
	})
}
