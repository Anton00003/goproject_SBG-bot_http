package data

import (
	"database/sql"
	//	"encoding/csv"
	"fmt"
	//	"io/ioutil"
	//	"os"
	//	s "strings"

	_ "github.com/lib/pq"
)

type DatabaseReader struct {
	Database *sql.DB
	Url      string
}

type Reader struct {
	filePath string
}

const (
	database_url = "host=localhost port=5431 user=db password=db dbname = db sslmode=disable"

// file         = "tmp/Birthbay.csv"
)

func NewDatabaseReader() *DatabaseReader {
	var database sql.DB
	return &DatabaseReader{Database: &database,
		Url: database_url}
}

//func NewReader() *Reader {
//	return &Reader{filePath: file}
//}

func (r *DatabaseReader) ReadDatabase() ([][]string, [][]byte, error) {
	b := make([][]byte, 0, 10)
	records := make([][]string, 0, 10)

	var err error
	r.Database, err = sql.Open("postgres", r.Url)
	if err != nil {
		fmt.Println("Ошибка открытия")
		fmt.Println(err)
		return records, b, err
	}
	fmt.Println("Успешное открытие")

	if err := r.Database.Ping(); err != nil {
		fmt.Println("Ошибка пинга")
		fmt.Println(err)
		return records, b, err
	}
	fmt.Println("Пинг прошел успешно")

	// Вывод из таблицы

	rows, _ := r.Database.Query("SELECT name, id, Date, subscribers FROM birthday")
	var name string
	var id string
	var date string
	var subscribers []byte

	for rows.Next() {
		record := make([]string, 3)
		rows.Scan(&name, &id, &date, &subscribers)
		//		fmt.Printf("%s %s %s %s\n", name, id, date, subscribers)
		//		fmt.Println(name)
		record[0] = name
		record[1] = id
		record[2] = date

		records = append(records, record)
		b = append(b, subscribers)
	}
	fmt.Println(records, b)
	return records, b, nil
}

func (r *DatabaseReader) WriteDatabase(record []string, b []byte) error {

	var err error
	r.Database, err = sql.Open("postgres", r.Url)
	if err != nil {
		fmt.Println("Ошибка открытия")
		fmt.Println(err)
		return err
	}
	fmt.Println("Успешное открытие")

	statement, err1 := r.Database.Prepare("INSERT INTO birthday (name, id, Date, Subscribers) VALUES ($1, $2, $3, $4)")
	if err1 != nil {
		fmt.Println("Ошибка записи в таблицу 1")
		fmt.Println(err1)
		return err1
	}
	//	statement.Exec("Aa", "1234", "2017-06-20", "Bb")
	statement.Exec(record[0], record[1], record[2], b)

	//+++++ Проверка ++++++++++++++++++++++++++++
	rows, _ := r.Database.Query("SELECT name, id, Date, subscribers FROM birthday")
	var name string
	var id string
	var date string
	var subscribers []byte

	for rows.Next() {
		rows.Scan(&name, &id, &date, &subscribers)
		fmt.Printf("%s %s %s %s\n", name, id, date, subscribers)
	}
	//+++++ Проверка ++++++++++++++++++++++++++++

	return nil
}

func (r *DatabaseReader) UpdateDatabase(nameUpdate string, record []byte) error {

	var err error
	r.Database, err = sql.Open("postgres", r.Url)
	if err != nil {
		fmt.Println("Ошибка открытия")
		fmt.Println(err)
		return err
	}
	fmt.Println("Успешное открытие")

	statement, err1 := r.Database.Prepare("UPDATE birthday SET Subscribers = $1 WHERE name = $2")
	if err1 != nil {
		fmt.Println("Ошибка удаления из таблицы")
		fmt.Println(err1)
		return err1
	}
	statement.Exec(record, nameUpdate)

	//+++++ Проверка ++++++++++++++++++++++++++++
	rows, _ := r.Database.Query("SELECT name, id, Date, subscribers FROM birthday")
	var name string
	var id string
	var date string
	var subscribers []byte

	for rows.Next() {
		rows.Scan(&name, &id, &date, &subscribers)
		fmt.Printf("%s %s %s %s\n", name, id, date, subscribers)
	}
	//+++++ Проверка ++++++++++++++++++++++++++++

	return nil
}
