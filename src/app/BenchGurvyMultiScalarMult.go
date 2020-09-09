package main

import "flag"
import "github.com/consensys/gurvy/bls381"
import "github.com/consensys/gurvy/bls381/fr"
import log "github.com/sirupsen/logrus"
import "runtime"
import "testing"
import "utils"

var _curveArg = utils.GetCurveArgument()
var _sizeArg = flag.Int("size", 32*32, "The size of the multi-scalar multiplication (MSM) benchmark")

var curve *bls381.Curve

func main() {
	flag.Parse()

	switch *_curveArg {
	case "bls12-381":
		curve = bls381.BLS381()
	case "bn254":
		log.Panicf("github.com/consensys/gurvy does not support easily switching between curves AFAICT")
	}

	log.Printf("size = %v\n", *_sizeArg)
	log.Printf("GOMAXPROCS = %v\n", runtime.GOMAXPROCS(-1))

	// See https://golang.org/pkg/testing/#BenchmarkResult
	results := testing.Benchmark(BenchmarkG1MultiScalarMult)
	utils.SummarizeResults(*_sizeArg, "G1 multiexp", "G1 exp", &results)
}

func getRandomG1Vec(n int) []bls381.G1Affine {
	var r bls381.G1Jac
	g := make([]bls381.G1Affine, n)

	for i := 0; i < n; i++ {
		// set g[i] to g^{r_i} for some random r_i
		var scalar fr.Element
		scalar.SetRandom()
		r.ScalarMulByGen(curve, scalar)

		// TODO: Apparently, this might be slow, so should time it (Indeed, 138 mus per element. Crazy.)
		// What would an application built on top of Gurvy would use in-memory: G1Affine or G1Jac?
		// Seems like ScalarMul only works with G1Jac, so that settles it.
		r.ToAffineFromJac(&g[i])
	}

	return g
}

func BenchmarkG1MultiScalarMult(b *testing.B) {
	var size int = *_sizeArg

	e := make([]fr.Element, size)
	g := getRandomG1Vec(size)

	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		//log.Printf("Iteration #%d out of %d, multiexp size %v\n", j, b.N, size)
		b.StopTimer()

		var r1 bls381.G1Jac
		//var r2 bls381.G1Jac

		for i := 0; i < size; i++ {
			// pick random Fr's
			e[i].SetRandom()

			//// compute the MSM inefficiently
			//var gij bls381.G1Jac
			//g[i].ToJacobian(&gij)

			//// compute g_i^{e[i]}
			//gij.ScalarMul(curve, &gij, e[i])

			//// accumulate r = r * g_i^{e[i]}
			//r2.Add(curve, &gij)
		}

		b.StartTimer()

		// compute the MSM efficiently
		<-r1.MultiExp(curve, g, e)

		//if !r1.Equal(&r2) {
		//	log.Panicf("Wrong MSM implementation! MSM returned %v, but expected %v", r1, r2)
		//} else {
		//	//log.Printf("r1: %v", r1)
		//	//log.Printf("r2: %v\n", r2)
		//}
	}
}
