package physarum

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"math/rand"
	"sort"

	"github.com/gonum/stat"
)

var DefaultPalette = []color.RGBA{
	HexColor(0xFA2B31),
	HexColor(0xFFBF1F),
	HexColor(0xFFF146),
	HexColor(0xABE319),
	HexColor(0x00C481),
}

func ShuffledPalette(palette []color.RGBA) []color.RGBA {
	result := make([]color.RGBA, len(palette))
	for i, j := range rand.Perm(len(result)) {
		result[i] = palette[j]
	}
	return result
}

func Image(w, h int, grids [][]float64, palette []color.RGBA, min, max, gamma float64) image.Image {
	im := image.NewRGBA(image.Rect(0, 0, w, h))

	minValues := make([]float64, len(grids))
	maxValues := make([]float64, len(grids))
	for i, grid := range grids {
		minValues[i] = min
		maxValues[i] = max
		if min == max {
			temp := make([]float64, len(grid))
			copy(temp, grid)
			sort.Float64s(temp)
			minValues[i] = 0
			// minValues[i] = stat.Quantile(0.01, stat.Empirical, temp, nil)
			maxValues[i] = stat.Quantile(0.99, stat.Empirical, temp, nil)
			c := palette[i]
			fmt.Printf("%d #%02X%02X%02X %.3f\n", i, c.R, c.G, c.B, maxValues[i])
		}
	}

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			index := y*w + x
			var r, g, b float64
			for i, grid := range grids {
				t := grid[index]
				t = (t - minValues[i]) / (maxValues[i] - minValues[i])
				if t < 0 {
					t = 0
				}
				if t > 1 {
					t = 1
				}
				if gamma != 1 {
					t = math.Pow(t, gamma)
				}
				c := palette[i]
				r += float64(c.R) * t
				g += float64(c.G) * t
				b += float64(c.B) * t
			}
			c := color.RGBA{
				uint8(math.Min(r, 255)),
				uint8(math.Min(g, 255)),
				uint8(math.Min(b, 255)),
				255,
			}
			im.SetRGBA(x, y, c)
		}
	}
	return im

}
