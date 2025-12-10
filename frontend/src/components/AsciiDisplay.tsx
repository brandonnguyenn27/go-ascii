import { Card } from './ui/card';
import type { ColorAsciiData } from '@/lib/api';

interface AsciiDisplayProps {
  asciiData: string | ColorAsciiData | null;
  isLoading: boolean;
  error: string | null;
  isColorMode: boolean;
  splitScreen?: boolean;
  originalImageUrl?: string | null;
}

export function AsciiDisplay({
  asciiData,
  isLoading,
  error,
  isColorMode,
  splitScreen = false,
  originalImageUrl,
}: AsciiDisplayProps) {
  if (isLoading) {
    return (
      <Card className="p-6">
        <div className="text-center text-muted-foreground">
          Converting image to ASCII...
        </div>
      </Card>
    );
  }

  if (error) {
    return (
      <Card className="p-6 border-destructive">
        <div className="text-destructive font-medium">Error: {error}</div>
      </Card>
    );
  }

  if (!asciiData) {
    return (
      <Card className="p-6">
        <div className="text-center text-muted-foreground">
          Upload an image and click Convert to see ASCII art
        </div>
      </Card>
    );
  }

  const renderAscii = () => {
    if (isColorMode && typeof asciiData !== 'string') {
      // Render colored ASCII
      return (
        <pre className="text-xs leading-tight font-mono whitespace-pre text-foreground">
          {asciiData.lines.map((line, lineIndex) => (
            <span key={lineIndex}>
              {line.map((char, charIndex) => {
                const rgb = `rgb(${char.r}, ${char.g}, ${char.b})`;
                return (
                  <span
                    key={charIndex}
                    style={{ color: rgb }}
                  >
                    {char.char}
                  </span>
                );
              })}
              {'\n'}
            </span>
          ))}
        </pre>
      );
    }

    // Render grayscale ASCII
    const asciiText = typeof asciiData === 'string' ? asciiData : '';
    return (
      <pre className="text-xs leading-tight font-mono whitespace-pre text-foreground">
        {asciiText}
      </pre>
    );
  };

  if (splitScreen && originalImageUrl) {
    return (
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <Card className="p-4">
          <h4 className="text-sm font-semibold mb-2">Original Image</h4>
          <div className="flex items-center justify-center bg-muted rounded-md overflow-hidden">
            <img
              src={originalImageUrl}
              alt="Original"
              className="max-w-full max-h-96 object-contain"
            />
          </div>
        </Card>
        <Card className="p-4 overflow-auto bg-background">
          <h4 className="text-sm font-semibold mb-2">ASCII Art</h4>
          {renderAscii()}
        </Card>
      </div>
    );
  }

  return (
    <Card className="p-4 overflow-auto bg-background">
      {renderAscii()}
    </Card>
  );
}

