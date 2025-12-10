package converter

import (
	"image"
	"image/color"
	"strings"
)

func BrightnessToChar(brightness uint8) string {
	pallete := " .:-=+*#%@"
	index := float64(brightness) / 255.0 * float64(len(pallete)-1)
	return string(pallete[int(index)])
}

func ConvertToASCII(img image.Image) string {
	bounds := img.Bounds()

	// strings.Builder is much more efficient than string concatenation
	// It preallocates memory and avoids creating new strings on each append
	var builder strings.Builder

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			brightness := img.At(x, y).(color.Gray).Y
			// WriteString appends to the builder without allocating new strings
			builder.WriteString(BrightnessToChar(brightness))
		}
		// Add newline at the end of each row
		builder.WriteString("\n")
	}

	// Convert the builder's contents to a final string
	return builder.String()
}
