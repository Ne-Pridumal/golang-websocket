import type { TMessage } from "@/models/message";
import type { TUser } from "@/models/user";

class SocketConnection {
  socket: WebSocket

  constructor() {
    this.socket = new WebSocket('ws://localhost:3000/1')
  }

  connect = (cb: (arg: MessageEvent<any>) => void) => {
    console.log('connecting', this.socket.url);

    this.socket.onopen = () => {
      console.log('Successfully Connected!');
    };

    this.socket.onmessage = (msg) => {
      console.log('new message')
      cb(msg);
    };

    this.socket.onclose = event => {
      console.log('Socket Closed Connection: ', event);
    };

    this.socket.onerror = error => {
      console.log('Socket Error: ', error);
    };
  };

  sendMsg = (msg: { type: string }) => {
    // send object as string
    console.log(msg);
    this.socket.send(JSON.stringify(msg));
  };

  connected = (user: TUser) => {
    this.socket.onopen = () => {
      console.log('Successfully Connected', user);
      // initiate mapping
      this.mapConnection(user);
    };
  };

  mapConnection = (user: TUser) => {
    console.log('mapping', user);
    this.socket.send(JSON.stringify({ type: 'bootup', user: user }));
  };
}


export default SocketConnection
