package jit

import "llvm.org/llvm/bindings/go/llvm"

func init() {
	// initialize llvm target
	llvm.InitializeNativeTarget()
	llvm.InitializeNativeAsmPrinter()
	llvm.InitializeAllAsmParsers()
}
