# golang-testcontainers-gherkin-setup

https://medium.com/@joseboretto/integration-tests-in-go-with-cucumber-testcontainers-and-httpmock-6e771f975de9

https://dev.to/joseboretto/integration-tests-in-go-with-cucumber-testcontainers-and-httpmock-5hb9

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
cd cmd/golang-testcontainers-gherkin-setup
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

# External services
This is a mock server. Check https://my-json-server.typicode.com/ for more information.

1. [Check ISBN](#check-isbn)
```shell
curl --location --request GET 'https://my-json-server.typicode.com/joseboretto/golang-testcontainers-gherkin-setup/isbn/0-061-96436-1'
```

1. [Send email](#send-email)
```shell
curl --location --request GET 'https://my-json-server.typicode.com/joseboretto/golang-testcontainers-gherkin-setup/sendEmail'
```
