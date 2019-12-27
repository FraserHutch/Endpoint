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
	// Endpoints. Technically only asked for the first, the others allow for unit tests.
	// Note that these could all share the base user endpoint - to differentiate between
	// get/delete and get all/delete all I could have specified that distinction in the json.
	// However, this to me is cleaner, and allows me to easily decouple from http.
	router.HandleFunc("/user/register", createUser).Methods("POST")
	router.HandleFunc("/user/get", getUser).Methods("GET")
	router.HandleFunc("/user/getAll", getAllUsers).Methods("GET")
	router.HandleFunc("/user/update", updateUser).Methods("PUT") // does NOT create if record not found
	router.HandleFunc("/user/delete", deleteUser).Methods("DELETE")
	router.HandleFunc("/user/deleteAll", deleteAllUsers).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
	releaseDB()
}
