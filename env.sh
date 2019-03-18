#!/usr/bin/env bash

source "goenv.sh"

function __build() {
    local args=(
        build "$@"
        -o bin/lex
        "github.com/evovetech/lex/cmd"
    )
    echo '$' "gollvm ${args[@]}"
    gollvm "${args[@]}"
}

function __test() {
    local args=(
        test "$@"
        "github.com/evovetech/lex/..."
    )
    echo '$' "gollvm ${args[@]}"
    gollvm "${args[@]}"
}
