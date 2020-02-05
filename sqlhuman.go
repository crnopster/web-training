package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type SQLHuman struct {
	db *sql.DB
}

func newSQLHuman() *SQLHuman {
	database, err := sql.Open("mysql", "root:1488@/mydb")
	if err != nil {
		log.Fatal(err.Error())
	}

	return &SQLHuman{
		db: database,
	}
}

func (sqlhuman *SQLHuman) Add(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")
	var human Human

	err := json.NewDecoder(r.Body).Decode(&human)
	if err != nil {
		log.Println(err.Error())
	}

	human.ID = uuid.New().String()
	prepared, err := sqlhuman.db.Prepare("INSERT INTO User(ID,Firstname,Lastname,Age) VALUES(?,?,?,?)")

	if err != nil {
		log.Println(err.Error())
	}

	_, err = prepared.Exec(human.ID, human.Firstname, human.Lastname, human.Age)
	if err != nil {
		log.Println(err.Error())
	}

	w.WriteHeader(http.StatusOK)
}

func (sqlhuman *SQLHuman) GetAll(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")
	var humanity []Human

	rows, err := sqlhuman.db.Query("SELECT ID,Firstname,Lastname,Age FROM User")
	if err != nil {
		log.Println(err.Error())
	}

	for rows.Next() {
		var h Human
		err = rows.Scan(&h.ID, &h.Firstname, &h.Lastname, &h.Age)

		if err != nil {
			log.Println(err.Error())
		}

		humanity = append(humanity, h)
	}

	err = json.NewEncoder(w).Encode(humanity)
	if err != nil {
		log.Println(err.Error())
	}
}

func (sqlhuman *SQLHuman) GetOne(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")
	var h Human

	params := mux.Vars(r)

	rows, err := sqlhuman.db.Query("SELECT ID,Firstname,Lastname,Age FROM User WHERE ID=?", params["ID"])

	if err != nil {
		log.Println(err.Error())
	}

	for rows.Next() {
		err = rows.Scan(&h.ID, &h.Firstname, &h.Lastname, &h.Age)
		if err != nil {
			log.Println(err.Error())
		}
	}

	err = json.NewEncoder(w).Encode(h)

	if err != nil {
		log.Println(err.Error())
	}
}

func (sqlhuman *SQLHuman) UpdateOne(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	var h Human

	err := json.NewDecoder(r.Body).Decode(&h)
	if err != nil {
		log.Println(err.Error())
	}

	h.ID = params["ID"]

	prepared, err := sqlhuman.db.Prepare("UPDATE User SET Firstname=?,Lastname=?,Age=? WHERE ID=?")
	if err != nil {
		log.Println(err.Error())
	}

	_, err = prepared.Exec(h.Firstname, h.Lastname, h.Age, h.ID)
	if err != nil {
		log.Println(err.Error())
	}

	w.WriteHeader(http.StatusOK)
}

func (sqlhuman *SQLHuman) DeleteOne(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	prepared, err := sqlhuman.db.Prepare("UPDATE User SET Salted=1 WHERE ID=?")
	if err != nil {
		log.Println(err.Error())
	}

	_, err = prepared.Exec(params["ID"])
	if err != nil {
		log.Println(err.Error())
	}

	w.WriteHeader(http.StatusOK)
}
