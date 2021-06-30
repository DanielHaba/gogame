package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

func test() bool {
	var err error
	f, err := os.Open("C:\\Windows\\Fonts\\Arial.ttf")
	if err != nil {
		panic(err)
	}
	fontFace, err := opentype.ParseReaderAt(f)
	if err != nil {
		panic(err)
	}
	face, err := opentype.NewFace(fontFace, &opentype.FaceOptions{
		Size: 72,
		DPI:  72,
	})
	if err != nil {
		panic(err)
	}
	adv, ok := face.GlyphAdvance('A')
	fmt.Printf("glyphAdvance: %d (%v)\n", adv, ok)
	bounds, adv, ok := face.GlyphBounds('A')
	fmt.Printf("glyphBounds: %+v / %d (%v)\n", bounds, adv, ok)
	metrics := face.Metrics()
	fmt.Printf("metrics: %+v\n", metrics)

	text := "lazy_grey fox jumped over big fence"
	padding := 80

	img := image.NewRGBA(image.Rect(0, 0, 800, 600))
	drawer := font.Drawer{
		Face: face,
		Src:  image.White,
		Dst:  img,
		Dot:  fixed.P(padding, padding),
	}
	drawRect(img, image.Rect(80, 80, 720, 520), color.RGBA{255, 0, 0, 255})

	for _, ch := range text {
		bounds, adv = drawer.BoundString(string(ch))
		//fmt.Printf("stringBounds: %+v\n", bounds)
		//fmt.Printf("stringBounds: %+v\n", boundsToRect(bounds))
		//fmt.Printf("stringBounds: %+v\n", boundsFloat64(bounds))

		//rect := metricsToRect(drawer.Dot, metrics, adv)
		rect := boundsToRect(bounds)
		fmt.Printf("'%c': %+v\n", ch, rect)
		drawRect(img, rect, color.RGBA{0, 255, 0, 255})
		rect = image.Rect(
			int(bounds.Min.X)>>6,
			int(bounds.Max.Y-metrics.Ascent)>>6,
			int(bounds.Min.X+adv)>>6,
			int(bounds.Min.Y)>>6,
		)

		drawer.DrawString(string(ch))
		drawRect(img, rect, color.RGBA{0, 0, 255, 255})
		//break
	}

	//bounds, _ = drawer.BoundString(text)
	//fmt.Printf("stringBounds: %+v\n", bounds)
	//fmt.Printf("stringBounds: %+v\n", boundsToRect(bounds))
	//fmt.Printf("stringBounds: %+v\n", boundsFloat64(bounds))
	//drawRect(img, boundsToRect(bounds), color.RGBA{0, 255, 0, 255})
	//drawer.DrawString(text)

	face.Close()
	f.Close()
	f, err = os.Create("out.jpeg")
	if err != nil {
		panic(err)
	}
	jpeg.Encode(f, img, nil)
	f.Close()
	return true
}

func boundsToRect(bounds fixed.Rectangle26_6) image.Rectangle {
	return image.Rect(
		int(bounds.Min.X)>>6,
		int(bounds.Min.Y)>>6,
		int(bounds.Max.X)>>6,
		int(bounds.Max.Y)>>6,
	)
}

func boundsFloat64(bounds fixed.Rectangle26_6) []float64 {
	return []float64{
		float64(bounds.Min.X) / 72,
		float64(bounds.Max.X) / 72,
		float64(bounds.Min.Y) / 72,
		float64(bounds.Max.Y) / 72,
	}
}

func metricsToRect(start fixed.Point26_6, metrics font.Metrics, advance fixed.Int26_6) image.Rectangle {
	return image.Rect(
		int(start.X)>>6,
		int(start.Y-metrics.Ascent)>>6,
		int(start.X+advance)>>6,
		int(start.Y+metrics.Descent)>>6,
	)
}

func drawRect(img draw.Image, rect image.Rectangle, color color.Color) {
	drawHLine(img, rect.Min.Y, rect.Min.X, rect.Max.X, color)
	drawHLine(img, rect.Max.Y, rect.Min.X, rect.Max.X, color)
	drawVLine(img, rect.Min.X, rect.Min.Y, rect.Max.Y, color)
	drawVLine(img, rect.Max.X, rect.Min.Y, rect.Max.Y, color)
}

func drawHLine(img draw.Image, y, x0, x1 int, color color.Color) {
	for x := x0; x <= x1; x++ {
		img.Set(x, y, color)
	}
}

func drawVLine(img draw.Image, x, y0, y1 int, color color.Color) {
	for y := y0; y <= y1; y++ {
		img.Set(x, y, color)
	}
}
