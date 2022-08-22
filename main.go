package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Bird struct {
	Species     string
	Description string
}

func main() {

	db, err := sql.Open("pgx", "postgresql://localhost:5432/bird_encyclopedia")

	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	db.SetMaxIdleConns(5)

	db.SetMaxOpenConns(10)

	db.SetConnMaxIdleTime(1 * time.Second)

	db.SetConnMaxLifetime(30 * time.Second)

	if err := db.Ping(); err != nil {
		log.Fatalf("Unable to reach database: %v", err)
	}
	fmt.Println("datbase is reachable")

	birdName := "eagle"

	row := db.QueryRow("SELECT bird, description FROM birds WHERE bird = $1 LIMIT $2", birdName, 1)

	bird := Bird{}

	if err := row.Scan(&bird.Species, &bird.Description); err != nil {
		log.Fatalf("COuld not scan row: %v", err)
	}
	fmt.Printf("found bird")

	rows, err := db.Query("SELECT bird, description FROM birds limit 10")
	if err != nil {
		log.Fatalf("could not execute query: %v", err)
	}

	birds := []Bird{}

	for rows.Next() {
		bird := Bird{}
		if err := rows.Scan(&bird.Species, &bird.Description); err != nil {
			log.Fatalf("could not scan now: %v", err)
		}
		birds = append(birds, bird)
	}

	fmt.Printf("found %d birds: %+v\n", len(birds), birds)

	_, err = db.Exec("DELETE FROM birds WHERE bird=$1", "rooster")
	if err != nil {
		log.Fatalf("could not delete row: %v", err)
	}

	newBird := Bird{
		Species:     "rooster",
		Description: "Wakes you up in the morning",
	}

	result, err := db.Exec("INSERT INTO birds (bird, description) VALUES ($1, $2)",
		newBird.Species, newBird.Description)

	if err != nil {
		log.Fatalf("could not insert row: %v", err)
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		log.Fatalf("could not get affected rows: %v", err)
	}

	fmt.Println("inserted", rowsAffected, "rows")

	ctx := context.Background()

	ctx, _ = context.WithTimeout(ctx, 300*time.Millisecond)

	_, err = db.QueryContext(ctx, "SELECT * from pg_sleep(1)")

	if err != nil {
		log.Fatalf("could not execute query: %v", err)
	}
}
