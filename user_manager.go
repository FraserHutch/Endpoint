package main

// This is the user manager - it roughly corresponds to the view part of MVP
// Abstracts the db implementation.
// Most of the model repsonse handling is a simple code mapping to HTTP status codes.

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

//// HANDLERS - these correspond one to one with the API declared in endpoint.go

// POST -> "/user/register"
func createUser(w http.ResponseWriter, r *http.Request) {
	log.Println("createUser(): invoked")
	var result UserOperationResult
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Invalid data - expected Username, Email, and Password for new user")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(result)
		return
	}

	// get user data from json.
	var newUser User
	json.Unmarshal(reqBody, &newUser)
	log.Printf("createUser(): request data: %v", newUser)

	// now update the db.
	var httpStatus int
	var retCode ModelStatusCode
	result.User, retCode, result.Reason = modelCreateUser(newUser)
	result.Status = ModelStatusText(retCode)

	// handle response.
	switch retCode {
	case ModelSuccess:
		httpStatus = http.StatusCreated
	case ModelDBCreateFailure:
		httpStatus = http.StatusInternalServerError
	default:
		log.Printf("createUser(): model returned unexpected status code %v", retCode)
		httpStatus = http.StatusInternalServerError
	}

	log.Printf("createUser(): returning %v -> %v", httpStatus, result)
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(result)
}

// PUT -> "/user/update"
func updateUser(w http.ResponseWriter, r *http.Request) {
	log.Println("updateUser(): invoked")
	var result UserOperationResult
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Invalid data - expected Username, Email, and Password for new user")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(result)
		return
	}

	var user User
	json.Unmarshal(reqBody, &user)
	log.Printf("updateUser(): request data: %v", user)

	// now update the db.
	var httpStatus int
	var retCode ModelStatusCode
	result.User, retCode, result.Reason = modelUpdateUser(user)
	result.Status = ModelStatusText(retCode)

	// handle response.
	switch retCode {
	case ModelSuccess:
		httpStatus = http.StatusOK
	case ModelDBUserNotFound:
		httpStatus = http.StatusNotFound
	case ModelDBUpdateFailure:
		httpStatus = http.StatusInternalServerError
	default:
		log.Printf("updateUser(): model returned unexpected status code %v", retCode)
		httpStatus = http.StatusInternalServerError
	}

	log.Printf("updateUser(): returning %v -> %v", httpStatus, result)
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(result)
}

// GET -> "/user/get"
func getUser(w http.ResponseWriter, r *http.Request) {
	log.Println("getUser(): invoked")
	var result UserOperationResult
	var httpStatus int

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Invalid data - expected Username")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(result)
		return
	}

	// get user data from json. Could use error handling.
	var userNameOp UserNameOperation
	json.Unmarshal(reqBody, &userNameOp)
	log.Printf("getUser(): request data: %v", userNameOp)

	// now retrieve from our db.
	var retCode ModelStatusCode
	result.User, retCode, result.Reason = modelGetUser(userNameOp.UserName)
	result.Status = ModelStatusText(retCode)

	// handle response.
	switch retCode {
	case ModelSuccess:
		httpStatus = http.StatusOK
	case ModelDBUserNotFound:
		httpStatus = http.StatusNotFound
	case ModelDBGetFailure:
		httpStatus = http.StatusInternalServerError
	default:
		log.Printf("getUser(): model returned unexpected status code %v", retCode)
		httpStatus = http.StatusInternalServerError

	}

	log.Printf("getUser(): returning %v -> %v", httpStatus, result)
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(result)
}

// GET -> "/user/getAll"
func getAllUsers(w http.ResponseWriter, r *http.Request) {
	log.Println("getAllUsers() invoked")
	var result UserGetAllOperationResult
	var httpStatus int

	// access db
	var retCode ModelStatusCode
	result.Users, retCode, result.Reason = modelGetAllUsers()
	result.Status = ModelStatusText(retCode)
	result.Count = len(result.Users)

	// handle response.
	switch retCode {
	case ModelSuccess:
		log.Println("getAllUsers(): found")
		httpStatus = http.StatusOK
	case ModelDBCreateFailure:
		log.Println("getAllUsers(): server error 1")
		httpStatus = http.StatusInternalServerError
	default:
		log.Printf("getAllUsers(): model returned unexpected status code %v", retCode)
		httpStatus = http.StatusInternalServerError
	}

	log.Printf("getAllUsers(): returning %v -> %v", httpStatus, result)
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(result)
}

// DELETE -> "/user/delete"
func deleteUser(w http.ResponseWriter, r *http.Request) {
	log.Println("deleteUser(): invoked")
	var result UserOperationResult
	var httpStatus int

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Invalid data - expected Username, Email, and Password for new user")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(result)
		return
	}

	// get user data from json.
	var userNameOp UserNameOperation
	json.Unmarshal(reqBody, &userNameOp)
	log.Printf("deleteUser(): request data: %v", userNameOp)

	// access db
	var retCode ModelStatusCode
	result.User, retCode, result.Reason = modelDeleteUser(userNameOp.UserName)
	result.Status = ModelStatusText(retCode)

	// handle response.
	switch retCode {
	case ModelSuccess:
		httpStatus = http.StatusOK
	case ModelDBCreateFailure:
		log.Println("deleteUser(): server error")
		httpStatus = http.StatusInternalServerError
	default:
		log.Printf("deleteUser(): model returned unexpected status code %v", retCode)
		httpStatus = http.StatusInternalServerError
	}

	log.Printf("deleteUser(): returning %v -> %v", httpStatus, result)
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(result)
}

// DELETE -> "/user/deleteAll"
func deleteAllUsers(w http.ResponseWriter, r *http.Request) {
	log.Println("deleteAllUsers(): invoked")
	var result SimpleOperationResult
	var httpStatus int

	// access db
	var retCode ModelStatusCode
	retCode, result.Reason = modelDeleteAllUsers()
	result.Status = ModelStatusText(retCode)

	// handle response.
	switch retCode {
	case ModelSuccess:
		httpStatus = http.StatusOK
	case ModelDBCreateFailure:
		httpStatus = http.StatusInternalServerError
		log.Println("deleteAllUsers(): server error 1")
	default:
		log.Printf("deleteAllUsers(): model returned unexpected status code %v", retCode)
		httpStatus = http.StatusInternalServerError
	}

	log.Printf("deleteAllUsers(): returning %v -> %v", httpStatus, result)
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(result)
}
