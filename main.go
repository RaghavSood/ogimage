package ogimage

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

type OgImage struct {
	Template image.Image
	Logo     image.Image
}

func NewOgImage(templateFile, logoFile string) (*OgImage, error) {
	template, err := loadImage(templateFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load template image: %v", err)
	}

	logo, err := loadImage(logoFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load logo image: %v", err)
	}

	return &OgImage{Template: template, Logo: logo}, nil
}

func loadImage(filePath string) (image.Image, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %v", filePath, err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image %s: %v", filePath, err)
	}

	return img, nil
}

func (og *OgImage) Generate(outputFile string) error {
	output := image.NewRGBA(og.Template.Bounds())
	draw.Draw(output, og.Template.Bounds(), og.Template, image.Point{}, draw.Over)
	draw.Draw(output, og.Logo.Bounds().Add(image.Point{10, 10}), og.Logo, image.Point{}, draw.Over)

	outFile, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file %s: %v", outputFile, err)
	}
	defer outFile.Close()

	return jpeg.Encode(outFile, output, nil)
}
