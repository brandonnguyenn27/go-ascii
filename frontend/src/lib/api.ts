const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:3000';

export interface ColoredChar {
  char: string;
  r: number;
  g: number;
  b: number;
}

export interface ColorAsciiData {
  lines: ColoredChar[][];
}

export interface GrayscaleAsciiResponse {
  ascii: string;
  originalSize?: number;
  originalWidth?: number;
  originalHeight?: number;
  asciiSize?: number;
}

export interface ColorAsciiResponse {
  lines: ColoredChar[][];
  originalSize?: number;
  originalWidth?: number;
  originalHeight?: number;
  asciiSize?: number;
}

export interface ErrorResponse {
  error: string;
}

/**
 * Converts an image to grayscale ASCII art
 * @param file The image file to convert
 * @param width Optional width in characters (20-300). If not provided or 0, uses original image size.
 * @param palette Optional palette type (normal, dense, sparse, unicode). Defaults to normal.
 * @returns Promise resolving to the ASCII string
 */
export async function convertToAscii(
  file: File,
  width?: number,
  palette?: string
): Promise<string> {
  const formData = new FormData();
  formData.append('image', file);
  if (width && width > 0) {
    formData.append('width', width.toString());
  }
  if (palette) {
    formData.append('palette', palette);
  }

  const response = await fetch(`${API_BASE_URL}/convert`, {
    method: 'POST',
    body: formData,
  });

  if (!response.ok) {
    const error: ErrorResponse = await response.json();
    throw new Error(error.error || `HTTP error! status: ${response.status}`);
  }

  const data: GrayscaleAsciiResponse = await response.json();
  // Store size info if available (could be used later)
  return data.ascii;
}

/**
 * Converts an image to colored ASCII art
 * @param file The image file to convert
 * @param width Optional width in characters (20-300). If not provided or 0, uses original image size.
 * @param palette Optional palette type (normal, dense, sparse, unicode). Defaults to normal.
 * @returns Promise resolving to the structured color ASCII data
 */
export async function convertToColorAscii(
  file: File,
  width?: number,
  palette?: string
): Promise<ColorAsciiData> {
  const formData = new FormData();
  formData.append('image', file);
  if (width && width > 0) {
    formData.append('width', width.toString());
  }
  if (palette) {
    formData.append('palette', palette);
  }

  const response = await fetch(`${API_BASE_URL}/convert/color`, {
    method: 'POST',
    body: formData,
  });

  if (!response.ok) {
    const error: ErrorResponse = await response.json();
    throw new Error(error.error || `HTTP error! status: ${response.status}`);
  }

  const data: ColorAsciiResponse = await response.json();
  return { lines: data.lines };
}

/**
 * Exports ASCII art as SVG
 * @param file The image file to convert and export
 * @param width Optional width in characters
 * @param palette Optional palette type
 * @param colorMode Whether to use color mode
 */
export async function exportToSVG(
  file: File,
  width?: number,
  palette?: string,
  colorMode?: boolean
): Promise<void> {
  const formData = new FormData();
  formData.append('image', file);
  if (width && width > 0) {
    formData.append('width', width.toString());
  }
  if (palette) {
    formData.append('palette', palette);
  }
  if (colorMode) {
    formData.append('color', 'true');
  }

  const response = await fetch(`${API_BASE_URL}/export/svg`, {
    method: 'POST',
    body: formData,
  });

  if (!response.ok) {
    const error: ErrorResponse = await response.json();
    throw new Error(error.error || `HTTP error! status: ${response.status}`);
  }

  // Extract filename from Content-Disposition header, or generate from original filename
  let filename = 'ascii-art.svg';
  const contentDisposition = response.headers.get('Content-Disposition');
  if (contentDisposition) {
    const filenameMatch = contentDisposition.match(/filename="?([^"]+)"?/);
    if (filenameMatch && filenameMatch[1]) {
      filename = filenameMatch[1];
    }
  }
  
  // Fallback: generate filename from original file if header not available
  if (filename === 'ascii-art.svg') {
    const originalName = file.name.replace(/\.[^/.]+$/, '');
    filename = `${originalName}_svg.svg`;
  }

  const blob = await response.blob();
  const url = window.URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = filename;
  document.body.appendChild(a);
  a.click();
  window.URL.revokeObjectURL(url);
  document.body.removeChild(a);
}

export interface VideoFrame {
  index: number;
  timestamp: number;
  ascii?: string;
  lines?: ColoredChar[][];
}

export interface VideoMetadata {
  originalSize: number;
  duration: number;
  originalFps: number;
  sampledFps: number;
  frameCount: number;
  width: number;
  height: number;
}

export interface VideoAsciiResponse {
  frames: VideoFrame[];
  metadata: VideoMetadata;
}

/**
 * Converts a video to grayscale ASCII art frames
 * @param file The video file to convert
 * @param width Optional width in characters (20-300)
 * @param palette Optional palette type (normal, dense, sparse, unicode)
 * @param fps Optional target frame rate (10-15 fps)
 * @returns Promise resolving to video frames and metadata
 */
export async function convertVideoToAscii(
  file: File,
  width?: number,
  palette?: string,
  fps?: number
): Promise<VideoAsciiResponse> {
  const formData = new FormData();
  formData.append('video', file);
  if (width && width > 0) {
    formData.append('width', width.toString());
  }
  if (palette) {
    formData.append('palette', palette);
  }
  if (fps && fps > 0) {
    formData.append('fps', fps.toString());
  }

  const response = await fetch(`${API_BASE_URL}/convert/video`, {
    method: 'POST',
    body: formData,
  });

  if (!response.ok) {
    const error: ErrorResponse = await response.json();
    throw new Error(error.error || `HTTP error! status: ${response.status}`);
  }

  const data: VideoAsciiResponse = await response.json();
  return data;
}

/**
 * Converts a video to colored ASCII art frames
 * @param file The video file to convert
 * @param width Optional width in characters (20-300)
 * @param palette Optional palette type (normal, dense, sparse, unicode)
 * @param fps Optional target frame rate (10-15 fps)
 * @returns Promise resolving to video frames and metadata
 */
export async function convertVideoToColorAscii(
  file: File,
  width?: number,
  palette?: string,
  fps?: number
): Promise<VideoAsciiResponse> {
  const formData = new FormData();
  formData.append('video', file);
  if (width && width > 0) {
    formData.append('width', width.toString());
  }
  if (palette) {
    formData.append('palette', palette);
  }
  if (fps && fps > 0) {
    formData.append('fps', fps.toString());
  }
  formData.append('color', 'true');

  const response = await fetch(`${API_BASE_URL}/convert/video`, {
    method: 'POST',
    body: formData,
  });

  if (!response.ok) {
    const error: ErrorResponse = await response.json();
    throw new Error(error.error || `HTTP error! status: ${response.status}`);
  }

  const data: VideoAsciiResponse = await response.json();
  return data;
}

