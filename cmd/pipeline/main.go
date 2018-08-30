package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"github.com/briandowns/spinner"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
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
		fmt.Printf("    pipeline inputFile outputFile ...\n")
		flag.PrintDefaults()
	}
	if len(os.Args) < 3 {
		fmt.Println("inputFile and outputFile are required")
		os.Exit(1)
	}
	inputFile := os.Args[1]
	outputFile := os.Args[2]
	limitPtr := flag.Int("word", 0, "limit on number of samples to process")
	flag.Parse()
	parse(inputFile, outputFile, *limitPtr)
}

func parse(inputFile string, outputFile string, limit int) {
	xmlFile, err := os.Open(inputFile)
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
	w.Write([]byte("["))
	sp := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	sp.Start()
	for {
		t, _ := decoder.Token()
		if t == nil || (limit > 0 && total >= limit) {
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
				w.Write(jsonData)

				total++
				if total < limit {
					w.Write([]byte(",\n"))
				}
				if total%1000 == 0 {
					sp.Suffix = ": " + strconv.Itoa(total) + " samples"
					//fmt.Printf("%d\n", total)
				}

			}
		}
	}
	sp.Stop()
	w.Write([]byte("]"))
	w.Close()
	err = ioutil.WriteFile(outputFile, b.Bytes(), 0666)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nExported %d samples to %s", total, outputFile)
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
