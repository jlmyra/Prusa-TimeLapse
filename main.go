package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

const (
	ServerPort = "8080"
)

func main() {
	// Create output directories if they don't exist
	if err := os.MkdirAll("output", 0755); err != nil {
		log.Fatal("Failed to create output directory:", err)
	}
	if err := os.MkdirAll("frames", 0755); err != nil {
		log.Fatal("Failed to create frames directory:", err)
	}

	// Set up HTTP routes
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/api/start", handleStart)
	http.HandleFunc("/api/stop", handleStop)
	http.HandleFunc("/api/status", handleStatus)
	http.HandleFunc("/api/videos", handleVideos)
	http.HandleFunc("/api/download/", handleDownload)
	http.HandleFunc("/api/delete/", handleDelete)
	http.HandleFunc("/api/stream", handleStream)

	// Start server
	addr := ":" + ServerPort
	fmt.Printf("🎬 Prusa-TimeLapse server starting on http://localhost:%s\n", ServerPort)
	fmt.Println("Press Ctrl+C to stop")

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

// handleHome serves the main web interface
func handleHome(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Prusa-TimeLapse</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            display: flex;
            justify-content: center;
            align-items: center;
            padding: 20px;
        }
        .container {
            background: white;
            border-radius: 20px;
            box-shadow: 0 20px 60px rgba(0,0,0,0.3);
            padding: 40px;
            max-width: 900px;
            width: 100%;
        }
        h1 {
            color: #333;
            margin-bottom: 10px;
            font-size: 2.5em;
        }
        .subtitle {
            color: #666;
            margin-bottom: 30px;
            font-size: 1.1em;
        }
        .form-group {
            margin-bottom: 20px;
        }
        label {
            display: block;
            margin-bottom: 8px;
            color: #555;
            font-weight: 500;
        }
        input[type="text"],
        input[type="number"],
        select {
            width: 100%;
            padding: 12px;
            border: 2px solid #e0e0e0;
            border-radius: 8px;
            font-size: 1em;
            transition: border-color 0.3s;
            background: white;
        }
        input:focus,
        select:focus {
            outline: none;
            border-color: #667eea;
        }
        .form-row {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 15px;
            margin-bottom: 20px;
        }
        .checkbox-label {
            display: flex;
            align-items: center;
            cursor: pointer;
            user-select: none;
        }
        .checkbox-label input[type="checkbox"] {
            width: auto;
            margin-right: 10px;
            cursor: pointer;
        }
        .checkbox-label span {
            color: #555;
        }
        .button-group {
            display: flex;
            gap: 10px;
            margin-top: 30px;
        }
        button {
            flex: 1;
            padding: 15px;
            border: none;
            border-radius: 8px;
            font-size: 1em;
            font-weight: 600;
            cursor: pointer;
            transition: all 0.3s;
        }
        .btn-start {
            background: #10b981;
            color: white;
        }
        .btn-start:hover {
            background: #059669;
            transform: translateY(-2px);
            box-shadow: 0 4px 12px rgba(16, 185, 129, 0.4);
        }
        .btn-stop {
            background: #ef4444;
            color: white;
        }
        .btn-stop:hover {
            background: #dc2626;
            transform: translateY(-2px);
            box-shadow: 0 4px 12px rgba(239, 68, 68, 0.4);
        }
        .btn-stop:disabled,
        .btn-start:disabled {
            opacity: 0.5;
            cursor: not-allowed;
            transform: none;
        }
        .status {
            margin-top: 20px;
            padding: 15px;
            border-radius: 8px;
            background: #f3f4f6;
            color: #374151;
        }
        .status.active {
            background: #d1fae5;
            color: #065f46;
        }
        .emoji {
            font-size: 1.5em;
            margin-right: 10px;
        }
        .videos-section {
            margin-top: 30px;
            padding-top: 30px;
            border-top: 2px solid #e0e0e0;
        }
        .videos-section h2 {
            color: #333;
            margin-bottom: 15px;
            font-size: 1.5em;
        }
        .video-list {
            max-height: 300px;
            overflow-y: auto;
        }
        .video-item {
            display: flex;
            justify-content: space-between;
            align-items: center;
            padding: 12px;
            margin-bottom: 10px;
            background: #f9fafb;
            border-radius: 8px;
            transition: background 0.2s;
            width: auto;
            white-space: nowrap;
        }
        .video-item:hover {
            background: #f3f4f6;
        }
        .video-info {
            flex: 1;
        }
        .video-name {
            font-weight: 600;
            color: #333;
            margin-bottom: 4px;
        }
        .video-meta {
            font-size: 0.85em;
            color: #666;
        }
        .video-actions {
            display: flex;
            gap: 8px;
        }
        .btn-small {
            padding: 8px 12px;
            font-size: 0.9em;
            border: none;
            border-radius: 6px;
            cursor: pointer;
            transition: all 0.2s;
            text-decoration: none;
            color: white;
        }
        .btn-download {
            background: #3b82f6;
        }
        .btn-download:hover {
            background: #2563eb;
        }
        .btn-delete {
            background: #ef4444;
        }
        .btn-delete:hover {
            background: #dc2626;
        }
        .no-videos {
            text-align: center;
            padding: 30px;
            color: #999;
        }
        .preview-section {
            margin-top: 20px;
            padding: 15px;
            background: #f9fafb;
            border-radius: 12px;
        }
        .preview-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 10px;
        }
        .preview-header h3 {
            color: #333;
            font-size: 1.2em;
            margin: 0;
        }
        .preview-toggle {
            padding: 6px 12px;
            background: #667eea;
            color: white;
            border: none;
            border-radius: 6px;
            cursor: pointer;
            font-size: 0.9em;
            transition: background 0.2s;
        }
        .preview-toggle:hover {
            background: #5568d3;
        }
        .preview-container {
            display: none;
            margin-top: 10px;
        }
        .preview-container.active {
            display: block;
        }
        .preview-container img {
            width: 100%;
            border-radius: 8px;
            background: #000;
        }
        .preview-info {
            margin-top: 8px;
            font-size: 0.85em;
            color: #666;
            text-align: center;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>🎬 Prusa-TimeLapse</h1>
        <p class="subtitle">Create time-lapse videos from your Prusa Buddy Camera</p>

        <div class="form-group">
            <label for="rtspUrl">RTSP Stream URL</label>
            <input type="text" id="rtspUrl" placeholder="rtsp://192.168.1.251/live"
                   value="rtsp://192.168.1.251/live">
        </div>

        <div class="form-group">
            <label for="interval">Capture Interval (seconds)</label>
            <input type="number" id="interval" min="1" max="3600" value="5">
        </div>

        <div class="form-row">
            <div class="form-group">
                <label for="fps">Video FPS</label>
                <select id="fps">
                    <option value="15">15 fps</option>
                    <option value="24">24 fps (Film)</option>
                    <option value="30" selected>30 fps (Default)</option>
                    <option value="60">60 fps (Smooth)</option>
                </select>
            </div>

            <div class="form-group">
                <label for="quality">Video Quality</label>
                <select id="quality">
                    <option value="high">High</option>
                    <option value="medium" selected>Medium</option>
                    <option value="low">Low (Smaller file)</option>
                </select>
            </div>
        </div>

        <div class="form-group">
            <label class="checkbox-label">
                <input type="checkbox" id="cleanupFrames" checked>
                <span>Auto-delete frames after video generation</span>
            </label>
        </div>

        <div class="button-group">
            <button class="btn-start" id="startBtn" onclick="startCapture()">Start Recording</button>
            <button class="btn-stop" id="stopBtn" onclick="stopCapture()" disabled>Stop Recording</button>
        </div>

        <div class="status" id="status">
            <span class="emoji">⏸️</span>
            <strong>Status:</strong> Idle
        </div>

        <div class="preview-section">
            <div class="preview-header">
                <h3>📷 Live Preview</h3>
                <button class="preview-toggle" onclick="togglePreview()">Show Preview</button>
            </div>
            <div class="preview-container" id="previewContainer">
                <img id="cameraStream" src="" alt="Camera stream will appear here">
                <div class="preview-info">Live stream from your Prusa camera (10 fps)</div>
            </div>
        </div>

        <div class="videos-section">
            <h2>📹 Your Timelapses</h2>
            <div class="video-list" id="videoList">
                <div class="no-videos">No videos yet. Create your first timelapse!</div>
            </div>
        </div>
    </div>

    <script>
        let statusInterval;

        function startCapture() {
            const rtspUrl = document.getElementById('rtspUrl').value;
            const interval = document.getElementById('interval').value;
            const fps = document.getElementById('fps').value;
            const quality = document.getElementById('quality').value;
            const cleanupFrames = document.getElementById('cleanupFrames').checked;

            if (!rtspUrl) {
                alert('Please enter an RTSP URL');
                return;
            }

            // Disable start button and show loading state
            document.getElementById('startBtn').disabled = true;
            document.getElementById('startBtn').textContent = 'Connecting...';

            fetch('/api/start', {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify({
                    rtspUrl: rtspUrl,
                    interval: parseInt(interval),
                    fps: parseInt(fps),
                    quality: quality,
                    cleanupFrames: cleanupFrames
                })
            })
            .then(res => res.json())
            .then(data => {
                if (data.success) {
                    document.getElementById('startBtn').textContent = 'Start Recording';
                    document.getElementById('stopBtn').disabled = false;
                    updateStatus();
                    statusInterval = setInterval(updateStatus, 2000);
                } else {
                    document.getElementById('startBtn').disabled = false;
                    document.getElementById('startBtn').textContent = 'Start Recording';
                    alert('Failed to start: ' + data.message);
                }
            })
            .catch(err => {
                document.getElementById('startBtn').disabled = false;
                document.getElementById('startBtn').textContent = 'Start Recording';
                alert('Error: ' + err.message);
            });
        }

        function stopCapture() {
            fetch('/api/stop', {method: 'POST'})
            .then(res => res.json())
            .then(data => {
                if (data.success) {
                    document.getElementById('startBtn').disabled = false;
                    document.getElementById('stopBtn').disabled = true;
                    clearInterval(statusInterval);
                    updateStatus();
                    // Refresh video list after a short delay (video generation takes time)
                    setTimeout(loadVideos, 3000);
                }
            })
            .catch(err => {
                alert('Error: ' + err.message);
            });
        }

        function updateStatus() {
            fetch('/api/status')
            .then(res => res.json())
            .then(data => {
                const statusDiv = document.getElementById('status');
                if (data.running) {
                    statusDiv.className = 'status active';
                    statusDiv.innerHTML =
                        '<span class="emoji">🎥</span>' +
                        '<strong>Status:</strong> Recording | ' +
                        '<strong>Frames:</strong> ' + data.frameCount + ' | ' +
                        '<strong>Duration:</strong> ' + data.duration;
                } else {
                    statusDiv.className = 'status';
                    statusDiv.innerHTML = '<span class="emoji">⏸️</span><strong>Status:</strong> Idle';
                }
            });
        }

        function loadVideos() {
            fetch('/api/videos')
            .then(res => res.json())
            .then(data => {
                const videoList = document.getElementById('videoList');
                if (data.videos && data.videos.length > 0) {
                    videoList.innerHTML = data.videos.map(video =>
                        '<div class="video-item">' +
                            '<div class="video-info">' +
                                '<div class="video-name">' + video.name + '</div>' +
                                '<div class="video-meta">' + video.size + ' • ' + video.date + '</div>' +
                            '</div>' +
                            '<div class="video-actions">' +
                                '<a href="/api/download/' + video.name + '" class="btn-small btn-download" download>Download</a>' +
                                '<button class="btn-small btn-delete" onclick="deleteVideo(\'' + video.name + '\')">Delete</button>' +
                            '</div>' +
                        '</div>'
                    ).join('');
                } else {
                    videoList.innerHTML = '<div class="no-videos">No videos yet. Create your first timelapse!</div>';
                }
            })
            .catch(err => {
                console.error('Error loading videos:', err);
            });
        }

        function deleteVideo(filename) {
            if (!confirm('Are you sure you want to delete ' + filename + '?')) {
                return;
            }

            fetch('/api/delete/' + filename, {method: 'DELETE'})
            .then(res => res.json())
            .then(data => {
                if (data.success) {
                    loadVideos();
                } else {
                    alert('Failed to delete video: ' + data.message);
                }
            })
            .catch(err => {
                alert('Error deleting video: ' + err.message);
            });
        }

        function formatBytes(bytes) {
            if (bytes === 0) return '0 Bytes';
            const k = 1024;
            const sizes = ['Bytes', 'KB', 'MB', 'GB'];
            const i = Math.floor(Math.log(bytes) / Math.log(k));
            return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i];
        }

        function togglePreview() {
            const container = document.getElementById('previewContainer');
            const button = document.querySelector('.preview-toggle');
            const img = document.getElementById('cameraStream');
            const rtspUrl = document.getElementById('rtspUrl').value;

            if (container.classList.contains('active')) {
                // Hide preview
                container.classList.remove('active');
                button.textContent = 'Show Preview';
                img.src = ''; // Stop stream
            } else {
                // Show preview
                container.classList.add('active');
                button.textContent = 'Hide Preview';
                // Start stream with current RTSP URL
                img.src = '/api/stream?url=' + encodeURIComponent(rtspUrl) + '&t=' + new Date().getTime();
            }
        }

        // Update status and videos on page load
        updateStatus();
        loadVideos();

        // Refresh video list every 30 seconds
        setInterval(loadVideos, 30000);
    </script>
