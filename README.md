# Book Store Service

This service is built with Go.

## Overview

- Authentication Login and Register
- Create, Read, Update, Delete operations for Books and Authors.
- List all Books for a specific Author.
- List all Authors for a specific Book.

## Requirements

- Golang version 1.20.2+

## Tools

- Postgresql
- GORM
- Gin Gonic (framework)
- JWT

## Run Code in lokal
```
local database setup based on environment variable
```

```
go mod tidy
```

```
go run main.go
```

## Building Binary

```
go build -o book-store main.go
```

## Run the Binary

```
./book-store
```