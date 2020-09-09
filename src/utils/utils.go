package utils

import "flag"

//import "fmt"
//import "time"
//import log "github.com/sirupsen/logrus"

func GetCurveArgument() *string {
	return flag.String("curve", "bls12-381", "Type of supported curve: bn254, bn254_snark, bls12-381")
}

// just use time.Duration::String(): i.e., d.String()
//func HumanizeDuration(d time.Duration) string {
//	units := []string{"mus", "ms", "secs", "mins", "hrs", "days", "years"}
//	numUnits := len(units)
//	result := float64(d.Microseconds())
//	i := 0
//
//	for result >= 1000.0 && i < 2 {
//		result /= 1000.0
//		i++
//	}
//
//	for result >= 60.0 && i >= 2 && i < 4 {
//		result /= 60.0
//		i++
//	}
//
//	if i == 4 && result >= 24.0 {
//		result /= 24.0
//		i++
//	}
//
//	if i == 5 && result >= 365.25 {
//		result /= 365.25
//		i++
//	}
//
//	if i >= numUnits {
//		log.Panicf("Expected i < numUnits, but %d >= %d", i, numUnits)
//	}
//
//	return fmt.Sprintf("%.2f %s", result, units[i])
//}
