package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/brandonnguyenn27/ascii-converter/pkg/converter"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// Define all flags
	serverMode := flag.Bool("server", false, "Start the REST API server")
	useColor := flag.Bool("color", false, "Enable colored ASCII output")
	width := flag.Int("width", 100, "Width of ASCII output in characters")
	palette := flag.String("palette", "normal", "Character palette: normal, dense, sparse, or unicode")

	flag.Parse()

	if *serverMode {
		startServer()
	} else {
		runCLI(*useColor, *width, *palette)
	}
}

func startServer() {
	app := fiber.New(fiber.Config{
		BodyLimit: 100 * 1024 * 1024, // 100MB max body size (for video uploads)
	})

	// Configure CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
		AllowMethods: "POST,OPTIONS",
		AllowHeaders: "Content-Type",
	}))

	app.Post("/convert", convertHandler)            // Grayscale ASCII (returns string)
	app.Post("/convert/color", convertColorHandler) // Colored ASCII (returns structured data)
	app.Post("/convert/video", convertVideoHandler) // Video to ASCII (returns frames array)
	app.Post("/export/svg", exportSVGHandler)       // Export ASCII as SVG

	log.Println("Server starting on :3000")
	log.Fatal(app.Listen(":3000"))
}

func convertHandler(c *fiber.Ctx) error {
	// Get the uploaded file
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing or invalid image file. Please upload an image using the 'image' field.",
		})
	}

	// Get file size
	fileSize := file.Size

	// Get optional width parameter (default: 100)
	width := 100
	if widthStr := c.FormValue("width"); widthStr != "" {
		if parsedWidth, err := strconv.Atoi(widthStr); err == nil && parsedWidth > 0 {
			width = parsedWidth
		}
	} else if widthStr := c.Query("width"); widthStr != "" {
		if parsedWidth, err := strconv.Atoi(widthStr); err == nil && parsedWidth > 0 {
			width = parsedWidth
		}
	}

	// Get optional palette parameter (default: normal)
	palette := c.FormValue("palette")
	if palette == "" {
		palette = c.Query("palette")
	}
	if palette == "" {
		palette = converter.PaletteNormal
	}

	// Open the uploaded file
	fileHeader, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to open uploaded file",
		})
	}
	defer fileHeader.Close()

	// Load the image from the reader
	img, err := converter.LoadImageFromReader(fileHeader)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Get original dimensions
	originalBounds := img.Bounds()
	originalWidth := originalBounds.Max.X - originalBounds.Min.X
	originalHeight := originalBounds.Max.Y - originalBounds.Min.Y

	// Resize the image
	resizedImg := converter.ResizeImage(img, width)

	// Convert to grayscale ASCII
	grayScaleImg := converter.ConvertToGrayscale(resizedImg)
	asciiImg := converter.ConvertToASCII(grayScaleImg, palette)

	// Calculate ASCII size in bytes
	// len() returns byte length, which correctly accounts for:
	// - ASCII palettes (normal/dense/sparse): 1 byte per character
	// - Unicode palette: 3 bytes per character (multi-byte UTF-8)
	asciiSize := len(asciiImg)

	// Return JSON response
	return c.JSON(fiber.Map{
		"ascii":          asciiImg,
		"originalSize":   fileSize,
		"originalWidth":  originalWidth,
		"originalHeight": originalHeight,
		"asciiSize":      asciiSize,
	})
}

