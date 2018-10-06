package main

import (
	"encoding/json"
	"fmt"
	"github.com/codetaming/indy-ingest/internal/model"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func init() {
	rootCmd.AddCommand(validateCmd)
}

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate json",
	Long:  `Validate json against a schema`,
	Run: func(cmd *cobra.Command, args []string) {
		client := &http.Client{}
		req, err := http.NewRequest("POST", fmt.Sprint(ingestServerURL, "/validate"), strings.NewReader("{}"))
		if err != nil {
			log.Printf("Error making request to %s", ingestServerURL)
			os.Exit(0)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Link", `<https://schema.humancellatlas.org/type/biomaterial/5.1.0/specimen_from_organism>; rel="describedby"`)
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error with request: %s", err.Error())
		}
		body, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Error: %d\n%s", resp.StatusCode, body)
		} else {
			var v model.ValidationResult
			json.Unmarshal(body, &v)
			fmt.Printf("Result: %s", v.Message)
		}
	},
}
