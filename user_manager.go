package main

// This is the user manager - it roughly corresponds to the view part of MVP

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// UserOperationResult  - rrequest and return block for create and update user operations
type UserOperationResult struct {
	Status string `json:"Status"`
	Reason string `json:"Reason"`
	User   User   `json:"User"`
}

// UserGetAllOperationResult - response object for GetAll users
type UserGetAllOperationResult struct {
	Status string `json:"Status"`
	Reason string `json:"Reason"`
	Count  int    `json:"Count"`
	Users  []User `json:"Users"`
}

// SimpleOperationResult  - request and return block for create and update user operations
type SimpleOperationResult struct {
	Status string `json:"Status"`
	Reason string `json:"Reason"`
}

// UserNameOperation  - request  block for get and delete user operations
type UserNameOperation struct {
	UserName string `json:"UserName"`
}

//// HANDLERS - these correcpond one to one with the API declared in endpoint.go
func createUser(w http.ResponseWriter, r *http.Request) {
	log.Println("createUser(): invoked")
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Invalid data - expected Username, Email, and Password for new user")
	}

	// gget user data from json. Could use error handling.
	var newUser User
	json.Unmarshal(reqBody, &newUser)
	log.Printf("createUser(): request data: %v", newUser)
	// now update the db.
	var result UserOperationResult
	var httpStatus int
	var retCode ModelStatusCode
	result.User, retCode, result.Reason = modelCreateUser(newUser)
	result.Status = ModelStatusText(retCode)

	// handle response.'
	switch retCode {
	case ModelSuccess:
		httpStatus = http.StatusOK
	case ModelDBCreateFailure:
		httpStatus = http.StatusInternalServerError
	}

	log.Printf("createUser(): returning %v -> %v", httpStatus, result)
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(result)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	log.Println("updateUser(): invoked")
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Invalid data - expected Username, Email, and Password for new user")
	}

	// gget user data from json. Could use error handling.
	var user User
	json.Unmarshal(reqBody, &user)
	log.Printf("updateUser(): request data: %v", user)
	// now update the db.
	var result UserOperationResult
	var httpStatus int
	var retCode ModelStatusCode
	result.User, retCode, result.Reason = modelUpdateUser(user)
	result.Status = ModelStatusText(retCode)

	// handle response.'
	switch retCode {
	case ModelSuccess:
		httpStatus = http.StatusOK
	case ModelDBUpdateFailure:
		httpStatus = http.StatusInternalServerError
	default:
		httpStatus = http.StatusInternalServerError
	}

	log.Printf("updateUser(): returning %v -> %v", httpStatus, result)
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(result)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	log.Println("getUser(): invoked")
	var result UserOperationResult
	var httpStatus int

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Invalid data - expected Username")
	}

	// gget user data from json. Could use error handling.
	var userNameOp UserNameOperation
	json.Unmarshal(reqBody, &userNameOp)
	log.Printf("getUser(): request data: %v", userNameOp)
	// now retrieve from our db.
	var retCode ModelStatusCode
	result.User, retCode, result.Reason = modelGetUser(userNameOp.UserName)
	result.Status = ModelStatusText(retCode)

	// handle response.'
	switch retCode {
	case ModelSuccess:
		httpStatus = http.StatusOK
	case ModelDBGetFailure:
		httpStatus = http.StatusInternalServerError
	}

	log.Printf("getUser(): returning %v -> %v", httpStatus, result)
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(result)
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	log.Println("getAllUsers() invoked")
	var result UserGetAllOperationResult
	var httpStatus int

	var retCode ModelStatusCode
	result.Users, retCode, result.Reason = modelGetAllUsers()
	result.Status = ModelStatusText(retCode)
	result.Count = len(result.Users)

	// handle response.'
	switch retCode {
	case ModelSuccess:
		log.Println("getAllUsers(): found")
		httpStatus = http.StatusOK
	case ModelDBCreateFailure:
		log.Println("getAllUsers(): server error 1")
		httpStatus = http.StatusInternalServerError
	default:
		log.Println("getAllUsers(): server error 2")
		httpStatus = http.StatusInternalServerError // should trap error
	}

	log.Printf("getAllUsers(): returning %v -> %v", httpStatus, result)
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(result)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	log.Println("deleteUser(): invoked")
	var result UserOperationResult
	var httpStatus int

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Invalid data - expected Username, Email, and Password for new user")
	}

	// get user data from json. Could use error handling.
	var userNameOp UserNameOperation
	json.Unmarshal(reqBody, &userNameOp)
	log.Printf("deleteUser(): request data: %v", userNameOp)

	var retCode ModelStatusCode
	result.User, retCode, result.Reason = modelDeleteUser(userNameOp.UserName)
	result.Status = ModelStatusText(retCode)

	// handle response.'
	switch retCode {
	case ModelSuccess:
		httpStatus = http.StatusOK
	case ModelDBCreateFailure:
		log.Println("deleteUser(): server error 1")
		httpStatus = http.StatusInternalServerError
	default:
		log.Println("deleteUser(): server error 2")
		httpStatus = http.StatusInternalServerError // should trap error
	}

	log.Printf("deleteUser(): returning %v -> %v", httpStatus, result)
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(result)
}

func deleteAllUsers(w http.ResponseWriter, r *http.Request) {
	log.Println("deleteAllUsers(): invoked")
	var result SimpleOperationResult
	var httpStatus int

	// now create in our db.
	var retCode ModelStatusCode
	retCode, result.Reason = modelDeleteAllUsers()
	result.Status = ModelStatusText(retCode)

	// handle response.'
	switch retCode {
	case ModelSuccess:
		httpStatus = http.StatusOK
	case ModelDBCreateFailure:
		httpStatus = http.StatusInternalServerError
		log.Println("deleteAllUsers(): server error 1")
	default:
		httpStatus = http.StatusInternalServerError
		log.Println("deleteAllUsers(): server error 2") // should trap error
	}

	log.Printf("deleteUser(): returning %v -> %v", httpStatus, result)
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(result)
}
