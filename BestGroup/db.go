package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	"database/sql"

	_ "github.com/lib/pq"
)

const (
	host     = "127.0.0.1"
	port     = 5432
	puser    = "postgres"
	password = "#PostgresDatabase"
	dbname   = "Schedule"
)

var (
	db *sql.DB
)

func NewDB() {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, puser, password, dbname)

	var err error
	// Connect to database
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
}

func QueryDB(query string) string {
	var rowf string
	sqlStatement := `SELECT tasks FROM test WHERE id=$1;`
	row, err := db.Query(sqlStatement, 1)
	if err != nil {
		loganswer("Error in query")
		return "Error in query"
	}
	defer row.Close()
	for row.Next() {
		switch err := row.Scan(&rowf); err {
		case sql.ErrNoRows:
			loganswer("No rows were returned!")
		case nil:
			return rowf
		default:
			loganswer("Error in query : " + err.Error())
			return "Error in query"
		}
	}
	return ""
}

func CreateTable(dbName string, dbusers string) string {
	name, extkey := Encrypt(dbName)
	statement := `
	INSERT INTO schedules (schedule, starting_date, ending_date, created_on, additional_key) VALUES($1,$2,$3,$4,$5)`
	// sT, err := time.Parse("", "")
	// if err != nil {
	// 	loganswer("Wrong time typo starting time")
	// 	return "Wrong time typo starting time"
	// }
	// eT, err := time.Parse("", "")
	// if err != nil {
	// 	loganswer("Wrong time typo ending time")
	// 	return "Wrong time typo ending time"
	// }
	db.Query(statement, name, time.Now(), time.Now(), time.Now(), extkey)
	statement = "INSERT INTO schedules_names (schedule) VALUES($1)"
	db.Query(statement, dbName)
	//reduce name to 62 characters and remove 10 characters from the beginning
	name = name[10:72]
	statement = "CREATE TABLE " + name + " (id SERIAL PRIMARY KEY, tasks json NOT NULL, roles json NOT NULL, reports JSON NOT NULL);"
	_, err := db.Query(statement)
	if err != nil {
		println(err.Error())
		loganswer("Error in query")
		return "Error in query"
	}
	statement = "CREATE TABLE " + name + "_users (id SERIAL PRIMARY KEY, username VARCHAR (100) UNIQUE NOT NULL, created_on TIMESTAMP NOT NULL, last_login TIMESTAMP NOT NULL);"
	db.Query(statement)
	statement = "INSERT INTO " + name + "_users (username, created_on, last_login) VALUES "
	splited := strings.Split(dbusers, ",")
	for numb := range splited {
		statement += "( $" + strconv.FormatInt(int64(numb), 32) + "," + time.Now().String() + "," + time.Now().String() + ")"
	}
	db.Query(statement, splited)

	return ""
}

func GetTable(name string, user string) string {
	if !CheckExist(name) {
		return "Schedule does not exist"
	}
	statement := "SELECT schedule, additional_key FROM schedules"
	rows, err := db.Query(statement)
	if err != nil {
		loganswer("Error in query")
		return "Error in query"
	}
	defer rows.Close()
	var schedule string
	var key string
	for rows.Next() {
		err = rows.Scan(&schedule, &key)
		if len(key)+len(name) == 32 {
			final, finalkey := Encrypt(name)
			if final == schedule && finalkey == key {
				break
			}
		}
		if err != nil {
			loganswer("no schedule found")
			return "no schedule found"
		}
	}
	return schedule
}

func CheckExist(name string) bool {
	statement := "SELECT * FROM schedules_names WHERE schedule=$1"
	rows, err := db.Query(statement, name)
	if err != nil {
		loganswer("Error in query")
		return false
	}
	defer rows.Close()
	var schedule string
	for rows.Next() {
		err = rows.Scan(&schedule)
		if err != nil {
			loganswer("Error in query")
			return false
		}
	}
	return true
}

func Encrypt(txt string) (string, string) {
	var additionalKey []byte
	if len(txt) < 32 {
		additionalKey = []byte(RandStringBytesRmndr(32 - len(txt)))
	} else if len(txt) > 32 {
		txt = txt[:32]
	}
	text := []byte(txt)
	key := text
	key = append(key, additionalKey...)

	// generate a new aes cipher using our 32 byte long key
	c, err := aes.NewCipher(key)
	// if there are any errors, handle them
	if err != nil {
		loganswer(err.Error())
	}

	// gcm or Galois/Counter Mode, is a mode of operation
	// for symmetric key cryptographic block ciphers
	// - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	gcm, err := cipher.NewGCM(c)
	// if any error generating new GCM
	// handle them
	if err != nil {
		loganswer(err.Error())
	}

	// creates a new byte array the size of the nonce
	// which must be passed to Seal
	nonce := make([]byte, gcm.NonceSize())
	// populates our nonce with a cryptographically secure
	// random sequence
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		loganswer(err.Error())
	}
	// here we encrypt our text using the Seal function
	// Seal encrypts and authenticates plaintext, authenticates the
	// additional data and appends the result to dst, returning the updated
	// slice. The nonce must be NonceSize() bytes long and unique for all
	// time, for a given key.
	return fmt.Sprintf("%x", gcm.Seal(nonce, nonce, text, nil)), string(additionalKey)
}