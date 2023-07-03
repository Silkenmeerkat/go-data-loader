package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var db *sql.DB

func main() {
	// Load in the `.env` file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to load env", err)
	}
	db, err = sql.Open("mysql", os.Getenv("DSN"))
	if err != nil {
		log.Fatal("failed to open db connection", err)
	}
	// GitHub repository URL
	repoURL := "https://api.github.com/repos/foundryvtt/pf2e/contents/packs/equipment"

	// Make the initial request to get the first page
	resp, err := http.Get(repoURL)
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		return
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Unexpected response status code:", resp.StatusCode)
		return
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Parse the response JSON
	var files []struct {
		ID     string `json:"_id"`
		Img    string `json:"img"`
		Name   string `json:"name"`
		System struct {
			BaseItem    interface{} `json:"baseItem"`
			ContainerID interface{} `json:"containerId"`
			Description struct {
				Value string `json:"value"`
			} `json:"description"`
			EquippedBulk struct {
				Value string `json:"value"`
			} `json:"equippedBulk"`
			Hardness int `json:"hardness"`
			HP       struct {
				BrokenThreshold int `json:"brokenThreshold"`
				Max             int `json:"max"`
				Value           int `json:"value"`
			} `json:"hp"`
			Level struct {
				Value int `json:"value"`
			} `json:"level"`
			NegateBulk struct {
				Value string `json:"value"`
			} `json:"negateBulk"`
			PreciousMaterial struct {
				Value string `json:"value"`
			} `json:"preciousMaterial"`
			PreciousMaterialGrade struct {
				Value string `json:"value"`
			} `json:"preciousMaterialGrade"`
			Price struct {
				Value struct {
					GP int `json:"gp"`
				} `json:"value"`
			} `json:"price"`
			Quantity int           `json:"quantity"`
			Rules    []interface{} `json:"rules"`
			Size     string        `json:"size"`
			Source   struct {
				Value string `json:"value"`
			} `json:"source"`
			StackGroup interface{} `json:"stackGroup"`
			Traits     struct {
				Rarity string   `json:"rarity"`
				Value  []string `json:"value"`
			} `json:"traits"`
			Usage struct {
				Value string `json:"value"`
			} `json:"usage"`
			Weight struct {
				Value string `json:"value"`
			} `json:"weight"`
		} `json:"system"`
		Type string `json:"type"`
	}

	if err := json.Unmarshal(body, &files); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// Loop over the files
	for _, file := range files {
		// Process each file here
		fmt.Println("File:", file.Name)
	}

	// Check if there are more pages
	linkHeader := resp.Header.Get("Link")
	if linkHeader != "" {
		// Extract the next page URL from the Link header
		nextPageURL := getNextPageURL(linkHeader)

		// Make requests for the subsequent pages
		for nextPageURL != "" {
			resp, err := http.Get(nextPageURL)
			if err != nil {
				fmt.Println("Error making HTTP request:", err)
				return
			}
			defer resp.Body.Close()

			// Read the response body
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error reading response body:", err)
				return
			}

			// Parse the response JSON
			if err := json.Unmarshal(body, &files); err != nil {
				fmt.Println("Error parsing JSON:", err)
				return
			}

			// Loop over the files in the subsequent page
			for _, file := range files {
				fmt.Println("Name: ", file.Name)
				// // Process each file here
				// fmt.Println("File:", file.Name)
				// fmt.Println("URL:", file.URL)
			}

			// Extract the next page URL for the next iteration
			nextPageURL = getNextPageURL(resp.Header.Get("Link"))
		}
	}
}

func getNextPageURL(linkHeader string) string {
	// Extracts the next page URL from the Link header
	links := parseLinkHeader(linkHeader)
	if links["next"] != "" {
		return links["next"]
	}
	return ""
}

func parseLinkHeader(linkHeader string) map[string]string {
	// Parses the Link header and returns a map of URLs
	links := make(map[string]string)
	entries := strings.Split(linkHeader, ",")
	for _, entry := range entries {
		parts := strings.Split(strings.TrimSpace(entry), ";")
		if len(parts) < 2 {
			continue
		}
		url := strings.Trim(parts[0], "<>")
		rel := strings.Trim(parts[1], " ")
		rel = strings.TrimPrefix(rel, "rel=")
		rel = strings.Trim(rel, "\"")
		links[rel] = url
	}
	return links
}
