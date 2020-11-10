package physarum

import "math"

const (
	trigTableSize = 4096
	trigFactor    = trigTableSize / (2 * math.Pi)
)

var (
	sinTable []float64
	cosTable []float64
)

func init() {
	sinTable = make([]float64, trigTableSize)
	cosTable = make([]float64, trigTableSize)
	for i := range sinTable {
		t := float64(i) / trigTableSize
		a := t * 2 * math.Pi
		sinTable[i] = math.Sin(a)
		cosTable[i] = math.Cos(a)
	}
}

func sin(t float64) float64 {
	i := int(t*trigFactor+trigTableSize) & (trigTableSize - 1)
	return sinTable[i]
}

func cos(t float64) float64 {
	i := int(t*trigFactor+trigTableSize) & (trigTableSize - 1)
	return cosTable[i]
}
