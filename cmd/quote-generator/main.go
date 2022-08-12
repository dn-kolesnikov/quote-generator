package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"io"
	"log"
	"net/http"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

//go:embed assets/bg+logo.png
var bgLogo []byte

//go:embed assets/UbuntuMono-R.ttf
var mainFont []byte

// default font size in inch
const fontSize = 62.0

func run() error {

	text, err := getForismaticQuote()
	if err != nil {
		return err
	}

	img, err := putTextToTemplateImage(text, fontSize)
	if err != nil {
		return err
	}

	if err := gg.SavePNG("out.png", img); err != nil {
		return err
	}

	return nil

}

func getForismaticQuote() (quote string, err error) {
	const url = "https://api.forismatic.com/api/1.0/?method=getQuote&format=json&lang=ru"

	type ForismaticQuote struct {
		QuoteText   string `json:"quoteText"`
		QuoteAuthor string `json:"quoteAuthor"`
	}

	var fq ForismaticQuote

	resp, err := http.Get(url)
	if resp != nil {
		defer func() {
			if _err := resp.Body.Close(); _err != nil {
				err = _err
			}
		}()
	}
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if err := json.Unmarshal(body, &fq); err != nil {
		return "", err
	}

	quote = fmt.Sprintf("%q\n%s\n", fq.QuoteText, fq.QuoteAuthor)

	return quote, nil
}

//
//
//
func putTextToTemplateImage(text string, fontSize float64) (image.Image, error) {

	img, _, err := image.Decode(bytes.NewReader(bgLogo))
	if err != nil {
		return nil, err
	}

	imgWidth := img.Bounds().Dx()
	imgHeight := img.Bounds().Dy()

	dc := gg.NewContext(imgWidth, imgHeight)
	dc.DrawImage(img, 0, 0)

	fnt, err := truetype.Parse(mainFont)
	if err != nil {
		return nil, err
	}
	face := truetype.NewFace(fnt, &truetype.Options{
		Size: fontSize,
	})

	dc.SetFontFace(face)

	x := float64(imgWidth / 2)
	y := float64(imgHeight / 2)
	maxWidth := float64(imgWidth - 10)

	dc.SetColor(color.Black)
	dc.DrawStringWrapped(text, x+3, y+3, 0.5, 0.5, maxWidth, 1.5, gg.AlignCenter)
	dc.SetColor(color.White)
	dc.DrawStringWrapped(text, x, y, 0.5, 0.5, maxWidth, 1.5, gg.AlignCenter)

	return dc.Image(), nil

}
