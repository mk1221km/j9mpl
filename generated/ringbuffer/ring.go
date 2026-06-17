package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const ringCapacity = 1024

// ringPush stores value at head, returns updated head and count.
func ringPush(buf []float64, head int, count int, value float64) (int, int) {
	buf[head] = value
	head = (head + 1) % ringCapacity
	if count < ringCapacity {
		count++
	}
	return head, count
}

// ringAvg returns arithmetic mean of stored elements. Returns 0 if empty.
func ringAvg(buf []float64, count int) float64 {
	if count == 0 {
		return 0.0
	}
	var sum float64
	for i := 0; i < count; i++ {
		sum += buf[i]
	}
	return sum / float64(count)
}

// ringReadRange copies up to n of the most recent elements into a new slice.
func ringReadRange(buf []float64, head int, count int, n int) []float64 {
	if n <= 0 || count == 0 {
		return nil
	}
	if n > count {
		n = count
	}
	res := make([]float64, n)
	for i := 0; i < n; i++ {
		idx := (head - n + i + ringCapacity) % ringCapacity
		res[i] = buf[idx]
	}
	return res
}

// ringUtilization returns ratio of stored count to capacity.
func ringUtilization(count int) float64 {
	return float64(count) / float64(ringCapacity)
}

func main() {
	var buf [ringCapacity]float64
	var head, count int

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}
		switch parts[0] {
		case "push":
			if len(parts) < 2 {
				continue
			}
			v, _ := strconv.ParseFloat(parts[1], 64)
			head, count = ringPush(buf[:], head, count, v)
		case "avg":
			fmt.Printf("%.1f\n", ringAvg(buf[:], count))
		case "readRange":
			n := 1
			if len(parts) > 1 {
				n, _ = strconv.Atoi(parts[1])
			}
			vals := ringReadRange(buf[:], head, count, n)
			for _, v := range vals {
				fmt.Printf("%.1f\n", v)
			}
		case "utilization":
			fmt.Printf("%.4f\n", ringUtilization(count))
		}
	}
}
