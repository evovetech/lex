package compiler

import (
	"fmt"
	"os"

	"llvm.org/llvm/bindings/go/llvm"
)

type Machine struct {
	llvm.TargetMachine
	targetData *llvm.TargetData
}

func NewDefaultMachine() (*Machine, error) {
	return NewMachine(llvm.DefaultTargetTriple())
}

func NewMachine(triple string) (*Machine, error) {
	target, err := llvm.GetTargetFromTriple(triple)
	if err != nil {
		return nil, err
	}
	return &Machine{TargetMachine: target.CreateTargetMachine(
		triple,
		"",
		"",
		llvm.CodeGenLevelNone,
		llvm.RelocDefault,
		llvm.CodeModelDefault,
	)}, nil
}

func (m *Machine) TargetData() (td llvm.TargetData) {
	if m.targetData == nil {
		td = m.CreateTargetData()
		m.targetData = &td
	} else {
		td = *m.targetData
	}
	return
}

func (m *Machine) NewCompiler(name string) CompilerMachine {
	c := GlobalContext().baseCompiler(name)
	mod := c.module
	mod.SetTarget(m.Triple())
	mod.SetDataLayout(m.TargetData().String())
	c.fpm = mod.newFunctionPassManager()
	return &mc{Machine: m, compiler: c}
}

type CompilerMachine interface {
	GetMachine() *Machine
	GetCompiler() Compiler
	Write(file string, fType llvm.CodeGenFileType) error
}

type mc struct {
	*Machine
	*compiler
}

func (mc *mc) GetMachine() *Machine {
	return mc.Machine
}

func (mc *mc) GetCompiler() Compiler {
	return mc.compiler
}

func (mc *mc) Write(file string, fType llvm.CodeGenFileType) (err error) {
	mc.AddAnalysisPasses(mc.fpm)
	//
	//var mb llvm.MemoryBuffer
	//if mb, err = mc.EmitToMemoryBuffer(mc.GetModule(), fType); err != nil {
	//	return
	//}
	//defer mb.Dispose()

	var f *os.File
	if f, err = os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0644); err != nil {
		return
	}
	defer f.Close()

	if err = llvm.WriteBitcodeToFile(mc.GetModule(), f); err != nil {
		fmt.Printf("wrote file %s", file)
	}
	return
}
