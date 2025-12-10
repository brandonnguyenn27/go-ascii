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
}

export interface ErrorResponse {
  error: string;
}

/**
 * Converts an image to grayscale ASCII art
 * @param file The image file to convert
 * @param width Optional width in characters (20-300). If not provided or 0, uses original image size.
 * @returns Promise resolving to the ASCII string
 */
export async function convertToAscii(
  file: File,
  width?: number
): Promise<string> {
  const formData = new FormData();
  formData.append('image', file);
  if (width && width > 0) {
    formData.append('width', width.toString());
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
  return data.ascii;
}

/**
 * Converts an image to colored ASCII art
 * @param file The image file to convert
 * @param width Optional width in characters (20-300). If not provided or 0, uses original image size.
 * @returns Promise resolving to the structured color ASCII data
 */
export async function convertToColorAscii(
  file: File,
  width?: number
): Promise<ColorAsciiData> {
  const formData = new FormData();
  formData.append('image', file);
  if (width && width > 0) {
    formData.append('width', width.toString());
  }

  const response = await fetch(`${API_BASE_URL}/convert/color`, {
    method: 'POST',
    body: formData,
  });

  if (!response.ok) {
    const error: ErrorResponse = await response.json();
    throw new Error(error.error || `HTTP error! status: ${response.status}`);
  }

  const data: ColorAsciiData = await response.json();
  return data;
}

