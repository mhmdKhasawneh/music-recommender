import React, { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'


function Login() {
    const [email, setEmail] = useState("")
    const [password, setPassword] = useState("")
    const navigate = useNavigate()

    const loginUser = async () => {
        const response = await fetch("http://localhost:8080/api/user/login", {
            method: "POST",
            mode: "cors",
            body: JSON.stringify({
                "email": email,
                "password": password
            })
        })
        if (response.ok) {
            const data = await response.json()
            localStorage.setItem("Authorization", data["token"])
            navigate('/')
        }
    }

    return (
        <div>
            <input
                type="text"
                placeholder='Email'
                onChange={(e) => setEmail(e.target.value)} />
            <br />
            <input
                type="password"
                placeholder='Password'
                onChange={(e) => setPassword(e.target.value)} />
            <br />
            <button onClick={loginUser}>LOGIN</button>
            <br /><br />
            <p>Not a user yet? <Link to="/signup">Sign up</Link></p>
        </div>
    )
}

export default Login
