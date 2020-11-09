package physarum

import (
	"image"
	"image/color"
	"math"
)

func Image(w, h int, colors [][]float64) image.Image {
	im := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			index := y*w + x
			var r, g, b float64
			for i, data := range colors {
				t := data[index] / 2
				if t > 1 {
					t = 1
				}
				t = math.Pow(t, 1/2.2)
				switch i {
				case 0:
					r = t
				case 1:
					b = t
				case 2:
					g = t
				}
			}
			c := color.RGBA{
				uint8(r * 255),
				uint8(g * 255),
				uint8(b * 255),
				255,
			}
			im.SetRGBA(x, y, c)
		}
	}
	return im

}
