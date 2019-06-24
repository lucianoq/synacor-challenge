package main

import (
	"bufio"
	"log"
	"os"
)

func main() {
	// Input from STDIN
	in := bufio.NewReader(os.Stdin)

	// Output to STDOUT
	out := os.Stdout

	// Execution dump in a file
	dumpOut, err := os.Create("debug.dump")
	if err != nil {
		log.Fatal(err)
	}
	defer dumpOut.Close()

	// Create VM
	vm := NewVM(in, out, dumpOut)

	// Read code from file
	f, err := os.Open("challenge.bin")
	if err != nil {
		log.Fatal(err)
	}

	vm.loadProgram(f)

	// Uncomment the following code to go to the beach
	// to find 7th and 8th tokens.
	hack(vm)

	vm.Run()
}

func hack(vm *VM) {
	// This is the result of Ackermann function
	// with input 4 and 1 and modulo 32768.
	vm.register[7] = 25734

	// crack to skip the "eight register" == 0 test
	// JT becomes JF
	vm.Memory[521] = 8

	// overwrite function 6027 beginning from the bloody recursive hell
	// to a:
	//   SET reg<0> 6
	//   RET
	vm.Memory[6027] = set
	vm.Memory[6028] = 32768 //register 0
	vm.Memory[6029] = 6     //found in the code @5491 -> EQ r1 r0 6
	vm.Memory[6030] = ret
}
