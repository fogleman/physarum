package physarum

import (
	"image"
	"image/png"
	"math"
	"os"
)

func Radians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

func Degrees(radians float64) float64 {
	return radians * 180 / math.Pi
}

func SavePNG(path string, im image.Image) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return png.Encode(file, im)
}

func mod(a, b float64) float64 {
	return math.Mod(math.Mod(a, b)+b, b)
}
