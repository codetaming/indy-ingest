package main

import (
	"github.com/spf13/cobra/cobra/cmd"
	"log"
	"os"
)

var (
	ingestServerURL = os.Getenv("INGEST_SERVER_URL")
)

func init() {
	if ingestServerURL == "" {
		log.Fatal("INGEST_SERVER_URL not set")
	}
}

func main() {
	cmd.Execute()
}
