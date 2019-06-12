package layergl

import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"strings"
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

func (shader shader) drawTexture(vao *vertexBuffer, texture *Texture) {
	texture.bind()
	shader.bind()
	vao.bind()

	gl.DrawElements(gl.TRIANGLES, int32(vao.count), gl.UNSIGNED_INT, nil)
}

func (shader shader) drawColor(vao *vertexBuffer, color Color) {
	shader.setUniformVec("color", float32(color.R), float32(color.G), float32(color.B), float32(color.A))
	vao.bind()

	gl.DrawElements(gl.TRIANGLES, int32(vao.count), gl.UNSIGNED_INT, nil)
}

func (shader shader) drawLines(vao *vertexBuffer, color Color) {
	shader.setUniformVec("color", float32(color.R), float32(color.G), float32(color.B), float32(color.A))
	vao.bind()

	gl.DrawElements(gl.LINE_STRIP, int32(vao.count), gl.UNSIGNED_INT, nil)
}

// Links vertex and fragment shaders.
func newShaderProgram(vs, fs string) shader {
	vertexShader, err := compileShader(vs, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragmentShader, err := compileShader(fs, gl.FRAGMENT_SHADER)
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

// Compiles a shader.
func compileShader(data string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	source := string(data)
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

// Shader sources:

const circleFrag = `
#version 330
out vec4 frag_color;

uniform vec4 color;
uniform vec3 circle;

const int aa = 1;

void main() {
    float d = distance(gl_FragCoord.xy, circle.xy);
    float k = 1-smoothstep(circle.z-(2*aa), circle[2], d);
    
    frag_color = vec4(color.xyz, k*color.w);
}
`

const polygonFrag = `
#version 330
out vec4 frag_color;

uniform vec4 color;

void main() {
    frag_color = color;
}
`

const textureFrag = `
#version 330
out vec4 frag_color;

in vec2 fragTexCoord;

uniform sampler2D tex;

void main() {
    frag_color = texture(tex, vec2(fragTexCoord.x, 1-fragTexCoord.y)); // Flip Y axis
}
`

const textureVert = `
#version 330
layout(location = 0)in vec3 vert;
layout(location = 1)in vec2 vertTexCoord;
out vec2 fragTexCoord;

uniform mat4 projection;

void main() {
    fragTexCoord = vertTexCoord;
    gl_Position = projection * vec4(vert, 1);
}
`

const vertexVert = `
#version 330
layout(location = 0) in vec2 vert;

uniform mat4 projection;

void main() {
    gl_Position = projection * vec4(vert, 0.0, 1.0);
}
`

const fontVert = `
#version 330
layout(location = 0) in vec2 vert;

//pass through to fragTexCoord
in vec2 vertTexCoord;

//window res
uniform vec2 resolution;

//pass to frag
out vec2 fragTexCoord;

void main() {
   // convert the rectangle from pixels to 0.0 to 1.0
   vec2 zeroToOne = vert / resolution;
   // convert from 0->1 to 0->2
   vec2 zeroToTwo = zeroToOne * 2.0;
   // convert from 0->2 to -1->+1 (clipspace)
   vec2 clipSpace = zeroToTwo - 1.0;
   fragTexCoord = vertTexCoord;
   gl_Position = vec4(clipSpace * vec2(1, -1), 0, 1);
}
`
