#!/usr/bin/env bash

function gopath() {
    local paths=( $( tr ':' '\n' <<< "$GOPATH" ) )
    printf %s "${paths[0]}"
}

function ldflags() {
    printf %s "-r $(gopath)/src/llvm.org/llvm/bindings/go/llvm/workdir/llvm_build/lib"
}

function build() {
    local args=(
        build "$@"
        -ldflags="$( ldflags )"
        -o ../bin/lex
        .
    )
    echo '$' "go ${args[@]}"
    cd cmd && go "${args[@]}"
}

function test() {
    local args=(
        test "$@"
        -ldflags="$( ldflags )"
        ./...
    )
    echo '$' "go ${args[@]}"
    go "${args[@]}"
}

export -f build
export -f test
