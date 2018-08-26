package main

import (
	"fmt"
	"github.com/auth0/go-jwt-middleware"
	"github.com/codetaming/indy-ingest/cmd/ingestd/handlers"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"os"
)

func init() {
	fmt.Println("AWS_REGION:", os.Getenv("AWS_REGION"))
	fmt.Println("DATASET_TABLE:", os.Getenv("DATASET_TABLE"))
	fmt.Println("METADATA_TABLE:", os.Getenv("METADATA_TABLE"))
	fmt.Println("METADATA_BUCKET:", os.Getenv("METADATA_BUCKET"))
}

func main() {
	r := mux.NewRouter()
	ar := mux.NewRouter()

	mw := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		},
		SigningMethod: jwt.SigningMethodRS256,
	})

	r.HandleFunc("/validate", handlers.Validate).Methods("POST")
	r.HandleFunc("/dataset", handlers.CreateDataset).Methods("POST")
	r.HandleFunc("/dataset", handlers.ListDatasets).Methods("GET")
	r.HandleFunc("/dataset/{datasetId}", handlers.GetDataset).Methods("GET")
	r.HandleFunc("/dataset/{datasetId}/metadata", handlers.AddMetadata).Methods("POST")
	r.HandleFunc("/dataset/{datasetId}/metadata", handlers.ListMetadata).Methods("GET")
	r.HandleFunc("/dataset/{datasetId}/metadata/{metadataId}", handlers.GetMetadata).Methods("GET")

	an := negroni.New(negroni.HandlerFunc(mw.HandlerWithNext), negroni.Wrap(ar))
	r.PathPrefix("/").Handler(an)

	n := negroni.Classic()
	n.UseHandler(r)
	n.Run(":9000")
}
