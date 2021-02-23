package main

import (
	"flag"
	"fmt"

	// the dot (.) means import everything into the global namespace, since I don't want to have to say mcl.G1 and mcl.G2
	"testing"

	. "github.com/alinush/go-mcl"
	"github.com/alinush/go-mcl/utils"
)

var _curveArg = utils.GetCurveArgument()
var _sizeArg = flag.Int("size", 2048, "The number of pairings to compute: \\prod_{i \\in [n]} e(P_i, Q_i)")

func main() {
	flag.Parse()

	InitFromString(*_curveArg)
	fmt.Printf("size = %v\n", *_sizeArg)

	// See https://golang.org/pkg/testing/#BenchmarkResult
	results := testing.Benchmark(BenchmarkMultiPairing)
	utils.SummarizeResults(*_sizeArg, "multi-pairing", "pairing", &results)
}

func BenchmarkMultiPairing(b *testing.B) {
	var size int = *_sizeArg

	g1Vec := make([]G1, size)
	g2Vec := make([]G2, size)

	for i := 0; i < b.N; i++ {
		//fmt.Printf("Iteration #%d out of %d\n", i, b.N)
		b.StopTimer()

		var r1, r2 GT

		// pick a random G1/G2
		g1Vec[0].Random()
		g2Vec[0].Random()

		r1.SetInt64(1)
		for i := 0; i < size; i++ {
			// the other random G1/G2s are just doublings of the first one
			if i > 0 {
				G1Dbl(&g1Vec[i], &g1Vec[i-1])
				G2Dbl(&g2Vec[i], &g2Vec[i-1])
			}

			// compute the multipairing inefficiently
			Pairing(&r2, &g1Vec[i], &g2Vec[i])
			GTMul(&r1, &r1, &r2)
		}

		b.StartTimer()

		// compute the multipairing efficiently
		MillerLoopVec(&r2, g1Vec, g2Vec)
		FinalExp(&r2, &r2)

		if !r1.IsEqual(&r2) {
			panic(fmt.Sprintf("Wrong MultiPairing! Efficient returned %v, but naive returned %v", r2, r1))
		}
	}
}
