package physarum

import (
	"fmt"
	"image/png"
	"math/rand"
	"time"
)

const (
	width      = 1024
	height     = 1024
	particles  = 1 << 20
	iterations = 500
	blurRadius = 1
	blurPasses = 2
	zoomFactor = 1
)

func one(model *Model, iterations int) {
	now := time.Now().UTC().UnixNano() / 1000
	path := fmt.Sprintf("out%d.png", now)
	fmt.Println()
	fmt.Println(path)
	fmt.Println(len(model.Particles), "particles")
	PrintConfigs(model.Configs)
	SummarizeConfigs(model.Configs)
	for i := 0; i < iterations; i++ {
		model.Step()
	}
	palette := RandomPalette()
	im := Image(model.W, model.H, model.Data(), palette, 0, 0, 1/2.2)
	SavePNG(path, im, png.DefaultCompression)
}

func frames(model *Model, rate int) {
	palette := RandomPalette()

	saveImage := func(path string, w, h int, grids [][]float32, ch chan bool) {
		max := particles / float32(width*height) * 20
		im := Image(w, h, grids, palette, 0, max, 1/2.2)
		SavePNG(path, im, png.BestSpeed)
		if ch != nil {
			ch <- true
		}
	}

	ch := make(chan bool, 1)
	ch <- true
	for i := 0; ; i++ {
		fmt.Println(i)
		model.Step()
		if i%rate == 0 {
			<-ch
			path := fmt.Sprintf("frame%08d.png", i/rate)
			go saveImage(path, model.W, model.H, model.Data(), ch)
		}
	}
}

func Run() {
	if false {
		n := 2 + rand.Intn(4)
		configs := RandomConfigs(n)
		model := NewModel(
			width, height, particles, blurRadius, blurPasses, zoomFactor, configs)
		frames(model, 3)
	}

	for {
		n := 2 + rand.Intn(4)
		configs := RandomConfigs(n)
		model := NewModel(
			width, height, particles, blurRadius, blurPasses, zoomFactor, configs)
		start := time.Now()
		one(model, iterations)
		fmt.Println(time.Since(start))
	}
}
