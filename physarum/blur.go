package physarum

func boxBlurH(src, dst []float32, w, h, r int, scale float32) {
	m := scale / float32(r+r+1)
	ww := w - (r*2 + 1)
	for i := 0; i < h; i++ {
		ti := i * w
		li := ti + w - 1 - r
		ri := ti + r
		val := src[li]
		for j := 0; j < r; j++ {
			val += src[li+j+1]
			val += src[ti+j]
		}
		for j := 0; j <= r; j++ {
			val += src[ri] - src[li]
			dst[ti] = val * m
			li++
			ri++
			ti++
		}
		li = i * w
		for j := 0; j < ww; j++ {
			val += src[ri] - src[li]
			dst[ti] = val * m
			li++
			ri++
			ti++
		}
		ri = i * w
		for j := 0; j < r; j++ {
			val += src[ri] - src[li]
			dst[ti] = val * m
			li++
			ri++
			ti++
		}
	}
}

func boxBlurV(src, dst []float32, w, h, r int, scale float32) {
	m := scale / float32(r+r+1)
	hh := h - (r*2 + 1)
	for i := 0; i < w; i++ {
		ti := i
		li := ti + (h-1-r)*w
		ri := ti + r*w
		val := src[li]
		for j := 0; j < r; j++ {
			val += src[li+(j+1)*w]
			val += src[ti+j*w]
		}
		for j := 0; j <= r; j++ {
			val += src[ri] - src[li]
			dst[ti] = val * m
			li += w
			ri += w
			ti += w
		}
		li = i
		for j := 0; j < hh; j++ {
			val += src[ri] - src[li]
			dst[ti] = val * m
			li += w
			ri += w
			ti += w
		}
		ri = i
		for j := 0; j < r; j++ {
			val += src[ri] - src[li]
			dst[ti] = val * m
			li += w
			ri += w
			ti += w
		}
	}
}

func boxBlur(src, tmp []float32, w, h, r int, scale float32) {
	boxBlurH(src, tmp, w, h, r, 1)
	boxBlurV(tmp, src, w, h, r, scale)
}
