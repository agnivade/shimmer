package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"strings"
	"syscall/js"

	"github.com/anthonynsimon/bild/adjust"
	"github.com/anthonynsimon/bild/imgio"
)

const jpegPrefix = "data:image/jpeg;base64,"
const pngPrefix = "data:image/jpeg;base64,"

func main() {
	var loadImgCb, brightnessCb, contrastCb js.Callback
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
		applyTransformation(reader)
	})

	js.Global.Get("document").
		Call("getElementById", "sourceImg").
		Call("addEventListener", "load", loadImgCb)

	brightnessCb = js.NewCallback(func(args []js.Value) {
		fmt.Println(args[0].Get("target").Get("value").Float())
	})
	js.Global.Get("document").
		Call("getElementById", "brightness").
		Call("addEventListener", "change", brightnessCb)

	contrastCb = js.NewCallback(func(args []js.Value) {
		fmt.Println(args[0].Get("target").Get("value").Float())
	})
	js.Global.Get("document").
		Call("getElementById", "brightness").
		Call("addEventListener", "change", contrastCb)

	// Just waiting indefinitely for now
	select {}
}

func applyTransformation(r io.Reader) {
	i, _, err := image.Decode(r)
	if err != nil {
		fmt.Println(err)
		return
	}
	var buf bytes.Buffer
	res := adjust.Brightness(i, 0.75)

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
}
