
# Simple Blog Service

A simple backend service made in golang, to simulate simple blog website usecases

## Features
* Simple auth (login, register) with JWT integration, with basic auth as auth module
* CRUD blog posts with the generated JWT
* CR blog posts comments with the generated JWT

## Tools Stack & prerequisites
* golang
* mysql (need to install it locally / use docker)
* JWT
* godotenv golang lib
* squirrel query builder golang lib

## DB Schema
WIP

## DB Indexes
* UNIQUE users.email : to speed up auth existing email register query check process
* INDEX users.created_at : to speed up get recently registered users query
* FK posts.author_id <- users.id : to speed up potential join queries
* INDEX posts.created_at : to speed up get recent blog posts timeline
* FK comments.post <- posts.id : to speed up potential join queries
* INDEX comments.created_at : to speed up get comments (mostly sorted by the recent comments)

## Quickstart
* Migrate the DB schema, the script is located at `migrations/migrations.sql` file
* Make local env file
```
 touch .env
```
* Fill the .env file with the desired variables based on your local setting
```
JWT_SECRET_KEY=<JWT_SECRET_KEY>
DB_USER=<DB_USER>
DB_NAME=<DB_NAME>
DB_PORT=3306
DB_PASSWORD=<DB_PASSWORD>
AUTH_V1_BASIC_AUTH_USERNAME=<AUTH_V1_BASIC_AUTH_USERNAME>
AUTH_V1_BASIC_AUTH_PASSWORD=<AUTH_V1_BASIC_AUTH_PASSWORD>
```
* Install the go required packages
```
go mod tidy
```
* Run the server
```
go run main.go
```
* Enjoy !

### Future Updates
* Add content categories & sections feature
* add reply comment (subcomment)
* Add upvote/downvote feature
* Add user contribution & rating point