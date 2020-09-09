#!/bin/bash
set -e

scriptdir=$(cd $(dirname $0); pwd -P)
sourcedir=$(cd $scriptdir/..; pwd -P)
. $scriptdir/shlibs/os.sh

run_bench() {
    curve=$1
    echo
    echo "Running with -curve=$curve"
    echo

    (
        cd $sourcedir/src/mcl && \
        go test -v -bench=BenchmarkG1 -benchtime=4000x -test.run=XXX -curve $curve
        go test -v -bench=BenchmarkG2 -benchtime=2000x -test.run=XXX -curve $curve
        go test -v -bench=BenchmarkPairing -benchtime=1000x -test.run=XXX -curve $curve
    )

    (
        cd $sourcedir/src/app && \
        export GOPATH=$sourcedir
        go run BenchMultiScalarMult.go -size 32 -curve $curve
        go run BenchMultiPairing.go -size 3 -curve $curve
        go run BenchMultiPairing.go -size 10 -curve $curve
    )
}

#
# NOTE: I disabled support for linking to the bn256 and bn384 libraries in src/mcl.go, so some curve benchmarks have been removed.
# Also, no need to pass -tags to 'go test'.
# (See original file here: https://github.com/alinush/mcl/blob/master/bench.sh)
# However, the BN254 curve can still be tested/benchmarked.
#
run_bench bn254
run_bench bn254_snark
run_bench bls12-381
