import React from "react";
import cx from "classnames";
import { ChatMessage } from "../types";
import "./MessageItem.scss";

interface MessageItemProps {
  currentUserName: string;
  message: ChatMessage;
}

const MessageItem: React.FC<MessageItemProps> = ({
  currentUserName,
  message,
}) => {
  const isSentByCurrentUser = message.fromUser === currentUserName;

  return (
    <div
      className={cx("messageitem-container", { sentbyme: isSentByCurrentUser })}
    >
      <div
        className={cx("messageitem-msg", {
          sentbyme: isSentByCurrentUser,
        })}
      >
        {message.content}
      </div>
      <div className="messageitem-fromuser">{message.fromUser}</div>
    </div>
  );
};

export default MessageItem;
