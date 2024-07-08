package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/lib/pq"
)

func main() {
	// PostgreSQL connection parameters
	const (
		host     = "localhost" // or the Docker container IP
		port     = 5432        // or the port your Docker container is exposed on
		user     = "postgres"
		password = "postgres"
		dbname   = "mediascraper"
	)

	// Connection string
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Connect to PostgreSQL database
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Ensure connection is valid
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to the database!")
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS url_metadata (
		id SERIAL PRIMARY KEY,
		real_url TEXT NOT NULL,
		stored_url TEXT NOT NULL,
		metadata JSONB,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`)

	if err != nil {
		panic(err)
	}
	fmt.Println("Table 'url_metadata' created or already exists.")
}
