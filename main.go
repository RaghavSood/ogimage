package ogimage

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	_ "image/png"

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
	FontData []byte
	FontFace font.Face
	FontSize float64
	Color    color.Color
	Point    image.Point
}

type Config struct {
	Position LogoPosition
	Padding  int
	Texts    []Text
}

func NewOgImage(templateData, logoData []byte) (*OgImage, error) {
	template, err := decodeImage(templateData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode template image: %w", err)
	}

	logo, err := decodeImage(logoData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode logo image: %w", err)
	}

	return &OgImage{Template: template, Logo: logo}, nil
}

func decodeImage(data []byte) (image.Image, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}
	return img, nil
}

func loadFont(fontData []byte, size float64) (font.Face, error) {
	f, err := opentype.Parse(fontData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse font data: %w", err)
	}

	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size: size,
		DPI:  72,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create font face: %w", err)
	}

	return face, nil
}

func (og *OgImage) Generate(config Config) ([]byte, error) {
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

	for _, text := range config.Texts {
		err := drawText(output, text)
		if err != nil {
			return nil, fmt.Errorf("failed to draw text: %w", err)
		}
	}

	var buf bytes.Buffer
	b := bufio.NewWriter(&buf)
	err := png.Encode(b, output)
	if err != nil {
		return nil, fmt.Errorf("failed to encode image: %w", err)
	}

	err = b.Flush()
	if err != nil {
		return nil, fmt.Errorf("failed to flush buffer: %w", err)
	}

	return buf.Bytes(), nil
}

func (og *OgImage) GenerateDefault(title, subtitle Text, padding int) ([]byte, error) {
	templateBounds := og.Template.Bounds()

	titlePosition := image.Point{padding, templateBounds.Max.Y/2 - 20}
	subtitlePosition := image.Point{padding, templateBounds.Max.Y/2 + 20}

	title.Point = titlePosition
	subtitle.Point = subtitlePosition

	config := Config{
		Position: BottomRight,
		Padding:  padding,
		Texts:    []Text{title, subtitle},
	}

	return og.Generate(config)
}

func drawText(img *image.RGBA, text Text) error {
	if text.FontFace == nil && text.FontData != nil {
		face, err := loadFont(text.FontData, text.FontSize)
		if err != nil {
			return fmt.Errorf("failed to load font: %w", err)
		}
		text.FontFace = face
	} else if text.FontFace == nil {
		return nil // No font face provided and no default specified
	}

	col := text.Color
	if col == nil {
		col = color.Black
	}

	point := fixed.Point26_6{
		X: fixed.I(text.Point.X),
		Y: fixed.I(text.Point.Y),
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
