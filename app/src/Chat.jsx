import React from 'react'
import { Redirect } from 'react-router-dom'

import { useAuthData } from './auth'
import UserList from './UserList'
import InfoBar from './InfoBar'

function sendJson(socket, obj) {
    socket.send(JSON.stringify(obj))
}

const Chat = () => {
    const auth = useAuthData()
    const [socket, setSocket] = React.useState()
    const [toUser, setToUser] = React.useState(null)
    const [message, setMessage] = React.useState('')

    React.useEffect(() => {
        if (!auth.userName) {
            return
        }

        let ws

        function start() {
            ws = new WebSocket(`ws://${auth.serverAddr}/ws`)
            ws.onopen = function (evt) {
                console.log("Login ...")

                sendJson(ws, {
                    cmd: "login",
                    body: { userId: auth.userId, userName: auth.userName }
                })

            }

            ws.onclose = () => {
                console.log("Closing ...")
                setTimeout(start, 5000)
            }
            ws.onmessage = (e) => {
                console.log("Received:", e.data)
                const resp = JSON.parse(e.data)
                if (resp?.cmd === 'login' && resp?.body?.ok) {
                    setSocket(ws)
                }
            }

        }

        start()


    }, [auth.userName, auth.userId, auth.serverAddr])


    if (!auth.isLogin) {
        return <Redirect to='/login' />
    }

    if (!socket) {
        return (
            <div>
                <h1>Hello {auth.username} !</h1>
                <p>Connecting to the server...</p>
            </div>
        )

    }

    const canSend = Boolean(toUser) && message !== ''
    const handleSend = () => {
        sendJson(socket, {
            cmd: "send",
            body: {
                userName: toUser.name,
                content: message
            }
        })
    }
    const handleChangeToUser = (user) => {
        setToUser(user)
    }
    const handleChangeMessage = (e) => {
        setMessage(e.target.value)
    }

    return (
        <div>
            <InfoBar
                toUser={toUser.name}
                serverAddr={auth.serverAddr}
            />
            <h1>Hello {auth.username} !</h1>
            <h3>ws: {auth.serverAddr}</h3>
            <div>
                {toUser ? toUser.name : 'No User Selected'}
            </div>
            <div>
                <input onChange={handleChangeMessage} value={message} placeholder='message' />
            </div>
            <div>
                <button onClick={handleSend} disabled={!canSend}>Send</button>
            </div>
            <UserList selectUser={handleChangeToUser} />
        </div>
    )
}

export default Chat