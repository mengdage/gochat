import React from 'react'

interface AuthState {
    isLogin: boolean
    userId: string
    userName: string
    serverAddr: string
}

const AuthContext = React.createContext<AuthState | undefined>(undefined)
const AuthDispatchContext = React.createContext<React.Dispatch<AuthActionType> | undefined>(undefined)

const defaultState: AuthState = {
    serverAddr: '',
    userId: '',
    userName: '',
    isLogin: false
}

type AuthActionType =
    | { type: 'login', body: { userId: string, userName: string, serverAddr: string } }

function authReducer(state: AuthState, action: AuthActionType) {
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

interface AuthProviderProps {}
const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
    const [auth, authDispatch] = React.useReducer(authReducer, defaultState)

    return (
        <AuthContext.Provider value={auth}>
            <AuthDispatchContext.Provider value={authDispatch}>
                {children}
            </AuthDispatchContext.Provider>
        </AuthContext.Provider>
    )

}

const useAuthData = (): AuthState => {
    const auth = React.useContext(AuthContext)
    if (auth === undefined) {
        throw Error("useAuthData must be used inside a AuthProvider")
    }

    return auth
}

const useAuthDispatchData = (): React.Dispatch<AuthActionType> => {
    const dispatch = React.useContext(AuthDispatchContext)
    if (dispatch === undefined) {
        throw Error("useAuthDispatchData must be used inside a AuthProvider")
    }

    return dispatch
}

export { AuthProvider, useAuthData, useAuthDispatchData }
