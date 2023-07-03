package main


func main() {
    // Load in the `.env` file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to load env", err)
	}
}