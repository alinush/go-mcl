#!/bin/bash
set -e
    
scriptdir=$(cd $(dirname $0); pwd -P)
sourcedir=$(cd $scriptdir/..; pwd -P)

echo
echo "You can access the docs at:"
echo "     http://localhost:6060/pkg/#thirdparty "
echo

(
    cd $sourcedir/
    GOPATH=$sourcedir godoc -v $@
)
