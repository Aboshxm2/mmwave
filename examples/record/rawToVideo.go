package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	width := 640
	height := 480
	frameRate := 30

	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer f.Close()

	cmd := exec.Command("ffmpeg",
		"-f", "rawvideo",
		"-pix_fmt", "rgb24",
		"-s", // size
		fmt.Sprintf("%dx%d", width, height),
		"-r", fmt.Sprintf("%d", frameRate),
		"-i", "pipe:0",
		"-c:v", "libx264",
		"-pix_fmt", "yuv420p",
		"actual.mp4",
	)

	cmd.Stdin = f
	cmd.Stderr = os.Stderr // optional: to see ffmpeg logs

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
