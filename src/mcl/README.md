mcl's Golang bindings
=====================

The `mcl.go` file was merged from two files in [herumi/mcl](https://github.com/herumi/mcl):

    1. [ffi/go/mcl/init.go](https://github.com/herumi/mcl/blob/master/ffi/go/mcl/init.go)
        - This file just had an `Init()` function, so any future changes to it are easy to integrate.
    1. [ffi/go/mcl/mcl.go](https://github.com/herumi/mcl/blob/master/ffi/go/mcl/mcl.go)
        - For this file, a manual diff will need to be done between herumi's version and our version to integrate changes

So to integrate new changes, do a:

    git clone https://github.com/herumi/mcl
    diff -rupN mcl/ffi/go/mcl/init.go init.go 
    diff -rupN mcl/ffi/go/mcl/mcl.go  mcl.go  

Changes in the `#cgo` directives in both these files will have to be manually integrated.
