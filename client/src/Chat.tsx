import React from "react";
import { Redirect } from "react-router-dom";

import { useAuthData } from "./auth";
import UserList from "./UserList";
import InfoBar from "./InfoBar";
import { ChatMessage, User } from "./types";
import MessageInput from "./MessageInput";

import "./Chat.scss";
import ChatMessages from "./messages/ChatMessages";
import { getApi } from "./api";

function sendJson(socket: WebSocket, obj: any) {
  socket.send(JSON.stringify(obj));
}

interface ChatProps {}

const Chat: React.FC<ChatProps> = () => {
  const auth = useAuthData();
  const [socket, setSocket] = React.useState<WebSocket>();
  const [toUser, setToUser] = React.useState<User>();
  const [messages, setMessages] = React.useState<ChatMessage[]>([]);

  React.useEffect(() => {
    if (!auth.userName) {
      return;
    }

    let ws: WebSocket;

    function start() {
      ws = new WebSocket(`ws://${auth.serverAddr}/ws`);
      ws.onopen = function () {
        console.log("Login ...");

        sendJson(ws, {
          cmd: "login",
          body: { userId: auth.userId, userName: auth.userName },
        });
      };

      ws.onclose = () => {
        console.log("Closing ...");
        setTimeout(start, 5000);
      };

      ws.onmessage = (e) => {
        console.log("Received:", e.data);
        const resp = JSON.parse(e.data);
        if (!resp.cmd) {
          return;
        }
        switch (resp.cmd) {
          case "login": {
            if (resp?.body?.ok) {
              setSocket(ws);
            }
            break;
          }
          case "recv": {
            if (resp?.body) {
              setMessages(ms=> [...ms, resp.body]);
            }
            break;
          }
        }
      };
    }

    start();
  }, []);

  React.useEffect(() => {
    async function fetchHistory() {
      const api = getApi();
      try {
        const resp = await api.get<ChatMessage[]>(
          `/user/conversation_history/${toUser?.name}`
        );
        console.log(resp.data);
        setMessages(resp.data);
      } catch (e) {
        console.error(e);
      }
    }

    if (toUser?.name) {
      fetchHistory()
    }

    return function clearMessages() {
      setMessages([])
    }
  }, [toUser?.name]);

  if (!auth.isLogin) {
    return <Redirect to="/login" />;
  }

  if (!socket) {
    return (
      <div>
        <h1>Hello {auth.userName} !</h1>
        <p>Connecting to the server...</p>
      </div>
    );
  }
  const handleChangeToUser = (user: User) => {
    setToUser(user);
  };

  if (!toUser) {
    return (
      <div className="chat-container">
        <UserList selectUser={handleChangeToUser} />
        <div>Select a user to start a chat.</div>
      </div>
    );
  }

  const handleSend = (message: string) => {
    sendJson(socket, {
      cmd: "send",
      body: {
        userName: toUser.name,
        content: message,
      },
    });
  };

  return (
    <div className="chat-container">
      <UserList
        selectUser={handleChangeToUser}
        selectedUserName={toUser.name}
      />
      <div className="chat-conversation">
        <InfoBar toUserName={toUser.name} serverAddr={auth.serverAddr} />

        <ChatMessages messages={messages} currentUserName={auth.userName} />

        <MessageInput onSubmit={handleSend} />
      </div>
    </div>
  );
};

export default Chat;
