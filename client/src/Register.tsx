import React from "react";
import { Redirect } from "react-router-dom";
import { initApi, getApi } from "./api";

import "./Register.scss";

const Register = () => {
  const [username, setUsername] = React.useState("");
  const [isResigered, setIsRegister] = React.useState(false);
  const [serverAddr, setServerAddr] = React.useState("localhost:8080");

  const handleRegister = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    
    initApi({ serverAddr });
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
    <div className="register-container">
      <h1 className="register-title">Register</h1>
      <form onSubmit={handleRegister}>
        <div className="register-row">
          <label>
            User name
            <input
              className="register-input"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
            />
          </label>
        </div>

        <div className="register-row">
          <label>
            Server Addr
            <input
              className="register-input"
              value={serverAddr}
              onChange={(e) => setServerAddr(e.target.value)}
            />
          </label>
        </div>

        <div className="register-row">
          <button className="register-submit" type="submit">
            Login
          </button>
        </div>
      </form>
    </div>
  );
};

export default Register;
