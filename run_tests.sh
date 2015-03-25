#!/usr/bin/env bash

v=""
if [[ $1 == '-v' ]]; then
    v="-test.v"
fi

# goapp test $v $(ls -d ./*/) | grep -v "no test files" | grep -v "martini" | grep -v -E "(INFO|devappserver2)"
go test --cover $v $(ls -d ./*/ | grep -v gae | grep -v static | grep -v remove_me ) | grep -v "no test files" | grep -v "martini"
