package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/auth0/go-jwt-middleware"
	"github.com/codetaming/indy-ingest/cmd/ingestd/handlers"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"net/http"
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
			// Verify 'aud' claim
			aud := os.Getenv("AUTH0_AUDIENCE")
			checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
			if !checkAud {
				return token, errors.New("invalid audience")
			}
			// Verify 'iss' claim
			iss := os.Getenv("AUTH0_ISSUER")
			checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
			if !checkIss {
				return token, errors.New("invalid issuer")
			}

			cert, err := getPemCert(token)
			if err != nil {
				panic(err.Error())
			}

			result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
			return result, nil
		},
		SigningMethod: jwt.SigningMethodRS256,
	})

	r.HandleFunc("/validate", handlers.Validate).Methods("POST")
	ar.HandleFunc("/dataset", handlers.CreateDataset).Methods("POST")
	ar.HandleFunc("/dataset", handlers.ListDatasets).Methods("GET")
	ar.HandleFunc("/dataset/{datasetId}", handlers.GetDataset).Methods("GET")
	ar.HandleFunc("/dataset/{datasetId}/metadata", handlers.AddMetadata).Methods("POST")
	ar.HandleFunc("/dataset/{datasetId}/metadata", handlers.ListMetadata).Methods("GET")
	ar.HandleFunc("/dataset/{datasetId}/metadata/{metadataId}", handlers.GetMetadata).Methods("GET")

	an := negroni.New(negroni.HandlerFunc(mw.HandlerWithNext), negroni.Wrap(ar))
	r.PathPrefix("/").Handler(an)

	n := negroni.Classic()
	n.UseHandler(r)
	n.Run(":9000")
}

type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

func getPemCert(token *jwt.Token) (string, error) {
	cert := ""
	resp, err := http.Get(os.Getenv("AUTH0_ISSUER") + ".well-known/jwks.json")

	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()

	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return cert, err
	}

	for k := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("unable to find appropriate key")
		return cert, err
	}

	return cert, nil
}
