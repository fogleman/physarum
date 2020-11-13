package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"runtime"
	"time"

	"github.com/fogleman/physarum/pkg/physarum"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

const (
	width  = 512
	height = 512
	pct    = 1
	scale  = 1
	gamma  = 1 / 2.2
	title  = "physarum"
)

func init() {
	runtime.LockOSThread()
}

type Texture struct {
	id  uint32
	buf []uint8
	acc []float32
	r   [][]float32
	g   [][]float32
	b   [][]float32
}

func NewTexture() *Texture {
	var id uint32
	gl.GenTextures(1, &id)
	gl.BindTexture(gl.TEXTURE_2D, id)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.BindTexture(gl.TEXTURE_2D, 0)
	return &Texture{id: id}
}

func (t *Texture) Init(width, height, count int, palette physarum.Palette, gamma float32) {
	const N = 256
	t.buf = make([]uint8, width*height*3)
	t.acc = make([]float32, width*height*3)
	t.r = make([][]float32, count)
	t.g = make([][]float32, count)
	t.b = make([][]float32, count)
	for i := 0; i < count; i++ {
		t.r[i] = make([]float32, N)
		t.g[i] = make([]float32, N)
		t.b[i] = make([]float32, N)
		c := palette[i]
		for j := 0; j < N; j++ {
			p := float32(j) / (N - 1)
			p = float32(math.Pow(float64(p), float64(gamma)))
			t.r[i][j] = float32(c.R) * p
			t.g[i][j] = float32(c.G) * p
			t.b[i][j] = float32(c.B) * p
		}
	}
}

func (t *Texture) update(data [][]float32) {
	const lo = 0
	const hi = 20
	for i := range t.acc {
		t.acc[i] = 0
	}
	f := float32(len(t.r[0]) - 1)
	m := 1 / float32(hi-lo)
	for i, grid := range data {
		for j, value := range grid {
			p := (value - lo) * m
			if p < 0 {
				p = 0
			}
			if p > 1 {
				p = 1
			}
			index := int(p * f)
			t.acc[j*3+0] += t.r[i][index]
			t.acc[j*3+1] += t.g[i][index]
			t.acc[j*3+2] += t.b[i][index]
		}
	}
	for i, value := range t.acc {
		if value > 255 {
			value = 255
		}
		t.buf[i] = uint8(value)
	}
}

func (t *Texture) draw(window *glfw.Window) {
	const padding = 0
	w, h := window.GetFramebufferSize()
	s1 := float32(w) / width
	s2 := float32(h) / height
	f := float32(1 - padding)
	var x, y float32
	if s1 >= s2 {
		x = f * s2 / s1
		y = f
	} else {
		x = f
		y = f * s1 / s2
	}
	gl.Begin(gl.QUADS)
	gl.TexCoord2f(0, 1)
	gl.Vertex2f(-x, -y)
	gl.TexCoord2f(1, 1)
	gl.Vertex2f(x, -y)
	gl.TexCoord2f(1, 0)
	gl.Vertex2f(x, y)
	gl.TexCoord2f(0, 0)
	gl.Vertex2f(-x, y)
	gl.End()
}

func (t *Texture) Draw(window *glfw.Window, data [][]float32) {
	t.update(data)
	gl.BindTexture(gl.TEXTURE_2D, t.id)
	gl.TexImage2D(
		gl.TEXTURE_2D, 0, gl.RGB, width, height,
		0, gl.RGB, gl.UNSIGNED_BYTE, gl.Ptr(t.buf))
	t.draw(window)
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

func makeModel() *physarum.Model {
	count := 2 + rand.Intn(4)
	configs := physarum.RandomConfigs(count)
	for i := range configs {
		configs[i].PopulationPercentage = pct
	}
	model := physarum.NewModel(width, height, configs)
	physarum.PrintConfigs(model.Configs)
	physarum.SummarizeConfigs(model.Configs)
	fmt.Println()
	return model
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	// initialize glfw
	if err := glfw.Init(); err != nil {
		log.Fatal(err)
	}
	defer glfw.Terminate()

	// create window
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	window, err := glfw.CreateWindow(
		width*scale, height*scale, title, nil, nil)
	if err != nil {
		log.Fatalln(err)
	}
	window.MakeContextCurrent()

	// initialize gl
	if err := gl.Init(); err != nil {
		log.Fatalln(err)
	}
	gl.Enable(gl.TEXTURE_2D)

	texture := NewTexture()

	var model *physarum.Model

	reset := func() {
		model = makeModel()
		texture.Init(
			width, height, len(model.Configs),
			physarum.RandomPalette(), gamma)
	}

	reset()

	window.SetKeyCallback(func(window *glfw.Window, key glfw.Key, code int, action glfw.Action, mods glfw.ModifierKey) {
		if action == glfw.Press && key == glfw.KeySpace {
			reset()
		}
	})

	// main loop
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT)
		model.Step()
		texture.Draw(window, model.Data())
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
