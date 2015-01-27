package annotator

import (
	"bytes"
	"image"
	"image/draw"
	"image/jpeg"

	"code.google.com/p/freetype-go/freetype"
	"code.google.com/p/freetype-go/freetype/raster"
	"code.google.com/p/freetype-go/freetype/truetype"
	"github.com/lucasb-eyer/go-colorful"
)

type Annotator struct {
	Font      *truetype.Font
	SrcFile   []byte
	FontSize  float64
	FontColor string
}

func (a Annotator) Annotate(upperText string, lowerText string) []byte {
	fontColor, _ := colorful.Hex(a.FontColor)
	fontMask := image.NewUniform(fontColor)

	srcImage, _ := jpeg.Decode(bytes.NewBuffer(a.SrcFile))
	srcBounds := srcImage.Bounds()
	srcHeight := srcBounds.Max.Y

	imageMask := image.NewRGBA(srcBounds)
	imageMaskClip := imageMask.Bounds()
	draw.Draw(imageMask, imageMaskClip, srcImage, image.ZP, draw.Src)

	c := freetype.NewContext()
	c.SetFont(a.Font)
	c.SetFontSize(a.FontSize)
	c.SetDPI(a.FontSize)
	c.SetSrc(fontMask)
	c.SetDst(imageMask)
	c.SetClip(imageMaskClip)

	drawLines := func(line string, point raster.Point) {
		c.DrawString(line, point)
	}

	upperPoint := freetype.Pt(10, 40)
	lowerPoint := freetype.Pt(10, srcHeight-10)

	drawLines(upperText, upperPoint)
	drawLines(lowerText, lowerPoint)

	dataBuffer := bytes.NewBuffer([]byte(""))
	jpeg.Encode(dataBuffer, imageMask, nil)

	return dataBuffer.Bytes()
}
