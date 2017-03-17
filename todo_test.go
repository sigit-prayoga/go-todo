package main

import "testing"
import "net/http/httptest"
import "bytes"

const (
	goodMock = `{
        "todo": "Get drunk",
        "done": false,
        "id": "682398ad-0486-4d8a-a33a-e7bd2950c49f"
    }`
	badMock = `{
        "todo": "Get drunk",
        "done": false
        "id": "682398ad-0486-4d8a-a33a-e7bd2950c49f"
    }`
)

// TestParseReq tests `getTodoFromRequest` function
func TestParseReq(t *testing.T) {
	b := []byte(goodMock)
	req := httptest.NewRequest("POST", "/todos/add", bytes.NewBuffer(b))
	todo, err := getTodoFromRequest(req)

	// Success test
	if err != nil {
		t.Error("Should not return error")
	}

	if todo == nil {
		t.Error("Should return todo, instead return nil")
	}

	if todo.Todo == nil {
		t.Error("Should have todo title, instead nil")
	}

	if *todo.Todo != "Get drunk" {
		t.Errorf(`Should return "Get drunk", instead return %s`, *todo.Todo)
	}

	if todo.Done != false {
		t.Errorf("Should return false, instead %b", todo.Done)
	}

	if todo.ID != "682398ad-0486-4d8a-a33a-e7bd2950c49f" {
		t.Errorf(`Should return this id "682398ad-0486-4d8a-a33a-e7bd2950c49f", instead return %s`, todo.ID)
	}

	// Failed test
	b = []byte(badMock)
	req = httptest.NewRequest("POST", "/todos/add", bytes.NewBuffer(b))

	todo, err = getTodoFromRequest(req)

	if err == nil {
		t.Error("Should return error, instead of nil")
	}
}
