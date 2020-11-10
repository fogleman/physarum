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

var Palette = []color.RGBA{
	HexColor(0xFA2B31),
	HexColor(0xFFBF1F),
	HexColor(0xFFF146),
	HexColor(0xABE319),
	HexColor(0x00C481),
}

func Image(w, h int, colors [][]float64, min, max, gamma float64) image.Image {
	rand.Shuffle(len(Palette), func(i, j int) {
		Palette[i], Palette[j] = Palette[j], Palette[i]
	})

	im := image.NewRGBA(image.Rect(0, 0, w, h))

	minValues := make([]float64, len(colors))
	maxValues := make([]float64, len(colors))
	for i, data := range colors {
		minValues[i] = min
		maxValues[i] = max
		if min == max {
			temp := make([]float64, len(data))
			copy(temp, data)
			sort.Float64s(temp)
			minValues[i] = stat.Quantile(0.01, stat.Empirical, temp, nil)
			maxValues[i] = stat.Quantile(0.99, stat.Empirical, temp, nil)
		}
		minValues[i] = 0
		c := Palette[i]
		fmt.Printf("%d #%02X%02X%02X %.3f\n", i, c.R, c.G, c.B, maxValues[i])
	}

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			index := y*w + x
			var r, g, b float64
			for i, data := range colors {
				t := data[index]
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
				c := Palette[i]
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
