package converter

import (
	"fmt"
	"strings"
)

// ConvertToSVG converts ASCII art to an SVG image
func ConvertToSVG(asciiArt string, coloredASCII *ColoredASCII, fontSize int) string {
	var svg strings.Builder
	
	lines := strings.Split(strings.TrimRight(asciiArt, "\n"), "\n")
	if coloredASCII != nil {
		lines = make([]string, len(coloredASCII.Lines))
		for i, line := range coloredASCII.Lines {
			var lineStr strings.Builder
			for _, char := range line {
				lineStr.WriteString(char.Char)
			}
			lines[i] = lineStr.String()
		}
	}
	
	if len(lines) == 0 {
		return ""
	}
	
	// Find the maximum line width (in characters)
	maxWidth := 0
	for _, line := range lines {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}
	
	// Calculate dimensions based on actual content
	charWidth := fontSize * 6 / 10 // Approximate character width (0.6 * fontSize)
	width := maxWidth * charWidth
	height := len(lines) * fontSize
	
	svg.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" style="background:black;">`, width, height))
	svg.WriteString("\n")
	
	if coloredASCII != nil {
		// Render colored ASCII
		y := fontSize
		for _, line := range coloredASCII.Lines {
			x := 0
			for _, char := range line {
				rgb := fmt.Sprintf("rgb(%d,%d,%d)", char.R, char.G, char.B)
				svg.WriteString(fmt.Sprintf(`<text x="%d" y="%d" fill="%s" font-family="monospace" font-size="%d">%s</text>`, 
					x, y, rgb, fontSize, escapeXML(char.Char)))
				x += fontSize * 6 / 10
			}
			y += fontSize
		}
	} else {
		// Render grayscale ASCII
		y := fontSize
		for _, line := range lines {
			svg.WriteString(fmt.Sprintf(`<text x="0" y="%d" fill="white" font-family="monospace" font-size="%d">%s</text>`, 
				y, fontSize, escapeXML(line)))
			y += fontSize
		}
	}
	
	svg.WriteString("\n</svg>")
	return svg.String()
}

// Helper function to escape XML special characters
func escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}

