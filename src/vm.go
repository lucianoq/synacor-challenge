package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"
)

const limit = 32768

type VM struct {
	Memory     []uint16
	ptr        uint16
	register   [8]uint16
	stack      stack
	reader     *bufio.Reader
	writer     io.Writer
	dumpWriter io.Writer
}

func NewVM(r *bufio.Reader, w io.Writer, dumpW io.Writer) *VM {
	return &VM{
		Memory:     make([]uint16, 0),
		stack:      make(stack, 0),
		reader:     r,
		writer:     w,
		dumpWriter: dumpW,
	}
}

func (vm *VM) loadProgram(r io.Reader) {
	buf := make([]byte, 2)

	for {
		n, err := r.Read(buf)
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatal(err)
		}
		if n != 2 {
			log.Fatal("no 2 bytes read")
		}

		vm.Memory = append(vm.Memory, binary.LittleEndian.Uint16(buf))
	}
}

func (vm *VM) Run() {
	for {
		opCode := vm.next()

		fmt.Fprintf(vm.dumpWriter, "[%d]\t%s\t", vm.ptr-1, names[opCode])

		switch opCode {
		case halt:
			// halt: 0
			//   stop execution and terminate the program
			return
		case set:
			vm.set()
		case push:
			vm.push()
		case pop:
			vm.pop()
		case eq:
			vm.eq()
		case gt:
			vm.gt()
		case jmp:
			vm.jmp()
		case jt:
			vm.jt()
		case jf:
			vm.jf()
		case add:
			vm.add()
		case mult:
			vm.mult()
		case mod:
			vm.mod()
		case and:
			vm.and()
		case or:
			vm.or()
		case not:
			vm.not()
		case rmem:
			vm.rmem()
		case wmem:
			vm.wmem()
		case call:
			vm.call()
		case ret:
			vm.ret()
		case out:
			vm.out()
		case in:
			vm.in()
		case noop:
			// noop: 21
			//   no operation
		default:
			log.Fatal("opcode unknown")
		}

		fmt.Fprintf(vm.dumpWriter, "\t%+v\n", vm.register)
	}
}

func (vm *VM) next() uint16 {
	x := vm.Memory[vm.ptr]
	vm.ptr++
	return x
}

func (vm *VM) readVal() uint16 {
	x := vm.next()
	if x < limit {
		fmt.Fprintf(vm.dumpWriter, "%d  ", x)
		return x
	}
	fmt.Fprintf(vm.dumpWriter, "<%d>{%d}  ", x-limit, vm.register[x-limit])
	return vm.register[x-limit]
}

func (vm *VM) readRegistryIndex() uint16 {
	x := vm.next()
	if x < limit {
		log.Fatalf("%d is not a registry index", vm.Memory[vm.ptr-1])
	}
	fmt.Fprintf(vm.dumpWriter, "<%d>  ", x-limit)
	return x - limit
}
