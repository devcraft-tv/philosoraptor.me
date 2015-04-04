package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"text/template"

	"crypto/md5"

	"code.google.com/p/freetype-go/freetype"
	"code.google.com/p/freetype-go/freetype/truetype"
	"github.com/devcraft-tv/philosoraptor/annotator"
	"github.com/devcraft-tv/philosoraptor/cache"
	"github.com/devcraft-tv/philosoraptor/uploader"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var font *truetype.Font
var templateFile []byte
var fileUploader *uploader.S3Uploader

type Cacher interface {
	Get(string) (string, error)
	Set(string, string) (bool, error)
}

var cacheDb Cacher

var htmlTemplates *template.Template

func main() {
	router := mux.NewRouter()
	htmlTemplates = template.Must(template.ParseGlob("templates/*"))
	font = readFont()
	templateFile = readTemplateFile()
	loadEnvVars()
	setUpS3Uploader()
	setUpCache()

	router.HandleFunc("/", homePage)
	router.HandleFunc("/generate", handleForm).Methods("POST")
	router.PathPrefix("/assets/").Handler(staticHandler())
	http.Handle("/", router)
	http.ListenAndServe(":8001", nil)
}

func loadEnvVars() {
	godotenv.Load()
}

func setUpCache() {
	connectionString := os.Getenv("REDIS_CONNECTION_STRING")
	password := os.Getenv("REDIS_PASSWORD")

	var err error
	cacheDb, err = cache.NewCache(connectionString, password)
	if err != nil {
		panic(err)
	}
}

func setUpS3Uploader() {
	fileUploader = &uploader.S3Uploader{
		AccessId:        os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		Bucket:          os.Getenv("S3_BUCKET"),
		Path:            os.Getenv("S3_BUCKET_PATH"),
	}
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

func readTemplateFile() []byte {
	templateFile, err := ioutil.ReadFile("./static/template.jpg")
	if err != nil {
		panic(err)
	}
	return templateFile
}

func staticHandler() http.Handler {
	return http.FileServer(http.Dir("static/"))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Println(os.Getenv("TEST_VARIABLE"))
	err := htmlTemplates.ExecuteTemplate(w, "index", nil)
	if err != nil {
		panic(err)
	}
}

func handleForm(w http.ResponseWriter, r *http.Request) {
	upperText := r.FormValue("upper_text")
	lowerText := r.FormValue("lower_text")
	annotator := newAnnotator()
	var err error

	fileName := hash(upperText + lowerText)
	url, err := cacheDb.Get(fileName)

	if err != nil {
		panic(err)
	}

	if url == "" {
		imageData := annotator.Annotate(upperText, lowerText)
		url, err = fileUploader.Upload(imageData, fileName)

		if err == nil {
			_, err = cacheDb.Set(fileName, url)
		}

		if err != nil {
			panic(err)
		}
	}

	http.Redirect(w, r, url, 301)
}

func hash(text string) (hashedText string) {
	t := md5.Sum([]byte(text))
	hashedText = hex.EncodeToString(t[:])
	return
}

func newAnnotator() annotator.Annotator {
	return annotator.Annotator{
		Font:      font,
		SrcFile:   templateFile,
		FontSize:  60,
		FontColor: "#000000",
	}
}
