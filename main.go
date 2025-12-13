package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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
            max-width: 600px;
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
        input[type="number"] {
            width: 100%;
            padding: 12px;
            border: 2px solid #e0e0e0;
            border-radius: 8px;
            font-size: 1em;
            transition: border-color 0.3s;
        }
        input:focus {
            outline: none;
            border-color: #667eea;
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
    </style>
</head>
<body>
    <div class="container">
        <h1>🎬 Prusa-TimeLapse</h1>
        <p class="subtitle">Create time-lapse videos from your Prusa Buddy Camera</p>

        <div class="form-group">
            <label for="rtspUrl">RTSP Stream URL</label>
            <input type="text" id="rtspUrl" placeholder="rtsp://192.168.1.100:8554/stream"
                   value="rtsp://192.168.1.100:8554/stream">
        </div>

        <div class="form-group">
            <label for="interval">Capture Interval (seconds)</label>
            <input type="number" id="interval" min="1" max="3600" value="5">
        </div>

        <div class="button-group">
            <button class="btn-start" id="startBtn" onclick="startCapture()">Start Recording</button>
            <button class="btn-stop" id="stopBtn" onclick="stopCapture()" disabled>Stop Recording</button>
        </div>

        <div class="status" id="status">
            <span class="emoji">⏸️</span>
            <strong>Status:</strong> Idle
        </div>
    </div>

    <script>
        let statusInterval;

        function startCapture() {
            const rtspUrl = document.getElementById('rtspUrl').value;
            const interval = document.getElementById('interval').value;

            if (!rtspUrl) {
                alert('Please enter an RTSP URL');
                return;
            }

            fetch('/api/start', {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify({
                    rtspUrl: rtspUrl,
                    interval: parseInt(interval)
                })
            })
            .then(res => res.json())
            .then(data => {
                if (data.success) {
                    document.getElementById('startBtn').disabled = true;
                    document.getElementById('stopBtn').disabled = false;
                    updateStatus();
                    statusInterval = setInterval(updateStatus, 2000);
                } else {
                    alert('Failed to start: ' + data.message);
                }
            })
            .catch(err => {
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

        // Update status on page load
        updateStatus();
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
