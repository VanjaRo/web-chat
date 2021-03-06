// let socket = new WebSocket("ws://localhost:8080/ws");
// console.log("attempting connection");

// socket.onopen = (event) => {
//   console.log("connection established");
// };

// socket.onclose = (event) => {
//   console.log("connection closed");
// };

// socket.onerror = (event) => {
//   console.log("connection error: ", event);
// };

let webSocket = {
  data() {
    return {
      ws: null,
      serverUrl: "ws://" + location.host + "/ws",
      roomInput: null,
      rooms: [],
      user: {
        username: "",
        password: "",
        token: "",
      },
      userRegister: {
        username: "",
        password: "",
      },
      friends: [],
      initialReconnectDelay: 1000,
      currentReconnectDelay: 0,
      maxReconnectDelay: 15000,
      loginError: "",
    };
  },
  methods: {
    connect() {
      this.connectToWebsocket();
    },

    async login() {
      let response = await fetch("http://" + location.host + "/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(this.user),
      });
      let data = await response.json();
      if (data.error) {
        this.loginError = data.error;
      } else {
        this.user.password = "";
        this.user.token = data.token;
        this.connect();
      }
    },
    async register() {
      let response = await fetch("http://" + location.host + "/register", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(this.userRegister),
      });
      let data = await response.json();
      if (data.error) {
        this.loginError = data.error;
      } else {
        this.user.username = this.userRegister.username;
        this.user.password = this.userRegister.password;
        this.login();
      }
    },
    connectToWebsocket() {
      if (this.user.token != "") {
        this.ws = new WebSocket(this.serverUrl + "?bearer=" + this.user.token);

        this.ws.addEventListener("open", (event) => {
          this.onWebsocketOpen(event);
        });
        this.ws.addEventListener("message", (event) => {
          this.handleNewMessage(event);
        });
        this.ws.addEventListener("close", (event) => {
          this.onWebsocketClose(event);
        });
      } else {
        alert("Please login first");
      }
    },
    onWebsocketClose() {
      this.ws = null;
      // reconnect clients in a random time between 1 and 15 seconds
      setTimeout(() => {
        this.reconnectToWebsocket();
      }, this.currentReconnectDelay + Math.floor(Math.random() * 3000));
    },
    reconnectToWebsocket() {
      if (this.currentReconnectDelay < this.maxReconnectDelay) {
        this.currentReconnectDelay *= 2;
      }
      this.connectToWebsocket();
    },
    handleNewMessage(event) {
      let data = event.data;
      // matching Windows and Unix newlines
      data = data.split(/\r?\n/);
      for (let i = 0; i < data.length; i++) {
        let msg = JSON.parse(data[i]);
        switch (msg.action) {
          case "send-message":
            this.handleChatMessage(msg);
            break;
          case "user-join":
            this.handleUserJoined(msg);
            break;
          case "user-left":
            this.handleUserLeft(msg);
            break;
          case "room-joined":
            this.handleRoomJoined(msg);
            break;
          default:
            break;
        }
        let room = this.findRoom(msg.target);
        if (typeof room !== "undefined") {
          room.messages.push(msg);
        }
      }
    },
    handleRoomJoined(msg) {
      room = msg.target;
      room.name = room.private ? msg.sender.name : room.name;
      room.messages = [];
      this.rooms.push(room);
    },
    handleChatMessage(msg) {
      let room = this.findRoom(msg.target.id);
      if (typeof room !== "undefined") {
        room.messages.push(msg);
      }
    },
    sendMessage(room) {
      if (room.newMessage !== "") {
        this.ws.send(
          JSON.stringify({
            action: "send-message",
            message: room.newMessage,
            target: {
              id: room.id,
              name: room.name,
            },
          })
        );
        room.newMessage = "";
      }
    },
    findRoom(roomId) {
      for (let i = 0; i < this.rooms.length; i++) {
        if (this.rooms[i].id === roomId) {
          return this.rooms[i];
        }
      }
    },
    joinRoom() {
      this.ws.send(
        JSON.stringify({ action: "join-room", message: this.roomInput })
      );
      this.roomInput = "";
    },
    joinRoomPrivate(room) {
      this.ws.send(
        JSON.stringify({ action: "join-room-private", message: room.id })
      );
    },
    leaveRoom(room) {
      this.ws.send(JSON.stringify({ action: "leave-room", room: room.name }));
      this.rooms.splice(this.rooms.indexOf(room), 1);
    },
    handleUserJoined(msg) {
      if (!this.userExists(msg.sender)) {
        this.users.push(msg.sender);
      }
    },
    userExists(user) {
      for (let i = 0; i < this.users.length; i++) {
        if (this.users[i].id === user.id) {
          return true;
        }
      }
      return false;
    },
    handleUserLeft(msg) {
      for (let i = 0; i < this.users.length; i++) {
        if (this.users[i].id == msg.sender.id) {
          this.users.splice(i, 1);
        }
      }
    },
    onWebsocketOpen() {
      this.currentReconnectDelay = this.initialReconnectDelay;
      console.log("connected to WS!");
    },
  },
};
Vue.createApp(webSocket).mount("#app");
