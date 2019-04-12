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
	*VertexObject
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

// Loads and creates new Texture object.
func NewTexture(fileName string, width, height float64) (texture *Texture, err error) {
	texture = new(Texture)
	texture.VertexObject = Rectangle(Rect{0, 0, width, height})
	texture.tex, err = loadImage(fileName)
	texture.width = float32(width)
	texture.height = float32(height)
	return
}

func (t *Texture) bind() {
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, t.tex)
}
