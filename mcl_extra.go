package mcl

import "fmt"
import log "github.com/sirupsen/logrus"
import "math/big"

var curve int   // the currently selected elliptic curve
var r big.Int   // the order r of the elliptic curve (and the field Fr)
var r1 Fr       // r-1
var br1 big.Int // r-1 as a big.Int
var g Fr        // generator of multiplicative subgroup G_{r-1} of Fr
var primes []Fr // prime factors of r-1
var maxK int    // Fr supports an nth root of unity for n = 2^{maxK}
var omega Fr    // a primitive nth root of unity (for the max n that Fr supports)

// Go's module initialization function
func init() {
	InitFromString("bls12-381")
}

// Returns the prime factors of r-1 where r is the order of the field.
// These are needed to find a generator of the multiplicative subgroup of the filed (of order r-1)
func initPrimeFactors() {
	// TODO: ideally, would've factored the order of Fr rather than manually writing these down
	multiplicity := make(map[int]int)
	switch curve {
	case BLS12_381:
		// r-1 = 2^32 * 3 * 11 * 19 * 10177 * 125527 * 859267 * 906349^2 * 2508409 * 2529403 * 52437899 * 254760293^2 (12 distinct prime factors)
		primes = make([]Fr, 12)
		i := 0
		primes[i].SetInt64(2)
		multiplicity[i] = 32
		i++
		primes[i].SetInt64(3)
		i++
		primes[i].SetInt64(11)
		i++
		primes[i].SetInt64(19)
		i++
		primes[i].SetInt64(10177)
		i++
		primes[i].SetInt64(125527)
		i++
		primes[i].SetInt64(859267)
		i++
		primes[i].SetInt64(906349)
		multiplicity[i] = 2
		i++
		primes[i].SetInt64(2508409)
		i++
		primes[i].SetInt64(2529403)
		i++
		primes[i].SetInt64(52437899)
		i++
		primes[i].SetInt64(254760293)
		multiplicity[i] = 2
		i++
		if i != 12 {
			log.Panicf("BLS12-381 was supposed to have 12 distinct prime factors")
		}
	case CurveSNARK1:
		// 2^28 * 3^2 * 13 * 29 * 983 * 11003 * 237073 * 405928799 * 1670836401704629 * 13818364434197438864469338081 (10 distinct prime factors)
		primes = make([]Fr, 10)
		i := 0
		primes[i].SetInt64(2)
		multiplicity[i] = 28
		i++
		primes[i].SetInt64(3)
		multiplicity[i] = 2
		i++
		primes[i].SetInt64(13)
		i++
		primes[i].SetInt64(29)
		i++
		primes[i].SetInt64(983)
		i++
		primes[i].SetInt64(11003)
		i++
		primes[i].SetInt64(237073)
		i++
		primes[i].SetInt64(405928799)
		i++
		primes[i].SetString("1670836401704629", 10)
		i++
		primes[i].SetString("13818364434197438864469338081", 10)
		i++
		if i != 10 {
			log.Panicf("CurveSNARK1 was supposed to have 12 distinct prime factors")
		}

	default:
		log.Panicf("We do not yet support this curve: %v", curve)
	}

	// Make sure we didn't make a mistake
	acc := big.NewInt(1)
	for i, p := range primes {
		bp := p.ToBigInt()
		pow, ok := multiplicity[i]

		if bp.ProbablyPrime(64) == false {
			log.Panicf("Expected %v to be a prime", p.GetString(10))
		}

		acc.Mul(acc, bp)
		// Some prime factors of r-1 have multiplicity 'pow': i.e., p^pow divides r-1
		if ok {
			for j := 0; j < pow-1; j++ {
				acc.Mul(acc, bp)
			}
		}
	}

	if acc.Cmp(&br1) != 0 {
		log.Panicf("Expected primes to multiply to r-1 = %v, but got %v instead", br1.String(), acc.String())
	} else {
		//log.Infof("Prime factorization of r-1 checks out!")
	}
}

