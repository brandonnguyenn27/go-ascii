import { useRef, useState } from 'react';
import { Button } from './ui/button';
import { Input } from './ui/input';
import { Label } from './ui/label';
import { Alert, AlertDescription } from './ui/alert';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from './ui/select';

interface VideoUploadProps {
  onFileSelect: (file: File) => void;
  selectedFile: File | null;
  fps: number;
  onFpsChange: (fps: number) => void;
}

const MAX_FILE_SIZE = 50 * 1024 * 1024; // 50MB
const ALLOWED_TYPES = ['video/webm'];
const ALLOWED_EXTENSIONS = ['.webm'];

export function VideoUpload({
  onFileSelect,
  selectedFile,
  fps,
  onFpsChange,
}: VideoUploadProps) {
  const [error, setError] = useState<string | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const validateFile = (file: File): string | null => {
    // Check file type by MIME type
    if (!ALLOWED_TYPES.includes(file.type)) {
      return `Invalid file type. Please upload a WebM video.`;
    }

    // Check file extension (case-insensitive)
    const extension = file.name
      .substring(file.name.lastIndexOf('.'))
      .toLowerCase();
    if (!ALLOWED_EXTENSIONS.includes(extension)) {
      return `Invalid file extension. Please upload a .webm file.`;
    }

    // Check file size
    if (file.size > MAX_FILE_SIZE) {
      const sizeMB = (file.size / (1024 * 1024)).toFixed(1);
      return `File size (${sizeMB}MB) exceeds 50MB limit. Please choose a smaller video.`;
    }

    return null;
  };

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file) {
      return;
    }

    setError(null);
    const validationError = validateFile(file);
    if (validationError) {
      setError(validationError);
      // Reset the input
      if (fileInputRef.current) {
        fileInputRef.current.value = '';
      }
      return;
    }

    onFileSelect(file);
  };

  const handleButtonClick = () => {
    fileInputRef.current?.click();
  };

  const formatFileSize = (bytes: number): string => {
    if (bytes < 1024 * 1024) {
      return `${(bytes / 1024).toFixed(1)} KB`;
    }
    return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
  };

  return (
    <div className="space-y-4">
      <div className="space-y-2">
        <Label htmlFor="video-upload">Video File</Label>
        <div className="flex gap-2">
          <Input
            id="video-upload"
            ref={fileInputRef}
            type="file"
            accept="video/webm"
            onChange={handleFileChange}
            className="hidden"
          />
          <Button type="button" onClick={handleButtonClick} variant="outline">
            Choose Video
          </Button>
          {selectedFile && (
            <div className="flex items-center px-3 text-sm text-muted-foreground">
              {selectedFile.name} ({formatFileSize(selectedFile.size)})
            </div>
          )}
        </div>
        {error && (
          <Alert variant="destructive">
            <AlertDescription>{error}</AlertDescription>
          </Alert>
        )}
        {!error && !selectedFile && (
          <p className="text-sm text-muted-foreground">
            Select a WebM video (max 50MB, ~20 seconds)
          </p>
        )}
      </div>

      <div className="space-y-2">
        <Label htmlFor="fps-select">Frame Rate</Label>
        <Select
          value={fps.toString()}
          onValueChange={(value) => onFpsChange(parseInt(value))}
        >
          <SelectTrigger id="fps-select" className="w-[180px]">
            <SelectValue placeholder="Select FPS" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="10">10 FPS (Recommended)</SelectItem>
            <SelectItem value="12">12 FPS</SelectItem>
            <SelectItem value="15">15 FPS</SelectItem>
          </SelectContent>
        </Select>
        <p className="text-sm text-muted-foreground">
          Higher FPS = smoother playback but slower processing
        </p>
      </div>
    </div>
  );
}
