import React from "react";
import { Redirect } from "react-router-dom";
import { initApi, getApi } from "./api";

const Register = () => {
  const [username, setUsername] = React.useState("");
  const [isResigered, setIsRegister] = React.useState(false);
  const [serverAddr, setServerAddr] = React.useState("localhost:8080");

  const handleRegister = async () => {
    initApi({serverAddr});
    const api = getApi();
    try {
      await api.post("/register", {
        name: username,
      });

      setIsRegister(true);
    } catch (e) {
      console.log(e);
    }
  };

  if (isResigered) {
    return <Redirect to="/login" />;
  }

  return (
    <div>
      <h1>Register</h1>
      <div>
        <label>
          User name:
          <input
            value={username}
            onChange={(e) => {
              setUsername(e.target.value);
            }}
          />
        </label>
      </div>
      <div>
        <label>
          Server Addr:
          <input
            value={serverAddr}
            onChange={(e) => {
              setServerAddr(e.target.value);
            }}
          />
        </label>
      </div>
      <button onClick={handleRegister}>Register</button>
    </div>
  );
};

export default Register;
