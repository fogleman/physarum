package physarum

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

func Radians(degrees float32) float32 {
	return degrees * math.Pi / 180
}

func Degrees(radians float32) float32 {
	return radians * 180 / math.Pi
}

func SavePNG(path string, im image.Image, level png.CompressionLevel) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	var encoder png.Encoder
	encoder.CompressionLevel = level
	return encoder.Encode(file, im)
}

func HexColor(x int) color.RGBA {
	r := uint8((x >> 16) & 0xff)
	g := uint8((x >> 8) & 0xff)
	b := uint8((x >> 0) & 0xff)
	return color.RGBA{r, g, b, 0xff}
}

func IsPowerOfTwo(x int) bool {
	return (x & (x - 1)) == 0
}

func Shift(x, size float32) float32 {
	if x < 0 {
		return x + size
	}
	if x >= size {
		return x - size
	}
	return x
}
