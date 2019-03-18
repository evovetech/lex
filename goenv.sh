#!/bin/bash

# Go env
if [[ -n "${__goenv}" ]]; then
    return 0
fi

# functions
function __run() {
    echo "$" "$@"
    # run
    "$@"
}

function __ldflags() {
    printf %s "-r `llvm-config --libdir`"
}

function gollvm() {
  local cmd="help"
  local args=()
  if [[ $# -gt 0 ]]; then
    cmd="${1}"
    args+=(
      -ldflags="`__ldflags`"
      "${@:2}"
    )
  else
    echo "must pass go argument. defaulting to help"
  fi
  __run llvm-go "${cmd}" "${args[@]}"
}

echo "goenv init!!"
eval "$( goenv init - )"
export GOVERSION="$( goenv version-name )"
export GOROOT="$( go env GOROOT )"
export GOPATH_MAIN="${HOME}/go/${GOVERSION}"
export GOPATH_DEV="${HOME}/Development/go"
export GOPATH="${GOPATH_MAIN}:${GOPATH_DEV}"
export GOBIN="${GOPATH_MAIN}/bin"
export GOLLVM_ROOTPATH="${GOPATH_MAIN}/src/llvm.org/llvm/bindings/go/llvm"
export GOLLVM_INCLUDEPATH="${GOLLVM_ROOTPATH}/include"
export GOLLVM_BUILDDIR="${GOLLVM_ROOTPATH}/workdir/llvm_build"
export GOLLVM_LIBPATH="${GOLLVM_BUILDDIR}/lib"
export GOLLVM_BINPATH="${GOLLVM_BUILDDIR}/bin"
export PATH="${PATH}:${GOBIN}:${GOPATH_DEV}/bin:${GOLLVM_BINPATH}"
export -f __run
export -f __ldflags
export -f gollvm
export __goenv="goenv"
