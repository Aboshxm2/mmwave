package serial

import (
	"bytes"
	"io"
	"log"

	"go.bug.st/serial"
)

var magicWord = []byte{0x02, 0x01, 0x04, 0x03, 0x06, 0x05, 0x08, 0x07}

// UARTReader reads data from a serial port and emits complete raw frames.
type UARTReader struct {
	port    serial.Port
	outChan chan []byte
	readBuf []byte
}

// NewUARTReader initializes a new UARTReader.
func NewUARTReader(portName string, baudRate int) (*UARTReader, error) {
	cfg := &serial.Mode{BaudRate: baudRate}
	s, err := serial.Open(portName, cfg)
	if err != nil {
		return nil, err
	}

	return &UARTReader{
		port:    s,
		outChan: make(chan []byte, 10),
		readBuf: make([]byte, 0, 4096),
	}, nil
}

// Start begins reading from the UART port.
func (r *UARTReader) Start() {
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := r.port.Read(buf)
			if err != nil {
				if err == io.EOF {
					continue
				}
				log.Println("UART read error:", err)
				break
			}
			if n > 0 {
				r.readBuf = append(r.readBuf, buf[:n]...)
				r.processBuffer()
			}
		}
	}()
}

func (r *UARTReader) Stop() {
	if r.port != nil {
		r.port.Close()
	}
	close(r.outChan)
}

func (r *UARTReader) OutChan() <-chan []byte {
	return r.outChan
}

// processBuffer looks for complete frames starting with the magic word.
func (r *UARTReader) processBuffer() {
	for {
		idx := bytes.Index(r.readBuf, magicWord)
		if idx == -1 {
			if len(r.readBuf) > 8192 {
				// Prevent infinite growth
				r.readBuf = r.readBuf[len(r.readBuf)-len(magicWord):]
			}
			return
		}

		if idx > 0 {
			// Drop junk before magic word
			r.readBuf = r.readBuf[idx:]
		}

		// Check if we have enough bytes to determine frame size (e.g. after header)
		if len(r.readBuf) < 12+8 { // 8 = magic word + 12 = first part of header
			return // wait for more data
		}

		// Frame length is at offset 12 after magic word (magic + version + length)
		packetLen := bytesToUint32(r.readBuf[12:16])

		if len(r.readBuf) < int(packetLen) {
			return // wait for more bytes
		}

		frame := make([]byte, packetLen)
		copy(frame, r.readBuf[:packetLen])
		r.readBuf = r.readBuf[packetLen:]

		r.outChan <- frame
	}
}

func bytesToUint32(b []byte) uint32 {
	if len(b) < 4 {
		return 0
	}
	return uint32(b[0]) |
		uint32(b[1])<<8 |
		uint32(b[2])<<16 |
		uint32(b[3])<<24
}
