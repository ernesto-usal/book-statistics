package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	pdfcontent "github.com/unidoc/unidoc/pdf/contentstream"
	pdf "github.com/unidoc/unidoc/pdf/model"
)

type PageNewWords struct {
	NumPage      int
	NumNewWords  uint
	ListNewWords []string
}

func checkErr(err error) {
	if err != nil {
		log.Fatal("ERROR:", err)
	}
}

// Function that analyses the data of a book
func countNewWordsByPage(inputPath string) (map[string]int, []PageNewWords, error) {
	// Map with the form word:appearances_in_the_book
	mapWordCount := make(map[string]int)
	// Slice with the data analysis of every page
	var pagesNewWords []PageNewWords

	f, err := os.Open(inputPath)
	if err != nil {
		return nil, nil, err
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

	totalPages, err := pdfReader.GetNumPages()
	checkErr(err)

	for numPage := 1; numPage <= totalPages; numPage++ {
		err := countNewWordsOfPage(pdfReader, numPage, mapWordCount, &pagesNewWords)
		checkErr(err)
	}

	return mapWordCount, pagesNewWords, nil
}

// Function that counts the new words that appear in a book page
func countNewWordsOfPage(pdfReader *pdf.PdfReader, currentPageNum int, mapWordCount map[string]int, pagesNewWords *[]PageNewWords) error {
	reg, err := regexp.Compile("[^a-zA-Z]+")
	checkErr(err)

	currentPage, err := pdfReader.GetPage(currentPageNum)
	checkErr(err)

	contentStreams, err := currentPage.GetContentStreams()
	checkErr(err)

	pageContentStr := ""
	for _, cstream := range contentStreams {
		pageContentStr += " " + cstream
	}

	cstreamParser := pdfcontent.NewContentStreamParser(pageContentStr)

	rawText, err := cstreamParser.ExtractText()
	checkErr(err)

	cleanedWords := reg.ReplaceAllString(strings.ToLower(rawText), " ")
	arrayWords := strings.Fields(cleanedWords)

	pageNewWords := PageNewWords{NumPage: currentPageNum, NumNewWords: 0, ListNewWords: []string{}}

	for _, word := range arrayWords {
		if mapWordCount[word] == 0 {
			pageNewWords.NumNewWords = pageNewWords.NumNewWords + 1
			pageNewWords.ListNewWords = append(pageNewWords.ListNewWords, word)
		}
		mapWordCount[word]++
	}
	*pagesNewWords = append(*pagesNewWords, pageNewWords)

	return nil
}

// Returns the name of the book from the file path
func getBookNameOfInputPath(inputPath string) string {
	inputPathSplitted := strings.Split(inputPath, "/")
	bookName := strings.Split(inputPathSplitted[len(inputPathSplitted)-1], ".")[0]
	return bookName
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: go run bookStatistics.go input.pdf\n")
		os.Exit(1)
	}
	inputPath := os.Args[1]
	bookName := getBookNameOfInputPath(inputPath)

	// Process the books data
	startTime := time.Now()
	mapWordCount, pagesNewWords, err := countNewWordsByPage(inputPath)
	executionTime := time.Since(startTime)
	fmt.Printf("Execution time -> %s\n", executionTime)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Write the results to files
	pagesNewWordsJSON, err := json.Marshal(pagesNewWords)
	wordsCountJSON, err := json.Marshal(mapWordCount)

	f1, err := os.Create("words-list-by-page-" + strings.ToUpper(bookName) + ".json")
	defer f1.Close()
	f1.Write(pagesNewWordsJSON)

	f2, err := os.Create("words-list-count-" + strings.ToUpper(bookName) + ".json")
	defer f2.Close()
	f2.Write(wordsCountJSON)
}
