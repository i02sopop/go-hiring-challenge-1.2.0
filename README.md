Go Hiring Challenge
===================

This repository contains a Go application for managing products and their prices, including functionalities for CRUD operations and seeding the database with initial data.

Project Structure
-----------------

1.	**cmd/**: Contains the main application and seed command entry points.

	-	`server/main.go`: The main application entry point, serves the REST API.
	-	`server/server.go`: Handles the api server to start or stop the server.
	-	`seed/main.go`: Command to seed the database with initial product data.

2.	**internal/api/**: Contains the application API handlers.

3.	**sql/**: Contains a very simple database migration scripts setup.

4.	**internal/model/**: Contains the data models used in the application.

5.	**internal/storage/**: Contains the storage layer of the application

6.	`.env`: Environment variables file for configuration.

Application Setup
-----------------

-	Ensure you have Go installed on your machine.
-	Ensure you have Docker installed on your machine.
-	Important makefile targets:
	-	`make build`: Build the binary.
	-	`make tidy`: will install all dependencies.
	-	`make docker-up`: will start the required infrastructure services via docker containers.
	-	`make seed`: ⚠️ Will destroy and re-create the database tables.
	-	`make docker-down`: Will stop the docker containers.
	-	`make test`: Will run the tests.
	-	`make test-update`: Will run the tests and update the snapshots.
	-	`make run`: Will start the application.
	-	`make lint`: Run the linter on the code.
	-	`make lint-fix`: Run the linter on the code and fix some of the easy linter errors.
	-	`make clean`: Clean some of the generated files.
	-	`make dependencies`: Install the go dependencies.

Follow up for the assignemnt here: [ASSIGNMENT.md](ASSIGNMENT.md)
