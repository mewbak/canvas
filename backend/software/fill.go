package softwarebackend

import (
	"image"
	"image/color"
	"math"

	"github.com/tfriedel6/canvas/backend/backendbase"
)

func (b *SoftwareBackend) Clear(pts [4][2]float64) {
	iterateTriangles(pts[:], func(tri [][2]float64) {
		b.fillTriangle(tri, func(x, y int) {
			if b.clip.AlphaAt(x, y).A == 0 {
				return
			}
			b.Image.SetRGBA(x, y, color.RGBA{})
		})
	})
}

func (b *SoftwareBackend) Fill(style *backendbase.FillStyle, pts [][2]float64) {
	b.clearMask()

	if lg := style.LinearGradient; lg != nil {
		lg := lg.(*LinearGradient)
		from := [2]float64{style.Gradient.X0, style.Gradient.Y0}
		dir := [2]float64{style.Gradient.X1 - style.Gradient.X0, style.Gradient.Y1 - style.Gradient.Y0}
		dirlen := math.Sqrt(dir[0]*dir[0] + dir[1]*dir[1])
		dir[0] /= dirlen
		dir[1] /= dirlen
		b.fillTriangles(pts, func(x, y int) {
			pos := [2]float64{float64(x) - from[0], float64(y) - from[1]}
			r := (pos[0]*dir[0] + pos[1]*dir[1]) / dirlen
			color := lg.data.ColorAt(r)
			b.Image.SetRGBA(x, y, mix(color, b.Image.RGBAAt(x, y)))
		})
	} else if rg := style.RadialGradient; rg != nil {
		rg := rg.(*RadialGradient)
		from := [2]float64{style.Gradient.X0, style.Gradient.Y0}
		to := [2]float64{style.Gradient.X1, style.Gradient.Y1}
		radFrom := style.Gradient.RadFrom
		radTo := style.Gradient.RadTo
		b.fillTriangles(pts, func(x, y int) {
			pos := [2]float64{float64(x), float64(y)}
			oa := 0.5 * math.Sqrt(
				math.Pow(-2.0*from[0]*from[0]+2.0*from[0]*to[0]+2.0*from[0]*pos[0]-2.0*to[0]*pos[0]-2.0*from[1]*from[1]+2.0*from[1]*to[1]+2.0*from[1]*pos[1]-2.0*to[1]*pos[1]+2.0*radFrom*radFrom-2.0*radFrom*radTo, 2.0)-
					4.0*(from[0]*from[0]-2.0*from[0]*pos[0]+pos[0]*pos[0]+from[1]*from[1]-2.0*from[1]*pos[1]+pos[1]*pos[1]-radFrom*radFrom)*
						(from[0]*from[0]-2.0*from[0]*to[0]+to[0]*to[0]+from[1]*from[1]-2.0*from[1]*to[1]+to[1]*to[1]-radFrom*radFrom+2.0*radFrom*radTo-radTo*radTo))
			ob := (from[0]*from[0] - from[0]*to[0] - from[0]*pos[0] + to[0]*pos[0] + from[1]*from[1] - from[1]*to[1] - from[1]*pos[1] + to[1]*pos[1] - radFrom*radFrom + radFrom*radTo)
			oc := (from[0]*from[0] - 2.0*from[0]*to[0] + to[0]*to[0] + from[1]*from[1] - 2.0*from[1]*to[1] + to[1]*to[1] - radFrom*radFrom + 2.0*radFrom*radTo - radTo*radTo)
			o1 := (-oa + ob) / oc
			o2 := (oa + ob) / oc
			if math.IsNaN(o1) && math.IsNaN(o2) {
				return
			}
			o := math.Max(o1, o2)
			col := rg.data.ColorAt(o)
			b.Image.SetRGBA(x, y, mix(col, b.Image.RGBAAt(x, y)))
			b.Image.SetRGBA(x, y, color.RGBA{R: 255, A: 255})
		})
	} else {
		b.fillTriangles(pts, func(x, y int) {
			b.Image.SetRGBA(x, y, mix(style.Color, b.Image.RGBAAt(x, y)))
		})
	}
}

func (b *SoftwareBackend) FillImageMask(style *backendbase.FillStyle, mask *image.Alpha, pts [4][2]float64) {
	mw := float64(mask.Bounds().Dx())
	mh := float64(mask.Bounds().Dy())
	b.fillQuad(pts, func(x, y int, sx2, sy2 float64) {
		sxi := int(mw * sx2)
		syi := int(mh * sy2)
		a := mask.AlphaAt(sxi, syi)
		if a.A == 0 {
			return
		}
		b.Image.SetRGBA(x, y, alphaColor(style.Color, a))
	})
}

func (b *SoftwareBackend) clearMask() {
	p := b.mask.Pix
	for i := range p {
		p[i] = 0
	}
}

func (b *SoftwareBackend) ClearClip() {
	p := b.clip.Pix
	for i := range p {
		p[i] = 255
	}
}

func (b *SoftwareBackend) Clip(pts [][2]float64) {
	p2 := b.mask.Pix
	for i := range p2 {
		p2[i] = 0
	}

	iterateTriangles(pts[:], func(tri [][2]float64) {
		b.fillTriangle(tri, func(x, y int) {
			b.mask.SetAlpha(x, y, color.Alpha{A: 255})
		})
	})

	p := b.clip.Pix
	for i := range p {
		if p2[i] == 0 {
			p[i] = 0
		}
	}
}
