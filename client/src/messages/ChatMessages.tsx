import React from "react";
import ScrollToBottom from "react-scroll-to-bottom";

import { ChatMessage } from "../types";
import MessageItem from "./MessageItem";

import "./ChatMessages.scss";

interface ChatMessagesProps {
  currentUserName: string;
  messages: ChatMessage[];

}

const ChatMessages: React.FC<ChatMessagesProps> = ({
  messages,
  currentUserName,
}) => {
  return (
    <ScrollToBottom className="chatmessage-container">
      {messages.map((m) => (
        <MessageItem
          key={m.createdAt}
          message={m}
          currentUserName={currentUserName}
        />
      ))}
    </ScrollToBottom>
  );
};

export default ChatMessages;
