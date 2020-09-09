package main

import "flag"
import "fmt"
import log "github.com/sirupsen/logrus"
import "mcl"
import "utils"

var _curveArg = utils.GetCurveArgument()
var _kArg = flag.Int("k", 32, "Finds a (2^k)th primitive root of unity")

func main() {
	log.SetLevel(log.DebugLevel)
	flag.Parse()

	k := *_kArg
	curveStr := *_curveArg

	// Make sure this is a curve we support (i.e., we've factorized r-1 so we can test the generator)
	switch curveStr {
	case "bls12-381":
	case "bn254_snark":
	default:
		log.Panicf("We do not yet support this curve: %s.", curveStr)
	}

	// Initialize mcl
	mcl.InitFromString(curveStr)

	// Find a generator for the multiplicative subgroup of the field
	g := mcl.RandomFieldGenerator()

	var n mcl.Fr
	// Compute n = 2^k
	mcl.FrPow2(&n, k)
	log.Printf("k = %v, n = %v", k, n.GetString(10))
	// Find an nth root of unity
	omega := mcl.GetRootOfUnityFromGen(&g, &n)

	fmt.Println()
	fmt.Println("Just copy the code below into InitFromString() in mcl_extra.go:")
	fmt.Println()
	fmt.Printf("// Generated using 'go run app/FindRootsOfUnity.go -curve %v'\n", curveStr)
	fmt.Printf("g.SetString(\"%s\", %v)\n", g.GetString(10), 10)
	fmt.Printf("omega.SetString(\"%s\", %v)\n", omega.GetString(10), 10)
	fmt.Println()
}
