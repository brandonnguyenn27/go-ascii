package converter

import (
	"fmt"
	"image"
	"strings"
)

// ColoredChar represents a single character with its RGB color
type ColoredChar struct {
	Char string `json:"char"`
	R    uint8  `json:"r"`
	G    uint8  `json:"g"`
	B    uint8  `json:"b"`
}

// ColoredASCII represents ASCII art with color information as structured data
// Each line is an array of ColoredChar objects
type ColoredASCII struct {
	Lines [][]ColoredChar `json:"lines"`
}

func RGBToANSI(r, g, b uint8) string {
	ansi := fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)
	return ansi
}

func ConvertToASCIIWithColor(img image.Image) string {
	bounds := img.Bounds()
	var builder strings.Builder

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			r8, g8, b8 := uint8(r>>8), uint8(g>>8), uint8(b>>8)

			brightness := RGBToGrayScale(r, g, b)
			char := BrightnessToChar(brightness)

			// Write colored character
			builder.WriteString(fmt.Sprintf("\033[38;2;%d;%d;%dm%s", r8, g8, b8, char))
		}
		// Reset color at end of line
		builder.WriteString("\033[0m\n")
	}

	return builder.String()
}

// ConvertToASCIIWithColorStructured converts an image to ASCII art with color information
// Returns structured data suitable for JSON serialization (for API responses)
func ConvertToASCIIWithColorStructured(img image.Image) ColoredASCII {
	bounds := img.Bounds()
	lines := make([][]ColoredChar, 0, bounds.Max.Y-bounds.Min.Y)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		line := make([]ColoredChar, 0, bounds.Max.X-bounds.Min.X)
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			r8, g8, b8 := uint8(r>>8), uint8(g>>8), uint8(b>>8)

			brightness := RGBToGrayScale(r, g, b)
			char := BrightnessToChar(brightness)

			line = append(line, ColoredChar{
				Char: char,
				R:    r8,
				G:    g8,
				B:    b8,
			})
		}
		lines = append(lines, line)
	}

	return ColoredASCII{Lines: lines}
}
