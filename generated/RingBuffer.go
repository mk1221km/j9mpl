package ringbuffer



type RingBuffer struct{}

method Push(value float64) public static
{
    if s.Count == len(s.Data) {
        s.Data[s.Head] = value
        s.Head = (s.Head + 1) % len(s.Data)
    } else {
        tail := (s.Head + s.Count) % len(s.Data)
        s.Data[tail] = value
        s.Count++
    }
}
