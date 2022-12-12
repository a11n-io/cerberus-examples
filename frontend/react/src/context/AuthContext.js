import {createContext, useContext, useEffect, useState} from "react";
import {useNavigate} from "react-router-dom";
import {CerberusContext} from "@a11n-io/cerberus-reactjs";
import useSessionStorageState from 'use-session-storage-state';

const AuthContext = createContext(null)

function AuthProvider(props) {
    const cerberusCtx = useContext(CerberusContext)

    const [user, setUser] = useSessionStorageState(`acme-user`, {defaultValue: null});

    const logout = () => {
        setUser(null)
        cerberusCtx.setApiTokenPair(null)
    }

    const login = (user) => {
        setUser(user)
        cerberusCtx.setApiTokenPair(user.cerberusTokenPair)
    }

    const value = {
        user: user,
        login: login,
        logout: logout,
    }

    return (
        <AuthContext.Provider value={value}>
            {props.children}
        </AuthContext.Provider>
    )
}

function AuthGuard(props) {
    const auth = useContext(AuthContext)
    const navigate = useNavigate()
    const {redirectTo = "", ...rest} = props

    useEffect(() => {
        if (redirectTo !== "" && !auth.user) {
            navigate(redirectTo)
        }
    }, [auth])

    if (!auth.user) {
        return <></>
    } else {

        return (
            <>
                {
                    rest.children
                }
            </>
        )
    }
}

export {AuthContext, AuthProvider, AuthGuard}