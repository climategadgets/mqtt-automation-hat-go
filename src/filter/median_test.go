package filter

import (
	"fmt"
	"testing"
)

// TestMedian1 is an edge case - this is not a filter but a pipe, but who are we to judge the client's choice
func TestMedian1(t *testing.T) {

	source := []float64{0.0, 1.0, 2.0, 3.0, 4.0, 5.0}
	result := []float64{0.0, 1.0, 2.0, 3.0, 4.0, 5.0}

	testMedian(t, 1, source, result)
}

// TestMedian3 tests the median filter of size 3
func TestMedian3(t *testing.T) {

	source := []float64{0.0, 1.0, 2.0, 3.0, 4.0, 5.0}
	result := []float64{0.0, 1.0, 1.0, 2.0, 3.0, 4.0}

	testMedian(t, 3, source, result)
}

// TestMedian5 tests the median filter of size 5
func TestMedian5(t *testing.T) {

	source := []float64{0.0, 1.0, 2.0, 3.0, 4.0, 5.0}
	result := []float64{0.0, 1.0, 2.0, 3.0, 2.0, 3.0}

	// VT: NOTE: actual buffer size will be 5
	testMedian(t, 4, source, result)
}

func testMedian(t *testing.T, size uint8, source []float64, expected []float64) {

	filter := NewMedianFilter(size)

	for offset, value := range source {
		if filtered := filter.Consume(value); filtered != expected[offset] {
			t.Fatalf(fmt.Sprintf("@%v: expected %v, got %v", offset, expected[offset], filtered))
		}
	}
}
