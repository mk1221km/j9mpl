package ringbuffer



type RingBuffer struct {
	Buf [1024]float64
	Head int
	Count int
}

func (s *RingBuffer) NewRingBuffer() *RingBuffer {
    return &RingBuffer{}
}

func (s *RingBuffer) Push(value float64) {
    s.Buf[s.Head] = value
    s.Head = (s.Head + 1) % len(s.Buf)
    if s.Count < len(s.Buf) {
        s.Count++
    }
}

func (s *RingBuffer) ReadRange(count int) []float64 {
    if count <= 0 || s.Count == 0 {
        return nil
    }
    if count > s.Count {
        count = s.Count
    }
    cap := len(s.Buf)
    startIdx := (s.Head - s.Count) % cap
    if startIdx < 0 {
        startIdx += cap
    }
    offset := s.Count - count
    actualStart := (startIdx + offset) % cap
    result := make([]float64, count)
    for i := 0; i < count; i++ {
        result[i] = s.Buf[(actualStart+i)%cap]
    }
    return result
}

func (s *RingBuffer) Avg() float64 {
    if s.Count == 0 {
        return 0
    }
    var sum float64
    size := len(s.Buf)
    for i := 0; i < s.Count; i++ {
        idx := (s.Head - s.Count + i + size) % size
        sum += s.Buf[idx]
    }
    return sum / float64(s.Count)
}

func (s *RingBuffer) Utilization() float64 {
    return float64(s.Count) / float64(len(s.Buf))
}
