package layergl

import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"io/ioutil"
	"strings"
)

const (
	shaderDir = "../shaders/"
)

type shader uint32

func (v shader) setUniformMat(name string, val []float32) error {
	v.bind()

	location := gl.GetUniformLocation(uint32(v), gl.Str(name+"\x00"))
	if location == -1 {
		err := fmt.Errorf("setUniformMat(\"%s\", %v): unable to find uniform location", name, val)
		fmt.Println(err)
		return err
	}

	switch len(val) {
	case 2 * 2:
		gl.UniformMatrix2fv(location, 1, false, &val[0])
	case 3 * 3:
		gl.UniformMatrix3fv(location, 1, false, &val[0])
	case 4 * 4:
		gl.UniformMatrix4fv(location, 1, false, &val[0])
	default:
		err := fmt.Errorf("setUniformMat(\"%s\", %v): wrong number of elements in matrix", name, val)
		fmt.Println(err)
		return err
	}

	return nil
}

func (v shader) setUniformVec(name string, val ...float32) error {
	v.bind()

	location := gl.GetUniformLocation(uint32(v), gl.Str(name+"\x00"))
	if location == -1 {
		err := fmt.Errorf("setUniformVec(\"%s\", %v): unable to find uniform location", name, val)
		fmt.Println(err)
		return err
	}

	switch len(val) {
	case 1:
		gl.Uniform1f(location, val[0])
	case 2:
		gl.Uniform2f(location, val[0], val[1])
	case 3:
		gl.Uniform3f(location, val[0], val[1], val[2])
	case 4:
		gl.Uniform4f(location, val[0], val[1], val[2], val[3])
	default:
		err := fmt.Errorf("setUniformVec(\"%s\", %v): wrong number of arguments", name, val)
		fmt.Println(err)
		return err
	}

	return nil
}

func (v shader) bind() {
	gl.UseProgram(uint32(v))
}

func (shader shader) drawTexture(vao vertexBuffer, texture Texture) {
	texture.bind()
	shader.bind()
	vao.bind()

	gl.DrawElements(gl.TRIANGLES, int32(vao.count), gl.UNSIGNED_INT, nil)
}

func (shader shader) drawColor(vao vertexBuffer, color Color) {
	shader.setUniformVec("color", color[:]...)
	vao.bind()

	gl.DrawElements(gl.TRIANGLES, int32(vao.count), gl.UNSIGNED_INT, nil)
}

func (shader shader) drawLines(vao vertexBuffer, color Color) {
	shader.setUniformVec("color", color[:]...)
	vao.bind()

	gl.DrawElements(gl.LINE_STRIP, int32(vao.count), gl.UNSIGNED_INT, nil)
}

func newShaderProgram(vs, fs string) shader {
	vertexShader, err := loadShader(vs, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragmentShader, err := loadShader(fs, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	shader := shader(program)

	return shader
}

func loadShader(fileName string, shaderType uint32) (uint32, error) {
	fileName = shaderDir + fileName
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	shader := gl.CreateShader(shaderType)

	source := string(b)
	csource, free := gl.Strs(source + "\x00")
	gl.ShaderSource(shader, 1, csource, nil)
	free()

	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength)+1)
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}
