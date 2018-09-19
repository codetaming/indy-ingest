package main

import (
	"encoding/json"
	"fmt"
	"github.com/codetaming/indy-ingest/internal/model"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {
	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://localhost:9000/validate", strings.NewReader("{}"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Link", `<https://schema.humancellatlas.org/type/biomaterial/5.1.0/specimen_from_organism>; rel="describedby"`)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error with request: %s", err.Error())
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		var v model.ValidationResult
		json.Unmarshal(body, &v)
		fmt.Printf("Result: %s", v.Message)
	}
}
