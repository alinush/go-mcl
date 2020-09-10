package utils

import "fmt"
import "testing"
import "time"
import "github.com/dustin/go-humanize"

func SummarizeResults(size int, multiop string, op string, r *testing.BenchmarkResult) {
	var totalTimeUsecs int64 = r.NsPerOp() / 1000

	fmt.Printf("Average time per %s (%d iters of size %v each): %v\n", multiop, r.N, humanize.Comma(int64(size)), time.Duration(totalTimeUsecs)*time.Microsecond)
	fmt.Printf("Average time per %s: %v\n", op, time.Duration(totalTimeUsecs/int64(size))*time.Microsecond)

	fmt.Println()
}
