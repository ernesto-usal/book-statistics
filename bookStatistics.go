package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	pdfcontent "github.com/unidoc/unidoc/pdf/contentstream"
	pdf "github.com/unidoc/unidoc/pdf/model"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: go run count-words.go input.pdf\n")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	startTime := time.Now()
	err := countNewWordsByPage(inputPath)
	executionTime := time.Since(startTime)
	fmt.Printf("Execution time -> %s", executionTime)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatal("ERROR:", err)
	}
}

func countNewWordsByPage(inputPath string) error {
	mapWordCount := make(map[string]int)
	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	pdfReader, err := pdf.NewPdfReader(f)
	checkErr(err)

	isEncrypted, err := pdfReader.IsEncrypted()
	checkErr(err)
	if isEncrypted {
		_, err = pdfReader.Decrypt([]byte(""))
		checkErr(err)
	}

	numPages, err := pdfReader.GetNumPages()
	checkErr(err)

	startTimeOuterLoop := time.Now()
	for i := 1; i <= numPages; i++ {
		startTimePageLoop := time.Now()
		newWordsInPage, countWordsPage, err := countNewWordsOfPage(pdfReader, i, mapWordCount)
		checkErr(err)

		executionTimePageLoop := time.Since(startTimePageLoop)
		fmt.Println("Page " +
			strconv.Itoa(i) + " - " +
			strconv.Itoa(newWordsInPage) +
			" new words of " +
			strconv.Itoa(countWordsPage) +
			" (Execution time regular loop-> " +
			executionTimePageLoop.String() +
			")")
	}
	executionTimeOuterLoop := time.Since(startTimeOuterLoop)
	fmt.Println("Execution time outer loop -> " + executionTimeOuterLoop.String())
	return nil
}

func countNewWordsOfPage(pdfReader *pdf.PdfReader, currentPageNum int, mapWordCount map[string]int) (int, int, error) {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	checkErr(err)
	newWordsInPage := 0

	currentPage, err := pdfReader.GetPage(currentPageNum)
	checkErr(err)

	contentStreams, err := currentPage.GetContentStreams()
	checkErr(err)

	pageContentStr := ""
	for _, cstream := range contentStreams {
		pageContentStr += cstream
	}

	cstreamParser := pdfcontent.NewContentStreamParser(pageContentStr)

	txt, err := cstreamParser.ExtractText()
	checkErr(err)

	cleanedWords := reg.ReplaceAllString(strings.ToLower(txt), " ")
	arrayWords := strings.Fields(cleanedWords)

	for _, word := range arrayWords {
		if mapWordCount[word] == 0 {
			mapWordCount[word]++
			newWordsInPage++
		} else {
			mapWordCount[word]++
		}
	}

	return newWordsInPage, len(arrayWords), nil
}
