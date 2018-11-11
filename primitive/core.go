package primitive

import (
	"image"
	"math"
)

func myComputeColor(target, current *image.RGBA, lines []Scanline, alpha int) Color {
	var rsum, gsum, bsum, count int64

	a := 0x101 * 255 / alpha

	for _, line := range lines {
		tr := TargetMemo[line.Y][line.X2][0]
		tg := TargetMemo[line.Y][line.X2][1]
		tb := TargetMemo[line.Y][line.X2][2]

		cr := CurrentMemo[line.Y][line.X2][0]
		cg := CurrentMemo[line.Y][line.X2][1]
		cb := CurrentMemo[line.Y][line.X2][2]

		if line.X1 > 0 {
			tr -= TargetMemo[line.Y][line.X1-1][0]
			tg -= TargetMemo[line.Y][line.X1-1][1]
			tb -= TargetMemo[line.Y][line.X1-1][2]

			cr -= CurrentMemo[line.Y][line.X1-1][0]
			cg -= CurrentMemo[line.Y][line.X1-1][1]
			cb -= CurrentMemo[line.Y][line.X1-1][2]
		}

		rsum += int64((tr-cr)*int64(a) + cr*0x101)
		gsum += int64((tg-cg)*int64(a) + cg*0x101)
		bsum += int64((tb-cb)*int64(a) + cb*0x101)
		count += int64(line.X2-line.X1)+1
	}

	if count == 0 {
		return Color{}
	}

	r := clampInt(int(rsum/count)>>8, 0, 255)
	g := clampInt(int(gsum/count)>>8, 0, 255)
	b := clampInt(int(bsum/count)>>8, 0, 255)

	return Color{r, g, b, alpha}
}

func computeColor(target, current *image.RGBA, lines []Scanline, alpha int) Color {
	var rsum, gsum, bsum, count int64
	a := 0x101 * 255 / alpha
	for _, line := range lines {
		i := target.PixOffset(line.X1, line.Y)
		for x := line.X1; x <= line.X2; x++ {
			tr := int(target.Pix[i])
			tg := int(target.Pix[i+1])
			tb := int(target.Pix[i+2])
			cr := int(current.Pix[i])
			cg := int(current.Pix[i+1])
			cb := int(current.Pix[i+2])
			i += 4
			rsum += int64((tr-cr)*a + cr*0x101)
			gsum += int64((tg-cg)*a + cg*0x101)
			bsum += int64((tb-cb)*a + cb*0x101)
			count++
		}
	}
	if count == 0 {
		return Color{}
	}
	r := clampInt(int(rsum/count)>>8, 0, 255)
	g := clampInt(int(gsum/count)>>8, 0, 255)
	b := clampInt(int(bsum/count)>>8, 0, 255)

	return Color{r, g, b, alpha}
}

func copyLines(dst, src *image.RGBA, lines []Scanline) {
	for _, line := range lines {
		a := dst.PixOffset(line.X1, line.Y)
		b := a + (line.X2-line.X1+1)*4
		copy(dst.Pix[a:b], src.Pix[a:b])
	}
}

func drawLines(im *image.RGBA, c Color, lines []Scanline) {
	const m = 0xffff
	sr, sg, sb, sa := c.NRGBA().RGBA()
	for _, line := range lines {
		ma := line.Alpha
		a := (m - sa*ma/m) * 0x101
		i := im.PixOffset(line.X1, line.Y)
		for x := line.X1; x <= line.X2; x++ {
			dr := uint32(im.Pix[i+0])
			dg := uint32(im.Pix[i+1])
			db := uint32(im.Pix[i+2])
			da := uint32(im.Pix[i+3])
			im.Pix[i+0] = uint8((dr*a + sr*ma) / m >> 8)
			im.Pix[i+1] = uint8((dg*a + sg*ma) / m >> 8)
			im.Pix[i+2] = uint8((db*a + sb*ma) / m >> 8)
			im.Pix[i+3] = uint8((da*a + sa*ma) / m >> 8)
			i += 4
		}
	}
}

