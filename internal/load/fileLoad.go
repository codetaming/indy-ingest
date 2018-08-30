package load

import (
	"bytes"
	"compress/gzip"
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
	w.Write(jsonData)
	*total++
	if limit == 0 || *total < limit {
		w.Write([]byte(",\n"))
	}
}

func (f *FileLoader) Finish() {
	w.Write([]byte("]"))
	w.Close()
	err := ioutil.WriteFile(f.outputFile, b.Bytes(), 0666)
	if err != nil {
		log.Fatal(err)
	}
}
