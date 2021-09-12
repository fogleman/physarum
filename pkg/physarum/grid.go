package physarum

import (
	"log"
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
		data[i] = 0.0 //01 //rand.Float32()
	}
	return &Grid{w, h, data, temp}
}

func (g *Grid) Index(x, y float32) int {
	i := int(x+float32(g.W)) & (g.W - 1)
	j := int(y+float32(g.H)) & (g.H - 1)
	return j*g.W + i
}

func (g *Grid) Get(x, y float32) float32 {
	return g.Data[g.Index(x, y)]
}

func (g *Grid) GetTemp(x, y float32) float32 {
	return g.Temp[g.Index(x, y)]
}

func (g *Grid) Add(x, y, a float32) {
	g.Data[g.Index(x, y)] += a
}

func (g *Grid) BoxBlur(radius, iterations int, decayFactor float32) {
	if iterations < 1 {
		for i := range g.Data {
			g.Data[i] *= decayFactor
		}
		return
	}
	for i := 1; i < iterations; i++ {
		boxBlur(g.Data, g.Temp, g.W, g.H, radius, 1)
	}
	boxBlur(g.Data, g.Temp, g.W, g.H, radius, decayFactor)
}
