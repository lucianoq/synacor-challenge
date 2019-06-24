package main

import (
	"fmt"
	"log"
	"os"
)

const (
	halt uint16 = iota
	set
	push
	pop
	eq
	gt
	jmp
	jt
	jf
	add
	mult
	mod
	and
	or
	not
	rmem
	wmem
	call
	ret
	out
	in
	noop
)

var names = []string{"HALT", "SET", "PUSH", "POP", "EQ", "GT", "JMP", "JT",
	"JF", "ADD", "MULT", "MOD", "AND", "OR", "NOT", "RMEM", "WMEM", "CALL",
	"RET", "OUT", "IN", "NOOP"}

// set: 1 a b
//   set vm.register <a> to the value of <b>
func (vm *VM) set() {
	a := vm.readRegistryIndex()
	b := vm.readVal()
	vm.register[a] = b
}

// push: 2 a
//   push <a> onto the stack
func (vm *VM) push() {
	a := vm.readVal()
	vm.stack.Push(a)
}

// pop: 3 a
//   remove the top element from the stack and write it into <a>; empty stack = error
func (vm *VM) pop() {
	a := vm.readRegistryIndex()
	x, err := vm.stack.Pop()
	if err == EmptyErr {
		os.Exit(0)
	}
	vm.register[a] = x
}

// eq: 4 a b c
//   set <a> to 1 if <b> is equal to <c>; set it to 0 otherwise
func (vm *VM) eq() {
	a := vm.readRegistryIndex()
	b := vm.readVal()
	c := vm.readVal()
	if b == c {
		vm.register[a] = 1
	} else {
		vm.register[a] = 0
	}
}

// gt: 5 a b c
//   set <a> to 1 if <b> is greater than <c>; set it to 0 otherwise
func (vm *VM) gt() {
	a := vm.readRegistryIndex()
	b := vm.readVal()
	c := vm.readVal()
	if b > c {
		vm.register[a] = 1
	} else {
		vm.register[a] = 0
	}
}

// jmp: 6 a
//   jump to <a>
func (vm *VM) jmp() {
	a := vm.readVal()
	vm.ptr = a
}

// jt: 7 a b
//   if <a> is nonzero, jump to <b>
func (vm *VM) jt() {
	a := vm.readVal()
	b := vm.readVal()
	if a != 0 {
		vm.ptr = b
	}
}

// jf: 8 a b
//   if <a> is zero, jump to <b>
func (vm *VM) jf() {
	a := vm.readVal()
	b := vm.readVal()
	if a == 0 {
		vm.ptr = b
	}
}

// add: 9 a b c
//   assign into <a> the sum of <b> and <c> (modulo 32768)
func (vm *VM) add() {
	a := vm.readRegistryIndex()
	b := vm.readVal()
	c := vm.readVal()
	vm.register[a] = (b + c) % limit
}

// mult: 10 a b c
//   store into <a> the product of <b> and <c> (modulo 32768)
func (vm *VM) mult() {
	a := vm.readRegistryIndex()
	b := vm.readVal()
	c := vm.readVal()

	b64 := uint64(b)
	c64 := uint64(c)
	vm.register[a] = uint16((b64 * c64) % limit)
}

// mod: 11 a b c
//   store into <a> the remainder of <b> divided by <c>
func (vm *VM) mod() {
	a := vm.readRegistryIndex()
	b := vm.readVal()
	c := vm.readVal()
	vm.register[a] = b % c
}

// and: 12 a b c
//   stores into <a> the bitwise and of <b> and <c>
func (vm *VM) and() {
	a := vm.readRegistryIndex()
	b := vm.readVal()
	c := vm.readVal()
	vm.register[a] = b & c
}

// or: 13 a b c
//   stores into <a> the bitwise or of <b> and <c>
func (vm *VM) or() {
	a := vm.readRegistryIndex()
	b := vm.readVal()
	c := vm.readVal()
	vm.register[a] = b | c
}

// not: 14 a b
//   stores 15-bit bitwise inverse of <b> in <a>
func (vm *VM) not() {
	a := vm.readRegistryIndex()
	b := vm.readVal()
	vm.register[a] = 0x7FFF ^ b
}

// rmem: 15 a b
//   read memory at address <b> and write it to <a>
func (vm *VM) rmem() {
	a := vm.readRegistryIndex()
	b := vm.readVal()
	vm.register[a] = vm.Memory[b]
}

// wmem: 16 a b
//   write the value from <b> into memory at address <a>
func (vm *VM) wmem() {
	a := vm.readVal()
	b := vm.readVal()
	vm.Memory[a] = b
}

// call: 17 a
//   write the address of the next instruction to the stack and jump to <a>
func (vm *VM) call() {
	a := vm.readVal()
	vm.stack.Push(vm.ptr)
	vm.ptr = a
}

// ret: 18
//   remove the top element from the stack and jump to it; empty stack = halt
func (vm *VM) ret() {
	x, err := vm.stack.Pop()
	if err == EmptyErr {
		os.Exit(0)
	}
	vm.ptr = x
}

// out: 19 a
//   write the character represented by ascii code <a> to the terminal
func (vm *VM) out() {
	a := vm.readVal()
	fmt.Fprintf(vm.writer, "%s", string(a))
}

// in: 20 a
//   read a character from the terminal and write its ascii code to <a>; it can
//   be assumed that once input starts, it will continue until a newline is
//   encountered; this means that you can safely read whole lines from the
//   keyboard and trust that they will be fully read
func (vm *VM) in() {
	a := vm.readRegistryIndex()
	char, _, err := vm.reader.ReadRune()
	if err != nil {
		log.Fatal("error reading")
	}
	vm.register[a] = uint16([]byte(string(char))[0])
}
