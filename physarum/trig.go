package physarum

import (
	"log"
	"math"
)

const (
	trigTableSize = 65536
	trigTableMask = trigTableSize - 1
	trigFactor    = trigTableSize / (2 * math.Pi)
)

var (
	sinTable []float32
	cosTable []float32
)

func init() {
	if !IsPowerOfTwo(trigTableSize) {
		log.Fatal("trigTableSize must be a power of two")
	}
	sinTable = make([]float32, trigTableSize)
	cosTable = make([]float32, trigTableSize)
	for i := range sinTable {
		t := float64(i) / trigTableSize
		a := t * 2 * math.Pi
		sinTable[i] = float32(math.Sin(a))
		cosTable[i] = float32(math.Cos(a))
	}
}

func sin(t float32) float32 {
	i := int(t*trigFactor+trigTableSize) & trigTableMask
	return sinTable[i]
}

func cos(t float32) float32 {
	i := int(t*trigFactor+trigTableSize) & trigTableMask
	return cosTable[i]
}
