import { Play, Pause } from 'lucide-react';
import { Button } from './ui/button';
import { Slider } from './ui/slider';

interface VideoControlsProps {
  isPlaying: boolean;
  currentFrame: number;
  totalFrames: number;
  onPlayPause: () => void;
  onSeek: (frame: number) => void;
}

export function VideoControls({
  isPlaying,
  currentFrame,
  totalFrames,
  onPlayPause,
  onSeek,
}: VideoControlsProps) {
  const handleSliderChange = (value: number[]) => {
    onSeek(value[0]);
  };

  return (
    <div className="space-y-3 p-4 bg-secondary/20 rounded-lg">
      <div className="flex items-center gap-4">
        <Button
          onClick={onPlayPause}
          variant="outline"
          size="icon"
          className="shrink-0"
        >
          {isPlaying ? (
            <Pause className="h-4 w-4" />
          ) : (
            <Play className="h-4 w-4" />
          )}
        </Button>

        <div className="flex-1">
          <Slider
            value={[currentFrame]}
            onValueChange={handleSliderChange}
            max={totalFrames - 1}
            step={1}
            className="w-full"
          />
        </div>

        <div className="text-sm text-muted-foreground shrink-0 min-w-[100px] text-right">
          Frame {currentFrame + 1} / {totalFrames}
        </div>
      </div>
    </div>
  );
}
