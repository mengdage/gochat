import React from "react";
import { useAuthData } from "./auth";
import { getApi } from "./api";

const UserList = (props) => {
  const { selectUser } = props;
  const auth = useAuthData();
  const [userList, setUserList] = React.useState([]);

  React.useEffect(() => {
    if (!auth.userName) {
      return;
    }

    const api = getApi();
    let toCancel;
    async function getUserList() {
      try {
        const resp = await api.get("/user/list", {
          headers: {
            Authorization: auth.userName,
          },
        });

        // console.log(resp.data)
        setUserList(resp.data);
        toCancel = setTimeout(getUserList, 60 * 1000);
      } catch (e) {
        console.log(e);
      }
    }

    getUserList();

    return () => {
      clearTimeout(toCancel);
    };
  }, [auth.userName]);

  if (!auth.isLogin) {
    return null;
  }
  return (
    <ul>
      {userList.map((user) => (
        <li
          key={user.id}
          onClick={() => {
            selectUser(user);
          }}
        >{`${user.name}`}</li>
      ))}
    </ul>
  );
};

export default UserList;
