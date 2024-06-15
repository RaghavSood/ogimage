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
		Content:  "Sample Title",
		FontFile: "menlo.ttf", // Specify your font file here
		FontSize: 40,
		Color:    color.RGBA{255, 0, 0, 255},
	}

	subtitle := ogimage.Text{
		Content:  "Sample Subtitle",
		FontFile: "menlo.ttf", // Specify your font file here
		FontSize: 20,
		Color:    color.RGBA{0, 0, 255, 255},
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
