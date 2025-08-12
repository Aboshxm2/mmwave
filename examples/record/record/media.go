package record

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/BurntSushi/toml"
)

type mediaConfig struct {
	Width       int    `toml:"width"`
	Height      int    `toml:"height"`
	Fps         int    `toml:"fps"`
	InputFormat string `toml:"inputFormat"`
	PixelFormat string `toml:"pixelFormat"` // TODO remove these constants
	DeviceName  string `toml:"deviceName"`  // TODO I think linux doesnt use "video="

}

var frameBytes = 640 * 480 * 3 // TODO

func ffmpegArgsForDevice() []string {

	var c mediaConfig
	_, err := toml.DecodeFile("./config.toml", &c)
	if err != nil {
		panic(err)
	}
	config := c
	frameBytes = config.Width * config.Height * 3

	return []string{
		"-f", config.InputFormat,
		"-i", config.DeviceName,
		"-pix_fmt", config.PixelFormat,
		"-s", fmt.Sprintf("%dx%d", config.Width, config.Height),
		"-r", strconv.Itoa(config.Fps),
		"-f", "rawvideo",
		"-",
	}
}

func StartFFmpeg() (*exec.Cmd, *bufio.Reader, error) {
	args := ffmpegArgsForDevice()
	cmd := exec.Command("ffmpeg", args...)
	fmt.Println(args)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}

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

func ReadFrames(reader *bufio.Reader) <-chan []byte {
	ch := make(chan []byte, 10)
	go func() {
		defer close(ch)
		for {
			frameData, err := readFrame(reader)
			if err != nil {
				return
			}
			ch <- frameData
		}
	}()

	return ch
}
