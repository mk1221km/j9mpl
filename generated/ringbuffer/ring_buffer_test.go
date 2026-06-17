package ringbuffer

import (
	"testing"
)

func FuzzRingBuffer(f *testing.F) {
	// Seed corpus: baseline valid operations
	f.Add(0.0, 1)
	f.Add(-1.0, 5)
	f.Add(100.0, 0)
	f.Add(1e10, 2048) // exceed capacity (1024)

	f.Fuzz(func(t *testing.T, value float64, count int) {
		rb := &RingBuffer{}
		rb.Push(value)
		rb.Push(value * 2)
		rb.Push(value * 3)

		result := rb.ReadRange(count)
		_ = result

		avg := rb.Avg()
		_ = avg

		util := rb.Utilization()
		_ = util
	})
}
