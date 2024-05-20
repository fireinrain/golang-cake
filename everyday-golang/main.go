package main

import (
	"fmt"
	"github.com/alexellis/hmac"
	"os"
)

func main() {
	name := os.Getenv("USER")

	fmt.Printf("Well done %s for having your first Go\n", name)
	//双引号字符串字面量需要对特定的字符转义
	input := []byte("https://github.com/alexellis/hmac")
	fmt.Println(input)
	digest := hmac.Sign(input, []byte(name))
	fmt.Printf("Digest: %x\n", digest)
	//反引号字符串字面量不需要对特殊字符转义

	input2 := []byte(`https://github.com/alexellis/hmac`)
	fmt.Println(input2)
	digest2 := hmac.Sign(input2, []byte(name))
	fmt.Printf("Digest: %x\n", digest2)

	//validat hmac
	err := hmac.Validate(input, fmt.Sprintf("sha1=%x", digest), name)

	if err != nil {
		panic(err)
	}
	fmt.Printf("Digest validated: %x\n", digest)

}
