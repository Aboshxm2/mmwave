package parser

type FrameHeader struct {
	Version        uint32
	PacketLength   uint32
	Platform       uint32
	FrameNumber    uint32
	CpuCycles      uint32
	NumDetectedObj uint32
	NumTLVs        uint32
	SubFrameNumber uint32
}

type Frame struct {
	Header FrameHeader
	TLVs   []TLV
}

type TLVHeader struct {
	Type   uint32
	Length uint32
}

type TLV struct {
	Header TLVHeader
	Value  []byte
}

type Point struct {
	X       float32
	Y       float32
	Z       float32
	Doppler float32
}

type PointUnit struct {
	ElevationUnit float32
	AzimuthUnit   float32
	DopplerUnit   float32
	RangeUnit     float32
	SnrUnit       float32
}

type CartesianPoint struct {
	Elevation int8
	Azimuth   int8
	Doppler   int16
	Range     int16
	Snr       int16
}

type CompressedPointCloud struct {
	Unit   PointUnit
	Points []CartesianPoint
}

type Target struct {
	ID         uint32
	X          float32
	Y          float32
	Z          float32
	VelX       float32
	VelY       float32
	VelZ       float32
	AccX       float32
	AccY       float32
	AccZ       float32
	Ec         [16]float32
	G          float32
	Confidence float32
}

type TargetHeight struct {
	ID   uint8
	MaxZ float32
	MinZ float32
}
