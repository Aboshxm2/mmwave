package record

import (
	"fmt"
	"os"
)

type Entry struct {
	FrameNumber, TID                                        uint32
	X, Y, Z, VelX, VelY, VelZ, AccX, AccY, AccZ, Confidence float32
}

type Logger struct {
	deviceCh <-chan []byte
	framesCh <-chan Entry

	videoFile  *os.File
	framesFile *os.File
}

func NewLogger(deviceCh <-chan []byte, framesCh <-chan Entry, path string) *Logger {
	v, err := os.Create(path + "raw_video")
	if err != nil {
		panic(err)
	}
	f, err := os.Create(path + "frames.csv")
	if err != nil {
		panic(err)
	}
	f.WriteString("frame,tid,posX,posY,posZ,velX,velY,velZ,accX,accY,accZ,confidenceLevel\n")

	return &Logger{deviceCh, framesCh, v, f}
}

func (l *Logger) Start() chan<- interface{} {
	closeCh := make(chan interface{})
	go func() {
		for {
			select {
			case <-closeCh:
				l.videoFile.Close()
				return
			case bytes := <-l.deviceCh:
				l.videoFile.Write(bytes)
			}
		}
	}()
	go func() {
		for {
			select {
			case <-closeCh:
				l.framesFile.Close()
				return
			case entry := <-l.framesCh:
				fmt.Fprintf(l.framesFile, "%d,%d,%.3f,%.3f,%.3f,%.3f,%.3f,%.3f,%.3f,%.3f,%.3f,%.3f,\n",
					entry.FrameNumber, entry.TID, entry.X, entry.Y, entry.Z, entry.VelX, entry.VelY,
					entry.VelZ, entry.AccX, entry.AccY, entry.AccZ, entry.Confidence)
			}
		}
	}()
	return closeCh
}
