# go-todo
Simple Todo list backend with Go

## Install all dependencies
```sh
$ go get
```

## Compile and run
```sh
$ go run main.go
```

## Install postgres
```sh
$ sudo apt-get install postgres

# Switch to 'postgres'
$ sudo su - postgres

# Start to work with psql
$ psql
```

## Create Todos table
Type this 
`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";` to work with UUID. And then, use 
`create table todos(todo text, done boolean default false, id uuid not null default uuid_generate_v4());` to create a `todos` table. 
To verify table, use `\d+ todos`.

## Test your application
```sh
# Get all todo
$ curl http://localhost:8383/todos

# Add new todo
$ curl -X POST -d '{"todo":"Learn more about Go"}' http://localhost:8383/todos

# Update todo
$ curl -X UPDATE -d '{"todo":"Learn more about Go", "id":"<some_id>", "done":true}' http://localhost:8383/todos

# Delete todo
$ curl -X DELETE -d '{"id":"<some_id>"}' http://localhost:8383/todos
```