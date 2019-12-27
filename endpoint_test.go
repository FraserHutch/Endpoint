package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
)

// Unit test client for endpoint.go

var baseURL = "http://localhost:8080/user/"

// Users - temporary (in memory) database for users
type testUsers []User

var myUsers = testUsers{
	{
		UserName: "Alfie",
		Email:    "alfie@some_office.org",
		Password: "passwrd1",
	},
	{
		UserName: "Joan",
		Email:    "joan@some_other_org.com",
		Password: "passwrd2",
	},
	{
		UserName: "Tony",
		Email:    "tones@somewhere_completly_different.net",
		Password: "passwrd3",
	},
}
var badUser = User{
	UserName: "",
	Email:    "luser@moron.com",
	Password: "passwrd1",
}

func deleteAll() bool {
	url := baseURL + "deleteAll"
	buf := new(bytes.Buffer)
	req, err := http.NewRequest("DELETE", url, buf)

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return true
}

// create a single user. On success, return true and the response user, otherwise return
// false and an error and an error string.
func testCreate(user User) (bool, string, UserOperationResult) {
	var createResp UserOperationResult // where we will  write the resp[onse object's user record
	log.Printf("    creating user %v", user.UserName)
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(user)
	req, err := http.NewRequest("POST", baseURL+"register", buf)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Sprintf("create request failed for user %+v: %v", user, err), createResp
	}
	defer resp.Body.Close()
	if resp.StatusCode != 201 {
		return false, fmt.Sprintf("create request failed for user %+v, expected request status code of 200, got  %+v", user, resp.StatusCode), createResp
	}
	log.Println("response Status:", resp.Status)
	log.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("response Body:", string(body))
	json.Unmarshal(body, &createResp)
	return true, "", createResp
}

// create a single user. On success, return true and the response user, otherwise return
// false and an error and an error string.
func testGet(user User) (bool, string, UserOperationResult) {
	var op = UserNameOperation{UserName: user.UserName}
	var getResp UserOperationResult // where we will  write the resp[onse object's user record
	log.Printf("    retrieving user %v", user.UserName)
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(op)
	req, err := http.NewRequest("GET", baseURL+"get", buf)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Sprintf("get request failed for user %+v: %v", user.UserName, err), getResp
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return false, fmt.Sprintf("get request failed for user %+v, expected request status code of 200, got  %+v", user, resp.StatusCode), getResp
	}
	log.Println("response Status:", resp.Status)
	log.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("response Body:", string(body))
	json.Unmarshal(body, &getResp)
	return true, "", getResp
}

// create a single user. On success, return true and the response user, otherwise return
// false and an error and an error string.
func testUpdate(user User) (bool, string, UserOperationResult) {
	var getResp UserOperationResult // where we will  write the resp[onse object's user record
	log.Printf("    updating user %v: %v", user.UserName, user)
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(user)
	req, err := http.NewRequest("PUT", baseURL+"update", buf)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Sprintf("update request failed for user %+v: %v", user.UserName, err), getResp
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return false, fmt.Sprintf("update request failed for user %+v, expected request status code of 200, got  %+v", user, resp.StatusCode), getResp
	}
	log.Println("response Status:", resp.Status)
	log.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("response Body:", string(body))
	json.Unmarshal(body, &getResp)
	return true, "", getResp
}

// create a single user. On success, return true and the response user, otherwise return
// false and an error and an error string.
func testDelete(userName string) (bool, string, UserOperationResult) {
	var op = UserNameOperation{UserName: userName}
	var getResp UserOperationResult // where we will  write the resp[onse object's user record
	log.Printf("    deleting user %v", userName)
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(op)
	req, err := http.NewRequest("DELETE", baseURL+"delete", buf)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Sprintf("get request failed for user %+v: %v", userName, err), getResp
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return false, fmt.Sprintf("get request failed for user %+v, expected request status code of 200, got  %+v", userName, resp.StatusCode), getResp
	}
	log.Println("response Status:", resp.Status)
	log.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("response Body:", string(body))
	json.Unmarshal(body, &getResp)
	return true, "", getResp
}

// create a single user. On success, return true and the response user, otherwise return
// false and an error and an error string.
func testGetAll() (bool, string, UserGetAllOperationResult) {
	var getResp UserGetAllOperationResult // where we will  write the resp[onse object's user record
	log.Println("    retrieving all users")
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode("")
	req, err := http.NewRequest("GET", baseURL+"getAll", buf)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Sprintf("get all request failed: %v", err), getResp
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return false, fmt.Sprintf("get all request failed, expected request status code of 200, got  %+v", resp.StatusCode), getResp
	}
	log.Println("response Status:", resp.Status)
	log.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("response Body:", string(body))
	json.Unmarshal(body, &getResp)
	return true, "", getResp
}

func findUserInArray(userName string, users []User) bool {
	// create all the users defined above.
	for _, user := range users {
		if user.UserName == userName {
			return true
		}
	}
	return false
}
func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	// TODO m- move deleteAll calls here
	os.Exit(m.Run())
}

// Test creating a single user record (register) and reading it back (get)
func TestCreate(t *testing.T) {

	log.Print("Starting unit test create1")
	if ret := deleteAll(); ret == false {
		t.Error("delete all request failed")
	}

	user := myUsers[0]
	if success, msg, userResp := testCreate(user); success == true {
		if userResp.User.UserName != user.UserName {
			t.Errorf("bad resp record: %v", userResp)
		} else {
			log.Printf("    created user %v", user.UserName)
		}
	} else {
		t.Error(msg)
	}

	// now read it back
	if success, msg, userResp := testGet(user); success == true {
		if userResp.User.UserName != user.UserName {
			t.Errorf("    bad resp record: %v", userResp)
		} else {
			log.Printf("    retrieved user %v", user.UserName)
		}
	} else {
		t.Error(msg)
	}
}

