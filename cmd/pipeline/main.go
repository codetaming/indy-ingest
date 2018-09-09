package main

import (
	"compress/gzip"
	"encoding/json"
	"encoding/xml"
	"flag"
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
	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		fmt.Printf("    pipeline inputFile ...\n")
		flag.PrintDefaults()
	}
	limitPtr := flag.Int("limit", 0, "limit on number of samples to process")
	outputFilePtr := flag.String("outputFile", "samples.json.gz", "output file name")
	flag.Parse()
	if len(flag.Args()) != 1 {
		fmt.Println("inputFile is required")
		os.Exit(1)
	}
	inputFile := flag.Args()[0]
	parse(inputFile, load.NewFileLoader(*outputFilePtr), *limitPtr)
}

func parse(inputFile string, loader load.Loader, limit int) {
	xmlFile, err := os.Open(inputFile)
	if err != nil {
		log.Fatal(err)
	}
	gz, err := gzip.NewReader(xmlFile)
	if err != nil {
		log.Fatal(err)
	}
	defer xmlFile.Close()
	defer gz.Close()
	decoder := xml.NewDecoder(gz)
	if err != nil {
		log.Fatal(err)
	}
	total := 0
	loader.Start()
	for {
		t, err := decoder.Token()
		if err != nil {
			log.Fatal(err)
		}
		if t == nil || (limit > 0 && total >= limit) {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			inElement := se.Name.Local
			if inElement == "SAMPLE" {
				var s XmlSample
				decoder.DecodeElement(&s, &se)
				js := convertToJson(s)
				jsonData, err := json.Marshal(js)
				if err != nil {
					log.Fatal("error:", err)
				}
				loader.Store(jsonData, &total, limit)
				if total%1000 == 0 {
					fmt.Printf("%d\n", total)
				}
			}
		}
	}
	loader.Finish()
	fmt.Printf("\nExported %d samples", total)
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
