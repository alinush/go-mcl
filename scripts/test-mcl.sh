#!/bin/bash
set -e

scriptdir=$(cd $(dirname $0); pwd -P)
sourcedir=$(cd $scriptdir/..; pwd -P)
. $scriptdir/shlibs/os.sh

cd $sourcedir/src/mcl/
go test -v -curve=bls12-381
#go test -v -curve=bn254 # NOTE: Disabled due to lack of support for roots of unity
go test -v -curve=bn254_snark
