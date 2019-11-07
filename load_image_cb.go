// +build js,wasm

package shimmer

import (
	"bytes"
	"image"
	"syscall/js"
)

func (s *Shimmer) setupOnImgLoadCb() {
	s.onImgLoadCb = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		array := args[0]
		s.inBuf = make([]uint8, array.Get("byteLength").Int())
		js.CopyBytesToGo(s.inBuf, array)

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