// Call this function before calling all the other operations.
// This function is not thread safe.
// 'curve' can be 'bn254', 'bn254_snark', 'fp382-1', 'fp382-2' or 'bls12-381', but the current library
// is linked with mcl in such a way that only 'bn254', 'bn254_snark' and 'bls12-381' are supported.
func InitFromString(curveStr string) {
	switch curveStr {
	case "":
		// so that when make.sh runs this we just exit
		fmt.Printf("No -curve argument given to 'go test', so assuming -curve=bls12-381...\n")
		curve = BLS12_381

		// Generated using 'go run app/FindRootsOfUnity.go -curve bls12-381'
		g.SetString("43354073823925847805380740654325203111959447648024832395752560736636494814639", 10)
		omega.SetString("29888530911446421229717110002157968030004381923472101974504948698264983864162", 10)
	case "bn254":
		curve = CurveFp254BNb
		// Really, only supports maxK = 2, so for all intents and purposes does not support roots of unity
		maxK = -1
	case "bn254_snark":
		curve = CurveSNARK1
		maxK = 28

		// Generated using 'go run app/FindRootsOfUnity.go -curve bn254_snark'
		g.SetString("2075519347835359061107806824385881338805436780408868160087873480743949887629", 10)
		omega.SetString("6269523429012533460690971800953088755243712115228691635134166863650918497857", 10)
	case "fp382-1":
		curve = CurveFp382_1
	case "fp382-2":
		curve = CurveFp382_2
	case "bls12-381":
		curve = BLS12_381
		maxK = 32
	default:
		log.Panicf("Unknown curve type: %s\n", curveStr)
	}

	// Initialize mcl
	InitMclHelper(curve)

	// The order r of the Fr = \mathbb{F}_r field "in the exponent" of the elliptic curve group
	var str string = GetCurveOrder()
	_, succ := r.SetString(str, 10)
	if !succ {
		log.Panicf("Error parsing big int: %v", str)
	}

	// The order r - 1 of the multiplicative subgroup G_{r-1} of Fr
	br1.Sub(&r, big.NewInt(1))
	r1 = *BigIntToFr(&br1)

	// WARNING: This needs br1 to be set
	if SupportsRootsOfUnity() {
		initPrimeFactors()
	}

	log.Debugf("Order of field \"in the exponent\" F_r:     %v", r.String())
	log.Debugf("Order of multiplicative subgroup G_{r-1}: %v", br1.String())

	log.Debugf("Selected curve: %s (%d)\n", curveStr, curve)
	log.Debugf("GetMaxOpUnitSize() = %d\n", GetMaxOpUnitSize())
	log.Debugf("GetFrUnitSize() = %d\n", GetFrUnitSize())
}

// Converts an Fr to a big.Int
func (x *Fr) ToBigInt() *big.Int {
	var str string = x.GetString(10)
	var b big.Int
	var succ bool
	_, succ = b.SetString(str, 10)
	if !succ {
		log.Panicf("Could not create big.Int from %v", str)
	}
	return &b
}

// Converts a big.Int to an Fr
func BigIntToFr(b *big.Int) *Fr {
	var x Fr
	x.SetString(b.String(), 10)
	return &x
}

// mcl does not provide a function to exponentiate Fr's, which we need when randomly picking generators of G_{r-1}.
// Thus, we implement it here (slowly) via Go's big.Int.
func FrModExp_Slow(base *Fr, exp *Fr) *Fr {
	// convert Fr's to big.Int's
	a := base.ToBigInt()
	e := exp.ToBigInt()

	// Do fast modular exponentiation via big.Int API
	var res big.Int
	res.Exp(a, e, &r)

	// Convert big.Int back to Fr and return it
	return BigIntToFr(&res)
}

// Returns the generator of the multiplicative subgroup of order r-1 of the field Fr
func RandomFieldGenerator() Fr {
	if len(primes) == 0 {
		log.Panicf("Your initialization code was supposed to call initPrimeFactors(). This likely isn't your fault.")
	}

	good := false
	var g Fr

	for !good {
		// Pick a random g in Fr and ensure g^{(r-1)/k} != 1 for all prime factors k
		// If so, then g is a generator of Fr.
		g.Random()

		good = true
		for _, k := range primes {
			// Compute exponent (r-1)/k
			var exp Fr
			FrDiv(&exp, &r1, &k)

			// Compute g^{(r-1)/k} (mod r)
			res := FrModExp_Slow(&g, &exp)

			// Make sure it is not equal to 1
			if res.IsOne() {
				//log.Warnf("Failed on g^{(r-1)/%v}. Trying another generator...", k.GetString(10))
				good = false
				break
			}
		}
	}

	return g
}

func SupportsRootsOfUnity() bool {
	return maxK != -1
}

// Returns a primitive nth root of unity in Fr (assuming it supports one)
// Here, g is a generator of G_{r-1} and r1 is r-1.
func GetRootOfUnity() Fr {
	if !SupportsRootsOfUnity() {
		log.Panicf("The currently selected curve does not support roots of unity")
	}
	return omega
}

func GetRootOfUnityFromGen(gen *Fr, n *Fr) Fr {
	if !SupportsRootsOfUnity() {
		log.Panicf("The currently selected curve does not support roots of unity")
	}

	// TODO: check n.isPowerOfTwo() and n.log2() <= maxK

	// Compute exponent (r-1)/n
	var exp Fr
	FrDiv(&exp, &r1, n)

	return *FrModExp_Slow(gen, &exp)
}
