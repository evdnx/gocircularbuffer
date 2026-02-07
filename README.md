# gocircularbuffer

Lightweight, dependency-free circular buffer focused on numerical data. It keeps track of running aggregates so you can compute rolling statistics without repeatedly walking the buffer.

## Features
- Fixed-size circular buffer that overwrites the oldest item when full
- Constant-time `Mean` and `StandardDeviation` via tracked sum and sum of squares
- Convenience helpers like `AddMany`, `ToSlice`, `GetOldest`, `GetLatest`, `Min`, `Max`, and `Clear`
- Exported error values (`ErrInvalidCapacity`, `ErrIndexOutOfBounds`, `ErrBufferEmpty`, `ErrInsufficientData`, `ErrInvalidValue`) for precise error handling
- Input validation: rejects `NaN` and `Inf` values to prevent statistical corruption
- Fully covered by unit tests

## Important Notes

### Input Validation
The buffer validates all input values to ensure data integrity:
- `NaN` (Not a Number) values are rejected
- `Inf` (positive or negative infinity) values are rejected
- These restrictions prevent corruption of statistical calculations (mean, standard deviation)

### Thread Safety
**This buffer is NOT thread-safe**. If you need concurrent access, you must provide your own synchronization (e.g., using `sync.Mutex`).

## Installation

```bash
go get github.com/evdnx/gocircularbuffer
```

## Usage

```go
package main

import (
	"fmt"
	"log"

	core "github.com/evdnx/gocircularbuffer"
)

func main() {
	buffer, err := core.NewCircularBuffer(5)
	if err != nil {
		panic(err)
	}

	err = buffer.AddMany(1, 2, 3, 4, 5, 6)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(buffer.ToSlice()) // [2 3 4 5 6]
	mean, _ := buffer.Mean()
	std, _ := buffer.StandardDeviation()
	max, _ := buffer.Max()

	fmt.Printf("mean=%.2f stddev=%.2f max=%.0f\n", mean, std, max)

	buffer.Clear()
}
```

## Testing

```bash
go test ./...
```

## License

This project is released under the MIT-0 License. See `LICENSE` for details.
