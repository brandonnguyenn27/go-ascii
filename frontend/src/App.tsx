import { useState } from 'react';
import { Button } from './components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from './components/ui/card';
import { FileUpload } from './components/FileUpload';
import { WidthControl } from './components/WidthControl';
import { ColorToggle } from './components/ColorToggle';
import { AsciiDisplay } from './components/AsciiDisplay';
import { Alert, AlertDescription } from './components/ui/alert';
import { convertToAscii, convertToColorAscii, type ColorAsciiData } from './lib/api';
import { copyAsciiToClipboard } from './lib/utils';
import './App.css';

function App() {
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [width, setWidth] = useState(100);
  const [widthEnabled, setWidthEnabled] = useState(false);
  const [colorMode, setColorMode] = useState(false);
  const [asciiResult, setAsciiResult] = useState<string | ColorAsciiData | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [copySuccess, setCopySuccess] = useState(false);

  const handleConvert = async () => {
    if (!selectedFile) {
      setError('Please select an image file first');
      return;
    }

    setIsLoading(true);
    setError(null);
    setCopySuccess(false);

    try {
      const widthToUse = widthEnabled && width > 0 ? width : undefined;
      if (colorMode) {
        const result = await convertToColorAscii(selectedFile, widthToUse);
        setAsciiResult(result);
      } else {
        const result = await convertToAscii(selectedFile, widthToUse);
        setAsciiResult(result);
      }
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to convert image';
      setError(errorMessage);
      setAsciiResult(null);
    } finally {
      setIsLoading(false);
    }
  };

  const handleCopy = async () => {
    if (!asciiResult) {
      return;
    }

    try {
      await copyAsciiToClipboard(asciiResult);
      setCopySuccess(true);
      setTimeout(() => setCopySuccess(false), 2000);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to copy to clipboard';
      setError(errorMessage);
    }
  };

  return (
    <div className="min-h-screen bg-background p-4 md:p-8">
      <div className="max-w-4xl mx-auto space-y-6">
        <Card>
          <CardHeader>
            <CardTitle>ASCII Art Converter</CardTitle>
            <CardDescription>
              Convert your images to ASCII art. Supports both grayscale and color modes.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-6">
            <FileUpload
              onFileSelect={setSelectedFile}
              selectedFile={selectedFile}
            />

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <WidthControl
                value={width}
                onChange={setWidth}
                enabled={widthEnabled}
                onToggle={setWidthEnabled}
              />
              <div className="flex items-end">
                <ColorToggle enabled={colorMode} onToggle={setColorMode} />
              </div>
            </div>

            <Button
              onClick={handleConvert}
              disabled={!selectedFile || isLoading}
              className="w-full"
            >
              {isLoading ? 'Converting...' : 'Convert to ASCII'}
            </Button>

            {copySuccess && (
              <Alert>
                <AlertDescription>Copied to clipboard!</AlertDescription>
              </Alert>
            )}

            <div className="space-y-2">
              <div className="flex items-center justify-between">
                <h3 className="text-lg font-semibold">Output</h3>
                {asciiResult && !colorMode && (
                  <Button
                    onClick={handleCopy}
                    variant="outline"
                    size="sm"
                  >
                    Copy to Clipboard
                  </Button>
                )}
                {asciiResult && colorMode && (
                  <Button
                    disabled
                    variant="outline"
                    size="sm"
                    title="Copy is disabled for color mode (color formatting cannot be preserved in plain text)"
                  >
                    Copy to Clipboard
                  </Button>
                )}
              </div>
              <AsciiDisplay
                asciiData={asciiResult}
                isLoading={isLoading}
                error={error}
                isColorMode={colorMode}
              />
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}

export default App;
