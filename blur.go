package physarum

import "math"

// TODO: wrap blur window at edges

func boxBlurH(src, dst []float64, w, h, r int, scale float64) {
	m := scale / float64(r+r+1)
	for i := 0; i < h; i++ {
		ti := i * w
		li := ti
		ri := ti + r
		fv := src[ti]
		lv := src[ti+w-1]
		val := float64(r+1) * fv
		for j := 0; j < r; j++ {
			val += src[ti+j]
		}
		for j := 0; j <= r; j++ {
			val += src[ri] - fv
			dst[ti] = val * m
			ri++
			ti++
		}
		for j := r + 1; j < w-r; j++ {
			val += src[ri] - src[li]
			dst[ti] = val * m
			li++
			ri++
			ti++
		}
		for j := w - r; j < w; j++ {
			val += lv - src[li]
			dst[ti] = val * m
			li++
			ti++
		}
	}
}

func boxBlurV(src, dst []float64, w, h, r int, scale float64) {
	m := scale / float64(r+r+1)
	for i := 0; i < w; i++ {
		ti := i
		li := ti
		ri := ti + r*w
		fv := src[ti]
		lv := src[ti+w*(h-1)]
		val := float64(r+1) * fv
		for j := 0; j < r; j++ {
			val += src[ti+j*w]
		}
		for j := 0; j <= r; j++ {
			val += src[ri] - fv
			dst[ti] = val * m
			ri += w
			ti += w
		}
		for j := r + 1; j < h-r; j++ {
			val += src[ri] - src[li]
			dst[ti] = val * m
			li += w
			ri += w
			ti += w
		}
		for j := h - r; j < h; j++ {
			val += lv - src[li]
			dst[ti] = val * m
			li += w
			ti += w
		}
	}
}

func boxBlur(src, tmp []float64, w, h, r int, scale float64) {
	boxBlurH(src, tmp, w, h, r, 1)
	boxBlurV(tmp, src, w, h, r, scale)
}

func boxesForGaussian(sigma float64, n int) []int {
	wIdeal := math.Sqrt((12 * sigma * sigma / float64(n)) + 1)
	wl := int(wIdeal)
	if wl%2 == 0 {
		wl--
	}
	wu := wl + 2
	mIdeal := (12*sigma*sigma - float64(n*wl*wl-4*n*wl-3*n)) / float64(-4*wl-4)
	m := int(math.Round(mIdeal))
	sizes := make([]int, n)
	for i := range sizes {
		if i < m {
			sizes[i] = wl
		} else {
			sizes[i] = wu
		}
	}
	return sizes
}

func gaussianBlur(src, tmp []float64, w, h, r int, scale float64) {
	boxes := boxesForGaussian(float64(r), 3)
	boxBlur(src, tmp, w, h, (boxes[0]-1)/2, 1)
	boxBlur(src, tmp, w, h, (boxes[1]-1)/2, 1)
	boxBlur(src, tmp, w, h, (boxes[2]-1)/2, scale)
}
