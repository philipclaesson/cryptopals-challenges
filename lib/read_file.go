package lib

import (
	"bufio"
	"fmt"
	"os"
)

func ReadDataFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer file.Close()

	payload := ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		payload += line
	}
	return payload, nil
}
