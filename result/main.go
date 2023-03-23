package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {

	fi, err := os.Open("./vmess.txt")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	var resultLine = []string{}
	for {
		line, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		if len(strings.Trim(string(line), "\r\n")) > 0 {
			resultLine = append(resultLine, string(line))
		}
	}
	joinStr := strings.Join(resultLine, "|")

	file, err := os.OpenFile("result.txt", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	defer file.Close()
	if err != nil {
		fmt.Println(err)
	}
	file.WriteString(joinStr)

}
