package load

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"log"
)

var b bytes.Buffer
var w *gzip.Writer

type FileLoader struct {
	outputFile string
}

func NewFileLoader(outputFile string) *FileLoader {
	f := new(FileLoader)
	f.outputFile = outputFile
	return f
}

func (f *FileLoader) Start() {
	w = gzip.NewWriter(&b)
	w.Write([]byte("["))
}

func (f *FileLoader) Store(jsonData []byte, total *int, limit int) {
	_, err := w.Write(jsonData)
	if err != nil {
		log.Fatal(err)
	}
	*total++
	if limit == 0 || *total < limit {
		w.Write([]byte(",\n"))
	} else {
		fmt.Printf("Limit %d reached\n", limit)
	}
}

func (f *FileLoader) Finish() {
	_, writeErr := w.Write([]byte("]"))
	if writeErr != nil {
		log.Fatal(writeErr)
	}
	closeError := w.Close()
	if closeError != nil {
		log.Fatal(closeError)
	}
	err := ioutil.WriteFile(f.outputFile, b.Bytes(), 0666)
	if err != nil {
		log.Fatal(err)
	}
}
