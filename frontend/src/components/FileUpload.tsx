import { useRef, useState } from 'react';
import { Button } from './ui/button';
import { Input } from './ui/input';
import { Label } from './ui/label';
import { Alert, AlertDescription } from './ui/alert';

interface FileUploadProps {
  onFileSelect: (file: File) => void;
  selectedFile: File | null;
}

const MAX_FILE_SIZE = 10 * 1024 * 1024; // 10MB
const ALLOWED_TYPES = ['image/png', 'image/jpeg'];
const ALLOWED_EXTENSIONS = ['.png', '.jpg', '.jpeg'];

export function FileUpload({ onFileSelect, selectedFile }: FileUploadProps) {
  const [error, setError] = useState<string | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const validateFile = (file: File): string | null => {
    // Check file type by MIME type
    if (!ALLOWED_TYPES.includes(file.type)) {
      return `Invalid file type. Please upload a PNG or JPEG image.`;
    }

    // Check file extension (case-insensitive)
    const extension = file.name
      .substring(file.name.lastIndexOf('.'))
      .toLowerCase();
    if (!ALLOWED_EXTENSIONS.includes(extension)) {
      return `Invalid file extension. Please upload a .png, .jpg, or .jpeg file.`;
    }

    // Check file size
    if (file.size > MAX_FILE_SIZE) {
      return `File size exceeds 10MB limit. Please choose a smaller image.`;
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

  return (
    <div className="space-y-2">
      <Label htmlFor="image-upload">Image File</Label>
      <div className="flex gap-2">
        <Input
          id="image-upload"
          ref={fileInputRef}
          type="file"
          accept="image/png,image/jpeg,image/jpg"
          onChange={handleFileChange}
          className="hidden"
        />
        <Button type="button" onClick={handleButtonClick} variant="outline">
          Choose File
        </Button>
        {selectedFile && (
          <div className="flex items-center px-3 text-sm text-muted-foreground">
            {selectedFile.name}
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
          Select a PNG or JPEG image (max 10MB)
        </p>
      )}
    </div>
  );
}

