package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"

	_ "github.com/go-sql-driver/mysql"
)

// // Define the structure of your JSON data
type items struct {
	Img       string `json:"img"`
	Name      string `json:"name"`
	System_id int
	System    System `json:"system"`
	Type      string `json:"type"`
	item_id   int
	FinalGP   float64
}

type System struct {
	Description           Description           `json:"description"`
	Category              string                `json: "category"`
	Level_type            Level                 `json:"level"`
	PreciousMaterial      PreciousMaterial      `json:"preciousMaterial"`
	PreciousMaterialGrade PreciousMaterialGrade `json:"preciousMaterialGrade"`
	Price                 Price                 `json:"price"`
	Rules                 []interface{}         `json:"rules"`
	Source                Source                `json:"source"`
	Traits                Traits                `json:"traits"`
	Weight                Weight                `json:"weight"`
}

type Description struct {
	Value string `json:"value"`
}

type Level struct {
	Value int `json:"value"`
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
		SP int `json:"sp"`
		CP int `json:"cp"`
		PP int `json:"pp"`
	} `json:"value"`
}

type Source struct {
	Value string `json:"value"`
}

type Traits struct {
	Rarity     string   `json:"rarity"`
	Trait_list []string `json:"value"`
}

type Weight struct {
	Value string `json:"value"`
}

var db *sql.DB
var err error

