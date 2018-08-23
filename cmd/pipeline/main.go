package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type Sample struct {
	Title            string            `xml:"TITLE" json:"name"`
	Accession        string            `xml:"accession,attr" json:"accession"`
	SampleAttributes []SampleAttribute `xml:"SAMPLE_ATTRIBUTES>SAMPLE_ATTRIBUTE" json:"characteristics"`
}

type SampleAttribute struct {
	Tag   string `xml:"TAG" json:"tag"`
	Value string `xml:"VALUE" json:"value"`
}

func main() {
	streamingParse()
}

func streamingParse() {
	xmlFile, err := os.Open("data/sample.xml.gz")
	if err != nil {
		fmt.Println(err)
	}

	gz, err := gzip.NewReader(xmlFile)

	if err != nil {
		log.Fatal(err)
	}

	defer xmlFile.Close()
	defer gz.Close()

	decoder := xml.NewDecoder(gz)
	total := 0
	var inElement string
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	limit := 1000000000
	w.Write([]byte("["))
	for {
		t, _ := decoder.Token()
		if t == nil || total >= limit {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			inElement = se.Name.Local
			if inElement == "SAMPLE" {
				var s Sample
				decoder.DecodeElement(&s, &se)
				jsonData, err := json.Marshal(s)
				if err != nil {
					fmt.Println("error:", err)
				}
				w.Write(jsonData)
				total++
				if total < limit {
					w.Write([]byte(",\n"))
				}
				if total%1000 == 0 {
					fmt.Printf("%d\n", total)
				}
			}
		}
	}
	w.Write([]byte("]"))
	w.Close()
	err = ioutil.WriteFile("samples.json.gz", b.Bytes(), 0666)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\n%d samples", total)
}
