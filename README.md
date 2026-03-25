# GlobalXtreme Go Backend Service

[![Release](https://img.shields.io/github/release/gin-gonic/gin.svg?style=flat-square)](https://github.com/globalxtreme/go-backend-service/releases)

## Overview

This is a boilerplate for any Go Projects in GlobalXtreme.
The purpose is to help speed up development and standarized the codes and folder structures.

## Features

The boilerplate already includes some of the commonly used features, like:
- Migration
- Seeder
- [gRPC](https://grpc.io) Integration
- Queue
- [RabbitMQ](https://www.rabbitmq.com) Integration
- Schedule

## ✨Getting Started

### Prerequisites
- **Go Version**: [Go](https://go.dev/) version [1.25](https://go.dev/doc/devel/release#go1.25.0) or above

### Installation

1. Download ["go-create-project.sh"](https://storage.globalxtreme-gateway.net/link/installations/go-create-project.sh) and place it in the project directory.
2. Navigate & allow script to be executed
```shell
cd /path/your-go-project-directory

chmod +x go-create-project.sh
```
3. Create new project.
```shell
./go-create-project.sh <project-name>
```


### Installation with Executable
If you want to install go-create-project globally and call it from any directory, use the following steps:

1. Install shc
```shell
brew install shc
```

2. Convert the shell script "go-create-project.sh" to a binary file
```shell
shc -f go-create-project.sh
```

3. Rename and move binary file into system's PATH
```shell
mv go-create-project.sh.x /usr/local/bin/go-create-project
```

4. Create new project anywhere using global command
```shell
go-create-project <project-and-path-name>
```

### Set Up
Set the environment variables before building project.
Update these variables in the `.env` file
```javascript
VERSION=<version>
SERVICE=<service-name>

DB_HOST=<host>
DB_PORT=<port>
DB_USERNAME=<uname>
DB_PASSWORD=<pass>
DB_DATABASE=<dbname>
```
Set the environment variables for RabbitMQ Connection
```javascript
DB_RABBITMQ_HOST=<host>
DB_RABBITMQ_PORT=<port>
DB_RABBITMQ_USERNAME=<uname>
DB_RABBITMQ_PASSWORD=<pass>
DB_RABBITMQ_DATABASE=<dbname>

RABBITMQ_GLOBAL_HOST=<rabbitmq-host>
RABBITMQ_GLOBAL_PORT=<rabbitmq-port>
RABBITMQ_GLOBAL_USER=<rabbitmq-uname>
RABBITMQ_GLOBAL_PASSWORD=<rabbitmq-pass>
```

### Running the Application
Run the application
```shell
go run cmd/main.go --dev

# or

go build -o application cmd/main.go

./application --dev
```
Runner commands, to run API, gRPC Server, and others.
```shell
# Migration
./application xtreme:migration

# Seeder
./application xtreme:seeder

# gRPC
./application xtreme:grpc

# Queue
./application xtreme:queue

# RabbitMQ
./application xtreme:rabbitmq

# Schedule
./application xtreme:schedule

# Custom Commands (Example)
./application dev-test
```
Generator commands to generate migration file, model, and others.
```shell
# Migration file
./application gen:migration <Name>

# Handler File
./application gen:handler <Name> --type=<web/mobile> --resource

# Model File
./application gen:model <Name> --migration

# Parser File
./application gen:parser <Name> --model
```
⚠️ Add `--dev` flag for development mode.

## �� Documentation
For documentation and Production deployment guide, you can read it on [Go-Lang Backend Service](https://www.notion.so/globalxtreme/Go-Lang-Backend-Service-527f335297b8465f838fc2598538dae7?pvs=4), which we will create later!
