package main

import (
	"bytes"
	"image"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"net/http"
	"text/template"

	"code.google.com/p/freetype-go/freetype"
	"code.google.com/p/freetype-go/freetype/raster"
	"code.google.com/p/freetype-go/freetype/truetype"
	"github.com/gorilla/mux"
	"github.com/lucasb-eyer/go-colorful"
)

var font *truetype.Font

var htmlTemplates *template.Template

func main() {
	router := mux.NewRouter()
	htmlTemplates = template.Must(template.ParseGlob("templates/*"))
	font = readFont()

	router.HandleFunc("/", homePage)
	router.HandleFunc("/generate", handleForm).Methods("POST")
	router.PathPrefix("/assets/").Handler(staticHandler())
	http.Handle("/", router)
	http.ListenAndServe(":8001", nil)
}

func readFont() *truetype.Font {
	rawFont, err := ioutil.ReadFile("./static/font.ttf")
	if err != nil {
		panic(err)
	}
	parsedFont, err := freetype.ParseFont(rawFont)
	if err != nil {
		panic(err)
	}
	return parsedFont
}

func staticHandler() http.Handler {
	return http.FileServer(http.Dir("static/"))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	err := htmlTemplates.ExecuteTemplate(w, "index", nil)
	if err != nil {
		panic(err)
	}
}

func handleForm(w http.ResponseWriter, r *http.Request) {
	upperText := r.FormValue("upper_text")
	lowerText := r.FormValue("lower_text")
	annotator := Annotator{}
	imageData := annotator.annotate(upperText, lowerText)

	w.Header().Set("Content-Type", "image/jpeg")
	w.Write([]byte(imageData))
}

type Annotator struct{}

func (a Annotator) annotate(upperText string, lowerText string) []byte {
	fontColor, _ := colorful.Hex("#000000")
	fontMask := image.NewUniform(fontColor)

	templateFile, _ := ioutil.ReadFile("./static/template.jpg")
	srcImage, _ := jpeg.Decode(bytes.NewBuffer(templateFile))
	srcBounds := srcImage.Bounds()
	srcHeight := srcBounds.Max.Y

	imageMask := image.NewRGBA(srcBounds)
	imageMaskClip := imageMask.Bounds()
	draw.Draw(imageMask, imageMaskClip, srcImage, image.ZP, draw.Src)

	c := freetype.NewContext()
	c.SetFont(font)
	c.SetFontSize(60)
	c.SetDPI(60)
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
