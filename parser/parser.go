package parser

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

var ErrInvalidMagicWord = errors.New("invalid magic word") // TODO

var ErrShortFrame = errors.New("frame too short")

func ParseFrames(readChan <-chan []byte) <-chan *Frame {
	outChan := make(chan *Frame, 10)

	go func() {
		defer close(outChan)

		for data := range readChan {
			frame, err := ParseFrame(data)
			if err != nil {
				fmt.Println("Error parsing frame:", err)
				continue
			}
			outChan <- frame
		}
	}()

	return outChan
}

func ParseFrame(data []byte) (*Frame, error) {
	if len(data) < 40 {
		return nil, fmt.Errorf("frame too short")
	}

	reader := bytes.NewReader(data[8:]) // skip magic word

	var header FrameHeader
	if err := binary.Read(reader, binary.LittleEndian, &header); err != nil {
		return nil, fmt.Errorf("failed to parse header: %w", err)
	}

	tlvs := make([]TLV, 0, header.NumTLVs)
	for i := 0; i < int(header.NumTLVs); i++ {
		var tlvHdr TLVHeader
		if err := binary.Read(reader, binary.LittleEndian, &tlvHdr); err != nil {
			return nil, fmt.Errorf("failed to read TLV header: %w", err)
		}

		tlvData := make([]byte, tlvHdr.Length)
		if _, err := reader.Read(tlvData); err != nil {
			return nil, fmt.Errorf("failed to read TLV data: %w", err)
		}

		tlv := TLV{
			Header: tlvHdr,
			Value:  tlvData,
		}
		tlvs = append(tlvs, tlv)
	}

	frame := &Frame{
		Header: header,
		TLVs:   tlvs,
	}
	return frame, nil
}
