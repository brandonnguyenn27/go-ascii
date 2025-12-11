import { useState, useEffect, useRef } from 'react';
import { VideoPlayer } from './VideoPlayer';
import { VideoControls } from './VideoControls';
import type { VideoAsciiResponse } from '../lib/api';

interface VideoDisplayProps {
  videoResult: VideoAsciiResponse;
  isColorMode: boolean;
}

export function VideoDisplay({ videoResult, isColorMode }: VideoDisplayProps) {
  const [currentFrame, setCurrentFrame] = useState(0);
  const [isPlaying, setIsPlaying] = useState(false);
  const lastFrameTimeRef = useRef(0);
  const animationIdRef = useRef<number | null>(null);

  const { frames, metadata } = videoResult;
  const frameInterval = 1000 / metadata.sampledFps; // milliseconds per frame

  // Playback loop using requestAnimationFrame
  useEffect(() => {
    if (!isPlaying || frames.length === 0) {
      if (animationIdRef.current !== null) {
        cancelAnimationFrame(animationIdRef.current);
        animationIdRef.current = null;
      }
      return;
    }

    const animate = (timestamp: number) => {
      if (timestamp - lastFrameTimeRef.current >= frameInterval) {
        setCurrentFrame((prev) => {
          const next = prev + 1;
          if (next >= frames.length) {
            // End of video - stop playing
            setIsPlaying(false);
            return 0; // Reset to start
          }
          return next;
        });
        lastFrameTimeRef.current = timestamp;
      }
      animationIdRef.current = requestAnimationFrame(animate);
    };

    animationIdRef.current = requestAnimationFrame(animate);

    return () => {
      if (animationIdRef.current !== null) {
        cancelAnimationFrame(animationIdRef.current);
      }
    };
  }, [isPlaying, frames.length, frameInterval]);

  const handlePlayPause = () => {
    if (!isPlaying && currentFrame >= frames.length - 1) {
      // If at the end, restart from beginning
      setCurrentFrame(0);
    }
    setIsPlaying(!isPlaying);
    lastFrameTimeRef.current = performance.now();
  };

  const handleSeek = (frame: number) => {
    setCurrentFrame(frame);
    setIsPlaying(false);
  };

  if (frames.length === 0) {
    return (
      <div className="flex items-center justify-center h-64 text-muted-foreground">
        No video frames available
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <div className="border rounded-lg overflow-hidden bg-black">
        <VideoPlayer
          frames={frames}
          currentFrame={currentFrame}
          isColorMode={isColorMode}
        />
      </div>

      <VideoControls
        isPlaying={isPlaying}
        currentFrame={currentFrame}
        totalFrames={frames.length}
        onPlayPause={handlePlayPause}
        onSeek={handleSeek}
      />

      <div className="text-sm text-muted-foreground space-y-1">
        <div className="flex justify-between">
          <span>Duration:</span>
          <span>{metadata.duration?.toFixed(1) || 'N/A'}s</span>
        </div>
        <div className="flex justify-between">
          <span>Frames:</span>
          <span>{metadata.frameCount} @ {metadata.sampledFps} FPS</span>
        </div>
        <div className="flex justify-between">
          <span>Dimensions:</span>
          <span>
            {metadata.width}x{metadata.height}
          </span>
        </div>
      </div>
    </div>
  );
}
