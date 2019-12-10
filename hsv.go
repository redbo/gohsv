package gohsv

import (
	"image"
	"image/color"
	"math"
)

// HSVtoRGB takes a color in HSV space with values hue(0.0 - 360.0),
// saturation (0 - 1.0) and value (0-1.0) and returns its representation
// in RGB color space, with values 0 - 0xFFFF.
func HSVtoRGB(h, s, v float64) (r, g, b uint32) {
	h, f := math.Modf(h / 60.0)
	p := uint32(math.Round((v * (1.0 - s)) * 0xffff))
	q := uint32(math.Round((v * (1.0 - (s * f))) * 0xffff))
	t := uint32(math.Round((v * (1.0 - (s * (1.0 - f)))) * 0xffff))
	vr := uint32(math.Round(v * 0xffff))
	switch int(h) {
	default:
		return vr, t, p
	case 1:
		return q, vr, p
	case 2:
		return p, vr, t
	case 3:
		return p, q, vr
	case 4:
		return t, p, vr
	case 5:
		return vr, p, q
	}
}

// RGBtoHSV takes a color in RGB space with values from 0 - 0xffff
// and returns the corresponding HSV representation as hue (0.0 - 360.0),
// saturation (0 - 1.0) and value (0 - 1.0).
func RGBtoHSV(r, g, b uint32) (h, s, v float64) {
	rf := float64(r) / 0xffff
	gf := float64(g) / 0xffff
	bf := float64(b) / 0xffff
	cmax := math.Max(rf, math.Max(gf, bf))
	cmin := math.Min(rf, math.Min(gf, bf))
	diff := cmax - cmin
	if cmax == cmin {
		h = 0.0
	} else if cmax == rf {
		h = math.Mod(60*((gf-bf)/diff)+360, 360)
	} else if cmax == gf {
		h = math.Mod(60*((bf-rf)/diff)+120, 360)
	} else if cmax == bf {
		h = math.Mod(60*((rf-gf)/diff)+240, 360)
	}

	if cmax == 0 {
		return h, 0, cmax
	}
	return h, (diff / cmax), cmax
}

// HSV represents a color in hue-saturation-value space.  Hue is a
// value between 0 and 360, saturation and value are represented as
// numbers between 0 and 1.0.
type HSV struct {
	H, S, V float64
}

func (c *HSV) RGBA() (uint32, uint32, uint32, uint32) {
	r, g, b := HSVtoRGB(c.H, c.S, c.V)
	return uint32(r), uint32(g), uint32(b), 0xFFFF
}

func hsvModel(c color.Color) color.Color {
	r, g, b, _ := c.RGBA()
	h, s, v := RGBtoHSV(r, g, b)
	return &HSV{H: h, S: s, V: v}
}

var HSVModel = color.ModelFunc(hsvModel)

// HSVImage is an in-memory image which stores image data
// in hue-saturation-value color space.
type HSVImage struct {
	Pix  []HSV
	Rect image.Rectangle
}

func (i *HSVImage) ColorModel() color.Model {
	return HSVModel
}

func (i *HSVImage) Bounds() image.Rectangle {
	return i.Rect
}

func (i *HSVImage) At(x, y int) color.Color {
	return &i.Pix[y*i.Rect.Dx()+x]
}

func (i *HSVImage) Set(x, y int, c color.Color) {
	r, g, b, _ := c.RGBA()
	h, s, v := RGBtoHSV(r, g, b)
	i.Pix[y*i.Rect.Dx()+x] = HSV{H: h, S: s, V: v}
}

var _ = image.Image(&HSVImage{}) // verify HSVImage satisfies the Image interface

// NewHSV returns a new HSVImage with the given bounds.
func NewHSV(r image.Rectangle) *HSVImage {
	return &HSVImage{
		Pix:  make([]HSV, r.Dx()*r.Dy()),
		Rect: r,
	}
}