func convertColorHandler(c *fiber.Ctx) error {
	// Get the uploaded file
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing or invalid image file. Please upload an image using the 'image' field.",
		})
	}

	// Get file size
	fileSize := file.Size

	// Get optional width parameter (default: 100)
	width := 100
	if widthStr := c.FormValue("width"); widthStr != "" {
		if parsedWidth, err := strconv.Atoi(widthStr); err == nil && parsedWidth > 0 {
			width = parsedWidth
		}
	} else if widthStr := c.Query("width"); widthStr != "" {
		if parsedWidth, err := strconv.Atoi(widthStr); err == nil && parsedWidth > 0 {
			width = parsedWidth
		}
	}

	// Get optional palette parameter (default: normal)
	palette := c.FormValue("palette")
	if palette == "" {
		palette = c.Query("palette")
	}
	if palette == "" {
		palette = converter.PaletteNormal
	}

	// Open the uploaded file
	fileHeader, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to open uploaded file",
		})
	}
	defer fileHeader.Close()

	// Load the image from the reader
	img, err := converter.LoadImageFromReader(fileHeader)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Get original dimensions
	originalBounds := img.Bounds()
	originalWidth := originalBounds.Max.X - originalBounds.Min.X
	originalHeight := originalBounds.Max.Y - originalBounds.Min.Y

	// Resize the image
	resizedImg := converter.ResizeImage(img, width)

	// Convert to colored ASCII with structured data
	coloredASCII := converter.ConvertToASCIIWithColorStructured(resizedImg, palette)

	// Calculate ASCII size by converting to JSON and measuring byte length
	// This accounts for the actual JSON representation size, which includes:
	// - Character data (varies by palette: ASCII=1 byte, Unicode=3 bytes per char)
	// - JSON structure overhead (brackets, commas, quotes, color values)
	jsonBytes, _ := json.Marshal(coloredASCII)
	asciiSize := len(jsonBytes)

	// Return JSON response with structured data and size info
	return c.JSON(fiber.Map{
		"lines":          coloredASCII.Lines,
		"originalSize":   fileSize,
		"originalWidth":  originalWidth,
		"originalHeight": originalHeight,
		"asciiSize":      asciiSize,
	})
}

func exportSVGHandler(c *fiber.Ctx) error {
	// Get the uploaded file
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing or invalid image file. Please upload an image using the 'image' field.",
		})
	}

	// Get optional width parameter (default: 100)
	width := 100
	if widthStr := c.FormValue("width"); widthStr != "" {
		if parsedWidth, err := strconv.Atoi(widthStr); err == nil && parsedWidth > 0 {
			width = parsedWidth
		}
	}

	// Get optional palette parameter (default: normal)
	palette := c.FormValue("palette")
	if palette == "" {
		palette = converter.PaletteNormal
	}

	// Get optional color mode
	useColor := c.FormValue("color") == "true" || c.Query("color") == "true"

	// Get optional fontSize (default: 12)
	fontSize := 12
	if fontSizeStr := c.FormValue("fontSize"); fontSizeStr != "" {
		if parsedSize, err := strconv.Atoi(fontSizeStr); err == nil && parsedSize > 0 {
			fontSize = parsedSize
		}
	}

	// Open the uploaded file
	fileHeader, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to open uploaded file",
		})
	}
	defer fileHeader.Close()

	// Load the image from the reader
	img, err := converter.LoadImageFromReader(fileHeader)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Resize the image
	resizedImg := converter.ResizeImage(img, width)

	var svg string
	if useColor {
		coloredASCII := converter.ConvertToASCIIWithColorStructured(resizedImg, palette)
		svg = converter.ConvertToSVG("", &coloredASCII, fontSize)
	} else {
		grayScaleImg := converter.ConvertToGrayscale(resizedImg)
		asciiImg := converter.ConvertToASCII(grayScaleImg, palette)
		svg = converter.ConvertToSVG(asciiImg, nil, fontSize)
	}

	// Generate filename from original file
	filename := generateExportFilename(file.Filename, "_svg")

	c.Set("Content-Type", "image/svg+xml")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	return c.SendString(svg)
}

