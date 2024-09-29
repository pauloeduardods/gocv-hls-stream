package main

import (
	"bytes"
	"image/jpeg"
	"log"
	"net/http"
	"os"
	"os/exec"

	"gocv.io/x/gocv"
)

func main() {
	if _, err := os.Stat("stream"); os.IsNotExist(err) {
		err := os.Mkdir("stream", 0755)
		if err != nil {
			log.Fatalf("Error when creating stream folder: %v", err)
		}
	}

	pipeName := "video_pipe"
	if _, err := os.Stat(pipeName); os.IsNotExist(err) {
		cmd := exec.Command("mkfifo", pipeName)
		err := cmd.Run()
		if err != nil {
			log.Fatalf("Error when create pipe: %v", err)
		}
	}

	ffmpegCmd := exec.Command("ffmpeg",
		"-y",
		"-f", "mjpeg",
		"-framerate", "30",
		"-i", pipeName,
		"-c:v", "libx264",
		"-preset", "ultrafast",
		"-f", "hls",
		"-hls_time", "1",
		"-hls_list_size", "5",
		"-hls_flags", "delete_segments",
		"stream/stream.m3u8",
	)

	ffmpegCmd.Stderr = os.Stderr
	ffmpegCmd.Stdout = os.Stdout

	err := ffmpegCmd.Start()
	if err != nil {
		log.Fatalf("Error when start ffmpeg: %v", err)
	}
	defer ffmpegCmd.Process.Kill()

	go func() {
		http.Handle("/", http.FileServer(http.Dir(".")))
		log.Println("Listen HTTP in :8080")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	webcam, err := gocv.OpenVideoCapture(0)
	if err != nil {
		log.Fatalf("Error opening webcam %v", err)
	}
	defer webcam.Close()

	webcam.Set(gocv.VideoCaptureFrameWidth, 640)
	webcam.Set(gocv.VideoCaptureFrameHeight, 480)
	webcam.Set(gocv.VideoCaptureFPS, 30)
	webcam.Set(gocv.VideoCaptureFOURCC, float64(webcam.ToCodec("MJPG")))

	pipeFile, err := os.OpenFile(pipeName, os.O_WRONLY, os.ModeNamedPipe)
	if err != nil {
		log.Fatalf("Error opening pipe: %v", err)
	}
	defer pipeFile.Close()

	img := gocv.NewMat()
	defer img.Close()

	for {
		if ok := webcam.Read(&img); !ok {
			log.Println("Error when reading camera")
			continue
		}
		if img.Empty() {
			log.Println("Error Empty image")

			continue
		}

		image, err := img.ToImage()
		if err != nil {
			log.Println("Error when getting image.Image from gocv.Mat", err)

			continue
		}

		buffer := new(bytes.Buffer)
		if err := jpeg.Encode(buffer, image, &jpeg.Options{Quality: 75}); err != nil {
			log.Println("Error when getting jpeg from image.Image", err)
			continue
		}

		data := buffer.Bytes()

		_, err = pipeFile.Write(data)
		if err != nil {
			log.Fatalf("Error writing frame in pipe: %v", err)
		}
	}
}
