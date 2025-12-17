# Prusa-TimeLapse***

A lightweight Go application to create beautiful time-lapse videos from Prusa Buddy Camera RTSP streams.

## Features

- 🎥 **Live RTSP Capture** - Connects to Prusa Buddy Camera streams
- 📹 **Live Camera Preview** - View real-time MJPEG stream from your camera
- ⏱️ **Configurable Intervals** - Set capture rate from 1-3600 seconds
- 🎬 **Automatic MP4 Generation** - Creates timelapse on stop with H.264 encoding
- ⚙️ **Customizable Video Settings** - Choose FPS (15/24/30/60) and quality (High/Medium/Low)
- 🧹 **Smart Frame Management** - Automatic frame cleanup after video generation
- 🌐 **Beautiful Web UI** - Modern, responsive interface with 900px wide viewport
- 📊 **Real-time Status** - Live frame count and duration tracking
- 📂 **Video Management** - List, download, and delete timelapses from the web UI
- ✅ **Connection Validation** - Pre-flight RTSP testing before capture starts
- 🔒 **Thread-safe** - Concurrent-safe capture sessions
- 💻 **Native macOS** - Optimized for Mac, single binary

## Quick Start

### Prerequisites

- **Go 1.21+** - [Download here](https://go.dev/dl/)
- **FFmpeg** - Install via Homebrew:
  ```bash
  brew install ffmpeg

Installation
git clone https://github.com/jlmyra/Prusa-TimeLapse.git
cd Prusa-TimeLapse
go build -o prusa-timelapse

Run
./prusa-timelapse

Open your browser to http://localhost:8080
Usage
📖 Read the complete usage guide for detailed instructions, troubleshooting, and tips.
Quick workflow:
Enter your camera's RTSP URL (default: rtsp://192.168.1.251/live)
(Optional) Click "Show Preview" to view live camera feed
Set capture interval (e.g., 30 seconds for a 4-hour print)
Choose video FPS and quality settings
Click "Start Recording"
Click "Stop Recording" when done
Download or delete videos from the "Your Timelapses" section
Project Architecture
main.go     - HTTP server, web UI, API endpoints, MJPEG streaming
capture.go  - RTSP capture logic, FFmpeg integration, timelapse generation
frames/     - Captured JPEG frames (auto-created)
output/     - Generated MP4 timelapses (auto-created)

API Endpoints:
GET / - Main web interface
POST /api/start - Start timelapse capture
POST /api/stop - Stop capture and generate video
GET /api/status - Get current capture status
GET /api/videos - List all generated videos
GET /api/download/:filename - Download video file
DELETE /api/delete/:filename - Delete video file
GET /api/stream - Live MJPEG stream from camera
Key Go Concepts Used:
Goroutines for background processing and streaming
Channels for stop signaling
Mutexes (RWMutex) for thread safety
HTTP server with multipart streaming
JSON API for frontend communication
Command execution for FFmpeg
JPEG frame boundary detection
Development Status
✅ Fully Functional - Ready to use for production timelapses!
Recently completed:
✅ Live camera preview with MJPEG streaming
✅ Frame cleanup options
✅ Custom video settings (FPS, quality)
✅ Download/delete videos from web UI
✅ Connection validation before capture
✅ Improved UI with larger preview window
Potential future enhancements:
 Multiple concurrent captures from different cameras
 Scheduled/timed captures
 Email notifications on completion
 Dark mode UI
 Motion detection triggers
 Cloud storage integration
Contributing
This is a learning project! Contributions, suggestions, and feedback welcome.
License
MIT License
