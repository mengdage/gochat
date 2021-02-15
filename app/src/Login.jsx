import React from "react";
import { Redirect } from "react-router-dom";
import { initApi, getApi } from "./api";
import { useAuthData, useAuthDispatchData } from "./auth";

const Login = () => {
  const [username, setUsername] = React.useState("");
  const [serverAddr, setServerAddr] = React.useState("localhost:8081");
  const authData = useAuthData();
  const authDispatch = useAuthDispatchData();

  const changeUsername = (e) => {
    setUsername(e.target.value);
  };
  const changeServerAddr = (e) => {
    setServerAddr(e.target.value);
  };

  const handleLogin = async () => {
    initApi(serverAddr)

    const api = getApi();

    try {
      const resp = await api.post("/login", {
        name: username,
      });

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
    <div>
      <h1>Login</h1>
      <div>
        <label>
          User name:
          <input value={username} onChange={changeUsername} />
        </label>
      </div>
      <div>
        <label>
          Server Addr:
          <input value={serverAddr} onChange={changeServerAddr} />
        </label>
      </div>
      <button onClick={handleLogin}>Login</button>
    </div>
  );
};

export default Login;
