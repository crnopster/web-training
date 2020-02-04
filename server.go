package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	addr := flag.String("a", ":80", "address of app")
	flag.Parse()

	var hm Humanity

	mainRoute := mux.NewRouter()
	apiRoute := mainRoute.PathPrefix("/api/v1").Subrouter()
	apiRoute.HandleFunc("/list", hm.GetUsers).Methods("GET")
	apiRoute.HandleFunc("/list", hm.AddUser).Methods("POST")
	apiRoute.HandleFunc("/list/{ID}", hm.GetUser).Methods("GET")
	apiRoute.HandleFunc("/list/{ID}", hm.UpdateUser).Methods("PUT")
	apiRoute.HandleFunc("/list/{ID}", hm.DeleteUser).Methods("DELETE")

	if err := http.ListenAndServe(*addr, mainRoute); err != nil {
		log.Fatal(err.Error())
	}
}
