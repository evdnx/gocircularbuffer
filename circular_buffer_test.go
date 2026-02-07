package core

import (
	"math"
	"testing"
)

func TestCircularBuffer(t *testing.T) {
	// Test creation with invalid capacity
	_, err := NewCircularBuffer(0)
	if err == nil {
		t.Error("Expected error when creating buffer with capacity 0")
	}

	// Test creation with valid capacity
	buffer, err := NewCircularBuffer(5)
	if err != nil {
		t.Fatalf("Failed to create buffer: %v", err)
	}

	// Test initial state
	if buffer.Size() != 0 {
		t.Errorf("Expected size 0, got %d", buffer.Size())
	}
	if buffer.Capacity() != 5 {
		t.Errorf("Expected capacity 5, got %d", buffer.Capacity())
	}
	if buffer.IsFull() {
		t.Error("New buffer should not be full")
	}

	// Test adding values
	if err := buffer.Add(1.0); err != nil {
		t.Fatalf("Failed to add value: %v", err)
	}
	if err := buffer.Add(2.0); err != nil {
		t.Fatalf("Failed to add value: %v", err)
	}
	if err := buffer.Add(3.0); err != nil {
		t.Fatalf("Failed to add value: %v", err)
	}

	if buffer.Size() != 3 {
		t.Errorf("Expected size 3, got %d", buffer.Size())
	}

	// Test getting values
	val, err := buffer.Get(0)
	if err != nil || val != 1.0 {
		t.Errorf("Expected 1.0 at index 0, got %f", val)
	}

	val, err = buffer.Get(2)
	if err != nil || val != 3.0 {
		t.Errorf("Expected 3.0 at index 2, got %f", val)
	}

	// Test out of bounds
	_, err = buffer.Get(3)
	if err == nil {
		t.Error("Expected error when getting out of bounds index")
	}

	// Test GetLatest and GetOldest
	latest, err := buffer.GetLatest()
	if err != nil || latest != 3.0 {
		t.Errorf("Expected latest value 3.0, got %f", latest)
	}

	oldest, err := buffer.GetOldest()
	if err != nil || oldest != 1.0 {
		t.Errorf("Expected oldest value 1.0, got %f", oldest)
	}

	// Test filling the buffer
	buffer.Add(4.0)
	buffer.Add(5.0)

	if !buffer.IsFull() {
		t.Error("Buffer should be full")
	}

	// Test overwriting
	buffer.Add(6.0)

	if buffer.Size() != 5 {
		t.Errorf("Expected size 5, got %d", buffer.Size())
	}

	oldest, _ = buffer.GetOldest()
	if oldest != 2.0 {
		t.Errorf("Expected oldest value 2.0 after overwrite, got %f", oldest)
	}

	latest, _ = buffer.GetLatest()
	if latest != 6.0 {
		t.Errorf("Expected latest value 6.0, got %f", latest)
	}

	meanAfterOverwrite, err := buffer.Mean()
	if err != nil || meanAfterOverwrite != 4.0 {
		t.Errorf("Expected mean 4.0 after overwrite, got %f", meanAfterOverwrite)
	}

	// Test ToSlice
	slice := buffer.ToSlice()
	expected := []float64{2.0, 3.0, 4.0, 5.0, 6.0}
	if len(slice) != len(expected) {
		t.Errorf("Expected slice length %d, got %d", len(expected), len(slice))
	}
	for i, v := range expected {
		if slice[i] != v {
			t.Errorf("Expected %f at index %d, got %f", v, i, slice[i])
		}
	}

	// Test Clear
	buffer.Clear()
	if buffer.Size() != 0 {
		t.Errorf("Expected size 0 after clear, got %d", buffer.Size())
	}

	// Test Mean, StandardDeviation, Min, Max on empty buffer
	_, err = buffer.Mean()
	if err == nil {
		t.Error("Expected error when calculating mean of empty buffer")
	}

	_, err = buffer.StandardDeviation()
	if err == nil {
		t.Error("Expected error when calculating standard deviation of empty buffer")
	}

	_, err = buffer.Min()
	if err == nil {
		t.Error("Expected error when calculating min of empty buffer")
	}

	_, err = buffer.Max()
	if err == nil {
		t.Error("Expected error when calculating max of empty buffer")
	}

	// Test statistical functions
	buffer.Add(2.0)
	buffer.Add(4.0)
	buffer.Add(6.0)
	buffer.Add(8.0)
	buffer.Add(10.0)

	mean, err := buffer.Mean()
	if err != nil || mean != 6.0 {
		t.Errorf("Expected mean 6.0, got %f", mean)
	}

	stdDev, err := buffer.StandardDeviation()
	if err != nil || stdDev < 3.15 || stdDev > 3.17 { // Approximately 3.16
		t.Errorf("Expected standard deviation ~3.16, got %f", stdDev)
	}

	min, err := buffer.Min()
	if err != nil || min != 2.0 {
		t.Errorf("Expected min 2.0, got %f", min)
	}

	max, err := buffer.Max()
	if err != nil || max != 10.0 {
		t.Errorf("Expected max 10.0, got %f", max)
	}
}

func TestAddMany(t *testing.T) {
	buffer, err := NewCircularBuffer(3)
	if err != nil {
		t.Fatalf("Failed to create buffer: %v", err)
	}

	buffer.AddMany(1, 2, 3, 4)

	values := buffer.ToSlice()
	expected := []float64{2, 3, 4}
	for i, v := range expected {
		if values[i] != v {
			t.Fatalf("Expected %f at index %d, got %f", v, i, values[i])
		}
	}

	mean, err := buffer.Mean()
	if err != nil || mean != 3 {
		t.Fatalf("Expected mean 3, got %f", mean)
	}
}

// TestInvalidValues tests that NaN and Inf are properly rejected
func TestInvalidValues(t *testing.T) {
	buffer, err := NewCircularBuffer(5)
	if err != nil {
		t.Fatalf("Failed to create buffer: %v", err)
	}

	// Test adding NaN
	err = buffer.Add(math.NaN())
	if err == nil {
		t.Error("Expected error when adding NaN")
	}
	if err != ErrInvalidValue {
		t.Errorf("Expected ErrInvalidValue, got %v", err)
	}

	// Test adding positive infinity
	err = buffer.Add(math.Inf(1))
	if err == nil {
		t.Error("Expected error when adding +Inf")
	}
	if err != ErrInvalidValue {
		t.Errorf("Expected ErrInvalidValue, got %v", err)
	}

	// Test adding negative infinity
	err = buffer.Add(math.Inf(-1))
	if err == nil {
		t.Error("Expected error when adding -Inf")
	}
	if err != ErrInvalidValue {
		t.Errorf("Expected ErrInvalidValue, got %v", err)
	}

	// Test AddMany with invalid value
	err = buffer.AddMany(1.0, 2.0, math.NaN(), 4.0)
	if err == nil {
		t.Error("Expected error when adding values including NaN")
	}
	if err != ErrInvalidValue {
		t.Errorf("Expected ErrInvalidValue, got %v", err)
	}

	// Verify buffer is still empty (atomicity - no values added on error)
	if buffer.Size() != 0 {
		t.Errorf("Expected size 0 after failed AddMany, got %d", buffer.Size())
	}

	// Test that valid values work
	err = buffer.AddMany(1.0, 2.0, 3.0)
	if err != nil {
		t.Errorf("Failed to add valid values: %v", err)
	}
	if buffer.Size() != 3 {
		t.Errorf("Expected size 3, got %d", buffer.Size())
	}
}
