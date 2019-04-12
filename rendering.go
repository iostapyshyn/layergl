package layergl

import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
)

var (
	polygonShader shader
	circleShader  shader
	textureShader shader
	vertBuffer    *vertexBuffer
)

func DrawTexture(d *Texture) {
	vertBuffer.loadVertexArray(d.vertexArray())
	vertBuffer.loadUVs([]float32{
		0.0, 0.0,
		0.0, 1.0,
		1.0, 0.0,
		1.0, 1.0,
	})
	textureShader.drawTexture(vertBuffer, d)
}

func DrawRect(rect Rect, color Color) {
	vertBuffer.loadVertexArray(rect.vertexArray())
	polygonShader.drawColor(vertBuffer, color)
}

func DrawVertexObject(d *VertexObject, color Color) {
	vertBuffer.loadVertexArray(d.vertexArray())
	polygonShader.drawColor(vertBuffer, color)
}

func DrawPoint(d Point, r float64, color Color) {
	rect := Rect{d.X - r, d.Y - r, d.X + r, d.Y + r}
	circleShader.setUniformVec("circle", float32(d.X), float32(d.Y), float32(r))
	vertBuffer.loadVertexArray(rect.vertexArray())
	circleShader.drawColor(vertBuffer, color)
}

func DrawLines(points []Point, color Color) {
	vertices := make([]float32, 0, len(points)*2)
	elements := make([]uint32, 0, len(points))

	for i, p := range points {
		vertices = append(vertices, float32(p.X))
		vertices = append(vertices, float32(p.Y))
		elements = append(elements, uint32(i))
	}

	vertBuffer.loadVertexArray(vertices, elements)
	polygonShader.drawLines(vertBuffer, color)
}

func Clear() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

func Init(width, height int) error {
	if err := gl.Init(); err != nil {
		return err
	}

	versionString := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL Version", versionString)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	gl.Enable(gl.MULTISAMPLE)

	// gl.Enable(gl.LINE_SMOOTH)
	// gl.Hint(gl.LINE_SMOOTH_HINT, gl.NICEST)
	// gl.Enable(gl.POLYGON_SMOOTH)
	// gl.Hint(gl.POLYGON_SMOOTH_HINT, gl.NICEST)

	gl.ClearColor(0, 0, 0, 1)

	vertBuffer = newVertexBuffer(128)

	polygonShader = newShaderProgram(vertexVert, polygonFrag)
	circleShader = newShaderProgram(vertexVert, circleFrag)
	textureShader = newShaderProgram(textureVert, textureFrag)

	projectionMatrix := orthoProjection(0, float32(width), 0, float32(height), -1, 1)
	polygonShader.setUniformMat("projection", projectionMatrix)
	circleShader.setUniformMat("projection", projectionMatrix)
	textureShader.setUniformMat("projection", projectionMatrix)

	textureShader.setUniformVec("tex", 0)

	return nil
}

func ClearColor(color Color) {
	gl.ClearColor(float32(color.R), float32(color.G), float32(color.B), float32(color.A))
}

func orthoProjection(left, right, bottom, top, near, far float32) []float32 {
	rml, tmb, fmn := (right - left), (top - bottom), (far - near)
	return []float32{
		2. / rml, 0, 0, 0,
		0, 2. / tmb, 0, 0,
		0, 0, 2. / fmn, 0,
		-(right + left) / rml, -(top + bottom) / tmb, -(far + near) / fmn, 1,
	}
}
