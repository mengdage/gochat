import React from "react";
import cx from "classnames";

import { useAuthData } from "./auth";
import { getApi } from "./api";
import { User } from "./types";

import "./UserList.scss";

interface UserListProps {
  selectUser: (user: User) => void;
  selectedUserName?: string;
}

const UserList: React.FC<UserListProps> = (props) => {
  const { selectUser, selectedUserName } = props;
  const auth = useAuthData();
  const [userList, setUserList] = React.useState<User[]>([]);

  React.useEffect(() => {
    if (!auth.userName) {
      return;
    }

    const api = getApi();
    let toCancel: NodeJS.Timeout;
    async function getUserList() {
      try {
        const resp = await api.get<User[]>("/user/list");

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

  let content: JSX.Element | null;

  if (!auth.isLogin) {
    content = null;
  } else {
    content = (
      <div className="userlist-list">
        {userList.map((user) => (
          <div
            className={cx("userlist-item", {
              selected: user.name === selectedUserName,
            })}
            key={user.id}
            onClick={() => {
              selectUser(user);
            }}
          >{`${user.name}`}</div>
        ))}
      </div>
    );
  }
  return (
    <div className="userlist-container">
      <h1>All Users</h1>
      {content}
    </div>
  );
};

export default UserList;
