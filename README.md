# go-api

First adventure using the GoLang programming language

### Pre-requisites
Create a file named .env in the project's root directory with the following content below:
```
DB_DRIVER=postgres
DB_HOST=<hostname>
DB_PORT=5432
DB_USER=user
DB_PASS=password
DB_NAME=database
WEB_SERVER_PORT=8000
JWT_SECRET=secret
JWT_EXPIRES_IN=1d
```

### GO RUN
```bash
❯ go run cmd/server/main.go
```