package main

import (
	"log"

	"github.com/RaghavSood/ogimage"
)

func main() {
	ogImage, err := ogimage.NewOgImage("template.png", "logo.png")
	if err != nil {
		log.Fatalf("failed to create ogimage: %v", err)
	}

	err = ogImage.Generate("output.jpg")
	if err != nil {
		log.Fatalf("failed to generate image: %v", err)
	}
}
