package ogimage

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
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

type Text struct {
	Content  string
	FontFile string
	FontFace font.Face
	FontSize float64
	Color    color.Color
}

type Config struct {
	Position LogoPosition
	Padding  int
	Title    Text
	Subtitle Text
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

func loadFont(filePath string, size float64) (font.Face, error) {
	fontBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	f, err := opentype.Parse(fontBytes)
	if err != nil {
		return nil, err
	}

	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size: size,
		DPI:  72,
	})
	if err != nil {
		return nil, err
	}

	return face, nil
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

	// Calculate positions for the title and subtitle
	titlePosition := image.Point{padding, templateBounds.Max.Y/2 - 20}
	subtitlePosition := image.Point{padding, templateBounds.Max.Y/2 + 20}

	err := drawText(output, config.Title, titlePosition)
	if err != nil {
		return err
	}

	err = drawText(output, config.Subtitle, subtitlePosition)
	if err != nil {
		return err
	}

	outFile, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer outFile.Close()

	b := bufio.NewWriter(outFile)
	err = jpeg.Encode(b, output, nil)
	if err != nil {
		return err
	}

	return b.Flush()
}

func drawText(img *image.RGBA, text Text, position image.Point) error {
	if text.FontFace == nil && text.FontFile != "" {
		face, err := loadFont(text.FontFile, text.FontSize)
		if err != nil {
			return err
		}
		text.FontFace = face
	} else if text.FontFace == nil {
		face, err := loadFont("path/to/default/font.ttf", text.FontSize)
		if err != nil {
			return err
		}
		text.FontFace = face
	}

	col := text.Color
	if col == nil {
		col = color.Black
	}

	point := fixed.Point26_6{
		X: fixed.I(position.X),
		Y: fixed.I(position.Y),
	}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: text.FontFace,
		Dot:  point,
	}
	d.DrawString(text.Content)
	return nil
}
