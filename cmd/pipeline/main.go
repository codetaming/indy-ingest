package main

import (
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
	Title      string            `xml:"TITLE"`
	Attributes []SampleAttribute `xml:"SAMPLE_ATTRIBUTES"`
}

type SampleAttribute struct {
	XMLName xml.Name `xml:"SAMPLE_ATTRIBUTE"`
	Tag     string   `xml:TAG`
	Value   string   `xml:VALUE`
}

func main() {
	xmlFile, err := os.Open("data/ena-samples.xml")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened users.xml")
	defer xmlFile.Close()
	byteValue, _ := ioutil.ReadAll(xmlFile)

	var sampleSet SampleSet

	err = xml.Unmarshal(byteValue, &sampleSet)
	if err != nil {
		fmt.Println(err)
	}

	for _, sample := range sampleSet.Samples {
		fmt.Println(sample.Title)
		for _, attribute := range sample.Attributes {
			fmt.Printf(`%s:%s\n`, attribute.Tag, attribute.Value)
		}
	}
}
