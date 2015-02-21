#!/usr/bin/env bash

v=""
if [[ $1 == '-v' ]]; then
    v="-test.v"
fi

goapp test $v $(ls -d ./*/) | grep -v "no test files" | grep -v "martini"
# cd app && goapp test $v
