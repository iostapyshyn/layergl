package layergl

import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
)

var (
	polygonShader shader
	circleShader  shader
	textureShader shader
	vertBuffer    vertexBuffer
)

func DrawTexture(d Texture) {
	vertBuffer.loadVertexArray(d.GetVertexArray())
	vertBuffer.loadUVs([]float32{
		0.0, 0.0,
		0.0, 1.0,
		1.0, 0.0,
		1.0, 1.0,
	})
	textureShader.drawTexture(vertBuffer, d)
}

func DrawPolygon(d VertexObject, color Color) {
	vertBuffer.loadVertexArray(d.GetVertexArray())
	polygonShader.drawColor(vertBuffer, color)
}

func DrawPoint(d Point, r float32, color Color) {
	circleShader.setUniformVec("circle", d.X, d.Y, r)
	vertBuffer.loadVertexArray(Square(d, r*2).GetVertexArray())
	circleShader.drawColor(vertBuffer, color)
}

func DrawLines(points []Point, color Color) {
	vo := VertexObject{}
	for i, v := range points {
		vo.Vertices = append(vo.Vertices, v)
		vo.Indices = append(vo.Indices, uint32(i))
	}

	vertBuffer.loadVertexArray(vo.GetVertexArray())
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
	gl.ClearColor(color[0], color[1], color[2], color[3])
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
