package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "developer"
	dbname   = "phone_exercise"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("failed opening connetion to db: %s", err)
	}

	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatalf("connection ping to db failed: %s", err)
	}

	seed(db)
	normalizeDBNumbers(db)
}

func normalizeDBNumbers(db *sql.DB) {
	rows, err := db.Query("select id, number from numbers")
	if err != nil {
		log.Fatalf("failed to get numbers from db: %s", err)
	}

	for rows.Next() {
		var id uint
		var num string
		err := rows.Scan(&id, &num)
		if err != nil {
			log.Fatalf("failed to get next number from db: %s", err)
		}

		normalizedNum := normalizeNumber(num)
		updateNormalizedRecord(db, id, normalizedNum)
	}
}

func updateNormalizedRecord(db *sql.DB, id uint, normalizedNum string) {
	duplicateNums, err := db.Query("select id from numbers where number = $1", normalizedNum)
	if err != nil {
		log.Fatalf("failed to query table for duplicate numbers: %s", err)
	}

	if duplicateNums.Next() {
		_, err = db.Exec("delete from numbers where id = $1", id)
		if err != nil {
			log.Fatalf("failed to delete duplicate number: %s", err)
		}
	} else {
		_, err = db.Exec("update numbers set number = $1 where id = $2", normalizedNum, id)
		if err != nil {
			log.Fatalf("failed normalizing number with id %d: %s", id, err)
		}
	}
}

func normalizeNumber(num string) string {
	var res string
	for _, ch := range num {
		if ch >= '0' && ch <= '9' {
			res += string(ch)
		}
	}

	return res
}

func seed(db *sql.DB) {
	resetTable(db)
	insertInitialData(db)
}

func resetTable(db *sql.DB) {
	_, err := db.Exec("SELECT 1 FROM numbers LIMIT 1")
	if err == nil {
		_, err := db.Exec("drop table numbers")
		if err != nil {
			log.Fatalf("failed to recreate table: %s", err)
		}
	}

	_, err = db.Exec("create table numbers (id serial primary key, number VARCHAR(50) not null)")
	if err != nil {
		log.Fatalf("failed to create new table: %s", err)
	}
}

func insertInitialData(db *sql.DB) {
	nums := []string{
		"1234567890",
		"123 456 7891",
		"(123) 456 7892",
		"(123) 456-7893",
		"123-456-7894",
		"123-456-7890",
		"1234567892",
		"(123)456-7892",
	}

	for _, v := range nums {
		_, err := db.Exec("insert into numbers (number) values ($1)", v)
		if err != nil {
			log.Fatalf("failed to insert %s into db: %s", v, err)
		}
	}
}
