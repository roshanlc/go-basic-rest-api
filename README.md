# go-rest-api
> A basic go based REST API built for learning purpose.
* No third-party packages/dependencies

## Implemented

* [✔] `GET /book` returns list of books as JSON
* [✔] `GET /book/{id}` returns details of a specific book as JSON
* [✔] `POST /book` accepts a new book to be added
* [✔] `POST /book` returns status 415 if content is not `application/json`

### Data types
A book object should look like this:

```json
{
    "id": 3,
    "title": "The Book of Life",
    "author": "Krishnamurti",
    "language": "English",
    "genres": [
        "non-fiction",
        "philosophy"
    ]
}
```

### Persistence

There is no persistence, a global variable is used with  RWMutex lock.

### <u>Run</u>

Since, no third-patry packages are used; just hit:
```go
go run main.go
```