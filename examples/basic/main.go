package main

import (
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
	}

	subtitle := ogimage.Text{
		Content:  "Tranasction 1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		FontFile: "menlo.ttf", // Specify your font file here
		FontSize: 25,
		Color:    color.White,
	}

	config := ogimage.Config{
		Position: ogimage.BottomRight,
		Padding:  10,
		Title:    title,
		Subtitle: subtitle,
	}

	err = ogImage.Generate("output.png", config)
	if err != nil {
		log.Fatalf("failed to generate image: %v", err)
	}
}
