// Uses the fairly standard and straigtforward Gorilla Mux
// go get -u github.com/gorilla/mux
//
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is my golang test home.")
}

func createEendpointsAndRun() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	if initDB() {
		log.Println("initialized memory model")
	} else {
		log.Fatal("Failed to initialize memory model")
	}
	// endpoints. Technically only asked for the first, the others allow for unit tests.
	//This is out API, and are RESTful
	router.HandleFunc("/user/register", createUser).Methods("POST")
	router.HandleFunc("/user/get", getUser).Methods("POST")
	router.HandleFunc("/user/getAll", getAllUsers).Methods("POST")
	router.HandleFunc("/user/update", updateUser).Methods("POST")
	router.HandleFunc("/user/delete", deleteUser).Methods("POST")
	router.HandleFunc("/user/deleteAll", deleteAllUsers).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
	releaseDB()
}
