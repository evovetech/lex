#!/usr/bin/env bash

function get_dir() {
    printf "%s" "$( cd "$( dirname "$@" )" >/dev/null && pwd )"
}

dir="$( get_dir "${BASH_SOURCE[0]}" )"

# build & compile average code to bitcode
cd "${dir}/.."
# ./build.sh
cat "${dir}/average.kl" | bin/lex compile

# create ll files
cd "${dir}"
llvm-dis output.bc -o output.ll
llvm-as output.ll -o output.bc
llc output.bc -filetype=asm
llc output.bc -filetype=obj

# compile main.cpp & output together
clang++ main.cpp output.o -o main
