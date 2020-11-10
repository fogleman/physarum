package physarum

import (
	"math/rand"
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
	// i := int(x+float64(g.W)) % g.W
	// j := int(y+float64(g.H)) % g.H
	i := int(x+float64(g.W)) & (g.W - 1)
	j := int(y+float64(g.H)) & (g.H - 1)
	index := j*g.W + i
	return g.Data[index]

	// if x < 0 {
	// 	x += float64(g.W)
	// }
	// if y < 0 {
	// 	y += float64(g.H)
	// }
	// x0 := int(x)
	// y0 := int(y)
	// x1 := x0 + 1
	// y1 := y0 + 1
	// x -= float64(x0)
	// y -= float64(y0)
	// if x0 >= g.W {
	// 	x0 -= g.W
	// }
	// if y0 >= g.H {
	// 	y0 -= g.H
	// }
	// if x1 >= g.W {
	// 	x1 -= g.W
	// }
	// if y1 >= g.H {
	// 	y1 -= g.H
	// }
	// var d float64
	// d += g.Data[x0+y0*g.W] * ((1 - x) * (1 - y))
	// d += g.Data[x0+y1*g.W] * ((1 - x) * y)
	// d += g.Data[x1+y0*g.W] * (x * (1 - y))
	// d += g.Data[x1+y1*g.W] * (x * y)
	// return d
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
	boxBlur(g.Data, g.Temp, g.W, g.H, 1, decayFactor)
	// gaussianBlur(g.Data, g.Temp, g.W, g.H, 1, decayFactor)
}
