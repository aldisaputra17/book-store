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

## Postman Documentation

[![Run in Postman](https://run.pstmn.io/button.svg)](https://app.getpostman.com/run-collection/16404807-2dd94ce8-d495-441c-98e1-970e6bf51fea?action=collection%2Ffork&collection-url=entityId%3D16404807-2dd94ce8-d495-441c-98e1-970e6bf51fea%26entityType%3Dcollection%26workspaceId%3D9722961b-ee27-4ce6-abf2-564c90301265)

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
