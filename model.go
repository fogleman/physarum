package physarum

import (
	"math"
	"math/rand"
	"runtime"
	"time"
)

type Model struct {
	W         int
	H         int
	Configs   []*Config
	Grids     []*Grid
	Particles []Particle
}

func NewModel(w, h int, configs []*Config) *Model {
	grids := make([]*Grid, len(configs))
	var particles []Particle
	for c, config := range configs {
		grids[c] = NewGrid(w, h)
		numParticles := int(float64(w*h) * config.PopulationPercentage)
		for i := 0; i < numParticles; i++ {
			x := rand.Float64() * float64(w)
			y := rand.Float64() * float64(h)
			a := rand.Float64() * 2 * math.Pi
			p := Particle{x, y, a, c}
			particles = append(particles, p)
		}
	}
	return &Model{w, h, configs, grids, particles}
}

func (m *Model) Step() {
	updateParticle := func(rnd *rand.Rand, i int) {
		p := m.Particles[i]
		config := m.Configs[p.C]

		// u := p.X / float64(m.W)
		// v := p.Y / float64(m.H)

		sensorDistance := config.SensorDistance
		sensorAngle := config.SensorAngle
		rotationAngle := config.RotationAngle
		stepDistance := config.StepDistance

		// sensorAngle = (u + 0.5) * math.Pi * 2
		// // sensorDistance = v * 64
		// rotationAngle = (v + 0.5) * math.Pi * 2

		// // rotationAngle *= (1 + rnd.NormFloat64()*0.2)
		// // stepDistance *= (1 + rnd.NormFloat64()*0.2)
		// // sensorAngle *= (1 + rnd.NormFloat64()*0.05)
		// // rotationAngle *= (1 + rnd.NormFloat64()*0.05)
		// sensorAngle += rnd.NormFloat64() * gg.Radians(15)
		// rotationAngle += rnd.NormFloat64() * gg.Radians(15)

		xc := p.X + math.Cos(p.A)*sensorDistance
		yc := p.Y + math.Sin(p.A)*sensorDistance
		xl := p.X + math.Cos(p.A-sensorAngle)*sensorDistance
		yl := p.Y + math.Sin(p.A-sensorAngle)*sensorDistance
		xr := p.X + math.Cos(p.A+sensorAngle)*sensorDistance
		yr := p.Y + math.Sin(p.A+sensorAngle)*sensorDistance
		C := m.Grids[p.C].Get(xc, yc)
		L := m.Grids[p.C].Get(xl, yl)
		R := m.Grids[p.C].Get(xr, yr)
		for c, grid := range m.Grids {
			if c != p.C {
				C -= grid.Get(xc, yc) * 1
				L -= grid.Get(xl, yl) * 1
				R -= grid.Get(xr, yr) * 1
			}
		}
		var da float64
		if C > L && C > R {
			// straight
		} else if C < L && C < R {
			// rotate randomly left or right
			da = rotationAngle
			if rnd.Intn(2) == 0 {
				da = -da
			}
		} else if L < R {
			// rotate right
			da = rotationAngle
		} else if R < L {
			// rotate left
			da = -rotationAngle
		} else {
			// straight
		}
		p.A += da
		p.X += math.Cos(p.A) * stepDistance
		p.Y += math.Sin(p.A) * stepDistance
		if p.X < 0 {
			p.X += float64(m.W)
		}
		if p.Y < 0 {
			p.Y += float64(m.H)
		}
		if p.X >= float64(m.W) {
			p.X -= float64(m.W)
		}
		if p.Y >= float64(m.H) {
			p.Y -= float64(m.H)
		}
		m.Particles[i] = p
	}

	updateParticles := func(wi, wn int, ch chan bool) {
		rnd := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
		n := len(m.Particles)
		for i := wi; i < n; i += wn {
			updateParticle(rnd, i)
		}
		ch <- true
	}

	wn := runtime.NumCPU()
	ch := make(chan bool, wn)
	for wi := 0; wi < wn; wi++ {
		go updateParticles(wi, wn, ch)
	}
	for wi := 0; wi < wn; wi++ {
		<-ch
	}

	for _, p := range m.Particles {
		config := m.Configs[p.C]
		m.Grids[p.C].Add(p.X, p.Y, config.DepositionAmount)
	}

	for c, config := range m.Configs {
		m.Grids[c].DiffuseAndDecay(config.DecayFactor)
	}
}

// func (m *Model) Image() image.Image {
// 	im := image.NewRGBA(image.Rect(0, 0, m.W, m.H))
// 	for y := 0; y < m.H; y++ {
// 		for x := 0; x < m.W; x++ {
// 			index := y*m.W + x
// 			var r, g, b float64
// 			for i, grid := range m.Grids {
// 				t := grid.Data[index] / 5
// 				if t > 1 {
// 					t = 1
// 				}
// 				t = math.Pow(t, 1/2.2)
// 				switch i {
// 				case 0:
// 					r = t
// 				case 1:
// 					b = t
// 				case 2:
// 					g = t
// 				}
// 			}
// 			c := color.RGBA{
// 				uint8(r * 255),
// 				uint8(g * 255),
// 				uint8(b * 255),
// 				255,
// 			}
// 			im.SetRGBA(x, y, c)
// 		}
// 	}
// 	return im
// }

func (m *Model) Colors() [][]float64 {
	result := make([][]float64, len(m.Grids))
	for i, grid := range m.Grids {
		result[i] = make([]float64, len(grid.Data))
		copy(result[i], grid.Data)
	}
	return result
}
