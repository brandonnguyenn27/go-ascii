# ASCII Converter

A Go-based ASCII art converter that can convert images to ASCII art. Supports both command-line interface, REST API server, and a modern web frontend.

## Features

- Convert images (PNG, JPEG) to ASCII art
- Support for colored and grayscale ASCII output
- Adjustable output width (with optional custom width control)
- Modern React web frontend with dark terminal theme
- REST API endpoint for frontend integration
- Command-line interface for local usage

## Installation

### Backend

1. Clone the repository
2. Navigate to the backend directory:
   ```bash
   cd backend
   ```
3. Install dependencies:
   ```bash
   go mod download
   ```

### Frontend

1. Navigate to the frontend directory:
   ```bash
   cd frontend
   ```
2. Install dependencies:
   ```bash
   npm install
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

### Web Frontend

The project includes a modern React frontend with a dark terminal-themed UI.

**Start the frontend:**

1. Make sure the backend server is running (see [Server Mode](#server-mode-rest-api))
2. In a new terminal, navigate to the frontend directory:
   ```bash
   cd frontend
   ```
3. Start the development server:
   ```bash
   npm run dev
   ```
4. Open your browser and navigate to `http://localhost:5173`

**Using the Frontend:**

- **Upload Image**: Click "Choose File" and select a PNG or JPEG image (max 10MB)
- **Custom Width**: Toggle "Custom Width" to enable width control (default: 100 characters)
  - Use the slider or input field to set width between 20-300 characters
- **Color Mode**: Toggle "Color Mode" to enable colored ASCII output
- **Convert**: Click "Convert to ASCII" to process your image
- **Copy**: For grayscale mode, use "Copy to Clipboard" to copy the ASCII art

**Build for Production:**

```bash
cd frontend
npm run build
```

The built files will be in the `dist/` directory.

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
├── frontend/
│   ├── src/
│   │   ├── components/      # React components
│   │   ├── lib/             # API client and utilities
│   │   └── App.tsx          # Main app component
│   ├── package.json
│   └── vite.config.ts
├── images/                   # Sample images
└── README.md
```

## Supported Image Formats

- PNG
- JPEG

## Dependencies

### Backend

- [Fiber](https://github.com/gofiber/fiber) - Web framework for REST API
- [nfnt/resize](https://github.com/nfnt/resize) - Image resizing library

### Frontend

- [React](https://react.dev/) - UI framework
- [Vite](https://vitejs.dev/) - Build tool and dev server
- [Tailwind CSS](https://tailwindcss.com/) - Styling
- [shadcn/ui](https://ui.shadcn.com/) - UI component library
