package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
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
	xmlFile, err := os.Open("data/ena-samples.xml")
	if err != nil {
		fmt.Println(err)
	}
	defer xmlFile.Close()
	decoder := xml.NewDecoder(xmlFile)
	total := 0
	var inElement string
	for {
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			inElement = se.Name.Local
			if inElement == "SAMPLE" {
				var s Sample
				decoder.DecodeElement(&s, &se)
				b, err := json.Marshal(s)
				if err != nil {
					fmt.Println("error:", err)
				}
				os.Stdout.Write(b)
				total++
			}
		}
	}
	fmt.Printf("\n%d samples", total)
}
