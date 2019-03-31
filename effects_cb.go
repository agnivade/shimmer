// +build js,wasm

package shimmer

import (
	"syscall/js"
	"time"

	"github.com/anthonynsimon/bild/adjust"
)

func (s *Shimmer) setupBrightnessCb() {
	s.brightnessCb = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// quick return if no source image is yet uploaded
		if s.sourceImg == nil {
			return nil
		}
		delta := args[0].Get("target").Get("valueAsNumber").Float()
		start := time.Now()
		res := adjust.Brightness(s.sourceImg, delta)
		s.updateImage(res, start)
		args[0].Call("preventDefault")
		return nil
	})
}

func (s *Shimmer) setupContrastCb() {
	s.contrastCb = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// quick return if no source image is yet uploaded
		if s.sourceImg == nil {
			return nil
		}
		delta := args[0].Get("target").Get("valueAsNumber").Float()
		start := time.Now()
		res := adjust.Contrast(s.sourceImg, delta)
		s.updateImage(res, start)
		args[0].Call("preventDefault")
		return nil
	})
}

func (s *Shimmer) setupHueCb() {
	s.hueCb = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// quick return if no source image is yet uploaded
		if s.sourceImg == nil {
			return nil
		}
		delta := args[0].Get("target").Get("valueAsNumber").Int()
		start := time.Now()
		res := adjust.Hue(s.sourceImg, delta)
		s.updateImage(res, start)
		args[0].Call("preventDefault")
		return nil
	})
}

func (s *Shimmer) setupSatCb() {
	s.satCb = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// quick return if no source image is yet uploaded
		if s.sourceImg == nil {
			return nil
		}
		delta := args[0].Get("target").Get("valueAsNumber").Float()
		start := time.Now()
		res := adjust.Saturation(s.sourceImg, delta)
		s.updateImage(res, start)
		args[0].Call("preventDefault")
		return nil
	})
}