func main() {
	//Check if the directory path is provided as an argument
	if len(os.Args) < 2 {
		log.Fatal("Please provide the absolute path to the directory as an argument.")
	}

	file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	defer file.Close()

	// Set the log output to the file
	log.SetOutput(file)

	// Get the directory path from command line argument
	dirPath := os.Args[1]

	err = godotenv.Load()
	if err != nil {
		log.Fatal("failed to load env", err)
	}

	db, err = sql.Open("mysql", os.Getenv("DSN"))
	if err != nil {
		log.Fatal("failed to open db connection", err)
	}

	//Start of Big Loop
	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(path) == ".json" {
			file, err := ioutil.ReadFile(path)
			if err != nil {
				log.Println("Failed to read file:", path, "-", err)
				return nil
			}

			// Parse JSON data into a struct
			var item items
			err = json.Unmarshal(file, &item)
			if err != nil {
				log.Println("Failed to parse JSON:", path, "-", err)
				return nil
			}
			var exists bool = checkForExisting(item.Name)
			if exists == false && item.Type != "kit" {

				fmt.Println(path)
				item.FinalGP = convertCurrencyToGP(item.System.Price.Value.GP, item.System.Price.Value.SP, item.System.Price.Value.CP, item.System.Price.Value.PP)
				item.System_id = getNextSystemId()

				writeItem(item.Img, item.Name, item.System_id, item.Type)

				writeSystem(item.System_id, item.System.Description.Value, item.System.Weight.Value, item.System.Level_type.Value, item.System.PreciousMaterial.Value, item.System.PreciousMaterialGrade.Value, item.FinalGP, item.System.Source.Value, item.System.Traits.Rarity)

				item.item_id = getItemId(item.Name)

				//Traits Loop
				for i := 0; i < len(item.System.Traits.Trait_list); i++ {
					fmt.Println(item.System.Traits.Trait_list[i])
					query := "select trait_id from traits where trait_name = ?"
					var trait_id int
					err = db.QueryRow(query, item.System.Traits.Trait_list[i]).Scan(&trait_id)

					if err != nil {
						log.Printf("Trait not found: Writing new Trait Id %v", item.System.Traits.Trait_list[i])

						_, err := db.Exec("INSERT INTO traits (trait_name) VALUES (?)", item.System.Traits.Trait_list[i])
						if err != nil {
							log.Println("Failed to insert Trait: ", item.System.Traits.Trait_list[i], "-", err)
							//IF writing of new trait was succesful, get trait_id
						} else {
							err = db.QueryRow(query, item.System.Traits.Trait_list[i]).Scan(&trait_id)
							if err != nil {
								log.Fatal("Failed to query recently inserted trait... This shouldn't happen ", item.System.Traits.Trait_list[i], "-", err)
								os.Exit(1)
							}
						}
					}
					//Regardless of path, write item trait id
					writeTraitsItem(item.item_id, trait_id)

				}
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal("Failed to iterate over files in directory:", err)
	}
}

func printItem(item items) {
	fmt.Println("item.item_id: " + string(item.item_id))
	fmt.Println("item.system.description: " + string(item.System.Description.Value))
	fmt.Println("Bulk: " + item.System.Weight.Value)
	fmt.Println(fmt.Sprintf("item.system.level.level: %v", item.System.Level_type.Value))
	fmt.Println("item.preciousMaterial: " + (item.System.PreciousMaterial.Value))
	fmt.Println("item.preciousmaterialGrade: " + string(item.System.PreciousMaterialGrade.Value))
	fmt.Println(fmt.Sprintf("GP final: %v", item.FinalGP))
	fmt.Println("SourceBook: " + item.System.Source.Value)
	fmt.Println("Rarity: " + item.System.Traits.Rarity)
	fmt.Println("item.name: " + string(item.Name))
	fmt.Println("item.img: " + string(item.Img))
	fmt.Println("item.Type: " + string(item.Type))
	fmt.Println("system_id: " + string(item.System_id))
	fmt.Println("Writing Item")

}

func writeTraitsItem(item_id int, trait_id int) {
	_, err := db.Exec("INSERT INTO item_traits (item_id, trait_id) VALUES (?,?)", item_id, trait_id)
	if err != nil {
		log.Fatal("We fucked up writing the traits you dumb fuck")
	}
}

func writeItem(Img string, Name string, System_id int, Type string) {
	_, err := db.Exec("INSERT INTO items (img, name, system_id, type) VALUES (?,?,?,?)", Img, Name, System_id, Type)
	if err != nil {
		log.Println("Failed to insert Item: ", Name, "-", err)
	}
}

func writeSystem(
	System_id int,
	Description string,
	bulk string,
	item_level int,
	PreciousMaterial string,
	PreciousMaterialGrade string,
	FinalGP float64,
	SourceBook string,
	Rarity string) {
	_, err := db.Exec("INSERT INTO system (system_id, description_value, bulk, level_value, precious_material_value, precious_material_grade_value, price_gp_value, source_book, rarity_value) VALUES (?,?,?,?,?,?,?,?,?)", System_id, Description, bulk, item_level, PreciousMaterial, PreciousMaterialGrade, FinalGP, SourceBook, Rarity)
	if err != nil {
		log.Println("Failed to insert System: ", System_id, "-", err)
	}
}

func getItemId(name string) int {
	query := `select item_id FROM items where name = ?`
	var item_id int
	err = db.QueryRow(query, name).Scan(&item_id)
	if err != nil {
		log.Fatal("WE FUCKED UP SON")
		os.Exit(1)
	}
	return (item_id)
}

func checkForExisting(name string) bool {
	query := `SELECT item_id FROM items WHERE items.name = ?`
	item_id := ""
	err = db.QueryRow(query, name).Scan(&item_id)
	if err != nil {
		return false
	}
	fmt.Println(item_id)

	return true
}

// SP = 1 = 0.1
//CP = 1 = .01
// PP = 1
//GP = 10.11
func convertCurrencyToGP(GP int, SP int, CP int, PP int) float64 {
	fmt.Println()
	var PP2GP float64 = float64(PP) * 10
	var SP2GP float64 = float64(SP) / 10
	var CP2GP float64 = float64(CP) / 100

	return float64(GP) + PP2GP + SP2GP + CP2GP
}

func getNextSystemId() int {
	query := `select MAX(system_id) from system`
	var system_id int
	err = db.QueryRow(query).Scan(&system_id)
	if err != nil {
		system_id = 0
	}
	return system_id + 1
}
