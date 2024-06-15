package main

import (
	"image"
	"image/color"
	"log"

	"github.com/RaghavSood/ogimage"
)

func main() {
	ogImage, err := ogimage.NewOgImage("template.png", "logo.png")
	if err != nil {
		log.Fatalf("failed to create ogimage: %v", err)
	}

	title := ogimage.Text{
		Content:  "1,232,232.12345678 BTC BURNED",
		FontFile: "menlo.ttf", // Specify your font file here
		FontSize: 64,
		Color:    color.White,
		Point:    image.Point{20, 305}, // Custom position for this text
	}

	subtitle := ogimage.Text{
		Content:  "Tranasction 1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		FontFile: "menlo.ttf", // Specify your font file here
		FontSize: 25,
		Color:    color.White,
		Point:    image.Point{20, 345}, // Custom position for this text
	}

	err = ogImage.GenerateDefault("output.png", title, subtitle, 10)
	if err != nil {
		log.Fatalf("failed to generate image: %v", err)
	}

	anotherText := ogimage.Text{
		Content:  "That's worth $1,234,567,890,239.12",
		FontFile: "menlo.ttf",
		FontSize: 30,
		Color:    color.White,
		Point:    image.Point{20, 450}, // Custom position for this text
	}

	config := ogimage.Config{
		Position: ogimage.BottomRight,
		Padding:  20,
		Texts:    []ogimage.Text{title, subtitle, anotherText},
	}

	err = ogImage.Generate("output2.png", config)
	if err != nil {
		log.Fatalf("failed to generate image: %v", err)
	}
}
