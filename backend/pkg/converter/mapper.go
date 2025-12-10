package converter

import (
	"image"
	"image/color"
	"strings"
)

// Palette types
const (
	PaletteNormal  = "normal"
	PaletteDense   = "dense"
	PaletteSparse  = "sparse"
	PaletteUnicode = "unicode"
)

// GetPalette returns the character palette string for the given palette type
func GetPalette(paletteType string) string {
	switch paletteType {
	case PaletteDense:
		return ".oO0@#"
	case PaletteSparse:
		return " .'`^\""
	case PaletteUnicode:
		// Unicode block characters: light shade, medium shade, dark shade, full block
		// Using explicit runes to ensure proper UTF-8 encoding
		return string([]rune{'\u2591', '\u2592', '\u2593', '\u2588'}) // ░▒▓█
	case PaletteNormal:
		fallthrough
	default:
		return " .:-=+*#%@"
	}
}

func BrightnessToChar(brightness uint8, palette string) string {
	// palette should be the actual character palette string, not the type
	if len(palette) == 0 {
		palette = GetPalette(PaletteNormal)
	}

	// Convert string to rune slice to handle multi-byte Unicode characters correctly
	runes := []rune(palette)
	if len(runes) == 0 {
		return " "
	}

	index := float64(brightness) / 255.0 * float64(len(runes)-1)
	return string(runes[int(index)])
}

func ConvertToASCII(img image.Image, palette string) string {
	bounds := img.Bounds()
	// Convert palette type to actual character palette
	charPalette := GetPalette(palette)

	// strings.Builder is much more efficient than string concatenation
	// It preallocates memory and avoids creating new strings on each append
	var builder strings.Builder

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			brightness := img.At(x, y).(color.Gray).Y
			// WriteString appends to the builder without allocating new strings
			builder.WriteString(BrightnessToChar(brightness, charPalette))
		}
		// Add newline at the end of each row
		builder.WriteString("\n")
	}

	// Convert the builder's contents to a final string
	return builder.String()
}
