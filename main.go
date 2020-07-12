package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	qrcode "github.com/skip2/go-qrcode"
)

const OUTPUT_DIR = os.Getenv("QR_OUTPUT_DIR") || "qr-codes"
const INPUT_DIR = os.Getenv("QR_INPUT_DIR") || "payloads"

// CreateDirIfNotExist - creates directory if it doesnt exist
func CreateDirIfNotExist(dir string) (e error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			e = err
		}
	}

	return
}

// ListDir - lists all files in directory
func ListDir(dir string) (fileNames []string, e error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		e = err
	}

	for _, f := range files {
		fileNames = append(fileNames, f.Name())
	}

	return
}

// ReadFile - Reads all lines in a file
func ReadFile(path string) (fileTextLines []string, e error) {
	readFile, err := os.Open(path)

	if err != nil {
		e = err
	}

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		fileTextLines = append(fileTextLines, fileScanner.Text())
	}

	readFile.Close()

	return
}

func main() {
	err := CreateDirIfNotExist(OUTPUT_DIR)
	if err != nil {
		log.Fatal(err)
	}

	fileNames, err := ListDir("./payloads")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range fileNames {
		fileLines, err := ReadFile(INPUT_DIR + "/" + file)
		if err != nil {
			log.Fatal(err)
		}

		for lineIndex, line := range fileLines {
			fileNoExt := strings.Replace(file, ".txt", "", -1)
			sLineIndex := strconv.Itoa(lineIndex)
			outputFilePath := OUTPUT_DIR + "/" + fileNoExt + "-" + sLineIndex + ".png"

			err := qrcode.WriteFile(line, qrcode.Medium, 256, outputFilePath)
			if err != nil {
				fmt.Println("Failed to generate QR Code from", file, "line:", lineIndex, "(", err, ")")
			}
		}
	}

	fmt.Println("Done generating QR Codes")
}
