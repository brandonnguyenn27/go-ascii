package converter

// VideoMetadata contains information about the original video
type VideoMetadata struct {
	OriginalSize int     `json:"originalSize"`
	Duration     float64 `json:"duration"`
	OriginalFps  float64 `json:"originalFps"`
	SampledFps   int     `json:"sampledFps"`
	FrameCount   int     `json:"frameCount"`
	Width        int     `json:"width"`
	Height       int     `json:"height"`
}

// FrameASCII represents a single video frame converted to ASCII (grayscale)
type FrameASCII struct {
	Index     int     `json:"index"`
	Timestamp float64 `json:"timestamp"`
	ASCII     string  `json:"ascii"`
}

// FrameColorASCII represents a single video frame converted to colored ASCII
type FrameColorASCII struct {
	Index     int             `json:"index"`
	Timestamp float64         `json:"timestamp"`
	Lines     [][]ColoredChar `json:"lines"`
}

// VideoAsciiResult is the response for grayscale video conversion
type VideoAsciiResult struct {
	Frames   []FrameASCII  `json:"frames"`
	Metadata VideoMetadata `json:"metadata"`
}

// VideoColorAsciiResult is the response for color video conversion
type VideoColorAsciiResult struct {
	Frames   []FrameColorASCII `json:"frames"`
	Metadata VideoMetadata     `json:"metadata"`
}
