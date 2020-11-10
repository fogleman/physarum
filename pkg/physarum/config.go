package physarum

import (
	"fmt"
	"math/rand"
)

type Config struct {
	PopulationPercentage float64
	SensorAngle          float64
	SensorDistance       float64
	RotationAngle        float64
	StepDistance         float64
	DepositionAmount     float64
	DecayFactor          float64
}

func RandomConfig() *Config {
	sensorAngle := rand.Float64() * Radians(90)
	sensorDistance := rand.Float64() * 64
	rotationAngle := rand.Float64() * Radians(90)
	stepDistance := 0.5 + rand.Float64()*2.5
	decayFactor := 0.1 + rand.Float64()*0.2
	return &Config{
		PopulationPercentage: 0.5,
		SensorAngle:          sensorAngle,
		SensorDistance:       sensorDistance,
		RotationAngle:        rotationAngle,
		StepDistance:         stepDistance,
		DepositionAmount:     5,
		DecayFactor:          decayFactor,
	}
}

func RandomConfigs(n int) []*Config {
	configs := make([]*Config, n)
	for i := range configs {
		configs[i] = RandomConfig()
	}
	return configs
}

func SummarizeConfigs(configs []*Config) {
	summarize := func(name string, getter func(i int) float64) {
		fmt.Printf("%s ", name)
		for i := 0; i < 17-len(name); i++ {
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

	summarize("StepDistance", func(i int) float64 {
		return configs[i].StepDistance
	})
	summarize("SensorDistance", func(i int) float64 {
		return configs[i].SensorDistance
	})
	summarize("SensorAngle", func(i int) float64 {
		return Degrees(configs[i].SensorAngle)
	})
	summarize("RotationAngle", func(i int) float64 {
		return Degrees(configs[i].RotationAngle)
	})
	summarize("DecayFactor", func(i int) float64 {
		return configs[i].DecayFactor
	})
}
