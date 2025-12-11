package converter

import (
	"encoding/json"
	"fmt"
	"image"
	"io"
	"os"
	"path/filepath"
	"strconv"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

const (
	MaxFrameCount = 200 // Maximum number of frames to process
	MaxDuration   = 20  // Maximum video duration in seconds
)

// ExtractFramesFromVideo extracts frames from a video file at the specified FPS
func ExtractFramesFromVideo(reader io.Reader, targetFps int, originalSize int) ([]image.Image, *VideoMetadata, error) {
	// Create a temporary file for the video
	tmpVideo, err := os.CreateTemp("", "video-*.webm")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create temp video file: %w", err)
	}
	tmpVideoPath := tmpVideo.Name()
	defer os.Remove(tmpVideoPath)

	// Copy the video data to the temp file
	_, err = io.Copy(tmpVideo, reader)
	if err != nil {
		tmpVideo.Close()
		return nil, nil, fmt.Errorf("failed to write video to temp file: %w", err)
	}
	tmpVideo.Close()

	// Probe video to get metadata
	metadata, err := probeVideo(tmpVideoPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to probe video: %w", err)
	}
	metadata.OriginalSize = originalSize
	metadata.SampledFps = targetFps

	// Calculate how many frames to extract
	totalFrames := int(metadata.Duration * float64(targetFps))
	if totalFrames > MaxFrameCount {
		totalFrames = MaxFrameCount
		fmt.Printf("Warning: limiting frames to %d (video is too long)\n", MaxFrameCount)
	}
	metadata.FrameCount = totalFrames

	// Create temp directory for extracted frames
	tmpDir, err := os.MkdirTemp("", "frames-*")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// Extract frames using FFmpeg
	// Use fps filter to extract frames at target FPS
	framePattern := filepath.Join(tmpDir, "frame_%04d.jpg")

	err = ffmpeg.Input(tmpVideoPath).
		Filter("fps", ffmpeg.Args{fmt.Sprintf("%d", targetFps)}).
		Output(framePattern, ffmpeg.KwArgs{
			"frames:v": totalFrames,
			"q:v":      2, // Quality (1-31, lower is better)
		}).
		OverWriteOutput().
		ErrorToStdOut().
		Run()

	if err != nil {
		return nil, nil, fmt.Errorf("failed to extract frames: %w", err)
	}

	// Load all extracted frames
	frames := make([]image.Image, 0, totalFrames)
	for i := 1; i <= totalFrames; i++ {
		framePath := filepath.Join(tmpDir, fmt.Sprintf("frame_%04d.jpg", i))

		// Check if file exists
		if _, err := os.Stat(framePath); os.IsNotExist(err) {
			fmt.Printf("Warning: frame %d not found, stopping extraction\n", i)
			break
		}

		img, err := LoadImage(framePath)
		if err != nil {
			fmt.Printf("Warning: failed to load frame %d: %v\n", i, err)
			continue
		}
		frames = append(frames, img)
	}

	if len(frames) == 0 {
		return nil, nil, fmt.Errorf("no frames were extracted from video")
	}

	metadata.FrameCount = len(frames)
	return frames, metadata, nil
}

// FFProbeFormat represents the format section of ffprobe output
type FFProbeFormat struct {
	Duration string `json:"duration"`
}

// FFProbeStream represents a stream in ffprobe output
type FFProbeStream struct {
	CodecType string `json:"codec_type"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	RFrameRate string `json:"r_frame_rate"`
}

// FFProbeOutput represents the full ffprobe JSON output
type FFProbeOutput struct {
	Streams []FFProbeStream `json:"streams"`
	Format  FFProbeFormat   `json:"format"`
}

// probeVideo uses FFmpeg to extract video metadata
func probeVideo(videoPath string) (*VideoMetadata, error) {
	// Use ffprobe to get video information
	data, err := ffmpeg.Probe(videoPath)
	if err != nil {
		return nil, fmt.Errorf("ffprobe failed: %w", err)
	}

	// Parse the JSON output
	var probeData FFProbeOutput
	err = json.Unmarshal([]byte(data), &probeData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ffprobe output: %w", err)
	}

	// Initialize metadata with defaults
	metadata := &VideoMetadata{
		Duration:    10.0,
		OriginalFps: 30.0,
		Width:       640,
		Height:      480,
	}

	// Parse duration from format
	if probeData.Format.Duration != "" {
		duration, err := strconv.ParseFloat(probeData.Format.Duration, 64)
		if err == nil {
			metadata.Duration = duration
		}
	}

	// Find the first video stream
	for _, stream := range probeData.Streams {
		if stream.CodecType == "video" {
			metadata.Width = stream.Width
			metadata.Height = stream.Height

			// Parse frame rate (format is typically "30/1" or "2997/100")
			if stream.RFrameRate != "" {
				var num, den float64
				_, err := fmt.Sscanf(stream.RFrameRate, "%f/%f", &num, &den)
				if err == nil && den > 0 {
					metadata.OriginalFps = num / den
				}
			}
			break
		}
	}

	return metadata, nil
}

// ProcessVideoToASCII converts all frames to grayscale ASCII
func ProcessVideoToASCII(frames []image.Image, width int, palette string) ([]FrameASCII, error) {
	result := make([]FrameASCII, 0, len(frames))

	for i, frame := range frames {
		// Resize the frame
		resized := ResizeImage(frame, width)

		// Convert to grayscale
		grayscale := ConvertToGrayscale(resized)

		// Convert to ASCII
		ascii := ConvertToASCII(grayscale, palette)

		// Calculate timestamp (assuming frames are evenly spaced)
		timestamp := float64(i) / 10.0 // Default to 10 fps spacing

		result = append(result, FrameASCII{
			Index:     i,
			Timestamp: timestamp,
			ASCII:     ascii,
		})
	}

	return result, nil
}

// ProcessVideoToColorASCII converts all frames to colored ASCII
func ProcessVideoToColorASCII(frames []image.Image, width int, palette string) ([]FrameColorASCII, error) {
	result := make([]FrameColorASCII, 0, len(frames))

	for i, frame := range frames {
		// Resize the frame
		resized := ResizeImage(frame, width)

		// Convert to colored ASCII (structured format)
		coloredASCII := ConvertToASCIIWithColorStructured(resized, palette)

		// Calculate timestamp (assuming frames are evenly spaced)
		timestamp := float64(i) / 10.0 // Default to 10 fps spacing

		result = append(result, FrameColorASCII{
			Index:     i,
			Timestamp: timestamp,
			Lines:     coloredASCII.Lines,
		})
	}

	return result, nil
}
