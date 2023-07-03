package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"

	_ "github.com/go-sql-driver/mysql"
)

// Define the structure of your JSON data
type PlanetData struct {
	ID     string `json:"_id"`
	Img    string `json:"img"`
	Name   string `json:"name"`
	System System `json:"system"`
	Type   string `json:"type"`
}

type System struct {
	BaseItem              interface{}           `json:"baseItem"`
	ContainerID           interface{}           `json:"containerId"`
	Description           Description           `json:"description"`
	EquippedBulk          EquippedBulk          `json:"equippedBulk"`
	Hardness              int                   `json:"hardness"`
	HP                    HP                    `json:"hp"`
	Level                 Level                 `json:"level"`
	NegateBulk            NegateBulk            `json:"negateBulk"`
	PreciousMaterial      PreciousMaterial      `json:"preciousMaterial"`
	PreciousMaterialGrade PreciousMaterialGrade `json:"preciousMaterialGrade"`
	Price                 Price                 `json:"price"`
	Quantity              int                   `json:"quantity"`
	Rules                 []interface{}         `json:"rules"`
	Size                  string                `json:"size"`
	Source                Source                `json:"source"`
	StackGroup            interface{}           `json:"stackGroup"`
	Traits                Traits                `json:"traits"`
	Usage                 Usage                 `json:"usage"`
	Weight                Weight                `json:"weight"`
}

type Description struct {
	Value string `json:"value"`
}

type EquippedBulk struct {
	Value string `json:"value"`
}

type HP struct {
	BrokenThreshold int `json:"brokenThreshold"`
	Max             int `json:"max"`
	Value           int `json:"value"`
}

type Level struct {
	Value int `json:"value"`
}

type NegateBulk struct {
	Value string `json:"value"`
}

type PreciousMaterial struct {
	Value string `json:"value"`
}

type PreciousMaterialGrade struct {
	Value string `json:"value"`
}

type Price struct {
	Value struct {
		GP int `json:"gp"`
	} `json:"value"`
}

type Source struct {
	Value string `json:"value"`
}

type Traits struct {
	Rarity string   `json:"rarity"`
	Value  []string `json:"value"`
}

type Usage struct {
	Value string `json:"value"`
}

type Weight struct {
	Value string `json:"value"`
}

var db *sql.DB

func main() {
	// Check if the directory path is provided as an argument
	if len(os.Args) < 2 {
		log.Fatal("Please provide the absolute path to the directory as an argument.")
	}

	// Get the directory path from command line argument
	dirPath := os.Args[1]

	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to load env", err)
	}

	// Open a connection to the database
	db, err = sql.Open("mysql", os.Getenv("DSN"))
	if err != nil {
		log.Fatal("failed to open db connection", err)
	}

	// Iterate over all files in the directory
	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		// Skip directories and process only JSON files
		if !info.IsDir() && filepath.Ext(path) == ".json" {
			// Read the JSON file
			file, err := ioutil.ReadFile(path)
			if err != nil {
				log.Println("Failed to read file:", path, "-", err)
				return nil
			}

			// Parse JSON data into a struct
			var planet PlanetData
			err = json.Unmarshal(file, &planet)
			if err != nil {
				log.Println("Failed to parse JSON:", path, "-", err)
				return nil
			}

			// // Insert the data into the MySQL database
			// _, err = db.Exec("INSERT INTO your_table (id, img, name, system_description, system_equipped_bulk, system_hardness, system_hp_max, system_level, system_price_gp, system_size, system_source, system_traits_rarity, system_usage_value, system_weight_value) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			// 	planet.ID, planet.Img, planet.Name, planet.System.Description.Value, planet.System.EquippedBulk.Value, planet.System.Hardness, planet.System.HP.Max, planet.System.Level.Value, planet.System.Price.Value.GP, planet.System.Size, planet.System.Source.Value, planet.System.Traits.Rarity, planet.System.Usage.Value, planet.System.Weight.Value)
			// if err != nil {
			// 	log.Println("Failed to insert data into database:", path, "-", err)
			// 	return nil
			// }

			//TODO
			//1. Check if item exists (by name!)
			// rows, err := db.Query(`SELECT * FROM items WHERE items.name="Abadar's Flawless Scale"`)
			// if err != nil {
			// 	log.Fatal("Failed to execute query:", err)
			// }
			// defer rows.Close()
			// fmt.Println(rows)

			//fmt.Println("Data inserted successfully from file:", path)
		}
		return nil
	})

	if err != nil {
		log.Fatal("Failed to iterate over files in directory:", err)
	}

	query := `SELECT * FROM items'`

	err = db.QueryRow(query).Scan()
	if err != nil {
		log.Fatal("fucked it up son", err)
	}
	// name := "Abadar's Flawless Scale"

	// checkForExisting(name)
}

// func checkForExisting(name) {
// 	query := "SELECT * FROM items"
// 	res, err := db.Query(query)
// 	defer res.Close()
// 	if err != nil {
// 		log.Fatal("fuckedUPSoon")
// 	}

// }
