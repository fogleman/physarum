package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"runtime"
	"sort"
	"time"

	"github.com/fogleman/physarum/pkg/physarum"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/gonum/stat"
)

const (
	width      = 512
	height     = 512
	particles  = 1 << 20
	blurRadius = 1
	blurPasses = 2
	zoomFactor = 1
	scale      = 1
	gamma      = 1 / 2.2
	title      = "physarum"
)

var Configs = []physarum.Config{
	// cyclones
	// 	{4, 0.87946403, 42.838207, 0.97047323, 2.8447638, 5, 0.29681, 1.4512},
	// 	{4, 1.7357124, 17.430664, 0.30490428, 2.1706762, 5, 0.27878627, 0.46232897},

	// dunes
	// {2, 0.99931663, 44.21652, 1.9704952, 1.4215798, 5, 0.1580779, 0.7574965},
	// {2, 1.9694986, 1.294038, 0.5384646, 1.1613986, 5, 0.21102181, 1.5123861},

	// dot grid
	// {1.3333334, 1.3433642, 49.39263, 0.91616887, 0.69644034, 5, 0.17888786, 0.2036435},
	// {1.3333334, 0.0856143, 1.6695175, 1.8827246, 2.3155663, 5, 0.14249614, 0.0026361942},
	// {1.3333334, 0.7959472, 33.977413, 0.5246451, 2.2891424, 5, 0.22549233, 1.4248871},

	// untitled
	// {1.7433162, 56.586433, 0.45428953, 0.78228176, 5, 0.19172272, 1.6682954},
	// {1.8340914, 1.6538872, 1.4098115, 1.6714363, 5, 0.17746642, 1.491355},
	// {0.0049473564, 13.269191, 0.033447478, 1.0102618, 5, 0.2197167, 1.6166985},
	// {0.37645763, 31.045816, 0.81319964, 2.5322618, 5, 0.10834738, 1.5170672},
	// {0.7355474, 14.832715, 0.2641479, 0.8953786, 5, 0.14977153, 0.14632958},

	// cool
	// {1.4107815, 61.27741, 0.49201587, 1.3007548, 5, 0.49895996, 1.0823951},
	// {1.1534524, 13.299458, 0.48315683, 1.8219115, 5, 0.41845483, 0.4055887},
	// {0.31089303, 60.62575, 1.0241486, 0.39942655, 5, 0.4576149, 0.24079543},
	// {0.40245488, 27.844227, 1.9592205, 0.5504824, 5, 0.19568197, 1.1694417},
	// {1.227412, 1.7987814, 0.39546785, 1.2640203, 5, 0.14201605, 0.77068233},
}

func init() {
	runtime.LockOSThread()
}

type Texture struct {
	w   int
	h   int
	id  uint32
	buf []uint8
	acc []float32
	r   [][]float32
	g   [][]float32
	b   [][]float32
	min []float32
	max []float32
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

func (t *Texture) Init(width, height, count int) {
	const N = 65536
	t.w = width
	t.h = height
	t.buf = make([]uint8, width*height*3)
	t.acc = make([]float32, width*height*3)
	t.r = make([][]float32, count)
	t.g = make([][]float32, count)
	t.b = make([][]float32, count)
	for i := 0; i < count; i++ {
		t.r[i] = make([]float32, N)
		t.g[i] = make([]float32, N)
		t.b[i] = make([]float32, N)
	}
	max := particles / float32(width*height) * 10
	t.min = make([]float32, count)
	t.max = make([]float32, count)
	for i := range t.min {
		t.min[i] = 0
		t.max[i] = max
	}
}

func (t *Texture) SetPalette(palette physarum.Palette, gamma float32) {
	count := len(t.r)
	N := len(t.r[0])
	for i := 0; i < count; i++ {
		c := palette[i]
		for j := 0; j < N; j++ {
			p := float32(j) / float32(N-1)
			p = float32(math.Pow(float64(p), float64(gamma)))
			t.r[i][j] = float32(c.R) * p
			t.g[i][j] = float32(c.G) * p
			t.b[i][j] = float32(c.B) * p
		}
	}
	palette.Print()
}

func (t *Texture) ShufflePalette() {
	rand.Shuffle(len(t.r), func(i, j int) {
		t.r[i], t.r[j] = t.r[j], t.r[i]
		t.g[i], t.g[j] = t.g[j], t.g[i]
		t.b[i], t.b[j] = t.b[j], t.b[i]
	})
}

func (t *Texture) AutoLevel(data [][]float32, minPercentile, maxPercentile float64) {
	for i, grid := range data {
		temp := make([]float64, len(grid))
		for j, v := range grid {
			temp[j] = float64(v)
		}
		sort.Float64s(temp)
		t.min[i] = float32(stat.Quantile(minPercentile, stat.Empirical, temp, nil))
		t.max[i] = float32(stat.Quantile(maxPercentile, stat.Empirical, temp, nil))
	}
}

func (t *Texture) update(data [][]float32) {
	for i := range t.acc {
		t.acc[i] = 0
	}
	f := float32(len(t.r[0]) - 1)
	for i, grid := range data {
		min, max := t.min[i], t.max[i]
		m := 1 / float32(max-min)
		for j, value := range grid {
			p := (value - min) * m
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
	s1 := float32(w) / float32(t.w)
	s2 := float32(h) / float32(t.h)
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
		gl.TEXTURE_2D, 0, gl.RGB, int32(t.w), int32(t.h),
		0, gl.RGB, gl.UNSIGNED_BYTE, gl.Ptr(t.buf))
	t.draw(window)
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

func makeModel() *physarum.Model {
	configs := physarum.RandomConfigs(2 + rand.Intn(4))
	if len(Configs) > 0 {
		configs = Configs
	}
	table := physarum.RandomAttractionTable(len(configs))
	model := physarum.NewModel(
		width, height, particles, blurRadius, blurPasses, zoomFactor,
		configs, table, "random")
	physarum.PrintConfigs(model.Configs, model.AttractionTable)
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

	var model *physarum.Model
	texture := NewTexture()

	reset := func() {
		model = makeModel()
		texture.Init(model.W, model.H, len(model.Configs))
		texture.SetPalette(physarum.RandomPalette(), gamma)
	}

	reset()

	window.SetKeyCallback(func(window *glfw.Window, key glfw.Key, code int, action glfw.Action, mods glfw.ModifierKey) {
		if action == glfw.Press {
			if key == glfw.KeySpace {
				reset()
			}
			if key == glfw.KeyR {
				model.StartOver()
			}
			if key == glfw.KeyP {
				texture.SetPalette(physarum.RandomPalette(), gamma)
			}
			if key == glfw.KeyO {
				texture.ShufflePalette()
			}
			if key == glfw.KeyA {
				texture.AutoLevel(model.Data(), 0.001, 0.999)
			}
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
