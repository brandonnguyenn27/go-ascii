package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/brandonnguyenn27/ascii-converter/pkg/converter"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// Define all flags
	serverMode := flag.Bool("server", false, "Start the REST API server")
	useColor := flag.Bool("color", false, "Enable colored ASCII output")
	width := flag.Int("width", 100, "Width of ASCII output in characters")

	flag.Parse()

	if *serverMode {
		startServer()
	} else {
		runCLI(*useColor, *width)
	}
}

func startServer() {
	app := fiber.New()

	// Configure CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
		AllowMethods: "POST,OPTIONS",
		AllowHeaders: "Content-Type",
	}))

	app.Post("/convert", convertHandler)            // Grayscale ASCII (returns string)
	app.Post("/convert/color", convertColorHandler) // Colored ASCII (returns structured data)

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

	// Convert to grayscale ASCII
	grayScaleImg := converter.ConvertToGrayscale(resizedImg)
	asciiImg := converter.ConvertToASCII(grayScaleImg)

	// Return JSON response
	return c.JSON(fiber.Map{
		"ascii": asciiImg,
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

	// Convert to colored ASCII with structured data
	coloredASCII := converter.ConvertToASCIIWithColorStructured(resizedImg)

	// Return JSON response with structured data
	return c.JSON(coloredASCII)
}

func runCLI(useColor bool, width int) {
	// Check if user provided an image path (after flags)
	if flag.NArg() < 1 {
		fmt.Println("Usage: go run main.go [flags] <image-path>")
		fmt.Println("       go run main.go --server  (to start API server)")
		fmt.Println("\nFlags:")
		flag.PrintDefaults()
		fmt.Println("\nExample: go run main.go -color -width 120 images/apple.png")
		fmt.Println("         go run main.go --server")
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
		asciiImg = converter.ConvertToASCIIWithColor(resizedImg)
	} else {
		grayScaleImg := converter.ConvertToGrayscale(resizedImg)
		asciiImg = converter.ConvertToASCII(grayScaleImg)
	}

	// Output the ASCII art
	fmt.Println(asciiImg)
}
