package main

// This is the user model - it roughly corresponds to the model part of MVP
// This implementation is an mySQL db
// go get -u github.com/go-sql-driver/mysql
// Status codes are defined in user_model_status.go

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// MyDB - mySql connection data.
type MyDB struct {
	dbName     string
	tableName  string
	connection *sql.DB
}

func (dbInfo *MyDB) closeDBConnection() {
	if dbInfo.connection != nil {
		dbInfo.connection.Close()
		dbInfo.connection = nil
	}
}

func (dbInfo MyDB) isValidDBConnection() bool {
	return dbInfo.connection != nil
}

func (dbInfo *MyDB) openDBConnection() bool {
	var err error
	dbInfo.connection, err = sql.Open("mysql", "root:@Bubba1111@tcp(127.0.0.1:3306)/"+dbInfo.dbName)
	if err != nil {
		log.Printf("openDBConnection(): ERROR: failed to open db %v: %v", dbInfo.dbName, err)
		return false
	}
	return true
}

var myDB = MyDB{dbName: "entrypoint", tableName: "usersTest"}

// User - basic user definition
//
// The UserName is the key. The id would usually be the primary key and the UserName a secondary key.
type User struct {
	ID       int    `json:"ID"`
	UserName string `json:"UserName"`
	Email    string `json:"Email"`
	Password string `json:"Password"`
}

// Simply determine if the requisite table exists, and if not, create it
func checkAndCreateTable() bool {
	if myDB.isValidDBConnection() == false {
		log.Printf("    no db connection")
		return false
	}
	query := fmt.Sprintf("Show tables like '%v'", myDB.tableName)
	tableResp, err := myDB.connection.Query(query)
	if err != nil {
		log.Printf("    error looking up table information for table %v", myDB.tableName)
		return false
	}
	count := 0
	for tableResp.Next() {
		count++
	}
	if count > 0 {
		return true
	}
	log.Printf("    table '%v' not found, will create", myDB.tableName)
	// I know, I know, I should data drive this from the User struct , ,gain, not that ambitious.
	createTableQuery := "create table " + myDB.tableName + " (ID int NOT NULL AUTO_INCREMENT, UserName varchar(255) NOT NULL UNIQUE, email varchar(255), password varchar(255), PRIMARY KEY (ID));"
	stmt, err := myDB.connection.Prepare(createTableQuery)
	if err != nil {
		log.Printf("    command to prepare statement to create table '%v' failed: %v", myDB.tableName, err)
	}
	if _, err = stmt.Exec(); err != nil {
		log.Printf("    command to create table '%v' failed: %v", myDB.tableName, err)
		return false
	}

	log.Printf("    table '%v' created successfully..", myDB.tableName)
	return true
}

func initDB() bool {
	// open our db
	log.Printf("initDB(): opening db %v", myDB.dbName)
	if myDB.openDBConnection() == false {
		return false
	}
	// check for existence of our table, attempt to create if not found
	if checkAndCreateTable() == false {
		return false
	}

	log.Println("initDB(): OK")
	return true
}

func releaseDB() {
	log.Println("releaseDB()")
	if myDB.connection != nil {
		myDB.connection.Close()
	}
	log.Println("releaseDB(): OK")
}

func isValidUser(user User) (bool, string) {
	var ret bool // false by default
	var reason string

	// I know, I know, it only returns the first failure, but we get the idea.
	if len(user.UserName) < 1 {
		reason = "empty user name"
	} else if len(user.Email) < 1 {
		reason = "invalid email"
	} else if len(user.Password) < 1 {
		reason = "invalid password"
	} else {
		ret = true // ok, now is acceptable
	}
	return ret, reason
}

func modelCreateUser(newUser User) (User, ModelStatusCode, string) {
	if myDB.isValidDBConnection() == false && myDB.openDBConnection() == false {
		log.Printf("isValidDBConnection(): no db connection")
		return newUser, ModelDBCreateFailure, "no db connection"
	}

	var retCode ModelStatusCode
	var reason string
	// test for valid record
	if isValid, errorStr := isValidUser(newUser); isValid == false {
		return newUser, ModelDBCreateFailure, errorStr
	}

	// ID is autoincremented
	query := fmt.Sprintf("INSERT into %v VALUES ( NULL, '%v', '%v', '%v' )", myDB.tableName, newUser.UserName, newUser.Email, newUser.Password)

	log.Println("    modelCreateUser(): using", query)
	insert, err := myDB.connection.Query(query)
	if err != nil {
		retCode = ModelDBCreateFailure
		reason = fmt.Sprintf("failed to insert newUser %v: %v", newUser, err)
		return newUser, retCode, reason
	}
	defer insert.Close()

	// todo - pyll newUser here from insert.
	return newUser, ModelSuccess, ""
}

