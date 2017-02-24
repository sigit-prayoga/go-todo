package main

import (
	"encoding/json"
	"log"
	"net/http"

	pg "gopkg.in/pg.v5"
)

var db *pg.DB

func main() {
	// Init DB connection
	connectDb()

	// Register all end-points
	http.HandleFunc("/", helloServer)
	http.HandleFunc("/todos", getTodos)
	http.HandleFunc("/todos/create", addTodo)

	// Start to listen
	http.ListenAndServe(":8383", nil)
}

func helloServer(w http.ResponseWriter, req *http.Request) {
	// Ping
	w.Write([]byte("Hello, client!"))
}

func getTodos(w http.ResponseWriter, req *http.Request) {
	// Init empty Todo
	todos := []ToDo{}

	// Do 'select' into DB
	err := db.Model(&todos).Select()
	if err != nil {
		log.Fatal("Error in Select", err)
	}
	// Send response back
	data, err := json.Marshal(todos)
	if err != nil {
		log.Fatal("Not making any sense", err)
	}
	w.Write(data)
}

func addTodo(w http.ResponseWriter, req *http.Request) {
	// Init new decoder from the Body
	decoder := json.NewDecoder(req.Body)
	var newTodo ToDo
	decodeErr := decoder.Decode(&newTodo)
	if decodeErr != nil {
		log.Fatal("SHOOT! ", decodeErr)
	}
	// Close eventually
	defer req.Body.Close()
	log.Println(newTodo.Label, newTodo.Description)

	// Insert into DB
	insertErr := db.Insert(&newTodo)
	if insertErr != nil {
		log.Fatal("WTF ?!", insertErr)
	}

	// Send response back
	data, err := json.Marshal(newTodo)
	if err != nil {
		log.Fatal("Not making any sense", err)
	}
	w.Write(data)
}

func connectDb() {
	// same as 'psql -U postgres -d postgres' in terminal, connect to localhost:5432
	db = pg.Connect(&pg.Options{
		User:     "postgres",
		Database: "postgres",
	})
}

// ToDo Model = 'todos' table
type ToDo struct {
	tableName   struct{} `sql:"todos"`
	Label       string
	Description string
}
