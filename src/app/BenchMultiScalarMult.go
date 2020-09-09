package main

import "flag"
import "fmt"
import . "mcl"
import "testing"
import "utils"

var _curveArg = utils.GetCurveArgument()
var _sizeArg = flag.Int("size", 32*32*32, "The size of the multi-scalar multiplication (MSM) benchmark")

func main() {
	flag.Parse()

	InitFromString(*_curveArg)
	fmt.Printf("size = %v\n", *_sizeArg)

	// See https://golang.org/pkg/testing/#BenchmarkResult
	results := testing.Benchmark(BenchmarkG1MultiScalarMult)
	utils.SummarizeResults(*_sizeArg, "G1 multiexp", "G1 exp", &results)

	results = testing.Benchmark(BenchmarkG2MultiScalarMult)
	utils.SummarizeResults(*_sizeArg, "G2 multiexp", "G2 exp", &results)
}

func BenchmarkG1MultiScalarMult(b *testing.B) {
	var size int = *_sizeArg

	gVec := make([]G1, size)
	cVec := make([]Fr, size)

	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		//fmt.Printf("Iteration #%d out of %d\n", j, b.N)
		b.StopTimer()

		var r1, r2 G1
		// pick a random G1
		gVec[0].Random()

		for i := 0; i < size; i++ {
			// pick random Fr's
			cVec[i].Random()

			// the other random G1s are just doublings of the first one
			if i > 0 {
				G1Dbl(&gVec[i], &gVec[i-1])
			}

			// compute the MSM inefficiently
			G1Mul(&r1, &gVec[i], &cVec[i])
			G1Add(&r2, &r2, &r1)
		}

		b.StartTimer()
		// compute the MSM efficiently
		G1MulVec(&r1, gVec, cVec)

		if !r1.IsEqual(&r2) {
			panic(fmt.Sprintf("Wrong G1MulVec! MSM returned %v, but expected %v", r1, r2))
		}
	}
}

func BenchmarkG2MultiScalarMult(b *testing.B) {
	var size int = *_sizeArg

	gVec := make([]G2, size)
	cVec := make([]Fr, size)

	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		//fmt.Printf("Iteration #%d out of %d\n", j, b.N)
		b.StopTimer()

		var r1, r2 G2
		// pick a random G2
		gVec[0].Random()

		for i := 0; i < size; i++ {
			// pick random Fr's
			cVec[i].Random()

			// the other random G2s are just doublings of the first one
			if i > 0 {
				G2Dbl(&gVec[i], &gVec[i-1])
			}

			// compute the MSM inefficiently
			G2Mul(&r1, &gVec[i], &cVec[i])
			G2Add(&r2, &r2, &r1)
		}

		b.StartTimer()
		// compute the MSM efficiently
		G2MulVec(&r1, gVec, cVec)

		if !r1.IsEqual(&r2) {
			panic(fmt.Sprintf("Wrong G2MulVec! MSM returned %v, but expected %v", r1, r2))
		}
	}
}
