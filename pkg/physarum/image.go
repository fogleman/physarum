package physarum

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"sort"

	"github.com/gonum/stat"
)

func Image(w, h int, grids [][]float32, palette Palette, min, max, gamma float32) image.Image {
	im := image.NewRGBA(image.Rect(0, 0, w, h))

	minValues := make([]float32, len(grids))
	maxValues := make([]float32, len(grids))
	for i, grid := range grids {
		minValues[i] = min
		maxValues[i] = max
		if min == max {
			temp := make([]float64, len(grid))
			for i, v := range grid {
				temp[i] = float64(v)
			}
			sort.Float64s(temp)
			minValues[i] = 0
			// minValues[i] = stat.Quantile(0.01, stat.Empirical, temp, nil)
			maxValues[i] = float32(stat.Quantile(0.99, stat.Empirical, temp, nil))
			c := palette[i]
			fmt.Printf("%d #%02X%02X%02X %.3f\n", i, c.R, c.G, c.B, maxValues[i])
		}
	}

	makeGammaLookupTable := func(gamma float32) func(t float32) float32 {
		const size = 65536
		table := make([]float32, size)
		for i := range table {
			t := float64(i) / (size - 1)
			table[i] = float32(math.Pow(t, float64(gamma)))
		}
		factor := float32(size - 1)
		return func(t float32) float32 {
			i := int(t * factor)
			return table[i]
		}
	}
	lookup := makeGammaLookupTable(gamma)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			index := y*w + x
			var r, g, b float32
			for i, grid := range grids {
				t := grid[index]
				t = (t - minValues[i]) / (maxValues[i] - minValues[i])
				if t < 0 {
					t = 0
				}
				if t > 1 {
					t = 1
				}
				t = lookup(t)
				c := palette[i]
				r += float32(c.R) * t
				g += float32(c.G) * t
				b += float32(c.B) * t
			}
			if r > 255 {
				r = 255
			}
			if g > 255 {
				g = 255
			}
			if b > 255 {
				b = 255
			}
			c := color.RGBA{uint8(r), uint8(g), uint8(b), 255}
			im.SetRGBA(x, y, c)
		}
	}
	return im

}
