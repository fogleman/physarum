package physarum

import (
	"fmt"
	"image/png"
	"math/rand"
	"path/filepath"
	"time"
)

// const (
// 	width      = 1024
// 	height     = 1024
// 	particles  = 1 << 22
// 	iterations = 400
// 	blurRadius = 1
// 	blurPasses = 2
// 	zoomFactor = 1
// )

const (
	width      = 4096
	height     = 2048
	particles  = 1 << 25
	iterations = 400
	blurRadius = 1
	blurPasses = 2
	zoomFactor = 1
)

func one(model *Model, iterations int) {
	now := time.Now().UTC().UnixNano() / 1000
	file := fmt.Sprintf("out%d.png", now)
	fmt.Println()
	fmt.Println(file)
	fmt.Println(len(model.Particles), "particles")
	PrintConfigs(model.Configs, model.AttractionTable)
	SummarizeConfigs(model.Configs)
	for i := 0; i < iterations; i++ {
		model.Step()
	}
	palette := RandomPalette()
	im := Image(model.W, model.H, model.Data(), palette, 0, 0, 1/2.2)
	SavePNG(".", file, im, png.DefaultCompression)
}

func frames(model *Model, rate int) {
	palette := RandomPalette()

	saveImage := func(path string, file string, w, h int, grids [][]float32) { //, ch chan bool) {
		max := particles / float32(width*height) * 20
		im := Image(w, h, grids, palette, 0, max, 1/2.2)
		SavePNG(path, file, im, png.BestSpeed)
		// if ch != nil {
		// 	ch <- true
		// }
	}

	now := time.Now().UTC().UnixNano() / 1000
	path := filepath.Join(".", "output", fmt.Sprintf("%d", now))

	fmt.Println(len(model.Particles), "particles")
	PrintConfigs(model.Configs, model.AttractionTable)
	SummarizeConfigs(model.Configs)

	// ch := make(chan bool, 1)
	// ch <- true
	for i := 0; ; i++ {
		if i%1000 == 0 {
			fmt.Println(i)
		}
		model.Step()
		if i%rate == 0 {
			// <-ch
			file := fmt.Sprintf("frame%08d.png", i/rate)
			// fmt.Println(path + " " + file)
			go saveImage(path, file, model.W, model.H, model.Data()) //, ch)
		}
	}
}

func Run() {
	if true {
		// n := 2 + rand.Intn(4)
		n := 1
		configs := RandomConfigs(n)
		table := RandomAttractionTable(n)
		model := NewModel(
			width, height, particles, blurRadius, blurPasses, zoomFactor,
			configs, table, "point")
		frames(model, 3)
	}

	for {
		n := 2 + rand.Intn(4)
		configs := RandomConfigs(n)
		table := RandomAttractionTable(n)
		model := NewModel(
			width, height, particles, blurRadius, blurPasses, zoomFactor,
			configs, table, "random")
		start := time.Now()
		one(model, iterations)
		fmt.Println(time.Since(start))
	}
}
