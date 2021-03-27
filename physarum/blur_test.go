package physarum

import (
	"testing"
)

func TestBoxBlurH(t *testing.T) {
	slow := func(src, dst []float32, w, h, r int, scale float32) {
		m := scale / float32(r+r+1)
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				var val float32
				for k := -r; k <= r; k++ {
					i := y*w + (x+w+k)%w
					val += src[i]
				}
				dst[y*w+x] = val * m
			}
		}
	}

	w := 10
	h := 10
	src := make([]float32, w*h)
	dst1 := make([]float32, w*h)
	dst2 := make([]float32, w*h)
	for i := range src {
		src[i] = float32(i)
	}
	for r := 0; r < 5; r++ {
		boxBlurH(src, dst1, w, h, r, 1)
		slow(src, dst2, w, h, r, 1)
		for i := range src {
			if dst1[i] != dst2[i] {
				t.Fatalf("got %v, want %v", dst1, dst2)
			}
		}
	}
}

func TestBoxBlurV(t *testing.T) {
	slow := func(src, dst []float32, w, h, r int, scale float32) {
		m := scale / float32(r+r+1)
		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				var val float32
				for k := -r; k <= r; k++ {
					i := x + ((y+h+k)%h)*w
					val += src[i]
				}
				dst[y*w+x] = val * m
			}
		}
	}

	w := 10
	h := 10
	src := make([]float32, w*h)
	dst1 := make([]float32, w*h)
	dst2 := make([]float32, w*h)
	for i := range src {
		src[i] = float32(i)
	}
	for r := 0; r < 5; r++ {
		boxBlurV(src, dst1, w, h, r, 1)
		slow(src, dst2, w, h, r, 1)
		for i := range src {
			if dst1[i] != dst2[i] {
				t.Fatalf("got %v, want %v", dst1, dst2)
			}
		}
	}
}
