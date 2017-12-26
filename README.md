## book-statistics-service

### Description

I love reading and, currently, I'm reading a lot of english books to improve my vocabulary and my language level in general.
One thing I have notice is that the appearance rhythm of new words as you move forward through pages and chapters is significantly different for each book: some books introduce a lot of new words in its beginning and the pace of introduction of new words slows down in the following parts; other books introduce less words in the beginning, but the pace of introduction of vocabulary through the rest of pages is constant.
To analyse the particular behaviour of different books concerning to the words appearance pace, I've made this service that, from a pdf book, returns somo useful statistics like the number of new words in each page, or the total number of appearances of all words contained in the book.


### Endpoints
- POST /pages-statistics (attach pdf file with "Content-Type: multipart/form-data"): returns a JSON with statistics of every page.
- POST /words-appearances (attach pdf file with "Content-Type: multipart/form-data"): returns a JSON with the total count of every word of the book.


### Requirements
- Go https://golang.org/
- dep https://github.com/golang/dep (dependency manager)


### Installation
1. Clone the repository to $GOPATH/src/ernesto-usal/
```git clone https://github.com/ernesto-usal/book-statistics-service.git```

2. Install the dependencies.
```dep ensure```

3. In the main folder of the project ($GOPATH/src/ernesto-usal/book-statistics-service) compile the service
```go build *.go```

4. Execute the generated binary to setup the service


### Notes
- Right now, the service works only with pdf files.
- The list of words in the book isn't "cleaned".