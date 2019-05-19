package filter

import (
	"sort"
)

// Median filter support data structure.
// The emphasis is on extremely long uptime, so initial memory investment into another buffer
// is negligible in comparison to effort reduction down the road.
type median struct {
	input   []float64 // Ring buffer for input samples
	offset  uint8     // Next element offset within the ring buffer
	sorted  []float64 // Sort scratchpad, to reduce memory chatter
	unknown uint8     // Number of unknown values remaining
	median  uint8     // Median offset
}

// NewMedianFilter creates a new filter instance with the buffer size being the next odd number up from the one given.
func NewMedianFilter(size uint8) Filter {

	if (size % 2) == 0 {
		// Enlarge the buffer to the next odd size if the number is even
		size++
	}

	return &median{
		input: make([]float64, size), offset: 0,
		sorted: make([]float64, size),
		unknown: size, median: uint8((size - 1) / 2)}
}

// Consume applies the median filter as described in https://en.wikipedia.org/wiki/Median_filter
func (f *median) Consume(sample float64) float64 {

	// Replace the current element and advance
	f.input[f.offset] = sample
	// Reset the offset if it wraps around
	f.offset = (f.offset + 1) % uint8(len(f.input))

	if f.unknown > 1 {
		f.unknown--

		// For simplicity, let's just return the input value if the buffer is not yet completely filled
		return sample
	}

	copy(f.sorted, f.input)
	sort.Float64s(f.sorted)

	return f.sorted[f.median]
}
