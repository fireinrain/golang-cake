package main

import (
	"fmt"
	"time"
)

func main() {
	for i := 0; i < 100; i++ {
		go func(numer int) {
			fmt.Println(numer)
		}(i)
	}
	time.Sleep(100 * time.Millisecond)
	fmt.Println("---------------------")
	for i := 0; i < 10; i++ {
		ii := i
		go func() {
			fmt.Println(ii)
		}()

	}
	time.Sleep(1 * time.Second)
}
