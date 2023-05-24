package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"testing"
)

// clean result ip file
func TestResult(t *testing.T) {
	// Open the text file
	file, err := os.Open("resultips.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Create an empty slice to store non-empty lines
	lines := []string{}

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Iterate over each line
	for scanner.Scan() {
		line := scanner.Text()

		// Check if the line is empty or contains only whitespace
		if len(line) == 0 || line == "" {
			continue // Skip empty lines
		}

		// Append non-empty lines to the slice
		lines = append(lines, line)
	}

	// Check for any scanner errors
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Print the non-empty lines
	for _, line := range lines {
		fmt.Println(line)
	}

}
