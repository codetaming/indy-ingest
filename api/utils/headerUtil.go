package utils

import (
	"errors"
	"github.com/tomnomnom/linkheader"
)

func ExtractSchemaUrl(headers map[string]string) (url string, err error) {
	linkHeader := headers["link"]
	if linkHeader == "" {
		linkHeader = headers["Link"]
	}
	if linkHeader == "" {
		return "", errors.New("Link header must be provided")
	}
	links := linkheader.Parse(linkHeader)
	var link linkheader.Link
	if len(links) == 1 {
		link = links[0]
	}
	return link.URL, nil
}
