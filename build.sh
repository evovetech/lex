#!/usr/bin/env bash

function gopath() {
    local paths=( $( tr ':' '\n' <<< "$GOPATH" ) )
    printf "${paths[0]}"
}

function build() {
    local args=(
        build "$@"
        -ldflags="-r $(gopath)/src/llvm.org/llvm/bindings/go/llvm/workdir/llvm_build/lib"
        -o ../bin/lex
        .
    )
    echo '$' "go ${args[@]}"
    cd cmd && go "${args[@]}"
}

# main
build "$@"
