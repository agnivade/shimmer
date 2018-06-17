package shimmer

import (
	"bytes"
	"encoding/base64"
	"testing"

	"github.com/anthonynsimon/bild/adjust"
	"github.com/anthonynsimon/bild/imgio"
)

func BenchmarkAdjustImage(b *testing.B) {
	img, err := imgio.Open("testdata/dragon.jpg")
	if err != nil {
		b.Error(err)
	}
	// var buf strings.Builder
	var buf bytes.Buffer
	// b64Enc := base64.NewEncoder(base64.StdEncoding, &buf)

	var sink string
	for i := 0; i < b.N; i++ {
		img2 := adjust.Brightness(img, 0.4)
		enc := imgio.JPEGEncoder(90)
		err = enc(&buf, img2)
		if err != nil {
			b.Error(err)
		}
		// b64Enc.Close()
		sink = base64.StdEncoding.EncodeToString(buf.Bytes())
		buf.Reset()
	}
	b.Log(sink[:10])
}
