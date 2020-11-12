package physarum

import (
	"math"
	"math/rand"
	"runtime"
)

type Model struct {
	W          int
	H          int
	Iterations int
	Configs    []*Config
	Grids      []*Grid
	Particles  []Particle
}

func NewModel(w, h int, configs []*Config) *Model {
	grids := make([]*Grid, len(configs))
	var particles []Particle
	for c, config := range configs {
		grids[c] = NewGrid(w, h)
		numParticles := int(float32(w*h) * config.PopulationPercentage)
		for i := 0; i < numParticles; i++ {
			x := rand.Float32() * float32(w)
			y := rand.Float32() * float32(h)
			a := rand.Float32() * 2 * math.Pi
			p := Particle{x, y, a, uint32(c)}
			particles = append(particles, p)
		}
	}
	return &Model{w, h, 0, configs, grids, particles}
}

func (m *Model) Step() {
	updateParticle := func(rnd *rand.Rand, i int) {
		p := m.Particles[i]
		config := m.Configs[p.C]

		// u := p.X / float32(m.W)
		// v := p.Y / float32(m.H)

		sensorDistance := config.SensorDistance
		sensorAngle := config.SensorAngle
		rotationAngle := config.RotationAngle
		stepDistance := config.StepDistance

		xc := p.X + cos(p.A)*sensorDistance
		yc := p.Y + sin(p.A)*sensorDistance
		xl := p.X + cos(p.A-sensorAngle)*sensorDistance
		yl := p.Y + sin(p.A-sensorAngle)*sensorDistance
		xr := p.X + cos(p.A+sensorAngle)*sensorDistance
		yr := p.Y + sin(p.A+sensorAngle)*sensorDistance
		var C, L, R float32
		for c, grid := range m.Grids {
			if uint32(c) == p.C {
				C += grid.Get(xc, yc)
				L += grid.Get(xl, yl)
				R += grid.Get(xr, yr)
			} else {
				C -= grid.Get(xc, yc)
				L -= grid.Get(xl, yl)
				R -= grid.Get(xr, yr)
			}
		}
		var da float32
		if C > L && C > R {
			// straight
		} else if C < L && C < R {
			// rotate randomly left or right
			if rnd.Intn(2) == 0 {
				da = rotationAngle
			} else {
				da = -rotationAngle
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
		p.A = Shift(p.A+da, 2*math.Pi)
		p.X = Shift(p.X+cos(p.A)*stepDistance, float32(m.W))
		p.Y = Shift(p.Y+sin(p.A)*stepDistance, float32(m.H))
		m.Particles[i] = p
	}

	updateParticles := func(wi, wn int, ch chan bool) {
		// seed := time.Now().UTC().UnixNano()
		seed := int64(m.Iterations)<<8 | int64(wi)
		rnd := rand.New(rand.NewSource(seed))
		n := len(m.Particles)
		for i := wi; i < n; i += wn {
			updateParticle(rnd, i)
		}
		ch <- true
	}

	updateGrids := func(c int, ch chan bool) {
		config := m.Configs[c]
		grid := m.Grids[c]
		for _, p := range m.Particles {
			if uint32(c) == p.C {
				grid.Add(p.X, p.Y, config.DepositionAmount)
			}
		}
		grid.BoxBlur(1, 1, config.DecayFactor)
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

	wn = len(m.Configs)
	for wi := 0; wi < wn; wi++ {
		go updateGrids(wi, ch)
	}
	for wi := 0; wi < wn; wi++ {
		<-ch
	}

	m.Iterations++
}

func (m *Model) Data() [][]float32 {
	result := make([][]float32, len(m.Grids))
	for i, grid := range m.Grids {
		result[i] = make([]float32, len(grid.Data))
		copy(result[i], grid.Data)
	}
	return result
}
