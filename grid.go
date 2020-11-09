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

type Grid struct {
	W    int
	H    int
	Data []float64
	Temp []float64
}

func NewGrid(w, h int) *Grid {
	data := make([]float64, w*h)
	temp := make([]float64, w*h)
	for i := range data {
		data[i] = rand.Float64()
	}
	return &Grid{w, h, data, temp}
}

func (g *Grid) Get(x, y float64) float64 {
	if x < 0 {
		x += float64(g.W)
	}
	if y < 0 {
		y += float64(g.H)
	}
	x0 := int(x)
	y0 := int(y)
	x1 := x0 + 1
	y1 := y0 + 1
	x -= float64(x0)
	y -= float64(y0)
	if x0 >= g.W {
		x0 -= g.W
	}
	if y0 >= g.H {
		y0 -= g.H
	}
	if x1 >= g.W {
		x1 -= g.W
	}
	if y1 >= g.H {
		y1 -= g.H
	}
	var d float64
	d += g.Data[x0+y0*g.W] * ((1 - x) * (1 - y))
	d += g.Data[x0+y1*g.W] * ((1 - x) * y)
	d += g.Data[x1+y0*g.W] * (x * (1 - y))
	d += g.Data[x1+y1*g.W] * (x * y)
	return d
}

func (g *Grid) Add(x, y, a float64) {
	if x < 0 {
		x += float64(g.W)
	}
	if y < 0 {
		y += float64(g.H)
	}
	x0 := int(x)
	y0 := int(y)
	x1 := x0 + 1
	y1 := y0 + 1
	x -= float64(x0)
	y -= float64(y0)
	if x0 >= g.W {
		x0 -= g.W
	}
	if y0 >= g.H {
		y0 -= g.H
	}
	if x1 >= g.W {
		x1 -= g.W
	}
	if y1 >= g.H {
		y1 -= g.H
	}
	g.Data[x0+y0*g.W] += a * ((1 - x) * (1 - y))
	g.Data[x0+y1*g.W] += a * ((1 - x) * y)
	g.Data[x1+y0*g.W] += a * (x * (1 - y))
	g.Data[x1+y1*g.W] += a * (x * y)
}

func (g *Grid) DiffuseAndDecay(decayFactor float64) {
	gaussianBlur(g.Data, g.Temp, g.W, g.H, 1, decayFactor)
}

func (g *Grid) Image(lo, hi, gamma float64) image.Image {
	// cm := colormap.Inferno
	// im := image.NewRGBA(image.Rect(0, 0, g.W, g.H))

	im := image.NewGray(image.Rect(0, 0, g.W, g.H))

	if hi <= 0 {
		copy(g.Temp, g.Data)
		sort.Float64s(g.Temp)
		hi = stat.Quantile(0.99, stat.Empirical, g.Temp, nil)
		fmt.Println(hi)
	}

	for y := 0; y < g.H; y++ {
		for x := 0; x < g.W; x++ {
			index := y*g.W + x
			t := g.Data[index] / hi
			if t > 1 {
				t = 1
			}
			if gamma != 1 {
				t = math.Pow(t, gamma)
			}
			// im.Set(x, y, cm.At(t))
			c := color.Gray{uint8(t * 0xff)}
			im.SetGray(x, y, c)
		}
	}

	return im
}
