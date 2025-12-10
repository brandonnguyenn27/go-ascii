import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"
import type { ColorAsciiData } from "./api"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

/**
 * Copies ASCII art to clipboard
 * @param asciiData - Either a string (grayscale) or ColorAsciiData (color mode)
 * @returns Promise that resolves when copy is complete
 */
export async function copyAsciiToClipboard(
  asciiData: string | ColorAsciiData
): Promise<void> {
  let textToCopy: string;

  if (typeof asciiData === 'string') {
    // Grayscale mode: copy the string directly
    textToCopy = asciiData;
  } else {
    // Color mode: convert structured data to plain text
    // Extract just the characters, ignoring color information
    textToCopy = asciiData.lines
      .map((line) => line.map((char) => char.char).join(''))
      .join('\n');
  }

  try {
    await navigator.clipboard.writeText(textToCopy);
  } catch (err) {
    throw new Error('Failed to copy to clipboard. Please try again.');
  }
}