func modelUpdateUser(user User) (User, ModelStatusCode, string) {
	// test for valid record
	if isValid, errorStr := isValidUser(user); isValid == false {
		return user, ModelDBUpdateFailure, errorStr
	}
	if myDB.isValidDBConnection() == false && myDB.openDBConnection() == false {
		log.Printf("modelUpdateUser(): no db connection")
		return user, ModelDBUpdateFailure, "no db connection"
	}

	query := fmt.Sprintf("UPDATE %v SET Email = '%v', Password = '%v' where UserName = '%v';",
		myDB.tableName, user.Email, user.Password, user.UserName)
	log.Printf("modelUpdateUser(): query: %v", query)
	res, err := myDB.connection.Exec(query)
	if err != nil {
		return user, ModelDBUpdateFailure, fmt.Sprintf("failed to update record for user '%v': %v", user.UserName, err)
	}

	numUpdated, err := res.RowsAffected() // we expect one row affected. If 0 we did not find the user.
	if err != nil {
		panic(err)
	}
	if numUpdated == 0 {
		return user, ModelDBUpdateFailure, fmt.Sprintf("user %v not found to update'", user.UserName)
	}
	if numUpdated != 1 {
		return user, ModelDBUpdateFailure,
			fmt.Sprintf("key error updating user '%v', %v instanced updated", user.UserName, numUpdated)
	}

	return user, ModelSuccess, ""
}

func modelGetUser(userName string) (User, ModelStatusCode, string) {
	var user User

	if len(userName) < 1 {
		return user, ModelDBGetFailure, "User name not supplied"
	}

	if myDB.isValidDBConnection() == false && myDB.openDBConnection() == false {
		log.Printf("modelGetUser(): no db connection")
		return user, ModelDBGetFailure, "no db connection"
	}
	query := fmt.Sprintf("SELECT ID, UserName, Email, Password from %v where UserName = '%v'", myDB.tableName, userName)
	log.Printf("modelGetUser(): query: %v", query)
	results, err := myDB.connection.Query(query)
	if err != nil {
		return user, ModelDBGetFailure, fmt.Sprintf("error retrieving record for user '%v': %v", userName, err)
	}
	defer results.Close()
	var cnt = 0
	for results.Next() {
		err = results.Scan(&user.ID, &user.UserName, &user.Email, &user.Password)
		if err != nil {
			return user, ModelDBGetFailure, fmt.Sprintf("failed to pull values from record for user '%v': %v", userName, err)
		}
		cnt++
	}
	if cnt < 1 {
		// user no ound
		return user, ModelDBUserNotFound, fmt.Sprintf("failed to retrieve record for user '%v': user not found", userName)
	}
	return user, ModelSuccess, ""
}

func modelGetAllUsers() ([]User, ModelStatusCode, string) {
	var users []User

	if myDB.isValidDBConnection() == false && myDB.openDBConnection() == false {
		log.Printf("modelGetAllUsers(): no db connection")
		return users, ModelDBGetFailure, "no db connection"
	}
	query := fmt.Sprintf("SELECT ID, UserName, Email, Password from %v", myDB.tableName)
	log.Printf("modelGetAllUsers(): query: %v", query)
	results, err := myDB.connection.Query(query)
	if err != nil {
		return users, ModelDBGetFailure, fmt.Sprintf("failed to retrieve records: %v", err)
	}
	defer results.Close()
	for results.Next() {
		var user User
		err = results.Scan(&user.ID, &user.UserName, &user.Email, &user.Password)
		if err != nil {
			return users, ModelDBGetFailure, fmt.Sprintf("failed to pull values from record: %v", err)
		}
		users = append(users, user)
	}
	return users, ModelSuccess, ""
}

func modelDeleteUser(userName string) (User, ModelStatusCode, string) {
	var user User

	if len(userName) < 1 {
		return user, ModelDBDeleteFailure, "User name not supplied"
	}

	if myDB.isValidDBConnection() == false && myDB.openDBConnection() == false {
		log.Printf("modelGetAllUsers(): no db connection")
		return user, ModelDBGetFailure, "no db connection"
	}

	// pull out old record. Ignore the errors, we'll try to delete it anyways if not found
	oldUser, ret, reason := modelGetUser(userName)
	if ret == ModelDBUserNotFound {
		return oldUser, ret, reason
	}

	query := fmt.Sprintf("DELETE from %v where UserName = '%v'", myDB.tableName, userName)
	log.Printf("modelDeleteUser(): query: %v", query)
	results, err := myDB.connection.Query(query)
	if err != nil {
		return user, ModelDBDeleteFailure, fmt.Sprintf("failed to delete record for user '%v': %v", userName, err)
	}
	defer results.Close()
	return oldUser, ModelSuccess, ""
}

func modelDeleteAllUsers() (ModelStatusCode, string) {
	var users []User

	if myDB.isValidDBConnection() == false && myDB.openDBConnection() == false {
		log.Printf("modelDeleteAllUsers(): no db connection")
		return ModelDBGetFailure, "no db connection"
	}
	query := fmt.Sprintf("TRUNCATE table %v;", myDB.tableName)
	log.Printf("modelDeleteAllUsers(): query: %v", query)
	results, err := myDB.connection.Query(query)
	if err != nil {
		return ModelDBGetFailure, fmt.Sprintf("failed to delete all records: %v", err)
	}
	defer results.Close()
	for results.Next() {
		var user User
		err = results.Scan(&user.ID, &user.UserName, &user.Email, &user.Password)
		if err != nil {
			return ModelDBGetFailure, fmt.Sprintf("failed to delete all records: %v", err)
		}
		users = append(users, user)
	}
	return ModelSuccess, ""
}
