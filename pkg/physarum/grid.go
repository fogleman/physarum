package physarum

import (
	"log"
	"math/rand"
)

type Grid struct {
	W    int
	H    int
	Data []float32
	Temp []float32
}

func NewGrid(w, h int) *Grid {
	if !IsPowerOfTwo(w) || !IsPowerOfTwo(h) {
		log.Fatal("grid dimensions must be a power of two")
	}
	data := make([]float32, w*h)
	temp := make([]float32, w*h)
	for i := range data {
		data[i] = rand.Float32()
	}
	return &Grid{w, h, data, temp}
}

func (g *Grid) Get(x, y float32) float32 {
	i := int(x+float32(g.W)) & (g.W - 1)
	j := int(y+float32(g.H)) & (g.H - 1)
	index := j*g.W + i
	return g.Data[index]
}

func (g *Grid) GetTemp(x, y float32) float32 {
	i := int(x+float32(g.W)) & (g.W - 1)
	j := int(y+float32(g.H)) & (g.H - 1)
	index := j*g.W + i
	return g.Temp[index]
}

func (g *Grid) Add(x, y, a float32) {
	x += float32(g.W)
	y += float32(g.H)
	x0 := int(x)
	y0 := int(y)
	x1 := x0 + 1
	y1 := y0 + 1
	u := x - float32(x0)
	v := y - float32(y0)
	x0 &= (g.W - 1)
	y0 &= (g.H - 1)
	x1 &= (g.W - 1)
	y1 &= (g.H - 1)
	g.Data[x0+y0*g.W] += a * ((1 - u) * (1 - v))
	g.Data[x0+y1*g.W] += a * ((1 - u) * v)
	g.Data[x1+y0*g.W] += a * (u * (1 - v))
	g.Data[x1+y1*g.W] += a * (u * v)
}

func (g *Grid) BoxBlur(radius, iterations int, decayFactor float32) {
	for i := 1; i < iterations; i++ {
		boxBlur(g.Data, g.Temp, g.W, g.H, radius, 1)
	}
	boxBlur(g.Data, g.Temp, g.W, g.H, radius, decayFactor)
}
