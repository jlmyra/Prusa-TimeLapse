#!/bin/bash

 

# Prusa-TimeLapse Camera Test Script

# This script helps verify your RTSP camera connection

 

echo "🎥 Prusa-TimeLapse Camera Test"

echo "================================"

echo ""

 

# Check if FFmpeg is installed

echo "1️⃣ Checking FFmpeg installation..."

if command -v ffmpeg &> /dev/null; then

    echo "   ✅ FFmpeg is installed"

    ffmpeg -version | head -n 1

else

    echo "   ❌ FFmpeg is NOT installed"

    echo "   Install with: brew install ffmpeg"

    exit 1

fi

 

echo ""

 

# Prompt for RTSP URL

echo "2️⃣ Testing RTSP stream connection..."

read -p "   Enter your Prusa camera RTSP URL (e.g., rtsp://192.168.1.100:8554/stream): " RTSP_URL

 

if [ -z "$RTSP_URL" ]; then

    echo "   ❌ No URL provided"

    exit 1

fi

 

echo "   Testing connection to: $RTSP_URL"

echo "   This may take a few seconds..."

echo ""

 

# Try to capture a test frame

ffmpeg -rtsp_transport tcp -i "$RTSP_URL" -vframes 1 -q:v 2 -y test-frame.jpg 2>&1 | tail -n 5

 

if [ -f "test-frame.jpg" ]; then

    echo ""

    echo "   ✅ SUCCESS! Test frame captured"

    echo "   📸 Saved as: test-frame.jpg"

    echo ""

    echo "   You can now use this URL in Prusa-TimeLapse!"

    echo "   Run: ./prusa-timelapse"

    echo "   Then open: http://localhost:8080"

else

    echo ""

    echo "   ❌ FAILED to capture frame"

    echo ""

    echo "   Troubleshooting:"

    echo "   • Verify camera IP address is correct"

    echo "   • Check that port 8554 is accessible"

    echo "   • Ensure camera is powered on and connected to network"

    echo "   • Try pinging the camera: ping <camera-ip>"

fi

 

echo ""