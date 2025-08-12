// main.go
package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"
	"os/exec"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

const (
	// Set these to your desired capture resolution and framerate
	width      = 640
	height     = 480
	fps        = 30
	frameBytes = width * height * 3 // RGB24
)

func ffmpegArgsForDevice(deviceName string) []string {
	// example:
	// ffmpeg -f dshow -i video="Integrated Camera" -pix_fmt rgb24 -s 640x480 -r 30 -f rawvideo -
	return []string{
		"-f", "dshow",
		"-i", "video=" + deviceName,
		"-pix_fmt", "rgb24",
		"-s", fmt.Sprintf("%dx%d", width, height),
		"-r", strconv.Itoa(fps),
		"-f", "rawvideo",
		"-",
	}
}

func startFFmpeg(deviceName string) (*exec.Cmd, *bufio.Reader, error) {
	args := ffmpegArgsForDevice(deviceName)
	cmd := exec.Command("ffmpeg", args...)
	fmt.Println(args)

	// we discard ffmpeg stderr to avoid clutter; optionally you can pipe and log it
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}

	// optional: redirect ffmpeg's stderr to this process's stderr for debugging
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return nil, nil, err
	}

	reader := bufio.NewReader(stdout)
	return cmd, reader, nil
}

func readFrame(reader *bufio.Reader) ([]byte, error) {
	frame := make([]byte, frameBytes)
	n := 0
	for n < frameBytes {
		readNow, err := reader.Read(frame[n:])
		if err != nil {
			return nil, err
		}
		n += readNow
	}
	return frame, nil
}

func Rgb24ToRGBA(imgBuf []byte) *image.RGBA {
	rect := image.Rect(0, 0, width, height)
	rgba := image.NewRGBA(rect)
	// image.RGBA.Pix layout: RGBA RGBA ...
	// ffmpeg gives RGBRGB...
	// We'll copy and append Alpha=255
	out := rgba.Pix
	j := 0
	for i := 0; i < len(imgBuf); i += 3 {
		out[j+0] = imgBuf[i+0] // R
		out[j+1] = imgBuf[i+1] // G
		out[j+2] = imgBuf[i+2] // B
		out[j+3] = 0xFF        // A
		j += 4
	}
	return rgba
}

func maina() {
	// --- EDIT THIS: put your camera device name exactly as ffmpeg shows it ---
	// Example: deviceName := `Integrated Camera` or `USB Camera`
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go \"Your Camera Device Name\"")
		fmt.Println("To list devices: ffmpeg -list_devices true -f dshow -i dummy")
		return
	}
	deviceName := os.Args[1]
	// -----------------------------------------------------------------------

	// Start ffmpeg capturing
	cmd, reader, err := startFFmpeg(deviceName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to start ffmpeg: %v\n", err)
		return
	}
	defer func() {
		// attempt graceful kill
		_ = cmd.Process.Kill()
		cmd.Wait()
	}()

	// Fyne app setup
	a := app.New()
	w := a.NewWindow("Camera Preview (ffmpeg -> Fyne)")
	w.Resize(fyne.NewSize(width, height+40))

	// placeholder image (black)
	rect := image.Rect(0, 0, width, height)
	initial := image.NewRGBA(rect)
	draw.Draw(initial, rect, &image.Uniform{color.Black}, image.Point{}, draw.Src)

	img := canvas.NewImageFromImage(initial)
	img.FillMode = canvas.ImageFillContain

	infoText := canvas.NewText("Press ESC or close window to stop", color.White)
	info := container.NewMax(img, container.NewVBox(container.NewHBox(infoText)))

	w.SetContent(info)
	w.Show()

	// Frame reading goroutine
	go func() {
		// buffer for a single frame
		for {
			frameBuf, err := readFrame(reader)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error reading frame: %v\n", err)
				// stop; schedule closing GUI
				a.Driver().DoFromGoroutine(func() {
					// show a simple error image
					msg := canvas.NewText("Stream ended / error", color.RGBA{R: 255, A: 255})
					w.SetContent(container.NewVBox(msg))
				}, false)
				return
			}

			// convert to RGBA image
			rgba := Rgb24ToRGBA(frameBuf)

			// Important: Fyne UI updates must be done on main thread
			// Use Driver().RunOnMain to safely update the image and refresh.
			a.Driver().DoFromGoroutine(func() {
				// Replace the image data and refresh
				img.Image = rgba
				img.Refresh()
			}, false)

			// optional: a tiny sleep to avoid hogging CPU if ffmpeg is faster
			// time.Sleep(time.Millisecond * 1)
		}
	}()

	// Close behavior: when window closes, kill ffmpeg and exit
	w.SetOnClosed(func() {
		_ = cmd.Process.Kill()
		cmd.Wait()
	})

	// Handle Esc key to close (simple)
	w.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
		if k.Name == fyne.KeyEscape {
			w.Close()
		}
	})

	// Run the GUI (blocks)
	a.Run()

	// wait a short while to let ffmpeg die
	time.Sleep(100 * time.Millisecond)
}
