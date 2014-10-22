#!/usr/bin/env bash

v=""
if [[ $1 == '-v' ]]; then
    v="-test.v"
fi

cd app && goapp test $v
