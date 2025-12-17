package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

// CaptureConfig holds the configuration for capturing frames
type CaptureConfig struct {
	RTSPUrl       string `json:"rtspUrl"`
	Interval      int    `json:"interval"`      // seconds between captures
	CleanupFrames bool   `json:"cleanupFrames"` // delete frames after video generation
	FPS           int    `json:"fps"`           // output video FPS (default 30)
	Quality       string `json:"quality"`       // video quality: "high", "medium", "low"
}

// CaptureSession represents an active capture session
type CaptureSession struct {
	Config     CaptureConfig
	Running    bool
	StartTime  time.Time
	FrameCount int
	StopChan   chan bool
	mu         sync.RWMutex
}

var (
	currentSession *CaptureSession
	sessionMutex   sync.Mutex
	streamProcesses = make(map[*exec.Cmd]bool)
	streamMutex     sync.Mutex
)

// StartCapture begins capturing frames from the RTSP stream
func StartCapture(config CaptureConfig) error {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	// Check if already running
	if currentSession != nil && currentSession.Running {
		return fmt.Errorf("capture already running")
	}

	// Validate configuration
	if config.RTSPUrl == "" {
		return fmt.Errorf("RTSP URL is required")
	}
	if config.Interval < 1 {
		return fmt.Errorf("capture interval must be at least 1 second")
	}

	// Validate FFmpeg is installed
	if err := checkFFmpeg(); err != nil {
		return fmt.Errorf("ffmpeg not found - please install with: brew install ffmpeg")
	}

	// Test RTSP connection before starting capture
	log.Println("Testing RTSP connection...")
	if err := testRTSPConnection(config.RTSPUrl); err != nil {
		return fmt.Errorf("cannot connect to camera: %w", err)
	}

	// Create new session
	session := &CaptureSession{
		Config:     config,
		Running:    true,
		StartTime:  time.Now(),
		FrameCount: 0,
		StopChan:   make(chan bool),
	}

	currentSession = session

	// Start capture in background
	go runCapture(session)

	return nil
}

// StopCapture stops the current capture session and generates timelapse
func StopCapture() error {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	if currentSession == nil || !currentSession.Running {
		return fmt.Errorf("no active capture session")
	}

	// Signal to stop
	close(currentSession.StopChan)
	currentSession.mu.Lock()
	currentSession.Running = false
	currentSession.mu.Unlock()

	// Generate timelapse video
	go generateTimelapse(currentSession)

	return nil
}

// GetStatus returns the current capture status
func GetStatus() map[string]interface{} {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	if currentSession == nil || !currentSession.Running {
		return map[string]interface{}{
			"running":    false,
			"frameCount": 0,
			"duration":   "0s",
		}
	}

	currentSession.mu.RLock()
	defer currentSession.mu.RUnlock()

	duration := time.Since(currentSession.StartTime).Round(time.Second)
	return map[string]interface{}{
		"running":    currentSession.Running,
		"frameCount": currentSession.FrameCount,
		"duration":   duration.String(),
	}
}

// runCapture performs the actual frame capture loop
func runCapture(session *CaptureSession) {
	ticker := time.NewTicker(time.Duration(session.Config.Interval) * time.Second)
	defer ticker.Stop()

	log.Printf("Starting capture from %s with %d second interval",
		session.Config.RTSPUrl, session.Config.Interval)

	// Capture first frame immediately
	captureFrame(session)

	for {
		select {
		case <-session.StopChan:
			log.Println("Capture stopped")
			return
		case <-ticker.C:
			captureFrame(session)
		}
	}
}

// captureFrame captures a single frame from the RTSP stream
func captureFrame(session *CaptureSession) {
	session.mu.Lock()
	frameNum := session.FrameCount
	session.mu.Unlock()

	// Generate filename with zero-padded frame number
	filename := fmt.Sprintf("frame_%05d.jpg", frameNum)
	filepath := filepath.Join("frames", filename)

	// Use FFmpeg to capture a single frame from RTSP stream
	cmd := exec.Command("ffmpeg",
		"-rtsp_transport", "tcp", // Use TCP for more reliable streaming
		"-i", session.Config.RTSPUrl,
		"-vframes", "1", // Capture only 1 frame
		"-q:v", "2", // High quality JPEG (2 is high quality)
		"-y", // Overwrite output file
		filepath,
	)

	// Capture output for debugging
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error capturing frame %d: %v\nOutput: %s", frameNum, err, string(output))
		return
	}

	session.mu.Lock()
	session.FrameCount++
	session.mu.Unlock()

	log.Printf("Captured frame %d -> %s", frameNum, filepath)
}

