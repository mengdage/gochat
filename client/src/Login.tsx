import React from "react";
import { Redirect } from "react-router-dom";
import { initApi, getApi } from "./api";
import { useAuthData, useAuthDispatchData } from "./auth";

import './Login.scss'

const Login = () => {
  const [username, setUsername] = React.useState("");
  const [serverAddr, setServerAddr] = React.useState("localhost:8080");
  const authData = useAuthData();
  const authDispatch = useAuthDispatchData();

  const changeUsername = (e: React.ChangeEvent<HTMLInputElement>) => {
    setUsername(e.target.value);
  };
  const changeServerAddr = (e: React.ChangeEvent<HTMLInputElement>) => {
    setServerAddr(e.target.value);
  };

  const handleLogin = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    initApi({ serverAddr });

    const api = getApi();

    try {
      const resp = await api.post("/login", {
        name: username,
      });
      initApi({ serverAddr, authorization: resp.data.user.name });

      authDispatch({
        type: "login",
        body: {
          userId: resp.data.user.id,
          userName: resp.data.user.name,
          serverAddr,
        },
      });
    } catch (e) {
      console.log(e);
    }
  };

  if (authData.isLogin) {
    return <Redirect to="/chat" />;
  }

  return (
    <div className="login-container">
      <h1 className="login-title">Login</h1>
      <form onSubmit={handleLogin}>
        <div className="login-row">
          <label>
            User name
            <input className="login-input" value={username} onChange={changeUsername} />
          </label>
        </div>

        <div className="login-row">
          <label>
            Server Addr
            <input className="login-input" value={serverAddr} onChange={changeServerAddr} />
          </label>
        </div>

        <div className="login-row">
          <button className="login-submit" type="submit">Login</button>
        </div>
      </form>
    </div>
  );
};

export default Login;
