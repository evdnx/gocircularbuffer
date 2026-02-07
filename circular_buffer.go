package core

import (
	"errors"
	"math"
)

var (
	// ErrInvalidCapacity indicates the buffer was created with an invalid capacity.
	ErrInvalidCapacity = errors.New("circularbuffer: capacity must be greater than 0")
	// ErrIndexOutOfBounds indicates the caller asked for an index that does not exist.
	ErrIndexOutOfBounds = errors.New("circularbuffer: index out of bounds")
	// ErrBufferEmpty indicates the buffer has no items.
	ErrBufferEmpty = errors.New("circularbuffer: buffer is empty")
	// ErrInsufficientData indicates more items are required for the operation.
	ErrInsufficientData = errors.New("circularbuffer: need at least two values to calculate standard deviation")
	// ErrInvalidValue indicates the caller tried to add an invalid value (NaN or Inf).
	ErrInvalidValue = errors.New("circularbuffer: value must be a finite number (not NaN or Inf)")
)

// CircularBuffer implements a fixed-size buffer that overwrites oldest data when full
type CircularBuffer struct {
	data       []float64
	capacity   int
	size       int
	head       int // Points to the next position to write
	sum        float64
	sumSquares float64
}

// NewCircularBuffer creates a new circular buffer with the specified capacity
func NewCircularBuffer(capacity int) (*CircularBuffer, error) {
	if capacity <= 0 {
		return nil, ErrInvalidCapacity
	}

	return &CircularBuffer{
		data:     make([]float64, capacity),
		capacity: capacity,
		size:     0,
		head:     0,
	}, nil
}

// Add adds a value to the buffer, overwriting the oldest value if the buffer is full.
// Returns an error if the value is NaN or Inf, which would corrupt statistical calculations.
func (cb *CircularBuffer) Add(value float64) error {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return ErrInvalidValue
	}

	if cb.size == cb.capacity {
		overwritten := cb.data[cb.head]
		cb.sum -= overwritten
		cb.sumSquares -= overwritten * overwritten
	} else {
		cb.size++
	}

	cb.data[cb.head] = value
	cb.sum += value
	cb.sumSquares += value * value
	cb.head = (cb.head + 1) % cb.capacity
	return nil
}

// AddMany appends several values to the buffer in order.
// Returns an error if any value is NaN or Inf. If an error occurs, no values are added.
func (cb *CircularBuffer) AddMany(values ...float64) error {
	// Validate all values first to ensure atomicity
	for _, value := range values {
		if math.IsNaN(value) || math.IsInf(value, 0) {
			return ErrInvalidValue
		}
	}
	
	// Add all values
	for _, value := range values {
		// We already validated, so this should never error
		_ = cb.Add(value)
	}
	return nil
}

// Get returns the value at the specified index (0 is the oldest, size-1 is the newest)
func (cb *CircularBuffer) Get(index int) (float64, error) {
	if index < 0 || index >= cb.size {
		return 0, ErrIndexOutOfBounds
	}

	// Calculate the actual index in the underlying array
	actualIndex := (cb.head - cb.size + index + cb.capacity) % cb.capacity
	return cb.data[actualIndex], nil
}

// GetLatest returns the most recently added value
func (cb *CircularBuffer) GetLatest() (float64, error) {
	if cb.size == 0 {
		return 0, ErrBufferEmpty
	}

	latestIndex := (cb.head - 1 + cb.capacity) % cb.capacity
	return cb.data[latestIndex], nil
}

// GetOldest returns the oldest value in the buffer
func (cb *CircularBuffer) GetOldest() (float64, error) {
	if cb.size == 0 {
		return 0, ErrBufferEmpty
	}

	oldestIndex := (cb.head - cb.size + cb.capacity) % cb.capacity
	return cb.data[oldestIndex], nil
}

// Size returns the current number of elements in the buffer
func (cb *CircularBuffer) Size() int {
	return cb.size
}

// Capacity returns the maximum capacity of the buffer
func (cb *CircularBuffer) Capacity() int {
	return cb.capacity
}

// IsFull returns true if the buffer is at full capacity
func (cb *CircularBuffer) IsFull() bool {
	return cb.size == cb.capacity
}

// Clear empties the buffer
func (cb *CircularBuffer) Clear() {
	for i := range cb.data {
		cb.data[i] = 0
	}
	cb.size = 0
	cb.head = 0
	cb.sum = 0
	cb.sumSquares = 0
}

// ToSlice returns all values in the buffer as a slice, from oldest to newest
func (cb *CircularBuffer) ToSlice() []float64 {
	if cb.size == 0 {
		return []float64{}
	}

	result := make([]float64, cb.size)
	for i := 0; i < cb.size; i++ {
		actualIndex := (cb.head - cb.size + i + cb.capacity) % cb.capacity
		result[i] = cb.data[actualIndex]
	}
	return result
}

// Mean calculates the mean of all values in the buffer
func (cb *CircularBuffer) Mean() (float64, error) {
	if cb.size == 0 {
		return 0, ErrBufferEmpty
	}

	return cb.sum / float64(cb.size), nil
}

// StandardDeviation calculates the standard deviation of all values in the buffer
func (cb *CircularBuffer) StandardDeviation() (float64, error) {
	if cb.size <= 1 {
		return 0, ErrInsufficientData
	}

	mean := cb.sum / float64(cb.size)
	meanOfSquares := cb.sumSquares / float64(cb.size)
	variance := (meanOfSquares - (mean * mean)) * float64(cb.size) / float64(cb.size-1)
	if variance < 0 {
		variance = 0 // Guard against negative zero due to floating point precision
	}
	return math.Sqrt(variance), nil
}

// Min returns the minimum value in the buffer
func (cb *CircularBuffer) Min() (float64, error) {
	if cb.size == 0 {
		return 0, ErrBufferEmpty
	}

	min := math.MaxFloat64
	for i := 0; i < cb.size; i++ {
		actualIndex := (cb.head - cb.size + i + cb.capacity) % cb.capacity
		if cb.data[actualIndex] < min {
			min = cb.data[actualIndex]
		}
	}
	return min, nil
}

// Max returns the maximum value in the buffer
func (cb *CircularBuffer) Max() (float64, error) {
	if cb.size == 0 {
		return 0, ErrBufferEmpty
	}

	max := -math.MaxFloat64
	for i := 0; i < cb.size; i++ {
		actualIndex := (cb.head - cb.size + i + cb.capacity) % cb.capacity
		if cb.data[actualIndex] > max {
			max = cb.data[actualIndex]
		}
	}
	return max, nil
}
