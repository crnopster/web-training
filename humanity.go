package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	jsoniter "github.com/json-iterator/go"
)

type Human struct {
	ID      string `json:"ID"`
	Name    string `json:"Name"`
	Surname string `json:"Surname"`
	Age     int    `json:"Age"`
}

type Humanity struct {
	humans []Human
}

func (hm *Humanity) AddUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var human Human

	errcheck(jsoniter.NewDecoder(r.Body).Decode(&human))
	human.ID = uuid.New().String()
	hm.humans = append(hm.humans, human)
	errcheck(jsoniter.NewEncoder(w).Encode(&human))
}

func (hm *Humanity) GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	errcheck(jsoniter.NewEncoder(w).Encode(hm.humans))
}

func (hm *Humanity) GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	for _, human := range hm.humans {
		if usercheck(human.ID, params["ID"]) {
			errcheck(jsoniter.NewEncoder(w).Encode(human))
			return
		}
	}
}

func (hm *Humanity) UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	for index, human := range hm.humans {
		if usercheck(human.ID, params["ID"]) {
			var h Human

			errcheck(jsoniter.NewDecoder(r.Body).Decode(&h))
			h.ID = human.ID
			errcheck(jsoniter.NewEncoder(w).Encode(h))
			hm.humans[index] = h
		}
	}
}

func (hm *Humanity) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	for index, human := range hm.humans {
		if usercheck(human.ID, params["ID"]) {
			hm.humans = append(hm.humans[:index], hm.humans[index+1:]...)
			break
		}
	}
}
