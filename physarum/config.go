package physarum

import (
	"fmt"
	"math/rand"
)

const (
	sensorAngleMin      = 0
	sensorAngleMax      = 120
	sensorDistanceMin   = 0
	sensorDistanceMax   = 64
	rotationAngleMin    = 0
	rotationAngleMax    = 120
	stepDistanceMin     = 0.2
	stepDistanceMax     = 2
	depositionAmountMin = 5
	depositionAmountMax = 5
	decayFactorMin      = 0.1
	decayFactorMax      = 0.1

	attractionFactorMean = 1
	attractionFactorStd  = 0.25
	repulsionFactorMean  = -1
	repulsionFactorStd   = 0.25
)

type Config struct {
	SensorAngle      float32
	SensorDistance   float32
	RotationAngle    float32
	StepDistance     float32
	DepositionAmount float32
	DecayFactor      float32
}

func RandomConfig() Config {
	uniform := func(min, max float32) float32 {
		return min + rand.Float32()*(max-min)
	}

	sensorAngle := Radians(uniform(sensorAngleMin, sensorAngleMax))
	sensorDistance := uniform(sensorDistanceMin, sensorDistanceMax)
	rotationAngle := Radians(uniform(rotationAngleMin, rotationAngleMax))
	stepDistance := uniform(stepDistanceMin, stepDistanceMax)
	depositionAmount := uniform(depositionAmountMin, depositionAmountMax)
	decayFactor := uniform(decayFactorMin, decayFactorMax)

	return Config{
		SensorAngle:      sensorAngle,
		SensorDistance:   sensorDistance,
		RotationAngle:    rotationAngle,
		StepDistance:     stepDistance,
		DepositionAmount: depositionAmount,
		DecayFactor:      decayFactor,
	}
}

func RandomConfigs(n int) []Config {
	configs := make([]Config, n)
	for i := range configs {
		configs[i] = RandomConfig()
	}
	return configs
}

func RandomAttractionTable(n int) [][]float32 {
	normal := func(mean, std float32) float32 {
		return mean + float32(rand.NormFloat64())*std
	}

	result := make([][]float32, n)
	for i := range result {
		result[i] = make([]float32, n)
		for j := range result[i] {
			if i == j {
				result[i][j] = normal(attractionFactorMean, attractionFactorStd)
			} else {
				result[i][j] = normal(repulsionFactorMean, repulsionFactorStd)
			}
		}
	}
	return result
}

func PrintConfigs(configs []Config, table [][]float32) {
	fmt.Println("configs = []Config{")
	for _, c := range configs {
		fmt.Printf("\tConfig{%v, %v, %v, %v, %v, %v},\n",
			c.SensorAngle,
			c.SensorDistance,
			c.RotationAngle,
			c.StepDistance,
			c.DepositionAmount,
			c.DecayFactor)
	}
	fmt.Println("}")
	fmt.Println("table = [][]float32{")
	for _, row := range table {
		fmt.Printf("\t{")
		for i, value := range row {
			if i > 0 {
				fmt.Printf(", ")
			}
			fmt.Printf("%v", value)
		}
		fmt.Println("},")
	}
	fmt.Println("}")
}

func SummarizeConfigs(configs []Config) {
	summarize := func(name string, getter func(i int) float32) {
		fmt.Printf("%s ", name)
		for i := 0; i < 18-len(name); i++ {
			fmt.Printf(".")
		}
		for i := range configs {
			if i != 0 {
				fmt.Printf(",")
			}
			fmt.Printf("% 8.3f", getter(i))
		}
		fmt.Printf("\n")
	}

	summarize("StepDistance", func(i int) float32 {
		return configs[i].StepDistance
	})
	summarize("SensorDistance", func(i int) float32 {
		return configs[i].SensorDistance
	})
	summarize("SensorAngle", func(i int) float32 {
		return Degrees(configs[i].SensorAngle)
	})
	summarize("RotationAngle", func(i int) float32 {
		return Degrees(configs[i].RotationAngle)
	})
	summarize("DecayFactor", func(i int) float32 {
		return configs[i].DecayFactor
	})
}
