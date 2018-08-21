package utils

import (
	"errors"
	"github.com/tomnomnom/linkheader"
)

func parseLinkHeader(linkHeader string) (url string, err error) {
	links := linkheader.Parse(linkHeader)
	var link linkheader.Link
	if len(links) == 1 {
		link = links[0]
	}
	return link.URL, nil
}

func ExtractSchemaUrlArray(headers map[string][]string) (url string, err error) {
	linkHeader := ""
	if value, ok := headers["link"]; ok {
		linkHeader = value[0]
	}
	if value, ok := headers["Link"]; ok {
		linkHeader = value[0]
	}
	if linkHeader == "" {
		return "", errors.New("Link header must be provided")
	}
	return parseLinkHeader(linkHeader)
}

func ExtractSchemaUrl(headers map[string]string) (url string, err error) {
	linkHeader := headers["link"]
	if linkHeader == "" {
		linkHeader = headers["Link"]
	}
	if linkHeader == "" {
		return "", errors.New("Link header must be provided")
	}
	return parseLinkHeader(linkHeader)
}
