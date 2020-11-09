package physarum

import (
	"fmt"
)

func one(model *Model, iterations int) {
	for i := 0; i < iterations; i++ {
		fmt.Println(i)
		model.Step()
	}
	SavePNG("out.png", Image(model.W, model.H, model.Colors()))
}

func frames(model *Model, rate int) {
	saveImage := func(path string, w, h int, colors [][]float64, ch chan bool) {
		im := Image(w, h, colors)
		SavePNG(path, im)
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
			path := fmt.Sprintf("out%08d.png", i/rate)
			go saveImage(path, model.W, model.H, model.Colors(), ch)
		}
	}
}

func Run() {
	model := NewModel(1920, 1080, DefaultConfigs)
	// one(model, 100)
	frames(model, 5)
}
