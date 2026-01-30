package main

import (
    "bufio"
    "fmt"
    "os"
		"io/fs"
)

func main() {
	file, err := os.OpenFile("./example.md", os.O_RDONLY, fs.ModeDevice)
	if err != nil {
		fmt.Errorf("error opening the file: %w", err)
		return
	}
	defer file.Close()


	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		// This will read line by line sepparated with an enter, here I need to check for which types I need to parse and create
		fmt.Println(scanner.Text())
	}
}
