package main

import (
	"fmt"

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

	fmt.Println("UARTReader started, waiting for frames...")

	for frame := range parser.ParseFrames(reader.OutChan()) {
		fmt.Println("Received frame with TLVs:", len(frame.TLVs))
		for _, tlv := range frame.TLVs {
			if tlv.Header.Type == parser.EXT_TARGET_LIST {
				targets, err := tlv.AsTargetList()
				if err != nil {
					fmt.Println("Error decoding target list:", err)
					continue
				}
				for _, target := range targets {
					fmt.Printf("Target ID: %d, X: %.2f, Y: %.2f, Confidence: %.2f\n", target.ID, target.X, target.Y, target.Confidence)
				}
			}
		}
	}

	fmt.Println("UARTReader stopped")
}
