package converter

import (
	"image"

	"github.com/nfnt/resize"
)

// ResizeImage resizes an image to the target width while maintaining aspect ratio.
// It also applies character aspect ratio correction (characters are ~2x taller than wide).
//
// Parameters:
//   - img: The source image to resize
//   - targetWidth: Desired width in characters (e.g., 80, 100, 120)
//
// Returns:
//   - A resized image ready for ASCII conversion
func ResizeImage(img image.Image, targetWidth int) image.Image {
	// Get original dimensions
	bounds := img.Bounds()
	originalWidth := bounds.Max.X - bounds.Min.X
	originalHeight := bounds.Max.Y - bounds.Min.Y

	// Calculate the scale factor based on target width
	scale := float64(targetWidth) / float64(originalWidth)

	// Calculate new height maintaining aspect ratio
	newHeight := int(float64(originalHeight) * scale)

	// Apply character aspect ratio correction
	// Terminal characters are roughly 2x taller than wide
	// So we reduce height by ~50% to prevent vertical stretching
	aspectRatio := 0.5
	newHeight = int(float64(newHeight) * aspectRatio)

	// Use the resize library with Lanczos3 interpolation
	// Lanczos3 provides high-quality results, good for downscaling
	// Other options: NearestNeighbor (fastest), Bilinear, Bicubic
	resizedImg := resize.Resize(uint(targetWidth), uint(newHeight), img, resize.Lanczos3)

	return resizedImg
}
