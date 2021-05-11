package main

// #cgo CFLAGS: -g -Wall
// #include <stdlib.h>
// #include "greeter.h"
import "C"

import (
	"fmt"
	"unsafe"
)

func main() {

	name := C.CString("Gopher")
	defer C.free(unsafe.Pointer(name))

	year := C.int(2018)
	number := C.greet(name, year)
	fmt.Println(number)
}
