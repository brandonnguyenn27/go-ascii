package converter

import (
	"image"
	"image/color"
)

func RGBToGrayScale(r, g, b uint32) uint8 {
	r = r >> 8
	g = g >> 8
	b = b >> 8
	return uint8((0.299 * float64(r)) + (0.587 * float64(g)) + (0.114 * float64(b)))
}

func ConvertToGrayscale(img image.Image) image.Image {
	bounds := img.Bounds()
	grayImg := image.NewGray(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			brightness := RGBToGrayScale(r, g, b)
			grayImg.Set(x, y, color.Gray{Y: brightness})
		}
	}
	return grayImg
}
