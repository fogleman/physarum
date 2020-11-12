package physarum

import (
	"fmt"
	"math/rand"
)

type Config struct {
	PopulationPercentage float32
	SensorAngle          float32
	SensorDistance       float32
	RotationAngle        float32
	StepDistance         float32
	DepositionAmount     float32
	DecayFactor          float32
}

func RandomConfig() *Config {
	sensorAngle := rand.Float32() * Radians(90)
	sensorDistance := rand.Float32() * 64
	rotationAngle := rand.Float32() * Radians(90)
	stepDistance := 0.5 + rand.Float32()*2.5
	decayFactor := 0.1 + rand.Float32()*0.2
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
	summarize := func(name string, getter func(i int) float32) {
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
