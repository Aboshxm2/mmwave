package main

import (
	"fmt"
	"os"
	"record/record"
	"time"

	"github.com/Aboshxm2/mmwave/parser"
	"github.com/Aboshxm2/mmwave/serial"
)

func main() {
	reader, err := serial.NewUARTReader("COM13", 921600)
	if err != nil {
		panic(err)
	}

	reader.Start()
	defer reader.Stop()

	cmd, deviceReader, err := record.StartFFmpeg()
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = cmd.Process.Kill()
		cmd.Wait()
	}()

	deviceChan := record.ReadFrames(deviceReader)
	framesChan := parser.ParseFrames(reader.OutChan())

	var (
		logger           record.Logger
		recordUntil      time.Time
		isRunning        bool
		closeChan        chan<- interface{}
		loggerDeviceChan = make(chan []byte, 10)
		loggerFramesChan = make(chan record.Entry, 10)
	)

	for {
		select {
		case frame := <-framesChan:
			fmt.Println("yesss")
			for _, tlv := range frame.TLVs {
				if tlv.Header.Type == parser.TRACKERPROC_3D_TARGET_LIST {
					targetList, err := tlv.AsTargetList()
					if err != nil {
						panic(err)
					}
					fmt.Println(len(*targetList))
					if len(*targetList) > 0 {
						recordUntil = time.Now().Add(time.Second * 3)

						if !isRunning {
							fmt.Println("Started Recording")
							isRunning = true
							i := 0
							entries, err := os.ReadDir("./output/")
							if err != nil {
								panic(err)
							}

							for range entries {
								i++
							}

							err = os.Mkdir(fmt.Sprint("./output/", i, "/"), 0755)
							if err != nil {
								panic(err)
							}
							logger = *record.NewLogger(loggerDeviceChan, loggerFramesChan, fmt.Sprint("./output/", i, "/"))
							closeChan = logger.Start()
						}
					}
					for _, d := range *targetList {
						entry := record.Entry{
							FrameNumber: frame.Header.FrameNumber, TID: d.ID,
							X: d.X, Y: d.Y, Z: d.Z, VelX: d.VelX, VelY: d.VelY, VelZ: d.VelZ, AccX: d.AccX, AccY: d.AccY, AccZ: d.AccZ, Confidence: d.Confidence,
						}
						loggerFramesChan <- entry
					}
				}
			}
		case b := <-deviceChan:
			if isRunning {
				loggerDeviceChan <- b
			}
		default:
			if recordUntil.Before(time.Now()) && isRunning {
				closeChan <- struct{}{}
				isRunning = false
				fmt.Println("Stoped Recording")
			}
		}
	}

	// logger := record.NewLogger(record.ReadFrames(deviceReader), ), "../../output/test/")

	// ch := logger.Start()

	// time.Sleep(time.Second * 4)
	// ch <- struct{}{}
}
