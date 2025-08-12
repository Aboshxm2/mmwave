# mmWave Library

A simple library for processing mmWave radar demo outputs. 

## Features

- **Serial Communication**: Read raw data frames from mmWave radar devices via UART.
- **Data Parsing**: Parse radar frames and extract TLVs (Type-Length-Value) for further processing.

## Installation
```bash
go get github.com/Aboshxm2/mmwave
```

## Usage
### Basic Example

```go
reader, err := serial.NewUARTReader("COM13", 921600)
if err != nil {
    panic(err)
}

reader.Start()
defer reader.Stop()

for frame := range parser.ParseFrames(reader.OutChan()) {
    fmt.Println("Received frame with TLVs:", len(frame.TLVs))
    for _, tlv := range frame.TLVs {
        if tlv.Header.Type == parser.EXT_TARGET_LIST {
            targets, err := tlv.AsTargetList()
            if err != nil {
                fmt.Println("Error decoding target list:", err)
                continue
            }
            for _, target := range *targets {
                fmt.Printf("Target ID: %d, X: %.2f, Y: %.2f, Confidence: %.2f\n", target.ID, target.X, target.Y, target.Confidence)
            }
        }
    }
}
```

