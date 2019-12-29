# Endpoint
Simple endpoint web service written in Go, backed by a mySQL DB. It includes unit tests.
I implemented an in-memory db and a mySQL db. The mySQL db is the one implemented in user_mode.go. To use the in memory one, just replace user_model.go with user_model_in_memorydb.go_opt and change the extension to .go.

If you wish to run this, you'll need to install Go of course, and then pull down a couple of packages that comprise my framework:
  go get -u github.com/gorilla/mux
  go get -u github.com/go-sql-driver/mysql
    
To run the server and tests open two explorers instances, both in <home>\go\src\endpoint. In one, type
  go build && endpoint
This executes the web server. In the other, type
  go test
