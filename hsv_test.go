package gohsv

import "testing"

func TestRoundtripColors(t *testing.T) {
	testCases := [][3]uint32{
		{30, 0xf, 0xf},
		{90, 0xff, 0xf0},
		{150, 0xff, 0xff},
		{210, 0xf, 0x0},
		{270, 0xff1, 0xf1},
		{330, 0xff0, 0xf0},
	}
	for _, tc := range testCases {
		h, s, v := RGBtoHSV(tc[0], tc[1], tc[2])
		r, g, b := HSVtoRGB(h, s, v)
		if r != tc[0] || g != tc[1] || b != tc[2] {
			t.Logf("Failed to roundtrip color: %d, %d, %d -> %d, %d, %d",
				tc[0], tc[1], tc[2], r, g, b)
			t.Fail()
		}
	}
}
