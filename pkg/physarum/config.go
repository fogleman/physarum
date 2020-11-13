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
	RepulsionFactor      float32
}

func RandomConfig() Config {
	sensorAngle := rand.Float32() * Radians(120)
	sensorDistance := rand.Float32() * 64
	rotationAngle := rand.Float32() * Radians(120)
	stepDistance := 0.5 + rand.Float32()*2.5
	decayFactor := 0.1 + rand.Float32()*0.2
	repulsionFactor := 1 + float32(rand.NormFloat64())*0.5
	return Config{
		PopulationPercentage: 1,
		SensorAngle:          sensorAngle,
		SensorDistance:       sensorDistance,
		RotationAngle:        rotationAngle,
		StepDistance:         stepDistance,
		DepositionAmount:     5,
		DecayFactor:          decayFactor,
		RepulsionFactor:      repulsionFactor,
	}
}

func RandomConfigs(n int) []Config {
	configs := make([]Config, n)
	for i := range configs {
		configs[i] = RandomConfig()
	}
	return configs
}

func PrintConfigs(configs []Config) {
	for _, c := range configs {
		fmt.Printf("Config{%v, %v, %v, %v, %v, %v, %v, %v},\n",
			c.PopulationPercentage,
			c.SensorAngle,
			c.SensorDistance,
			c.RotationAngle,
			c.StepDistance,
			c.DepositionAmount,
			c.DecayFactor,
			c.RepulsionFactor)
	}
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
	summarize("RepulsionFactor", func(i int) float32 {
		return configs[i].RepulsionFactor
	})
}
