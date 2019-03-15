package layergl

import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

type Texture struct {
	VertexObject
	width, height float32
	tex           uint32
}

func loadImage(fileName string) (uint32, error) {
	imgFile, err := os.Open(fileName)
	if err != nil {
		return 0, err
	}
	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	if err != nil {
		return 0, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return 0, fmt.Errorf("unsupported stride")
	}

	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	gl.TexImage2D(
		gl.TEXTURE_2D, 0, gl.RGBA,
		int32(rgba.Rect.Size().X), int32(rgba.Rect.Size().Y), 0,
		gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix))

	return texture, nil
}

func NewTexture(fileName string, width, height float32) (texture Texture, err error) {
	texture.VertexObject = Rectangle(Point{0, 0}, width, height)
	texture.tex, err = loadImage(fileName)
	texture.width = width
	texture.height = height
	return
}

func (t Texture) bind() {
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, t.tex)
}

func (t Texture) Move(x, y float32) Texture {
	t.VertexObject = t.VertexObject.Move(x, y)
	return t
}

func (t Texture) CenterAt(point Point) Texture {
	t.VertexObject = t.VertexObject.CenterAt(point)
	return t
}

func (t Texture) RotateDeg(angle float64) Texture {
	t.VertexObject = t.VertexObject.RotateDeg(angle)
	return t
}

func (t Texture) RotateRad(angle float64) Texture {
	t.VertexObject = t.VertexObject.RotateRad(angle)
	return t
}

func (t Texture) Scale(scale float64) Texture {
	t.VertexObject = t.VertexObject.Scale(scale)
	return t
}
