import React from "react"
import { Redirect } from "react-router-dom"
import { getApi } from './api'

const Register = () => {
    const [username, setUsername] = React.useState('')
    const [isResigered, setIsRegister] = React.useState(false)

    const changeUsername = (e) => {
        setUsername(e.target.value)
    }

    const handleLogin = async () => {
        const api = getApi()
        try {

            await api.post('/register', {
                name: username
            })

            setIsRegister(true)

        } catch (e) {
            console.log(e)
        }

    }

    if (isResigered) {
        return <Redirect to='/login' />
    }

    return (
        <div>
            <h1>Register</h1>
            <input value={username} onChange={changeUsername} />
            <button onClick={handleLogin}>Register</button>
        </div>
    )
}

export default Register
