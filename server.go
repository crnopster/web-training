package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type human struct {
	ID      string
	Name    string
	Surname string
}

type humanity struct {
	humans []human
}

func (hm *humanity) Add(w http.ResponseWriter, r *http.Request) {
	h := human{
		ID:      uuid.New().String(),
		Name:    r.FormValue("Name"),
		Surname: r.FormValue("Surname"),
	}
	hm.humans = append(hm.humans, h)
	log.Println("Add ", h)
}

func (hm *humanity) Read(w http.ResponseWriter, r *http.Request) {
	for _, h := range hm.humans {
		fmt.Fprintf(w, "Name: %v, Surname: %v with ID: %v\n", h.Name, h.Surname, h.ID)
		log.Println("Read ", h)
	}
}

func (hm *humanity) Get(w http.ResponseWriter, r *http.Request) {
	m := mux.Vars(r)
	for _, h := range hm.humans {
		if h.ID != m["ID"] {
			continue
		}

		fmt.Fprintf(w, "Name: %v, Surname: %v\n", h.Name, h.Surname)
		log.Println("Get ", h)
	}
}

func (hm *humanity) Update(w http.ResponseWriter, r *http.Request) {
	m := mux.Vars(r)
	for _, h := range hm.humans {
		if h.ID != m["ID"] {
			continue
		}

		h.Name = r.FormValue("Name")
		h.Surname = r.FormValue("Surname")
		fmt.Fprintf(w, "Update %v with %v, %v", h.ID, h.Name, h.Surname)
	}
}

func (hm *humanity) Delete(w http.ResponseWriter, r *http.Request) {
	m := mux.Vars(r)
	for n, h := range hm.humans {
		if h.ID != m["ID"] {
			continue
		}

		hm.humans = append(hm.humans[:n], hm.humans[n+1:]...)
	}
}

func main() {
	addr := flag.String("a", ":8080", "address of app")
	flag.Parse()

	hm := new(humanity)

	mainRoute := mux.NewRouter()

	apiRoute := mainRoute.PathPrefix("/api/v1").Subrouter()

	// Create
	apiRoute.HandleFunc("/list", hm.Add).Methods("POST")
	// READ
	apiRoute.HandleFunc("/list", hm.Read).Methods("GET")
	// GET
	apiRoute.HandleFunc("/list/{ID}", hm.Get).Methods("GET")
	// UPDATE
	apiRoute.HandleFunc("/list/{ID}", hm.Update).Methods("PUT")
	// DELETE
	apiRoute.HandleFunc("/list/{ID}", hm.Delete).Methods("DELETE")

	http.Handle("/", apiRoute)

	log.Println("server started")

	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err.Error())
	}
}
