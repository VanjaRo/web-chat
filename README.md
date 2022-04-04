# Web chat

## About

This is a webchat application. During the development, I focused on implementing the backend architecture and WebSocket communication.

## Architecture and design

- As the head of the system, I have a **WebSocket Server**. It controls Users' actions and Rooms' communication.
- **Rooms** are the places where the communication happens. There are two types of them: Public and Private.
  - **Public** rooms can be accessed by their names.
  - **Private** rooms are 1 on 1 rooms. Users can access them only by Partner's invitation.
- Whole communication happens through the **publish-subscribe pattern** by sending **Messages**. Joining a Room, User subscribes to the updates in this room. I use **Redis pub-sub mechanism** to implement that pattern.
- Messages are processed through **Go channels** by concurrently running **goroutines**.
- Users and Rooms are stored in a DB.
  - I use _repository pattern_ for the **data access** mechanism to add an abstraction layer between objects and DB commands.
  - You can run multiple instances of the application using the same DB and Redis.
- I use JWT to authenticate the user. I generate it on the login stage and then pass it back and forth, from frontend to backend. Authentication middleware controls that process.
- User password is secured with the argon2 library.

## Run

To run the app, make sure you have installed:

- Golang (this app is compatible with the 1.16 version) to run the server
- Redis (or Docker + docker-compose) to run Pub/Sub mechanism

### Redis

Start the Redis by yourself on the default port(6379) or use docker-compose (make sure your Docker **daemon is running**):

```sh
docker-compose up
```

### App

To run the application, you need to download the repo and run the command:

```go
go run ./ --addr=your_address
```

The command will start the backend server and make it listen to "your_address".
To access the frontend –– go to "your_address".

### Secrets

There is a JWT secret stored in the ".env" file. Include one in the format I used in the "example.env" file.

## ToDo

- Friends support
  - Adding/Removing from friends
  - Starting private chat only with friends
- Storing messages of an active WS connection
