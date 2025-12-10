import { Card } from './ui/card';
import type { ColorAsciiData } from '@/lib/api';

interface AsciiDisplayProps {
  asciiData: string | ColorAsciiData | null;
  isLoading: boolean;
  error: string | null;
  isColorMode: boolean;
}

export function AsciiDisplay({
  asciiData,
  isLoading,
  error,
  isColorMode,
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

  if (isColorMode && typeof asciiData !== 'string') {
    // Render colored ASCII
    return (
      <Card className="p-4 overflow-auto bg-background">
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
      </Card>
    );
  }

  // Render grayscale ASCII
  const asciiText = typeof asciiData === 'string' ? asciiData : '';
  return (
    <Card className="p-4 overflow-auto bg-background">
      <pre className="text-xs leading-tight font-mono whitespace-pre text-foreground">
        {asciiText}
      </pre>
    </Card>
  );
}

