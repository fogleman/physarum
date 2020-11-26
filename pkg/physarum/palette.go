package physarum

import (
	"fmt"
	"image/color"
	"math/rand"
)

type Palette []color.RGBA

func ShuffledPalette(palette Palette) Palette {
	result := make(Palette, len(palette))
	for i, j := range rand.Perm(len(result)) {
		result[i] = palette[j]
	}
	return result
}

func RandomPalette() Palette {
	palette := Palettes[rand.Intn(len(Palettes))]
	return ShuffledPalette(palette)
}

func (p Palette) Print() {
	for _, c := range p {
		fmt.Printf("HexColor(0x%02X%02X%02X),\n", c.R, c.G, c.B)
	}
	fmt.Println()
}

var Palettes = []Palette{
	Palette{
		HexColor(0xFA2B31),
		HexColor(0xFFBF1F),
		HexColor(0xFFF146),
		HexColor(0xABE319),
		HexColor(0x00C481),
	},
	Palette{
		HexColor(0x004358),
		HexColor(0x1F8A70),
		HexColor(0xBEDB39),
		HexColor(0xFFE11A),
		HexColor(0xFD7400),
	},
	Palette{
		HexColor(0x334D5C),
		HexColor(0x45B29D),
		HexColor(0xEFC94C),
		HexColor(0xE27A3F),
		HexColor(0xDF5A49),
	},
	Palette{
		HexColor(0xFF8000),
		HexColor(0xFFD933),
		HexColor(0xCCCC52),
		HexColor(0x8FB359),
		HexColor(0x192B33),
	},
	Palette{
		HexColor(0x730046),
		HexColor(0xBFBB11),
		HexColor(0xFFC200),
		HexColor(0xE88801),
		HexColor(0xC93C00),
	},
	Palette{
		HexColor(0xE6DD00),
		HexColor(0x8CB302),
		HexColor(0x008C74),
		HexColor(0x004C66),
		HexColor(0x332B40),
	},
	Palette{
		HexColor(0xF15A5A),
		HexColor(0xF0C419),
		HexColor(0x4EBA6F),
		HexColor(0x2D95BF),
		HexColor(0x955BA5),
	},
	Palette{
		HexColor(0xF41C54),
		HexColor(0xFF9F00),
		HexColor(0xFBD506),
		HexColor(0xA8BF12),
		HexColor(0x00AAB5),
	},
}
