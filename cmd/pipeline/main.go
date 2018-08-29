package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/codetaming/indy-ingest/internal/load"
	"log"
	"os"
)

type XmlSample struct {
	Title            string               `xml:"TITLE"`
	Accession        string               `xml:"accession,attr" `
	SampleAttributes []XmlSampleAttribute `xml:"SAMPLE_ATTRIBUTES>SAMPLE_ATTRIBUTE"`
}

type XmlSampleAttribute struct {
	Tag   string `xml:"TAG"`
	Value string `xml:"VALUE"`
}

type JsonSample struct {
	Title           string            `json:"name"`
	Accession       string            `json:"accession"`
	Characteristics map[string]string `json:"characteristics"`
}

func main() {
	streamingParse()
}

func streamingParse() {
	xmlFile, err := os.Open("data/ena-samples.xml.gz")
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
	limit := 1000
	loader := load.NewFileLoader()
	for {
		t, _ := decoder.Token()
		if t == nil || total >= limit {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			inElement = se.Name.Local
			if inElement == "SAMPLE" {
				var s XmlSample
				decoder.DecodeElement(&s, &se)
				js := convertToJson(s)
				jsonData, err := json.Marshal(js)
				if err != nil {
					fmt.Println("error:", err)
				}
				loader.Load(jsonData)
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
	err = loader.Close()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\n%d samples", total)
}

func convertToJson(xs XmlSample) JsonSample {
	characteristics := make(map[string]string)
	for _, sa := range xs.SampleAttributes {
		characteristics[sa.Tag] = sa.Value
	}
	js := JsonSample{
		Title:           xs.Title,
		Accession:       xs.Accession,
		Characteristics: characteristics,
	}
	return js
}
