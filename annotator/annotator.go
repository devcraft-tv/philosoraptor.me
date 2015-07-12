package annotator

import (
	"bytes"
	"image"
	"image/draw"
	"image/jpeg"

	"code.google.com/p/freetype-go/freetype"
	"code.google.com/p/freetype-go/freetype/raster"
	"code.google.com/p/freetype-go/freetype/truetype"
	"github.com/devcraft-tv/philosoraptor/line_breaker"
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

	margin := 10
	leftPoint := margin
	topPoint := margin

	lineHeight := int(a.FontSize)
	simpleLineBreaker := line_breaker.SimpleLineBreaker{}

	upperLines := simpleLineBreaker.GetLines(upperText, 30)
	for idx, line := range upperLines {
		point := freetype.Pt(leftPoint, topPoint+(lineHeight*(idx+1)))
		drawLines(c, line, point)
	}

	lowerLines := simpleLineBreaker.GetLines(lowerText, 30)
	lowerTopPoint := srcHeight - margin - len(lowerLines)*lineHeight
	for idx, line := range lowerLines {
		point := freetype.Pt(leftPoint, lowerTopPoint+(lineHeight*(idx+1)))
		drawLines(c, line, point)
	}

	dataBuffer := bytes.NewBuffer([]byte(""))
	jpeg.Encode(dataBuffer, imageMask, nil)

	return dataBuffer.Bytes()
}

func drawLines(context *freetype.Context, line string, point raster.Point) {
	context.DrawString(line, point)
}
