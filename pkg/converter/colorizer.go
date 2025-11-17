package converter

import (
	"fmt"
	"image"
	"strings"
)

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
