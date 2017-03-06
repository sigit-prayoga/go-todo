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
	db = initDB()

	// Register all end-points
	http.HandleFunc("/", helloServer)
	http.HandleFunc("/todos", getTodos)
	http.HandleFunc("/todos/add", addTodo)
	http.HandleFunc("/todos/update", updateTodo)

	log.Println("Listening...")
	log.Fatal(http.ListenAndServe(":8383", nil))
}

func helloServer(w http.ResponseWriter, req *http.Request) {
	// Ping
	w.Write([]byte("Hello, client!"))
}

func getTodos(w http.ResponseWriter, req *http.Request) {
	// Init empty Todo
	todos := []*Todo{}

	// Do 'select' into DB
	err := db.Model(&todos).Select()
	if err != nil {
		log.Fatal("Error in Select", err)
	}

	writeResponse(w, &todos, http.StatusOK)
}

func addTodo(w http.ResponseWriter, req *http.Request) {
	newTodo := getTodoFromRequest(req)

	// Insert into DB
	insertErr := db.Insert(&newTodo)
	if insertErr != nil {
		log.Fatal("WTF ?!", insertErr)
	}

	writeResponse(w, &newTodo, http.StatusOK)
}

func updateTodo(w http.ResponseWriter, req *http.Request) {
	todo := getTodoFromRequest(req)

	// only 'done' column will get updated based on 'id'
	_, err := db.Model(&todo).Column("done").Update()
	if err != nil {
		log.Fatal("Update FAILED!", err)
	}

	writeResponse(w, &todo, http.StatusOK)
}

func initDB() *pg.DB {
	// same as 'psql -U postgres -d postgres' in terminal, connect to localhost:5432
	db := pg.Connect(&pg.Options{
		User:     "postgres",
		Database: "postgres",
	})
	return db
}

func getTodoFromRequest(r *http.Request) Todo {
	// Init new decoder from the Body
	decoder := json.NewDecoder(r.Body)
	var newTodo Todo
	decodeErr := decoder.Decode(&newTodo)
	if decodeErr != nil {
		log.Fatal("SHOOT! ", decodeErr)
	}
	// Close eventually
	defer r.Body.Close()

	return newTodo
}

func writeResponse(w http.ResponseWriter, v interface{}, code int) {
	data, err := json.Marshal(v)
	if err != nil {
		log.Printf("Marshalling error: %v", err)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(data)
}

// Todo Model = 'todos' table
type Todo struct {
	tableName struct{} `sql:"todos"`
	Todo      *string  `json:"todo"`
	Done      bool     `json:"done"`
	ID        string   `json:"id"`
}