</body>
</html>`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

// handleStart starts the time-lapse capture
func handleStart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read and parse request body
	var config CaptureConfig
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"success": false, "message": "Invalid request: %s"}`, err.Error())
		return
	}

	// Validate configuration
	if config.RTSPUrl == "" {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"success": false, "message": "RTSP URL is required"}`)
		return
	}
	if config.Interval < 1 {
		config.Interval = 5 // Default to 5 seconds
	}

	// Start capture
	if err := StartCapture(config); err != nil {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"success": false, "message": "Failed to start capture: %s"}`, err.Error())
		return
	}

	log.Printf("Started capture: URL=%s, Interval=%ds", config.RTSPUrl, config.Interval)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, `{"success": true, "message": "Capture started successfully"}`)
}

// handleStop stops the time-lapse capture
func handleStop(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := StopCapture(); err != nil {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"success": false, "message": "Failed to stop capture: %s"}`, err.Error())
		return
	}

	log.Println("Stopped capture, generating timelapse video...")
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, `{"success": true, "message": "Capture stopped, generating video..."}`)
}

// handleStatus returns the current status
func handleStatus(w http.ResponseWriter, r *http.Request) {
	status := GetStatus()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(status); err != nil {
		log.Printf("Error encoding status: %v", err)
		fmt.Fprint(w, `{"running": false, "frameCount": 0, "duration": "0s"}`)
	}
}

// handleVideos lists all timelapse videos
func handleVideos(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir("output")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"videos": []}`)
		return
	}

	type VideoInfo struct {
		Name string `json:"name"`
		Size string `json:"size"`
		Date string `json:"date"`
	}

	var videos []VideoInfo
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".mp4") {
			continue
		}

		info, err := file.Info()
		if err != nil {
			continue
		}

		videos = append(videos, VideoInfo{
			Name: file.Name(),
			Size: formatBytes(info.Size()),
			Date: info.ModTime().Format("Jan 2, 2006 3:04 PM"),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"videos": videos})
}

