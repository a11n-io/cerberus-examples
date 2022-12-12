import useFetch from "../../hooks/useFetch";
import Input from "../../uikit/Input";
import Btn from "../../uikit/Btn";
import {Link, useNavigate} from "react-router-dom";
import Loader from "../../uikit/Loader";
import {useContext, useState} from "react";
import {encode} from "base-64"
import {AuthContext} from "../../context/AuthContext";

export default function Login() {
    const auth = useContext(AuthContext)
    const navigate = useNavigate()
    const {post, loading} = useFetch("/")
    const [email, setEmail] = useState()
    const [password, setPassword] = useState()

    function handleLogin(e) {
        e.preventDefault()

        const basicAuth = "Basic " + encode(email+":"+password)

        post("auth/login", {
            email: email,
            password: password
        }, {"Authorization": basicAuth})
            .then(r => {
                auth.login(r)
                navigate("/")
            })
            .catch(e => console.error(e))
    }

    function handleEmailChanged(e) {
        setEmail(e.target.value)
    }

    function handlePasswordChanged(e) {
        setPassword(e.target.value)
    }

    if (loading) {
        return <Loader/>
    }

    return <>
        <form onSubmit={handleLogin}>
            <Input required placeholder="Email" onChange={handleEmailChanged}/>
            <Input required placeholder="Password" type="password" onChange={handlePasswordChanged}/>
            <Btn type="submit">Login</Btn>
        </form>
        <Link to="/register">Register</Link>
    </>
}