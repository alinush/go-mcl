go-mcl
======

This is a wrapper for `herumi/mcl`'s elliptic curve library that you can import easily into your projects via a simple:

```
    import "github.com/alinush/go-mcl"
```

## Limitations

I disabled support for linking to the `bn256` and `bn384` libraries in `src/mcl.go`, so some curve benchmarks have been removed.
However, `bn256` and its benchmarks still work.

You can read more about our changes in [this README](src/mcl/README.md).

## Documentation

You can see the documentation by running the doc server:

    ./scripts/run-doc-server.sh

Then, just browse to [http://localhost:6060/pkg/#thirdparty](http://localhost:6060/pkg/#thirdparty).

## Dependencies

    ./scripts/install-deps.sh

## Tests and benchmarks

First, set up the Go environment:

    source scripts/set-env.sh

To run all tests:
    
    ./scripts/test-mcl.sh

To run all benchmarks:

    ./scripts/bench-mcl.sh

We include instructions for running benchmarks for a different Golang-based elliptic curve library called [Gurvy](https://github.com/consensys/gurvy) in [this README](src/app/README.md).
