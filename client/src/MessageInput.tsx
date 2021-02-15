import React from "react";
import cx from "classnames";
import "./MessageInput.scss";

interface MessageInputProps {
  onSubmit: (m: string) => void;
}
const MessageInput: React.FC<MessageInputProps> = ({ onSubmit }) => {
  const [message, setMessage] = React.useState("");
  const canSend = Boolean(message);

  return (
    <form
      className="messageinput-form"
      onSubmit={(e) => {
        e.preventDefault();
        if (!canSend) {
          return;
        }

        onSubmit(message);
        setMessage("");
      }}
    >
      <input
        className="messageinput-input"
        value={message}
        onChange={(e) => {
          setMessage(e.target.value);
        }}
      />
      <button
        className={cx("messageinput-send", { disabled: !canSend })}
        disabled={!canSend}
        type="submit"
      >
        Send
      </button>
    </form>
  );
};

export default MessageInput;
