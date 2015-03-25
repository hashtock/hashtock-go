#!/usr/bin/env bash

v=""
if [[ $1 == '-v' ]]; then
    v="-test.v"
fi

go test $v $(ls -d ./*/ | grep -v static ) | grep -v "no test files" | grep -v "martini"