// handleDownload serves video files for download
func handleDownload(w http.ResponseWriter, r *http.Request) {
	filename := strings.TrimPrefix(r.URL.Path, "/api/download/")
	if filename == "" {
		http.Error(w, "Filename required", http.StatusBadRequest)
		return
	}

	// Security: prevent directory traversal
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") {
		http.Error(w, "Invalid filename", http.StatusBadRequest)
		return
	}

	filepath := "output/" + filename
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", "video/mp4")
	http.ServeFile(w, r, filepath)
}

// handleDelete deletes a video file
func handleDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	filename := strings.TrimPrefix(r.URL.Path, "/api/delete/")
	if filename == "" {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"success": false, "message": "Filename required"}`)
		return
	}

	// Security: prevent directory traversal
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"success": false, "message": "Invalid filename"}`)
		return
	}

	filepath := "output/" + filename
	if err := os.Remove(filepath); err != nil {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"success": false, "message": "Failed to delete file: %s"}`, err.Error())
		return
	}

	log.Printf("Deleted video: %s", filename)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, `{"success": true}`)
}

// formatBytes converts bytes to human-readable format
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	sizes := []string{"KB", "MB", "GB", "TB"}
	return fmt.Sprintf("%.1f %s", float64(bytes)/float64(div), sizes[exp])
}

// handleStream streams MJPEG from RTSP camera
func handleStream(w http.ResponseWriter, r *http.Request) {
	rtspUrl := r.URL.Query().Get("url")
	if rtspUrl == "" {
		rtspUrl = "rtsp://192.168.1.251/live" // Default camera URL
	}

	// Set headers for MJPEG stream
	w.Header().Set("Content-Type", "multipart/x-mixed-replace; boundary=frame")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	w.Header().Set("Pragma", "no-cache")

	w.Header().Set("Expires", "0")

	// Start ffmpeg to convert RTSP to MJPEG

	cmd := exec.Command("ffmpeg",

		"-rtsp_transport", "tcp",

		"-i", rtspUrl,

		"-f", "mjpeg",

		"-q:v", "3", // Quality (2-31, lower is better)

		"-r", "5", // 5 fps for stream (reduced for stability)

		"-vf", "scale=640:-1", // Scale down for better performance

		"-",
	)

	stdout, err := cmd.StdoutPipe()

	if err != nil {

		log.Printf("Error creating stdout pipe: %v", err)

		http.Error(w, "Failed to create stream", http.StatusInternalServerError)

		return

	}

	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {

		log.Printf("Error starting ffmpeg: %v", err)

		http.Error(w, "Failed to start stream", http.StatusInternalServerError)

		return

	}

	defer func() {

		cmd.Process.Kill()

		cmd.Wait()

	}()

	// Log ffmpeg errors in background

	go func() {

		if stderr != nil {

			io.Copy(os.Stderr, stderr)

		}

	}()

	// Read MJPEG stream frame by frame

	reader := stdout

	buffer := make([]byte, 0, 512*1024) // 512KB buffer for accumulating data

	tempBuf := make([]byte, 4096) // 4KB temporary read buffer

	for {

		// Check if client disconnected

		select {

		case <-r.Context().Done():

			log.Println("Client disconnected from stream")

			return

		default:

		}

		// Read chunk from ffmpeg

		n, err := reader.Read(tempBuf)

		if err != nil {

			if err != io.EOF {

				log.Printf("Error reading stream: %v", err)

			}

			break

		}

		// Append to buffer

		buffer = append(buffer, tempBuf[:n]...)

		// Look for complete JPEG frames (starts with 0xFF 0xD8, ends with 0xFF 0xD9)

		for {

			// Find JPEG start marker

			startIdx := -1

			for i := 0; i < len(buffer)-1; i++ {

				if buffer[i] == 0xFF && buffer[i+1] == 0xD8 {

					startIdx = i

					break

				}

			}

			if startIdx == -1 {

				// No start marker found, keep first byte and discard rest

				if len(buffer) > 1 {

					buffer = buffer[len(buffer)-1:]

				}

				break

			}

			// Find JPEG end marker after start

			endIdx := -1

			for i := startIdx + 2; i < len(buffer)-1; i++ {

				if buffer[i] == 0xFF && buffer[i+1] == 0xD9 {

					endIdx = i + 2 // Include the end marker

					break

				}

			}

			if endIdx == -1 {

				// Incomplete frame, wait for more data

				// But keep buffer from start marker

				buffer = buffer[startIdx:]

				break

			}

			// Extract complete JPEG frame

			frame := buffer[startIdx:endIdx]

			// Write frame to client

			_, err := fmt.Fprintf(w, "--frame\r\nContent-Type: image/jpeg\r\nContent-Length: %d\r\n\r\n", len(frame))

			if err != nil {

				log.Printf("Error writing frame header: %v", err)

				return

			}

			_, err = w.Write(frame)

			if err != nil {

				log.Printf("Error writing frame data: %v", err)

				return

			}

			_, err = fmt.Fprint(w, "\r\n")

			if err != nil {

				return

			}

			// Flush the response

			if flusher, ok := w.(http.Flusher); ok {

				flusher.Flush()

			}

			// Remove processed frame from buffer

			buffer = buffer[endIdx:]

		}

		// Prevent buffer from growing too large

		if len(buffer) > 1024*1024 { // 1MB limit

			log.Println("Buffer too large, resetting")

			buffer = buffer[len(buffer)-4096:] // Keep last 4KB

		}

	}

	log.Println("Stream ended")

}
