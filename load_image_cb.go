// +build js,wasm

package shimmer

import (
	"encoding/base64"
	"image"
	"strings"
	"syscall/js"
)

const (
	jpegPrefix = "data:image/jpeg;base64,"
	pngPrefix  = "data:image/png;base64,"
)

func (s *Shimmer) setupOnImgLoadCb() {
	s.onImgLoadCb = js.NewCallback(func(args []js.Value) {
		source := js.Global.
			Get("document").Call("getElementById", "sourceImg").
			Get("src").String()

		switch {
		case strings.HasPrefix(source, jpegPrefix):
			source = strings.TrimPrefix(source, jpegPrefix)
		case strings.HasPrefix(source, pngPrefix):
			source = strings.TrimPrefix(source, pngPrefix)
		default:
			s.log("unrecognized image format")
			return
		}

		reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(source))
		var err error
		s.sourceImg, _, err = image.Decode(reader)
		if err != nil {
			s.log(err.Error())
			return
		}
		s.log("Ready for operations")

		// reset brightness and contrast sliders
		js.Global.Get("document").
			Call("getElementById", "brightness").
			Set("value", 0)

		js.Global.Get("document").
			Call("getElementById", "contrast").
			Set("value", 0)
	})
}
