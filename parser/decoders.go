package parser

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func (t *TLV) AsPointCloud() (*[]Point, error) {
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

	return &points, nil
}

func (t *TLV) AsCompressedPointCloud() (*CompressedPointCloud, error) {
	if t.Header.Type != COMPRESSED_POINTS {
		return nil, fmt.Errorf("TLV type %d is not CompressedPoints", t.Header.Type)
	}

	var cp CompressedPointCloud
	reader := bytes.NewReader(t.Value)

	if err := binary.Read(reader, binary.LittleEndian, &cp.Unit); err != nil {
		return nil, fmt.Errorf("failed to read PointUnit: %w", err)
	}

	for reader.Len() > 0 {
		var point CartesianPoint
		if err := binary.Read(reader, binary.LittleEndian, &point); err != nil {
			return nil, fmt.Errorf("failed to read CartesianPoint: %w", err)
		}
		cp.Points = append(cp.Points, point)
	}

	return &cp, nil
}

func (t *TLV) AsTargetList() (*[]Target, error) {
	// if t.Header.Type != EXT_TARGET_LIST {
	// 	return nil, fmt.Errorf("TLV type %d is not TargetList", t.Header.Type)
	// }

	var targets []Target
	reader := bytes.NewReader(t.Value)

	for reader.Len() > 0 {
		var target Target
		if err := binary.Read(reader, binary.LittleEndian, &target); err != nil {
			return nil, fmt.Errorf("failed to read Target: %w", err)
		}
		targets = append(targets, target)
	}

	return &targets, nil
}

func (t *TLV) AsTargetIndex() ([]uint8, error) {
	// if t.Header.Type != EXT_TARGET_INDEX {
	// 	return nil, fmt.Errorf("TLV type %d is not TargetIndex", t.Header.Type)
	// } // TODO

	var indices []uint8
	reader := bytes.NewReader(t.Value)

	for reader.Len() > 0 {
		var index uint8
		if err := binary.Read(reader, binary.LittleEndian, &index); err != nil {
			return nil, fmt.Errorf("failed to read TargetIndex: %w", err)
		}
		indices = append(indices, index)
	}

	return indices, nil
}

func (t *TLV) AsPresenceIndecation() ([]uint8, error) {
	if t.Header.Type != PRESCENCE_INDICATION {
		return nil, fmt.Errorf("TLV type %d is not PresenceIndication", t.Header.Type)
	}

	var presence []uint8
	reader := bytes.NewReader(t.Value)

	for reader.Len() > 0 {
		var index uint8
		if err := binary.Read(reader, binary.LittleEndian, &index); err != nil {
			return nil, fmt.Errorf("failed to read PresenceIndication: %w", err)
		}
		presence = append(presence, index)
	}

	return presence, nil
}

func (t *TLV) AsTargetHeight() (*[]TargetHeight, error) {
	if t.Header.Type != TRACKERPROC_TARGET_HEIGHT {
		return nil, fmt.Errorf("TLV type %d is not TargetHeight", t.Header.Type)
	}

	var heights []TargetHeight
	reader := bytes.NewReader(t.Value)

	for reader.Len() > 0 {
		var height TargetHeight
		if err := binary.Read(reader, binary.LittleEndian, &height); err != nil {
			return nil, fmt.Errorf("failed to read TargetHeight: %w", err)
		}
		heights = append(heights, height)
	}

	return &heights, nil
}
