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

type LogoPosition int

const (
	TopLeft LogoPosition = iota
	TopRight
	BottomLeft
	BottomRight
	Center
)

type Config struct {
	Position LogoPosition
	Padding  int
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

func (og *OgImage) Generate(outputFile string, config Config) error {
	output := image.NewRGBA(og.Template.Bounds())
	draw.Draw(output, og.Template.Bounds(), og.Template, image.Point{}, draw.Over)

	logoBounds := og.Logo.Bounds()
	templateBounds := og.Template.Bounds()

	padding := config.Padding
	if padding < 0 {
		padding = 0
	}

	var position image.Point
	switch config.Position {
	case TopLeft:
		position = image.Point{padding, padding}
	case TopRight:
		position = image.Point{templateBounds.Max.X - logoBounds.Max.X - padding, padding}
	case BottomLeft:
		position = image.Point{padding, templateBounds.Max.Y - logoBounds.Max.Y - padding}
	case BottomRight:
		position = image.Point{templateBounds.Max.X - logoBounds.Max.X - padding, templateBounds.Max.Y - logoBounds.Max.Y - padding}
	case Center:
		position = image.Point{(templateBounds.Max.X - logoBounds.Max.X) / 2, (templateBounds.Max.Y - logoBounds.Max.Y) / 2}
	}

	draw.Draw(output, logoBounds.Add(position), og.Logo, image.Point{}, draw.Over)

	outFile, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer outFile.Close()

	return jpeg.Encode(outFile, output, nil)
}
