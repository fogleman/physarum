package physarum

import (
	"fmt"
	"math"
	"testing"
)

func TestTrigLookupTables(t *testing.T) {
	const N = 10000000
	const tolerance = 1e-4
	var maxCosError, maxSinError float64
	for i := 0; i <= N; i++ {
		p := float32(i) / N
		a := (p*20 - 10) * math.Pi
		cosError := math.Abs(float64(cos(a)) - math.Cos(float64(a)))
		sinError := math.Abs(float64(sin(a)) - math.Sin(float64(a)))
		maxCosError = math.Max(maxCosError, cosError)
		maxSinError = math.Max(maxSinError, sinError)
		if cosError > tolerance {
			t.Fatalf("cos(%v) = %v, math.Cos(%v) = %v (%g)", a, cos(a), a, math.Cos(float64(a)), cosError)
		}
		if sinError > tolerance {
			t.Fatalf("sin(%v) = %v, math.Sin(%v) = %v (%g)", a, sin(a), a, math.Sin(float64(a)), sinError)
		}
	}
	fmt.Println(maxCosError)
	fmt.Println(maxSinError)
}