// generateTimelapse creates a timelapse video from captured frames
func generateTimelapse(session *CaptureSession) {
	log.Println("Generating timelapse video...")

	// Generate output filename with timestamp
	timestamp := session.StartTime.Format("2006-01-02_15-04-05")
	outputFile := filepath.Join("output", fmt.Sprintf("timelapse_%s.mp4", timestamp))

	// Determine FPS (default to 30)
	fps := session.Config.FPS
	if fps <= 0 {
		fps = 30
	}

	// Determine CRF based on quality setting
	// CRF: 18 = high quality, 23 = default, 28 = lower quality
	crf := "23" // default
	switch session.Config.Quality {
	case "high":
		crf = "18"
	case "low":
		crf = "28"
	default:
		crf = "23" // medium
	}

	// Use FFmpeg to create timelapse video
	// -framerate: Output video FPS
	// -pattern_type glob: Use glob pattern to match files
	// -i "frames/frame_*.jpg": Input pattern
	// -c:v libx264: Use H.264 codec
	// -pix_fmt yuv420p: Pixel format for compatibility
	// -crf: Quality (lower = better)
	cmd := exec.Command("ffmpeg",
		"-framerate", fmt.Sprintf("%d", fps),
		"-pattern_type", "glob",
		"-i", "frames/frame_*.jpg",
		"-c:v", "libx264",
		"-pix_fmt", "yuv420p",
		"-crf", crf,
		"-y",
		outputFile,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error generating timelapse: %v\nOutput: %s", err, string(output))
		return
	}

	log.Printf("Timelapse video created: %s (FPS: %d, Quality: %s)", outputFile, fps, session.Config.Quality)
	log.Printf("Total frames: %d, Duration: %v",
		session.FrameCount,
		time.Since(session.StartTime).Round(time.Second))

	// Clean up frames if requested
	if session.Config.CleanupFrames {
		log.Println("Cleaning up frame files...")
		cleanupFrames()
	} else {
		log.Println("Frame files preserved in frames/ directory")
	}
}

// cleanupFrames removes all captured frame images
func cleanupFrames() {
	matches, err := filepath.Glob("frames/frame_*.jpg")
	if err != nil {
		log.Printf("Error finding frames to clean up: %v", err)
		return
	}

	for _, file := range matches {
		if err := os.Remove(file); err != nil {
			log.Printf("Error removing frame %s: %v", file, err)
		}
	}

	log.Printf("Cleaned up %d frame files", len(matches))
}

// checkFFmpeg verifies that ffmpeg is installed and accessible
func checkFFmpeg() error {
	cmd := exec.Command("ffmpeg", "-version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg is not installed or not in PATH")
	}
	return nil
}

// testRTSPConnection tests if the RTSP URL is accessible
func testRTSPConnection(rtspUrl string) error {
	// Try to capture a single test frame with a 10-second timeout
	cmd := exec.Command("ffmpeg",
		"-rtsp_transport", "tcp",
		"-i", rtspUrl,
		"-vframes", "1",
		"-f", "null",
		"-",
	)

	// Set a timeout for the test
	done := make(chan error, 1)
	go func() {
		done <- cmd.Run()
	}()

	select {
	case err := <-done:
		if err != nil {
			return fmt.Errorf("RTSP connection failed - check camera IP and URL format")
		}
		return nil
	case <-time.After(10 * time.Second):
		cmd.Process.Kill()
		return fmt.Errorf("RTSP connection timeout - camera not responding")
	}
}

// ParseCaptureConfig parses capture configuration from JSON
func ParseCaptureConfig(data []byte) (CaptureConfig, error) {
	var config CaptureConfig
	err := json.Unmarshal(data, &config)
	return config, err
}

// RegisterStreamProcess adds a stream process to the tracking map
func RegisterStreamProcess(cmd *exec.Cmd) {
	streamMutex.Lock()
	defer streamMutex.Unlock()
	streamProcesses[cmd] = true
}

// UnregisterStreamProcess removes a stream process from the tracking map
func UnregisterStreamProcess(cmd *exec.Cmd) {
	streamMutex.Lock()
	defer streamMutex.Unlock()
	delete(streamProcesses, cmd)
}

// KillAllStreamProcesses forcefully kills all active stream processes
func KillAllStreamProcesses() {
	streamMutex.Lock()
	defer streamMutex.Unlock()

	for cmd := range streamProcesses {
		if cmd.Process != nil {
			log.Printf("Killing stream process PID %d", cmd.Process.Pid)
			cmd.Process.Kill()
		}
	}
	// Clear the map
	streamProcesses = make(map[*exec.Cmd]bool)
}
