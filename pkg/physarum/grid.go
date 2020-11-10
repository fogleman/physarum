package physarum

import (
	"log"
	"math/rand"
)

type Grid struct {
	W    int
	H    int
	Data []float64
	Temp []float64
}

func NewGrid(w, h int) *Grid {
	if !IsPowerOfTwo(w) || !IsPowerOfTwo(h) {
		log.Fatal("grid dimensions must be a power of two")
	}
	data := make([]float64, w*h)
	temp := make([]float64, w*h)
	for i := range data {
		data[i] = rand.Float64()
	}
	return &Grid{w, h, data, temp}
}

func (g *Grid) Get(x, y float64) float64 {
	i := int(x+float64(g.W)) & (g.W - 1)
	j := int(y+float64(g.H)) & (g.H - 1)
	index := j*g.W + i
	return g.Data[index]
}

func (g *Grid) Add(x, y, a float64) {
	x += float64(g.W)
	y += float64(g.H)
	x0 := int(x)
	y0 := int(y)
	x1 := x0 + 1
	y1 := y0 + 1
	u := x - float64(x0)
	v := y - float64(y0)
	x0 &= (g.W - 1)
	y0 &= (g.H - 1)
	x1 &= (g.W - 1)
	y1 &= (g.H - 1)
	g.Data[x0+y0*g.W] += a * ((1 - u) * (1 - v))
	g.Data[x0+y1*g.W] += a * ((1 - u) * v)
	g.Data[x1+y0*g.W] += a * (u * (1 - v))
	g.Data[x1+y1*g.W] += a * (u * v)
}

func (g *Grid) BoxBlur(radius, iterations int, decayFactor float64) {
	for i := 1; i < iterations; i++ {
		boxBlur(g.Data, g.Temp, g.W, g.H, radius, 1)
	}
	boxBlur(g.Data, g.Temp, g.W, g.H, radius, decayFactor)
}

func (g *Grid) GaussianBlur(radius int, decayFactor float64) {
	gaussianBlur(g.Data, g.Temp, g.W, g.H, radius, decayFactor)
}
