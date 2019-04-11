package layergl

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"log"
)

type vertexBuffer struct {
	vao, vbo, uvbo, ebo uint32
	vboSize, eboSize    int
	count               int
}

const t32Bytes = 4

func newVertexBuffer(bufferSize int) *vertexBuffer {
	if bufferSize <= 0 {
		bufferSize = 1
	}

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, bufferSize*t32Bytes, gl.Ptr(nil), gl.DYNAMIC_DRAW)

	// vert attribute
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 0, nil)

	var uvbo uint32
	gl.GenBuffers(1, &uvbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, uvbo)
	gl.BufferData(gl.ARRAY_BUFFER, bufferSize*t32Bytes, gl.Ptr(nil), gl.DYNAMIC_DRAW)

	// texCoord attribute
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 0, nil)

	var ebo uint32
	gl.GenBuffers(1, &ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, bufferSize*t32Bytes, gl.Ptr(nil), gl.DYNAMIC_DRAW)

	return &vertexBuffer{vao, vbo, uvbo, ebo, bufferSize, bufferSize, 0}
}

func (v *vertexBuffer) loadUVs(uv []float32) {
	if len(uv) > v.vboSize {
		log.Println("Error loading UVs.")
		return
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, v.uvbo)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(uv)*t32Bytes, gl.Ptr(uv))
}

func (v *vertexBuffer) loadVertexArray(vertices []float32, elements []uint32) {
	gl.BindVertexArray(v.vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, v.vbo)
	if len(vertices) > v.vboSize {
		for len(vertices) > v.vboSize {
			v.vboSize *= 2
		}
		log.Println("Reallocating VBO:", v.vboSize)
		gl.BufferData(gl.ARRAY_BUFFER, v.vboSize*t32Bytes, gl.Ptr(vertices), gl.DYNAMIC_DRAW)
	} else {
		gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(vertices)*t32Bytes, gl.Ptr(vertices))
	}

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, v.ebo)
	if len(elements) > v.eboSize {
		for len(elements) > v.eboSize {
			v.eboSize *= 2
		}
		log.Println("Reallocating EBO:", v.eboSize)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, v.eboSize*t32Bytes, gl.Ptr(elements), gl.DYNAMIC_DRAW)
	} else {
		gl.BufferSubData(gl.ELEMENT_ARRAY_BUFFER, 0, len(elements)*t32Bytes, gl.Ptr(elements))
	}

	v.count = len(elements)
}

func (v *vertexBuffer) bind() {
	gl.BindVertexArray(v.vao)
}
