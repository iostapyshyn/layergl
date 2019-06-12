package layergl

import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/draw"
	"image/jpeg"
	"io"
	"io/ioutil"
	"os"
)

const maxchar = 128

type Font struct {
	char     []*character
	vao, vbo uint32
	texture  uint32
}

type character struct {
	tex    uint32
	w, h   int32
	adv    int32
	bH, bV int32
}

func loadFont(r io.Reader, scale int32) (*Font, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	ttf, err := truetype.Parse(data)
	if err != nil {
		return nil, err
	}

	f := new(Font)
	f.char = make([]*character, 0, maxchar)

	for ch := 0; ch < maxchar; ch++ {
		char := new(character)

		ttfFace := truetype.NewFace(ttf, &truetype.Options{
			Size:    float64(scale),
			DPI:     72,
			Hinting: font.HintingFull,
		})

		gBnd, gAdv, ok := ttfFace.GlyphBounds(rune(ch))
		if ok != true {
			return nil, fmt.Errorf("The face does not contain a glyph for %v.", rune(ch))
		}

		gw := int32((gBnd.Max.X - gBnd.Min.X) >> 6)
		gh := int32((gBnd.Max.Y - gBnd.Min.Y) >> 6)

		if gw == 0 || gh == 0 {
			gBnd := ttf.Bounds(fixed.Int26_6(scale))
			gw := int32((gBnd.Max.X - gBnd.Min.X) >> 6)
			gh := int32((gBnd.Max.Y - gBnd.Min.Y) >> 6)

			if gw == 0 || gh == 0 {
				gw = 1
				gh = 1
			}
		}

		gDescent := -(int32(gBnd.Min.Y) >> 6)
		gAscent := (int32(gBnd.Max.Y) >> 6)

		char.w = gw
		char.h = gh
		char.adv = int32(gAdv)
		char.bV = gAscent
		char.bH = (int32(gBnd.Min.X) >> 6)

		// Create and image to draw glyph
		fg, bg := image.White, image.Black
		rect := image.Rect(0, 0, int(gw), int(gh))
		rgba := image.NewRGBA(rect)
		draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)

		c := freetype.NewContext()
		c.SetDPI(72)
		c.SetFont(ttf)
		c.SetFontSize(float64(scale))
		c.SetClip(rgba.Bounds())
		c.SetDst(rgba)
		c.SetSrc(fg)
		c.SetHinting(font.HintingFull)

		// set point
		px := 0 - (int(gBnd.Min.X) >> 6)
		py := int(gDescent)
		pt := freetype.Pt(px, py)

		_, err := c.DrawString(string(rune(ch)), pt)
		if err != nil {
			return nil, err
		}

		out, err := os.Create("./output.jpg")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		var opt jpeg.Options

		opt.Quality = 80
		// ok, write out the data into the new JPEG file

		err = jpeg.Encode(out, rgba, &opt) // put quality to 80%
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		var texture uint32
		gl.GenTextures(1, &texture)
		gl.BindTexture(gl.TEXTURE_2D, texture)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		//		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(rgba.Rect.Dx()), int32(rgba.Rect.Dy()), 0,
		//			gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix))
		gl.TexImage2D(
			gl.TEXTURE_2D, 0, gl.RGBA,
			int32(rgba.Rect.Size().X), int32(rgba.Rect.Size().Y), 0,
			gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix))

		char.tex = texture

		f.char = append(f.char, char)
	}

	return f, nil
}

func LoadFont(file string, scale int32) (*Font, error) {
	fd, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	return loadFont(fd, scale)
}
