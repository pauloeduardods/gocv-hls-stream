# Webcam HLS Streaming

This project captures video from your webcam, processes it using FFmpeg, and serves it as an HLS stream, which can be viewed in a web browser.
## How to Run the Project

1. Install Dependencies

    Make sure you have the following installed:
      - Go
      - OpenCV with Go bindings (gocv)
      - FFmpeg

   Install the Go bindings for OpenCV:

    `go get -u -d gocv.io/x/gocv`
   
3. Build the Project

    `go build -o webcam-hls main.go`

4. Run the Application

    Execute the binary to start the video capture and streaming:

    `./webcam-hls`

5. Open the Stream in a Browser

    Open your browser and go to:

      `http://localhost:8080/index.html`

    The video stream should start playing automatically.

## Folder Structure

    main.go: Go application for capturing video and handling streaming.
    index.html: HTML page for viewing the HLS stream.
    player.js: JavaScript file to handle HLS playback using hls.js.
    stream/: Directory where HLS segments are stored.
