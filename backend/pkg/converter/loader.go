package converter

import (
	"fmt"
	"image"
	"io"

	// Import image format decoders
	// The _ means "import for side effects only"
	// These packages register their decoders with the image package
	_ "image/jpeg"
	_ "image/png"
	"os"
)

// LoadImage reads an image file from the given path and decodes it.
// It returns the decoded image and any error encountered.
// Supported formats: JPEG, PNG (can add GIF, WebP, etc. by importing their packages)
func LoadImage(filePath string) (image.Image, error) {
	// Step 1: Open the file
	// os.Open returns a *os.File which implements io.Reader
	file, err := os.Open(filePath)
	if err != nil {
		// Wrap the error with context about what we were trying to do
		return nil, fmt.Errorf("failed to open image file: %w", err)
	}
	// defer ensures the file is closed when the function returns
	// This happens even if we return early due to an error
	defer file.Close()

	// Step 2: Decode the image
	// image.Decode automatically detects the format by reading the file header
	// It then uses the appropriate registered decoder (JPEG, PNG, etc.)
	img, format, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// Step 3: Optional - log what we decoded (useful for debugging)
	// In production, you might want to use a proper logger
	fmt.Printf("Successfully loaded %s image\n", format)

	// Step 4: Return the decoded image
	// img is of type image.Image (an interface)
	// The actual concrete type depends on the format (e.g., *image.YCbCr for JPEG)
	return img, nil
}

// LoadImageFromReader reads an image from an io.Reader and decodes it.
// It returns the decoded image and any error encountered.
// Supported formats: JPEG, PNG (can add GIF, WebP, etc. by importing their packages)
// This function is useful for API endpoints that receive image data via HTTP requests.
func LoadImageFromReader(reader io.Reader) (image.Image, error) {
	// Decode the image from the reader
	// image.Decode automatically detects the format by reading the file header
	// It then uses the appropriate registered decoder (JPEG, PNG, etc.)
	img, format, err := image.Decode(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// Optional - log what we decoded (useful for debugging)
	fmt.Printf("Successfully loaded %s image\n", format)

	// Return the decoded image
	return img, nil
}
