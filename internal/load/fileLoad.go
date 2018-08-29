package load

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

type FileLoader struct {
	l Loader
	c Closer
}

var b bytes.Buffer
var w gzip.Writer

func NewFileLoader() *FileLoader {
	l := new(FileLoader)
	return l
}

func (FileLoader) Load(json []byte) {
	w.Write(json)
}

func (FileLoader) Close() error {
	w.Write([]byte("]"))
	w.Close()
	err := ioutil.WriteFile("data/samples.json.gz", b.Bytes(), 0666)
	return err
}

func (l *FileLoader) init() {
	*l = FileLoader{
		l: l,
	}
	w := gzip.NewWriter(&b)
	w.Write([]byte("["))
}
