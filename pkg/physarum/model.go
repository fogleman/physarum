package physarum

import (
	"math"
	"math/rand"
	"runtime"
	"sync"
)

type Model struct {
	W int
	H int

	BlurRadius int
	BlurPasses int

	Configs   []Config
	Grids     []*Grid
	Particles []Particle

	Iteration int
}

func NewModel(w, h, numParticles, blurRadius, blurPasses int, configs []Config) *Model {
	grids := make([]*Grid, len(configs))
	numParticlesPerConfig := int(math.Ceil(
		float64(numParticles) / float64(len(configs))))
	actualNumParticles := numParticlesPerConfig * len(configs)
	particles := make([]Particle, 0, actualNumParticles)
	for c := range configs {
		grids[c] = NewGrid(w, h)
		for i := 0; i < numParticlesPerConfig; i++ {
			x := rand.Float32() * float32(w)
			y := rand.Float32() * float32(h)
			a := rand.Float32() * 2 * math.Pi
			p := Particle{x, y, a, uint32(c)}
			particles = append(particles, p)
		}
	}
	return &Model{w, h, blurRadius, blurPasses, configs, grids, particles, 0}
}

func (m *Model) Step() {
	updateParticle := func(rnd *rand.Rand, i int) {
		p := m.Particles[i]
		config := m.Configs[p.C]
		grid := m.Grids[p.C]

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
		C := grid.GetTemp(xc, yc)
		L := grid.GetTemp(xl, yl)
		R := grid.GetTemp(xr, yr)

		da := rotationAngle * direction(rnd, C, L, R)
		// da := rotationAngle * weightedDirection(rnd, C, L, R)
		p.A = Shift(p.A+da, 2*math.Pi)
		p.X = Shift(p.X+cos(p.A)*stepDistance, float32(m.W))
		p.Y = Shift(p.Y+sin(p.A)*stepDistance, float32(m.H))
		m.Particles[i] = p
	}

	updateParticles := func(wi, wn int, wg *sync.WaitGroup) {
		seed := int64(m.Iteration)<<8 | int64(wi)
		rnd := rand.New(rand.NewSource(seed))
		n := len(m.Particles)
		batch := int(math.Ceil(float64(n) / float64(wn)))
		i0 := wi * batch
		i1 := i0 + batch
		if wi == wn-1 {
			i1 = n
		}
		for i := i0; i < i1; i++ {
			updateParticle(rnd, i)
		}
		wg.Done()
	}

	updateGrids := func(c int, wg *sync.WaitGroup) {
		config := m.Configs[c]
		grid := m.Grids[c]
		for _, p := range m.Particles {
			if uint32(c) == p.C {
				grid.Add(p.X, p.Y, config.DepositionAmount)
			}
		}
		grid.BoxBlur(m.BlurRadius, m.BlurPasses, config.DecayFactor)
		wg.Done()
	}

	combineGrids := func(c int, wg *sync.WaitGroup) {
		config := m.Configs[c]
		grid := m.Grids[c]
		repulsionFactor := config.RepulsionFactor
		copy(grid.Temp, grid.Data)
		for i, other := range m.Grids {
			if i == c {
				continue
			}
			for j, value := range other.Data {
				grid.Temp[j] -= value * repulsionFactor
			}
		}
		wg.Done()
	}

	var wg sync.WaitGroup

	// step 1: combine grids
	for i := range m.Configs {
		wg.Add(1)
		go combineGrids(i, &wg)
	}
	wg.Wait()

	// step 2: move particles
	wn := runtime.NumCPU()
	for wi := 0; wi < wn; wi++ {
		wg.Add(1)
		go updateParticles(wi, wn, &wg)
	}
	wg.Wait()

	// step 3: deposit, blur, and decay
	for i := range m.Configs {
		wg.Add(1)
		go updateGrids(i, &wg)
	}
	wg.Wait()

	m.Iteration++
}

func (m *Model) Data() [][]float32 {
	result := make([][]float32, len(m.Grids))
	for i, grid := range m.Grids {
		result[i] = make([]float32, len(grid.Data))
		copy(result[i], grid.Data)
	}
	return result
}

func direction(rnd *rand.Rand, C, L, R float32) float32 {
	if C > L && C > R {
		return 0
	} else if C < L && C < R {
		return float32((rnd.Int63()&1)<<1 - 1)
	} else if L < R {
		return 1
	} else if R < L {
		return -1
	}
	return 0
}

func weightedDirection(rnd *rand.Rand, C, L, R float32) float32 {
	W := [3]float32{C, L, R}
	D := [3]float32{0, -1, 1}

	if W[0] > W[1] {
		W[0], W[1] = W[1], W[0]
		D[0], D[1] = D[1], D[0]
	}
	if W[0] > W[2] {
		W[0], W[2] = W[2], W[0]
		D[0], D[2] = D[2], D[0]
	}
	if W[1] > W[2] {
		W[1], W[2] = W[2], W[1]
		D[1], D[2] = D[2], D[1]
	}

	a := W[1] - W[0]
	b := W[2] - W[1]
	if rnd.Float32()*(a+b) < a {
		return D[1]
	}
	return D[2]
}
