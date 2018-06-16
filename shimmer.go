package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"strings"
	"syscall/js"

	"github.com/anthonynsimon/bild/adjust"
	"github.com/anthonynsimon/bild/imgio"
)

const jpegPrefix = "data:image/jpeg;base64,"
const pngPrefix = "data:image/jpeg;base64,"

func main() {
	var loadImgCb, brightnessCb, contrastCb js.Callback

	var i image.Image
	var err error
	var buf bytes.Buffer
	// TODO: explicitly close callback when done
	loadImgCb = js.NewCallback(func(args []js.Value) {
		source := js.Global.
			Get("document").Call("getElementById", "sourceImg").
			Get("src").String()

		switch {
		case strings.HasPrefix(source, jpegPrefix):
			source = strings.TrimPrefix(source, jpegPrefix)
		case strings.HasPrefix(source, pngPrefix):
			source = strings.TrimPrefix(source, pngPrefix)
		default:
			// TODO: log this in the status div
			fmt.Println("unrecognized image format")
			return
		}

		reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(source))
		i, _, err = image.Decode(reader)
		if err != nil {
			fmt.Println(err)
			return
		}
		// TODO: log this in status div
		fmt.Println("Ready for operations")
		// TODO: reset brightness and contrast sliders
	})

	js.Global.Get("document").
		Call("getElementById", "sourceImg").
		Call("addEventListener", "load", loadImgCb)

	brightnessCb = js.NewEventCallback(js.PreventDefault, func(ev js.Value) {
		b := ev.Get("target").Get("value").Float()
		res := adjust.Brightness(i, b)

		buf.Reset()
		enc := imgio.JPEGEncoder(90)
		err = enc(&buf, res)
		if err != nil {
			fmt.Println(err)
			return
		}
		// Updating the image
		js.Global.Get("document").
			Call("getElementById", "targetImg").
			Set("src", jpegPrefix+base64.StdEncoding.EncodeToString(buf.Bytes()))
	})
	js.Global.Get("document").
		Call("getElementById", "brightness").
		Call("addEventListener", "change", brightnessCb)

	contrastCb = js.NewEventCallback(js.PreventDefault, func(ev js.Value) {
		c := ev.Get("target").Get("value").Float()
		fmt.Println(c)
	})
	js.Global.Get("document").
		Call("getElementById", "contrast").
		Call("addEventListener", "change", contrastCb)

	// Just waiting indefinitely for now
	select {}
}
