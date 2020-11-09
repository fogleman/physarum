package physarum

type Config struct {
	PopulationPercentage float64
	SensorAngle          float64
	SensorDistance       float64
	RotationAngle        float64
	StepDistance         float64
	DepositionAmount     float64
	DecayFactor          float64
}

var DefaultConfigs = []*Config{
	&Config{
		PopulationPercentage: 0.25,
		SensorAngle:          Radians(22.5),
		SensorDistance:       9,
		RotationAngle:        Radians(45),
		StepDistance:         1,
		DepositionAmount:     5,
		DecayFactor:          0.1,
	},
	&Config{
		PopulationPercentage: 0.25,
		SensorAngle:          Radians(22.5),
		SensorDistance:       32,
		RotationAngle:        Radians(45),
		StepDistance:         1,
		DepositionAmount:     5,
		DecayFactor:          0.1,
	},
}
