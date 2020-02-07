#!/bin/bash

set -e -o pipefail
trap '[ "$?" -eq 0 ] || echo "Error Line:<$LINENO> Error Function:<${FUNCNAME}>"' EXIT
cd `dirname $0`
CURRENT=`pwd`

function bench
{
  testenv
  go test -v $(go list ./... | grep -v vendor) -run none -bench . -benchtime 10s -benchmem
}

function test
{
  testenv
  go test -v $(go list ./... | grep -v vendor) --count 1 -race -covermode=atomic -timeout 120s
}

function testenv
{
     if [ -e $CURRENT/local_env.sh ]; then
         source $CURRENT/local_env.sh
     fi
}

CMD=$1
shift
$CMD $*
