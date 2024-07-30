package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"sample-broker/osbapi"
)

const (
	hardcodedUserName = "broker-user"
	hardcodedPassword = "broker-password"
)

func main() {
	http.HandleFunc("/", helloWorldHandler)
	http.HandleFunc("/v2/catalog", getCatalogHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Listening on port %s\n", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}

func helloWorldHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintln(w, "Hi, I'm the sample broker!")
}

func getCatalogHandler(w http.ResponseWriter, r *http.Request) {
	if status, err := checkCredentials(w, r); err != nil {
		w.WriteHeader(status)
		fmt.Fprintf(w, "Credentials check failed: %v", err)
		return
	}

	catalog := osbapi.Catalog{
		Services: []osbapi.Service{{
			Name:        "sample-service",
			Id:          "edfd6e50-aa59-4688-b5bf-b21e2ab27cdb",
			Description: "A sample service that does nothing",
			Plans: []osbapi.Plan{{
				Id:          "ebf1c1df-fefb-479b-9231-ddf700a37b58",
				Name:        "sample",
				Description: "Sample plan",
				Free:        true,
				Bindable:    true,
			}},
		}},
	}

	catalogBytes, err := json.Marshal(catalog)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to marshal catalog: %v", err)
		return
	}

	fmt.Fprintln(w, string(catalogBytes))
}

func checkCredentials(_ http.ResponseWriter, r *http.Request) (int, error) {
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) == 0 {
		return http.StatusUnauthorized, errors.New("Authorization request header not specified")
	}

	headerSplit := strings.Split(authHeader, " ")
	if len(headerSplit) != 2 {
		return http.StatusUnauthorized, errors.New("Could not parse Authorization request header")
	}

	if headerSplit[0] != "Basic" {
		return http.StatusUnauthorized, errors.New("Unsupported Authorization request header scheme. Only 'Basic' is supported")
	}

	credBytes, err := base64.StdEncoding.DecodeString(headerSplit[1])
	if err != nil {
		return http.StatusUnauthorized, errors.New("Failed to decode Authorization request header")
	}

	creds := strings.Split(string(credBytes), ":")
	if len(creds) != 2 {
		return http.StatusUnauthorized, errors.New("Failed to extract user credentials from Authorization request header")
	}

	username := creds[0]
	password := creds[1]

	if username != hardcodedUserName || password != hardcodedPassword {
		return http.StatusForbidden, fmt.Errorf("Incorrect credentials: user %q, password %q. Use %q as username and %q as password", username, password, hardcodedUserName, hardcodedPassword)
	}

	return -1, nil
}
