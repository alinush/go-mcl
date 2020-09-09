# app/ folder

This folder has an application for finding a primitive $n$th root of unity in the finite field $\mathbb{F}_p$ associated with the elliptic curve of order $p$.
It also has a few benchmark applications.

Please see details below.

## Finding roots of unity

`FindRootsOfUnity.go` can be used to:

 1. Pick a random generator $g$ for the multiplicative subgroup $G_{r-1}$ of the finite field $\mathbb{F}_r$ associated with the elliptic curve group.
 2. Pick a random $N$th primitive root of unity where $N=2^k$ is the maximum supported in $\mathbb{F}_r$.

Specifically, pick a supported curve (e.g., `bls12-381`) and do:

    . ../../scripts/set-env.sh
    go run FindRootsOfUnity.go -curve=bls12-381

This will output some Go code that sets the generator and the primitive root of unity.
You can then copy & paste this code into the right place in `src/mcl/mcl_extra.go`

## mcl multi-exp and multi-pairing benchmarks

To run individual benchmarks for multi-exps and multi-pairings:

    . ../../scripts/set-env.sh

    go run BenchMultiScalarMult.go -size 1000
    go run BenchMultiPairing.go -size 1000

## Gurvy multiexp

Gurvy has a very fast multiexp implementation, but it uses multi-threading.
The gurvy multithreaded numbers are much faster on a 16 core machine than single-threaded mcl (e.g., 3 mus per G1 exp in Gurvy multiexp of 1024 versus 33 mus in mcl).

    go run BenchGurvyMultiScalarMult.go -size 1024

Without multi-threading, gurvy is still 2x faster than mcl (15 mus versus 33 mus per G1 exp in a multiexp of size 2048).

    GOMAXPROCS=1 go run BenchGurvyMultiScalarMult.go -size=2048

Also, single-threaded gurvy scales much better with larger multiexp size than mcl, which always does 33 mus.
