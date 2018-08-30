package jit

import "llvm.org/llvm/bindings/go/llvm"

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
