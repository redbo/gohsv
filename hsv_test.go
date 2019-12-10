package gohsv

import (
	"image"
	"image/color"
	"image/draw"
	"math"
	"testing"
)

func TestRGBtoHSV(t *testing.T) {
	// https://en.wikipedia.org/wiki/HSL_and_HSV#Examples
	examples := []struct {
		color   uint32
		h, v, s float64
	}{
		{0xFF0000, 0.0, 1, 1},
		{0xBFBF00, 60.0, 0.75, 1},
		{0x008000, 120.0, 0.5, 1},
		{0x80FFFF, 180.0, 1, 0.5},
		{0x8080FF, 240.0, 1, 0.5},
		{0xBF40BF, 300.0, 0.75, 0.667},
		{0xA0A424, 61.8, 0.643, 0.779},
		{0x411BEA, 251.1, 0.918, 0.887},
		{0x1EAC41, 134.9, 0.675, 0.828},
		{0xF0C80E, 49.5, 0.941, 0.944},
		{0xB430E5, 283.7, 0.897, 0.792},
		{0xED7651, 14.3, 0.931, 0.661},
		{0xFEF888, 56.9, 0.998, 0.467},
		{0x19CB97, 162.4, 0.795, 0.875},
		{0x362698, 248.3, 0.597, 0.75},
		{0x7E7EB8, 240.5, 0.721, 0.316},
	}
	withinTolerance := func(a, b float64) bool {
		// test if the numbers differ by less than 1% or some tiny absolute number
		diff := math.Abs(a - b)
		return diff < a/100.0 || diff < 0.000001
	}
	for _, ex := range examples {
		r := ((ex.color >> 16) & 0xFF) * 0xffff
		g := ((ex.color >> 8) & 0xFF) * 0xffff
		b := (ex.color & 0xFF) * 0xffff
		h, s, v := RGBtoHSV(r, g, b)
		if !withinTolerance(ex.h, h) {
			t.Logf("Hue incorrect for %x (have: %f, want: %f)", ex.color, h, ex.h)
			t.Fail()
		}
		if !withinTolerance(ex.s, s) {
			t.Logf("Sat incorrect for %x (have: %f, want: %f)", ex.color, s, ex.s)
			t.Fail()
		}
		if !withinTolerance(ex.v, v) {
			t.Logf("Value incorrect for %x (have: %f, want: %f)", ex.color, v, ex.v)
			t.Fail()
		}
	}
}

func TestRoundtripColors(t *testing.T) {
	testCases := []uint32{
		0x0, 0x1, 0x9fc9, 0x3d69, 0x5e69,
		0xcc25, 0x84b6, 0xe095, 0x3d6d,
		0xfa91, 0x1c93, 0xe829, 0x9f5,
		0xbca2, 0xb678, 0xc1a4,
	}
	for _, tr := range testCases {
		for _, tg := range testCases {
			for _, tb := range testCases {
				h, s, v := RGBtoHSV(tr, tg, tb)
				r, g, b := HSVtoRGB(h, s, v)
				if r != tr || g != tg || b != tb {
					t.Logf("Failed to roundtrip color: %d, %d, %d -> %d, %d, %d",
						tr, tg, tb, r, g, b)
					t.Fail()
				}
			}
		}
	}
}

func TestHSVImage(t *testing.T) {
	// make a random image
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			img.Set(x, y, color.RGBA{
				R: uint8(x * 20),
				G: uint8(y * 20),
				B: uint8((x + y) * 10),
				A: 255,
			})
		}
	}
	// roundtrip it through hsv
	imgcopy := image.NewRGBA(img.Bounds())
	hsv := NewHSV(img.Bounds())
	draw.Draw(hsv, hsv.Bounds(), img, image.ZP, draw.Over)
	draw.Draw(imgcopy, imgcopy.Bounds(), hsv, image.ZP, draw.Over)
	// make sure it didn't change
	for i := range img.Pix {
		if imgcopy.Pix[i] != img.Pix[i] {
			t.Logf("Failed to roundtrip image: byte %d", i)
			t.Fail()
		}
	}
}
