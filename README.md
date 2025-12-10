# ASCII Converter

A Go-based ASCII art converter that can convert images to ASCII art. Supports both command-line interface and REST API server modes.

## Features

- Convert images (PNG, JPEG) to ASCII art
- Support for colored and grayscale ASCII output
- Adjustable output width
- REST API endpoint for frontend integration
- Command-line interface for local usage

## Installation

1. Clone the repository
2. Navigate to the backend directory:
   ```bash
   cd backend
   ```
3. Install dependencies:
   ```bash
   go mod download
   ```

## Usage

### Terminal Mode (CLI)

Run the converter directly from the command line to convert images locally.

**Basic usage:**

```bash
cd backend
go run main.go <image-path>
```

**With flags:**

```bash
go run main.go [flags] <image-path>
```

#### Flags

- `-color` (boolean): Enable colored ASCII output. Default: `false`
- `-width` (int): Width of ASCII output in characters. Default: `100`
- `-server` (boolean): Start the REST API server instead of CLI mode. Default: `false`

#### Examples

```bash
# Basic conversion (grayscale, width 100)
go run main.go ../images/apple.png

# Colored output with custom width
go run main.go -color -width 120 ../images/apple.png

# Grayscale with custom width
go run main.go -width 80 ../images/ryan.png
```

### Server Mode (REST API)

Start the REST API server to accept image uploads via HTTP requests.

**Start the server:**

```bash
cd backend
go run main.go --server
```

The server will start on `http://localhost:3000`

#### API Endpoints

##### POST `/convert`

Converts an uploaded image to grayscale ASCII art (returns plain text string).

**Request:**

- Method: `POST`
- Content-Type: `multipart/form-data`
- Body: Form data with `image` field containing the image file
- Optional: `width` parameter (form field or query param) - default: `100`

**Response:**

- Content-Type: `application/json`
- Success (200): `{"ascii": "..."}`
- Error (400/500): `{"error": "error message"}`

**Example using curl:**

```bash
curl -X POST http://localhost:3000/convert \
  -F "image=@../images/apple.png"

# With custom width
curl -X POST http://localhost:3000/convert \
  -F "image=@../images/apple.png" \
  -F "width=120"
```

**Example response:**

```json
{
  "ascii": " .:-=+*#%@\n..."
}
```

##### POST `/convert/color`

Converts an uploaded image to colored ASCII art (returns structured data with RGB color information).

**Request:**

- Method: `POST`
- Content-Type: `multipart/form-data`
- Body: Form data with `image` field containing the image file
- Optional: `width` parameter (form field or query param) - default: `100`

**Response:**

- Content-Type: `application/json`
- Success (200): Structured data with lines and character color information
- Error (400/500): `{"error": "error message"}`

**Example using curl:**

```bash
curl -X POST http://localhost:3000/convert/color \
  -F "image=@../images/apple.png"

# With custom width
curl -X POST http://localhost:3000/convert/color \
  -F "image=@../images/apple.png" \
  -F "width=120"
```

**Example response:**

```json
{
  "lines": [
    [
      { "char": " ", "r": 255, "g": 255, "b": 255 },
      { "char": ".", "r": 200, "g": 200, "b": 200 },
      { "char": ":", "r": 150, "g": 150, "b": 150 }
    ],
    [
      { "char": "#", "r": 100, "g": 50, "b": 25 },
      { "char": "@", "r": 50, "g": 25, "b": 10 }
    ]
  ]
}
```

**Frontend Integration Example (React):**

The `/convert/color` endpoint returns structured data that's easy to render without any ANSI parsing:

```jsx
function ColoredAsciiDisplay({ coloredData }) {
  return (
    <pre style={{ fontFamily: "monospace", fontSize: "8px", lineHeight: "1" }}>
      {coloredData.lines.map((line, lineIdx) => (
        <div key={lineIdx}>
          {line.map((char, charIdx) => (
            <span
              key={charIdx}
              style={{
                color: `rgb(${char.r}, ${char.g}, ${char.b})`,
              }}
            >
              {char.char}
            </span>
          ))}
        </div>
      ))}
    </pre>
  );
}

// Usage
const formData = new FormData();
formData.append("image", fileInput.files[0]);

const response = await fetch("http://localhost:3000/convert/color", {
  method: "POST",
  body: formData,
});
const data = await response.json();
<ColoredAsciiDisplay coloredData={data} />;
```

## Project Structure

```
ascii-converter/
├── backend/
│   ├── main.go              # Main entry point (CLI + Server)
│   ├── go.mod
│   ├── go.sum
│   └── pkg/
│       └── converter/
│           ├── colorizer.go  # Colored ASCII conversion
│           ├── grayscale.go  # Grayscale conversion
│           ├── loader.go     # Image loading utilities
│           ├── mapper.go     # Brightness to character mapping
│           └── resizer.go   # Image resizing
├── images/                   # Sample images
└── README.md
```

## Supported Image Formats

- PNG
- JPEG

## Dependencies

- [Fiber](https://github.com/gofiber/fiber) - Web framework for REST API
- [nfnt/resize](https://github.com/nfnt/resize) - Image resizing library
