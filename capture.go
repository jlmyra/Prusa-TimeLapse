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

	RTSPUrl  string `json:"rtspUrl"`

	Interval int    `json:"interval"` // seconds between captures

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

)

 

// StartCapture begins capturing frames from the RTSP stream

func StartCapture(config CaptureConfig) error {

	sessionMutex.Lock()

	defer sessionMutex.Unlock()

 

	// Check if already running

	if currentSession != nil && currentSession.Running {

		return fmt.Errorf("capture already running")

	}

 

	// Validate FFmpeg is installed

	if err := checkFFmpeg(); err != nil {

		return fmt.Errorf("ffmpeg not found: %w", err)

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

 

	// Use FFmpeg to create timelapse video

	// -framerate 30: Output video will play at 30fps

	// -pattern_type glob: Use glob pattern to match files

	// -i "frames/frame_*.jpg": Input pattern

	// -c:v libx264: Use H.264 codec

	// -pix_fmt yuv420p: Pixel format for compatibility

	// -crf 23: Quality (lower = better, 23 is default)

	cmd := exec.Command("ffmpeg",

		"-framerate", "30",

		"-pattern_type", "glob",

		"-i", "frames/frame_*.jpg",

		"-c:v", "libx264",

		"-pix_fmt", "yuv420p",

		"-crf", "23",

		"-y",

		outputFile,

	)

 

	output, err := cmd.CombinedOutput()

	if err != nil {

		log.Printf("Error generating timelapse: %v\nOutput: %s", err, string(output))

		return

	}

 

	log.Printf("Timelapse video created: %s", outputFile)

	log.Printf("Total frames: %d, Duration: %v",

		session.FrameCount,

		time.Since(session.StartTime).Round(time.Second))

 

	// Optionally clean up frames

	// cleanupFrames()

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

 

// ParseCaptureConfig parses capture configuration from JSON

func ParseCaptureConfig(data []byte) (CaptureConfig, error) {

	var config CaptureConfig

	err := json.Unmarshal(data, &config)

	return config, err

}