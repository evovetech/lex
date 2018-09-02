package jit

import (
	"sync"

	"github.com/evovetech/lex/compiler"
	"llvm.org/llvm/bindings/go/llvm"
)

type JIT struct {
	OptLevel int

	machine *compiler.Machine
	pm      llvm.PassManager

	init sync.Once
}

func (j *JIT) Init() *JIT {
	j.init.Do(j.initialize)
	return j
}

func (j *JIT) initialize() {
	// setup target stuff
	var err error
	j.machine, err = compiler.NewDefaultMachine()
	if err != nil {
		panic(err)
	}

	passManager := llvm.NewPassManager()
	passBuilder := llvm.NewPassManagerBuilder()
	if j.OptLevel > 0 {
		passBuilder.SetOptLevel(j.OptLevel)
		passBuilder.Populate(passManager)
	}
	j.pm = passManager
}