// Test creating a single user record (register) and reading it back (get)
func TestCreateAll(t *testing.T) {
	log.Print("Starting unit test createAll")
	if ret := deleteAll(); ret == false {
		t.Error("delete all request failed")
	}

	// create all the users defined above.
	for ndx, user := range myUsers {
		log.Printf("    creating user %v of %v: %v", ndx, len(myUsers), user.UserName)
		if success, msg, userResp := testCreate(user); success == true {
			if userResp.User.UserName != user.UserName {
				t.Errorf("    bad resp record: %v", userResp)
			} else {
				log.Println("    retrieved OK")
			}
		} else {
			t.Error(msg)
		}
	}
	// now read them all back
	for ndx, user := range myUsers {
		log.Printf("    retrieving user %v of %v: %v", ndx, len(myUsers), user.UserName)
		if success, msg, userResp := testGet(user); success == true {
			if userResp.User.UserName != user.UserName {
				t.Errorf("    bad resp record: %v", userResp)
			} else {
				log.Println("    retrieved OK")
			}
		} else {
			t.Error(msg)
		}
	}
	// now test get all to verify we have the expected number of records.
	log.Println("    retrieving all users")
	if success, msg, userResp := testGetAll(); success == true {
		if userResp.Count != len(myUsers) {
			t.Errorf("    expected %v users, got %v", len(myUsers), userResp.Count)
			// simple iteration to checl that we have the proper users.
			for _, user := range userResp.Users {
				if findUserInArray(user.UserName, myUsers) == false {
					t.Errorf("    user %v not found in response list", user.UserName)
				}
			}
		} else {
			log.Println("    retrieved OK")
		}
	} else {
		t.Error(msg)
	}

}

func TestUpdate(t *testing.T) {
	log.Print("Starting unit test Update")
	if ret := deleteAll(); ret == false {
		t.Error("delete all request failed")
	}

	// create all the users defined above.
	for ndx, user := range myUsers {
		log.Printf("    creating user %v of %v: %v", ndx, len(myUsers), user.UserName)
		if success, msg, userResp := testCreate(user); success == true {
			if userResp.User.UserName != user.UserName {
				t.Errorf("    bad resp record: %v", userResp)
			} else {
				log.Println("    retrieved OK")
			}
		} else {
			t.Error(msg)
		}
	}

	// now test get all to verify we have the expected number of records.
	log.Println("    retrieving all users")
	if success, msg, userResp := testGetAll(); success == true {
		if userResp.Count != len(myUsers) {
			t.Errorf("    expected %v users, got %v", len(myUsers), userResp.Count)
			// simple iteration to checl that we have the proper users.
			for _, user := range userResp.Users {
				if findUserInArray(user.UserName, myUsers) == false {
					t.Errorf("    user %v not found in response list", user.UserName)
				}
			}
		} else {
			log.Println("    retrieved OK")
		}
	} else {
		t.Error(msg)
	}

	// now, let's change a passwrd to something silly.
	var user = myUsers[1]
	var newPassword = "newPassword"
	user.Password = newPassword
	if success, msg, userResp := testUpdate(user); success == true {
		if userResp.User.UserName != user.UserName {
			t.Errorf("    update failed, expected username %v, got %v", user.UserName, userResp.User.UserName)
		}
		if userResp.User.Password != newPassword {
			t.Errorf("    update failed, expected password %v, got %v", newPassword, userResp.User.Password)
		} else {
			log.Println("    updated OK")
		}
	} else {
		t.Error(msg)
	}

}

func TestDelete(t *testing.T) {
	log.Print("TestDelete():Starting unit test Update")
	if ret := deleteAll(); ret == false {
		t.Error("delete all request failed")
	}

	// create all the users defined above.
	for ndx, user := range myUsers {
		log.Printf("TestDelete():creating user %v of %v: %v", ndx, len(myUsers), user.UserName)
		if success, msg, userResp := testCreate(user); success == true {
			if userResp.User.UserName != user.UserName {
				t.Errorf("TestDelete():bad resp record: %v", userResp)
			} else {
				log.Println("TestDelete():created OK")
			}
		} else {
			t.Error(msg)
		}
	}

	// DELETE USER
	var expectedUsers []User
	expectedUsers = append(expectedUsers, myUsers[0])
	var deletedUser = myUsers[1]
	expectedUsers = append(expectedUsers, myUsers[2])

	log.Println("TestDelete():deleting user", deletedUser.UserName)
	if success, msg, userResp := testDelete(deletedUser.UserName); success == true {
		if userResp.User.UserName != deletedUser.UserName {
			t.Errorf("TestDelete():bad resp record: %v", userResp)
		} else {
			log.Println("TestDelete():deleted OK")
		}
	} else {
		t.Error(msg)
	}

	// now we expect two and only two records, and they should match expectedUsers above.
	log.Println("TestDelete():testing remaining users - retrieving all")
	if success, msg, userResp := testGetAll(); success == true {
		if userResp.Count != len(expectedUsers) {
			t.Errorf("TestDelete():expected %v users, got %v", len(expectedUsers), userResp.Count)
			// simple iteration to checl that we have the proper users.
			for _, user := range userResp.Users {
				if findUserInArray(user.UserName, expectedUsers) == false {
					t.Errorf("TestDelete():user %v not found in response list", user.UserName)
				}
			}
		} else {
			log.Println("    retrieved OK")
		}
	} else {
		t.Error(msg)
	}
}
