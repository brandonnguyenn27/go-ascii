interface SizeDisplayProps {
  originalSize?: number;
  asciiSize?: number;
}

function formatBytes(bytes: number): string {
  if (bytes === 0) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + ' ' + sizes[i];
}

export function SizeDisplay({ originalSize, asciiSize }: SizeDisplayProps) {
  if (!originalSize && !asciiSize) {
    return null;
  }

  const reduction = originalSize && asciiSize
    ? ((1 - asciiSize / originalSize) * 100).toFixed(1)
    : null;

  return (
    <div className="text-sm text-muted-foreground space-y-1">
      {originalSize && (
        <div>Original: {formatBytes(originalSize)}</div>
      )}
      {asciiSize && (
        <div>ASCII: {formatBytes(asciiSize)}</div>
      )}
      {reduction && (
        <div className="text-green-600 dark:text-green-400">
          Size reduction: {reduction}%
        </div>
      )}
    </div>
  );
}