func differenceFull(a, b *image.RGBA) float64 {
	size := a.Bounds().Size()
	w, h := size.X, size.Y
	var total uint64
	for y := 0; y < h; y++ {
		i := a.PixOffset(0, y)
		for x := 0; x < w; x++ {
			ar := int(a.Pix[i])
			ag := int(a.Pix[i+1])
			ab := int(a.Pix[i+2])
			aa := int(a.Pix[i+3])
			br := int(b.Pix[i])
			bg := int(b.Pix[i+1])
			bb := int(b.Pix[i+2])
			ba := int(b.Pix[i+3])
			i += 4
			dr := ar - br
			dg := ag - bg
			db := ab - bb
			da := aa - ba
			total += uint64(dr*dr + dg*dg + db*db + da*da)
		}
	}
	return math.Sqrt(float64(total)/float64(w*h*4)) / 255
}

func myDifferencePartial(target, before, after *image.RGBA, score float64, lines []Scanline) float64 {
	size := target.Bounds().Size()
	w, h := size.X, size.Y
	total := uint64(math.Pow(score*255, 2) * float64(w*h*4))

	ai := after.PixOffset(lines[0].X1, lines[0].Y)
	ar := int64(after.Pix[ai])
	ag := int64(after.Pix[ai+1])
	ab := int64(after.Pix[ai+2])

	for _, line := range lines {
		n := int64(line.X2 - line.X1 + 1)

		tar := TargetMemo[line.Y][line.X2][0]
		tag := TargetMemo[line.Y][line.X2][1]
		tab := TargetMemo[line.Y][line.X2][2]
		tbr := CurrentTargetMemo[line.Y][line.X2][0]
		tbg := CurrentTargetMemo[line.Y][line.X2][1]
		tbb := CurrentTargetMemo[line.Y][line.X2][2]
		bsr := CurrentSquaredMemo[line.Y][line.X2][0]
		bsg := CurrentSquaredMemo[line.Y][line.X2][1]
		bsb := CurrentSquaredMemo[line.Y][line.X2][2]

		if line.X1 > 0 {
			tar -= TargetMemo[line.Y][line.X1-1][0]
			tag -= TargetMemo[line.Y][line.X1-1][1]
			tab -= TargetMemo[line.Y][line.X1-1][2]
			tbr -= CurrentTargetMemo[line.Y][line.X1-1][0]
			tbg -= CurrentTargetMemo[line.Y][line.X1-1][1]
			tbb -= CurrentTargetMemo[line.Y][line.X1-1][2]
			bsr -= CurrentSquaredMemo[line.Y][line.X1-1][0]
			bsg -= CurrentSquaredMemo[line.Y][line.X1-1][1]
			bsb -= CurrentSquaredMemo[line.Y][line.X1-1][2]
		}

		TA := tar * ar + tag * ag + tab * ab
		As := (ar*ar + ab*ab + ag*ag) * n
		TB := tbr + tbg + tbb
		Bs := bsr + bsg + bsb

		total = uint64(int64(total) - 2 * TA + 2 * TB + As - Bs)
	}

	result := math.Sqrt(float64(total)/float64(w*h*4)) / 255

	return result
}

func differencePartial(target, before, after *image.RGBA, score float64, lines []Scanline) float64 {
	size := target.Bounds().Size()
	w, h := size.X, size.Y
	total := uint64(math.Pow(score*255, 2) * float64(w*h*4))
	for _, line := range lines {
		i := target.PixOffset(line.X1, line.Y)
		for x := line.X1; x <= line.X2; x++ {
			tr := int(target.Pix[i])
			tg := int(target.Pix[i+1])
			tb := int(target.Pix[i+2])
			ta := int(target.Pix[i+3])
			br := int(before.Pix[i])
			bg := int(before.Pix[i+1])
			bb := int(before.Pix[i+2])
			ba := int(before.Pix[i+3])
			ar := int(after.Pix[i])
			ag := int(after.Pix[i+1])
			ab := int(after.Pix[i+2])
			aa := int(after.Pix[i+3])
			i += 4
			dr1 := tr - br
			dg1 := tg - bg
			db1 := tb - bb
			da1 := ta - ba
			dr2 := tr - ar
			dg2 := tg - ag
			db2 := tb - ab
			da2 := ta - aa
			total -= uint64(dr1*dr1 + dg1*dg1 + db1*db1 + da1*da1)
			total += uint64(dr2*dr2 + dg2*dg2 + db2*db2 + da2*da2)
		}
	}

	return math.Sqrt(float64(total)/float64(w*h*4)) / 255
}
