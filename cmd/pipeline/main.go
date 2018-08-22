package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

type SampleSet struct {
	XMLName xml.Name `xml:"SAMPLE_SET"`
	Samples []Sample `xml:"SAMPLE"`
}

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
	xmlFile, err := os.Open("data/ena-samples.xml")
	if err != nil {
		fmt.Println(err)
	}
	defer xmlFile.Close()
	byteValue, _ := ioutil.ReadAll(xmlFile)

	var sampleSet SampleSet

	err = xml.Unmarshal(byteValue, &sampleSet)
	if err != nil {
		fmt.Println(err)
	}

	for _, sample := range sampleSet.Samples {
		b, err := json.Marshal(sample)
		if err != nil {
			fmt.Println("error:", err)
		}
		os.Stdout.Write(b)
	}
}
