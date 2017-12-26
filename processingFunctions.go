/* Set of functions that process the data of a book*/

package main

import (
	"io"
	"regexp"
	"strings"

	pdfcontent "github.com/unidoc/unidoc/pdf/contentstream"
	pdf "github.com/unidoc/unidoc/pdf/model"
)

// Function that analyses the data of a book
func processBookData(file io.ReadSeeker) (map[string]int, []PageStatistics, error) {
	wordsAppearances := make(map[string]int)
	var pagesStatistics []PageStatistics

	pdfReader, err := pdf.NewPdfReader(file)
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
		err := processPageData(pdfReader, numPage, wordsAppearances, &pagesStatistics)
		checkErr(err)
	}

	return wordsAppearances, pagesStatistics, nil
}

// Function that analysis the data of a page
func processPageData(pdfReader *pdf.PdfReader, currentPageNum int, wordsAppearances map[string]int, pagesStatistics *[]PageStatistics) error {
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

	pageNewWords := PageStatistics{NumPage: currentPageNum, NumNewWords: 0, ListNewWords: []string{}}

	for _, word := range arrayWords {
		if wordsAppearances[word] == 0 {
			pageNewWords.NumNewWords = pageNewWords.NumNewWords + 1
			pageNewWords.ListNewWords = append(pageNewWords.ListNewWords, word)
		}
		wordsAppearances[word]++
	}
	*pagesStatistics = append(*pagesStatistics, pageNewWords)

	return nil
}
