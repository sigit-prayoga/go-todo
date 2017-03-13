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

## Test your application
```sh
# Get all todo
$ curl http://localhost:8383/todos

# Add new todo
$ curl -X POST -d '{"todo":"Learn more about Go"}' http://localhost:8383/todos/add

# Update todo
$ curl -X POST -d '{"id":"<some_id>"}' http://localhost:8383/todos/update
```

**Note** Make sure you have posgres installed in your computer, with username `postgres` and db name `postgres`
```sh
Macbook:~ sigitprayoga$ psql -U postgres
psql (9.5.3)
Type "help" for help.

postgres=# \d+
                    List of relations
 Schema | Name  | Type  |  Owner   | Size  | Description
--------+-------+-------+----------+-------+-------------
 public | todos | table | postgres | 16 kB |
(1 row)

```