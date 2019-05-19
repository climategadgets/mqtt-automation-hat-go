package filter

type Filter interface {
	// Consume takes an input value, and returns a filtered value.
	Consume(sample float64) float64
}
