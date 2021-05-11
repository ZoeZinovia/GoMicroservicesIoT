package main

// #include <stdio.h>
// #include <stdlib.h>
// int greet(const char *name, int year) {
//     int n = 2;
//     return n;
// }
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
