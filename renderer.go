package layergl

import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"time"
)

type Renderer struct {
	polygonShader           shader
	circleShader            shader
	textureShader           shader
	vertexBuffer            vertexBuffer
	lastTime                int64
	frames, FramesPerSecond int
}

func (ren Renderer) DrawTexture(d Texture) {
	ren.vertexBuffer.loadVertexArray(d.GetVertexArray())
	ren.vertexBuffer.loadUVs([]float32{
		0.0, 0.0,
		0.0, 1.0,
		1.0, 0.0,
		1.0, 1.0,
	})
	ren.textureShader.drawTexture(ren.vertexBuffer, d)
}

func (ren Renderer) DrawPolygon(d VertexObject, color Color) {
	ren.vertexBuffer.loadVertexArray(d.GetVertexArray())
	ren.polygonShader.drawColor(ren.vertexBuffer, color)
}

func (ren Renderer) DrawPoint(d Point, r float32, color Color) {
	ren.circleShader.setUniformVec("circle", d.X, d.Y, r)
	ren.vertexBuffer.loadVertexArray(Square(d, r*2).GetVertexArray())
	ren.circleShader.drawColor(ren.vertexBuffer, color)
}

func (ren Renderer) DrawLines(points []Point, color Color) {
	vo := VertexObject{}
	for i, v := range points {
		vo.Vertices = append(vo.Vertices, v)
		vo.Indices = append(vo.Indices, uint32(i))
	}

	ren.vertexBuffer.loadVertexArray(vo.GetVertexArray())
	ren.polygonShader.drawLines(ren.vertexBuffer, color)
}

func (ren *Renderer) Clear() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	ren.UpdateFPS()
}

func CreateRenderer(width, height int) (Renderer, error) {
	if err := gl.Init(); err != nil {
		return Renderer{}, err
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

	ren := Renderer{}

	ren.vertexBuffer = newVertexBuffer(128)

	ren.polygonShader = newShaderProgram(vertexVert, polygonFrag)
	ren.circleShader = newShaderProgram(vertexVert, circleFrag)
	ren.textureShader = newShaderProgram(textureVert, textureFrag)

	projectionMatrix := orthoProjection(0, float32(width), 0, float32(height), -1, 1)
	ren.polygonShader.setUniformMat("projection", projectionMatrix)
	ren.circleShader.setUniformMat("projection", projectionMatrix)
	ren.textureShader.setUniformMat("projection", projectionMatrix)

	ren.textureShader.setUniformVec("tex", 0)

	ren.UpdateFPS()

	return ren, nil
}

func (ren Renderer) ClearColor(color Color) {
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

func (ren *Renderer) UpdateFPS() {
	nowTime := time.Now().UnixNano()
	ren.frames++
	if nowTime-ren.lastTime >= int64(time.Second) {
		ren.lastTime = nowTime
		ren.FramesPerSecond = ren.frames
		ren.frames = 0
	}
}
