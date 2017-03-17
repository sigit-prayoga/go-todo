package main

import (
	"encoding/json"
	"log"
	"net/http"

	"fmt"

	"github.com/sigit-prayoga/printj"
	pg "gopkg.in/pg.v5"
)

var db *pg.DB

func main() {

	// Init DB connection
	db = initDB()

	h := http.NewServeMux()

	// Register all end-points
	h.HandleFunc("/", errorHandler(helloServer))
	h.HandleFunc("/todos", errorHandler(requestTodo))

	log.Println("Listening...")
	log.Fatal(http.ListenAndServe(":8383", h))
}

func errorHandler(f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("handling %q: %v", r.RequestURI, err)
		}
	}
}

func helloServer(res http.ResponseWriter, req *http.Request) error {
	// Ping
	res.Write([]byte("Hello, client!"))
	return nil
}

func requestTodo(res http.ResponseWriter, req *http.Request) error {
	var err error
	switch req.Method {
	case "POST":
		err = addTodo(res, req)
	case "PUT":
		err = updateTodo(res, req)
	case "DELETE":
		err = deleteTodo(res, req)
	default:
		err = getTodos(res, req)
	}

	return err
}

func getTodos(w http.ResponseWriter, req *http.Request) error {
	// Init empty Todo
	todos := []*Todo{}

	// Do 'select' into DB
	err := db.Model(&todos).Select()
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	writeResponse(w, &todos, http.StatusOK)

	return nil
}

func addTodo(w http.ResponseWriter, req *http.Request) error {
	newTodo, err := getTodoFromRequest(req)

	// Insert into DB
	err = db.Insert(&newTodo)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	writeResponse(w, &newTodo, http.StatusOK)

	return nil
}

func updateTodo(w http.ResponseWriter, req *http.Request) error {
	todo, err := getTodoFromRequest(req)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	// only 'done' column will get updated based on 'id'
	_, err = db.Model(todo).Column("done").Update()
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	writeResponse(w, &todo, http.StatusOK)

	return nil
}

func deleteTodo(w http.ResponseWriter, req *http.Request) error {
	todo, err := getTodoFromRequest(req)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	err = db.Delete(todo)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	// just give it back to them!
	writeResponse(w, &todo, http.StatusOK)

	return nil
}

func initDB() *pg.DB {
	// same as 'psql -U postgres -d postgres' in terminal, connect to localhost:5432
	db := pg.Connect(&pg.Options{
		User:     "postgres",
		Database: "postgres",
	})
	return db
}

func getTodoFromRequest(r *http.Request) (*Todo, error) {
	// Init new decoder from the Body
	decoder := json.NewDecoder(r.Body)
	var newTodo Todo
	err := decoder.Decode(&newTodo)
	if err != nil {
		return nil, err
	}
	// Close eventually
	defer r.Body.Close()

	return &newTodo, nil
}

func writeResponse(w http.ResponseWriter, v interface{}, code int) {
	printj.Print(v, true, "Response")
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
