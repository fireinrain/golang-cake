package main

import (
	"fmt"
	"unsafe"
)

func main() {
	var a float64 = 12.3434342342
	b := float32(a)
	fmt.Println(b)
	var name string = "你好"
	size := unsafe.Sizeof(name)
	fmt.Printf("name occupy %d bytes\n", size)

}

func Bigger(a any, b any) {

}
func CreatUser(name, address string, age int) {

}
func CreateUser2(age int, name, adress string) {

}
