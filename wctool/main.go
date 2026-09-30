package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	args := os.Args
	var argsValues []string
	filename := ""

	for _, val := range args[1:] {
		if strings.HasPrefix(val, "-") {
			argsValues = append(argsValues, val)
		} else {
			if filename != "" {
				panic(fmt.Sprintf("Only one file can be observed, given: %s", val))
			}
			filename = val
		}
	}

	var file *os.File
	var err error

	if filename == "" {
		file, err = os.CreateTemp("", "stdin-buf-*")
		defer os.Remove(file.Name())

		var fi os.FileInfo
		fi, err = os.Stdin.Stat()

		if fi.Size() == 0 {
			panic("STDIN is empty")
		}

		io.Copy(file, os.Stdin)
	} else {
		file, err = os.Open(filename)
	}

	if err != nil {
		panic(err)
	}

	defer file.Close()

	result := ""

	if len(argsValues) == 0 {
		argsValues = []string{"-l", "-w", "-c"}
	}

	for _, val := range argsValues {
		switch val {
		case "-c":
			fi, err := file.Stat()
			if err != nil {
				panic(err)
			}
			result = result + fmt.Sprintf("%d ", fi.Size())
		case "-l":
			lines := 0
			scanFile(file, func(scanner *bufio.Scanner) {
				lines += 1
			})

			result = result + fmt.Sprintf("%d ", lines)
		case "-w":
			wordsCount := 0

			scanFile(file, func(scanner *bufio.Scanner) {
				words := strings.Fields(scanner.Text())
				wordsCount += len(words)
			})

			result = result + fmt.Sprintf("%d ", wordsCount)
		}
	}

	fmt.Printf("%s%s\n", result, filename)
}

func scanFile(file *os.File, cb func(*bufio.Scanner)) error {
	file.Seek(0, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return err
		}

		cb(scanner)
	}

	return nil
}
