# golang-crud-api

# Resources

- "Receiver has a generic name"
    - https://blog.devgenius.io/how-should-you-name-your-receivers-in-go-aec60abd7f67
    - https://github.com/golang/go/wiki/CodeReviewComments#receiver-names
- Linters
    - https://pre-commit.com/
        - brew install pre-commit
        - pre-commit install
- Race condition
    - https://dev.to/antsmartian/creating-a-simple-concurrent-database-in-go-2l7j
    - https://pavledjuric.medium.com/thread-safety-in-golang-47fa856fb8bb

# Run locally

1. Start docker compose

```bash
cd local
docker-compose up
```

2. Set environment variables

```bash
DATABASE_USER=user
DATABASE_PASSWORD=password
DATABASE_NAME=db
DATABASE_HOST=localhost
DATABASE_PORT=5432
CHECK_ISBN_CLIENT_HOST=https://my-json-server.typicode.com/joseboretto/golang-testcontainers-gherkin-setup
```

3. Run the application

```bash
cd cmd/golang-crud-api
go run main.go
```

# API Documentation
1. [Create book](#create-book)
```shell
curl --location --request POST 'http://localhost:8000/api/v1/createBook' \
--header 'Content-Type: application/json' \
--data '{
    "title": "title",
    "total_pages": 10,
    "isbn": "0-061-96436-0"
}'
```
2. [Get books](#get-books)
```shell
curl --location --request GET 'http://localhost:8000/api/v1/getBooks'
```
