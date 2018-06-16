package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/jpeg"
	"image/png"
	"io"
	"strings"
	"syscall/js"
)

const jpegPrefix = "data:image/jpeg;base64,"
const pngPrefix = "data:image/jpeg;base64,"

func main() {
	var loadImgCb js.Callback
	// TODO: explicitly close callback when done
	loadImgCb = js.NewCallback(func(args []js.Value) {
		// this does not show up - investigate
		// fmt.Println(args[0].Get("target"))

		source := js.Global.
			Get("document").Call("getElementById", "sourceImg").
			Get("src").String()
		if strings.HasPrefix(source, jpegPrefix) {
			source = strings.TrimPrefix(source, jpegPrefix)
		} else {
			fmt.Println("Unrecognized image format")
			return
		}
		reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(source))
		convertToPNG(reader)
	})
	js.Global.Get("document").
		Call("getElementById", "sourceImg").
		Call("addEventListener", "load", loadImgCb)
	fmt.Println("hello wasm")

	// Just waiting indefinitely for now
	select {}
}

func convertToPNG(r io.Reader) {
	i, err := jpeg.Decode(r)
	if err != nil {
		fmt.Println(err)
		return
	}
	var buf bytes.Buffer
	err = png.Encode(&buf, i)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Updating the image
	js.Global.Get("document").
		Call("getElementById", "targetImg").
		Set("src", pngPrefix+base64.StdEncoding.EncodeToString(buf.Bytes()))
}
