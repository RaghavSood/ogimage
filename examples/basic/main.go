package main

import (
	"image"
	"image/color"
	"io/ioutil"
	"log"

	"github.com/RaghavSood/ogimage"
)

func main() {
	templateData, err := ioutil.ReadFile("template.png")
	if err != nil {
		log.Fatalf("failed to read template file: %v", err)
	}

	logoData, err := ioutil.ReadFile("logo.png")
	if err != nil {
		log.Fatalf("failed to read logo file: %v", err)
	}

	fontData, err := ioutil.ReadFile("menlo.ttf")
	if err != nil {
		log.Fatalf("failed to read font file: %v", err)
	}

	ogImage, err := ogimage.NewOgImage(templateData, logoData)
	if err != nil {
		log.Fatalf("failed to create ogimage: %v", err)
	}

	title := ogimage.Text{
		Content:  "1,232,232.12345678 BTC BURNED",
		FontData: fontData,
		FontSize: 64,
		Color:    color.White,
		Point:    image.Point{20, 305}, // Custom position for this text
	}

	subtitle := ogimage.Text{
		Content:  "Tranasction 1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		FontData: fontData,
		FontSize: 25,
		Color:    color.White,
		Point:    image.Point{20, 345}, // Custom position for this text
	}

	imageData, err := ogImage.GenerateDefault(title, subtitle, 10)
	if err != nil {
		log.Fatalf("failed to generate image: %v", err)
	}

	err = ioutil.WriteFile("output.png", imageData, 0644)
	if err != nil {
		log.Fatalf("failed to write output file: %v", err)
	}

	anotherText := ogimage.Text{
		Content:  "That's worth $1,234,567,890,239.12",
		FontData: fontData,
		FontSize: 30,
		Color:    color.White,
		Point:    image.Point{20, 450}, // Custom position for this text
	}

	config := ogimage.Config{
		Position: ogimage.BottomRight,
		Padding:  20,
		Texts:    []ogimage.Text{title, subtitle, anotherText},
	}

	imageDataMultiple, err := ogImage.Generate(config)
	if err != nil {
		log.Fatalf("failed to generate image: %v", err)
	}

	err = ioutil.WriteFile("output_multiple.png", imageDataMultiple, 0644)
	if err != nil {
		log.Fatalf("failed to write output file: %v", err)
	}
}
