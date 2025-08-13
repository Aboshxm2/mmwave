package parser

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func (t *TLV) AsPointCloud() ([]Point, error) {
	if t.Header.Type != DETECTED_POINTS {
		return nil, fmt.Errorf("TLV type %d is not DetectedPoints", t.Header.Type)
	}

	var points []Point
	reader := bytes.NewReader(t.Value)

	for reader.Len() > 0 {
		var point Point
		if err := binary.Read(reader, binary.LittleEndian, &point); err != nil {
			return nil, fmt.Errorf("failed to read DetectedPoint: %w", err)
		}
	}

	return points, nil
}

func (t *TLV) AsCompressedPointCloud() (cp CompressedPointCloud, err error) {
	if t.Header.Type != COMPRESSED_POINTS {
		err = fmt.Errorf("TLV type %d is not CompressedPoints", t.Header.Type)
		return
	}

	reader := bytes.NewReader(t.Value)

	if err = binary.Read(reader, binary.LittleEndian, &cp.Unit); err != nil {
		err = fmt.Errorf("failed to read PointUnit: %w", err)
		return
	}

	for reader.Len() > 0 {
		var point CartesianPoint
		if err = binary.Read(reader, binary.LittleEndian, &point); err != nil {
			err = fmt.Errorf("failed to read CartesianPoint: %w", err)
			return
		}
		cp.Points = append(cp.Points, point)
	}

	return cp, nil
}

func (t *TLV) AsTargetList() ([]Target, error) {
	var targets []Target
	reader := bytes.NewReader(t.Value)

	for reader.Len() > 0 {
		var target Target
		if err := binary.Read(reader, binary.LittleEndian, &target); err != nil {
			return nil, fmt.Errorf("failed to read Target: %w", err)
		}
		targets = append(targets, target)
	}

	return targets, nil
}

func (t *TLV) AsUint8Slice() ([]uint8, error) {
	var slice []uint8
	reader := bytes.NewReader(t.Value)

	for reader.Len() > 0 {
		var value uint8
		if err := binary.Read(reader, binary.LittleEndian, &value); err != nil {
			return nil, fmt.Errorf("failed to read uint8 value: %w", err)
		}
		slice = append(slice, value)
	}

	return slice, nil
}

func (t *TLV) AsTargetHeight() ([]TargetHeight, error) {
	var heights []TargetHeight
	reader := bytes.NewReader(t.Value)

	for reader.Len() > 0 {
		var height TargetHeight
		if err := binary.Read(reader, binary.LittleEndian, &height); err != nil {
			return nil, fmt.Errorf("failed to read TargetHeight: %w", err)
		}
		heights = append(heights, height)
	}

	return heights, nil
}
