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

type HSV struct {
	H, S, V float64
}

func (c *HSV) RGBA() (uint32, uint32, uint32, uint32) {
	r, g, b := HSVtoRGB(c.H, c.S, c.V)
	return uint32(r), uint32(g), uint32(b), 0xFFFF
}

// HSVImage is an in-memory image whose At method returns HSV values.
type HSVImage struct {
	Pix    []HSV
	Stride int
	Rect   image.Rectangle
}

func (i *HSVImage) ColorModel() color.Model {
	return color.ModelFunc(
		func(c color.Color) color.Color {
			r, g, b, _ := c.RGBA()
			h, s, v := RGBtoHSV(r, g, b)
			return &HSV{
				H: h,
				S: s,
				V: v,
			}
		})
}

func (i *HSVImage) Bounds() image.Rectangle {
	return image.Rectangle{}
}

func (i *HSVImage) At(x, y int) color.Color {
	return nil
}

// verify HSVImage satisfies the Image interface
var _ = image.Image(&HSVImage{})
