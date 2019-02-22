package goglbackend

import (
	"unsafe"

	"github.com/go-gl/gl/v3.2-core/gl"
)

func (b *GoGLBackend) ClearClip() {
	gl.StencilMask(0xFF)
	gl.Clear(gl.STENCIL_BUFFER_BIT)
}

func (b *GoGLBackend) Clip(pts [][2]float64) {
	b.ptsBuf = b.ptsBuf[:0]
	b.ptsBuf = append(b.ptsBuf,
		0, 0,
		0, float32(b.fh),
		float32(b.fw), float32(b.fh),
		float32(b.fw), 0)
	for _, pt := range pts {
		b.ptsBuf = append(b.ptsBuf, float32(pt[0]), float32(pt[1]))
	}

	mode := uint32(gl.TRIANGLES)
	if len(pts) == 4 {
		mode = gl.TRIANGLE_FAN
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, b.buf)
	gl.BufferData(gl.ARRAY_BUFFER, len(b.ptsBuf)*4, unsafe.Pointer(&b.ptsBuf[0]), gl.STREAM_DRAW)
	gl.VertexAttribPointer(b.sr.Vertex, 2, gl.FLOAT, false, 0, nil)

	gl.UseProgram(b.sr.ID)
	gl.Uniform4f(b.sr.Color, 1, 1, 1, 1)
	gl.Uniform2f(b.sr.CanvasSize, float32(b.fw), float32(b.fh))
	gl.EnableVertexAttribArray(b.sr.Vertex)

	gl.ColorMask(false, false, false, false)

	gl.StencilMask(0x04)
	gl.StencilFunc(gl.ALWAYS, 4, 0x04)
	gl.StencilOp(gl.REPLACE, gl.REPLACE, gl.REPLACE)
	gl.DrawArrays(mode, 4, int32(len(pts)))

	gl.StencilMask(0x02)
	gl.StencilFunc(gl.EQUAL, 0, 0x06)
	gl.StencilOp(gl.KEEP, gl.INVERT, gl.INVERT)
	gl.DrawArrays(gl.TRIANGLE_FAN, 0, 4)

	gl.StencilMask(0x04)
	gl.StencilFunc(gl.ALWAYS, 0, 0x04)
	gl.StencilOp(gl.ZERO, gl.ZERO, gl.ZERO)
	gl.DrawArrays(gl.TRIANGLE_FAN, 0, 4)

	gl.DisableVertexAttribArray(b.sr.Vertex)

	gl.ColorMask(true, true, true, true)
	gl.StencilOp(gl.KEEP, gl.KEEP, gl.KEEP)
	gl.StencilMask(0xFF)
	gl.StencilFunc(gl.EQUAL, 0, 0xFF)
}
