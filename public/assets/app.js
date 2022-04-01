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
      serverUrl: "ws://localhost:8080/ws",
      messages: [],
      newMessage: "",
    };
  },
  methods: {
    connectToWebsocket() {
      this.ws = new WebSocket(this.serverUrl);
      this.ws.addEventListener("open", (event) => {
        this.onWebsocketOpen(event);
      });
      this.ws.addEventListener("message", (event) => {
        this.handleNewMessage(event);
      });
    },
    handleNewMessage(event) {
      let data = event.data;
      // matching Windows and Unix newlines
      data = data.split(/\r?\n/);
      for (let i = 0; i < data.length; i++) {
        let msg = JSON.parse(data[i]);
        this.messages.push(msg);
      }
    },
    sendMessage() {
      if (this.newMessage !== "") {
        this.ws.send(JSON.stringify({ message: this.newMessage }));
        this.newMessage = "";
      }
    },
    onWebsocketOpen() {
      console.log("connected to WS!");
    },
  },
  mounted() {
    this.connectToWebsocket();
  },
};
Vue.createApp(webSocket).mount("#app");
