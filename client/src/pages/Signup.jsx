import React, { useState } from 'react'
import { useNavigate } from 'react-router-dom'

function Signup() {
    const [email, setEmail] = useState("")
    const [password, setPassword] = useState("")
    const navigate = useNavigate()

    const signup = async () => {
        const response = await fetch("http://localhost:8080/api/user/signup", {
            mode: "cors",
            method: "POST",
            body: JSON.stringify({
                "email": email,
                "password": password,
            })
        })
        if (response.ok) {
            window.alert("Signup successful ig.. go login")
            navigate('/login')
        }
        else { window.alert("Whoops.. can't sign you up\nStatus " + response.status) }
    }
    return (
        <div>
            <input type="text" placeholder='Email'
                onChange={(e) => setEmail(e.target.value)} />
            <br />
            <input type="password" placeholder='Password'
                onChange={(e) => setPassword(e.target.value)} />
            <br />
            <button onClick={signup}>SIGNUP</button>
        </div>
    )
}

export default Signup