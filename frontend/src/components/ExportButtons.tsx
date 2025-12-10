import { Button } from './ui/button';
import { exportToSVG } from '@/lib/api';

interface ExportButtonsProps {
  file: File | null;
  width?: number;
  palette: string;
  colorMode: boolean;
  disabled?: boolean;
}

export function ExportButtons({
  file,
  width,
  palette,
  colorMode,
  disabled,
}: ExportButtonsProps) {
  const handleExportSVG = async () => {
    if (!file) return;

    try {
      await exportToSVG(file, width, palette, colorMode);
    } catch (err) {
      console.error('Failed to export SVG:', err);
      alert('Failed to export SVG. Please try again.');
    }
  };

  return (
    <Button
      onClick={handleExportSVG}
      disabled={disabled || !file}
      variant="outline"
      size="sm"
    >
      Export SVG
    </Button>
  );
}

