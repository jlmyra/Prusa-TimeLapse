# Prusa-TimeLapse

 

A Go-based application to create time-lapse videos from Prusa Buddy Camera RTSP streams.

 

## Features

 

- 🎥 Capture frames from Prusa Buddy Camera RTSP stream

- ⏱️ Configurable capture intervals

- 🎬 Generate MP4 time-lapse videos

- 🌐 Web-based interface for easy control

- 💻 Runs locally on macOS

 

## Prerequisites

 

- Go 1.21 or higher

- FFmpeg installed on your system

 

### Installing FFmpeg on macOS

 

```bash

brew install ffmpeg

```

 

## Installation

 

```bash

git clone https://github.com/jlmyra/Prusa-TimeLapse.git

cd Prusa-TimeLapse

go mod download

```

 

## Usage

 

```bash

go run main.go

```

 

Then open your browser to `http://localhost:8080`

 

## Configuration

 

- **RTSP URL**: Enter your Prusa Buddy Camera RTSP stream URL

- **Capture Interval**: Set how often to capture frames (in seconds)

- **Output Directory**: Where to save time-lapse videos

 

## Project Status

 

🚧 Under active development

A lightweight Go application to create beautiful time-lapse videos from Prusa Buddy Camera RTSP streams.

 

## Features

 

- 🎥 **Live RTSP Capture** - Connects to Prusa Buddy Camera streams

- ⏱️ **Configurable Intervals** - Set capture rate from 1-3600 seconds

- 🎬 **Automatic MP4 Generation** - Creates timelapse on stop

- 🌐 **Beautiful Web UI** - Modern, responsive interface

- 📊 **Real-time Status** - Live frame count and duration tracking

- 🔒 **Thread-safe** - Concurrent-safe capture sessions

- 💻 **Native macOS** - Optimized for Mac, single binary

 

## Quick Start

 

### Prerequisites

 

- **Go 1.21+** - [Download here](https://go.dev/dl/)

- **FFmpeg** - Install via Homebrew:

  ```bash

  brew install ffmpeg

  ```

 

### Installation

 

```bash

git clone https://github.com/jlmyra/Prusa-TimeLapse.git

cd Prusa-TimeLapse

go build -o prusa-timelapse

```

 

### Run

 

```bash

./prusa-timelapse

```

 

Open your browser to **http://localhost:8080**

 

## Usage

 

📖 **[Read the complete usage guide](USAGE.md)** for detailed instructions, troubleshooting, and tips.

 

**Quick workflow:**

1. Enter your camera's RTSP URL (e.g., `rtsp://192.168.1.100:8554/stream`)

2. Set capture interval (e.g., 30 seconds for a 4-hour print)

3. Click "Start Recording"

4. Click "Stop Recording" when done

5. Find your timelapse in `output/timelapse_YYYY-MM-DD_HH-MM-SS.mp4`

 

## Project Architecture

 

```

main.go     - HTTP server, web UI, API endpoints

capture.go  - RTSP capture logic, FFmpeg integration, timelapse generation

frames/     - Captured JPEG frames (auto-created)

output/     - Generated MP4 timelapses (auto-created)

```

 

**Key Go Concepts Used:**

- Goroutines for background processing

- Channels for stop signaling

- Mutexes for thread safety

- HTTP server for web interface

- JSON API for frontend communication

- Command execution for FFmpeg

 

## Development Status

 

✅ **Fully Functional** - Ready to use for production timelapses!

 

**Future enhancements:**

- [ ] Multiple concurrent captures

- [ ] Frame cleanup options

- [ ] Custom video settings (FPS, quality, resolution)

- [ ] Dark mode UI

- [ ] Download video from web UI

 

## Contributing

 

This is a learning project! Contributions, suggestions, and feedback welcome.

## License

 

MIT License