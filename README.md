# Web chat

## About

This is a web chat application. During it's development I was mostly focused on implementing the backend part.

## Architecture and design

- As a head of the system I have a **WebSocket Server**. It controlls Users' actions and Rooms communication.
- **Rooms** are the places, where the communication happens. There are two types of them: Public and Private.
  - **Public** rooms can be accessed by their names.
  - **Private** roooms are 1 on 1 rooms. They can only be accessed by an invitation.
- Whole communication happens through **publish-subscribe pattern** by sending **Messages**. Joining a Room, User subscribes to the updates in this room. It is implemented by **Redis pub-sub mechanism**.
- Messages are processed through **Go channels** by a cuncurrenly running **goroutines**.
- Users and Rooms are stored in a DB.
  - The **data access** mecanism is done through _interfaces_ and _repository pattern_. It's done to minimise the amount of code you need to rewrite to swap the DB.
  - You can run multiple instances of the application using same DB and Redis.
- The authentication is done by JWT. It is passed to frontend and verified by the Authentication middleware.
- User password is secured with argon2 library.

## Run

To run the app make sure you have installed:

- Golang (this app was written with 1.16 version) to run server
- Redis (or Docker + docker-compose) to run Pub/Sub mechanism

### Redis

Start the Redis by yourself on the default port(6379) or use docker-compose (make shure you Docker **daemon is running**):

```sh
docker-compose up
```

### App

To run the application you need to download the repo and run the command:

```go
go run ./ --addr=your_address
```

This will start the backend server and make it listen to "your_address".  
To access the frontend –– go to "your_address".

### Secrets

There is a JWT secret stored in ".env" file. Make shure to include one in format I used in "example.env" file.

## ToDo

- Friends support
  - Adding/Removing from friends
  - Starting private chat only with friends
- Storing messages of an active WS connection
