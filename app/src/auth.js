import React from 'react'

const AuthContext = React.createContext(undefined)
const AuthDispatchContext = React.createContext(undefined)

const defaultState = {isLogin: false}
function authReducer(state, action) {
    switch (action.type) {
        case 'login': {
            const { userId, userName, serverAddr } = action.body
            return { userId, userName, isLogin: true, serverAddr }
        }
        default: {
            return state
        }
    }
}
const AuthProvider = ({ children }) => {
    const [auth, authDispatch] = React.useReducer(authReducer, defaultState)
    return (
        <AuthContext.Provider value={auth}>
            <AuthDispatchContext.Provider value={authDispatch}>
                {children}
            </AuthDispatchContext.Provider>
        </AuthContext.Provider>
    )

}

const useAuthData = () => {
    const auth = React.useContext(AuthContext)
    if (auth === undefined) {
        throw Error("useAuthData must be used inside a AuthProvider")
    }

    return auth
}

const useAuthDispatchData = () => {
    const dispatch = React.useContext(AuthDispatchContext)
    if (dispatch === undefined) {
        throw Error("useAuthDispatchData must be used inside a AuthProvider")
    }

    return dispatch
}

export { AuthProvider, useAuthData, useAuthDispatchData }