func convertVideoHandler(c *fiber.Ctx) error {
	// Get the uploaded video file
	file, err := c.FormFile("video")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing or invalid video file. Please upload a video using the 'video' field.",
		})
	}

	// Get file size
	fileSize := file.Size

	// Validate file size (max 50MB)
	const maxFileSize = 50 * 1024 * 1024 // 50MB
	if fileSize > maxFileSize {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("Video file too large. Maximum size is %d MB.", maxFileSize/(1024*1024)),
		})
	}

	// Get optional width parameter (default: 100)
	width := 100
	if widthStr := c.FormValue("width"); widthStr != "" {
		if parsedWidth, err := strconv.Atoi(widthStr); err == nil && parsedWidth > 0 {
			width = parsedWidth
		}
	}

	// Get optional palette parameter (default: normal)
	palette := c.FormValue("palette")
	if palette == "" {
		palette = converter.PaletteNormal
	}

	// Get optional fps parameter (default: 10)
	fps := 10
	if fpsStr := c.FormValue("fps"); fpsStr != "" {
		if parsedFps, err := strconv.Atoi(fpsStr); err == nil && parsedFps > 0 && parsedFps <= 15 {
			fps = parsedFps
		}
	}

	// Get optional color mode (default: false)
	useColor := c.FormValue("color") == "true"

	// Open the uploaded file
	fileHeader, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to open uploaded file",
		})
	}
	defer fileHeader.Close()

	// Extract frames from video
	frames, metadata, err := converter.ExtractFramesFromVideo(fileHeader, fps, int(fileSize))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to extract frames: %v", err),
		})
	}

	// Get first frame dimensions for metadata
	if len(frames) > 0 {
		firstFrame := frames[0]
		bounds := firstFrame.Bounds()
		metadata.Width = bounds.Max.X - bounds.Min.X
		metadata.Height = bounds.Max.Y - bounds.Min.Y
	}

	// Convert frames to ASCII
	if useColor {
		// Color mode
		colorFrames, err := converter.ProcessVideoToColorASCII(frames, width, palette)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": fmt.Sprintf("Failed to convert frames: %v", err),
			})
		}

		// Update timestamps based on actual FPS
		for i := range colorFrames {
			colorFrames[i].Timestamp = float64(i) / float64(fps)
		}

		return c.JSON(converter.VideoColorAsciiResult{
			Frames:   colorFrames,
			Metadata: *metadata,
		})
	} else {
		// Grayscale mode
		asciiFrames, err := converter.ProcessVideoToASCII(frames, width, palette)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": fmt.Sprintf("Failed to convert frames: %v", err),
			})
		}

		// Update timestamps based on actual FPS
		for i := range asciiFrames {
			asciiFrames[i].Timestamp = float64(i) / float64(fps)
		}

		return c.JSON(converter.VideoAsciiResult{
			Frames:   asciiFrames,
			Metadata: *metadata,
		})
	}
}

// generateExportFilename creates a filename by appending suffix before the extension
func generateExportFilename(originalFilename, suffix string) string {
	// Remove path if present, get just the filename
	filename := filepath.Base(originalFilename)

	// Get extension
	ext := filepath.Ext(filename)

	// Remove extension
	nameWithoutExt := strings.TrimSuffix(filename, ext)

	// Append suffix and extension
	return nameWithoutExt + suffix + ext
}

func runCLI(useColor bool, width int, palette string) {
	// Check if user provided an image path (after flags)
	if flag.NArg() < 1 {
		fmt.Println("Usage: go run main.go [flags] <image-path>")
		fmt.Println("       go run main.go --server  (to start API server)")
		fmt.Println("\nFlags:")
		flag.PrintDefaults()
		fmt.Println("\nExample: go run main.go -color -width 120 -palette dense images/apple.png")
		fmt.Println("         go run main.go --server")
		os.Exit(1)
	}

	// Validate palette
	validPalettes := map[string]bool{
		converter.PaletteNormal:  true,
		converter.PaletteDense:   true,
		converter.PaletteSparse:  true,
		converter.PaletteUnicode: true,
	}
	if !validPalettes[palette] {
		fmt.Printf("Error: Invalid palette '%s'. Valid options: normal, dense, sparse, unicode\n", palette)
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
	resizedImg := converter.ResizeImage(img, width)

	// Convert to ASCII (color or grayscale)
	var asciiImg string
	if useColor {
		asciiImg = converter.ConvertToASCIIWithColor(resizedImg, palette)
	} else {
		grayScaleImg := converter.ConvertToGrayscale(resizedImg)
		asciiImg = converter.ConvertToASCII(grayScaleImg, palette)
	}

	// Output the ASCII art
	fmt.Println(asciiImg)
}
