package gohsv

import "testing"

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
