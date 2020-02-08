package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const offset = 10

type SQLHuman struct {
	db *sql.DB
}

func newSQLHuman() *SQLHuman {
	s := fmt.Sprintf("%v:%v@/%v", os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASS"), os.Getenv("MYSQL_DB"))
	database, err := sql.Open("mysql", s)

	if err != nil {
		log.Println(err)
	}

	return &SQLHuman{
		db: database,
	}
}

func (sqlhuman *SQLHuman) Add(w http.ResponseWriter, r *http.Request) {
	var human Human

	err := json.NewDecoder(r.Body).Decode(&human)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	human.ID = uuid.New().String()

	_, err = sqlhuman.db.Exec(
		"INSERT INTO User(ID,Firstname,Lastname,Age) VALUES(?,?,?,?)",
		human.ID, human.Firstname, human.Lastname, human.Age,
	)
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (sqlhuman *SQLHuman) GetAll(w http.ResponseWriter, r *http.Request) {
	var humanity []Human

	params := mux.Vars(r)
	page, err := strconv.Atoi(params["PAGE"])

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	page--

	offsetnum := page * offset

	rows, err := sqlhuman.db.Query("SELECT ID,Firstname,Lastname,Age FROM User LIMIT ?,?", offsetnum, offset)
	if err != nil {
		log.Println(err)
		return
	}

	for rows.Next() {
		var h Human
		err = rows.Scan(&h.ID, &h.Firstname, &h.Lastname, &h.Age)

		if err != nil {
			log.Println(err)
			continue
		}

		if err = rows.Err(); err != nil {
			log.Println(err)
			continue
		}

		humanity = append(humanity, h)
	}

	err = json.NewEncoder(w).Encode(humanity)
	if err != nil {
		log.Println(err)
	}
}

func (sqlhuman *SQLHuman) GetOne(w http.ResponseWriter, r *http.Request) {
	var h Human

	params := mux.Vars(r)

	row := sqlhuman.db.QueryRow("SELECT ID,Firstname,Lastname,Age FROM User WHERE ID=?", params["ID"])

	err := row.Scan(&h.ID, &h.Firstname, &h.Lastname, &h.Age)
	if err != nil {
		log.Println(err)
		return
	}

	err = json.NewEncoder(w).Encode(h)
	if err != nil {
		log.Println(err.Error())
	}
}

func (sqlhuman *SQLHuman) UpdateOne(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var h Human

	err := json.NewDecoder(r.Body).Decode(&h)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	h.ID = params["ID"]

	_, err = sqlhuman.db.Exec(
		"UPDATE User SET Firstname=?,Lastname=?,Age=? WHERE ID=?", h.Firstname, h.Lastname, h.Age, h.ID,
	)
	if err != nil {
		log.Println(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (sqlhuman *SQLHuman) DeleteOne(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	_, err := sqlhuman.db.Exec("UPDATE User SET Salted=1 WHERE ID=?", params["ID"])
	if err != nil {
		log.Println(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
