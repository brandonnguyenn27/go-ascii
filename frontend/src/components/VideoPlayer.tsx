import { useEffect, useRef } from 'react';
import type { VideoFrame, ColoredChar } from '../lib/api';

interface VideoPlayerProps {
  frames: VideoFrame[];
  currentFrame: number;
  isColorMode: boolean;
}

export function VideoPlayer({
  frames,
  currentFrame,
  isColorMode,
}: VideoPlayerProps) {
  const containerRef = useRef<HTMLPreElement>(null);

  // Render the current frame
  useEffect(() => {
    if (!containerRef.current || frames.length === 0) return;

    const frame = frames[currentFrame];
    if (!frame) return;

    if (isColorMode && frame.lines) {
      // Render colored ASCII
      renderColorFrame(frame.lines, containerRef.current);
    } else if (frame.ascii) {
      // Render grayscale ASCII
      containerRef.current.textContent = frame.ascii;
    }
  }, [frames, currentFrame, isColorMode]);

  const renderColorFrame = (
    lines: ColoredChar[][],
    container: HTMLPreElement
  ) => {
    // Clear existing content
    container.innerHTML = '';

    // Create a document fragment for better performance
    const fragment = document.createDocumentFragment();

    lines.forEach((line, lineIndex) => {
      line.forEach((char) => {
        const span = document.createElement('span');
        span.textContent = char.char;
        span.style.color = `rgb(${char.r}, ${char.g}, ${char.b})`;
        fragment.appendChild(span);
      });

      // Add newline after each line except the last
      if (lineIndex < lines.length - 1) {
        fragment.appendChild(document.createTextNode('\n'));
      }
    });

    container.appendChild(fragment);
  };

  if (frames.length === 0) {
    return (
      <div className="flex items-center justify-center h-64 text-muted-foreground">
        No frames to display
      </div>
    );
  }

  return (
    <div className="w-full overflow-auto">
      <pre
        ref={containerRef}
        className="font-mono text-xs leading-tight whitespace-pre"
        style={{
          color: isColorMode ? undefined : 'white',
          margin: 0,
          padding: '1rem',
          backgroundColor: 'black',
        }}
      />
    </div>
  );
}
