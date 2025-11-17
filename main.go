package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/brandonnguyenn27/ascii-converter/pkg/converter"
)

func main() {
	// Define flags
	useColor := flag.Bool("color", false, "Enable colored ASCII output")
	width := flag.Int("width", 100, "Width of ASCII output in characters")

	// Parse flags
	flag.Parse()

	// Check if user provided an image path (after flags)
	if flag.NArg() < 1 {
		fmt.Println("Usage: go run main.go [flags] <image-path>")
		fmt.Println("Flags:")
		flag.PrintDefaults()
		fmt.Println("\nExample: go run main.go -color -width 120 images/apple.png")
		os.Exit(1)
	}

	// Get the image path (first non-flag argument)
	imagePath := flag.Arg(0)

	// Load the image
	img, err := converter.LoadImage(imagePath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Resize the image
	resizedImg := converter.ResizeImage(img, *width)

	// Convert to ASCII (color or grayscale)
	var asciiImg string
	if *useColor {
		asciiImg = converter.ConvertToASCIIWithColor(resizedImg)
	} else {
		grayScaleImg := converter.ConvertToGrayscale(resizedImg)
		asciiImg = converter.ConvertToASCII(grayScaleImg)
	}

	// Output the ASCII art
	fmt.Println(asciiImg)
}
